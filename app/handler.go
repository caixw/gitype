// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/caixw/typing/helper"
	"github.com/issue9/logs"
	"github.com/issue9/middleware/host"
	"github.com/issue9/middleware/recovery"
)

func (a *app) buildHandler() http.Handler {
	h := a.buildDomains(a.buildHeader(a.mux))

	h = recovery.New(h, func(w http.ResponseWriter, msg interface{}) {
		logs.Error(msg)
		helper.StatusError(w, http.StatusInternalServerError)
	})

	// 将 pprof 包装在最外层
	return a.buildPprof(h)
}

func (a *app) buildDomains(h http.Handler) http.Handler {
	if len(a.conf.Domains) == 0 {
		return h
	}

	return host.New(h, a.conf.Domains...)
}

func (a *app) buildHeader(h http.Handler) http.Handler {
	if len(a.conf.Headers) == 0 {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range a.conf.Headers {
			w.Header().Set(k, v)
		}
		h.ServeHTTP(w, r)
	})
}

// 根据 Config.Pprof 决定是否包装调试地址，调用前请确认是否已经开启 Pprof 选项
func (a *app) buildPprof(h http.Handler) http.Handler {
	if !a.conf.Pprof {
		return h
	}

	logs.Debug("开启了调试功能，地址为：", debugPprof)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, debugPprof) {
			h.ServeHTTP(w, r)
			return
		}

		switch r.URL.Path[len(debugPprof):] {
		case "cmdline":
			pprof.Cmdline(w, r)
		case "profile":
			pprof.Profile(w, r)
		case "symbol":
			pprof.Symbol(w, r)
		case "trace":
			pprof.Trace(w, r)
		default:
			pprof.Index(w, r)
		}
	}) // end return http.HandlerFunc
}
