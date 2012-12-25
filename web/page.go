package web

import (
	"bytes"
	"fmt"
	"golanger.com/framework/cookiesession"
	"golanger.com/framework/filesession"
	"golanger.com/framework/i18n"
	"golanger.com/framework/log"
	"golanger.com/framework/session"
	"golanger.com/framework/urlmanage"
	"golanger.com/framework/validate"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Page struct {
	*site
	Config
	Document
	Controller          map[string]interface{}
	DefaultController   interface{}
	NotFoundtController interface{}
	CurrentController   string
	CurrentAction       string
	Template            string
	MAX_FORM_SIZE       int64
	GET                 map[string]string
	POST                map[string]string
	COOKIE              map[string]string
	SESSION             map[string]interface{}
	ONCE_SESSION        interface{}
	COOKIE_SESSION      map[string]interface{}
	LANG                map[string]string
	Session             *session.SessionManager
	FileSession         *filesession.SessionManager
	CookieSession       *cookiesession.SessionManager
	I18n                *i18n.I18nManager
	Validation          validate.Validation
	UrlManage           *urlmanage.UrlManage
	currentPath         string
	currentFileName     string
}

type PageParam struct {
	MaxFormSize       int64
	SessionDir        string
	CookieName        string
	CookieSessionName string
	CookieSessionKey  string
	I18nName          string
	Expires           int
	TimerDuration     string
}

func NewPage(param PageParam) Page {
	if param.MaxFormSize <= 0 {
		param.MaxFormSize = 2 << 20 // 2MB => 2的20次方 乘以 2 =》 2 * 1024 * 1024
	}

	return Page{
		site: &site{
			base: &base{
				header: map[string][]string{},
			},
			templateFunc:  template.FuncMap{},
			templateCache: map[string]templateCache{},
			Version:       strconv.Itoa(time.Now().Year()),
		},
		Controller: map[string]interface{}{},

		Config: NewConfig(),
		Document: Document{
			Css:  map[string]string{},
			Js:   map[string]string{},
			Img:  map[string]string{},
			Func: template.FuncMap{},
		},
		MAX_FORM_SIZE: param.MaxFormSize,
		Session:       session.New(param.CookieName, param.Expires, param.TimerDuration),
		FileSession:   filesession.New(param.CookieName, param.Expires, param.SessionDir, param.TimerDuration),
		CookieSession: cookiesession.New(param.CookieSessionName, param.CookieSessionKey),
		I18n:          i18n.New(param.I18nName),
		UrlManage:     urlmanage.New(),
	}
}

func (p *Page) Init(w http.ResponseWriter, r *http.Request) {
	p.site.base.mutex.Lock()
	p.site.Init(w, r)
	p.site.base.mutex.Unlock()

	p.site.base.rmutex.RLock()
	p.GET = p.site.base.getHttpGet(r)
	p.POST = p.site.base.getHttpPost(r, p.MAX_FORM_SIZE)
	p.COOKIE = p.site.base.getHttpCookie(r)
	if p.site.supportSession {
		switch p.Config.SessionType {
		case "file":
			p.SESSION = p.FileSession.Get(w, r)
		case "memory":
			p.SESSION = p.Session.Get(w, r)
		default:
			p.SESSION = p.Session.Get(w, r)
		}

		var ok bool
		p.ONCE_SESSION, ok = p.SESSION["__ONCE"]
		if ok {
			log.Debug("<Page.Init> ", `p.SESSION["__ONCE"] to set:`, p.ONCE_SESSION)
			delete(p.SESSION, "__ONCE")
			log.Debug("<Page.Init> ", `p.SESSION["__ONCE"] to delete:`, p.SESSION["__ONCE"])
		}
	}

	if p.site.supportCookieSession {
		p.COOKIE_SESSION = p.CookieSession.Get(r)
	}

	if p.site.supportI18n {
		p.LANG = func() map[string]string {
			l := strings.TrimSpace(r.Header.Get("Accept-Language"))
			if i := strings.Index(l, ","); i != -1 {
				l = l[:i]
			}

			return p.I18n.Lang(l)
		}()
	}

	if p.site.base.header != nil || len(p.site.base.header) > 0 {
		log.Debug("<Page.Init> ", "p.site.base.header:", p.site.base.header)
		for t, s := range p.site.base.header {
			for _, v := range s {
				w.Header().Add(t, v)
			}
		}
	}
	p.site.base.rmutex.RUnlock()
}

func (p *Page) SetDefaultController(i interface{}) *Page {
	p.DefaultController = i

	return p
}

func (p *Page) SetNotFoundController(i interface{}) *Page {
	p.NotFoundtController = i

	return p
}

func (p *Page) RegisterController(relUrlPath string, i interface{}) *Page {
	if _, ok := p.Controller[relUrlPath]; !ok {
		p.Controller[relUrlPath] = i
	}

	return p
}

func (p *Page) UpdateController(oldUrlPath, relUrlPath string, i interface{}) *Page {
	delete(p.Controller, oldUrlPath)
	p.Controller[relUrlPath] = i

	return p
}

func (p *Page) GetController(urlPath string) interface{} {
	var relUrlPath string
	if strings.HasPrefix(urlPath, p.site.Root) {
		relUrlPath = urlPath[len(p.site.Root):]
	} else {
		relUrlPath = urlPath
	}

	i, ok := p.Controller[relUrlPath]
	if !ok {
		i = p.NotFoundtController
	}

	return i
}

func (p *Page) Load(configPath string) {
	p.Config.Load(configPath)
	p.reset(false)
}

func (p *Page) setGlobalTpl(globalTplModTime int64, reset bool) {
	if globalTpls, err := filepath.Glob(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Config.TemplateGlobalDirectory + p.Config.TemplateGlobalFile); err == nil && globalTpls != nil && len(globalTpls) > 0 {
		var buf bytes.Buffer
		for _, tpl := range globalTpls {
			if b, err := ioutil.ReadFile(tpl); err == nil {
				buf.Write(b)
			} else {
				log.Error("<Page.setGlobalTpl> ", err)
			}
		}

		sBuf := strings.TrimSpace(buf.String())
		if sBuf != "" {
			p.site.SetTemplateCacheObject("globalTpl", buf.String(), globalTplModTime)
			if reset {
				newT := template.New("globalTpl").Funcs(p.site.templateFunc)
				if t, err := newT.Parse(sBuf); t != nil {
					p.site.globalTemplate = t
				} else {
					log.Error("<Page.setGlobalTpl> ", err)
				}
			} else {
				if t, err := p.site.globalTemplate.Parse(sBuf); t != nil {
					p.site.globalTemplate = t
				} else {
					log.Error("<Page.setGlobalTpl> ", err)
				}
			}
		}
	}
}

func (p *Page) reset(update bool) {
	p.setLog(p.Config.SupportLog, p.Config.LogWriteTo, p.Config.LogLevel)
	p.setUrlManage(p.Config.SupportUrlManage, p.Config.UrlManageRule)

	if update {
		if p.site.supportSession != p.Config.SupportSession {
			p.site.supportSession = p.Config.SupportSession
		}

		if p.site.supportCookieSession != p.Config.SupportCookieSession {
			p.site.supportCookieSession = p.Config.SupportCookieSession
		}

		if p.site.supportI18n != p.Config.SupportI18n {
			p.site.supportI18n = p.Config.SupportI18n
		}

		if p.Document.Theme != p.Config.Theme {
			p.Document.Theme = p.Config.Theme
		}

		if p.Document.Static != p.Config.SiteRoot+p.Config.StaticDirectory {
			p.Document.Static = p.Config.SiteRoot + p.Config.StaticDirectory
		}

		if p.site.Root == p.Config.SiteRoot {
			return
		} else {
			log.Info("<Page.reset> ", "p.Config.SiteRoot be changed:", p.Config.SiteRoot)
			p.SetDefaultController(p.GetController(p.Config.IndexDirectory))
			p.UpdateController(p.site.Root, p.Config.SiteRoot, p.DefaultController)
			p.site.Root = p.Config.SiteRoot
		}
	} else {
		p.site.supportSession = p.Config.SupportSession
		p.site.supportCookieSession = p.Config.SupportCookieSession
		p.site.supportI18n = p.Config.SupportI18n
		p.Document.Theme = p.Config.Theme
		p.site.Root = p.Config.SiteRoot
		p.Document.Static = p.site.Root + p.Config.StaticDirectory
		p.SetDefaultController(p.GetController(p.Config.IndexDirectory))
		p.RegisterController(p.site.Root, p.DefaultController)
		p.site.globalTemplate = template.New("globalTpl").Funcs(p.site.templateFunc)
	}

	if globalCssFi, err := os.Stat(p.Config.StaticCssDirectory + "/global/"); err == nil && globalCssFi.IsDir() {
		DcssPath := p.Config.StaticCssDirectory + "global/"
		p.Document.Css["global"] = p.site.Root + DcssPath[len(p.Config.AssetsDirectory):]
		log.Debug("<Page.reset> ", `p.Document.Css["global"]:`, p.Document.Css["global"])

		if _, err := os.Stat(DcssPath + "global.css"); err == nil {
			p.Document.GlobalCssFile = p.Document.Css["global"] + "global.css"
		}
	}

	if globalJsFi, err := os.Stat(p.Config.StaticJsDirectory + "/global/"); err == nil && globalJsFi.IsDir() {
		DjsPath := p.Config.StaticJsDirectory + "global/"
		p.Document.Js["global"] = p.site.Root + DjsPath[len(p.Config.AssetsDirectory):]
		log.Debug("<Page.reset> ", `p.Document.Js["global"]:`, p.Document.GlobalCssFile)
		if _, err := os.Stat(DjsPath + "global.js"); err == nil {
			p.Document.GlobalJsFile = p.Document.Js["global"] + "global.js"
		}
	}

	if globalImgFi, err := os.Stat(p.Config.StaticImgDirectory + "/global/"); err == nil && globalImgFi.IsDir() {
		DimgPath := p.Config.StaticImgDirectory + "global/"
		p.Document.Img["global"] = p.site.Root + DimgPath[len(p.Config.AssetsDirectory):]
		log.Debug("<Page.reset> ", `p.Document.Img["global"]:`, p.Document.Img["global"])
	}

	if globalTplFi, err := os.Stat(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Config.TemplateGlobalDirectory); err != nil {
		log.Error("<Page.reset> ", err)
	} else {
		p.setGlobalTpl(globalTplFi.ModTime().Unix(), false)
	}
}

func (p *Page) setUrlManage(manage bool, rules []string) {
	if !manage {
		p.UrlManage.Stop()
	} else {
		if len(rules) == 0 {
			p.UrlManage.Stop()
		} else {
			p.UrlManage.Start()
			ruleContent := strings.Join(rules, "\n")
			p.UrlManage.LoadRule(ruleContent, true)
		}
	}
}

func (p *Page) setLog(hasLog bool, writeTo, level string) {
	if !hasLog {
		log.SetLevel(log.LEVEL_DISABLE)
	} else {
		if level == "" {
			log.SetLevel(log.LEVEL_DEFAULT)
		} else {
			var lv int
			lvs := strings.Split(level, ",")
			for _, v := range lvs {
				if l, lok := log.LevelText[strings.ToLower(strings.TrimSpace(v))]; lok {
					lv |= l
				}
			}

			log.SetLevel(lv)
		}

		switch writeTo {
		default:
			//包括console
		}
	}
}

func (p *Page) setCurrentInfo(path string) {
	urlPath, fileName := filepath.Split(path)
	if urlPath == p.site.Root {
		urlPath = p.site.Root + p.Config.IndexDirectory
	}

	if fileName == "" {
		fileName = p.Config.IndexPage
	}

	p.currentPath = urlPath
	log.Debug("<Page.setCurrentInfo> ", "p.currentPath:", p.currentPath)
	p.currentFileName = fileName
	log.Debug("<Page.setCurrentInfo> ", "p.currentFileName:", p.currentFileName)
	p.CurrentController = urlPath[len(p.site.Root):]
	log.Debug("<Page.setCurrentInfo> ", "p.CurrentController:", p.CurrentController)
	p.CurrentAction = strings.Replace(strings.Title(strings.Replace(p.currentFileName[:len(p.currentFileName)-len(filepath.Ext(p.currentFileName))], "_", " ", -1)), " ", "", -1)
	log.Debug("<Page.setCurrentInfo> ", "p.CurrentAction:", p.CurrentAction)
}

func (p *Page) callMethod(tpc reflect.Type, vpc reflect.Value, action string, rvr reflect.Value, rvw reflect.Value) []reflect.Value {
	arv := []reflect.Value{}
	if rm, ok := tpc.MethodByName(action); ok {
		mt := rm.Type
		log.Debug("<Page.callMethod> ", mt.String()+".NumIn():", mt.NumIn())
		switch mt.NumIn() {
		case 2:
			if mt.In(1) == rvr.Type() {
				arv = vpc.MethodByName(action).Call([]reflect.Value{rvr})
			} else {
				arv = vpc.MethodByName(action).Call([]reflect.Value{rvw})
			}
		case 3:
			arv = vpc.MethodByName(action).Call([]reflect.Value{rvw, rvr})
		default:
			arv = vpc.MethodByName(action).Call([]reflect.Value{})
		}
	}

	return arv
}

func (p *Page) filterMethod(methodName string) bool {
	b := false
	b = !strings.HasPrefix(methodName, "Before_") && !strings.HasPrefix(methodName, "After_") && !strings.HasPrefix(methodName, "Filter_")
	return b
}

func (p *Page) filterDoMethod(tpc reflect.Type, vpc reflect.Value, action string, rvr reflect.Value, rvw reflect.Value) {
	log.Debug("<Page.filterDoMethod> ", "action:", action)
	if aarv := p.callMethod(tpc, vpc, action, rvr, rvw); len(aarv) == 1 {
		if ma, tok := aarv[0].Interface().([]map[string]string); tok {
			log.Debug("<Page.filterDoMethod> ", "FilterMap:", ma)
			for _, filter := range ma {
				doFilter := false
				doMethod := false
				doParam := false
				method, mok := filter["_FILTER"]

				if mok {
					filterType, found := filter[p.CurrentAction]
					if !found && filter["_ALL"] == "allow" {
						doFilter = true
					} else if found && filterType == "allow" {
						doFilter = true
					}
				}

				if doFilter {
					if rm, ok := filter["_METHOD"]; !ok {
						doMethod = true
					} else {
						reqMethods := strings.Split(rm, ",")
						for _, v := range reqMethods {
							r := rvr.Interface().(*http.Request)
							if r.Method == strings.ToUpper(v) {
								doMethod = true
								break
							}
						}
					}

					if rmp, ok := filter["_PARAM"]; !ok {
						doParam = true
					} else {
						reqParams := strings.Split(rmp, ",")
						for _, v := range reqParams {
							_, vgok := p.GET[v]
							_, vpok := p.POST[v]
							if vgok || vpok {
								doParam = true
								break
							}
						}
					}

					if doMethod && doParam {
						p.callMethod(tpc, vpc, "Filter_"+method, rvr, rvw)
					}
				}
			}
		}
	}
}

func (p *Page) routeController(i interface{}, w http.ResponseWriter, r *http.Request) {
	pageOriController := p.GetController(p.currentPath)
	rv := reflect.ValueOf(pageOriController)
	rvw, rvr := reflect.ValueOf(w), reflect.ValueOf(r)
	rt := rv.Type()
	vpc := reflect.New(rt)

	iv := reflect.ValueOf(i).Elem()
	vapc := vpc.Elem().FieldByName("Application")
	vapc.Set(iv)

	tapc := vapc.Type()
	if _, found := tapc.FieldByName("RW"); found {
		vapc.FieldByName("RW").Set(rvw)
	}

	if _, found := tapc.FieldByName("R"); found {
		vapc.FieldByName("R").Set(rvr)
	}

	vppc := vapc.FieldByName("Page")
	ppc := vppc.Addr().Interface().(*Page)
	log.Debug("<Page.routeController> ", "ppc.Config.LoadStaticHtmlWithLogic:", ppc.Config.LoadStaticHtmlWithLogic)
	log.Debug("<Page.routeController> ", "ppc.Config.AutoGenerateHtml:", ppc.Config.AutoGenerateHtml)

	if !ppc.Config.LoadStaticHtmlWithLogic && ppc.Config.AutoGenerateHtml {
		tplFi, tplErr := os.Stat(ppc.Config.TemplateDirectory + ppc.Config.ThemeDirectory + ppc.Template)
		if tplErr != nil {
			log.Error("<Page.routeController> ", tplErr)
		} else {
			ppc.site.base.rmutex.RLock()
			tmplCache := ppc.site.GetTemplateCache(ppc.Template)
			ppc.site.base.rmutex.RUnlock()
			if tplFi.ModTime().Unix() > tmplCache.ModTime {
				goto DO_ROUTER
			}
		}

		htmlFile := ""
		ppc.site.base.rmutex.RLock()
		assetsHtmlDir := ppc.Config.AssetsDirectory + ppc.Config.HtmlDirectory
		if strings.HasPrefix(ppc.Template, ppc.Config.IndexDirectory) {
			htmlFile = assetsHtmlDir + strings.Replace(ppc.Template, ppc.Config.IndexDirectory, "", 1)
		} else {
			htmlFile = assetsHtmlDir + ppc.Template
		}

		if r.URL.RawQuery != "" {
			htmlFile += "?" + r.URL.RawQuery
		}

		ppc.site.base.rmutex.RUnlock()
		if htmlFi, htmlErr := os.Stat(htmlFile); htmlErr == nil {
			if ppc.checkHtmlDoWrite(tplFi, htmlFi, htmlErr) {
				goto DO_ROUTER
			}

			htmlContent, err := ioutil.ReadFile(htmlFile)
			if err == nil {
				w.Write(htmlContent)
				return
			} else {
				goto DO_ROUTER
			}
		} else {
			goto DO_ROUTER
		}

	}

DO_ROUTER:
	tpc := vpc.Type()
	if ppc.CurrentAction != "Init" {
		ppc.callMethod(tpc, vpc, "Init", rvr, rvw)
	}

	if _, ok := tpc.MethodByName(ppc.CurrentAction); ok && ppc.filterMethod(ppc.CurrentAction) {
		ppc.filterDoMethod(tpc, vpc, "Before_", rvr, rvw)
		ppc.callMethod(tpc, vpc, "Before_"+ppc.CurrentAction, rvr, rvw)

		if ppc.Document.Close == false {
			ppc.callMethod(tpc, vpc, ppc.CurrentAction, rvr, rvw)
		}

		ppc.callMethod(tpc, vpc, "After_"+ppc.CurrentAction, rvr, rvw)
		ppc.filterDoMethod(tpc, vpc, "After_", rvr, rvw)
	} else {
		if !strings.Contains(tpc.String(), "Page404") {
			notFountRV := reflect.ValueOf(ppc.NotFoundtController)
			notFountRT := notFountRV.Type()
			vnpc := reflect.New(notFountRT)
			vanpc := vnpc.Elem().FieldByName("Application")
			vanpc.Set(vapc)
			vpnpc := vanpc.FieldByName("Page")
			vpnpc.Set(vppc)

			ppc = vpnpc.Addr().Interface().(*Page)
			tnpc := vnpc.Type()
			ppc.callMethod(tnpc, vnpc, "Init", rvr, rvw)
		}
	}

	if ppc.site.supportSession {
		switch ppc.Config.SessionType {
		case "file":
			ppc.FileSession.Set(ppc.SESSION, w, r)
		case "memory":
			ppc.Session.Set(w, r)
		default:
			ppc.Session.Set(w, r)
		}
	}

	if ppc.site.supportCookieSession {
		ppc.CookieSession.Set(ppc.COOKIE_SESSION, w, r)
	}

	if ppc.Config.SupportTemplate {
		ppc.setStaticDocument()
		ppc.routeTemplate(w, r)
	}
}

func (p *Page) setStaticDocument() {
	fileNameNoExt := p.currentFileName[:len(p.currentFileName)-len(filepath.Ext(p.currentFileName))]
	p.site.base.rmutex.RLock()
	p.site.base.rmutex.RUnlock()

	cssFi, cssErr := os.Stat(p.Config.StaticCssDirectory + p.currentPath)
	jsFi, jsErr := os.Stat(p.Config.StaticJsDirectory + p.currentPath)
	imgFi, imgErr := os.Stat(p.Config.StaticImgDirectory + p.currentPath)

	if cssErr == nil && cssFi.IsDir() {
		cssPath := strings.Trim(p.currentPath, "/")
		DcssPath := p.Config.StaticCssDirectory + cssPath + "/"
		p.Document.Css[cssPath] = p.site.Root + DcssPath[len(p.Config.AssetsDirectory):]
		log.Debug("<Page.setStaticDocument> ", "p.Document.Css["+cssPath+"]:", p.Document.Css[cssPath])

		_, errgcss := os.Stat(DcssPath + "global.css")
		_, errcss := os.Stat(DcssPath + fileNameNoExt + ".css")

		if errgcss == nil {
			p.Document.GlobalIndexCssFile = p.Document.Css[cssPath] + "global.css"
			log.Debug("<Page.setStaticDocument> ", "p.Document.GlobalIndexCssFile:", p.Document.GlobalIndexCssFile)
		}

		if errcss == nil {
			p.Document.IndexCssFile = p.Document.Css[cssPath] + fileNameNoExt + ".css"
			log.Debug("<Page.setStaticDocument> ", "p.Document.IndexCssFile:", p.Document.IndexCssFile)
		}
	}

	if jsErr == nil && jsFi.IsDir() {
		jsPath := strings.Trim(p.currentPath, "/")
		DjsPath := p.Config.StaticJsDirectory + jsPath + "/"
		p.Document.Js[jsPath] = p.site.Root + DjsPath[len(p.Config.AssetsDirectory):]
		log.Debug("<Page.setStaticDocument> ", "p.Document.Js["+jsPath+"]:", p.Document.Js[jsPath])

		_, errgjs := os.Stat(DjsPath + "global.js")
		_, errjs := os.Stat(DjsPath + fileNameNoExt + ".js")

		if errgjs == nil {
			p.Document.GlobalIndexJsFile = p.Document.Js[jsPath] + "global.js"
			log.Debug("<Page.setStaticDocument> ", "p.Document.GlobalIndexJsFile:", p.Document.GlobalIndexJsFile)
		}

		if errjs == nil {
			p.Document.IndexJsFile = p.Document.Js[jsPath] + fileNameNoExt + ".js"
			log.Debug("<Page.setStaticDocument> ", "p.Document.IndexJsFile:", p.Document.IndexJsFile)
		}
	}

	if imgErr == nil && imgFi.IsDir() {
		imgPath := strings.Trim(p.currentPath, "/")
		DimgPath := p.Config.StaticImgDirectory + imgPath + "/"
		p.Document.Img[imgPath] = p.site.Root + DimgPath[len(p.Config.AssetsDirectory):]
		log.Debug("<Page.setStaticDocument> ", "p.Document.Img["+imgPath+"]:", p.Document.Img[imgPath])
	}
}

func (p *Page) routeTemplate(w http.ResponseWriter, r *http.Request) {
	if p.Config.AutoGenerateHtml {
		p.Document.GenerateHtml = true
	}

	if p.Document.Close == false && p.Document.Hide == false {
		if globalTplFi, err := os.Stat(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Config.TemplateGlobalDirectory); err != nil {
			log.Error("<Page.routeTemplate> ", err)
		} else {
			if globalTplCache := p.site.GetTemplateCache("globalTpl"); globalTplCache.ModTime > 0 {
				if globalTplFi.ModTime().Unix() > globalTplCache.ModTime {
					p.setGlobalTpl(globalTplFi.ModTime().Unix(), true)
				}
			}
		}

		if tplFi, err := os.Stat(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Template); err != nil {
			p.site.base.rmutex.RLock()
			tmplCache := p.site.GetTemplateCache(p.Template)
			if tmplCache.ModTime > 0 {
				p.site.base.mutex.Lock()
				p.DelTemplateCache(p.Template)
				p.site.base.mutex.Unlock()
				log.Info("<Page.routeTemplate> ", "Delete Template Cache:", p.Template)
			}
			p.site.base.rmutex.RUnlock()
		} else {
			p.site.base.rmutex.RLock()
			tmplCache := p.site.GetTemplateCache(p.Template)

			if tplFi.ModTime().Unix() > tmplCache.ModTime {
				p.site.base.mutex.Lock()
				p.SetTemplateCache(p.Template, p.Config.TemplateDirectory+p.Config.ThemeDirectory+p.Template)
				p.site.base.mutex.Unlock()

				tmplCache = p.site.GetTemplateCache(p.Template)
			}

			globalTemplate, _ := p.site.globalTemplate.Clone()
			pageTemplate, err := globalTemplate.New(filepath.Base(p.Template)).Parse(tmplCache.Content)
			p.site.base.rmutex.RUnlock()

			if err != nil {
				log.Error("<Page.routeTemplate> ", err)
				w.Write([]byte(fmt.Sprint(err)))
			} else {
				templateVar := map[string]interface{}{
					"G":        p.GET,
					"P":        p.POST,
					"S":        p.SESSION,
					"O_S":      p.ONCE_SESSION,
					"C":        p.COOKIE,
					"CS":       p.COOKIE_SESSION,
					"D":        p.Document,
					"L":        p.LANG,
					"Config":   p.Config.M,
					"Template": p.Template,
				}

				p.site.base.rmutex.RLock()
				templateVar["Siteroot"] = p.site.Root
				templateVar["Version"] = p.site.Version
				p.site.base.rmutex.RUnlock()

				if !p.Document.GenerateHtml {
					err := pageTemplate.Execute(w, templateVar)
					if err != nil {
						log.Error("<Page.routeTemplate> ", err)
						w.Write([]byte(fmt.Sprint(err)))
					}
				} else {
					htmlFile := ""
					p.site.base.rmutex.RLock()
					assetsHtmlDir := p.Config.AssetsDirectory + p.Config.HtmlDirectory
					if strings.HasPrefix(p.Template, p.Config.IndexDirectory) {
						htmlFile = assetsHtmlDir + strings.Replace(p.Template, p.Config.IndexDirectory, "", 1)
					} else {
						htmlFile = assetsHtmlDir + p.Template
					}

					if r.URL.RawQuery != "" {
						htmlFile += "?" + r.URL.RawQuery
					}

					p.site.base.rmutex.RUnlock()
					htmlDir := filepath.Dir(htmlFile)
					if htmlDirFi, err := os.Stat(htmlDir); err != nil || !htmlDirFi.IsDir() {
						os.MkdirAll(htmlDir, 0777)
						log.Info("<Page.routeTemplate> ", "MkdirAll:", htmlDir)
					}

					htmlFi, htmlErr := os.Stat(htmlFile)
					if p.checkHtmlDoWrite(tplFi, htmlFi, htmlErr) {
						if file, err := os.OpenFile(htmlFile, os.O_CREATE|os.O_WRONLY, 0777); err != nil {
							log.Error("<Page.routeTemplate> ", err)
						} else {
							if p.Config.AutoJumpToHtml || p.Config.ChangeSiteRoot {
								templateVar["Siteroot"] = p.site.Root + p.Config.HtmlDirectory
							}

							pageTemplate.Execute(file, templateVar)
						}
					}

					if p.Config.AutoJumpToHtml {
						p.site.base.rmutex.RLock()
						http.Redirect(w, r, p.site.Root+htmlFile[2:], http.StatusFound)
						p.site.base.rmutex.RUnlock()
					} else if p.Config.AutoLoadStaticHtml && htmlFi != nil && htmlErr == nil {
						htmlContent, err := ioutil.ReadFile(htmlFile)
						if err == nil {
							w.Write(htmlContent)
						} else {
							log.Error("<Page.routeTemplate> ", err)
						}
					} else {
						err := pageTemplate.Execute(w, templateVar)
						if err != nil {
							log.Error("<Page.routeTemplate> ", err)
						}
					}
				}
			}
		}
	}
}

func (p *Page) checkHtmlDoWrite(tplFi, htmlFi os.FileInfo, htmlErr error) bool {
	var doWrite bool
	log.Debug("<Page.checkHtmlDoWrite> ", "p.Config.AutoGenerateHtmlCycleTime:", p.Config.AutoGenerateHtmlCycleTime)
	if p.Config.AutoGenerateHtmlCycleTime <= 0 {
		doWrite = true
	} else {
		if htmlErr != nil {
			doWrite = true
		} else {
			log.Debug("<Page.checkHtmlDoWrite> ", "tplFi.ModTime().Unix():", tplFi.ModTime().Unix())
			log.Debug("<Page.checkHtmlDoWrite> ", "htmlFi.ModTime().Unix():", htmlFi.ModTime().Unix())

			switch {
			case tplFi.ModTime().Unix() >= htmlFi.ModTime().Unix():
				doWrite = true
			case time.Now().Unix()-htmlFi.ModTime().Unix() >= p.Config.AutoGenerateHtmlCycleTime:
				doWrite = true
			default:
				globalTplCache := p.site.GetTemplateCache("globalTpl")
				log.Debug("<Page.checkHtmlDoWrite> ", `globalTplCache.ModTime:`, globalTplCache.ModTime)
				if globalTplCache.ModTime > 0 && globalTplCache.ModTime >= htmlFi.ModTime().Unix() {
					doWrite = true
				}
			}
		}
	}

	return doWrite
}

func (p *Page) handleRootStatic(files string) {
	aFile := strings.Split(files, ",")
	for _, file := range aFile {
		http.HandleFunc(p.site.Root+file, func(w http.ResponseWriter, r *http.Request) {
			staticPath := p.Config.AssetsDirectory + file
			log.Debug("<Page.handleRootStatic> ", "staticPath:", staticPath)
			http.ServeFile(w, r, staticPath)
		})
	}
}

func (p *Page) handleStatic() {
	StaticHtmlDir := p.Config.SiteRoot + p.Config.HtmlDirectory
	http.HandleFunc(StaticHtmlDir, func(w http.ResponseWriter, r *http.Request) {
		if p.UrlManage.Manage() {
			newUrl := p.UrlManage.ReWrite(w, r)
			if newUrl == "redirect" {
				p.site.base.mutex.Unlock()
				return
			} else {
				r.URL, _ = url.Parse(newUrl)
			}
		}
		
		staticPath := p.Config.AssetsDirectory + p.Config.HtmlDirectory + r.URL.Path[len(StaticHtmlDir):]
		if r.URL.RawQuery != "" {
			staticPath += "?" + r.URL.RawQuery
		}

		log.Debug("<Page.handleStatic> ", "staticPath:", staticPath)
		http.ServeFile(w, r, staticPath)
	})

	http.HandleFunc(p.Document.Static, func(w http.ResponseWriter, r *http.Request) {
		staticPath := p.Config.AssetsDirectory + p.Config.StaticDirectory + r.URL.Path[len(p.Document.Static):]
		log.Debug("<Page.handleStatic> ", "staticPath:", staticPath)
		http.ServeFile(w, r, staticPath)
	})
}

func (p *Page) handleRoute(i interface{}) {
	http.HandleFunc(p.site.Root, func(w http.ResponseWriter, r *http.Request) {
		p.site.base.mutex.Lock()
		if p.Config.Reload() {
			p.reset(true)
		}

		if p.UrlManage.Manage() {
			newUrl := p.UrlManage.ReWrite(w, r)
			if newUrl == "redirect" {
				p.site.base.mutex.Unlock()
				return
			} else {
				r.URL, _ = url.Parse(newUrl)
			}
		}

		if p.site.supportI18n {
			if err := p.I18n.Setup(p.Config.DefaultLocalePath, p.Config.DefaultLanguage); err != nil {
				log.Panic("<Page.handleRoute> ", "I18n(Setup):", err)
			}
		}

		p.setCurrentInfo(r.URL.Path)
		p.Template = p.CurrentController + p.currentFileName
		log.Debug("<Page.handleRoute> ", "p.Template:", p.Template)
		p.site.base.mutex.Unlock()

		p.routeController(i, w, r)
	})
}

func (p *Page) ListenAndServe(addr string, i interface{}) {
	if p.Config.SupportStatic {
		p.handleRootStatic(p.Config.RootStaticFiles)
		p.handleStatic()
	}

	p.handleRoute(i)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("<Page.ListenAndServe> ", err)
	}
}
