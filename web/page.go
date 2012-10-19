package web

import (
	"fmt"
	"golanger.com/framework/session"
	"log"
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
	Session             *session.SessionManager
	currentPath         string
	currentFileName     string
}

type PageParam struct {
	MaxFormSize   int64
	CookieName    string
	Expires       int
	TimerDuration string
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
		p.SESSION = p.Session.Get(w, r)
	}

	if p.site.base.header != nil || len(p.site.base.header) > 0 {
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

func (p *Page) reset(update bool) {
	if update {
		if p.site.supportSession != p.Config.SupportSession {
			p.site.supportSession = p.Config.SupportSession
		}

		if p.Document.Theme != p.Config.Theme {
			p.Document.Theme = p.Config.Theme
		}

		if p.Document.Static != p.Config.SiteRoot+p.Config.StaticDirectory[2:] {
			p.Document.Static = p.Config.SiteRoot + p.Config.StaticDirectory[2:]
		}

		if p.site.Root == p.Config.SiteRoot {
			return
		} else {
			p.SetDefaultController(p.GetController(p.Config.IndexDirectory))
			p.UpdateController(p.site.Root, p.Config.SiteRoot, p.DefaultController)
			p.site.Root = p.Config.SiteRoot
		}
	} else {
		p.site.supportSession = p.Config.SupportSession
		p.Document.Theme = p.Config.Theme
		p.site.Root = p.Config.SiteRoot
		p.Document.Static = p.site.Root + p.Config.StaticDirectory[2:]
		p.SetDefaultController(p.GetController(p.Config.IndexDirectory))
		p.RegisterController(p.site.Root, p.DefaultController)
		p.site.globalTemplate = template.New("globalTpl").Funcs(p.site.templateFunc)
	}

	siteRootRightTrim := p.site.Root[:len(p.site.Root)-1]

	if globalCssFi, err := os.Stat(p.Config.StaticCssDirectory + "/global/"); err == nil && globalCssFi.IsDir() {
		DcssPath := p.Config.StaticCssDirectory + "global/"
		p.Document.Css["global"] = siteRootRightTrim + DcssPath[1:]
		if _, err := os.Stat(DcssPath + "global.css"); err == nil {
			p.Document.GlobalCssFile = p.Document.Css["global"] + "global.css"
		}
	}

	if globalJsFi, err := os.Stat(p.Config.StaticJsDirectory + "/global/"); err == nil && globalJsFi.IsDir() {
		DjsPath := p.Config.StaticJsDirectory + "global/"
		p.Document.Js["global"] = siteRootRightTrim + DjsPath[1:]
		if _, err := os.Stat(DjsPath + "global.js"); err == nil {
			p.Document.GlobalJsFile = p.Document.Js["global"] + "global.js"
		}
	}

	if globalImgFi, err := os.Stat(p.Config.StaticImgDirectory + "/global/"); err == nil && globalImgFi.IsDir() {
		DimgPath := p.Config.StaticImgDirectory + "global/"
		p.Document.Img["global"] = siteRootRightTrim + DimgPath[1:]
	}

	if t, _ := p.site.globalTemplate.ParseGlob(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Config.TemplateGlobalDirectory + p.Config.TemplateGlobalFile); t != nil {
		p.site.globalTemplate = t
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
	p.currentFileName = fileName
	p.CurrentController = urlPath[len(p.site.Root):]
	p.CurrentAction = strings.Replace(strings.Title(strings.Replace(p.currentFileName[:len(p.currentFileName)-len(filepath.Ext(p.currentFileName))], "_", " ", -1)), " ", "", -1)
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
	tpc := vpc.Type()

	if ppc.CurrentAction != "Init" {
		if rm, ok := tpc.MethodByName("Init"); ok {
			mt := rm.Type
			switch mt.NumIn() {
			case 2:
				if mt.In(1) == rvr.Type() {
					vpc.MethodByName("Init").Call([]reflect.Value{rvr})

				} else {
					vpc.MethodByName("Init").Call([]reflect.Value{rvw})
				}
			case 3:
				vpc.MethodByName("Init").Call([]reflect.Value{rvw, rvr})
			default:
				vpc.MethodByName("Init").Call([]reflect.Value{})
			}
		}
	}


	if rm, ok := tpc.MethodByName(ppc.CurrentAction); ok {
		ctl := p.CurrentController[:len(p.CurrentController)-1]
		ca := p.ControllerAction
		bfok,afok := true,true
		if ca["Before"+ctl][ppc.CurrentAction] == "deny" || 
			(ca["Before"+ctl]["All_"] == "deny" && ca["Before"+ctl][ppc.CurrentAction] != "allow"){
			 bfok = false
		}
		if ca["After"+ctl][ppc.CurrentAction] == "deny" || 
			(ca["After"+ctl]["All_"] == "deny" && ca["After"+ctl][ppc.CurrentAction] != "allow"){
			 afok = false
		}
		
		if bam, ok := tpc.MethodByName("BeforeAction"); ok && bfok{
			excuteActionFunc(bam ,vpc,"BeforeAction", rvr , rvw ,ppc)
		}
		if btam, ok := tpc.MethodByName("Before"+ppc.CurrentAction); ok {
			excuteActionFunc(btam ,vpc,"Before"+ppc.CurrentAction, rvr , rvw ,ppc)
		}
		
		excuteActionFunc(rm ,vpc,ppc.CurrentAction , rvr , rvw ,ppc)
		
		if aam, ok := tpc.MethodByName("AfterAction"); ok && afok{
			excuteActionFunc(aam ,vpc,"AfterAction", rvr , rvw ,ppc)
		}
		if atam, ok := tpc.MethodByName("After"+ppc.CurrentAction); ok {
			excuteActionFunc(atam ,vpc,"After"+ppc.CurrentAction, rvr , rvw ,ppc)
		}
		
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

			if rm, ok := tnpc.MethodByName("Init"); ok {
				mt := rm.Type
				switch mt.NumIn() {
				case 2:
					if mt.In(1) == rvr.Type() {
						vnpc.MethodByName("Init").Call([]reflect.Value{rvr})
					} else {
						vnpc.MethodByName("Init").Call([]reflect.Value{rvw})
					}
				case 3:
					vnpc.MethodByName("Init").Call([]reflect.Value{rvw, rvr})
				default:
					vnpc.MethodByName("Init").Call([]reflect.Value{})
				}
			}
		}
	}

	if ppc.Config.SupportTemplate {
		ppc.setStaticDocument()
		ppc.routeTemplate(w, r)
	}

}

func (p *Page) setStaticDocument() {
	fileNameNoExt := p.currentFileName[:len(p.currentFileName)-len(filepath.Ext(p.currentFileName))]
	p.site.base.rmutex.RLock()
	siteRootRightTrim := p.site.Root[:len(p.site.Root)-1]
	p.site.base.rmutex.RUnlock()

	cssFi, cssErr := os.Stat(p.Config.StaticCssDirectory + p.currentPath)
	jsFi, jsErr := os.Stat(p.Config.StaticJsDirectory + p.currentPath)
	imgFi, imgErr := os.Stat(p.Config.StaticImgDirectory + p.currentPath)

	if cssErr == nil && cssFi.IsDir() {
		cssPath := strings.Trim(p.currentPath, "/")
		DcssPath := p.Config.StaticCssDirectory + cssPath + "/"
		p.Document.Css[cssPath] = siteRootRightTrim + DcssPath[1:]

		_, errgcss := os.Stat(DcssPath + "global.css")
		_, errcss := os.Stat(DcssPath + fileNameNoExt + ".css")

		if errgcss == nil {
			p.Document.GlobalIndexCssFile = p.Document.Css[cssPath] + "global.css"
		}

		if errcss == nil {
			p.Document.IndexCssFile = p.Document.Css[cssPath] + fileNameNoExt + ".css"
		}

	}

	if jsErr == nil && jsFi.IsDir() {
		jsPath := strings.Trim(p.currentPath, "/")
		DjsPath := p.Config.StaticJsDirectory + jsPath + "/"
		p.Document.Js[jsPath] = siteRootRightTrim + DjsPath[1:]

		_, errgjs := os.Stat(DjsPath + "global.js")
		_, errjs := os.Stat(DjsPath + fileNameNoExt + ".js")

		if errgjs == nil {
			p.Document.GlobalIndexJsFile = p.Document.Js[jsPath] + "global.js"
		}

		if errjs == nil {
			p.Document.IndexJsFile = p.Document.Js[jsPath] + fileNameNoExt + ".js"
		}
	}

	if imgErr == nil && imgFi.IsDir() {
		imgPath := strings.Trim(p.currentPath, "/")
		DimgPath := p.Config.StaticImgDirectory + imgPath + "/"
		p.Document.Img[imgPath] = siteRootRightTrim + DimgPath[1:]
	}
}

func (p *Page) routeTemplate(w http.ResponseWriter, r *http.Request) {
	if p.Config.AutoGenerateHtml {
		p.Document.GenerateHtml = true
	}

	if p.Document.Close == false && p.Document.Hide == false {
		if tplFi, err := os.Stat(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Template); err == nil {

			p.site.base.rmutex.RLock()
			tmplCache := p.GetTemplateCache(p.Template)

			if tplFi.ModTime().Unix() > tmplCache.ModTime {
				p.site.base.mutex.Lock()
				p.SetTemplateCache(p.Template, p.Config.TemplateDirectory+p.Config.ThemeDirectory+p.Template)
				p.site.base.mutex.Unlock()

				tmplCache = p.GetTemplateCache(p.Template)
			}

			globalTemplate, _ := p.site.globalTemplate.Clone()
			pageTemplate, err := globalTemplate.New(filepath.Base(p.Template)).Parse(tmplCache.Content)
			p.site.base.rmutex.RUnlock()

			if err == nil {
				templateVar := map[string]interface{}{
					"G":        p.GET,
					"P":        p.POST,
					"S":        p.SESSION,
					"C":        p.COOKIE,
					"D":        p.Document,
					"Config":   p.Config.M,
					"Template": p.Template,
				}

				p.site.base.rmutex.RLock()
				templateVar["Siteroot"] = p.site.Root
				templateVar["Version"] = p.site.Version
				p.site.base.rmutex.RUnlock()

				if p.Document.GenerateHtml {
					p.site.base.rmutex.RLock()
					htmlFile := p.Config.StaticDirectory + p.Config.HtmlDirectory + p.site.Root + p.Template
					p.site.base.rmutex.RUnlock()

					htmlDir := filepath.Dir(htmlFile)
					if htmlDirFi, err := os.Stat(htmlDir); err != nil || !htmlDirFi.IsDir() {
						os.MkdirAll(htmlDir, 0777)
					}

					var doWrite bool
					if p.Config.AutoGenerateHtml {
						if p.Config.AutoGenerateHtmlCycleTime <= 0 {
							doWrite = true
						} else {
							if htmlFi, err := os.Stat(htmlFile); err != nil {
								doWrite = true
							} else {
								switch {
								case tplFi.ModTime().Unix() >= htmlFi.ModTime().Unix():
									doWrite = true
								case tplFi.ModTime().Unix() >= htmlFi.ModTime().Unix():
									doWrite = true
								case time.Now().Unix()-htmlFi.ModTime().Unix() >= p.Config.AutoGenerateHtmlCycleTime:
									doWrite = true
								default:
									globalTplFi, err := os.Stat(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Config.TemplateGlobalDirectory)
									if err == nil {
										if globalTplFi.ModTime().Unix() >= htmlFi.ModTime().Unix() {
											doWrite = true
										}
									}
								}
							}
						}
					}

					if doWrite {
						if file, err := os.OpenFile(htmlFile, os.O_CREATE|os.O_WRONLY, 0777); err == nil {
							templateVar["Siteroot"] = p.Config.SiteRoot + htmlDir + "/"
							pageTemplate.Execute(file, templateVar)
						}
					}

					if p.Config.AutoJumpToHtml {
						p.site.base.rmutex.RLock()
						http.Redirect(w, r, p.site.Root+htmlFile[2:], http.StatusFound)
						p.site.base.rmutex.RUnlock()
					} else {
						err := pageTemplate.Execute(w, templateVar)
						if err != nil {
							log.Println(err)
						}
					}
				} else {
					err := pageTemplate.Execute(w, templateVar)
					if err != nil {
						log.Println(err)
						w.Write([]byte(fmt.Sprint(err)))
					}
				}
			} else {
				log.Println(err)
				w.Write([]byte(fmt.Sprint(err)))
			}
		} else {
			p.site.base.rmutex.RLock()
			tmplCache := p.GetTemplateCache(p.Template)
			if tmplCache.ModTime > 0 {
				p.site.base.mutex.Lock()
				p.DelTemplateCache(p.Template)
				p.site.base.mutex.Unlock()
			}
			p.site.base.rmutex.RUnlock()
		}
	}
}

func (p *Page) HandleFavicon() {
	http.HandleFunc(p.site.Root+"favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		staticPath := p.Config.StaticDirectory + p.Config.ThemeDirectory + "favicon.ico"
		http.ServeFile(w, r, staticPath)
	})
}

func (p *Page) HandleStatic() {
	http.HandleFunc(p.Document.Static, func(w http.ResponseWriter, r *http.Request) {
		staticPath := p.Config.StaticDirectory + r.URL.Path[len(p.Document.Static):]
		http.ServeFile(w, r, staticPath)
	})
}

func (p *Page) handleRoute(i interface{}) {
	http.HandleFunc(p.site.Root, func(w http.ResponseWriter, r *http.Request) {
		p.site.base.mutex.Lock()
		if p.Config.Reload() {
			p.reset(true)
		}

		p.setCurrentInfo(r.URL.Path)
		p.Template = p.CurrentController + p.currentFileName
		p.site.base.mutex.Unlock()

		p.routeController(i, w, r)
	})
}

func (p *Page) ListenAndServe(addr string, i interface{}) {
	p.handleRoute(i)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

/*
cookie[0] => name string
cookie[1] => value string
cookie[2] => expires string
cookie[3] => path string
cookie[4] => domain string
*/
func (p *Page) SetCookie(w http.ResponseWriter, args ...interface{}) {
	if len(args) < 2 {
		return
	}

	const LEN = 5
	var cookie = [LEN]interface{}{}

	for k, v := range args {
		if k >= LEN {
			break
		}

		cookie[k] = v
	}

	var (
		name    string
		value   string
		expires int
		path    string
		domain  string
	)

	if v, ok := cookie[0].(string); ok {
		name = v
	} else {
		return
	}

	if v, ok := cookie[1].(string); ok {
		value = v
	} else {
		return
	}

	if v, ok := cookie[2].(int); ok {
		expires = v
	}

	if v, ok := cookie[3].(string); ok {
		path = v
	}

	if v, ok := cookie[4].(string); ok {
		domain = v
	}

	pCookie := &http.Cookie{
		Name:   name,
		Value:  url.QueryEscape(value),
		Path:   path,
		Domain: domain,
	}

	if expires > 0 {
		d, _ := time.ParseDuration(strconv.Itoa(expires) + "s")
		pCookie.Expires = time.Now().Add(d)
	}

	http.SetCookie(w, pCookie)
	

	if expires > 0 {
		p.COOKIE[pCookie.Name] = pCookie.Value
	} else {
		delete(p.COOKIE, pCookie.Name)
	}
}

func excuteActionFunc(rm reflect.Method,vpc reflect.Value,action string, rvr reflect.Value, rvw reflect.Value ,ppc *Page){

	if ppc.Document.Close == false {
		mt := rm.Type
		switch mt.NumIn() {
		case 2:
			if mt.In(1) == rvr.Type() {
				vpc.MethodByName(action).Call([]reflect.Value{rvr})
			} else {
				vpc.MethodByName(action).Call([]reflect.Value{rvw})
			}
		case 3:
			vpc.MethodByName(action).Call([]reflect.Value{rvw, rvr})
		default:
			vpc.MethodByName(action).Call([]reflect.Value{})
		}
	}
}

