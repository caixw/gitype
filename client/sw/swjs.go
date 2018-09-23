// 请勿修改此文件

package sw

var swjs = []byte(`// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

'use strict';

const versions = new Map();

{{VERSIONS}}

self.addEventListener('install', (event)=>{
    log.info('sw install');
    event.waitUntil(onInstall());
});

self.addEventListener('fetch', (event)=>{
    log.info('sw fetch');
    event.respondWith(onFetch(event));
});

self.addEventListener('activate', (event)=>{
    log.info('sw activate');
    event.waitUntil(onActivate());
})

function onInstall() {
    const ps = [];
    for (let [key, vals] of versions) {
        const p = caches.open(key).then((cache)=>{
            return cache.addAll(vals);
        })

        ps.push(p);
    }

    return Promise.all(ps);
}

function onFetch(event) {
    caches.match(event.request).then((resp)=>{
        return resp || fetch(e.request);
    }).catch((err)=>{
        console.log('sw fetch error:', err);
        return fetch(e.request);
    })
}

function onActivate() {
    caches.keys().then((cachesName)=>{
        const ps = [];

        cachesName.forEach((name)=>{
            if (!cacheExists(name)) {
                const p = caches.delete(name);
                ps.push(p);
            }
        })
        return Promise.all(ps);
    });
}

function cacheExists(name) {
    for(let [key, val] of versions) {
        if (key == name) {
            return true;
        }
    }

    return false;
}
`)
