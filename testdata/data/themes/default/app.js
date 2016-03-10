"use strict";

// 根据与页面顶部的距离，控制是否显示top按钮。
$(window).on('scroll', function(){
    var button = $('#top');
    if($(document).scrollTop() > 30){
        button.fadeIn();
    }else{
        button.fadeOut();
    }
}).trigger('scroll'); // end $(window).onscroll
