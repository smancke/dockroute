package dockrouted

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"net/http/httputil"
)

type Proxy struct {
	backend    Backend
	hostSuffix string
}

func NewProxy(b Backend) *Proxy {
	return &Proxy{
		backend:    b,
		hostSuffix: ".proxy",
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSuffix(r.Host, p.hostSuffix)
	host, port, err := p.backend.GetService(name)
	if err != nil {
		rLog(r).WithError(err).Error("Error fetching service.")
		httpError(w, http.StatusBadGateway)
		return
	}

	director := func(req *http.Request) {
		rLog(r).Info("req before")
		if r.TLS != nil {
			req.URL.Scheme = "https"
		} else {
			req.URL.Scheme = "http"
		}
		req.URL.Host = host + ":" + port
		req.Host = req.URL.Host
		rLog(r).Info("req after")
	}
	reverseProxy := &httputil.ReverseProxy{Director: director}
	reverseProxy.ServeHTTP(w, r)
	rLog(r).WithField("target_host", host).
		WithField("target_port", host).
		Info("served.")
}

func rLog(r *http.Request) *log.Entry {
	return log.WithFields(log.Fields{
		"RemoteAddr": r.RemoteAddr,
		"URL":        r.URL,
		"Method":     r.Method,
		"Host":       r.Host,
	})
}

func httpError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
