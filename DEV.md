# dev

### 代码中的一些常见修改

typing的配置文件仅提供了一些基本配置，若要达到完全的自定义，需要修改代码，
以下列出了一些常见的功能修改。

##### 修改URL结构

所有URL结构的调整，都可以通过修改core/url.go中的相关函数完成。


##### 修改密码加密算法

密码加密算法默认为md5，可通过修改core.HashPassword。


##### 修改配置文件路径

配置文件路径包括app.json和logs.xml两个文件，分别对应core包的LogConfigPath和ConfigPath两个常量。


##### 禁用feed相关内容

在main.go中删除与feed包相关的代码即可。




### 模板制作
