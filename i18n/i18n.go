package i18n

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type I18nManager struct {
	I18nName        string
	localePath      string
	defaultLanguage string
	rmutex          sync.RWMutex
	mutex           sync.Mutex
	languages       map[string]map[string]string
	lastModTime     map[string]int64
}

func New(name string) *I18nManager {
	if name == "" {
		name = "GoLangerI18n"
	}

	i := &I18nManager{
		I18nName:    name,
		languages:   map[string]map[string]string{},
		lastModTime: map[string]int64{},
	}

	return i
}

func (i *I18nManager) Setup(path, lang string) error {
	i.localePath = path
	i.defaultLanguage = lang

	return i.Load(i.defaultLanguage)
}

func (i *I18nManager) Load(lang string) error {
	i.rmutex.RLock()
	_, found := i.languages[lang]
	oldLastModTime, mFound := i.lastModTime[lang]
	i.rmutex.RUnlock()

	langFile := i.localePath + lang
	newer := false
	dataFi, err := os.Stat(langFile)
	if err == nil {
		lastModTime := dataFi.ModTime().Unix()
		if !mFound {
			i.rmutex.Lock()
			i.lastModTime[lang] = lastModTime
			i.rmutex.Unlock()
			newer = true
		} else {
			if lastModTime > oldLastModTime {
				newer = true
			}
		}

		if found && !newer {
			return nil
		}

		data, _ := ioutil.ReadFile(langFile)
		m := map[string]string{}
		err = json.Unmarshal(data, &m)
		if err == nil {
			i.mutex.Lock()
			i.languages[lang] = m
			i.mutex.Unlock()
		}
	}

	return err
}

func (i *I18nManager) Lang(lang string) map[string]string {
	lang = strings.ToLower(lang)
	err := i.Load(lang)
	if err != nil {
		lang = i.defaultLanguage
	}

	i.rmutex.RLock()
	msgs := i.languages[lang]
	i.rmutex.RUnlock()

	return msgs
}
