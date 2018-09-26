// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/vars"
)

// 定时更新所有文章的过时信息的一个服务。
//
// 文章一旦过时，则提示信息中的过时天数会每天变化，
// outdatedServer 即为按一定时间启动的定时器，
// 启动频率可以通过 vars.OutdatedFrequency 进行调整。
type outdatedServer struct {
	postsTicker     *time.Ticker
	postsTickerDone chan bool
	duration        time.Duration
	posts           []*Post
}

func (d *Data) initOutdatedServer(conf *loader.Config) {
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
	go d.runTickerServer()
}

func (d *Data) runTickerServer() {
	srv := d.outdatedServer
	for {
		select {
		case <-srv.postsTicker.C:
			d.updateOutdated()
		case <-srv.postsTickerDone:
			return
		}
	}
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
		if post.Outdated.Type != loader.OutdatedTypeCreated &&
			post.Outdated.Type != loader.OutdatedTypeModified {
			continue
		}

		outdated := now.Sub(post.Outdated.Date)
		if outdated >= srv.duration {
			post.Outdated.Days = int(outdated.Hours()) / 24
		}
	}

	d.setUpdated(now)
}
