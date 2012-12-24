package web

import (
	"bytes"
	"encoding/json"
	"golanger.com/framework/log"
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
	SessionType               string                 `json:"SessionType"`
	RootStaticFiles           string                 `json:"RootStaticFiles"`
	DefaultLanguage           string                 `json:"DefaultLanguage"`
	DefaultLocalePath         string                 `json:"DefaultLocalePath"`
	AutoGenerateHtml          bool                   `json:"AutoGenerateHtml"`
	AutoGenerateHtmlCycleTime int64                  `json:"AutoGenerateHtmlCycleTime"`
	AutoLoadStaticHtml        bool                   `json:"AutoLoadStaticHtml"`
	LoadStaticHtmlWithLogic   bool                   `json:"LoadStaticHtmlWithLogic"`
	ChangeSiteRoot            bool                   `json:"ChangeSiteRoot"`
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
		LogWriteTo:              "console",
		SessionType:             "memory",
		RootStaticFiles:         "favicon.ico",
		TemplateDirectory:       "./view/",
		TemporaryDirectory:      "./tmp/",
		AssetsDirectory:         "./assets/",
		StaticDirectory:         "static/",
		DefaultLanguage:         "zh-cn",
		DefaultLocalePath:       "./locale/",
		ThemeDirectory:          "theme/",
		Theme:                   "default",
		StaticCssDirectory:      "css/",
		StaticJsDirectory:       "js/",
		StaticImgDirectory:      "img/",
		HtmlDirectory:           "html/",
		UploadDirectory:         "upload/",
		TemplateGlobalDirectory: "_global/",
		TemplateGlobalFile:      "*",
		IndexDirectory:          "index/",
		IndexPage:               "index.html",
		SiteRoot:                "/",
		Environment:             map[string]string{},
		Database:                map[string]string{},
		UrlManageRule:           []string{},
	}
}

func (c *Config) format(configPath string) []byte {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("<Config.format> error: ", err)
	}

	return regexpSpace.ReplaceAll(regexpNote.ReplaceAll(data, []byte(``)), []byte(``))
}

func (c *Config) readDir(configDir string) []byte {
	fis, err := ioutil.ReadDir(configDir)
	if err != nil {
		log.Fatal("<Config.readDir> error: ", err)
	}

	lfis := len(fis)
	chContent := make(chan []byte, lfis)

	for _, fi := range fis {
		go func(chContent chan []byte, configPath string) {
			chContent <- c.format(configPath)
		}(chContent, configDir+"/"+fi.Name())
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

func (c *Config) Load(configDir string) {
	data := c.readDir(configDir)
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

	c.configDir = configDir
	dataFi, _ := os.Stat(configDir)
	c.configLastModTime = dataFi.ModTime().Unix()
}

func (c *Config) Reload() bool {
	var b bool
	configDir := c.configDir
	dataFi, _ := os.Stat(configDir)
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

	return b
}
