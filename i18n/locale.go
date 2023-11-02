package i18n

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

type Locale interface {
	Language() string
	Format(string, ...any) string
	Get(string) string
	GetOrDefault(string, string) string
	Exists(string) bool
	Refresh()
}

type Store interface {
	AddLocale(...Locale)
	GetLocale(string) (Locale, bool)
	SetDefaultLocale(string) error
	GetDefaultLocale() Locale
}

type defaultLocale struct{}

func (*defaultLocale) Language() string {
	return "default"
}

func (*defaultLocale) Format(string, ...any) string {
	return ""
}

func (*defaultLocale) Get(string) string {
	return ""
}

func (*defaultLocale) GetOrDefault(string, string) string {
	return ""
}

func (*defaultLocale) Exists(string) bool {
	return false
}

func (*defaultLocale) Refresh() {}

type SimpleLocale struct {
	m    map[string]string
	lang string
}

func NewSimpleLocale(data map[string]string, lang string) *SimpleLocale {
	if data == nil {
		data = make(map[string]string)
	}
	return &SimpleLocale{
		m:    data,
		lang: lang,
	}
}

func (l *SimpleLocale) Language() string {
	return l.lang
}

func (l *SimpleLocale) Format(k string, args ...any) string {
	return fmt.Sprintf(l.Get(k), args...)
}

func (l *SimpleLocale) Get(k string) string {
	return l.m[k]
}

func (l *SimpleLocale) GetOrDefault(k string, d string) string {
	v, b := l.m[k]
	if b {
		return v
	}
	return d
}

func (l *SimpleLocale) Exists(k string) bool {
	_, b := l.m[k]
	return b
}

func (l *SimpleLocale) Refresh() {}

func NewImmutableLocaleFromIniFile(fileName, lang string) (Locale, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return NewImmutableLocaleFromIni(content, lang)
}

func NewImmutableLocaleFromIni(content []byte, lang string) (Locale, error) {
	iniContent, err := ini.Load(content)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, 8)
	sections := iniContent.Sections()
	for _, section := range sections {
		for _, key := range section.Keys() {
			var k string
			if section.Name() == "" || section.Name() == "DEFAULT" {
				k = key.Name()
			} else {
				k = section.Name() + "." + key.Name()
			}
			m[k] = key.Value()
		}
	}
	return NewSimpleLocale(m, lang), nil
}
