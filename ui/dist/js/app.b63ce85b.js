(function(e){function n(n){for(var r,c,a=n[0],i=n[1],l=n[2],f=0,s=[];f<a.length;f++)c=a[f],o[c]&&s.push(o[c][0]),o[c]=0;for(r in i)Object.prototype.hasOwnProperty.call(i,r)&&(e[r]=i[r]);p&&p(n);while(s.length)s.shift()();return u.push.apply(u,l||[]),t()}function t(){for(var e,n=0;n<u.length;n++){for(var t=u[n],r=!0,c=1;c<t.length;c++){var i=t[c];0!==o[i]&&(r=!1)}r&&(u.splice(n--,1),e=a(a.s=t[0]))}return e}var r={},o={app:0},u=[];function c(e){return a.p+"js/"+({}[e]||e)+"."+{"chunk-3131":"d0d758f4","chunk-1d29":"2cb6bd59",d353:"9d5db562","chunk-4e8d":"cf0c5c64"}[e]+".js"}function a(n){if(r[n])return r[n].exports;var t=r[n]={i:n,l:!1,exports:{}};return e[n].call(t.exports,t,t.exports,a),t.l=!0,t.exports}a.e=function(e){var n=[],t=o[e];if(0!==t)if(t)n.push(t[2]);else{var r=new Promise(function(n,r){t=o[e]=[n,r]});n.push(t[2]=r);var u,i=document.getElementsByTagName("head")[0],l=document.createElement("script");l.charset="utf-8",l.timeout=120,a.nc&&l.setAttribute("nonce",a.nc),l.src=c(e),u=function(n){l.onerror=l.onload=null,clearTimeout(f);var t=o[e];if(0!==t){if(t){var r=n&&("load"===n.type?"missing":n.type),u=n&&n.target&&n.target.src,c=new Error("Loading chunk "+e+" failed.\n("+r+": "+u+")");c.type=r,c.request=u,t[1](c)}o[e]=void 0}};var f=setTimeout(function(){u({type:"timeout",target:l})},12e4);l.onerror=l.onload=u,i.appendChild(l)}return Promise.all(n)},a.m=e,a.c=r,a.d=function(e,n,t){a.o(e,n)||Object.defineProperty(e,n,{enumerable:!0,get:t})},a.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},a.t=function(e,n){if(1&n&&(e=a(e)),8&n)return e;if(4&n&&"object"===typeof e&&e&&e.__esModule)return e;var t=Object.create(null);if(a.r(t),Object.defineProperty(t,"default",{enumerable:!0,value:e}),2&n&&"string"!=typeof e)for(var r in e)a.d(t,r,function(n){return e[n]}.bind(null,r));return t},a.n=function(e){var n=e&&e.__esModule?function(){return e["default"]}:function(){return e};return a.d(n,"a",n),n},a.o=function(e,n){return Object.prototype.hasOwnProperty.call(e,n)},a.p="/",a.oe=function(e){throw console.error(e),e};var i=window["webpackJsonp"]=window["webpackJsonp"]||[],l=i.push.bind(i);i.push=n,i=i.slice();for(var f=0;f<i.length;f++)n(i[f]);var p=l;u.push([0,"chunk-vendors"]),t()})({0:function(e,n,t){e.exports=t("56d7")},"106f":function(e,n,t){},"56d7":function(e,n,t){"use strict";t.r(n);t("cadf"),t("551c"),t("f466"),t("579f"),t("587a");var r=t("a026"),o=t("9f7b"),u=function(){var e=this,n=e.$createElement,t=e._self._c||n;return t("router-view")},c=[],a={name:"app"},i=a,l=(t("5c0b"),t("2877")),f=Object(l["a"])(i,u,c,!1,null,null,null),p=f.exports,s=t("a18c");r["default"].use(o["a"]),new r["default"]({el:"#app",router:s["a"],template:"<App/>",components:{App:p}})},"5c0b":function(e,n,t){"use strict";var r=t("106f"),o=t.n(r);o.a},a18c:function(e,n,t){"use strict";t("cadf"),t("551c");var r=t("a026"),o=t("8c4f"),u=function(){return t.e("chunk-4e8d").then(t.bind(null,"e8c5"))},c=function(){return Promise.all([t.e("chunk-3131"),t.e("chunk-1d29")]).then(t.bind(null,"aa72"))},a=function(){return Promise.all([t.e("chunk-3131"),t.e("d353")]).then(t.bind(null,"d353"))};r["default"].use(o["a"]),n["a"]=new o["a"]({mode:"hash",linkActiveClass:"open active",scrollBehavior:function(){return{y:0}},routes:[{path:"/",redirect:"/devices",name:"Home",component:u,props:!0,children:[{path:"devices",name:"Devices",component:c,props:!0},{path:"newdevice",name:"Newdevice",component:a,props:!0}]}]})}});
//# sourceMappingURL=app.b63ce85b.js.map