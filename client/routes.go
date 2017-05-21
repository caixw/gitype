// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"path/filepath"

	"github.com/issue9/mux"
	"github.com/issue9/utils"
)

func (c *Client) initRoutes() error {
	urls := c.data.Config.URLS

	c.mux.Prefix(urls.Root).GetFunc(urls.Post+"/{slug}"+urls.Suffix, c.pre(c.getPost)).
		GetFunc(urls.Posts+urls.Suffix, c.pre(c.getPosts)).
		GetFunc(urls.Tag+"/{slug}"+urls.Suffix, c.pre(c.getTag)).
		GetFunc(urls.Tags+urls.Suffix, c.pre(c.getTags)).
		GetFunc(urls.Search+urls.Suffix, c.pre(c.getSearch)).
		GetFunc(urls.Media+"/*", c.pre(c.getMedia)).
		GetFunc(urls.Themes+"/*", c.pre(c.getThemes)).
		GetFunc("/*", c.pre(c.getRaws))
	return nil
}

func (c *Client) getPost(w http.ResponseWriter, r *http.Request) {
	ps := mux.GetParams(r)
	if ps == nil {
		// TODO
	}
}

func (c *Client) getPosts(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getTag(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getTags(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getThemes(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getSearch(w http.ResponseWriter, r *http.Request) {
	//
}

// 获取媒体文件
//
// /media/2015/intro-php/content.html ==> /posts/2015/intro-php/content.html
func (c *Client) getMedia(w http.ResponseWriter, r *http.Request) {
	media := c.data.Config.URLS.Media
	dir := http.Dir(c.data.PostsPath(""))
	http.StripPrefix(media, http.FileServer(dir)).ServeHTTP(w, r)
}

// 读取根下的文件
func (c *Client) getRaws(w http.ResponseWriter, r *http.Request) {
	urls := c.data.Config.URLS
	if r.URL.Path == urls.Root || r.URL.Path == urls.Root+"/" {
		c.getPosts(w, r)
		return
	}

	root := http.Dir(a.path.DataRaws)
	if !utils.FileExists(filepath.Join(a.path.DataRaws, r.URL.Path)) {
		a.renderStatusCode(w, http.StatusNotFound)
		return
	}
	prefix := a.data.URLS.Root + "/"
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)

}
