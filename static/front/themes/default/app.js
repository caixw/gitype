"use strict";

// 控制top按钮的显示
window.addEventListener('scroll', function(){
    var button = document.getElementById('return-top');
    var t = document.documentElement.scrollTop || document.body.scrollTop;
    if(t>30){
        button.style.display='block';
    }else{
        button.style.display='none';
    }
}, false);
