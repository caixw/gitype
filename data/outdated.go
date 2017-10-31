// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/caixw/gitype/vars"
)

type outdatedServer struct {
	updated         time.Time // 最后更新时间
	etag            string
	postsTicker     *time.Ticker
	postsTickerDone chan bool
	duration        time.Duration
	posts           []*Post
}

func (d *Data) newOutdatedServer(conf *config) {
	if conf.Outdated == 0 {
		return
	}

	srv := &outdatedServer{
		updated:         time.Now(),
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
}

// StartOutdatedService 开始 Outdated 更新服务
func (d *Data) StartOutdatedService() {
	d.outdatedServer.run()
}

func (srv *outdatedServer) run() {
	go func() {
		for {
			select {
			case <-srv.postsTicker.C:
				srv.updateOutdated()
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

func (srv *outdatedServer) updateOutdated() {
	if srv.duration == 0 {
		return
	}

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

	srv.updated = now
	srv.etag = vars.Etag(now)
}

// Updated 最后的更新时间
func (d *Data) Updated() time.Time {
	return d.outdatedServer.updated
}

// Etag 根据更新时间生成的 etag 值
func (d *Data) Etag() string {
	return d.outdatedServer.etag
}
