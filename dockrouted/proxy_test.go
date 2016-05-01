package dockrouted

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func init() {
	log.SetFormatter(&log.TextFormatter{DisableColors: true, DisableTimestamp: true})
}

func Test_Proxy_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := assert.New(t)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bar"))
	}))
	defer s.Close()
	url, err := url.Parse(s.URL)
	a.NoError(err)

	backend := NewMockBackend(ctrl)
	backend.EXPECT().GetService("foo.com").Return(url.Host, nil)

	p := NewProxy(backend)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://foo.com/xyz", nil)
	p.ServeHTTP(w, r)

	a.Equal("bar", w.Body.String())
}

func Test_Proxy_ServeHTTP_LookupErrro(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := assert.New(t)

	backend := NewMockBackend(ctrl)
	backend.EXPECT().GetService("foo.com").Return("", errors.New("no such host"))

	p := NewProxy(backend)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://foo.com/xyz", nil)
	p.ServeHTTP(w, r)

	a.Equal(502, w.Code)
}
