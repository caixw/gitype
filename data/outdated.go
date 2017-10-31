// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/caixw/gitype/vars"
)

type outdatedServer struct {
	postsTicker     *time.Ticker
	postsTickerDone chan bool
	duration        time.Duration
	posts           []*Post
}

func (d *Data) initOutdatedServer(conf *config) {
	if conf.Outdated == 0 {
		return
	}

	srv := &outdatedServer{
		postsTicker:     time.NewTicker(vars.OutdatedFrequency),
		postsTickerDone: make(chan bool, 1),
		duration:        conf.Outdated,
		posts:           make([]*Post, 0, len(d.Posts)),
	}

	for _, post := range d.Posts {
		if post.Outdated == nil {
			continue
		}

		srv.posts = append(srv.posts, post)
	}

	d.outdatedServer = srv

	d.updateOutdated()
	d.runTickerServer()
}

func (d *Data) runTickerServer() {
	go func() {
		srv := d.outdatedServer
		for {
			select {
			case <-srv.postsTicker.C:
				d.updateOutdated()
			case <-srv.postsTickerDone:
				return
			}
		}
	}()
}

func (srv *outdatedServer) stop() {
	if srv.postsTicker != nil {
		srv.postsTicker.Stop()
		srv.postsTickerDone <- true
	}
}

func (d *Data) updateOutdated() {
	srv := d.outdatedServer

	now := time.Now()
	for _, post := range srv.posts {
		if post.Outdated.Type != outdatedTypeCreated &&
			post.Outdated.Type != outdatedTypeModified {
			continue
		}

		outdated := now.Sub(post.Outdated.Date)
		if outdated >= srv.duration {
			post.Outdated.Days = int(outdated.Hours()) / 24
		}
	}

	d.setUpdated(now)
}
