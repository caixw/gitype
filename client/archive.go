// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"sort"
	"time"

	"github.com/caixw/typing/data"
)

type archive struct {
	date  int64
	Title string
	Posts []*data.Post
}

func (client *Client) initArchives() error {
	archives := make([]*archive, 0, 10)

	for _, post := range client.data.Posts {
		t := time.Unix(post.Created, 0)
		var date int64

		switch client.data.Config.Archive.Type {
		case data.ArchiveTypeMonth:
			date = time.Date(t.Year(), t.Month(), 2, 0, 0, 0, 0, t.Location()).Unix()
		case data.ArchiveTypeYear:
			date = time.Date(t.Year(), 2, 0, 0, 0, 0, 0, t.Location()).Unix()
		}

		found := false
		for _, archive := range archives {
			if archive.date == date {
				archive.Posts = append(archive.Posts, post)
				found = true
				break
			}
		}
		if !found {
			archives = append(archives, &archive{
				date:  date,
				Title: time.Unix(date, 0).Format(client.data.Config.Archive.Format),
				Posts: []*data.Post{post},
			})
		}
	} // end for

	sort.SliceStable(archives, func(i, j int) bool {
		return archives[i].date > archives[j].date
	})

	client.archives = archives

	return nil
}
