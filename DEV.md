# dev

### 代码中的一些常见修改

typing的配置文件仅提供了一些基本配置，若要达到完全的自定义，需要修改代码，
以下列出了一些常见的功能修改。


##### 修改URL结构

所有URL结构的调整，都可以通过修改core/url.go中的相关函数完成。


##### 修改密码加密算法

密码加密算法默认为md5，可通过修改core.HashPassword。


##### 修改配置文件路径

配置文件路径包括app.json和logs.xml两个文件，分别对应app包的LogConfigPath和ConfigPath两个常量。


##### 禁用feed相关内容

在main.go中删除与feed包相关的代码即可。


##### 修正atom, sitemap, rss等文件名称

修正feed包中的相关常量名称。


##### 自定义模板

所有的模板都在static目录下，其中static/admin对应后的模板；
static/front/themes下对应的是前台的各个主题。


##### 默认的安装数据

配置文件在app/install.go中；
数据库定义及大部分的默认数据models/install.go中；
配置项在options/install.go中。



### 目录结构

```
|--- admin 后台的逻辑处理代码
|
|--- app 处理启动时需要初始化的内容
|
|--- options 与数据库加载的配置内容的相关功能
|
|--- models 模块定义文件
|
|--- util 被其它包引用的核心代码
|
|--- feed 与sitemap, rss, atom相关的代码
|
|--- static 静态文件，包括后台的模板和所有的主题
|      |
|      |--- admin 后台的静态代码
|      |
|      |--- front 网站的前端页面内容
|             |
|             |--- themes 主题模板所在的目录
|
|--- themes 主题的逻辑处理代码
```


### 模板制作
