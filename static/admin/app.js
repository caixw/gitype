"use strict";

// app.js typing的路由配置和公用函数
// @copyright 2015 by caixw
// @link https://github.com/caixw/typing

// BUG(caixw): 每个页面必须得有一个pageInit函数，否则会调用上一个页面的pageInit函数。

// api的统一前缀
var adminAPIPrefix = '/admin/api';
var frontAPIPrefix = '/api';
var themeURLPrefix = '/themes';

// 每页加载的记录数量
var size = 25;

// 标题设置项
var titleSuffix = 'typing';
var titleSeparator = '-';

// 提示信息关闭时间
var messageTimeout = 5000;

// 将tpl模板的内容解析并插入到container中，data为传递给模板的数据。
function loadTemplate(containerSelector, templateSelector, data) {
    var source = $(templateSelector).html();
    var tpl = Handlebars.compile(source);
    $(containerSelector).append(tpl(data));
}

// 加载指定的模板页面，该页页会自动包含菜单等内容。
// 若未登录，是会自动跳转到登录页面。
// 除了模板名称之外，还可参传递其它任何参数给loadPage，一般为路由匹配项上的参数。
// 若存在这些参数，则会尝试调用加载页面的pageInit函数来做一些初始化。
function loadPage(template) {
    if (!isLogin()){
        redirect('login')
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
}

// 加载登录页面，若已经登录，则会跳转到dashboard页面。
function loadLoginPage() {
    if (isLogin()){
        redirect('dashboard');
        return
    }

    $('body').load('./login.html', function(){
        if (typeof(pageInit) == 'function'){
            pageInit.apply(null);
        }
    });
}

// 判断是否已经登录。
function isLogin() {
    return !!window.sessionStorage.token;
}

// 跳转到指定页页。
function redirect(fregment) {
    if(fregment.charAt(0) != '#'){
        fregment = '#' + fregment;
    }
    window.location.href = fregment;
}

// 设置当前页面的标题
function setTitle(title) {
    if (!title){
        title = titleSuffix
    }else{
        title = title + titleSeparator + titleSuffix;
    }
    $('html head title').html(title);
}

// form中所有包含name的任意元素将被构建成一个JSON
function buildJSON(form) {
    var ret = {};
    form.find('*[name]').each(function(index, elem){
        ret[$(elem).attr('name')] = $(elem).val();
    });
    return JSON.stringify(ret);
}

function del(settings) {
    settings.method  = 'DELETE';
    settings.headers = {'Authorization': window.sessionStorage.token};
    return $.ajax(settings);
}

function post(settings) {
    settings.method  = 'POST';
    settings.headers = {'Authorization': window.sessionStorage.token};
    return $.ajax(settings);
}

function put(settings) {
    settings.method  = 'PUT';
    settings.headers = {'Authorization': window.sessionStorage.token};
    return $.ajax(settings);
}

function get(settings) {
    settings.method  = 'GET';
    settings.headers = {'Authorization': window.sessionStorage.token};
    return $.ajax(settings);
}

function patch(settings) {
    settings.method  = 'PATCH';
    settings.headers = {'Authorization': window.sessionStorage.token};
    return $.ajax(settings);
}

// 向elem元素输出color颜色的信息
function showMessage(color, message) {
    var div = $('<div id="message" class="ui '+color+' visible message">'+message+'</div>');
    div.hide();
    $('.message-row').prepend(div);
    div.slideDown();

    window.setTimeout(function() {
        div.slideUp(function(){div.remove();});
    }, messageTimeout);
}

$(document).ready(function() {
    $.ajaxSetup({
        'contentType': 'application/json;charset=utf-8', // 请求提交类型
        'dataType':    'json',                           // 服务器返回类型
    });

    var routes = {
        "":
                function(){ loadLoginPage(); },
        "login":
                function(){ loadLoginPage(); },
        "logout":
                function(){ $('body').load('logout.html'); },
        "dashboard":
                function(){ loadPage('dashboard.html'); },
        "settings/system":
                function(){ loadPage('settings-system.html'); },
        "settings/users":
                function(){ loadPage('settings-users.html'); },
        "settings/themes":
                function(){ loadPage('settings-themes.html'); },
        "settings/sitemap":
                function(){ loadPage('settings-sitemap.html'); },
        "metas/tags":
                function(){ loadPage('metas-tags.html'); },
        "metas/cats":
                function(){ loadPage('metas-cats.html'); },
        "posts/list":
                function(){ loadPage('posts-list.html'); },
        "posts/edit/:id":
                function(){ loadPage('posts-edit.html'); },
        "comments/list":
                function(){ loadPage('comments-list.html'); },
    };

    var router = Router(routes).init();
});
