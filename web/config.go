package web

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
)

type Config struct {
	SupportSession            bool                   `json:"SupportSession"`
	AutoGenerateHtml          bool                   `json:"AutoGenerateHtml"`
	AutoGenerateHtmlCycleTime int64                  `json:"AutoGenerateHtmlCycleTime"`
	AutoJumpToHtml            bool                   `json:"AutoJumpToHtml"`
	Debug                     bool                   `json:"Debug"`
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
	M                         map[string]interface{} `json:"Custom"`
	configPath                string
	configLastModTime         int64
}

func NewConfig() Config {
	return Config{
		TemplateDirectory:       "./view/",
		TemporaryDirectory:      "./tmp/",
		StaticDirectory:         "./static/",
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
	}
}

func (c *Config) format(configPath string) []byte {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	return regexp.MustCompile(`#.*\n`).ReplaceAll(data, []byte("\n"))
}

func (c *Config) Load(configPath string) {
	data := c.format(configPath)

	err := json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}

	c.UploadDirectory = c.StaticDirectory + c.UploadDirectory
	c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
	c.StaticCssDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
	c.StaticJsDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
	c.StaticImgDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory

	c.configPath = configPath
	dataFi, _ := os.Stat(configPath)
	c.configLastModTime = dataFi.ModTime().Unix()
}

func (c *Config) Reload() bool {
	var b bool
	configPath := c.configPath
	dataFi, _ := os.Stat(configPath)
	if dataFi.ModTime().Unix() > c.configLastModTime {
		data := c.format(configPath)
		*c = NewConfig()
		json.Unmarshal(data, c)
		c.configPath = configPath
		c.configLastModTime = dataFi.ModTime().Unix()
		c.UploadDirectory = c.StaticDirectory + c.UploadDirectory
		c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
		c.StaticCssDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
		c.StaticJsDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
		c.StaticImgDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory
		b = true
	}

	return b
}
