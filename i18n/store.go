package i18n

import (
	"errors"
	"github.com/LeeZXin/zsf-utils/concurrentutil"
	"github.com/LeeZXin/zsf-utils/maputil"
)

var (
	defaultStore = NewDefaultStore()
)

type storeImpl struct {
	m *maputil.ConcurrentMap[string, Locale]
	l *concurrentutil.Value[Locale]
}

func NewDefaultStore() Store {
	l := concurrentutil.NewValue[Locale]()
	var d Locale = new(defaultLocale)
	l.Store(d)
	return &storeImpl{
		m: maputil.NewConcurrentMap[string, Locale](nil),
		l: l,
	}
}

func (s *storeImpl) AddLocale(l ...Locale) {
	for _, i := range l {
		s.m.Store(i.Language(), i)
	}
}

func (s *storeImpl) GetLocale(k string) (Locale, bool) {
	return s.m.Load(k)
}

func (s *storeImpl) SetDefaultLocale(k string) error {
	l, b := s.GetLocale(k)
	if !b {
		return errors.New("locale not found")
	}
	s.l.Store(l)
	return nil
}

func (s *storeImpl) GetDefaultLocale() Locale {
	return s.l.Load()
}

func (s *storeImpl) Format(k string, args ...any) string {
	return s.l.Load().Format(k, args...)
}

func (s *storeImpl) Get(k string) string {
	return s.l.Load().Get(k)
}

func (s *storeImpl) GetOrDefault(k string, d string) string {
	return s.l.Load().GetOrDefault(k, d)
}

func (s *storeImpl) Exists(k string) bool {
	return s.l.Load().Exists(k)
}

func AddLocale(l ...Locale) {
	defaultStore.AddLocale(l...)
}

func GetLocale(k string) (Locale, bool) {
	return defaultStore.GetLocale(k)
}

func GetDefaultLocale() Locale {
	return defaultStore.GetDefaultLocale()
}

func SetDefaultLocale(k string) error {
	return defaultStore.SetDefaultLocale(k)
}

func Get(k string) string {
	return GetDefaultLocale().Get(k)
}

func Format(k string, args ...any) string {
	return GetDefaultLocale().Format(k, args)
}

func GetOrDefault(k string, d string) string {
	return GetDefaultLocale().GetOrDefault(k, d)
}

func Exists(k string) bool {
	return GetDefaultLocale().Exists(k)
}
