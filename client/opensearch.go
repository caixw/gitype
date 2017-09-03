// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"

	"github.com/caixw/typing/vars"
)

func (client *Client) initOpensearch() error {
	if client.data.Config.Opensearch == nil {
		return nil
	}

	if err := client.buildOpensearch(); err != nil {
		return err
	}

	conf := client.data.Config
	client.patterns = append(client.patterns, conf.Opensearch.URL)
	client.mux.GetFunc(conf.Opensearch.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, conf.Opensearch.Type)
		w.Write(client.opensearch)
	}))

	return nil
}

// 用于生成一个符合 atom 规范的 XML 文本。
func (client *Client) buildOpensearch() error {
	w := newWrite()
	o := client.data.Config.Opensearch

	w.WriteStartElement("OpenSearchDescription", map[string]string{
		"xmlns": "http://a9.com/-/spec/opensearch/1.1/",
	})

	w.WriteElement("InputEncoding", "UTF-8", nil)
	w.WriteElement("OutputEncoding", "UTF-8", nil)
	w.WriteElement("ShortName", o.ShortName, nil)
	w.WriteElement("Description", o.Description, nil)

	if len(o.LongName) > 0 {
		w.WriteElement("LongName", o.LongName, nil)
	}

	if o.Image != nil {
		w.WriteElement("Image", o.Image.URL, map[string]string{
			"type": o.Image.Type,
		})
	}

	w.WriteCloseElement("Url", map[string]string{
		"type":     client.data.Config.Type,
		"template": vars.SearchURL("{searchTerms}", 0),
	})

	w.WriteElement("Developer", vars.AppName, nil)
	w.WriteElement("Language", client.data.Config.Language, nil)

	w.WriteEndElement("OpenSearchDescription")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	client.opensearch = bs

	return nil
}
