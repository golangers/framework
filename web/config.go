package web

import (
	"bytes"
	"encoding/json"
	"golanger.com/config"
	"golanger.com/log"
	"io/ioutil"
	"os"
	"regexp"
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

func (c *Config) Init() Config {
	c.preSet()
	return *c
}

func (c *Config) preSet() {
	c.UploadDirectory = c.AssetsDirectory + c.StaticDirectory + c.UploadDirectory
	c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
	c.StaticCssDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
	c.StaticJsDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
	c.StaticImgDirectory = c.AssetsDirectory + c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory
}

func (c *Config) set(configDir string, configLastModTime int64) {
	c.configDir = configDir
	c.configLastModTime = dataFi.ModTime().Unix()
}

func (c *Config) LoadData(data string) {
	conf := config.Data(data).Load(c)
	c.preSet()
}

func (c *Config) Load(configDir string) {
	conf := config.Dir(configDir).Load(c)
	configDir = conf.Target()
	if dataFi, err := os.Stat(configDir); err == nil {
		c.set(configDir, dataFi.ModTime().Unix())
		c.preSet()
	} else {
		log.Fatal("<Config.Load> error:", err)
	}
}

func (c *Config) Reload() bool {
	var b bool
	configDir := c.configDir
	if configDir == "" {
		return false
	}

	if dataFi, err := os.Stat(configDir); err == nil {
		if dataFi.ModTime().Unix() > c.configLastModTime {
			cm := NewConfig()
			cm.Load(configDir)
			*c = cm
			b = true
		}
	}

	return b
}
