package browserstorage

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
)

type Storage struct {
	browsers map[int64]*rod.Browser
}

func New() *Storage {
	return &Storage{
		browsers: make(map[int64]*rod.Browser),
	}
}

func (s *Storage) GetOrCreate(id int64, proxyStr string) (*rod.Browser, error) {
	candidate, ok := s.browsers[id]
	if ok {
		return candidate, nil
	}

	parts := strings.Split(proxyStr, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid proxy format, expected ip:port@username:password")
	}

	authParts := strings.Split(parts[1], ":")
	if len(authParts) != 2 {
		return nil, fmt.Errorf("invalid username:password format")
	}

	proxy := parts[0]
	user := authParts[0]
	password := authParts[1]
	l := launcher.New().Headless(false)
	l = l.Set(flags.ProxyServer, proxy)

	controlURL, _ := l.Launch()
	browser := rod.New().ControlURL(controlURL).MustConnect()

	go browser.MustHandleAuth(user, password)()

	s.browsers[id] = browser

	return browser, nil
}

func (s *Storage) Remove(id int64) error {
	candidate, ok := s.browsers[id]
	if !ok {
		return fmt.Errorf("браузер не существует")
	}

	if err := candidate.Close(); err != nil {
		return err
	}

	delete(s.browsers, id)

	return nil
}
