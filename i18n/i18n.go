package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
)

type I18nManager struct {
	I18nName        string
	localePath      string
	defaultLanguage string
	rmutex          sync.RWMutex
	mutex           sync.Mutex
	languages       map[string]map[string]string
}

func New(name, path, lang string) *I18nManager {
	if name == "" {
		name = "GoLangerI18n"
	}

	if path == "" {
		path = "./"
	}
	if lang == "" {
		lang = "zh-CN"
	}

	i := &I18nManager{
		I18nName:        name,
		localePath:      path,
		defaultLanguage: lang,
		languages:       map[string]map[string]string{},
	}

	return i
}

func loadFile(filepath string) (data []byte, err error) {
	if filepath == "" {
		return nil, errors.New("Filepath Error")
	}
	data, err = ioutil.ReadFile(filepath)
	return
}

func (i *I18nManager) loadLanguageFile(lang string) (err error) {
	if lang == "" {
		lang = i.defaultLanguage
	}

	file := fmt.Sprintf("%s/%s.json", i.localePath, lang)
	data, err := loadFile(file)
	if err != nil {
		return errors.New("Error: Loading Language File " + file)
	}
	i.mutex.Lock()
	defer i.mutex.Unlock()
	err = json.Unmarshal(data, &i.languages)
	return
}

func (i *I18nManager) Get(lang, key string) string {
	i.rmutex.RLock()
	if _, ok := i.languages[lang]; !ok {
		i.rmutex.RUnlock()
		// Load The Language File
		err := i.loadLanguageFile(lang)
		if err != nil {
			return ""
		}
	}
	i.rmutex.RUnlock()

	i.rmutex.RLock()
	defer i.rmutex.RUnlock()
	if msgs, ok := i.languages[lang]; ok {
		if value, ok := msgs[key]; ok {
			return value
		}
	}
	return ""
}
