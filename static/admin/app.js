"use strict";

// app.js typing的路由配置和公用函数
// @copyright 2015 by caixw
// @link https://github.com/caixw/typing

// App 封装了后台的一些公共操作函数。
// 通过options参数可更改以下内容：
// - titleSuffix    标题后缀
// - titleSeparator 标题分隔符
// - messageTimeout 提示信息关闭时间
// - adminAPIPrefix 后台操作api地址的前缀
// - frontAPIPrefix 前端操作api地址的前缀
function App(options) {
    var defaults = {
        titleSuffix:    'typing',
        titleSeparator: '-',
        adminAPIPrefix: '/admin/api',
        frontAPIPrefix: '/api',
        messageTimeout: 5000
    };
    var opt = $.extend({}, defaults, options);

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

    // 执行一条delete的restful api操作。
    this.delete = function(settings) {
        settings.contentType = 'application/json;charset=utf-8';
        settings.dataType    = 'json';
        settings.method      = 'DELETE';
        settings.headers     = {'Authorization': window.sessionStorage.token};
        return $.ajax(settings);
    };

    // 执行一条post的restful api操作。
    this.post = function(settings) {
        settings.contentType = 'application/json;charset=utf-8';
        settings.dataType    = 'json';
        settings.method      = 'POST';
        settings.headers     = {'Authorization': window.sessionStorage.token};
        return $.ajax(settings);
    };

    // 执行一条put的restful api操作。
    this.put = function(settings) {
        settings.contentType = 'application/json;charset=utf-8';
        settings.dataType    = 'json';
        settings.method      = 'PUT';
        settings.headers     = {'Authorization': window.sessionStorage.token};
        return $.ajax(settings);
    };

    // 执行一条get的restful api操作。
    this.get = function(settings) {
        settings.contentType = 'application/json;charset=utf-8';
        settings.dataType    = 'json';
        settings.method      = 'GET';
        settings.headers     = {'Authorization': window.sessionStorage.token};
        return $.ajax(settings);
    };

    // 执行一条patch的restful api操作。
    this.patch = function(settings) {
        settings.contentType = 'application/json;charset=utf-8';
        settings.dataType    = 'json';
        settings.method      = 'PATCH';
        settings.headers     = {'Authorization': window.sessionStorage.token};
        return $.ajax(settings);
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
    // 除了模板名称之外，还可参传递其它任何参数给loadPage，一般为路由匹配项上的参数。
    // 若存在这些参数，则会尝试调用加载页面的pageInit函数来做一些初始化。
    this.loadBodyPage = function(template) {
        if (!this.isLogin()){
            this.redirect('login')
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
        if (this.isLogin()){
            this.redirect('dashboard');
            return
        }

        this.loadPage('login.html');
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

    // form中所有包含name的任意元素将被构建成一个JSON
    this.buildJSON = function(form) {
        var ret = {};
        form.find('*[name]').each(function(index, elem){
            var v = $(elem).val();
            // TODO 考虑其它类型的情况
            if ($(elem).attr('type') == 'number'){
                v = parseInt(v);
            }
            ret[$(elem).attr('name')] = v;
        });
        return JSON.stringify(ret);
    };

    // 开始监听路由。
    this.listen = function(routes) {
        var router = Router(routes).init();
    };
} // end App
