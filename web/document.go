package web

import (
	"text/template"
)

type Document struct {
	Close              bool
	GenerateHtml       bool
	Static             string
	Theme              string
	Css                map[string]string
	Js                 map[string]string
	Img                map[string]string
	GlobalCssFile      string
	GlobalJsFile       string
	GlobalIndexCssFile string
	GlobalIndexJsFile  string
	IndexCssFile       string
	IndexJsFile        string
	Hide               bool
	Func               template.FuncMap
	Title              string
	Subtitle           string
	Header             string
	Body               interface{}
	Footer             string
}

/*func (d *Document) Reset() {
	document := d
	globalCss, okCss := document.Css["global"]
	globalJs, okJs := document.Js["global"]
	globalImg, okImg := document.Img["global"]
	*d = Document{
		Static:        document.Static,
		Theme:         document.Theme,
		GlobalCssFile: document.GlobalCssFile,
		GlobalJsFile:  document.GlobalJsFile,
		Css:           map[string]string{},
		Js:            map[string]string{},
		Img:           map[string]string{},
		Func:          template.FuncMap{},
	}

	if okCss {
		d.Css["global"] = globalCss
	}

	if okJs {
		d.Js["global"] = globalJs
	}

	if okImg {
		d.Img["global"] = globalImg
	}
}*/
