"use strict";function _typeof(t){return(_typeof="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(t){return typeof t}:function(t){return t&&"function"==typeof Symbol&&t.constructor===Symbol&&t!==Symbol.prototype?"symbol":typeof t})(t)}function _typeof(t){return(_typeof="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(t){return typeof t}:function(t){return t&&"function"==typeof Symbol&&t.constructor===Symbol&&t!==Symbol.prototype?"symbol":typeof t})(t)}!function(t){"function"==typeof define&&define.amd?define(t):"undefined"!=typeof module&&void 0!==module.exports?module.exports=t():"undefined"!=typeof Package?Sortable=t():window.Sortable=t()}(function(){var w,C,E,I,N,O,f,p,B,q,L,d,i,R,l,a,A,t,m={},o=/\s+/g,k="Sortable"+(new Date).getTime(),g=window,c=g.document,r=g.parseInt,s=!!("draggable"in c.createElement("div")),h=((t=c.createElement("x")).style.cssText="pointer-events:auto","auto"===t.style.pointerEvents),$=!1,v=Math.abs,u=[],j=e(function(t,e,n){if(n&&e.scroll){var i,o,a,r,s=e.scrollSensitivity,l=e.scrollSpeed,d=t.clientX,c=t.clientY,h=window.innerWidth,u=window.innerHeight;if(p!==n&&(f=e.scroll,p=n,!0===f)){f=n;do{if(f.offsetWidth<f.scrollWidth||f.offsetHeight<f.scrollHeight)break}while(f=f.parentNode)}f&&(o=(i=f).getBoundingClientRect(),a=(v(o.right-d)<=s)-(v(o.left-d)<=s),r=(v(o.bottom-c)<=s)-(v(o.top-c)<=s)),a||r||(r=(u-c<=s)-(c<=s),((a=(h-d<=s)-(d<=s))||r)&&(i=g)),m.vx===a&&m.vy===r&&m.el===i||(m.el=i,m.vx=a,m.vy=r,clearInterval(m.pid),i&&(m.pid=setInterval(function(){i===g?g.scrollTo(g.pageXOffset+a*l,g.pageYOffset+r*l):(r&&(i.scrollTop+=r*l),a&&(i.scrollLeft+=a*l))},24)))}},30),b=function(t){var e=t.group;e&&"object"==_typeof(e)||(e=t.group={name:e}),["pull","put"].forEach(function(t){t in e||(e[t]=!0)}),t.groups=" "+e.name+(e.put.join?" "+e.put.join(" "):"")+" "};function _(t,e){if(!t||!t.nodeType||1!==t.nodeType)throw"Sortable: `el` must be HTMLElement, and not "+{}.toString.call(t);this.el=t,this.options=e=V({},e),t[k]=this;var n={group:Math.random(),sort:!0,disabled:!1,store:null,handle:null,scroll:!0,scrollSensitivity:30,scrollSpeed:10,draggable:/[uo]l/i.test(t.nodeName)?"li":">*",ghostClass:"sortable-ghost",chosenClass:"sortable-chosen",ignore:"a, img",filter:null,animation:0,setData:function(t,e){t.setData("Text",e.textContent)},dropBubble:!1,dragoverBubble:!1,dataIdAttr:"data-id",delay:0,forceFallback:!1,fallbackClass:"sortable-fallback",fallbackOnBody:!1};for(var i in n)!(i in e)&&(e[i]=n[i]);for(var o in b(e),this)"_"===o.charAt(0)&&(this[o]=this[o].bind(this));this.nativeDraggable=!e.forceFallback&&s,y(t,"mousedown",this._onTapStart),y(t,"touchstart",this._onTapStart),this.nativeDraggable&&(y(t,"dragover",this),y(t,"dragenter",this)),u.push(this._onDragOver),e.store&&this.sort(e.store.get(this))}function M(t){I&&I.state!==t&&(Y(I,"display",t?"none":""),!t&&I.state&&N.insertBefore(I,w),I.state=t)}function X(t,e,n){if(t){n=n||c;var i=(e=e.split(".")).shift().toUpperCase(),o=new RegExp("\\s("+e.join("|")+")(?=\\s)","g");do{if(">*"===i&&t.parentNode===n||(""===i||t.nodeName.toUpperCase()==i)&&(!e.length||((" "+t.className+" ").match(o)||[]).length==e.length))return t}while(t!==n&&(t=t.parentNode))}return null}function y(t,e,n){t.addEventListener(e,n,!1)}function D(t,e,n){t.removeEventListener(e,n,!1)}function S(t,e,n){if(t)if(t.classList)t.classList[n?"add":"remove"](e);else{var i=(" "+t.className+" ").replace(o," ").replace(" "+e+" "," ");t.className=(i+(n?" "+e:"")).replace(o," ")}}function Y(t,e,n){var i=t&&t.style;if(i){if(void 0===n)return c.defaultView&&c.defaultView.getComputedStyle?n=c.defaultView.getComputedStyle(t,""):t.currentStyle&&(n=t.currentStyle),void 0===e?n:n[e];e in i||(e="-webkit-"+e),i[e]=n+("string"==typeof n?"":"px")}}function T(t,e,n){if(t){var i=t.getElementsByTagName(e),o=0,a=i.length;if(n)for(;o<a;o++)n(i[o],o);return i}return[]}function x(t,e,n,i,o,a,r){var s=c.createEvent("Event"),l=(t||e[k]).options,d="on"+n.charAt(0).toUpperCase()+n.substr(1);s.initEvent(n,!0,!0),s.to=e,s.from=o||e,s.item=i||e,s.clone=I,s.oldIndex=a,s.newIndex=r,e.dispatchEvent(s),l[d]&&l[d].call(t,s)}function F(t,e,n,i,o,a){var r,s,l=t[k],d=l.options.onMove;return(r=c.createEvent("Event")).initEvent("move",!0,!0),r.to=e,r.from=t,r.dragged=n,r.draggedRect=i,r.related=o||e,r.relatedRect=a||e.getBoundingClientRect(),t.dispatchEvent(r),d&&(s=d.call(l,r)),s}function U(t){t.draggable=!1}function H(){$=!1}function P(t){for(var e=t.tagName+t.className+t.src+t.href+t.textContent,n=e.length,i=0;n--;)i+=e.charCodeAt(n);return i.toString(36)}function W(t){var e=0;if(!t||!t.parentNode)return-1;for(;t&&(t=t.previousElementSibling);)"TEMPLATE"!==t.nodeName.toUpperCase()&&e++;return e}function e(t,e){var n,i;return function(){void 0===n&&(n=arguments,i=this,setTimeout(function(){1===n.length?t.call(i,n[0]):t.apply(i,n),n=void 0},e))}}function V(t,e){if(t&&e)for(var n in e)e.hasOwnProperty(n)&&(t[n]=e[n]);return t}return _.prototype={constructor:_,_onTapStart:function(t){var e=this,n=this.el,i=this.options,o=t.type,a=t.touches&&t.touches[0],r=(a||t).target,s=r,l=i.filter;if(!("mousedown"===o&&0!==t.button||i.disabled)&&(r=X(r,i.draggable,n))){if(d=W(r),"function"==typeof l){if(l.call(this,t,r,this))return x(e,s,"filter",r,n,d),void t.preventDefault()}else if(l&&(l=l.split(",").some(function(t){if(t=X(s,t.trim(),n))return x(e,t,"filter",r,n,d),!0})))return void t.preventDefault();i.handle&&!X(s,i.handle,n)||this._prepareDragStart(t,a,r)}},_prepareDragStart:function(t,e,n){var i,o=this,a=o.el,r=o.options,s=a.ownerDocument;n&&!w&&n.parentNode===a&&(l=t,N=a,C=(w=n).parentNode,O=w.nextSibling,R=r.group,i=function(){o._disableDelayedDrag(),w.draggable=!0,S(w,o.options.chosenClass,!0),o._triggerDragStart(e)},r.ignore.split(",").forEach(function(t){T(w,t.trim(),U)}),y(s,"mouseup",o._onDrop),y(s,"touchend",o._onDrop),y(s,"touchcancel",o._onDrop),r.delay?(y(s,"mouseup",o._disableDelayedDrag),y(s,"touchend",o._disableDelayedDrag),y(s,"touchcancel",o._disableDelayedDrag),y(s,"mousemove",o._disableDelayedDrag),y(s,"touchmove",o._disableDelayedDrag),o._dragStartTimer=setTimeout(i,r.delay)):i())},_disableDelayedDrag:function(){var t=this.el.ownerDocument;clearTimeout(this._dragStartTimer),D(t,"mouseup",this._disableDelayedDrag),D(t,"touchend",this._disableDelayedDrag),D(t,"touchcancel",this._disableDelayedDrag),D(t,"mousemove",this._disableDelayedDrag),D(t,"touchmove",this._disableDelayedDrag)},_triggerDragStart:function(t){t?(l={target:w,clientX:t.clientX,clientY:t.clientY},this._onDragStart(l,"touch")):this.nativeDraggable?(y(w,"dragend",this),y(N,"dragstart",this._onDragStart)):this._onDragStart(l,!0);try{c.selection?c.selection.empty():window.getSelection().removeAllRanges()}catch(t){}},_dragStarted:function(){N&&w&&(S(w,this.options.ghostClass,!0),x(_.active=this,N,"start",w,N,d))},_emulateDragOver:function(){if(a){if(this._lastX===a.clientX&&this._lastY===a.clientY)return;this._lastX=a.clientX,this._lastY=a.clientY,h||Y(E,"display","none");var t=c.elementFromPoint(a.clientX,a.clientY),e=t,n=" "+this.options.group.name,i=u.length;if(e)do{if(e[k]&&-1<e[k].options.groups.indexOf(n)){for(;i--;)u[i]({clientX:a.clientX,clientY:a.clientY,target:t,rootEl:e});break}t=e}while(e=e.parentNode);h||Y(E,"display","")}},_onTouchMove:function(t){if(l){_.active||this._dragStarted(),this._appendGhost();var e=t.touches?t.touches[0]:t,n=e.clientX-l.clientX,i=e.clientY-l.clientY,o=t.touches?"translate3d("+n+"px,"+i+"px,0)":"translate("+n+"px,"+i+"px)";A=!0,a=e,Y(E,"webkitTransform",o),Y(E,"mozTransform",o),Y(E,"msTransform",o),Y(E,"transform",o),t.preventDefault()}},_appendGhost:function(){if(!E){var t,e=w.getBoundingClientRect(),n=Y(w),i=this.options;S(E=w.cloneNode(!0),i.ghostClass,!1),S(E,i.fallbackClass,!0),Y(E,"top",e.top-r(n.marginTop,10)),Y(E,"left",e.left-r(n.marginLeft,10)),Y(E,"width",e.width),Y(E,"height",e.height),Y(E,"opacity","0.8"),Y(E,"position","fixed"),Y(E,"zIndex","100000"),Y(E,"pointerEvents","none"),i.fallbackOnBody&&c.body.appendChild(E)||N.appendChild(E),t=E.getBoundingClientRect(),Y(E,"width",2*e.width-t.width),Y(E,"height",2*e.height-t.height)}},_onDragStart:function(t,e){var n=t.dataTransfer,i=this.options;this._offUpEvents(),"clone"==R.pull&&(Y(I=w.cloneNode(!0),"display","none"),N.insertBefore(I,w)),e?("touch"===e?(y(c,"touchmove",this._onTouchMove),y(c,"touchend",this._onDrop),y(c,"touchcancel",this._onDrop)):(y(c,"mousemove",this._onTouchMove),y(c,"mouseup",this._onDrop)),this._loopId=setInterval(this._emulateDragOver,50)):(n&&(n.effectAllowed="move",i.setData&&i.setData.call(this,n,w)),y(c,"drop",this),setTimeout(this._dragStarted,0))},_onDragOver:function(t){var e,n,i,o,a,r,s=this.el,l=this.options,d=l.group,c=d.put,h=R===d,u=l.sort;if(void 0!==t.preventDefault&&(t.preventDefault(),!l.dragoverBubble&&t.stopPropagation()),A=!0,R&&!l.disabled&&(h?u||(i=!N.contains(w)):R.pull&&c&&(R.name===d.name||c.indexOf&&~c.indexOf(R.name)))&&(void 0===t.rootEl||t.rootEl===this.el)){if(j(t,l,this.el),$)return;if(e=X(t.target,l.draggable,s),n=w.getBoundingClientRect(),i)return M(!0),void(I||O?N.insertBefore(w,I||O):u||N.appendChild(w));if(0===s.children.length||s.children[0]===E||s===t.target&&(o=t,a=s.lastElementChild,r=a.getBoundingClientRect(),e=(5<o.clientY-(r.top+r.height)||5<o.clientX-(r.right+r.width))&&a)){if(e){if(e.animated)return;p=e.getBoundingClientRect()}M(h),!1!==F(N,s,w,n,e,p)&&(w.contains(s)||(s.appendChild(w),C=s),this._animate(n,w),e&&this._animate(p,e))}else if(e&&!e.animated&&e!==w&&void 0!==e.parentNode[k]){B!==e&&(q=Y(B=e),L=Y(e.parentNode));var f,p=e.getBoundingClientRect(),m=p.right-p.left,g=p.bottom-p.top,v=/left|right|inline/.test(q.cssFloat+q.display)||"flex"==L.display&&0===L["flex-direction"].indexOf("row"),b=e.offsetWidth>w.offsetWidth,_=e.offsetHeight>w.offsetHeight,y=.5<(v?(t.clientX-p.left)/m:(t.clientY-p.top)/g),D=e.nextElementSibling,S=F(N,s,w,n,e,p);if(!1!==S){if($=!0,setTimeout(H,30),M(h),1===S||-1===S)f=1===S;else if(v){var T=w.offsetTop,x=e.offsetTop;f=T===x?e.previousElementSibling===w&&!b||y&&b:T<x}else f=D!==w&&!_||y&&_;w.contains(s)||(f&&!D?s.appendChild(w):e.parentNode.insertBefore(w,f?D:e)),C=w.parentNode,this._animate(n,w),this._animate(p,e)}}}},_animate:function(t,e){var n=this.options.animation;if(n){var i=e.getBoundingClientRect();Y(e,"transition","none"),Y(e,"transform","translate3d("+(t.left-i.left)+"px,"+(t.top-i.top)+"px,0)"),e.offsetWidth,Y(e,"transition","all "+n+"ms"),Y(e,"transform","translate3d(0,0,0)"),clearTimeout(e.animated),e.animated=setTimeout(function(){Y(e,"transition",""),Y(e,"transform",""),e.animated=!1},n)}},_offUpEvents:function(){var t=this.el.ownerDocument;D(c,"touchmove",this._onTouchMove),D(t,"mouseup",this._onDrop),D(t,"touchend",this._onDrop),D(t,"touchcancel",this._onDrop)},_onDrop:function(t){var e=this.el,n=this.options;clearInterval(this._loopId),clearInterval(m.pid),clearTimeout(this._dragStartTimer),D(c,"mousemove",this._onTouchMove),this.nativeDraggable&&(D(c,"drop",this),D(e,"dragstart",this._onDragStart)),this._offUpEvents(),t&&(A&&(t.preventDefault(),!n.dropBubble&&t.stopPropagation()),E&&E.parentNode.removeChild(E),w&&(this.nativeDraggable&&D(w,"dragend",this),U(w),S(w,this.options.ghostClass,!1),S(w,this.options.chosenClass,!1),N!==C?0<=(i=W(w))&&(x(null,C,"sort",w,N,d,i),x(this,N,"sort",w,N,d,i),x(null,C,"add",w,N,d,i),x(this,N,"remove",w,N,d,i)):(I&&I.parentNode.removeChild(I),w.nextSibling!==O&&0<=(i=W(w))&&(x(this,N,"update",w,N,d,i),x(this,N,"sort",w,N,d,i))),_.active&&(null!==i&&-1!==i||(i=d),x(this,N,"end",w,N,d,i),this.save())),N=w=C=E=O=I=f=p=l=a=A=i=B=q=R=_.active=null)},handleEvent:function(t){var e=t.type;"dragover"===e||"dragenter"===e?w&&(this._onDragOver(t),function(t){t.dataTransfer&&(t.dataTransfer.dropEffect="move");t.preventDefault()}(t)):"drop"!==e&&"dragend"!==e||this._onDrop(t)},toArray:function(){for(var t,e=[],n=this.el.children,i=0,o=n.length,a=this.options;i<o;i++)X(t=n[i],a.draggable,this.el)&&e.push(t.getAttribute(a.dataIdAttr)||P(t));return e},sort:function(t){var i={},o=this.el;this.toArray().forEach(function(t,e){var n=o.children[e];X(n,this.options.draggable,o)&&(i[t]=n)},this),t.forEach(function(t){i[t]&&(o.removeChild(i[t]),o.appendChild(i[t]))})},save:function(){var t=this.options.store;t&&t.set(this)},closest:function(t,e){return X(t,e||this.options.draggable,this.el)},option:function(t,e){var n=this.options;if(void 0===e)return n[t];n[t]=e,"group"===t&&b(n)},destroy:function(){var t=this.el;t[k]=null,D(t,"mousedown",this._onTapStart),D(t,"touchstart",this._onTapStart),this.nativeDraggable&&(D(t,"dragover",this),D(t,"dragenter",this)),Array.prototype.forEach.call(t.querySelectorAll("[draggable]"),function(t){t.removeAttribute("draggable")}),u.splice(u.indexOf(this._onDragOver),1),this._onDrop(),this.el=t=null}},_.utils={on:y,off:D,css:Y,find:T,is:function(t,e){return!!X(t,e,t)},extend:V,throttle:e,closest:X,toggleClass:S,index:W},_.create=function(t,e){return new _(t,e)},_.version="1.4.2",_}),function(t){"function"==typeof define&&define.amd?define(["jquery"],t):"object"===("undefined"==typeof exports?"undefined":_typeof(exports))?t(require("jquery")):t(jQuery)}(function(l){var n=l("body"),o="qor.chooser.sortable",t="enable."+o,e="click."+o,d=".select2-container",c=".select2-search__field",h=".qor-dragable__list",u=".qor-dragable__list-data",a=".qor-dragable__button-add",r=".qor-bottomsheets",s="is_selected",i='[data-select-modal="many_sortable"]';function f(t,e){this.$element=l(t),this.options=l.extend({},f.DEFAULTS,l.isPlainObject(e)&&e),this.init()}return f.prototype={constructor:f,init:function(){var n=this.$element,t=n.data(),e=n.parents(".qor-dragable"),i=n.data("placeholder"),o=this,a={minimumResultsForSearch:8,dropdownParent:n.parent()};this.$selector=e.find(u),this.$sortableList=e.find(h);var r=(this.$parent=e).find(h)[0];if(this.sortable=window.Sortable.create(r,{animation:150,handle:".qor-dragable__list-handle",filter:".qor-dragable__list-delete",dataIdAttr:"data-index",onFilter:function(t){var e=l(t.item);e.remove(),o.removeItemsFromList(e.data())},onUpdate:function(){o.renderOption()}}),t.remoteData){var s=window.getSelect2AjaxDynamicURL;a.ajax=l.fn.select2.ajaxCommonOptions(t),s&&l.isFunction(s)?a.ajax.url=function(){return s(t)}:a.ajax.url=t.remoteUrl,a.templateResult=function(t){var e=n.parents(".qor-field").find('[name="select2-result-template"]');return l.fn.select2.ajaxFormatResult(t,e)},a.templateSelection=function(t){if(t.loading)return t.text;var e=n.parents(".qor-field").find('[name="select2-selection-template"]');return l.fn.select2.ajaxFormatResult(t,e)}}n.is("select")&&(n.on("change",function(){l(c).attr("placeholder",i)}).on("select2:select",function(t){var e=t.params.data;e.value=e.Name||e.text||e.Text||e.Title||e.Code,o.addItems(e)}).on("select2:unselect",function(t){o.removeItems(t.params.data)}),n.select2(a),e.find(d).hide(),l(c).attr("placeholder",i)),this.bind()},bind:function(){this.$parent.on(e,a,this.show.bind(this)),l(document).on(e,i,this.openSortable.bind(this))},unbind:function(){this.$parent.off(e,a,this.show),l(document).off(e,i,this.openSortable)},openSortable:function(t){var e=l(t.target).data();this.BottomSheets=n.data("qor.bottomsheets"),this.selectedIconTmpl=l('[name="select-many-selected-icon"]').html(),e.ingoreSubmit=!0,e.url=e.selectListingUrl,e.selectDefaultCreating&&(e.url=e.selectCreatingUrl),this.BottomSheets.open(e,this.handleBottomSelect.bind(this))},show:function(){var t=this.$parent.find(d),e=this.$element,n=e.val(),i=[],o=this.$parent.find(".qor-dragable__list > li");n.length&&(o.each(function(){i.push(String(l(this).data("index")))}),i.length&&n.forEach(function(t){i.includes(t)||e.find("option:selected[data-index='".concat(t,"']")).prop("selected",!1)})),this.$element.trigger("change"),t.show(),this.$parent.find(a).hide(),setTimeout(function(){t.find(c).click()},100)},handleBottomSelect:function(){var t=l(r),e={onSelect:this.onSelectResults.bind(this),onSubmit:this.onSubmitResults.bind(this)};t.qorSelectCore(e).addClass("qor-bottomsheets__select-many"),this.initItems()},onSelectResults:function(t){var e=t.$clickElement,n=e.find("td:first"),i=this.collectData(t);l(h).find('li[data-index="'+i.id+'"]').length?(this.removeItems(i),e.removeClass(s),n.find(".qor-select__select-icon").remove()):(this.addItems(i),e.addClass(s),n.append(this.selectedIconTmpl))},onSubmitResults:function(t){this.addItems(this.collectData(t),!0)},collectData:function(t){var e=this.$element.data("remote-data-primary-key"),n={};return n.id=t[e]||t.primaryKey||t.Id||t.ID,n.value=t.Name||t.text||t.Text||t.Title||t.Code||n.id,n},initItems:function(){var n,t=l(r).find("tbody tr"),i=this.selectedIconTmpl,o=[];this.$sortableList.find("[data-index]").each(function(){o.push(l(this).data("index"))}),t.each(function(){var t=l(this),e=t.find("td:first");n=t.data().primaryKey,"-1"!=o.indexOf(n)&&(t.addClass(s),e.append(i))})},renderItem:function(t){return window.Mustache.render(f.LIST_HTML,t)},renderOption:function(){var t=this.sortable.toArray(),e=this.$parent.find(u);e.empty(),window._.each(t,function(t){e.append(window.Mustache.render(f.OPTION_HTML,{value:t}))})},removeItems:function(t){l(h).find('li[data-index="'+t.id+'"]').remove(),this.renderOption()},removeItemsFromList:function(t){this.renderOption();var e=t.index,n=t.value;e&&l('.select2-selection__choice[item-id="'.concat(e,'"]')).length?l('.select2-selection__choice[item-id="'.concat(e,'"]')).find(".select2-selection__choice__remove").click():n&&l('.select2-selection__choice[title="'.concat(n,'"]')).length&&l('.select2-selection__choice[title="'.concat(n,'"]')).find(".select2-selection__choice__remove").click()},addItems:function(t,e){this.$sortableList.append(this.renderItem(t)),this.renderOption(),e&&this.BottomSheets.hide()},destroy:function(){this.sortable.destroy(),this.unbind(),this.$element.select2("destroy").removeData(o)}},f.DEFAULTS={},f.LIST_HTML='<li data-index="[[id]]" data-result-id="[[_resultId]]" data-value="[[value]]"><span>[[value]]</span><div><i class="material-icons qor-dragable__list-delete">clear</i><i class="material-icons qor-dragable__list-handle">drag_handle</i></div></li>',f.OPTION_HTML='<option selected value="[[value]]"></option>',f.plugin=function(i){return this.each(function(){var t,e=l(this),n=e.data(o);if(!n){if(/destroy/.test(i))return;e.data(o,n=new f(this,i))}"string"==typeof i&&l.isFunction(t=n[i])&&t.apply(n)})},l(function(){var e='[data-toggle="qor.chooser.sortable"]';l(document).on("disable.qor.chooser.sortable",function(t){f.plugin.call(l(e,t.target),"destroy")}).on(t,function(t){f.plugin.call(l(e,t.target))}).triggerHandler(t)}),f});