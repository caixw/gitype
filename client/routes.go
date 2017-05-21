// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import "net/http"

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
	//
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

func (c *Client) getMedia(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getRaws(w http.ResponseWriter, r *http.Request) {
	//
}
