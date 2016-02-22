"use strict";System.register(["underscore","./view","./collection"],function(t,e){function n(t,e){if(!(t instanceof e))throw new TypeError("Cannot call a class as a function")}var i,o,r,s,a;return{setters:[function(t){i=t.extend},function(t){o=t["default"]},function(t){r=t["default"]}],execute:function(){s=function(){function t(t,e){for(var n=0;n<e.length;n++){var i=e[n];i.enumerable=i.enumerable||!1,i.configurable=!0,"value"in i&&(i.writable=!0),Object.defineProperty(t,i.key,i)}}return function(e,n,i){return n&&t(e.prototype,n),i&&t(e,i),e}}(),a=function(){function t(){var e=arguments.length<=0||void 0===arguments[0]?{}:arguments[0];n(this,t),this.attrs=e,this.views=[],this.changeHooks={}}return s(t,[{key:"get",value:function(t){return this.attrs[t]}},{key:"set",value:function(t,e){this.attrs[t]=e,this.execChangeHooks(t)}},{key:"setAttrs",value:function(t){i(this.attrs,t);for(var e in t)this.execChangeHooks(e)}},{key:"append",value:function(t,e){this.attrs[t]?this.attrs[t].push(e):this.attrs[t]=[e],this.execChangeHooks(t)}},{key:"extend",value:function(t,e){this.attrs[t]?i(this.attrs[t],e):this.attrs[t]=e,this.execChangeHooks(t)}},{key:"attach",value:function(t){this.views.push(t)}},{key:"detach",value:function(t){this.views.splice(this.views.indexOf(t),1)}},{key:"remove",value:function(){this.collection&&this.collection.remove(this);var t=!0,e=!1,n=void 0;try{for(var i,o=this.views[Symbol.iterator]();!(t=(i=o.next()).done);t=!0){var r=i.value;r.remove()}}catch(s){e=!0,n=s}finally{try{!t&&o["return"]&&o["return"]()}finally{if(e)throw n}}}},{key:"onChange",value:function(t,e){this.changeHooks[t]?this.changeHooks[t].push(e):this.changeHooks[t]=[e]}},{key:"execChangeHooks",value:function(t){if(this.changeHooks[t]){var e=this.get(t),n=!0,i=!1,o=void 0;try{for(var r,s=this.changeHooks[t][Symbol.iterator]();!(n=(r=s.next()).done);n=!0){var a=r.value;a(e)}}catch(u){i=!0,o=u}finally{try{!n&&s["return"]&&s["return"]()}finally{if(i)throw o}}}}}]),t}(),t("default",a)}}});
//# sourceMappingURL=../maps/client/model.js.map