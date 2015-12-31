"use strict";

// 根据与页面顶部的距离，控制是否显示top按钮。
$(window).on('scroll', function(){
    var button = $('#return-top');
    if($(document).scrollTop() > 30){
        button.fadeIn();
    }else{
        button.fadeOut();
    }
}).trigger('scroll'); // end $(window).onscroll


/***************** 根据不同的屏宽调整显示内容 ********************/

var result = window.matchMedia("(max-width:600px)");
result.addListener(sizeChange);

sizeChange(result)

// 根据r是否匹配，控制具体的显示规则
function sizeChange(r){
    if (r.matches){
        $('#sidebar .topbar .title a').before('<span class="menu-toggle">&#9776;</span>');

        $('#sidebar .topbar .menu-toggle').on('click', function(){
            $('#sidebar .menus').fadeToggle();
        });
    }else{
        $('#sidebar .topbar.menu-toggle').remove();
    }
}
