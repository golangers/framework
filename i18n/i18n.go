package i18n

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
}

func New(name, path, lang string) *I18nManager {
	if name == "" {
		name = "GoLangerI18n"
	}

	if path == "" {
		path = "./locale/"
	}
	if lang == "" {
		lang = "zh-cn"
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

func (i *I18nManager) loadLanguageFile(lang string) error {
	if lang == "" {
		lang = i.defaultLanguage
	}

	i.rmutex.RLock()
	_, ok := i.languages[i.defaultLanguage]
	if ok {
		i.rmutex.RUnlock()
		return nil
	}

	i.rmutex.RUnlock()

	file := i.localePath + lang
	data, err := loadFile(file)
	if err != nil {
		return errors.New("Error: Loading Language File " + file)
	}
	i.mutex.Lock()
	err = json.Unmarshal(data, &i.languages)
	i.mutex.Unlock()

	return err
}

func (i *I18nManager) Lang(l string) map[string]string {
	if l == "" {
		l = i.defaultLanguage
	} else {
		l = strings.ToLower(l)
	}

	i.rmutex.RLock()
	defer i.rmutex.RUnlock()
	msgs, found := i.languages[l]
	if !found {
		// Load The Language File
		err := i.loadLanguageFile(l)
		if err != nil {
			return map[string]string{}
		}
		msgs = i.languages[l]
	}

	return msgs
}

func (i *I18nManager) Get(lang, key string) string {
	if lang == "" {
		lang = i.defaultLanguage
	} else {
		lang = strings.ToLower(lang)
	}

	targetLang := ""
	i.rmutex.RLock()
	msgs, found := i.languages[lang]
	if !found {
		// Load The Language File
		err := i.loadLanguageFile(lang)
		if err != nil {
			i.rmutex.RUnlock()
			return key
		}

		msgs = i.languages[lang]
	}

	targetLang, ok := msgs[key]
	i.rmutex.RUnlock()

	if !ok {
		targetLang = key
	}

	return targetLang
}
