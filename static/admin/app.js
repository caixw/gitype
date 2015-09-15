"use strict";

// app.js typing的路由配置和公用函数
// @copyright 2015 by caixw
// @link https://github.com/caixw/typing

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

    $('body').load('./login.html');
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
    var suffix = 'typing';
    if (!title){
        title = suffix
    }else{
        title = title + '-' + suffix;
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
    $('.message-row').prepend(div);
    div.hide()
    div.slideDown();

    window.setTimeout(function() {
        div.slideUp(function(){div.remove();});
    }, 5000);
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
        "metas/tags":
                function(){ loadPage('metas-tags.html'); },
    };

    var router = Router(routes).init();
});
