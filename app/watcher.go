// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caixw/gitype/path"
	"github.com/issue9/logs"
	fsnotify "gopkg.in/fsnotify.v1"
)

// 初始化一个文件监视器
func (a *app) initWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	paths, err := recursivePaths(a.path)
	if err != nil {
		watcher.Close()
		return nil, err
	}

	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			watcher.Close()
			return nil, err
		}
	}

	return watcher, nil
}

// 递归查找 path.DataDir 每个目录下的子目录。
func recursivePaths(path *path.Path) ([]string, error) {
	ret := []string{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() && strings.Index(path, "/.") < 0 {
			ret = append(ret, path)
		}
		return nil
	}

	if err := filepath.Walk(path.DataDir, walk); err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *app) watch(watcher *fsnotify.Watcher) {
	go func() {
		for {
			fmt.Println("t")
			select {
			case event := <-watcher.Events:
				fmt.Println("1")
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					logs.Debug("watcher.Events:忽略 CHMOD 事件:", event)
					continue
				}

				if time.Now().Sub(a.client.Created()) <= 1*time.Second { // 已经记录
					logs.Debug("watcher.Events:更新太频繁，该监控事件被忽略:", event)
					continue
				}

				logs.Debug("watcher.Events:触发加载事件:", event)

				go func() {
					if err := a.reload(); err != nil {
						logs.Error(err) // 异步事件，直接输出错误日志
					}
				}()
			case err := <-watcher.Errors:
				fmt.Println("2")
				logs.Error(err)
				return // 出错就结束
			} // end select
		} // end for
	}()
}
