package session

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"time"

	utls "github.com/refraction-networking/utls"
)

type Session struct {
	*http.Client
}

func New(opts ...Option) *Session {
	session := &Session{
		Client: &http.Client{},
	}
	for _, opt := range opts {
		opt(session)
	}
	return session
}

type Option func(*Session)

// WithClientHelloID is used to set tls config
func WithClientHelloID(clientHelloID utls.ClientHelloID) Option {
	c, err := utls.UTLSIdToSpec(clientHelloID)
	if err != nil {
		slog.Error("tls config error", "err", err)
		return func(s *Session) {}
	}
	transport := &http.Transport{
		DisableKeepAlives: false,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         c.TLSVersMin,
			MaxVersion:         c.TLSVersMax,
			CipherSuites:       c.CipherSuites,
			ClientSessionCache: tls.NewLRUClientSessionCache(32),
		},
	}
	return func(s *Session) {
		s.Client.Transport = transport
	}
}

// WithTransport is used to set transport
func WithTransport(transport *http.Transport) Option {
	return func(s *Session) {
		s.Client.Transport = transport
	}
}

// WithTimeout is used to set timeout
func WithTimeout(timeout time.Duration) Option {
	return func(s *Session) {
		s.Client.Timeout = timeout
	}
}

// WithCookieJar is used to set cookie jar
func WithCookieJar(jar http.CookieJar) Option {
	return func(s *Session) {
		s.Client.Jar = jar
	}
}
