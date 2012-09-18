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
