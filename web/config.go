package web

import (
	"bytes"
	"encoding/json"
	"golanger.com/log"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	regexpNote  *regexp.Regexp = regexp.MustCompile(`#.*`)
	regexpSpace *regexp.Regexp = regexp.MustCompile(`[\n\v\f\r]*`)
)

type Config struct {
	SupportLog                bool                   `json:"SupportLog"`
	LogWriteTo                string                 `json:"LogWriteTo"`
	LogLevel                  string                 `json:"LogLevel"`
	SupportTemplate           bool                   `json:"SupportTemplate"`
	SupportSession            bool                   `json:"SupportSession"`
	SupportCookieSession      bool                   `json:"SupportCookieSession"`
	SupportI18n               bool                   `json:"SupportI18n"`
	SupportStatic             bool                   `json:"SupportStatic"`
	SupportUrlManage          bool                   `json:"SupportUrlManage"`
	SupportUrlManageWithCache bool                   `json:"SupportUrlManageWithCache"`
	SessionType               string                 `json:"SessionType"`
	RootStaticFiles           string                 `json:"RootStaticFiles"`
	DefaultLocalePath         string                 `json:"DefaultLocalePath"`
	DefaultLanguage           string                 `json:"DefaultLanguage"`
	AutoGenerateHtml          bool                   `json:"AutoGenerateHtml"`
	AutoGenerateHtmlCycleTime int64                  `json:"AutoGenerateHtmlCycleTime"`
	AutoLoadStaticHtml        bool                   `json:"AutoLoadStaticHtml"`
	LoadStaticHtmlWithLogic   bool                   `json:"LoadStaticHtmlWithLogic"`
	ChangeSiteRoot            bool                   `json:"ChangeSiteRoot"`
	AccessHtml                bool                   `json:"AccessHtml"`
	AutoJumpToHtml            bool                   `json:"AutoJumpToHtml"`
	AssetsDirectory           string                 `json:"AssetsDirectory"`
	StaticDirectory           string                 `json:"StaticDirectory"`
	ThemeDirectory            string                 `json:"ThemeDirectory"`
	Theme                     string                 `json:"Theme"`
	StaticCssDirectory        string                 `json:"StaticCssDirectory"`
	StaticJsDirectory         string                 `json:"StaticJsDirectory"`
	StaticImgDirectory        string                 `json:"StaticImgDirectory"`
	HtmlDirectory             string                 `json:"HtmlDirectory"`
	TemplateDirectory         string                 `json:"TemplateDirectory"`
	TemplateGlobalDirectory   string                 `json:"TemplateGlobalDirectory"`
	TemplateGlobalFile        string                 `json:"TemplateGlobalFile"`
	TemporaryDirectory        string                 `json:"TemporaryDirectory"`
	UploadDirectory           string                 `json:"UploadDirectory"`
	IndexDirectory            string                 `json:"IndexDirectory"`
	IndexPage                 string                 `json:"IndexPage"`
	SiteRoot                  string                 `json:"SiteRoot"`
	Environment               map[string]string      `json:"Environment"`
	Database                  map[string]string      `json:"Database"`
	UrlManageRule             []string               `json:"UrlManageRule"`
	M                         map[string]interface{} `json:"Custom"`
	configDir                 string
	configLastModTime         int64
}

func NewConfig() Config {
	return Config{
		SupportUrlManageWithCache: true,
		LogWriteTo:                "console",
		SessionType:               "memory",
		RootStaticFiles:           "favicon.ico",
		TemplateDirectory:         "./view/",
		TemporaryDirectory:        "./tmp/",
		AssetsDirectory:           "./assets/",
		StaticDirectory:           "static/",
		DefaultLocalePath:         "./config/locale/",
		DefaultLanguage:           "zh-cn",
		ThemeDirectory:            "theme/",
		Theme:                     "default",
		StaticCssDirectory:        "css/",
		StaticJsDirectory:         "js/",
		StaticImgDirectory:        "img/",
		HtmlDirectory:             "html/",
		UploadDirectory:           "upload/",
		TemplateGlobalDirectory:   "_global/",
		TemplateGlobalFile:        "*",
		IndexDirectory:            "index/",
		IndexPage:                 "index.html",
		SiteRoot:                  "/",
		Environment:               map[string]string{},
		Database:                  map[string]string{},
		UrlManageRule:             []string{},
	}
}

func (c *Config) format(configPath string) []byte {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("<Config.format> error: ", err)
	}

	return bytes.TrimSpace(regexpSpace.ReplaceAll(regexpNote.ReplaceAll(data, []byte(``)), []byte(``)))
}

func (c *Config) readDir(configDir string) []byte {
	fis, err := ioutil.ReadDir(configDir)
	if err != nil {
		log.Fatal("<Config.readDir> error: ", err)
	}

	lfis := len(fis)
	chContent := make(chan []byte, lfis)

	for _, fi := range fis {
		fiName := fi.Name()
		if fi.IsDir() || fiName[0] == '.' {
			lfis--
			continue
		}

		go func(chContent chan []byte, configPath string) {
			chContent <- c.format(configPath)
		}(chContent, configDir+"/"+fiName)
	}

	contentBuf := bytes.NewBufferString(`{`)
	for i := 1; i <= lfis; i++ {
		content := <-chContent
		if len(content) == 0 {
			continue
		}

		contentBuf.Write(content)
		if i < lfis {
			contentBuf.WriteString(",")
		}
	}

	contentBuf.WriteString(`}`)

	return contentBuf.Bytes()
}

func (c *Config) load(data []byte) {
	err := json.Unmarshal(data, c)
	if err != nil {
		log.Debug("<Config.Load> jsonData: ", string(data))
		log.Fatal("<Config.Load> error: ", err)
	}

	c.UploadDirectory = c.AssetsDirectory + c.StaticDirectory + c.UploadDirectory
	c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
	c.StaticCssDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
	c.StaticJsDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
	c.StaticImgDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory
}

func (c *Config) LoadData(data string) {
	c.load([]byte("{" + data + "}"))
}

func (c *Config) Load(configDir string) {
	data := c.readDir(configDir)
	c.load(data)
	c.configDir = configDir
	dataFi, _ := os.Stat(c.configDir)
	c.configLastModTime = dataFi.ModTime().Unix()
}

func (c *Config) Reload() bool {
	var b bool
	configDir := c.configDir
	if configDir == "" {
		return false
	}

	if dataFi, err := os.Stat(configDir); err == nil {
		if dataFi.ModTime().Unix() > c.configLastModTime {
			data := c.readDir(configDir)
			*c = NewConfig()
			json.Unmarshal(data, c)
			c.configDir = configDir
			c.configLastModTime = dataFi.ModTime().Unix()
			c.UploadDirectory = c.AssetsDirectory + c.StaticDirectory + c.UploadDirectory
			c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
			c.StaticCssDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
			c.StaticJsDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
			c.StaticImgDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory
			b = true
		}
	}

	return b
}
