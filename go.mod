module github.com/caixw/gitype

require (
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/issue9/assert v1.0.0
	github.com/issue9/is v1.0.0
	github.com/issue9/logs v1.0.0
	github.com/issue9/mux v1.0.0
	github.com/issue9/utils v1.0.0
	github.com/issue9/version v1.0.0
	github.com/issue9/web v0.16.2
	golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33 // indirect
	golang.org/x/text v0.3.0
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/yaml.v2 v2.2.1
)

replace (
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180905080454-ebe1bf3edb33
	golang.org/x/text => github.com/golang/text v0.3.0
)
