"use strict";

// app.js typing的路由配置和公用函数
// @copyright 2015 by caixw
// @link https://github.com/caixw/typing

// App 封装了后台的一些公共操作函数。
// 通过options参数可更改以下内容，除非有特殊需求，否则无需更改这些内容：
// - titleSuffix     标题后缀
// - titleSeparator  标题分隔符
// - messageTimeout  提示信息关闭时间
// - adminAPIPrefix  后台操作api地址的前缀
// - frontAPIPrefix  前端操作api地址的前缀
// - uploadURLPrefix 上传文件的根目录
function App(options) {
    var defaults = {
        titleSuffix:    'typing',
        titleSeparator: '-',
        adminAPIPrefix: '/admin/api',
        frontAPIPrefix: '/api',
        uploadURLPrefix:'/uploads',
        messageTimeout: 5000
    };
    var opt = $.extend({}, defaults, options);
    var self = this;

    // 设置标题，若值为空，则只显示opt.titleSuffix。
    this.setTitle = function(title) {
        if (!title){
            title = opt.titleSuffix
        }else{
            title = title + opt.titleSeparator + opt.titleSuffix;
        }
        $('html>head>title').html(title);
    };

    this.frontAPI = function(url) {
        return opt.frontAPIPrefix + url;
    };

    this.adminAPI = function(url) {
        return opt.adminAPIPrefix + url;
    };

    this.uploadURL = function(url) {
        return opt.uploadURLPrefix + url;
    }

    // 执行一个ajax操作，提交和返回数据均为json
    function ajax(settings) {
        settings.contentType = 'application/json;charset=utf-8';
        settings.dataType    = 'json';
        settings.headers     = {'Authorization': window.sessionStorage.token};

        return $.ajax(settings).fail(function(jqXHR, textStatus, errorThrown){
            if (jqXHR.status == 401) {
                window.sessionStorage.token = '';
                self.redirect('login');
                return;
            }

            var msg = '访问资源<'+settings.url+'>时发生以下错误：'+jqXHR.status;
            self.showMessage('red', msg);
        });
    }

    // 执行一条delete的restful api操作。
    this.delete = function(settings) {
        settings.method = 'DELETE';
        return ajax(settings);
    };

    // 执行一条post的restful api操作。
    this.post = function(settings) {
        settings.method = 'POST';
        return ajax(settings);
    };

    // 执行一条put的restful api操作。
    this.put = function(settings) {
        settings.method = 'PUT';
        return ajax(settings);
    };

    // 执行一条get的restful api操作。
    this.get = function(settings) {
        settings.method = 'GET';
        return ajax(settings);
    };

    // 执行一条patch的restful api操作。
    this.patch = function(settings) {
        settings.method = 'PATCH';
        return ajax(settings);
    };

    // 在页面右下解显示一提示信息，该信息会自动消失。
    this.showMessage = function(color, message) {
        var div = $('<div id="message" class="ui '+color+' visible message">'+message+'</div>');
        div.hide();
        $('.message-row').prepend(div);
        div.slideDown();

        window.setTimeout(function() {
            div.slideUp(function(){div.remove();});
        }, opt.messageTimeout);
    };

    // 将tpl模板的内容解析并插入到container中，data为传递给模板的数据。
    this.loadTemplate = function(containerSelector, templateSelector, data) {
        var source = $(templateSelector).html();
        var tpl = Handlebars.compile(source);
        $(containerSelector).append(tpl(data));
    };

    // 加载指定的模板页面，该页面会自动包含菜单等内容。
    // 若未登录，是会自动跳转到登录页面。
    // template用于指定需要加载的模板，如果还有额外的其它参数，将会被传递给页面的pageInit函数。
    this.loadBodyPage = function(template) {
        if (!self.isLogin()){
            self.redirect('login')
            return
        }

        var args = [];
        Array.prototype.push.apply(args, arguments);
        args.shift();

        $('body').load('./body.html', function(){
            $('#content').load(template, function(){
                if (typeof(pageInit) == 'function'){
                    pageInit.apply(null, args);
                }
            });
        });
    };

    // 加载指定的模板到body元素中。
    // template用于指定需要加载的模板，如果还有额外的其它参数，将会被传递给页面的pageInit函数。
    this.loadPage = function(template) {
        var args = [];
        Array.prototype.push.apply(args, arguments);
        args.shift();

        $('body').load(template, function(){
            if (typeof(pageInit) == 'function'){
                pageInit.apply(null, args);
            }
        });
    };

    // 加载登录页面，若已经登录，则会跳转到dashboard页面。
    this.loadLoginPage = function() {
        if (self.isLogin()){
            self.redirect('dashboard');
            return
        }

        self.loadPage('login.html');
    };


    // 判断是否已经登录。
    this.isLogin = function() {
        return !!window.sessionStorage.token;
    };

    // 跳转到指定页页。
    this.redirect = function(fregment) {
        if(fregment.charAt(0) != '#'){
            fregment = '#' + fregment;
        }
        window.location.href = fregment;
    };

    // form中所有包含name的任意元素将被构建成一个Object
    this.buildObject = function(form) {
        var obj = {};
        form.find('*[name]').each(function(index, elem){
            var v = $(elem).val();
            var type = $(elem).attr('data-type');
            switch(type){
            case 'int':
                v = parseInt(v);
                break;
            case 'float':
                v = parseFloat(v);
                break;
            case 'bool':
                v = $(elem).prop('checked');
                break;
            }
            obj[$(elem).attr('name')] = v;
        });
        return obj
    };

    // form中所有包含name的任意元素将被构建成一个JSON
    this.buildJSON = function(form) {
        return JSON.stringify(this.buildObject(form));
    };

    // 开始监听路由。
    this.listen = function(routes) {
        var router = Router(routes).init();
    };

    // 执行单个文件上传操作
    this.upload = function(mediaSelector){
        var files = $(mediaSelector)[0].files;
        var formdata = new FormData();
        formdata.append('media', files[0]);

        return $.ajax({
            'type':        'POST',
            'data':        formdata,
            'url':         self.adminAPI('/media'),
            'contentType': false,
            'processData': false,
            'headers':    {'Authorization': window.sessionStorage.token}
        }).done(function(){
            self.showMessage('green', '上传文件成功');
        }).fail(function(jqXHR, textStatus, errorThrown){
            self.showMessage('red', '上传文件失败:'+errorThrown);
        });
    } // end upload
} // end App
