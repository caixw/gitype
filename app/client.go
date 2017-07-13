// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"net/http"
)

func (a *app) initFeeds() {
	conf := a.buf.Data.Config

	if conf.RSS != nil {
		a.mux.GetFunc(conf.RSS.URL, a.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.buf.RSS)
		}))
	}

	if conf.Atom != nil {
		a.mux.GetFunc(conf.Atom.URL, a.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.buf.Atom)
		}))
	}

	if conf.Sitemap != nil {
		a.mux.GetFunc(conf.Sitemap.URL, a.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.buf.Sitemap)
		}))
	}
}

func (a *app) removeFeeds() {
	conf := a.buf.Data.Config

	if conf.RSS != nil {
		a.mux.Remove(conf.RSS.URL)
	}

	if conf.Atom != nil {
		a.mux.Remove(conf.Atom.URL)
	}

	if conf.Sitemap != nil {
		a.mux.Remove(conf.Sitemap.URL)
	}
}
