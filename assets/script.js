/**
 * Minified by jsDelivr using Terser v5.13.1.
 * Original file: /npm/codemirror@5.65.4/lib/codemirror.js
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
!function(e,t){"object"==typeof exports&&"undefined"!=typeof module?module.exports=t():"function"==typeof define&&define.amd?define(t):(e=e||self).CodeMirror=t()}(this,(function(){"use strict";var e=navigator.userAgent,t=navigator.platform,r=/gecko\/\d/i.test(e),n=/MSIE \d/.test(e),i=/Trident\/(?:[7-9]|\d{2,})\..*rv:(\d+)/.exec(e),o=/Edge\/(\d+)/.exec(e),l=n||i||o,s=l&&(n?document.documentMode||6:+(o||i)[1]),a=!o&&/WebKit\//.test(e),u=a&&/Qt\/\d+\.\d+/.test(e),c=!o&&/Chrome\//.test(e),h=/Opera\//.test(e),f=/Apple Computer/.test(navigator.vendor),d=/Mac OS X 1\d\D([8-9]|\d\d)\D/.test(e),p=/PhantomJS/.test(e),g=f&&(/Mobile\/\w+/.test(e)||navigator.maxTouchPoints>2),v=/Android/.test(e),m=g||v||/webOS|BlackBerry|Opera Mini|Opera Mobi|IEMobile/i.test(e),y=g||/Mac/.test(t),b=/\bCrOS\b/.test(e),w=/win/i.test(t),x=h&&e.match(/Version\/(\d*\.\d*)/);x&&(x=Number(x[1])),x&&x>=15&&(h=!1,a=!0);var C=y&&(u||h&&(null==x||x<12.11)),S=r||l&&s>=9;function L(e){return new RegExp("(^|\\s)"+e+"(?:$|\\s)\\s*")}var k,T=function(e,t){var r=e.className,n=L(t).exec(r);if(n){var i=r.slice(n.index+n[0].length);e.className=r.slice(0,n.index)+(i?n[1]+i:"")}};function M(e){for(var t=e.childNodes.length;t>0;--t)e.removeChild(e.firstChild);return e}function N(e,t){return M(e).appendChild(t)}function O(e,t,r,n){var i=document.createElement(e);if(r&&(i.className=r),n&&(i.style.cssText=n),"string"==typeof t)i.appendChild(document.createTextNode(t));else if(t)for(var o=0;o<t.length;++o)i.appendChild(t[o]);return i}function A(e,t,r,n){var i=O(e,t,r,n);return i.setAttribute("role","presentation"),i}function D(e,t){if(3==t.nodeType&&(t=t.parentNode),e.contains)return e.contains(t);do{if(11==t.nodeType&&(t=t.host),t==e)return!0}while(t=t.parentNode)}function W(){var e;try{e=document.activeElement}catch(t){e=document.body||null}for(;e&&e.shadowRoot&&e.shadowRoot.activeElement;)e=e.shadowRoot.activeElement;return e}function H(e,t){var r=e.className;L(t).test(r)||(e.className+=(r?" ":"")+t)}function F(e,t){for(var r=e.split(" "),n=0;n<r.length;n++)r[n]&&!L(r[n]).test(t)&&(t+=" "+r[n]);return t}k=document.createRange?function(e,t,r,n){var i=document.createRange();return i.setEnd(n||e,r),i.setStart(e,t),i}:function(e,t,r){var n=document.body.createTextRange();try{n.moveToElementText(e.parentNode)}catch(e){return n}return n.collapse(!0),n.moveEnd("character",r),n.moveStart("character",t),n};var P=function(e){e.select()};function E(e){var t=Array.prototype.slice.call(arguments,1);return function(){return e.apply(null,t)}}function I(e,t,r){for(var n in t||(t={}),e)!e.hasOwnProperty(n)||!1===r&&t.hasOwnProperty(n)||(t[n]=e[n]);return t}function R(e,t,r,n,i){null==t&&-1==(t=e.search(/[^\s\u00a0]/))&&(t=e.length);for(var o=n||0,l=i||0;;){var s=e.indexOf("\t",o);if(s<0||s>=t)return l+(t-o);l+=s-o,l+=r-l%r,o=s+1}}g?P=function(e){e.selectionStart=0,e.selectionEnd=e.value.length}:l&&(P=function(e){try{e.select()}catch(e){}});var z=function(){this.id=null,this.f=null,this.time=0,this.handler=E(this.onTimeout,this)};function B(e,t){for(var r=0;r<e.length;++r)if(e[r]==t)return r;return-1}z.prototype.onTimeout=function(e){e.id=0,e.time<=+new Date?e.f():setTimeout(e.handler,e.time-+new Date)},z.prototype.set=function(e,t){this.f=t;var r=+new Date+e;(!this.id||r<this.time)&&(clearTimeout(this.id),this.id=setTimeout(this.handler,e),this.time=r)};var G={toString:function(){return"CodeMirror.Pass"}},U={scroll:!1},V={origin:"*mouse"},K={origin:"+move"};function j(e,t,r){for(var n=0,i=0;;){var o=e.indexOf("\t",n);-1==o&&(o=e.length);var l=o-n;if(o==e.length||i+l>=t)return n+Math.min(l,t-i);if(i+=o-n,n=o+1,(i+=r-i%r)>=t)return n}}var X=[""];function Y(e){for(;X.length<=e;)X.push($(X)+" ");return X[e]}function $(e){return e[e.length-1]}function _(e,t){for(var r=[],n=0;n<e.length;n++)r[n]=t(e[n],n);return r}function q(){}function Z(e,t){var r;return Object.create?r=Object.create(e):(q.prototype=e,r=new q),t&&I(t,r),r}var Q=/[\u00df\u0587\u0590-\u05f4\u0600-\u06ff\u3040-\u309f\u30a0-\u30ff\u3400-\u4db5\u4e00-\u9fcc\uac00-\ud7af]/;function J(e){return/\w/.test(e)||e>""&&(e.toUpperCase()!=e.toLowerCase()||Q.test(e))}function ee(e,t){return t?!!(t.source.indexOf("\\w")>-1&&J(e))||t.test(e):J(e)}function te(e){for(var t in e)if(e.hasOwnProperty(t)&&e[t])return!1;return!0}var re=/[\u0300-\u036f\u0483-\u0489\u0591-\u05bd\u05bf\u05c1\u05c2\u05c4\u05c5\u05c7\u0610-\u061a\u064b-\u065e\u0670\u06d6-\u06dc\u06de-\u06e4\u06e7\u06e8\u06ea-\u06ed\u0711\u0730-\u074a\u07a6-\u07b0\u07eb-\u07f3\u0816-\u0819\u081b-\u0823\u0825-\u0827\u0829-\u082d\u0900-\u0902\u093c\u0941-\u0948\u094d\u0951-\u0955\u0962\u0963\u0981\u09bc\u09be\u09c1-\u09c4\u09cd\u09d7\u09e2\u09e3\u0a01\u0a02\u0a3c\u0a41\u0a42\u0a47\u0a48\u0a4b-\u0a4d\u0a51\u0a70\u0a71\u0a75\u0a81\u0a82\u0abc\u0ac1-\u0ac5\u0ac7\u0ac8\u0acd\u0ae2\u0ae3\u0b01\u0b3c\u0b3e\u0b3f\u0b41-\u0b44\u0b4d\u0b56\u0b57\u0b62\u0b63\u0b82\u0bbe\u0bc0\u0bcd\u0bd7\u0c3e-\u0c40\u0c46-\u0c48\u0c4a-\u0c4d\u0c55\u0c56\u0c62\u0c63\u0cbc\u0cbf\u0cc2\u0cc6\u0ccc\u0ccd\u0cd5\u0cd6\u0ce2\u0ce3\u0d3e\u0d41-\u0d44\u0d4d\u0d57\u0d62\u0d63\u0dca\u0dcf\u0dd2-\u0dd4\u0dd6\u0ddf\u0e31\u0e34-\u0e3a\u0e47-\u0e4e\u0eb1\u0eb4-\u0eb9\u0ebb\u0ebc\u0ec8-\u0ecd\u0f18\u0f19\u0f35\u0f37\u0f39\u0f71-\u0f7e\u0f80-\u0f84\u0f86\u0f87\u0f90-\u0f97\u0f99-\u0fbc\u0fc6\u102d-\u1030\u1032-\u1037\u1039\u103a\u103d\u103e\u1058\u1059\u105e-\u1060\u1071-\u1074\u1082\u1085\u1086\u108d\u109d\u135f\u1712-\u1714\u1732-\u1734\u1752\u1753\u1772\u1773\u17b7-\u17bd\u17c6\u17c9-\u17d3\u17dd\u180b-\u180d\u18a9\u1920-\u1922\u1927\u1928\u1932\u1939-\u193b\u1a17\u1a18\u1a56\u1a58-\u1a5e\u1a60\u1a62\u1a65-\u1a6c\u1a73-\u1a7c\u1a7f\u1b00-\u1b03\u1b34\u1b36-\u1b3a\u1b3c\u1b42\u1b6b-\u1b73\u1b80\u1b81\u1ba2-\u1ba5\u1ba8\u1ba9\u1c2c-\u1c33\u1c36\u1c37\u1cd0-\u1cd2\u1cd4-\u1ce0\u1ce2-\u1ce8\u1ced\u1dc0-\u1de6\u1dfd-\u1dff\u200c\u200d\u20d0-\u20f0\u2cef-\u2cf1\u2de0-\u2dff\u302a-\u302f\u3099\u309a\ua66f-\ua672\ua67c\ua67d\ua6f0\ua6f1\ua802\ua806\ua80b\ua825\ua826\ua8c4\ua8e0-\ua8f1\ua926-\ua92d\ua947-\ua951\ua980-\ua982\ua9b3\ua9b6-\ua9b9\ua9bc\uaa29-\uaa2e\uaa31\uaa32\uaa35\uaa36\uaa43\uaa4c\uaab0\uaab2-\uaab4\uaab7\uaab8\uaabe\uaabf\uaac1\uabe5\uabe8\uabed\udc00-\udfff\ufb1e\ufe00-\ufe0f\ufe20-\ufe26\uff9e\uff9f]/;function ne(e){return e.charCodeAt(0)>=768&&re.test(e)}function ie(e,t,r){for(;(r<0?t>0:t<e.length)&&ne(e.charAt(t));)t+=r;return t}function oe(e,t,r){for(var n=t>r?-1:1;;){if(t==r)return t;var i=(t+r)/2,o=n<0?Math.ceil(i):Math.floor(i);if(o==t)return e(o)?t:r;e(o)?r=o:t=o+n}}var le=null;function se(e,t,r){var n;le=null;for(var i=0;i<e.length;++i){var o=e[i];if(o.from<t&&o.to>t)return i;o.to==t&&(o.from!=o.to&&"before"==r?n=i:le=i),o.from==t&&(o.from!=o.to&&"before"!=r?n=i:le=i)}return null!=n?n:le}var ae=function(){var e=/[\u0590-\u05f4\u0600-\u06ff\u0700-\u08ac]/,t=/[stwN]/,r=/[LRr]/,n=/[Lb1n]/,i=/[1n]/;function o(e,t,r){this.level=e,this.from=t,this.to=r}return function(l,s){var a="ltr"==s?"L":"R";if(0==l.length||"ltr"==s&&!e.test(l))return!1;for(var u,c=l.length,h=[],f=0;f<c;++f)h.push((u=l.charCodeAt(f))<=247?"bbbbbbbbbtstwsbbbbbbbbbbbbbbssstwNN%%%NNNNNN,N,N1111111111NNNNNNNLLLLLLLLLLLLLLLLLLLLLLLLLLNNNNNNLLLLLLLLLLLLLLLLLLLLLLLLLLNNNNbbbbbbsbbbbbbbbbbbbbbbbbbbbbbbbbb,N%%%%NNNNLNNNNN%%11NLNNN1LNNNNNLLLLLLLLLLLLLLLLLLLLLLLNLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLN".charAt(u):1424<=u&&u<=1524?"R":1536<=u&&u<=1785?"nnnnnnNNr%%r,rNNmmmmmmmmmmmrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrmmmmmmmmmmmmmmmmmmmmmnnnnnnnnnn%nnrrrmrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrmmmmmmmnNmmmmmmrrmmNmmmmrr1111111111".charAt(u-1536):1774<=u&&u<=2220?"r":8192<=u&&u<=8203?"w":8204==u?"b":"L");for(var d=0,p=a;d<c;++d){var g=h[d];"m"==g?h[d]=p:p=g}for(var v=0,m=a;v<c;++v){var y=h[v];"1"==y&&"r"==m?h[v]="n":r.test(y)&&(m=y,"r"==y&&(h[v]="R"))}for(var b=1,w=h[0];b<c-1;++b){var x=h[b];"+"==x&&"1"==w&&"1"==h[b+1]?h[b]="1":","!=x||w!=h[b+1]||"1"!=w&&"n"!=w||(h[b]=w),w=x}for(var C=0;C<c;++C){var S=h[C];if(","==S)h[C]="N";else if("%"==S){var L=void 0;for(L=C+1;L<c&&"%"==h[L];++L);for(var k=C&&"!"==h[C-1]||L<c&&"1"==h[L]?"1":"N",T=C;T<L;++T)h[T]=k;C=L-1}}for(var M=0,N=a;M<c;++M){var O=h[M];"L"==N&&"1"==O?h[M]="L":r.test(O)&&(N=O)}for(var A=0;A<c;++A)if(t.test(h[A])){var D=void 0;for(D=A+1;D<c&&t.test(h[D]);++D);for(var W="L"==(A?h[A-1]:a),H=W==("L"==(D<c?h[D]:a))?W?"L":"R":a,F=A;F<D;++F)h[F]=H;A=D-1}for(var P,E=[],I=0;I<c;)if(n.test(h[I])){var R=I;for(++I;I<c&&n.test(h[I]);++I);E.push(new o(0,R,I))}else{var z=I,B=E.length,G="rtl"==s?1:0;for(++I;I<c&&"L"!=h[I];++I);for(var U=z;U<I;)if(i.test(h[U])){z<U&&(E.splice(B,0,new o(1,z,U)),B+=G);var V=U;for(++U;U<I&&i.test(h[U]);++U);E.splice(B,0,new o(2,V,U)),B+=G,z=U}else++U;z<I&&E.splice(B,0,new o(1,z,I))}return"ltr"==s&&(1==E[0].level&&(P=l.match(/^\s+/))&&(E[0].from=P[0].length,E.unshift(new o(0,0,P[0].length))),1==$(E).level&&(P=l.match(/\s+$/))&&($(E).to-=P[0].length,E.push(new o(0,c-P[0].length,c)))),"rtl"==s?E.reverse():E}}();function ue(e,t){var r=e.order;return null==r&&(r=e.order=ae(e.text,t)),r}var ce=[],he=function(e,t,r){if(e.addEventListener)e.addEventListener(t,r,!1);else if(e.attachEvent)e.attachEvent("on"+t,r);else{var n=e._handlers||(e._handlers={});n[t]=(n[t]||ce).concat(r)}};function fe(e,t){return e._handlers&&e._handlers[t]||ce}function de(e,t,r){if(e.removeEventListener)e.removeEventListener(t,r,!1);else if(e.detachEvent)e.detachEvent("on"+t,r);else{var n=e._handlers,i=n&&n[t];if(i){var o=B(i,r);o>-1&&(n[t]=i.slice(0,o).concat(i.slice(o+1)))}}}function pe(e,t){var r=fe(e,t);if(r.length)for(var n=Array.prototype.slice.call(arguments,2),i=0;i<r.length;++i)r[i].apply(null,n)}function ge(e,t,r){return"string"==typeof t&&(t={type:t,preventDefault:function(){this.defaultPrevented=!0}}),pe(e,r||t.type,e,t),xe(t)||t.codemirrorIgnore}function ve(e){var t=e._handlers&&e._handlers.cursorActivity;if(t)for(var r=e.curOp.cursorActivityHandlers||(e.curOp.cursorActivityHandlers=[]),n=0;n<t.length;++n)-1==B(r,t[n])&&r.push(t[n])}function me(e,t){return fe(e,t).length>0}function ye(e){e.prototype.on=function(e,t){he(this,e,t)},e.prototype.off=function(e,t){de(this,e,t)}}function be(e){e.preventDefault?e.preventDefault():e.returnValue=!1}function we(e){e.stopPropagation?e.stopPropagation():e.cancelBubble=!0}function xe(e){return null!=e.defaultPrevented?e.defaultPrevented:0==e.returnValue}function Ce(e){be(e),we(e)}function Se(e){return e.target||e.srcElement}function Le(e){var t=e.which;return null==t&&(1&e.button?t=1:2&e.button?t=3:4&e.button&&(t=2)),y&&e.ctrlKey&&1==t&&(t=3),t}var ke,Te,Me=function(){if(l&&s<9)return!1;var e=O("div");return"draggable"in e||"dragDrop"in e}();function Ne(e){if(null==ke){var t=O("span","​");N(e,O("span",[t,document.createTextNode("x")])),0!=e.firstChild.offsetHeight&&(ke=t.offsetWidth<=1&&t.offsetHeight>2&&!(l&&s<8))}var r=ke?O("span","​"):O("span"," ",null,"display: inline-block; width: 1px; margin-right: -1px");return r.setAttribute("cm-text",""),r}function Oe(e){if(null!=Te)return Te;var t=N(e,document.createTextNode("AخA")),r=k(t,0,1).getBoundingClientRect(),n=k(t,1,2).getBoundingClientRect();return M(e),!(!r||r.left==r.right)&&(Te=n.right-r.right<3)}var Ae,De=3!="\n\nb".split(/\n/).length?function(e){for(var t=0,r=[],n=e.length;t<=n;){var i=e.indexOf("\n",t);-1==i&&(i=e.length);var o=e.slice(t,"\r"==e.charAt(i-1)?i-1:i),l=o.indexOf("\r");-1!=l?(r.push(o.slice(0,l)),t+=l+1):(r.push(o),t=i+1)}return r}:function(e){return e.split(/\r\n?|\n/)},We=window.getSelection?function(e){try{return e.selectionStart!=e.selectionEnd}catch(e){return!1}}:function(e){var t;try{t=e.ownerDocument.selection.createRange()}catch(e){}return!(!t||t.parentElement()!=e)&&0!=t.compareEndPoints("StartToEnd",t)},He="oncopy"in(Ae=O("div"))||(Ae.setAttribute("oncopy","return;"),"function"==typeof Ae.oncopy),Fe=null;var Pe={},Ee={};function Ie(e,t){arguments.length>2&&(t.dependencies=Array.prototype.slice.call(arguments,2)),Pe[e]=t}function Re(e){if("string"==typeof e&&Ee.hasOwnProperty(e))e=Ee[e];else if(e&&"string"==typeof e.name&&Ee.hasOwnProperty(e.name)){var t=Ee[e.name];"string"==typeof t&&(t={name:t}),(e=Z(t,e)).name=t.name}else{if("string"==typeof e&&/^[\w\-]+\/[\w\-]+\+xml$/.test(e))return Re("application/xml");if("string"==typeof e&&/^[\w\-]+\/[\w\-]+\+json$/.test(e))return Re("application/json")}return"string"==typeof e?{name:e}:e||{name:"null"}}function ze(e,t){t=Re(t);var r=Pe[t.name];if(!r)return ze(e,"text/plain");var n=r(e,t);if(Be.hasOwnProperty(t.name)){var i=Be[t.name];for(var o in i)i.hasOwnProperty(o)&&(n.hasOwnProperty(o)&&(n["_"+o]=n[o]),n[o]=i[o])}if(n.name=t.name,t.helperType&&(n.helperType=t.helperType),t.modeProps)for(var l in t.modeProps)n[l]=t.modeProps[l];return n}var Be={};function Ge(e,t){I(t,Be.hasOwnProperty(e)?Be[e]:Be[e]={})}function Ue(e,t){if(!0===t)return t;if(e.copyState)return e.copyState(t);var r={};for(var n in t){var i=t[n];i instanceof Array&&(i=i.concat([])),r[n]=i}return r}function Ve(e,t){for(var r;e.innerMode&&(r=e.innerMode(t))&&r.mode!=e;)t=r.state,e=r.mode;return r||{mode:e,state:t}}function Ke(e,t,r){return!e.startState||e.startState(t,r)}var je=function(e,t,r){this.pos=this.start=0,this.string=e,this.tabSize=t||8,this.lastColumnPos=this.lastColumnValue=0,this.lineStart=0,this.lineOracle=r};function Xe(e,t){if((t-=e.first)<0||t>=e.size)throw new Error("There is no line "+(t+e.first)+" in the document.");for(var r=e;!r.lines;)for(var n=0;;++n){var i=r.children[n],o=i.chunkSize();if(t<o){r=i;break}t-=o}return r.lines[t]}function Ye(e,t,r){var n=[],i=t.line;return e.iter(t.line,r.line+1,(function(e){var o=e.text;i==r.line&&(o=o.slice(0,r.ch)),i==t.line&&(o=o.slice(t.ch)),n.push(o),++i})),n}function $e(e,t,r){var n=[];return e.iter(t,r,(function(e){n.push(e.text)})),n}function _e(e,t){var r=t-e.height;if(r)for(var n=e;n;n=n.parent)n.height+=r}function qe(e){if(null==e.parent)return null;for(var t=e.parent,r=B(t.lines,e),n=t.parent;n;t=n,n=n.parent)for(var i=0;n.children[i]!=t;++i)r+=n.children[i].chunkSize();return r+t.first}function Ze(e,t){var r=e.first;e:do{for(var n=0;n<e.children.length;++n){var i=e.children[n],o=i.height;if(t<o){e=i;continue e}t-=o,r+=i.chunkSize()}return r}while(!e.lines);for(var l=0;l<e.lines.length;++l){var s=e.lines[l].height;if(t<s)break;t-=s}return r+l}function Qe(e,t){return t>=e.first&&t<e.first+e.size}function Je(e,t){return String(e.lineNumberFormatter(t+e.firstLineNumber))}function et(e,t,r){if(void 0===r&&(r=null),!(this instanceof et))return new et(e,t,r);this.line=e,this.ch=t,this.sticky=r}function tt(e,t){return e.line-t.line||e.ch-t.ch}function rt(e,t){return e.sticky==t.sticky&&0==tt(e,t)}function nt(e){return et(e.line,e.ch)}function it(e,t){return tt(e,t)<0?t:e}function ot(e,t){return tt(e,t)<0?e:t}function lt(e,t){return Math.max(e.first,Math.min(t,e.first+e.size-1))}function st(e,t){if(t.line<e.first)return et(e.first,0);var r=e.first+e.size-1;return t.line>r?et(r,Xe(e,r).text.length):function(e,t){var r=e.ch;return null==r||r>t?et(e.line,t):r<0?et(e.line,0):e}(t,Xe(e,t.line).text.length)}function at(e,t){for(var r=[],n=0;n<t.length;n++)r[n]=st(e,t[n]);return r}je.prototype.eol=function(){return this.pos>=this.string.length},je.prototype.sol=function(){return this.pos==this.lineStart},je.prototype.peek=function(){return this.string.charAt(this.pos)||void 0},je.prototype.next=function(){if(this.pos<this.string.length)return this.string.charAt(this.pos++)},je.prototype.eat=function(e){var t=this.string.charAt(this.pos);if("string"==typeof e?t==e:t&&(e.test?e.test(t):e(t)))return++this.pos,t},je.prototype.eatWhile=function(e){for(var t=this.pos;this.eat(e););return this.pos>t},je.prototype.eatSpace=function(){for(var e=this.pos;/[\s\u00a0]/.test(this.string.charAt(this.pos));)++this.pos;return this.pos>e},je.prototype.skipToEnd=function(){this.pos=this.string.length},je.prototype.skipTo=function(e){var t=this.string.indexOf(e,this.pos);if(t>-1)return this.pos=t,!0},je.prototype.backUp=function(e){this.pos-=e},je.prototype.column=function(){return this.lastColumnPos<this.start&&(this.lastColumnValue=R(this.string,this.start,this.tabSize,this.lastColumnPos,this.lastColumnValue),this.lastColumnPos=this.start),this.lastColumnValue-(this.lineStart?R(this.string,this.lineStart,this.tabSize):0)},je.prototype.indentation=function(){return R(this.string,null,this.tabSize)-(this.lineStart?R(this.string,this.lineStart,this.tabSize):0)},je.prototype.match=function(e,t,r){if("string"!=typeof e){var n=this.string.slice(this.pos).match(e);return n&&n.index>0?null:(n&&!1!==t&&(this.pos+=n[0].length),n)}var i=function(e){return r?e.toLowerCase():e};if(i(this.string.substr(this.pos,e.length))==i(e))return!1!==t&&(this.pos+=e.length),!0},je.prototype.current=function(){return this.string.slice(this.start,this.pos)},je.prototype.hideFirstChars=function(e,t){this.lineStart+=e;try{return t()}finally{this.lineStart-=e}},je.prototype.lookAhead=function(e){var t=this.lineOracle;return t&&t.lookAhead(e)},je.prototype.baseToken=function(){var e=this.lineOracle;return e&&e.baseToken(this.pos)};var ut=function(e,t){this.state=e,this.lookAhead=t},ct=function(e,t,r,n){this.state=t,this.doc=e,this.line=r,this.maxLookAhead=n||0,this.baseTokens=null,this.baseTokenPos=1};function ht(e,t,r,n){var i=[e.state.modeGen],o={};wt(e,t.text,e.doc.mode,r,(function(e,t){return i.push(e,t)}),o,n);for(var l=r.state,s=function(n){r.baseTokens=i;var s=e.state.overlays[n],a=1,u=0;r.state=!0,wt(e,t.text,s.mode,r,(function(e,t){for(var r=a;u<e;){var n=i[a];n>e&&i.splice(a,1,e,i[a+1],n),a+=2,u=Math.min(e,n)}if(t)if(s.opaque)i.splice(r,a-r,e,"overlay "+t),a=r+2;else for(;r<a;r+=2){var o=i[r+1];i[r+1]=(o?o+" ":"")+"overlay "+t}}),o),r.state=l,r.baseTokens=null,r.baseTokenPos=1},a=0;a<e.state.overlays.length;++a)s(a);return{styles:i,classes:o.bgClass||o.textClass?o:null}}function ft(e,t,r){if(!t.styles||t.styles[0]!=e.state.modeGen){var n=dt(e,qe(t)),i=t.text.length>e.options.maxHighlightLength&&Ue(e.doc.mode,n.state),o=ht(e,t,n);i&&(n.state=i),t.stateAfter=n.save(!i),t.styles=o.styles,o.classes?t.styleClasses=o.classes:t.styleClasses&&(t.styleClasses=null),r===e.doc.highlightFrontier&&(e.doc.modeFrontier=Math.max(e.doc.modeFrontier,++e.doc.highlightFrontier))}return t.styles}function dt(e,t,r){var n=e.doc,i=e.display;if(!n.mode.startState)return new ct(n,!0,t);var o=function(e,t,r){for(var n,i,o=e.doc,l=r?-1:t-(e.doc.mode.innerMode?1e3:100),s=t;s>l;--s){if(s<=o.first)return o.first;var a=Xe(o,s-1),u=a.stateAfter;if(u&&(!r||s+(u instanceof ut?u.lookAhead:0)<=o.modeFrontier))return s;var c=R(a.text,null,e.options.tabSize);(null==i||n>c)&&(i=s-1,n=c)}return i}(e,t,r),l=o>n.first&&Xe(n,o-1).stateAfter,s=l?ct.fromSaved(n,l,o):new ct(n,Ke(n.mode),o);return n.iter(o,t,(function(r){pt(e,r.text,s);var n=s.line;r.stateAfter=n==t-1||n%5==0||n>=i.viewFrom&&n<i.viewTo?s.save():null,s.nextLine()})),r&&(n.modeFrontier=s.line),s}function pt(e,t,r,n){var i=e.doc.mode,o=new je(t,e.options.tabSize,r);for(o.start=o.pos=n||0,""==t&&gt(i,r.state);!o.eol();)vt(i,o,r.state),o.start=o.pos}function gt(e,t){if(e.blankLine)return e.blankLine(t);if(e.innerMode){var r=Ve(e,t);return r.mode.blankLine?r.mode.blankLine(r.state):void 0}}function vt(e,t,r,n){for(var i=0;i<10;i++){n&&(n[0]=Ve(e,r).mode);var o=e.token(t,r);if(t.pos>t.start)return o}throw new Error("Mode "+e.name+" failed to advance stream.")}ct.prototype.lookAhead=function(e){var t=this.doc.getLine(this.line+e);return null!=t&&e>this.maxLookAhead&&(this.maxLookAhead=e),t},ct.prototype.baseToken=function(e){if(!this.baseTokens)return null;for(;this.baseTokens[this.baseTokenPos]<=e;)this.baseTokenPos+=2;var t=this.baseTokens[this.baseTokenPos+1];return{type:t&&t.replace(/( |^)overlay .*/,""),size:this.baseTokens[this.baseTokenPos]-e}},ct.prototype.nextLine=function(){this.line++,this.maxLookAhead>0&&this.maxLookAhead--},ct.fromSaved=function(e,t,r){return t instanceof ut?new ct(e,Ue(e.mode,t.state),r,t.lookAhead):new ct(e,Ue(e.mode,t),r)},ct.prototype.save=function(e){var t=!1!==e?Ue(this.doc.mode,this.state):this.state;return this.maxLookAhead>0?new ut(t,this.maxLookAhead):t};var mt=function(e,t,r){this.start=e.start,this.end=e.pos,this.string=e.current(),this.type=t||null,this.state=r};function yt(e,t,r,n){var i,o,l=e.doc,s=l.mode,a=Xe(l,(t=st(l,t)).line),u=dt(e,t.line,r),c=new je(a.text,e.options.tabSize,u);for(n&&(o=[]);(n||c.pos<t.ch)&&!c.eol();)c.start=c.pos,i=vt(s,c,u.state),n&&o.push(new mt(c,i,Ue(l.mode,u.state)));return n?o:new mt(c,i,u.state)}function bt(e,t){if(e)for(;;){var r=e.match(/(?:^|\s+)line-(background-)?(\S+)/);if(!r)break;e=e.slice(0,r.index)+e.slice(r.index+r[0].length);var n=r[1]?"bgClass":"textClass";null==t[n]?t[n]=r[2]:new RegExp("(?:^|\\s)"+r[2]+"(?:$|\\s)").test(t[n])||(t[n]+=" "+r[2])}return e}function wt(e,t,r,n,i,o,l){var s=r.flattenSpans;null==s&&(s=e.options.flattenSpans);var a,u=0,c=null,h=new je(t,e.options.tabSize,n),f=e.options.addModeClass&&[null];for(""==t&&bt(gt(r,n.state),o);!h.eol();){if(h.pos>e.options.maxHighlightLength?(s=!1,l&&pt(e,t,n,h.pos),h.pos=t.length,a=null):a=bt(vt(r,h,n.state,f),o),f){var d=f[0].name;d&&(a="m-"+(a?d+" "+a:d))}if(!s||c!=a){for(;u<h.start;)i(u=Math.min(h.start,u+5e3),c);c=a}h.start=h.pos}for(;u<h.pos;){var p=Math.min(h.pos,u+5e3);i(p,c),u=p}}var xt=!1,Ct=!1;function St(e,t,r){this.marker=e,this.from=t,this.to=r}function Lt(e,t){if(e)for(var r=0;r<e.length;++r){var n=e[r];if(n.marker==t)return n}}function kt(e,t){for(var r,n=0;n<e.length;++n)e[n]!=t&&(r||(r=[])).push(e[n]);return r}function Tt(e,t){if(t.full)return null;var r=Qe(e,t.from.line)&&Xe(e,t.from.line).markedSpans,n=Qe(e,t.to.line)&&Xe(e,t.to.line).markedSpans;if(!r&&!n)return null;var i=t.from.ch,o=t.to.ch,l=0==tt(t.from,t.to),s=function(e,t,r){var n;if(e)for(var i=0;i<e.length;++i){var o=e[i],l=o.marker;if(null==o.from||(l.inclusiveLeft?o.from<=t:o.from<t)||o.from==t&&"bookmark"==l.type&&(!r||!o.marker.insertLeft)){var s=null==o.to||(l.inclusiveRight?o.to>=t:o.to>t);(n||(n=[])).push(new St(l,o.from,s?null:o.to))}}return n}(r,i,l),a=function(e,t,r){var n;if(e)for(var i=0;i<e.length;++i){var o=e[i],l=o.marker;if(null==o.to||(l.inclusiveRight?o.to>=t:o.to>t)||o.from==t&&"bookmark"==l.type&&(!r||o.marker.insertLeft)){var s=null==o.from||(l.inclusiveLeft?o.from<=t:o.from<t);(n||(n=[])).push(new St(l,s?null:o.from-t,null==o.to?null:o.to-t))}}return n}(n,o,l),u=1==t.text.length,c=$(t.text).length+(u?i:0);if(s)for(var h=0;h<s.length;++h){var f=s[h];if(null==f.to){var d=Lt(a,f.marker);d?u&&(f.to=null==d.to?null:d.to+c):f.to=i}}if(a)for(var p=0;p<a.length;++p){var g=a[p];if(null!=g.to&&(g.to+=c),null==g.from)Lt(s,g.marker)||(g.from=c,u&&(s||(s=[])).push(g));else g.from+=c,u&&(s||(s=[])).push(g)}s&&(s=Mt(s)),a&&a!=s&&(a=Mt(a));var v=[s];if(!u){var m,y=t.text.length-2;if(y>0&&s)for(var b=0;b<s.length;++b)null==s[b].to&&(m||(m=[])).push(new St(s[b].marker,null,null));for(var w=0;w<y;++w)v.push(m);v.push(a)}return v}function Mt(e){for(var t=0;t<e.length;++t){var r=e[t];null!=r.from&&r.from==r.to&&!1!==r.marker.clearWhenEmpty&&e.splice(t--,1)}return e.length?e:null}function Nt(e){var t=e.markedSpans;if(t){for(var r=0;r<t.length;++r)t[r].marker.detachLine(e);e.markedSpans=null}}function Ot(e,t){if(t){for(var r=0;r<t.length;++r)t[r].marker.attachLine(e);e.markedSpans=t}}function At(e){return e.inclusiveLeft?-1:0}function Dt(e){return e.inclusiveRight?1:0}function Wt(e,t){var r=e.lines.length-t.lines.length;if(0!=r)return r;var n=e.find(),i=t.find(),o=tt(n.from,i.from)||At(e)-At(t);if(o)return-o;var l=tt(n.to,i.to)||Dt(e)-Dt(t);return l||t.id-e.id}function Ht(e,t){var r,n=Ct&&e.markedSpans;if(n)for(var i=void 0,o=0;o<n.length;++o)(i=n[o]).marker.collapsed&&null==(t?i.from:i.to)&&(!r||Wt(r,i.marker)<0)&&(r=i.marker);return r}function Ft(e){return Ht(e,!0)}function Pt(e){return Ht(e,!1)}function Et(e,t){var r,n=Ct&&e.markedSpans;if(n)for(var i=0;i<n.length;++i){var o=n[i];o.marker.collapsed&&(null==o.from||o.from<t)&&(null==o.to||o.to>t)&&(!r||Wt(r,o.marker)<0)&&(r=o.marker)}return r}function It(e,t,r,n,i){var o=Xe(e,t),l=Ct&&o.markedSpans;if(l)for(var s=0;s<l.length;++s){var a=l[s];if(a.marker.collapsed){var u=a.marker.find(0),c=tt(u.from,r)||At(a.marker)-At(i),h=tt(u.to,n)||Dt(a.marker)-Dt(i);if(!(c>=0&&h<=0||c<=0&&h>=0)&&(c<=0&&(a.marker.inclusiveRight&&i.inclusiveLeft?tt(u.to,r)>=0:tt(u.to,r)>0)||c>=0&&(a.marker.inclusiveRight&&i.inclusiveLeft?tt(u.from,n)<=0:tt(u.from,n)<0)))return!0}}}function Rt(e){for(var t;t=Ft(e);)e=t.find(-1,!0).line;return e}function zt(e,t){var r=Xe(e,t),n=Rt(r);return r==n?t:qe(n)}function Bt(e,t){if(t>e.lastLine())return t;var r,n=Xe(e,t);if(!Gt(e,n))return t;for(;r=Pt(n);)n=r.find(1,!0).line;return qe(n)+1}function Gt(e,t){var r=Ct&&t.markedSpans;if(r)for(var n=void 0,i=0;i<r.length;++i)if((n=r[i]).marker.collapsed){if(null==n.from)return!0;if(!n.marker.widgetNode&&0==n.from&&n.marker.inclusiveLeft&&Ut(e,t,n))return!0}}function Ut(e,t,r){if(null==r.to){var n=r.marker.find(1,!0);return Ut(e,n.line,Lt(n.line.markedSpans,r.marker))}if(r.marker.inclusiveRight&&r.to==t.text.length)return!0;for(var i=void 0,o=0;o<t.markedSpans.length;++o)if((i=t.markedSpans[o]).marker.collapsed&&!i.marker.widgetNode&&i.from==r.to&&(null==i.to||i.to!=r.from)&&(i.marker.inclusiveLeft||r.marker.inclusiveRight)&&Ut(e,t,i))return!0}function Vt(e){for(var t=0,r=(e=Rt(e)).parent,n=0;n<r.lines.length;++n){var i=r.lines[n];if(i==e)break;t+=i.height}for(var o=r.parent;o;o=(r=o).parent)for(var l=0;l<o.children.length;++l){var s=o.children[l];if(s==r)break;t+=s.height}return t}function Kt(e){if(0==e.height)return 0;for(var t,r=e.text.length,n=e;t=Ft(n);){var i=t.find(0,!0);n=i.from.line,r+=i.from.ch-i.to.ch}for(n=e;t=Pt(n);){var o=t.find(0,!0);r-=n.text.length-o.from.ch,r+=(n=o.to.line).text.length-o.to.ch}return r}function jt(e){var t=e.display,r=e.doc;t.maxLine=Xe(r,r.first),t.maxLineLength=Kt(t.maxLine),t.maxLineChanged=!0,r.iter((function(e){var r=Kt(e);r>t.maxLineLength&&(t.maxLineLength=r,t.maxLine=e)}))}var Xt=function(e,t,r){this.text=e,Ot(this,t),this.height=r?r(this):1};function Yt(e){e.parent=null,Nt(e)}Xt.prototype.lineNo=function(){return qe(this)},ye(Xt);var $t={},_t={};function qt(e,t){if(!e||/^\s*$/.test(e))return null;var r=t.addModeClass?_t:$t;return r[e]||(r[e]=e.replace(/\S+/g,"cm-$&"))}function Zt(e,t){var r=A("span",null,null,a?"padding-right: .1px":null),n={pre:A("pre",[r],"CodeMirror-line"),content:r,col:0,pos:0,cm:e,trailingSpace:!1,splitSpaces:e.getOption("lineWrapping")};t.measure={};for(var i=0;i<=(t.rest?t.rest.length:0);i++){var o=i?t.rest[i-1]:t.line,l=void 0;n.pos=0,n.addToken=Jt,Oe(e.display.measure)&&(l=ue(o,e.doc.direction))&&(n.addToken=er(n.addToken,l)),n.map=[],rr(o,n,ft(e,o,t!=e.display.externalMeasured&&qe(o))),o.styleClasses&&(o.styleClasses.bgClass&&(n.bgClass=F(o.styleClasses.bgClass,n.bgClass||"")),o.styleClasses.textClass&&(n.textClass=F(o.styleClasses.textClass,n.textClass||""))),0==n.map.length&&n.map.push(0,0,n.content.appendChild(Ne(e.display.measure))),0==i?(t.measure.map=n.map,t.measure.cache={}):((t.measure.maps||(t.measure.maps=[])).push(n.map),(t.measure.caches||(t.measure.caches=[])).push({}))}if(a){var s=n.content.lastChild;(/\bcm-tab\b/.test(s.className)||s.querySelector&&s.querySelector(".cm-tab"))&&(n.content.className="cm-tab-wrap-hack")}return pe(e,"renderLine",e,t.line,n.pre),n.pre.className&&(n.textClass=F(n.pre.className,n.textClass||"")),n}function Qt(e){var t=O("span","•","cm-invalidchar");return t.title="\\u"+e.charCodeAt(0).toString(16),t.setAttribute("aria-label",t.title),t}function Jt(e,t,r,n,i,o,a){if(t){var u,c=e.splitSpaces?function(e,t){if(e.length>1&&!/  /.test(e))return e;for(var r=t,n="",i=0;i<e.length;i++){var o=e.charAt(i);" "!=o||!r||i!=e.length-1&&32!=e.charCodeAt(i+1)||(o=" "),n+=o,r=" "==o}return n}(t,e.trailingSpace):t,h=e.cm.state.specialChars,f=!1;if(h.test(t)){u=document.createDocumentFragment();for(var d=0;;){h.lastIndex=d;var p=h.exec(t),g=p?p.index-d:t.length-d;if(g){var v=document.createTextNode(c.slice(d,d+g));l&&s<9?u.appendChild(O("span",[v])):u.appendChild(v),e.map.push(e.pos,e.pos+g,v),e.col+=g,e.pos+=g}if(!p)break;d+=g+1;var m=void 0;if("\t"==p[0]){var y=e.cm.options.tabSize,b=y-e.col%y;(m=u.appendChild(O("span",Y(b),"cm-tab"))).setAttribute("role","presentation"),m.setAttribute("cm-text","\t"),e.col+=b}else"\r"==p[0]||"\n"==p[0]?((m=u.appendChild(O("span","\r"==p[0]?"␍":"␤","cm-invalidchar"))).setAttribute("cm-text",p[0]),e.col+=1):((m=e.cm.options.specialCharPlaceholder(p[0])).setAttribute("cm-text",p[0]),l&&s<9?u.appendChild(O("span",[m])):u.appendChild(m),e.col+=1);e.map.push(e.pos,e.pos+1,m),e.pos++}}else e.col+=t.length,u=document.createTextNode(c),e.map.push(e.pos,e.pos+t.length,u),l&&s<9&&(f=!0),e.pos+=t.length;if(e.trailingSpace=32==c.charCodeAt(t.length-1),r||n||i||f||o||a){var w=r||"";n&&(w+=n),i&&(w+=i);var x=O("span",[u],w,o);if(a)for(var C in a)a.hasOwnProperty(C)&&"style"!=C&&"class"!=C&&x.setAttribute(C,a[C]);return e.content.appendChild(x)}e.content.appendChild(u)}}function er(e,t){return function(r,n,i,o,l,s,a){i=i?i+" cm-force-border":"cm-force-border";for(var u=r.pos,c=u+n.length;;){for(var h=void 0,f=0;f<t.length&&!((h=t[f]).to>u&&h.from<=u);f++);if(h.to>=c)return e(r,n,i,o,l,s,a);e(r,n.slice(0,h.to-u),i,o,null,s,a),o=null,n=n.slice(h.to-u),u=h.to}}}function tr(e,t,r,n){var i=!n&&r.widgetNode;i&&e.map.push(e.pos,e.pos+t,i),!n&&e.cm.display.input.needsContentAttribute&&(i||(i=e.content.appendChild(document.createElement("span"))),i.setAttribute("cm-marker",r.id)),i&&(e.cm.display.input.setUneditable(i),e.content.appendChild(i)),e.pos+=t,e.trailingSpace=!1}function rr(e,t,r){var n=e.markedSpans,i=e.text,o=0;if(n)for(var l,s,a,u,c,h,f,d=i.length,p=0,g=1,v="",m=0;;){if(m==p){a=u=c=s="",f=null,h=null,m=1/0;for(var y=[],b=void 0,w=0;w<n.length;++w){var x=n[w],C=x.marker;if("bookmark"==C.type&&x.from==p&&C.widgetNode)y.push(C);else if(x.from<=p&&(null==x.to||x.to>p||C.collapsed&&x.to==p&&x.from==p)){if(null!=x.to&&x.to!=p&&m>x.to&&(m=x.to,u=""),C.className&&(a+=" "+C.className),C.css&&(s=(s?s+";":"")+C.css),C.startStyle&&x.from==p&&(c+=" "+C.startStyle),C.endStyle&&x.to==m&&(b||(b=[])).push(C.endStyle,x.to),C.title&&((f||(f={})).title=C.title),C.attributes)for(var S in C.attributes)(f||(f={}))[S]=C.attributes[S];C.collapsed&&(!h||Wt(h.marker,C)<0)&&(h=x)}else x.from>p&&m>x.from&&(m=x.from)}if(b)for(var L=0;L<b.length;L+=2)b[L+1]==m&&(u+=" "+b[L]);if(!h||h.from==p)for(var k=0;k<y.length;++k)tr(t,0,y[k]);if(h&&(h.from||0)==p){if(tr(t,(null==h.to?d+1:h.to)-p,h.marker,null==h.from),null==h.to)return;h.to==p&&(h=!1)}}if(p>=d)break;for(var T=Math.min(d,m);;){if(v){var M=p+v.length;if(!h){var N=M>T?v.slice(0,T-p):v;t.addToken(t,N,l?l+a:a,c,p+N.length==m?u:"",s,f)}if(M>=T){v=v.slice(T-p),p=T;break}p=M,c=""}v=i.slice(o,o=r[g++]),l=qt(r[g++],t.cm.options)}}else for(var O=1;O<r.length;O+=2)t.addToken(t,i.slice(o,o=r[O]),qt(r[O+1],t.cm.options))}function nr(e,t,r){this.line=t,this.rest=function(e){for(var t,r;t=Pt(e);)e=t.find(1,!0).line,(r||(r=[])).push(e);return r}(t),this.size=this.rest?qe($(this.rest))-r+1:1,this.node=this.text=null,this.hidden=Gt(e,t)}function ir(e,t,r){for(var n,i=[],o=t;o<r;o=n){var l=new nr(e.doc,Xe(e.doc,o),o);n=o+l.size,i.push(l)}return i}var or=null;var lr=null;function sr(e,t){var r=fe(e,t);if(r.length){var n,i=Array.prototype.slice.call(arguments,2);or?n=or.delayedCallbacks:lr?n=lr:(n=lr=[],setTimeout(ar,0));for(var o=function(e){n.push((function(){return r[e].apply(null,i)}))},l=0;l<r.length;++l)o(l)}}function ar(){var e=lr;lr=null;for(var t=0;t<e.length;++t)e[t]()}function ur(e,t,r,n){for(var i=0;i<t.changes.length;i++){var o=t.changes[i];"text"==o?fr(e,t):"gutter"==o?pr(e,t,r,n):"class"==o?dr(e,t):"widget"==o&&gr(e,t,n)}t.changes=null}function cr(e){return e.node==e.text&&(e.node=O("div",null,null,"position: relative"),e.text.parentNode&&e.text.parentNode.replaceChild(e.node,e.text),e.node.appendChild(e.text),l&&s<8&&(e.node.style.zIndex=2)),e.node}function hr(e,t){var r=e.display.externalMeasured;return r&&r.line==t.line?(e.display.externalMeasured=null,t.measure=r.measure,r.built):Zt(e,t)}function fr(e,t){var r=t.text.className,n=hr(e,t);t.text==t.node&&(t.node=n.pre),t.text.parentNode.replaceChild(n.pre,t.text),t.text=n.pre,n.bgClass!=t.bgClass||n.textClass!=t.textClass?(t.bgClass=n.bgClass,t.textClass=n.textClass,dr(e,t)):r&&(t.text.className=r)}function dr(e,t){!function(e,t){var r=t.bgClass?t.bgClass+" "+(t.line.bgClass||""):t.line.bgClass;if(r&&(r+=" CodeMirror-linebackground"),t.background)r?t.background.className=r:(t.background.parentNode.removeChild(t.background),t.background=null);else if(r){var n=cr(t);t.background=n.insertBefore(O("div",null,r),n.firstChild),e.display.input.setUneditable(t.background)}}(e,t),t.line.wrapClass?cr(t).className=t.line.wrapClass:t.node!=t.text&&(t.node.className="");var r=t.textClass?t.textClass+" "+(t.line.textClass||""):t.line.textClass;t.text.className=r||""}function pr(e,t,r,n){if(t.gutter&&(t.node.removeChild(t.gutter),t.gutter=null),t.gutterBackground&&(t.node.removeChild(t.gutterBackground),t.gutterBackground=null),t.line.gutterClass){var i=cr(t);t.gutterBackground=O("div",null,"CodeMirror-gutter-background "+t.line.gutterClass,"left: "+(e.options.fixedGutter?n.fixedPos:-n.gutterTotalWidth)+"px; width: "+n.gutterTotalWidth+"px"),e.display.input.setUneditable(t.gutterBackground),i.insertBefore(t.gutterBackground,t.text)}var o=t.line.gutterMarkers;if(e.options.lineNumbers||o){var l=cr(t),s=t.gutter=O("div",null,"CodeMirror-gutter-wrapper","left: "+(e.options.fixedGutter?n.fixedPos:-n.gutterTotalWidth)+"px");if(s.setAttribute("aria-hidden","true"),e.display.input.setUneditable(s),l.insertBefore(s,t.text),t.line.gutterClass&&(s.className+=" "+t.line.gutterClass),!e.options.lineNumbers||o&&o["CodeMirror-linenumbers"]||(t.lineNumber=s.appendChild(O("div",Je(e.options,r),"CodeMirror-linenumber CodeMirror-gutter-elt","left: "+n.gutterLeft["CodeMirror-linenumbers"]+"px; width: "+e.display.lineNumInnerWidth+"px"))),o)for(var a=0;a<e.display.gutterSpecs.length;++a){var u=e.display.gutterSpecs[a].className,c=o.hasOwnProperty(u)&&o[u];c&&s.appendChild(O("div",[c],"CodeMirror-gutter-elt","left: "+n.gutterLeft[u]+"px; width: "+n.gutterWidth[u]+"px"))}}}function gr(e,t,r){t.alignable&&(t.alignable=null);for(var n=L("CodeMirror-linewidget"),i=t.node.firstChild,o=void 0;i;i=o)o=i.nextSibling,n.test(i.className)&&t.node.removeChild(i);mr(e,t,r)}function vr(e,t,r,n){var i=hr(e,t);return t.text=t.node=i.pre,i.bgClass&&(t.bgClass=i.bgClass),i.textClass&&(t.textClass=i.textClass),dr(e,t),pr(e,t,r,n),mr(e,t,n),t.node}function mr(e,t,r){if(yr(e,t.line,t,r,!0),t.rest)for(var n=0;n<t.rest.length;n++)yr(e,t.rest[n],t,r,!1)}function yr(e,t,r,n,i){if(t.widgets)for(var o=cr(r),l=0,s=t.widgets;l<s.length;++l){var a=s[l],u=O("div",[a.node],"CodeMirror-linewidget"+(a.className?" "+a.className:""));a.handleMouseEvents||u.setAttribute("cm-ignore-events","true"),br(a,u,r,n),e.display.input.setUneditable(u),i&&a.above?o.insertBefore(u,r.gutter||r.text):o.appendChild(u),sr(a,"redraw")}}function br(e,t,r,n){if(e.noHScroll){(r.alignable||(r.alignable=[])).push(t);var i=n.wrapperWidth;t.style.left=n.fixedPos+"px",e.coverGutter||(i-=n.gutterTotalWidth,t.style.paddingLeft=n.gutterTotalWidth+"px"),t.style.width=i+"px"}e.coverGutter&&(t.style.zIndex=5,t.style.position="relative",e.noHScroll||(t.style.marginLeft=-n.gutterTotalWidth+"px"))}function wr(e){if(null!=e.height)return e.height;var t=e.doc.cm;if(!t)return 0;if(!D(document.body,e.node)){var r="position: relative;";e.coverGutter&&(r+="margin-left: -"+t.display.gutters.offsetWidth+"px;"),e.noHScroll&&(r+="width: "+t.display.wrapper.clientWidth+"px;"),N(t.display.measure,O("div",[e.node],null,r))}return e.height=e.node.parentNode.offsetHeight}function xr(e,t){for(var r=Se(t);r!=e.wrapper;r=r.parentNode)if(!r||1==r.nodeType&&"true"==r.getAttribute("cm-ignore-events")||r.parentNode==e.sizer&&r!=e.mover)return!0}function Cr(e){return e.lineSpace.offsetTop}function Sr(e){return e.mover.offsetHeight-e.lineSpace.offsetHeight}function Lr(e){if(e.cachedPaddingH)return e.cachedPaddingH;var t=N(e.measure,O("pre","x","CodeMirror-line-like")),r=window.getComputedStyle?window.getComputedStyle(t):t.currentStyle,n={left:parseInt(r.paddingLeft),right:parseInt(r.paddingRight)};return isNaN(n.left)||isNaN(n.right)||(e.cachedPaddingH=n),n}function kr(e){return 50-e.display.nativeBarWidth}function Tr(e){return e.display.scroller.clientWidth-kr(e)-e.display.barWidth}function Mr(e){return e.display.scroller.clientHeight-kr(e)-e.display.barHeight}function Nr(e,t,r){if(e.line==t)return{map:e.measure.map,cache:e.measure.cache};if(e.rest){for(var n=0;n<e.rest.length;n++)if(e.rest[n]==t)return{map:e.measure.maps[n],cache:e.measure.caches[n]};for(var i=0;i<e.rest.length;i++)if(qe(e.rest[i])>r)return{map:e.measure.maps[i],cache:e.measure.caches[i],before:!0}}}function Or(e,t,r,n){return Wr(e,Dr(e,t),r,n)}function Ar(e,t){if(t>=e.display.viewFrom&&t<e.display.viewTo)return e.display.view[cn(e,t)];var r=e.display.externalMeasured;return r&&t>=r.lineN&&t<r.lineN+r.size?r:void 0}function Dr(e,t){var r=qe(t),n=Ar(e,r);n&&!n.text?n=null:n&&n.changes&&(ur(e,n,r,on(e)),e.curOp.forceUpdate=!0),n||(n=function(e,t){var r=qe(t=Rt(t)),n=e.display.externalMeasured=new nr(e.doc,t,r);n.lineN=r;var i=n.built=Zt(e,n);return n.text=i.pre,N(e.display.lineMeasure,i.pre),n}(e,t));var i=Nr(n,t,r);return{line:t,view:n,rect:null,map:i.map,cache:i.cache,before:i.before,hasHeights:!1}}function Wr(e,t,r,n,i){t.before&&(r=-1);var o,a=r+(n||"");return t.cache.hasOwnProperty(a)?o=t.cache[a]:(t.rect||(t.rect=t.view.text.getBoundingClientRect()),t.hasHeights||(!function(e,t,r){var n=e.options.lineWrapping,i=n&&Tr(e);if(!t.measure.heights||n&&t.measure.width!=i){var o=t.measure.heights=[];if(n){t.measure.width=i;for(var l=t.text.firstChild.getClientRects(),s=0;s<l.length-1;s++){var a=l[s],u=l[s+1];Math.abs(a.bottom-u.bottom)>2&&o.push((a.bottom+u.top)/2-r.top)}}o.push(r.bottom-r.top)}}(e,t.view,t.rect),t.hasHeights=!0),o=function(e,t,r,n){var i,o=Pr(t.map,r,n),a=o.node,u=o.start,c=o.end,h=o.collapse;if(3==a.nodeType){for(var f=0;f<4;f++){for(;u&&ne(t.line.text.charAt(o.coverStart+u));)--u;for(;o.coverStart+c<o.coverEnd&&ne(t.line.text.charAt(o.coverStart+c));)++c;if((i=l&&s<9&&0==u&&c==o.coverEnd-o.coverStart?a.parentNode.getBoundingClientRect():Er(k(a,u,c).getClientRects(),n)).left||i.right||0==u)break;c=u,u-=1,h="right"}l&&s<11&&(i=function(e,t){if(!window.screen||null==screen.logicalXDPI||screen.logicalXDPI==screen.deviceXDPI||!function(e){if(null!=Fe)return Fe;var t=N(e,O("span","x")),r=t.getBoundingClientRect(),n=k(t,0,1).getBoundingClientRect();return Fe=Math.abs(r.left-n.left)>1}(e))return t;var r=screen.logicalXDPI/screen.deviceXDPI,n=screen.logicalYDPI/screen.deviceYDPI;return{left:t.left*r,right:t.right*r,top:t.top*n,bottom:t.bottom*n}}(e.display.measure,i))}else{var d;u>0&&(h=n="right"),i=e.options.lineWrapping&&(d=a.getClientRects()).length>1?d["right"==n?d.length-1:0]:a.getBoundingClientRect()}if(l&&s<9&&!u&&(!i||!i.left&&!i.right)){var p=a.parentNode.getClientRects()[0];i=p?{left:p.left,right:p.left+nn(e.display),top:p.top,bottom:p.bottom}:Fr}for(var g=i.top-t.rect.top,v=i.bottom-t.rect.top,m=(g+v)/2,y=t.view.measure.heights,b=0;b<y.length-1&&!(m<y[b]);b++);var w=b?y[b-1]:0,x=y[b],C={left:("right"==h?i.right:i.left)-t.rect.left,right:("left"==h?i.left:i.right)-t.rect.left,top:w,bottom:x};i.left||i.right||(C.bogus=!0);e.options.singleCursorHeightPerLine||(C.rtop=g,C.rbottom=v);return C}(e,t,r,n),o.bogus||(t.cache[a]=o)),{left:o.left,right:o.right,top:i?o.rtop:o.top,bottom:i?o.rbottom:o.bottom}}var Hr,Fr={left:0,right:0,top:0,bottom:0};function Pr(e,t,r){for(var n,i,o,l,s,a,u=0;u<e.length;u+=3)if(s=e[u],a=e[u+1],t<s?(i=0,o=1,l="left"):t<a?o=(i=t-s)+1:(u==e.length-3||t==a&&e[u+3]>t)&&(i=(o=a-s)-1,t>=a&&(l="right")),null!=i){if(n=e[u+2],s==a&&r==(n.insertLeft?"left":"right")&&(l=r),"left"==r&&0==i)for(;u&&e[u-2]==e[u-3]&&e[u-1].insertLeft;)n=e[2+(u-=3)],l="left";if("right"==r&&i==a-s)for(;u<e.length-3&&e[u+3]==e[u+4]&&!e[u+5].insertLeft;)n=e[(u+=3)+2],l="right";break}return{node:n,start:i,end:o,collapse:l,coverStart:s,coverEnd:a}}function Er(e,t){var r=Fr;if("left"==t)for(var n=0;n<e.length&&(r=e[n]).left==r.right;n++);else for(var i=e.length-1;i>=0&&(r=e[i]).left==r.right;i--);return r}function Ir(e){if(e.measure&&(e.measure.cache={},e.measure.heights=null,e.rest))for(var t=0;t<e.rest.length;t++)e.measure.caches[t]={}}function Rr(e){e.display.externalMeasure=null,M(e.display.lineMeasure);for(var t=0;t<e.display.view.length;t++)Ir(e.display.view[t])}function zr(e){Rr(e),e.display.cachedCharWidth=e.display.cachedTextHeight=e.display.cachedPaddingH=null,e.options.lineWrapping||(e.display.maxLineChanged=!0),e.display.lineNumChars=null}function Br(){return c&&v?-(document.body.getBoundingClientRect().left-parseInt(getComputedStyle(document.body).marginLeft)):window.pageXOffset||(document.documentElement||document.body).scrollLeft}function Gr(){return c&&v?-(document.body.getBoundingClientRect().top-parseInt(getComputedStyle(document.body).marginTop)):window.pageYOffset||(document.documentElement||document.body).scrollTop}function Ur(e){var t=Rt(e).widgets,r=0;if(t)for(var n=0;n<t.length;++n)t[n].above&&(r+=wr(t[n]));return r}function Vr(e,t,r,n,i){if(!i){var o=Ur(t);r.top+=o,r.bottom+=o}if("line"==n)return r;n||(n="local");var l=Vt(t);if("local"==n?l+=Cr(e.display):l-=e.display.viewOffset,"page"==n||"window"==n){var s=e.display.lineSpace.getBoundingClientRect();l+=s.top+("window"==n?0:Gr());var a=s.left+("window"==n?0:Br());r.left+=a,r.right+=a}return r.top+=l,r.bottom+=l,r}function Kr(e,t,r){if("div"==r)return t;var n=t.left,i=t.top;if("page"==r)n-=Br(),i-=Gr();else if("local"==r||!r){var o=e.display.sizer.getBoundingClientRect();n+=o.left,i+=o.top}var l=e.display.lineSpace.getBoundingClientRect();return{left:n-l.left,top:i-l.top}}function jr(e,t,r,n,i){return n||(n=Xe(e.doc,t.line)),Vr(e,n,Or(e,n,t.ch,i),r)}function Xr(e,t,r,n,i,o){function l(t,l){var s=Wr(e,i,t,l?"right":"left",o);return l?s.left=s.right:s.right=s.left,Vr(e,n,s,r)}n=n||Xe(e.doc,t.line),i||(i=Dr(e,n));var s=ue(n,e.doc.direction),a=t.ch,u=t.sticky;if(a>=n.text.length?(a=n.text.length,u="before"):a<=0&&(a=0,u="after"),!s)return l("before"==u?a-1:a,"before"==u);function c(e,t,r){return l(r?e-1:e,1==s[t].level!=r)}var h=se(s,a,u),f=le,d=c(a,h,"before"==u);return null!=f&&(d.other=c(a,f,"before"!=u)),d}function Yr(e,t){var r=0;t=st(e.doc,t),e.options.lineWrapping||(r=nn(e.display)*t.ch);var n=Xe(e.doc,t.line),i=Vt(n)+Cr(e.display);return{left:r,right:r,top:i,bottom:i+n.height}}function $r(e,t,r,n,i){var o=et(e,t,r);return o.xRel=i,n&&(o.outside=n),o}function _r(e,t,r){var n=e.doc;if((r+=e.display.viewOffset)<0)return $r(n.first,0,null,-1,-1);var i=Ze(n,r),o=n.first+n.size-1;if(i>o)return $r(n.first+n.size-1,Xe(n,o).text.length,null,1,1);t<0&&(t=0);for(var l=Xe(n,i);;){var s=Jr(e,l,i,t,r),a=Et(l,s.ch+(s.xRel>0||s.outside>0?1:0));if(!a)return s;var u=a.find(1);if(u.line==i)return u;l=Xe(n,i=u.line)}}function qr(e,t,r,n){n-=Ur(t);var i=t.text.length,o=oe((function(t){return Wr(e,r,t-1).bottom<=n}),i,0);return{begin:o,end:i=oe((function(t){return Wr(e,r,t).top>n}),o,i)}}function Zr(e,t,r,n){return r||(r=Dr(e,t)),qr(e,t,r,Vr(e,t,Wr(e,r,n),"line").top)}function Qr(e,t,r,n){return!(e.bottom<=r)&&(e.top>r||(n?e.left:e.right)>t)}function Jr(e,t,r,n,i){i-=Vt(t);var o=Dr(e,t),l=Ur(t),s=0,a=t.text.length,u=!0,c=ue(t,e.doc.direction);if(c){var h=(e.options.lineWrapping?tn:en)(e,t,r,o,c,n,i);s=(u=1!=h.level)?h.from:h.to-1,a=u?h.to:h.from-1}var f,d,p=null,g=null,v=oe((function(t){var r=Wr(e,o,t);return r.top+=l,r.bottom+=l,!!Qr(r,n,i,!1)&&(r.top<=i&&r.left<=n&&(p=t,g=r),!0)}),s,a),m=!1;if(g){var y=n-g.left<g.right-n,b=y==u;v=p+(b?0:1),d=b?"after":"before",f=y?g.left:g.right}else{u||v!=a&&v!=s||v++,d=0==v?"after":v==t.text.length?"before":Wr(e,o,v-(u?1:0)).bottom+l<=i==u?"after":"before";var w=Xr(e,et(r,v,d),"line",t,o);f=w.left,m=i<w.top?-1:i>=w.bottom?1:0}return $r(r,v=ie(t.text,v,1),d,m,n-f)}function en(e,t,r,n,i,o,l){var s=oe((function(s){var a=i[s],u=1!=a.level;return Qr(Xr(e,et(r,u?a.to:a.from,u?"before":"after"),"line",t,n),o,l,!0)}),0,i.length-1),a=i[s];if(s>0){var u=1!=a.level,c=Xr(e,et(r,u?a.from:a.to,u?"after":"before"),"line",t,n);Qr(c,o,l,!0)&&c.top>l&&(a=i[s-1])}return a}function tn(e,t,r,n,i,o,l){var s=qr(e,t,n,l),a=s.begin,u=s.end;/\s/.test(t.text.charAt(u-1))&&u--;for(var c=null,h=null,f=0;f<i.length;f++){var d=i[f];if(!(d.from>=u||d.to<=a)){var p=Wr(e,n,1!=d.level?Math.min(u,d.to)-1:Math.max(a,d.from)).right,g=p<o?o-p+1e9:p-o;(!c||h>g)&&(c=d,h=g)}}return c||(c=i[i.length-1]),c.from<a&&(c={from:a,to:c.to,level:c.level}),c.to>u&&(c={from:c.from,to:u,level:c.level}),c}function rn(e){if(null!=e.cachedTextHeight)return e.cachedTextHeight;if(null==Hr){Hr=O("pre",null,"CodeMirror-line-like");for(var t=0;t<49;++t)Hr.appendChild(document.createTextNode("x")),Hr.appendChild(O("br"));Hr.appendChild(document.createTextNode("x"))}N(e.measure,Hr);var r=Hr.offsetHeight/50;return r>3&&(e.cachedTextHeight=r),M(e.measure),r||1}function nn(e){if(null!=e.cachedCharWidth)return e.cachedCharWidth;var t=O("span","xxxxxxxxxx"),r=O("pre",[t],"CodeMirror-line-like");N(e.measure,r);var n=t.getBoundingClientRect(),i=(n.right-n.left)/10;return i>2&&(e.cachedCharWidth=i),i||10}function on(e){for(var t=e.display,r={},n={},i=t.gutters.clientLeft,o=t.gutters.firstChild,l=0;o;o=o.nextSibling,++l){var s=e.display.gutterSpecs[l].className;r[s]=o.offsetLeft+o.clientLeft+i,n[s]=o.clientWidth}return{fixedPos:ln(t),gutterTotalWidth:t.gutters.offsetWidth,gutterLeft:r,gutterWidth:n,wrapperWidth:t.wrapper.clientWidth}}function ln(e){return e.scroller.getBoundingClientRect().left-e.sizer.getBoundingClientRect().left}function sn(e){var t=rn(e.display),r=e.options.lineWrapping,n=r&&Math.max(5,e.display.scroller.clientWidth/nn(e.display)-3);return function(i){if(Gt(e.doc,i))return 0;var o=0;if(i.widgets)for(var l=0;l<i.widgets.length;l++)i.widgets[l].height&&(o+=i.widgets[l].height);return r?o+(Math.ceil(i.text.length/n)||1)*t:o+t}}function an(e){var t=e.doc,r=sn(e);t.iter((function(e){var t=r(e);t!=e.height&&_e(e,t)}))}function un(e,t,r,n){var i=e.display;if(!r&&"true"==Se(t).getAttribute("cm-not-content"))return null;var o,l,s=i.lineSpace.getBoundingClientRect();try{o=t.clientX-s.left,l=t.clientY-s.top}catch(e){return null}var a,u=_r(e,o,l);if(n&&u.xRel>0&&(a=Xe(e.doc,u.line).text).length==u.ch){var c=R(a,a.length,e.options.tabSize)-a.length;u=et(u.line,Math.max(0,Math.round((o-Lr(e.display).left)/nn(e.display))-c))}return u}function cn(e,t){if(t>=e.display.viewTo)return null;if((t-=e.display.viewFrom)<0)return null;for(var r=e.display.view,n=0;n<r.length;n++)if((t-=r[n].size)<0)return n}function hn(e,t,r,n){null==t&&(t=e.doc.first),null==r&&(r=e.doc.first+e.doc.size),n||(n=0);var i=e.display;if(n&&r<i.viewTo&&(null==i.updateLineNumbers||i.updateLineNumbers>t)&&(i.updateLineNumbers=t),e.curOp.viewChanged=!0,t>=i.viewTo)Ct&&zt(e.doc,t)<i.viewTo&&dn(e);else if(r<=i.viewFrom)Ct&&Bt(e.doc,r+n)>i.viewFrom?dn(e):(i.viewFrom+=n,i.viewTo+=n);else if(t<=i.viewFrom&&r>=i.viewTo)dn(e);else if(t<=i.viewFrom){var o=pn(e,r,r+n,1);o?(i.view=i.view.slice(o.index),i.viewFrom=o.lineN,i.viewTo+=n):dn(e)}else if(r>=i.viewTo){var l=pn(e,t,t,-1);l?(i.view=i.view.slice(0,l.index),i.viewTo=l.lineN):dn(e)}else{var s=pn(e,t,t,-1),a=pn(e,r,r+n,1);s&&a?(i.view=i.view.slice(0,s.index).concat(ir(e,s.lineN,a.lineN)).concat(i.view.slice(a.index)),i.viewTo+=n):dn(e)}var u=i.externalMeasured;u&&(r<u.lineN?u.lineN+=n:t<u.lineN+u.size&&(i.externalMeasured=null))}function fn(e,t,r){e.curOp.viewChanged=!0;var n=e.display,i=e.display.externalMeasured;if(i&&t>=i.lineN&&t<i.lineN+i.size&&(n.externalMeasured=null),!(t<n.viewFrom||t>=n.viewTo)){var o=n.view[cn(e,t)];if(null!=o.node){var l=o.changes||(o.changes=[]);-1==B(l,r)&&l.push(r)}}}function dn(e){e.display.viewFrom=e.display.viewTo=e.doc.first,e.display.view=[],e.display.viewOffset=0}function pn(e,t,r,n){var i,o=cn(e,t),l=e.display.view;if(!Ct||r==e.doc.first+e.doc.size)return{index:o,lineN:r};for(var s=e.display.viewFrom,a=0;a<o;a++)s+=l[a].size;if(s!=t){if(n>0){if(o==l.length-1)return null;i=s+l[o].size-t,o++}else i=s-t;t+=i,r+=i}for(;zt(e.doc,r)!=r;){if(o==(n<0?0:l.length-1))return null;r+=n*l[o-(n<0?1:0)].size,o+=n}return{index:o,lineN:r}}function gn(e){for(var t=e.display.view,r=0,n=0;n<t.length;n++){var i=t[n];i.hidden||i.node&&!i.changes||++r}return r}function vn(e){e.display.input.showSelection(e.display.input.prepareSelection())}function mn(e,t){void 0===t&&(t=!0);var r=e.doc,n={},i=n.cursors=document.createDocumentFragment(),o=n.selection=document.createDocumentFragment(),l=e.options.$customCursor;l&&(t=!0);for(var s=0;s<r.sel.ranges.length;s++)if(t||s!=r.sel.primIndex){var a=r.sel.ranges[s];if(!(a.from().line>=e.display.viewTo||a.to().line<e.display.viewFrom)){var u=a.empty();if(l){var c=l(e,a);c&&yn(e,c,i)}else(u||e.options.showCursorWhenSelecting)&&yn(e,a.head,i);u||wn(e,a,o)}}return n}function yn(e,t,r){var n=Xr(e,t,"div",null,null,!e.options.singleCursorHeightPerLine),i=r.appendChild(O("div"," ","CodeMirror-cursor"));if(i.style.left=n.left+"px",i.style.top=n.top+"px",i.style.height=Math.max(0,n.bottom-n.top)*e.options.cursorHeight+"px",/\bcm-fat-cursor\b/.test(e.getWrapperElement().className)){var o=jr(e,t,"div",null,null),l=o.right-o.left;i.style.width=(l>0?l:e.defaultCharWidth())+"px"}if(n.other){var s=r.appendChild(O("div"," ","CodeMirror-cursor CodeMirror-secondarycursor"));s.style.display="",s.style.left=n.other.left+"px",s.style.top=n.other.top+"px",s.style.height=.85*(n.other.bottom-n.other.top)+"px"}}function bn(e,t){return e.top-t.top||e.left-t.left}function wn(e,t,r){var n=e.display,i=e.doc,o=document.createDocumentFragment(),l=Lr(e.display),s=l.left,a=Math.max(n.sizerWidth,Tr(e)-n.sizer.offsetLeft)-l.right,u="ltr"==i.direction;function c(e,t,r,n){t<0&&(t=0),t=Math.round(t),n=Math.round(n),o.appendChild(O("div",null,"CodeMirror-selected","position: absolute; left: "+e+"px;\n                             top: "+t+"px; width: "+(null==r?a-e:r)+"px;\n                             height: "+(n-t)+"px"))}function h(t,r,n){var o,l,h=Xe(i,t),f=h.text.length;function d(r,n){return jr(e,et(t,r),"div",h,n)}function p(t,r,n){var i=Zr(e,h,null,t),o="ltr"==r==("after"==n)?"left":"right";return d("after"==n?i.begin:i.end-(/\s/.test(h.text.charAt(i.end-1))?2:1),o)[o]}var g=ue(h,i.direction);return function(e,t,r,n){if(!e)return n(t,r,"ltr",0);for(var i=!1,o=0;o<e.length;++o){var l=e[o];(l.from<r&&l.to>t||t==r&&l.to==t)&&(n(Math.max(l.from,t),Math.min(l.to,r),1==l.level?"rtl":"ltr",o),i=!0)}i||n(t,r,"ltr")}(g,r||0,null==n?f:n,(function(e,t,i,h){var v="ltr"==i,m=d(e,v?"left":"right"),y=d(t-1,v?"right":"left"),b=null==r&&0==e,w=null==n&&t==f,x=0==h,C=!g||h==g.length-1;if(y.top-m.top<=3){var S=(u?w:b)&&C,L=(u?b:w)&&x?s:(v?m:y).left,k=S?a:(v?y:m).right;c(L,m.top,k-L,m.bottom)}else{var T,M,N,O;v?(T=u&&b&&x?s:m.left,M=u?a:p(e,i,"before"),N=u?s:p(t,i,"after"),O=u&&w&&C?a:y.right):(T=u?p(e,i,"before"):s,M=!u&&b&&x?a:m.right,N=!u&&w&&C?s:y.left,O=u?p(t,i,"after"):a),c(T,m.top,M-T,m.bottom),m.bottom<y.top&&c(s,m.bottom,null,y.top),c(N,y.top,O-N,y.bottom)}(!o||bn(m,o)<0)&&(o=m),bn(y,o)<0&&(o=y),(!l||bn(m,l)<0)&&(l=m),bn(y,l)<0&&(l=y)})),{start:o,end:l}}var f=t.from(),d=t.to();if(f.line==d.line)h(f.line,f.ch,d.ch);else{var p=Xe(i,f.line),g=Xe(i,d.line),v=Rt(p)==Rt(g),m=h(f.line,f.ch,v?p.text.length+1:null).end,y=h(d.line,v?0:null,d.ch).start;v&&(m.top<y.top-2?(c(m.right,m.top,null,m.bottom),c(s,y.top,y.left,y.bottom)):c(m.right,m.top,y.left-m.right,m.bottom)),m.bottom<y.top&&c(s,m.bottom,null,y.top)}r.appendChild(o)}function xn(e){if(e.state.focused){var t=e.display;clearInterval(t.blinker);var r=!0;t.cursorDiv.style.visibility="",e.options.cursorBlinkRate>0?t.blinker=setInterval((function(){e.hasFocus()||kn(e),t.cursorDiv.style.visibility=(r=!r)?"":"hidden"}),e.options.cursorBlinkRate):e.options.cursorBlinkRate<0&&(t.cursorDiv.style.visibility="hidden")}}function Cn(e){e.hasFocus()||(e.display.input.focus(),e.state.focused||Ln(e))}function Sn(e){e.state.delayingBlurEvent=!0,setTimeout((function(){e.state.delayingBlurEvent&&(e.state.delayingBlurEvent=!1,e.state.focused&&kn(e))}),100)}function Ln(e,t){e.state.delayingBlurEvent&&!e.state.draggingText&&(e.state.delayingBlurEvent=!1),"nocursor"!=e.options.readOnly&&(e.state.focused||(pe(e,"focus",e,t),e.state.focused=!0,H(e.display.wrapper,"CodeMirror-focused"),e.curOp||e.display.selForContextMenu==e.doc.sel||(e.display.input.reset(),a&&setTimeout((function(){return e.display.input.reset(!0)}),20)),e.display.input.receivedFocus()),xn(e))}function kn(e,t){e.state.delayingBlurEvent||(e.state.focused&&(pe(e,"blur",e,t),e.state.focused=!1,T(e.display.wrapper,"CodeMirror-focused")),clearInterval(e.display.blinker),setTimeout((function(){e.state.focused||(e.display.shift=!1)}),150))}function Tn(e){for(var t=e.display,r=t.lineDiv.offsetTop,n=Math.max(0,t.scroller.getBoundingClientRect().top),i=t.lineDiv.getBoundingClientRect().top,o=0,a=0;a<t.view.length;a++){var u=t.view[a],c=e.options.lineWrapping,h=void 0,f=0;if(!u.hidden){if(i+=u.line.height,l&&s<8){var d=u.node.offsetTop+u.node.offsetHeight;h=d-r,r=d}else{var p=u.node.getBoundingClientRect();h=p.bottom-p.top,!c&&u.text.firstChild&&(f=u.text.firstChild.getBoundingClientRect().right-p.left-1)}var g=u.line.height-h;if((g>.005||g<-.005)&&(i<n&&(o-=g),_e(u.line,h),Mn(u.line),u.rest))for(var v=0;v<u.rest.length;v++)Mn(u.rest[v]);if(f>e.display.sizerWidth){var m=Math.ceil(f/nn(e.display));m>e.display.maxLineLength&&(e.display.maxLineLength=m,e.display.maxLine=u.line,e.display.maxLineChanged=!0)}}}Math.abs(o)>2&&(t.scroller.scrollTop+=o)}function Mn(e){if(e.widgets)for(var t=0;t<e.widgets.length;++t){var r=e.widgets[t],n=r.node.parentNode;n&&(r.height=n.offsetHeight)}}function Nn(e,t,r){var n=r&&null!=r.top?Math.max(0,r.top):e.scroller.scrollTop;n=Math.floor(n-Cr(e));var i=r&&null!=r.bottom?r.bottom:n+e.wrapper.clientHeight,o=Ze(t,n),l=Ze(t,i);if(r&&r.ensure){var s=r.ensure.from.line,a=r.ensure.to.line;s<o?(o=s,l=Ze(t,Vt(Xe(t,s))+e.wrapper.clientHeight)):Math.min(a,t.lastLine())>=l&&(o=Ze(t,Vt(Xe(t,a))-e.wrapper.clientHeight),l=a)}return{from:o,to:Math.max(l,o+1)}}function On(e,t){var r=e.display,n=rn(e.display);t.top<0&&(t.top=0);var i=e.curOp&&null!=e.curOp.scrollTop?e.curOp.scrollTop:r.scroller.scrollTop,o=Mr(e),l={};t.bottom-t.top>o&&(t.bottom=t.top+o);var s=e.doc.height+Sr(r),a=t.top<n,u=t.bottom>s-n;if(t.top<i)l.scrollTop=a?0:t.top;else if(t.bottom>i+o){var c=Math.min(t.top,(u?s:t.bottom)-o);c!=i&&(l.scrollTop=c)}var h=e.options.fixedGutter?0:r.gutters.offsetWidth,f=e.curOp&&null!=e.curOp.scrollLeft?e.curOp.scrollLeft:r.scroller.scrollLeft-h,d=Tr(e)-r.gutters.offsetWidth,p=t.right-t.left>d;return p&&(t.right=t.left+d),t.left<10?l.scrollLeft=0:t.left<f?l.scrollLeft=Math.max(0,t.left+h-(p?0:10)):t.right>d+f-3&&(l.scrollLeft=t.right+(p?0:10)-d),l}function An(e,t){null!=t&&(Hn(e),e.curOp.scrollTop=(null==e.curOp.scrollTop?e.doc.scrollTop:e.curOp.scrollTop)+t)}function Dn(e){Hn(e);var t=e.getCursor();e.curOp.scrollToPos={from:t,to:t,margin:e.options.cursorScrollMargin}}function Wn(e,t,r){null==t&&null==r||Hn(e),null!=t&&(e.curOp.scrollLeft=t),null!=r&&(e.curOp.scrollTop=r)}function Hn(e){var t=e.curOp.scrollToPos;t&&(e.curOp.scrollToPos=null,Fn(e,Yr(e,t.from),Yr(e,t.to),t.margin))}function Fn(e,t,r,n){var i=On(e,{left:Math.min(t.left,r.left),top:Math.min(t.top,r.top)-n,right:Math.max(t.right,r.right),bottom:Math.max(t.bottom,r.bottom)+n});Wn(e,i.scrollLeft,i.scrollTop)}function Pn(e,t){Math.abs(e.doc.scrollTop-t)<2||(r||ai(e,{top:t}),En(e,t,!0),r&&ai(e),ni(e,100))}function En(e,t,r){t=Math.max(0,Math.min(e.display.scroller.scrollHeight-e.display.scroller.clientHeight,t)),(e.display.scroller.scrollTop!=t||r)&&(e.doc.scrollTop=t,e.display.scrollbars.setScrollTop(t),e.display.scroller.scrollTop!=t&&(e.display.scroller.scrollTop=t))}function In(e,t,r,n){t=Math.max(0,Math.min(t,e.display.scroller.scrollWidth-e.display.scroller.clientWidth)),(r?t==e.doc.scrollLeft:Math.abs(e.doc.scrollLeft-t)<2)&&!n||(e.doc.scrollLeft=t,hi(e),e.display.scroller.scrollLeft!=t&&(e.display.scroller.scrollLeft=t),e.display.scrollbars.setScrollLeft(t))}function Rn(e){var t=e.display,r=t.gutters.offsetWidth,n=Math.round(e.doc.height+Sr(e.display));return{clientHeight:t.scroller.clientHeight,viewHeight:t.wrapper.clientHeight,scrollWidth:t.scroller.scrollWidth,clientWidth:t.scroller.clientWidth,viewWidth:t.wrapper.clientWidth,barLeft:e.options.fixedGutter?r:0,docHeight:n,scrollHeight:n+kr(e)+t.barHeight,nativeBarWidth:t.nativeBarWidth,gutterWidth:r}}var zn=function(e,t,r){this.cm=r;var n=this.vert=O("div",[O("div",null,null,"min-width: 1px")],"CodeMirror-vscrollbar"),i=this.horiz=O("div",[O("div",null,null,"height: 100%; min-height: 1px")],"CodeMirror-hscrollbar");n.tabIndex=i.tabIndex=-1,e(n),e(i),he(n,"scroll",(function(){n.clientHeight&&t(n.scrollTop,"vertical")})),he(i,"scroll",(function(){i.clientWidth&&t(i.scrollLeft,"horizontal")})),this.checkedZeroWidth=!1,l&&s<8&&(this.horiz.style.minHeight=this.vert.style.minWidth="18px")};zn.prototype.update=function(e){var t=e.scrollWidth>e.clientWidth+1,r=e.scrollHeight>e.clientHeight+1,n=e.nativeBarWidth;if(r){this.vert.style.display="block",this.vert.style.bottom=t?n+"px":"0";var i=e.viewHeight-(t?n:0);this.vert.firstChild.style.height=Math.max(0,e.scrollHeight-e.clientHeight+i)+"px"}else this.vert.scrollTop=0,this.vert.style.display="",this.vert.firstChild.style.height="0";if(t){this.horiz.style.display="block",this.horiz.style.right=r?n+"px":"0",this.horiz.style.left=e.barLeft+"px";var o=e.viewWidth-e.barLeft-(r?n:0);this.horiz.firstChild.style.width=Math.max(0,e.scrollWidth-e.clientWidth+o)+"px"}else this.horiz.style.display="",this.horiz.firstChild.style.width="0";return!this.checkedZeroWidth&&e.clientHeight>0&&(0==n&&this.zeroWidthHack(),this.checkedZeroWidth=!0),{right:r?n:0,bottom:t?n:0}},zn.prototype.setScrollLeft=function(e){this.horiz.scrollLeft!=e&&(this.horiz.scrollLeft=e),this.disableHoriz&&this.enableZeroWidthBar(this.horiz,this.disableHoriz,"horiz")},zn.prototype.setScrollTop=function(e){this.vert.scrollTop!=e&&(this.vert.scrollTop=e),this.disableVert&&this.enableZeroWidthBar(this.vert,this.disableVert,"vert")},zn.prototype.zeroWidthHack=function(){var e=y&&!d?"12px":"18px";this.horiz.style.height=this.vert.style.width=e,this.horiz.style.pointerEvents=this.vert.style.pointerEvents="none",this.disableHoriz=new z,this.disableVert=new z},zn.prototype.enableZeroWidthBar=function(e,t,r){e.style.pointerEvents="auto",t.set(1e3,(function n(){var i=e.getBoundingClientRect();("vert"==r?document.elementFromPoint(i.right-1,(i.top+i.bottom)/2):document.elementFromPoint((i.right+i.left)/2,i.bottom-1))!=e?e.style.pointerEvents="none":t.set(1e3,n)}))},zn.prototype.clear=function(){var e=this.horiz.parentNode;e.removeChild(this.horiz),e.removeChild(this.vert)};var Bn=function(){};function Gn(e,t){t||(t=Rn(e));var r=e.display.barWidth,n=e.display.barHeight;Un(e,t);for(var i=0;i<4&&r!=e.display.barWidth||n!=e.display.barHeight;i++)r!=e.display.barWidth&&e.options.lineWrapping&&Tn(e),Un(e,Rn(e)),r=e.display.barWidth,n=e.display.barHeight}function Un(e,t){var r=e.display,n=r.scrollbars.update(t);r.sizer.style.paddingRight=(r.barWidth=n.right)+"px",r.sizer.style.paddingBottom=(r.barHeight=n.bottom)+"px",r.heightForcer.style.borderBottom=n.bottom+"px solid transparent",n.right&&n.bottom?(r.scrollbarFiller.style.display="block",r.scrollbarFiller.style.height=n.bottom+"px",r.scrollbarFiller.style.width=n.right+"px"):r.scrollbarFiller.style.display="",n.bottom&&e.options.coverGutterNextToScrollbar&&e.options.fixedGutter?(r.gutterFiller.style.display="block",r.gutterFiller.style.height=n.bottom+"px",r.gutterFiller.style.width=t.gutterWidth+"px"):r.gutterFiller.style.display=""}Bn.prototype.update=function(){return{bottom:0,right:0}},Bn.prototype.setScrollLeft=function(){},Bn.prototype.setScrollTop=function(){},Bn.prototype.clear=function(){};var Vn={native:zn,null:Bn};function Kn(e){e.display.scrollbars&&(e.display.scrollbars.clear(),e.display.scrollbars.addClass&&T(e.display.wrapper,e.display.scrollbars.addClass)),e.display.scrollbars=new Vn[e.options.scrollbarStyle]((function(t){e.display.wrapper.insertBefore(t,e.display.scrollbarFiller),he(t,"mousedown",(function(){e.state.focused&&setTimeout((function(){return e.display.input.focus()}),0)})),t.setAttribute("cm-not-content","true")}),(function(t,r){"horizontal"==r?In(e,t):Pn(e,t)}),e),e.display.scrollbars.addClass&&H(e.display.wrapper,e.display.scrollbars.addClass)}var jn=0;function Xn(e){var t;e.curOp={cm:e,viewChanged:!1,startHeight:e.doc.height,forceUpdate:!1,updateInput:0,typing:!1,changeObjs:null,cursorActivityHandlers:null,cursorActivityCalled:0,selectionChanged:!1,updateMaxLine:!1,scrollLeft:null,scrollTop:null,scrollToPos:null,focus:!1,id:++jn,markArrays:null},t=e.curOp,or?or.ops.push(t):t.ownsGroup=or={ops:[t],delayedCallbacks:[]}}function Yn(e){var t=e.curOp;t&&function(e,t){var r=e.ownsGroup;if(r)try{!function(e){var t=e.delayedCallbacks,r=0;do{for(;r<t.length;r++)t[r].call(null);for(var n=0;n<e.ops.length;n++){var i=e.ops[n];if(i.cursorActivityHandlers)for(;i.cursorActivityCalled<i.cursorActivityHandlers.length;)i.cursorActivityHandlers[i.cursorActivityCalled++].call(null,i.cm)}}while(r<t.length)}(r)}finally{or=null,t(r)}}(t,(function(e){for(var t=0;t<e.ops.length;t++)e.ops[t].cm.curOp=null;!function(e){for(var t=e.ops,r=0;r<t.length;r++)$n(t[r]);for(var n=0;n<t.length;n++)_n(t[n]);for(var i=0;i<t.length;i++)qn(t[i]);for(var o=0;o<t.length;o++)Zn(t[o]);for(var l=0;l<t.length;l++)Qn(t[l])}(e)}))}function $n(e){var t=e.cm,r=t.display;!function(e){var t=e.display;!t.scrollbarsClipped&&t.scroller.offsetWidth&&(t.nativeBarWidth=t.scroller.offsetWidth-t.scroller.clientWidth,t.heightForcer.style.height=kr(e)+"px",t.sizer.style.marginBottom=-t.nativeBarWidth+"px",t.sizer.style.borderRightWidth=kr(e)+"px",t.scrollbarsClipped=!0)}(t),e.updateMaxLine&&jt(t),e.mustUpdate=e.viewChanged||e.forceUpdate||null!=e.scrollTop||e.scrollToPos&&(e.scrollToPos.from.line<r.viewFrom||e.scrollToPos.to.line>=r.viewTo)||r.maxLineChanged&&t.options.lineWrapping,e.update=e.mustUpdate&&new oi(t,e.mustUpdate&&{top:e.scrollTop,ensure:e.scrollToPos},e.forceUpdate)}function _n(e){e.updatedDisplay=e.mustUpdate&&li(e.cm,e.update)}function qn(e){var t=e.cm,r=t.display;e.updatedDisplay&&Tn(t),e.barMeasure=Rn(t),r.maxLineChanged&&!t.options.lineWrapping&&(e.adjustWidthTo=Or(t,r.maxLine,r.maxLine.text.length).left+3,t.display.sizerWidth=e.adjustWidthTo,e.barMeasure.scrollWidth=Math.max(r.scroller.clientWidth,r.sizer.offsetLeft+e.adjustWidthTo+kr(t)+t.display.barWidth),e.maxScrollLeft=Math.max(0,r.sizer.offsetLeft+e.adjustWidthTo-Tr(t))),(e.updatedDisplay||e.selectionChanged)&&(e.preparedSelection=r.input.prepareSelection())}function Zn(e){var t=e.cm;null!=e.adjustWidthTo&&(t.display.sizer.style.minWidth=e.adjustWidthTo+"px",e.maxScrollLeft<t.doc.scrollLeft&&In(t,Math.min(t.display.scroller.scrollLeft,e.maxScrollLeft),!0),t.display.maxLineChanged=!1);var r=e.focus&&e.focus==W();e.preparedSelection&&t.display.input.showSelection(e.preparedSelection,r),(e.updatedDisplay||e.startHeight!=t.doc.height)&&Gn(t,e.barMeasure),e.updatedDisplay&&ci(t,e.barMeasure),e.selectionChanged&&xn(t),t.state.focused&&e.updateInput&&t.display.input.reset(e.typing),r&&Cn(e.cm)}function Qn(e){var t=e.cm,r=t.display,n=t.doc;if(e.updatedDisplay&&si(t,e.update),null==r.wheelStartX||null==e.scrollTop&&null==e.scrollLeft&&!e.scrollToPos||(r.wheelStartX=r.wheelStartY=null),null!=e.scrollTop&&En(t,e.scrollTop,e.forceScroll),null!=e.scrollLeft&&In(t,e.scrollLeft,!0,!0),e.scrollToPos){var i=function(e,t,r,n){var i;null==n&&(n=0),e.options.lineWrapping||t!=r||(r="before"==t.sticky?et(t.line,t.ch+1,"before"):t,t=t.ch?et(t.line,"before"==t.sticky?t.ch-1:t.ch,"after"):t);for(var o=0;o<5;o++){var l=!1,s=Xr(e,t),a=r&&r!=t?Xr(e,r):s,u=On(e,i={left:Math.min(s.left,a.left),top:Math.min(s.top,a.top)-n,right:Math.max(s.left,a.left),bottom:Math.max(s.bottom,a.bottom)+n}),c=e.doc.scrollTop,h=e.doc.scrollLeft;if(null!=u.scrollTop&&(Pn(e,u.scrollTop),Math.abs(e.doc.scrollTop-c)>1&&(l=!0)),null!=u.scrollLeft&&(In(e,u.scrollLeft),Math.abs(e.doc.scrollLeft-h)>1&&(l=!0)),!l)break}return i}(t,st(n,e.scrollToPos.from),st(n,e.scrollToPos.to),e.scrollToPos.margin);!function(e,t){if(!ge(e,"scrollCursorIntoView")){var r=e.display,n=r.sizer.getBoundingClientRect(),i=null;if(t.top+n.top<0?i=!0:t.bottom+n.top>(window.innerHeight||document.documentElement.clientHeight)&&(i=!1),null!=i&&!p){var o=O("div","​",null,"position: absolute;\n                         top: "+(t.top-r.viewOffset-Cr(e.display))+"px;\n                         height: "+(t.bottom-t.top+kr(e)+r.barHeight)+"px;\n                         left: "+t.left+"px; width: "+Math.max(2,t.right-t.left)+"px;");e.display.lineSpace.appendChild(o),o.scrollIntoView(i),e.display.lineSpace.removeChild(o)}}}(t,i)}var o=e.maybeHiddenMarkers,l=e.maybeUnhiddenMarkers;if(o)for(var s=0;s<o.length;++s)o[s].lines.length||pe(o[s],"hide");if(l)for(var a=0;a<l.length;++a)l[a].lines.length&&pe(l[a],"unhide");r.wrapper.offsetHeight&&(n.scrollTop=t.display.scroller.scrollTop),e.changeObjs&&pe(t,"changes",t,e.changeObjs),e.update&&e.update.finish()}function Jn(e,t){if(e.curOp)return t();Xn(e);try{return t()}finally{Yn(e)}}function ei(e,t){return function(){if(e.curOp)return t.apply(e,arguments);Xn(e);try{return t.apply(e,arguments)}finally{Yn(e)}}}function ti(e){return function(){if(this.curOp)return e.apply(this,arguments);Xn(this);try{return e.apply(this,arguments)}finally{Yn(this)}}}function ri(e){return function(){var t=this.cm;if(!t||t.curOp)return e.apply(this,arguments);Xn(t);try{return e.apply(this,arguments)}finally{Yn(t)}}}function ni(e,t){e.doc.highlightFrontier<e.display.viewTo&&e.state.highlight.set(t,E(ii,e))}function ii(e){var t=e.doc;if(!(t.highlightFrontier>=e.display.viewTo)){var r=+new Date+e.options.workTime,n=dt(e,t.highlightFrontier),i=[];t.iter(n.line,Math.min(t.first+t.size,e.display.viewTo+500),(function(o){if(n.line>=e.display.viewFrom){var l=o.styles,s=o.text.length>e.options.maxHighlightLength?Ue(t.mode,n.state):null,a=ht(e,o,n,!0);s&&(n.state=s),o.styles=a.styles;var u=o.styleClasses,c=a.classes;c?o.styleClasses=c:u&&(o.styleClasses=null);for(var h=!l||l.length!=o.styles.length||u!=c&&(!u||!c||u.bgClass!=c.bgClass||u.textClass!=c.textClass),f=0;!h&&f<l.length;++f)h=l[f]!=o.styles[f];h&&i.push(n.line),o.stateAfter=n.save(),n.nextLine()}else o.text.length<=e.options.maxHighlightLength&&pt(e,o.text,n),o.stateAfter=n.line%5==0?n.save():null,n.nextLine();if(+new Date>r)return ni(e,e.options.workDelay),!0})),t.highlightFrontier=n.line,t.modeFrontier=Math.max(t.modeFrontier,n.line),i.length&&Jn(e,(function(){for(var t=0;t<i.length;t++)fn(e,i[t],"text")}))}}var oi=function(e,t,r){var n=e.display;this.viewport=t,this.visible=Nn(n,e.doc,t),this.editorIsHidden=!n.wrapper.offsetWidth,this.wrapperHeight=n.wrapper.clientHeight,this.wrapperWidth=n.wrapper.clientWidth,this.oldDisplayWidth=Tr(e),this.force=r,this.dims=on(e),this.events=[]};function li(e,t){var r=e.display,n=e.doc;if(t.editorIsHidden)return dn(e),!1;if(!t.force&&t.visible.from>=r.viewFrom&&t.visible.to<=r.viewTo&&(null==r.updateLineNumbers||r.updateLineNumbers>=r.viewTo)&&r.renderedView==r.view&&0==gn(e))return!1;fi(e)&&(dn(e),t.dims=on(e));var i=n.first+n.size,o=Math.max(t.visible.from-e.options.viewportMargin,n.first),l=Math.min(i,t.visible.to+e.options.viewportMargin);r.viewFrom<o&&o-r.viewFrom<20&&(o=Math.max(n.first,r.viewFrom)),r.viewTo>l&&r.viewTo-l<20&&(l=Math.min(i,r.viewTo)),Ct&&(o=zt(e.doc,o),l=Bt(e.doc,l));var s=o!=r.viewFrom||l!=r.viewTo||r.lastWrapHeight!=t.wrapperHeight||r.lastWrapWidth!=t.wrapperWidth;!function(e,t,r){var n=e.display;0==n.view.length||t>=n.viewTo||r<=n.viewFrom?(n.view=ir(e,t,r),n.viewFrom=t):(n.viewFrom>t?n.view=ir(e,t,n.viewFrom).concat(n.view):n.viewFrom<t&&(n.view=n.view.slice(cn(e,t))),n.viewFrom=t,n.viewTo<r?n.view=n.view.concat(ir(e,n.viewTo,r)):n.viewTo>r&&(n.view=n.view.slice(0,cn(e,r)))),n.viewTo=r}(e,o,l),r.viewOffset=Vt(Xe(e.doc,r.viewFrom)),e.display.mover.style.top=r.viewOffset+"px";var u=gn(e);if(!s&&0==u&&!t.force&&r.renderedView==r.view&&(null==r.updateLineNumbers||r.updateLineNumbers>=r.viewTo))return!1;var c=function(e){if(e.hasFocus())return null;var t=W();if(!t||!D(e.display.lineDiv,t))return null;var r={activeElt:t};if(window.getSelection){var n=window.getSelection();n.anchorNode&&n.extend&&D(e.display.lineDiv,n.anchorNode)&&(r.anchorNode=n.anchorNode,r.anchorOffset=n.anchorOffset,r.focusNode=n.focusNode,r.focusOffset=n.focusOffset)}return r}(e);return u>4&&(r.lineDiv.style.display="none"),function(e,t,r){var n=e.display,i=e.options.lineNumbers,o=n.lineDiv,l=o.firstChild;function s(t){var r=t.nextSibling;return a&&y&&e.display.currentWheelTarget==t?t.style.display="none":t.parentNode.removeChild(t),r}for(var u=n.view,c=n.viewFrom,h=0;h<u.length;h++){var f=u[h];if(f.hidden);else if(f.node&&f.node.parentNode==o){for(;l!=f.node;)l=s(l);var d=i&&null!=t&&t<=c&&f.lineNumber;f.changes&&(B(f.changes,"gutter")>-1&&(d=!1),ur(e,f,c,r)),d&&(M(f.lineNumber),f.lineNumber.appendChild(document.createTextNode(Je(e.options,c)))),l=f.node.nextSibling}else{var p=vr(e,f,c,r);o.insertBefore(p,l)}c+=f.size}for(;l;)l=s(l)}(e,r.updateLineNumbers,t.dims),u>4&&(r.lineDiv.style.display=""),r.renderedView=r.view,function(e){if(e&&e.activeElt&&e.activeElt!=W()&&(e.activeElt.focus(),!/^(INPUT|TEXTAREA)$/.test(e.activeElt.nodeName)&&e.anchorNode&&D(document.body,e.anchorNode)&&D(document.body,e.focusNode))){var t=window.getSelection(),r=document.createRange();r.setEnd(e.anchorNode,e.anchorOffset),r.collapse(!1),t.removeAllRanges(),t.addRange(r),t.extend(e.focusNode,e.focusOffset)}}(c),M(r.cursorDiv),M(r.selectionDiv),r.gutters.style.height=r.sizer.style.minHeight=0,s&&(r.lastWrapHeight=t.wrapperHeight,r.lastWrapWidth=t.wrapperWidth,ni(e,400)),r.updateLineNumbers=null,!0}function si(e,t){for(var r=t.viewport,n=!0;;n=!1){if(n&&e.options.lineWrapping&&t.oldDisplayWidth!=Tr(e))n&&(t.visible=Nn(e.display,e.doc,r));else if(r&&null!=r.top&&(r={top:Math.min(e.doc.height+Sr(e.display)-Mr(e),r.top)}),t.visible=Nn(e.display,e.doc,r),t.visible.from>=e.display.viewFrom&&t.visible.to<=e.display.viewTo)break;if(!li(e,t))break;Tn(e);var i=Rn(e);vn(e),Gn(e,i),ci(e,i),t.force=!1}t.signal(e,"update",e),e.display.viewFrom==e.display.reportedViewFrom&&e.display.viewTo==e.display.reportedViewTo||(t.signal(e,"viewportChange",e,e.display.viewFrom,e.display.viewTo),e.display.reportedViewFrom=e.display.viewFrom,e.display.reportedViewTo=e.display.viewTo)}function ai(e,t){var r=new oi(e,t);if(li(e,r)){Tn(e),si(e,r);var n=Rn(e);vn(e),Gn(e,n),ci(e,n),r.finish()}}function ui(e){var t=e.gutters.offsetWidth;e.sizer.style.marginLeft=t+"px",sr(e,"gutterChanged",e)}function ci(e,t){e.display.sizer.style.minHeight=t.docHeight+"px",e.display.heightForcer.style.top=t.docHeight+"px",e.display.gutters.style.height=t.docHeight+e.display.barHeight+kr(e)+"px"}function hi(e){var t=e.display,r=t.view;if(t.alignWidgets||t.gutters.firstChild&&e.options.fixedGutter){for(var n=ln(t)-t.scroller.scrollLeft+e.doc.scrollLeft,i=t.gutters.offsetWidth,o=n+"px",l=0;l<r.length;l++)if(!r[l].hidden){e.options.fixedGutter&&(r[l].gutter&&(r[l].gutter.style.left=o),r[l].gutterBackground&&(r[l].gutterBackground.style.left=o));var s=r[l].alignable;if(s)for(var a=0;a<s.length;a++)s[a].style.left=o}e.options.fixedGutter&&(t.gutters.style.left=n+i+"px")}}function fi(e){if(!e.options.lineNumbers)return!1;var t=e.doc,r=Je(e.options,t.first+t.size-1),n=e.display;if(r.length!=n.lineNumChars){var i=n.measure.appendChild(O("div",[O("div",r)],"CodeMirror-linenumber CodeMirror-gutter-elt")),o=i.firstChild.offsetWidth,l=i.offsetWidth-o;return n.lineGutter.style.width="",n.lineNumInnerWidth=Math.max(o,n.lineGutter.offsetWidth-l)+1,n.lineNumWidth=n.lineNumInnerWidth+l,n.lineNumChars=n.lineNumInnerWidth?r.length:-1,n.lineGutter.style.width=n.lineNumWidth+"px",ui(e.display),!0}return!1}function di(e,t){for(var r=[],n=!1,i=0;i<e.length;i++){var o=e[i],l=null;if("string"!=typeof o&&(l=o.style,o=o.className),"CodeMirror-linenumbers"==o){if(!t)continue;n=!0}r.push({className:o,style:l})}return t&&!n&&r.push({className:"CodeMirror-linenumbers",style:null}),r}function pi(e){var t=e.gutters,r=e.gutterSpecs;M(t),e.lineGutter=null;for(var n=0;n<r.length;++n){var i=r[n],o=i.className,l=i.style,s=t.appendChild(O("div",null,"CodeMirror-gutter "+o));l&&(s.style.cssText=l),"CodeMirror-linenumbers"==o&&(e.lineGutter=s,s.style.width=(e.lineNumWidth||1)+"px")}t.style.display=r.length?"":"none",ui(e)}function gi(e){pi(e.display),hn(e),hi(e)}function vi(e,t,n,i){var o=this;this.input=n,o.scrollbarFiller=O("div",null,"CodeMirror-scrollbar-filler"),o.scrollbarFiller.setAttribute("cm-not-content","true"),o.gutterFiller=O("div",null,"CodeMirror-gutter-filler"),o.gutterFiller.setAttribute("cm-not-content","true"),o.lineDiv=A("div",null,"CodeMirror-code"),o.selectionDiv=O("div",null,null,"position: relative; z-index: 1"),o.cursorDiv=O("div",null,"CodeMirror-cursors"),o.measure=O("div",null,"CodeMirror-measure"),o.lineMeasure=O("div",null,"CodeMirror-measure"),o.lineSpace=A("div",[o.measure,o.lineMeasure,o.selectionDiv,o.cursorDiv,o.lineDiv],null,"position: relative; outline: none");var u=A("div",[o.lineSpace],"CodeMirror-lines");o.mover=O("div",[u],null,"position: relative"),o.sizer=O("div",[o.mover],"CodeMirror-sizer"),o.sizerWidth=null,o.heightForcer=O("div",null,null,"position: absolute; height: 50px; width: 1px;"),o.gutters=O("div",null,"CodeMirror-gutters"),o.lineGutter=null,o.scroller=O("div",[o.sizer,o.heightForcer,o.gutters],"CodeMirror-scroll"),o.scroller.setAttribute("tabIndex","-1"),o.wrapper=O("div",[o.scrollbarFiller,o.gutterFiller,o.scroller],"CodeMirror"),o.wrapper.setAttribute("translate","no"),l&&s<8&&(o.gutters.style.zIndex=-1,o.scroller.style.paddingRight=0),a||r&&m||(o.scroller.draggable=!0),e&&(e.appendChild?e.appendChild(o.wrapper):e(o.wrapper)),o.viewFrom=o.viewTo=t.first,o.reportedViewFrom=o.reportedViewTo=t.first,o.view=[],o.renderedView=null,o.externalMeasured=null,o.viewOffset=0,o.lastWrapHeight=o.lastWrapWidth=0,o.updateLineNumbers=null,o.nativeBarWidth=o.barHeight=o.barWidth=0,o.scrollbarsClipped=!1,o.lineNumWidth=o.lineNumInnerWidth=o.lineNumChars=null,o.alignWidgets=!1,o.cachedCharWidth=o.cachedTextHeight=o.cachedPaddingH=null,o.maxLine=null,o.maxLineLength=0,o.maxLineChanged=!1,o.wheelDX=o.wheelDY=o.wheelStartX=o.wheelStartY=null,o.shift=!1,o.selForContextMenu=null,o.activeTouch=null,o.gutterSpecs=di(i.gutters,i.lineNumbers),pi(o),n.init(o)}oi.prototype.signal=function(e,t){me(e,t)&&this.events.push(arguments)},oi.prototype.finish=function(){for(var e=0;e<this.events.length;e++)pe.apply(null,this.events[e])};var mi=0,yi=null;function bi(e){var t=e.wheelDeltaX,r=e.wheelDeltaY;return null==t&&e.detail&&e.axis==e.HORIZONTAL_AXIS&&(t=e.detail),null==r&&e.detail&&e.axis==e.VERTICAL_AXIS?r=e.detail:null==r&&(r=e.wheelDelta),{x:t,y:r}}function wi(e){var t=bi(e);return t.x*=yi,t.y*=yi,t}function xi(e,t){var n=bi(t),i=n.x,o=n.y,l=yi;0===t.deltaMode&&(i=t.deltaX,o=t.deltaY,l=1);var s=e.display,u=s.scroller,c=u.scrollWidth>u.clientWidth,f=u.scrollHeight>u.clientHeight;if(i&&c||o&&f){if(o&&y&&a)e:for(var d=t.target,p=s.view;d!=u;d=d.parentNode)for(var g=0;g<p.length;g++)if(p[g].node==d){e.display.currentWheelTarget=d;break e}if(i&&!r&&!h&&null!=l)return o&&f&&Pn(e,Math.max(0,u.scrollTop+o*l)),In(e,Math.max(0,u.scrollLeft+i*l)),(!o||o&&f)&&be(t),void(s.wheelStartX=null);if(o&&null!=l){var v=o*l,m=e.doc.scrollTop,b=m+s.wrapper.clientHeight;v<0?m=Math.max(0,m+v-50):b=Math.min(e.doc.height,b+v+50),ai(e,{top:m,bottom:b})}mi<20&&0!==t.deltaMode&&(null==s.wheelStartX?(s.wheelStartX=u.scrollLeft,s.wheelStartY=u.scrollTop,s.wheelDX=i,s.wheelDY=o,setTimeout((function(){if(null!=s.wheelStartX){var e=u.scrollLeft-s.wheelStartX,t=u.scrollTop-s.wheelStartY,r=t&&s.wheelDY&&t/s.wheelDY||e&&s.wheelDX&&e/s.wheelDX;s.wheelStartX=s.wheelStartY=null,r&&(yi=(yi*mi+r)/(mi+1),++mi)}}),200)):(s.wheelDX+=i,s.wheelDY+=o))}}l?yi=-.53:r?yi=15:c?yi=-.7:f&&(yi=-1/3);var Ci=function(e,t){this.ranges=e,this.primIndex=t};Ci.prototype.primary=function(){return this.ranges[this.primIndex]},Ci.prototype.equals=function(e){if(e==this)return!0;if(e.primIndex!=this.primIndex||e.ranges.length!=this.ranges.length)return!1;for(var t=0;t<this.ranges.length;t++){var r=this.ranges[t],n=e.ranges[t];if(!rt(r.anchor,n.anchor)||!rt(r.head,n.head))return!1}return!0},Ci.prototype.deepCopy=function(){for(var e=[],t=0;t<this.ranges.length;t++)e[t]=new Si(nt(this.ranges[t].anchor),nt(this.ranges[t].head));return new Ci(e,this.primIndex)},Ci.prototype.somethingSelected=function(){for(var e=0;e<this.ranges.length;e++)if(!this.ranges[e].empty())return!0;return!1},Ci.prototype.contains=function(e,t){t||(t=e);for(var r=0;r<this.ranges.length;r++){var n=this.ranges[r];if(tt(t,n.from())>=0&&tt(e,n.to())<=0)return r}return-1};var Si=function(e,t){this.anchor=e,this.head=t};function Li(e,t,r){var n=e&&e.options.selectionsMayTouch,i=t[r];t.sort((function(e,t){return tt(e.from(),t.from())})),r=B(t,i);for(var o=1;o<t.length;o++){var l=t[o],s=t[o-1],a=tt(s.to(),l.from());if(n&&!l.empty()?a>0:a>=0){var u=ot(s.from(),l.from()),c=it(s.to(),l.to()),h=s.empty()?l.from()==l.head:s.from()==s.head;o<=r&&--r,t.splice(--o,2,new Si(h?c:u,h?u:c))}}return new Ci(t,r)}function ki(e,t){return new Ci([new Si(e,t||e)],0)}function Ti(e){return e.text?et(e.from.line+e.text.length-1,$(e.text).length+(1==e.text.length?e.from.ch:0)):e.to}function Mi(e,t){if(tt(e,t.from)<0)return e;if(tt(e,t.to)<=0)return Ti(t);var r=e.line+t.text.length-(t.to.line-t.from.line)-1,n=e.ch;return e.line==t.to.line&&(n+=Ti(t).ch-t.to.ch),et(r,n)}function Ni(e,t){for(var r=[],n=0;n<e.sel.ranges.length;n++){var i=e.sel.ranges[n];r.push(new Si(Mi(i.anchor,t),Mi(i.head,t)))}return Li(e.cm,r,e.sel.primIndex)}function Oi(e,t,r){return e.line==t.line?et(r.line,e.ch-t.ch+r.ch):et(r.line+(e.line-t.line),e.ch)}function Ai(e){e.doc.mode=ze(e.options,e.doc.modeOption),Di(e)}function Di(e){e.doc.iter((function(e){e.stateAfter&&(e.stateAfter=null),e.styles&&(e.styles=null)})),e.doc.modeFrontier=e.doc.highlightFrontier=e.doc.first,ni(e,100),e.state.modeGen++,e.curOp&&hn(e)}function Wi(e,t){return 0==t.from.ch&&0==t.to.ch&&""==$(t.text)&&(!e.cm||e.cm.options.wholeLineUpdateBefore)}function Hi(e,t,r,n){function i(e){return r?r[e]:null}function o(e,r,i){!function(e,t,r,n){e.text=t,e.stateAfter&&(e.stateAfter=null),e.styles&&(e.styles=null),null!=e.order&&(e.order=null),Nt(e),Ot(e,r);var i=n?n(e):1;i!=e.height&&_e(e,i)}(e,r,i,n),sr(e,"change",e,t)}function l(e,t){for(var r=[],o=e;o<t;++o)r.push(new Xt(u[o],i(o),n));return r}var s=t.from,a=t.to,u=t.text,c=Xe(e,s.line),h=Xe(e,a.line),f=$(u),d=i(u.length-1),p=a.line-s.line;if(t.full)e.insert(0,l(0,u.length)),e.remove(u.length,e.size-u.length);else if(Wi(e,t)){var g=l(0,u.length-1);o(h,h.text,d),p&&e.remove(s.line,p),g.length&&e.insert(s.line,g)}else if(c==h)if(1==u.length)o(c,c.text.slice(0,s.ch)+f+c.text.slice(a.ch),d);else{var v=l(1,u.length-1);v.push(new Xt(f+c.text.slice(a.ch),d,n)),o(c,c.text.slice(0,s.ch)+u[0],i(0)),e.insert(s.line+1,v)}else if(1==u.length)o(c,c.text.slice(0,s.ch)+u[0]+h.text.slice(a.ch),i(0)),e.remove(s.line+1,p);else{o(c,c.text.slice(0,s.ch)+u[0],i(0)),o(h,f+h.text.slice(a.ch),d);var m=l(1,u.length-1);p>1&&e.remove(s.line+1,p-1),e.insert(s.line+1,m)}sr(e,"change",e,t)}function Fi(e,t,r){!function e(n,i,o){if(n.linked)for(var l=0;l<n.linked.length;++l){var s=n.linked[l];if(s.doc!=i){var a=o&&s.sharedHist;r&&!a||(t(s.doc,a),e(s.doc,n,a))}}}(e,null,!0)}function Pi(e,t){if(t.cm)throw new Error("This document is already in use.");e.doc=t,t.cm=e,an(e),Ai(e),Ei(e),e.options.direction=t.direction,e.options.lineWrapping||jt(e),e.options.mode=t.modeOption,hn(e)}function Ei(e){("rtl"==e.doc.direction?H:T)(e.display.lineDiv,"CodeMirror-rtl")}function Ii(e){this.done=[],this.undone=[],this.undoDepth=e?e.undoDepth:1/0,this.lastModTime=this.lastSelTime=0,this.lastOp=this.lastSelOp=null,this.lastOrigin=this.lastSelOrigin=null,this.generation=this.maxGeneration=e?e.maxGeneration:1}function Ri(e,t){var r={from:nt(t.from),to:Ti(t),text:Ye(e,t.from,t.to)};return Vi(e,r,t.from.line,t.to.line+1),Fi(e,(function(e){return Vi(e,r,t.from.line,t.to.line+1)}),!0),r}function zi(e){for(;e.length;){if(!$(e).ranges)break;e.pop()}}function Bi(e,t,r,n){var i=e.history;i.undone.length=0;var o,l,s=+new Date;if((i.lastOp==n||i.lastOrigin==t.origin&&t.origin&&("+"==t.origin.charAt(0)&&i.lastModTime>s-(e.cm?e.cm.options.historyEventDelay:500)||"*"==t.origin.charAt(0)))&&(o=function(e,t){return t?(zi(e.done),$(e.done)):e.done.length&&!$(e.done).ranges?$(e.done):e.done.length>1&&!e.done[e.done.length-2].ranges?(e.done.pop(),$(e.done)):void 0}(i,i.lastOp==n)))l=$(o.changes),0==tt(t.from,t.to)&&0==tt(t.from,l.to)?l.to=Ti(t):o.changes.push(Ri(e,t));else{var a=$(i.done);for(a&&a.ranges||Ui(e.sel,i.done),o={changes:[Ri(e,t)],generation:i.generation},i.done.push(o);i.done.length>i.undoDepth;)i.done.shift(),i.done[0].ranges||i.done.shift()}i.done.push(r),i.generation=++i.maxGeneration,i.lastModTime=i.lastSelTime=s,i.lastOp=i.lastSelOp=n,i.lastOrigin=i.lastSelOrigin=t.origin,l||pe(e,"historyAdded")}function Gi(e,t,r,n){var i=e.history,o=n&&n.origin;r==i.lastSelOp||o&&i.lastSelOrigin==o&&(i.lastModTime==i.lastSelTime&&i.lastOrigin==o||function(e,t,r,n){var i=t.charAt(0);return"*"==i||"+"==i&&r.ranges.length==n.ranges.length&&r.somethingSelected()==n.somethingSelected()&&new Date-e.history.lastSelTime<=(e.cm?e.cm.options.historyEventDelay:500)}(e,o,$(i.done),t))?i.done[i.done.length-1]=t:Ui(t,i.done),i.lastSelTime=+new Date,i.lastSelOrigin=o,i.lastSelOp=r,n&&!1!==n.clearRedo&&zi(i.undone)}function Ui(e,t){var r=$(t);r&&r.ranges&&r.equals(e)||t.push(e)}function Vi(e,t,r,n){var i=t["spans_"+e.id],o=0;e.iter(Math.max(e.first,r),Math.min(e.first+e.size,n),(function(r){r.markedSpans&&((i||(i=t["spans_"+e.id]={}))[o]=r.markedSpans),++o}))}function Ki(e){if(!e)return null;for(var t,r=0;r<e.length;++r)e[r].marker.explicitlyCleared?t||(t=e.slice(0,r)):t&&t.push(e[r]);return t?t.length?t:null:e}function ji(e,t){var r=function(e,t){var r=t["spans_"+e.id];if(!r)return null;for(var n=[],i=0;i<t.text.length;++i)n.push(Ki(r[i]));return n}(e,t),n=Tt(e,t);if(!r)return n;if(!n)return r;for(var i=0;i<r.length;++i){var o=r[i],l=n[i];if(o&&l)e:for(var s=0;s<l.length;++s){for(var a=l[s],u=0;u<o.length;++u)if(o[u].marker==a.marker)continue e;o.push(a)}else l&&(r[i]=l)}return r}function Xi(e,t,r){for(var n=[],i=0;i<e.length;++i){var o=e[i];if(o.ranges)n.push(r?Ci.prototype.deepCopy.call(o):o);else{var l=o.changes,s=[];n.push({changes:s});for(var a=0;a<l.length;++a){var u=l[a],c=void 0;if(s.push({from:u.from,to:u.to,text:u.text}),t)for(var h in u)(c=h.match(/^spans_(\d+)$/))&&B(t,Number(c[1]))>-1&&($(s)[h]=u[h],delete u[h])}}}return n}function Yi(e,t,r,n){if(n){var i=e.anchor;if(r){var o=tt(t,i)<0;o!=tt(r,i)<0?(i=t,t=r):o!=tt(t,r)<0&&(t=r)}return new Si(i,t)}return new Si(r||t,t)}function $i(e,t,r,n,i){null==i&&(i=e.cm&&(e.cm.display.shift||e.extend)),Ji(e,new Ci([Yi(e.sel.primary(),t,r,i)],0),n)}function _i(e,t,r){for(var n=[],i=e.cm&&(e.cm.display.shift||e.extend),o=0;o<e.sel.ranges.length;o++)n[o]=Yi(e.sel.ranges[o],t[o],null,i);Ji(e,Li(e.cm,n,e.sel.primIndex),r)}function qi(e,t,r,n){var i=e.sel.ranges.slice(0);i[t]=r,Ji(e,Li(e.cm,i,e.sel.primIndex),n)}function Zi(e,t,r,n){Ji(e,ki(t,r),n)}function Qi(e,t,r){var n=e.history.done,i=$(n);i&&i.ranges?(n[n.length-1]=t,eo(e,t,r)):Ji(e,t,r)}function Ji(e,t,r){eo(e,t,r),Gi(e,e.sel,e.cm?e.cm.curOp.id:NaN,r)}function eo(e,t,r){(me(e,"beforeSelectionChange")||e.cm&&me(e.cm,"beforeSelectionChange"))&&(t=function(e,t,r){var n={ranges:t.ranges,update:function(t){this.ranges=[];for(var r=0;r<t.length;r++)this.ranges[r]=new Si(st(e,t[r].anchor),st(e,t[r].head))},origin:r&&r.origin};return pe(e,"beforeSelectionChange",e,n),e.cm&&pe(e.cm,"beforeSelectionChange",e.cm,n),n.ranges!=t.ranges?Li(e.cm,n.ranges,n.ranges.length-1):t}(e,t,r));var n=r&&r.bias||(tt(t.primary().head,e.sel.primary().head)<0?-1:1);to(e,no(e,t,n,!0)),r&&!1===r.scroll||!e.cm||"nocursor"==e.cm.getOption("readOnly")||Dn(e.cm)}function to(e,t){t.equals(e.sel)||(e.sel=t,e.cm&&(e.cm.curOp.updateInput=1,e.cm.curOp.selectionChanged=!0,ve(e.cm)),sr(e,"cursorActivity",e))}function ro(e){to(e,no(e,e.sel,null,!1))}function no(e,t,r,n){for(var i,o=0;o<t.ranges.length;o++){var l=t.ranges[o],s=t.ranges.length==e.sel.ranges.length&&e.sel.ranges[o],a=oo(e,l.anchor,s&&s.anchor,r,n),u=oo(e,l.head,s&&s.head,r,n);(i||a!=l.anchor||u!=l.head)&&(i||(i=t.ranges.slice(0,o)),i[o]=new Si(a,u))}return i?Li(e.cm,i,t.primIndex):t}function io(e,t,r,n,i){var o=Xe(e,t.line);if(o.markedSpans)for(var l=0;l<o.markedSpans.length;++l){var s=o.markedSpans[l],a=s.marker,u="selectLeft"in a?!a.selectLeft:a.inclusiveLeft,c="selectRight"in a?!a.selectRight:a.inclusiveRight;if((null==s.from||(u?s.from<=t.ch:s.from<t.ch))&&(null==s.to||(c?s.to>=t.ch:s.to>t.ch))){if(i&&(pe(a,"beforeCursorEnter"),a.explicitlyCleared)){if(o.markedSpans){--l;continue}break}if(!a.atomic)continue;if(r){var h=a.find(n<0?1:-1),f=void 0;if((n<0?c:u)&&(h=lo(e,h,-n,h&&h.line==t.line?o:null)),h&&h.line==t.line&&(f=tt(h,r))&&(n<0?f<0:f>0))return io(e,h,t,n,i)}var d=a.find(n<0?-1:1);return(n<0?u:c)&&(d=lo(e,d,n,d.line==t.line?o:null)),d?io(e,d,t,n,i):null}}return t}function oo(e,t,r,n,i){var o=n||1,l=io(e,t,r,o,i)||!i&&io(e,t,r,o,!0)||io(e,t,r,-o,i)||!i&&io(e,t,r,-o,!0);return l||(e.cantEdit=!0,et(e.first,0))}function lo(e,t,r,n){return r<0&&0==t.ch?t.line>e.first?st(e,et(t.line-1)):null:r>0&&t.ch==(n||Xe(e,t.line)).text.length?t.line<e.first+e.size-1?et(t.line+1,0):null:new et(t.line,t.ch+r)}function so(e){e.setSelection(et(e.firstLine(),0),et(e.lastLine()),U)}function ao(e,t,r){var n={canceled:!1,from:t.from,to:t.to,text:t.text,origin:t.origin,cancel:function(){return n.canceled=!0}};return r&&(n.update=function(t,r,i,o){t&&(n.from=st(e,t)),r&&(n.to=st(e,r)),i&&(n.text=i),void 0!==o&&(n.origin=o)}),pe(e,"beforeChange",e,n),e.cm&&pe(e.cm,"beforeChange",e.cm,n),n.canceled?(e.cm&&(e.cm.curOp.updateInput=2),null):{from:n.from,to:n.to,text:n.text,origin:n.origin}}function uo(e,t,r){if(e.cm){if(!e.cm.curOp)return ei(e.cm,uo)(e,t,r);if(e.cm.state.suppressEdits)return}if(!(me(e,"beforeChange")||e.cm&&me(e.cm,"beforeChange"))||(t=ao(e,t,!0))){var n=xt&&!r&&function(e,t,r){var n=null;if(e.iter(t.line,r.line+1,(function(e){if(e.markedSpans)for(var t=0;t<e.markedSpans.length;++t){var r=e.markedSpans[t].marker;!r.readOnly||n&&-1!=B(n,r)||(n||(n=[])).push(r)}})),!n)return null;for(var i=[{from:t,to:r}],o=0;o<n.length;++o)for(var l=n[o],s=l.find(0),a=0;a<i.length;++a){var u=i[a];if(!(tt(u.to,s.from)<0||tt(u.from,s.to)>0)){var c=[a,1],h=tt(u.from,s.from),f=tt(u.to,s.to);(h<0||!l.inclusiveLeft&&!h)&&c.push({from:u.from,to:s.from}),(f>0||!l.inclusiveRight&&!f)&&c.push({from:s.to,to:u.to}),i.splice.apply(i,c),a+=c.length-3}}return i}(e,t.from,t.to);if(n)for(var i=n.length-1;i>=0;--i)co(e,{from:n[i].from,to:n[i].to,text:i?[""]:t.text,origin:t.origin});else co(e,t)}}function co(e,t){if(1!=t.text.length||""!=t.text[0]||0!=tt(t.from,t.to)){var r=Ni(e,t);Bi(e,t,r,e.cm?e.cm.curOp.id:NaN),po(e,t,r,Tt(e,t));var n=[];Fi(e,(function(e,r){r||-1!=B(n,e.history)||(yo(e.history,t),n.push(e.history)),po(e,t,null,Tt(e,t))}))}}function ho(e,t,r){var n=e.cm&&e.cm.state.suppressEdits;if(!n||r){for(var i,o=e.history,l=e.sel,s="undo"==t?o.done:o.undone,a="undo"==t?o.undone:o.done,u=0;u<s.length&&(i=s[u],r?!i.ranges||i.equals(e.sel):i.ranges);u++);if(u!=s.length){for(o.lastOrigin=o.lastSelOrigin=null;;){if(!(i=s.pop()).ranges){if(n)return void s.push(i);break}if(Ui(i,a),r&&!i.equals(e.sel))return void Ji(e,i,{clearRedo:!1});l=i}var c=[];Ui(l,a),a.push({changes:c,generation:o.generation}),o.generation=i.generation||++o.maxGeneration;for(var h=me(e,"beforeChange")||e.cm&&me(e.cm,"beforeChange"),f=function(r){var n=i.changes[r];if(n.origin=t,h&&!ao(e,n,!1))return s.length=0,{};c.push(Ri(e,n));var o=r?Ni(e,n):$(s);po(e,n,o,ji(e,n)),!r&&e.cm&&e.cm.scrollIntoView({from:n.from,to:Ti(n)});var l=[];Fi(e,(function(e,t){t||-1!=B(l,e.history)||(yo(e.history,n),l.push(e.history)),po(e,n,null,ji(e,n))}))},d=i.changes.length-1;d>=0;--d){var p=f(d);if(p)return p.v}}}}function fo(e,t){if(0!=t&&(e.first+=t,e.sel=new Ci(_(e.sel.ranges,(function(e){return new Si(et(e.anchor.line+t,e.anchor.ch),et(e.head.line+t,e.head.ch))})),e.sel.primIndex),e.cm)){hn(e.cm,e.first,e.first-t,t);for(var r=e.cm.display,n=r.viewFrom;n<r.viewTo;n++)fn(e.cm,n,"gutter")}}function po(e,t,r,n){if(e.cm&&!e.cm.curOp)return ei(e.cm,po)(e,t,r,n);if(t.to.line<e.first)fo(e,t.text.length-1-(t.to.line-t.from.line));else if(!(t.from.line>e.lastLine())){if(t.from.line<e.first){var i=t.text.length-1-(e.first-t.from.line);fo(e,i),t={from:et(e.first,0),to:et(t.to.line+i,t.to.ch),text:[$(t.text)],origin:t.origin}}var o=e.lastLine();t.to.line>o&&(t={from:t.from,to:et(o,Xe(e,o).text.length),text:[t.text[0]],origin:t.origin}),t.removed=Ye(e,t.from,t.to),r||(r=Ni(e,t)),e.cm?function(e,t,r){var n=e.doc,i=e.display,o=t.from,l=t.to,s=!1,a=o.line;e.options.lineWrapping||(a=qe(Rt(Xe(n,o.line))),n.iter(a,l.line+1,(function(e){if(e==i.maxLine)return s=!0,!0})));n.sel.contains(t.from,t.to)>-1&&ve(e);Hi(n,t,r,sn(e)),e.options.lineWrapping||(n.iter(a,o.line+t.text.length,(function(e){var t=Kt(e);t>i.maxLineLength&&(i.maxLine=e,i.maxLineLength=t,i.maxLineChanged=!0,s=!1)})),s&&(e.curOp.updateMaxLine=!0));(function(e,t){if(e.modeFrontier=Math.min(e.modeFrontier,t),!(e.highlightFrontier<t-10)){for(var r=e.first,n=t-1;n>r;n--){var i=Xe(e,n).stateAfter;if(i&&(!(i instanceof ut)||n+i.lookAhead<t)){r=n+1;break}}e.highlightFrontier=Math.min(e.highlightFrontier,r)}})(n,o.line),ni(e,400);var u=t.text.length-(l.line-o.line)-1;t.full?hn(e):o.line!=l.line||1!=t.text.length||Wi(e.doc,t)?hn(e,o.line,l.line+1,u):fn(e,o.line,"text");var c=me(e,"changes"),h=me(e,"change");if(h||c){var f={from:o,to:l,text:t.text,removed:t.removed,origin:t.origin};h&&sr(e,"change",e,f),c&&(e.curOp.changeObjs||(e.curOp.changeObjs=[])).push(f)}e.display.selForContextMenu=null}(e.cm,t,n):Hi(e,t,n),eo(e,r,U),e.cantEdit&&oo(e,et(e.firstLine(),0))&&(e.cantEdit=!1)}}function go(e,t,r,n,i){var o;n||(n=r),tt(n,r)<0&&(r=(o=[n,r])[0],n=o[1]),"string"==typeof t&&(t=e.splitLines(t)),uo(e,{from:r,to:n,text:t,origin:i})}function vo(e,t,r,n){r<e.line?e.line+=n:t<e.line&&(e.line=t,e.ch=0)}function mo(e,t,r,n){for(var i=0;i<e.length;++i){var o=e[i],l=!0;if(o.ranges){o.copied||((o=e[i]=o.deepCopy()).copied=!0);for(var s=0;s<o.ranges.length;s++)vo(o.ranges[s].anchor,t,r,n),vo(o.ranges[s].head,t,r,n)}else{for(var a=0;a<o.changes.length;++a){var u=o.changes[a];if(r<u.from.line)u.from=et(u.from.line+n,u.from.ch),u.to=et(u.to.line+n,u.to.ch);else if(t<=u.to.line){l=!1;break}}l||(e.splice(0,i+1),i=0)}}}function yo(e,t){var r=t.from.line,n=t.to.line,i=t.text.length-(n-r)-1;mo(e.done,r,n,i),mo(e.undone,r,n,i)}function bo(e,t,r,n){var i=t,o=t;return"number"==typeof t?o=Xe(e,lt(e,t)):i=qe(t),null==i?null:(n(o,i)&&e.cm&&fn(e.cm,i,r),o)}function wo(e){this.lines=e,this.parent=null;for(var t=0,r=0;r<e.length;++r)e[r].parent=this,t+=e[r].height;this.height=t}function xo(e){this.children=e;for(var t=0,r=0,n=0;n<e.length;++n){var i=e[n];t+=i.chunkSize(),r+=i.height,i.parent=this}this.size=t,this.height=r,this.parent=null}Si.prototype.from=function(){return ot(this.anchor,this.head)},Si.prototype.to=function(){return it(this.anchor,this.head)},Si.prototype.empty=function(){return this.head.line==this.anchor.line&&this.head.ch==this.anchor.ch},wo.prototype={chunkSize:function(){return this.lines.length},removeInner:function(e,t){for(var r=e,n=e+t;r<n;++r){var i=this.lines[r];this.height-=i.height,Yt(i),sr(i,"delete")}this.lines.splice(e,t)},collapse:function(e){e.push.apply(e,this.lines)},insertInner:function(e,t,r){this.height+=r,this.lines=this.lines.slice(0,e).concat(t).concat(this.lines.slice(e));for(var n=0;n<t.length;++n)t[n].parent=this},iterN:function(e,t,r){for(var n=e+t;e<n;++e)if(r(this.lines[e]))return!0}},xo.prototype={chunkSize:function(){return this.size},removeInner:function(e,t){this.size-=t;for(var r=0;r<this.children.length;++r){var n=this.children[r],i=n.chunkSize();if(e<i){var o=Math.min(t,i-e),l=n.height;if(n.removeInner(e,o),this.height-=l-n.height,i==o&&(this.children.splice(r--,1),n.parent=null),0==(t-=o))break;e=0}else e-=i}if(this.size-t<25&&(this.children.length>1||!(this.children[0]instanceof wo))){var s=[];this.collapse(s),this.children=[new wo(s)],this.children[0].parent=this}},collapse:function(e){for(var t=0;t<this.children.length;++t)this.children[t].collapse(e)},insertInner:function(e,t,r){this.size+=t.length,this.height+=r;for(var n=0;n<this.children.length;++n){var i=this.children[n],o=i.chunkSize();if(e<=o){if(i.insertInner(e,t,r),i.lines&&i.lines.length>50){for(var l=i.lines.length%25+25,s=l;s<i.lines.length;){var a=new wo(i.lines.slice(s,s+=25));i.height-=a.height,this.children.splice(++n,0,a),a.parent=this}i.lines=i.lines.slice(0,l),this.maybeSpill()}break}e-=o}},maybeSpill:function(){if(!(this.children.length<=10)){var e=this;do{var t=new xo(e.children.splice(e.children.length-5,5));if(e.parent){e.size-=t.size,e.height-=t.height;var r=B(e.parent.children,e);e.parent.children.splice(r+1,0,t)}else{var n=new xo(e.children);n.parent=e,e.children=[n,t],e=n}t.parent=e.parent}while(e.children.length>10);e.parent.maybeSpill()}},iterN:function(e,t,r){for(var n=0;n<this.children.length;++n){var i=this.children[n],o=i.chunkSize();if(e<o){var l=Math.min(t,o-e);if(i.iterN(e,l,r))return!0;if(0==(t-=l))break;e=0}else e-=o}}};var Co=function(e,t,r){if(r)for(var n in r)r.hasOwnProperty(n)&&(this[n]=r[n]);this.doc=e,this.node=t};function So(e,t,r){Vt(t)<(e.curOp&&e.curOp.scrollTop||e.doc.scrollTop)&&An(e,r)}Co.prototype.clear=function(){var e=this.doc.cm,t=this.line.widgets,r=this.line,n=qe(r);if(null!=n&&t){for(var i=0;i<t.length;++i)t[i]==this&&t.splice(i--,1);t.length||(r.widgets=null);var o=wr(this);_e(r,Math.max(0,r.height-o)),e&&(Jn(e,(function(){So(e,r,-o),fn(e,n,"widget")})),sr(e,"lineWidgetCleared",e,this,n))}},Co.prototype.changed=function(){var e=this,t=this.height,r=this.doc.cm,n=this.line;this.height=null;var i=wr(this)-t;i&&(Gt(this.doc,n)||_e(n,n.height+i),r&&Jn(r,(function(){r.curOp.forceUpdate=!0,So(r,n,i),sr(r,"lineWidgetChanged",r,e,qe(n))})))},ye(Co);var Lo=0,ko=function(e,t){this.lines=[],this.type=t,this.doc=e,this.id=++Lo};function To(e,t,r,n,i){if(n&&n.shared)return function(e,t,r,n,i){(n=I(n)).shared=!1;var o=[To(e,t,r,n,i)],l=o[0],s=n.widgetNode;return Fi(e,(function(e){s&&(n.widgetNode=s.cloneNode(!0)),o.push(To(e,st(e,t),st(e,r),n,i));for(var a=0;a<e.linked.length;++a)if(e.linked[a].isParent)return;l=$(o)})),new Mo(o,l)}(e,t,r,n,i);if(e.cm&&!e.cm.curOp)return ei(e.cm,To)(e,t,r,n,i);var o=new ko(e,i),l=tt(t,r);if(n&&I(n,o,!1),l>0||0==l&&!1!==o.clearWhenEmpty)return o;if(o.replacedWith&&(o.collapsed=!0,o.widgetNode=A("span",[o.replacedWith],"CodeMirror-widget"),n.handleMouseEvents||o.widgetNode.setAttribute("cm-ignore-events","true"),n.insertLeft&&(o.widgetNode.insertLeft=!0)),o.collapsed){if(It(e,t.line,t,r,o)||t.line!=r.line&&It(e,r.line,t,r,o))throw new Error("Inserting collapsed marker partially overlapping an existing one");Ct=!0}o.addToHistory&&Bi(e,{from:t,to:r,origin:"markText"},e.sel,NaN);var s,a=t.line,u=e.cm;if(e.iter(a,r.line+1,(function(n){u&&o.collapsed&&!u.options.lineWrapping&&Rt(n)==u.display.maxLine&&(s=!0),o.collapsed&&a!=t.line&&_e(n,0),function(e,t,r){var n=r&&window.WeakSet&&(r.markedSpans||(r.markedSpans=new WeakSet));n&&e.markedSpans&&n.has(e.markedSpans)?e.markedSpans.push(t):(e.markedSpans=e.markedSpans?e.markedSpans.concat([t]):[t],n&&n.add(e.markedSpans)),t.marker.attachLine(e)}(n,new St(o,a==t.line?t.ch:null,a==r.line?r.ch:null),e.cm&&e.cm.curOp),++a})),o.collapsed&&e.iter(t.line,r.line+1,(function(t){Gt(e,t)&&_e(t,0)})),o.clearOnEnter&&he(o,"beforeCursorEnter",(function(){return o.clear()})),o.readOnly&&(xt=!0,(e.history.done.length||e.history.undone.length)&&e.clearHistory()),o.collapsed&&(o.id=++Lo,o.atomic=!0),u){if(s&&(u.curOp.updateMaxLine=!0),o.collapsed)hn(u,t.line,r.line+1);else if(o.className||o.startStyle||o.endStyle||o.css||o.attributes||o.title)for(var c=t.line;c<=r.line;c++)fn(u,c,"text");o.atomic&&ro(u.doc),sr(u,"markerAdded",u,o)}return o}ko.prototype.clear=function(){if(!this.explicitlyCleared){var e=this.doc.cm,t=e&&!e.curOp;if(t&&Xn(e),me(this,"clear")){var r=this.find();r&&sr(this,"clear",r.from,r.to)}for(var n=null,i=null,o=0;o<this.lines.length;++o){var l=this.lines[o],s=Lt(l.markedSpans,this);e&&!this.collapsed?fn(e,qe(l),"text"):e&&(null!=s.to&&(i=qe(l)),null!=s.from&&(n=qe(l))),l.markedSpans=kt(l.markedSpans,s),null==s.from&&this.collapsed&&!Gt(this.doc,l)&&e&&_e(l,rn(e.display))}if(e&&this.collapsed&&!e.options.lineWrapping)for(var a=0;a<this.lines.length;++a){var u=Rt(this.lines[a]),c=Kt(u);c>e.display.maxLineLength&&(e.display.maxLine=u,e.display.maxLineLength=c,e.display.maxLineChanged=!0)}null!=n&&e&&this.collapsed&&hn(e,n,i+1),this.lines.length=0,this.explicitlyCleared=!0,this.atomic&&this.doc.cantEdit&&(this.doc.cantEdit=!1,e&&ro(e.doc)),e&&sr(e,"markerCleared",e,this,n,i),t&&Yn(e),this.parent&&this.parent.clear()}},ko.prototype.find=function(e,t){var r,n;null==e&&"bookmark"==this.type&&(e=1);for(var i=0;i<this.lines.length;++i){var o=this.lines[i],l=Lt(o.markedSpans,this);if(null!=l.from&&(r=et(t?o:qe(o),l.from),-1==e))return r;if(null!=l.to&&(n=et(t?o:qe(o),l.to),1==e))return n}return r&&{from:r,to:n}},ko.prototype.changed=function(){var e=this,t=this.find(-1,!0),r=this,n=this.doc.cm;t&&n&&Jn(n,(function(){var i=t.line,o=qe(t.line),l=Ar(n,o);if(l&&(Ir(l),n.curOp.selectionChanged=n.curOp.forceUpdate=!0),n.curOp.updateMaxLine=!0,!Gt(r.doc,i)&&null!=r.height){var s=r.height;r.height=null;var a=wr(r)-s;a&&_e(i,i.height+a)}sr(n,"markerChanged",n,e)}))},ko.prototype.attachLine=function(e){if(!this.lines.length&&this.doc.cm){var t=this.doc.cm.curOp;t.maybeHiddenMarkers&&-1!=B(t.maybeHiddenMarkers,this)||(t.maybeUnhiddenMarkers||(t.maybeUnhiddenMarkers=[])).push(this)}this.lines.push(e)},ko.prototype.detachLine=function(e){if(this.lines.splice(B(this.lines,e),1),!this.lines.length&&this.doc.cm){var t=this.doc.cm.curOp;(t.maybeHiddenMarkers||(t.maybeHiddenMarkers=[])).push(this)}},ye(ko);var Mo=function(e,t){this.markers=e,this.primary=t;for(var r=0;r<e.length;++r)e[r].parent=this};function No(e){return e.findMarks(et(e.first,0),e.clipPos(et(e.lastLine())),(function(e){return e.parent}))}function Oo(e){for(var t=function(t){var r=e[t],n=[r.primary.doc];Fi(r.primary.doc,(function(e){return n.push(e)}));for(var i=0;i<r.markers.length;i++){var o=r.markers[i];-1==B(n,o.doc)&&(o.parent=null,r.markers.splice(i--,1))}},r=0;r<e.length;r++)t(r)}Mo.prototype.clear=function(){if(!this.explicitlyCleared){this.explicitlyCleared=!0;for(var e=0;e<this.markers.length;++e)this.markers[e].clear();sr(this,"clear")}},Mo.prototype.find=function(e,t){return this.primary.find(e,t)},ye(Mo);var Ao=0,Do=function(e,t,r,n,i){if(!(this instanceof Do))return new Do(e,t,r,n,i);null==r&&(r=0),xo.call(this,[new wo([new Xt("",null)])]),this.first=r,this.scrollTop=this.scrollLeft=0,this.cantEdit=!1,this.cleanGeneration=1,this.modeFrontier=this.highlightFrontier=r;var o=et(r,0);this.sel=ki(o),this.history=new Ii(null),this.id=++Ao,this.modeOption=t,this.lineSep=n,this.direction="rtl"==i?"rtl":"ltr",this.extend=!1,"string"==typeof e&&(e=this.splitLines(e)),Hi(this,{from:o,to:o,text:e}),Ji(this,ki(o),U)};Do.prototype=Z(xo.prototype,{constructor:Do,iter:function(e,t,r){r?this.iterN(e-this.first,t-e,r):this.iterN(this.first,this.first+this.size,e)},insert:function(e,t){for(var r=0,n=0;n<t.length;++n)r+=t[n].height;this.insertInner(e-this.first,t,r)},remove:function(e,t){this.removeInner(e-this.first,t)},getValue:function(e){var t=$e(this,this.first,this.first+this.size);return!1===e?t:t.join(e||this.lineSeparator())},setValue:ri((function(e){var t=et(this.first,0),r=this.first+this.size-1;uo(this,{from:t,to:et(r,Xe(this,r).text.length),text:this.splitLines(e),origin:"setValue",full:!0},!0),this.cm&&Wn(this.cm,0,0),Ji(this,ki(t),U)})),replaceRange:function(e,t,r,n){go(this,e,t=st(this,t),r=r?st(this,r):t,n)},getRange:function(e,t,r){var n=Ye(this,st(this,e),st(this,t));return!1===r?n:""===r?n.join(""):n.join(r||this.lineSeparator())},getLine:function(e){var t=this.getLineHandle(e);return t&&t.text},getLineHandle:function(e){if(Qe(this,e))return Xe(this,e)},getLineNumber:function(e){return qe(e)},getLineHandleVisualStart:function(e){return"number"==typeof e&&(e=Xe(this,e)),Rt(e)},lineCount:function(){return this.size},firstLine:function(){return this.first},lastLine:function(){return this.first+this.size-1},clipPos:function(e){return st(this,e)},getCursor:function(e){var t=this.sel.primary();return null==e||"head"==e?t.head:"anchor"==e?t.anchor:"end"==e||"to"==e||!1===e?t.to():t.from()},listSelections:function(){return this.sel.ranges},somethingSelected:function(){return this.sel.somethingSelected()},setCursor:ri((function(e,t,r){Zi(this,st(this,"number"==typeof e?et(e,t||0):e),null,r)})),setSelection:ri((function(e,t,r){Zi(this,st(this,e),st(this,t||e),r)})),extendSelection:ri((function(e,t,r){$i(this,st(this,e),t&&st(this,t),r)})),extendSelections:ri((function(e,t){_i(this,at(this,e),t)})),extendSelectionsBy:ri((function(e,t){_i(this,at(this,_(this.sel.ranges,e)),t)})),setSelections:ri((function(e,t,r){if(e.length){for(var n=[],i=0;i<e.length;i++)n[i]=new Si(st(this,e[i].anchor),st(this,e[i].head||e[i].anchor));null==t&&(t=Math.min(e.length-1,this.sel.primIndex)),Ji(this,Li(this.cm,n,t),r)}})),addSelection:ri((function(e,t,r){var n=this.sel.ranges.slice(0);n.push(new Si(st(this,e),st(this,t||e))),Ji(this,Li(this.cm,n,n.length-1),r)})),getSelection:function(e){for(var t,r=this.sel.ranges,n=0;n<r.length;n++){var i=Ye(this,r[n].from(),r[n].to());t=t?t.concat(i):i}return!1===e?t:t.join(e||this.lineSeparator())},getSelections:function(e){for(var t=[],r=this.sel.ranges,n=0;n<r.length;n++){var i=Ye(this,r[n].from(),r[n].to());!1!==e&&(i=i.join(e||this.lineSeparator())),t[n]=i}return t},replaceSelection:function(e,t,r){for(var n=[],i=0;i<this.sel.ranges.length;i++)n[i]=e;this.replaceSelections(n,t,r||"+input")},replaceSelections:ri((function(e,t,r){for(var n=[],i=this.sel,o=0;o<i.ranges.length;o++){var l=i.ranges[o];n[o]={from:l.from(),to:l.to(),text:this.splitLines(e[o]),origin:r}}for(var s=t&&"end"!=t&&function(e,t,r){for(var n=[],i=et(e.first,0),o=i,l=0;l<t.length;l++){var s=t[l],a=Oi(s.from,i,o),u=Oi(Ti(s),i,o);if(i=s.to,o=u,"around"==r){var c=e.sel.ranges[l],h=tt(c.head,c.anchor)<0;n[l]=new Si(h?u:a,h?a:u)}else n[l]=new Si(a,a)}return new Ci(n,e.sel.primIndex)}(this,n,t),a=n.length-1;a>=0;a--)uo(this,n[a]);s?Qi(this,s):this.cm&&Dn(this.cm)})),undo:ri((function(){ho(this,"undo")})),redo:ri((function(){ho(this,"redo")})),undoSelection:ri((function(){ho(this,"undo",!0)})),redoSelection:ri((function(){ho(this,"redo",!0)})),setExtending:function(e){this.extend=e},getExtending:function(){return this.extend},historySize:function(){for(var e=this.history,t=0,r=0,n=0;n<e.done.length;n++)e.done[n].ranges||++t;for(var i=0;i<e.undone.length;i++)e.undone[i].ranges||++r;return{undo:t,redo:r}},clearHistory:function(){var e=this;this.history=new Ii(this.history),Fi(this,(function(t){return t.history=e.history}),!0)},markClean:function(){this.cleanGeneration=this.changeGeneration(!0)},changeGeneration:function(e){return e&&(this.history.lastOp=this.history.lastSelOp=this.history.lastOrigin=null),this.history.generation},isClean:function(e){return this.history.generation==(e||this.cleanGeneration)},getHistory:function(){return{done:Xi(this.history.done),undone:Xi(this.history.undone)}},setHistory:function(e){var t=this.history=new Ii(this.history);t.done=Xi(e.done.slice(0),null,!0),t.undone=Xi(e.undone.slice(0),null,!0)},setGutterMarker:ri((function(e,t,r){return bo(this,e,"gutter",(function(e){var n=e.gutterMarkers||(e.gutterMarkers={});return n[t]=r,!r&&te(n)&&(e.gutterMarkers=null),!0}))})),clearGutter:ri((function(e){var t=this;this.iter((function(r){r.gutterMarkers&&r.gutterMarkers[e]&&bo(t,r,"gutter",(function(){return r.gutterMarkers[e]=null,te(r.gutterMarkers)&&(r.gutterMarkers=null),!0}))}))})),lineInfo:function(e){var t;if("number"==typeof e){if(!Qe(this,e))return null;if(t=e,!(e=Xe(this,e)))return null}else if(null==(t=qe(e)))return null;return{line:t,handle:e,text:e.text,gutterMarkers:e.gutterMarkers,textClass:e.textClass,bgClass:e.bgClass,wrapClass:e.wrapClass,widgets:e.widgets}},addLineClass:ri((function(e,t,r){return bo(this,e,"gutter"==t?"gutter":"class",(function(e){var n="text"==t?"textClass":"background"==t?"bgClass":"gutter"==t?"gutterClass":"wrapClass";if(e[n]){if(L(r).test(e[n]))return!1;e[n]+=" "+r}else e[n]=r;return!0}))})),removeLineClass:ri((function(e,t,r){return bo(this,e,"gutter"==t?"gutter":"class",(function(e){var n="text"==t?"textClass":"background"==t?"bgClass":"gutter"==t?"gutterClass":"wrapClass",i=e[n];if(!i)return!1;if(null==r)e[n]=null;else{var o=i.match(L(r));if(!o)return!1;var l=o.index+o[0].length;e[n]=i.slice(0,o.index)+(o.index&&l!=i.length?" ":"")+i.slice(l)||null}return!0}))})),addLineWidget:ri((function(e,t,r){return function(e,t,r,n){var i=new Co(e,r,n),o=e.cm;return o&&i.noHScroll&&(o.display.alignWidgets=!0),bo(e,t,"widget",(function(t){var r=t.widgets||(t.widgets=[]);if(null==i.insertAt?r.push(i):r.splice(Math.min(r.length,Math.max(0,i.insertAt)),0,i),i.line=t,o&&!Gt(e,t)){var n=Vt(t)<e.scrollTop;_e(t,t.height+wr(i)),n&&An(o,i.height),o.curOp.forceUpdate=!0}return!0})),o&&sr(o,"lineWidgetAdded",o,i,"number"==typeof t?t:qe(t)),i}(this,e,t,r)})),removeLineWidget:function(e){e.clear()},markText:function(e,t,r){return To(this,st(this,e),st(this,t),r,r&&r.type||"range")},setBookmark:function(e,t){var r={replacedWith:t&&(null==t.nodeType?t.widget:t),insertLeft:t&&t.insertLeft,clearWhenEmpty:!1,shared:t&&t.shared,handleMouseEvents:t&&t.handleMouseEvents};return To(this,e=st(this,e),e,r,"bookmark")},findMarksAt:function(e){var t=[],r=Xe(this,(e=st(this,e)).line).markedSpans;if(r)for(var n=0;n<r.length;++n){var i=r[n];(null==i.from||i.from<=e.ch)&&(null==i.to||i.to>=e.ch)&&t.push(i.marker.parent||i.marker)}return t},findMarks:function(e,t,r){e=st(this,e),t=st(this,t);var n=[],i=e.line;return this.iter(e.line,t.line+1,(function(o){var l=o.markedSpans;if(l)for(var s=0;s<l.length;s++){var a=l[s];null!=a.to&&i==e.line&&e.ch>=a.to||null==a.from&&i!=e.line||null!=a.from&&i==t.line&&a.from>=t.ch||r&&!r(a.marker)||n.push(a.marker.parent||a.marker)}++i})),n},getAllMarks:function(){var e=[];return this.iter((function(t){var r=t.markedSpans;if(r)for(var n=0;n<r.length;++n)null!=r[n].from&&e.push(r[n].marker)})),e},posFromIndex:function(e){var t,r=this.first,n=this.lineSeparator().length;return this.iter((function(i){var o=i.text.length+n;if(o>e)return t=e,!0;e-=o,++r})),st(this,et(r,t))},indexFromPos:function(e){var t=(e=st(this,e)).ch;if(e.line<this.first||e.ch<0)return 0;var r=this.lineSeparator().length;return this.iter(this.first,e.line,(function(e){t+=e.text.length+r})),t},copy:function(e){var t=new Do($e(this,this.first,this.first+this.size),this.modeOption,this.first,this.lineSep,this.direction);return t.scrollTop=this.scrollTop,t.scrollLeft=this.scrollLeft,t.sel=this.sel,t.extend=!1,e&&(t.history.undoDepth=this.history.undoDepth,t.setHistory(this.getHistory())),t},linkedDoc:function(e){e||(e={});var t=this.first,r=this.first+this.size;null!=e.from&&e.from>t&&(t=e.from),null!=e.to&&e.to<r&&(r=e.to);var n=new Do($e(this,t,r),e.mode||this.modeOption,t,this.lineSep,this.direction);return e.sharedHist&&(n.history=this.history),(this.linked||(this.linked=[])).push({doc:n,sharedHist:e.sharedHist}),n.linked=[{doc:this,isParent:!0,sharedHist:e.sharedHist}],function(e,t){for(var r=0;r<t.length;r++){var n=t[r],i=n.find(),o=e.clipPos(i.from),l=e.clipPos(i.to);if(tt(o,l)){var s=To(e,o,l,n.primary,n.primary.type);n.markers.push(s),s.parent=n}}}(n,No(this)),n},unlinkDoc:function(e){if(e instanceof Ml&&(e=e.doc),this.linked)for(var t=0;t<this.linked.length;++t){if(this.linked[t].doc==e){this.linked.splice(t,1),e.unlinkDoc(this),Oo(No(this));break}}if(e.history==this.history){var r=[e.id];Fi(e,(function(e){return r.push(e.id)}),!0),e.history=new Ii(null),e.history.done=Xi(this.history.done,r),e.history.undone=Xi(this.history.undone,r)}},iterLinkedDocs:function(e){Fi(this,e)},getMode:function(){return this.mode},getEditor:function(){return this.cm},splitLines:function(e){return this.lineSep?e.split(this.lineSep):De(e)},lineSeparator:function(){return this.lineSep||"\n"},setDirection:ri((function(e){var t;("rtl"!=e&&(e="ltr"),e!=this.direction)&&(this.direction=e,this.iter((function(e){return e.order=null})),this.cm&&Jn(t=this.cm,(function(){Ei(t),hn(t)})))}))}),Do.prototype.eachLine=Do.prototype.iter;var Wo=0;function Ho(e){var t=this;if(Fo(t),!ge(t,e)&&!xr(t.display,e)){be(e),l&&(Wo=+new Date);var r=un(t,e,!0),n=e.dataTransfer.files;if(r&&!t.isReadOnly())if(n&&n.length&&window.FileReader&&window.File)for(var i=n.length,o=Array(i),s=0,a=function(){++s==i&&ei(t,(function(){var e={from:r=st(t.doc,r),to:r,text:t.doc.splitLines(o.filter((function(e){return null!=e})).join(t.doc.lineSeparator())),origin:"paste"};uo(t.doc,e),Qi(t.doc,ki(st(t.doc,r),st(t.doc,Ti(e))))}))()},u=function(e,r){if(t.options.allowDropFileTypes&&-1==B(t.options.allowDropFileTypes,e.type))a();else{var n=new FileReader;n.onerror=function(){return a()},n.onload=function(){var e=n.result;/[\x00-\x08\x0e-\x1f]{2}/.test(e)||(o[r]=e),a()},n.readAsText(e)}},c=0;c<n.length;c++)u(n[c],c);else{if(t.state.draggingText&&t.doc.sel.contains(r)>-1)return t.state.draggingText(e),void setTimeout((function(){return t.display.input.focus()}),20);try{var h=e.dataTransfer.getData("Text");if(h){var f;if(t.state.draggingText&&!t.state.draggingText.copy&&(f=t.listSelections()),eo(t.doc,ki(r,r)),f)for(var d=0;d<f.length;++d)go(t.doc,"",f[d].anchor,f[d].head,"drag");t.replaceSelection(h,"around","paste"),t.display.input.focus()}}catch(e){}}}}function Fo(e){e.display.dragCursor&&(e.display.lineSpace.removeChild(e.display.dragCursor),e.display.dragCursor=null)}function Po(e){if(document.getElementsByClassName){for(var t=document.getElementsByClassName("CodeMirror"),r=[],n=0;n<t.length;n++){var i=t[n].CodeMirror;i&&r.push(i)}r.length&&r[0].operation((function(){for(var t=0;t<r.length;t++)e(r[t])}))}}var Eo=!1;function Io(){var e;Eo||(he(window,"resize",(function(){null==e&&(e=setTimeout((function(){e=null,Po(Ro)}),100))})),he(window,"blur",(function(){return Po(kn)})),Eo=!0)}function Ro(e){var t=e.display;t.cachedCharWidth=t.cachedTextHeight=t.cachedPaddingH=null,t.scrollbarsClipped=!1,e.setSize()}for(var zo={3:"Pause",8:"Backspace",9:"Tab",13:"Enter",16:"Shift",17:"Ctrl",18:"Alt",19:"Pause",20:"CapsLock",27:"Esc",32:"Space",33:"PageUp",34:"PageDown",35:"End",36:"Home",37:"Left",38:"Up",39:"Right",40:"Down",44:"PrintScrn",45:"Insert",46:"Delete",59:";",61:"=",91:"Mod",92:"Mod",93:"Mod",106:"*",107:"=",109:"-",110:".",111:"/",145:"ScrollLock",173:"-",186:";",187:"=",188:",",189:"-",190:".",191:"/",192:"`",219:"[",220:"\\",221:"]",222:"'",224:"Mod",63232:"Up",63233:"Down",63234:"Left",63235:"Right",63272:"Delete",63273:"Home",63275:"End",63276:"PageUp",63277:"PageDown",63302:"Insert"},Bo=0;Bo<10;Bo++)zo[Bo+48]=zo[Bo+96]=String(Bo);for(var Go=65;Go<=90;Go++)zo[Go]=String.fromCharCode(Go);for(var Uo=1;Uo<=12;Uo++)zo[Uo+111]=zo[Uo+63235]="F"+Uo;var Vo={};function Ko(e){var t,r,n,i,o=e.split(/-(?!$)/);e=o[o.length-1];for(var l=0;l<o.length-1;l++){var s=o[l];if(/^(cmd|meta|m)$/i.test(s))i=!0;else if(/^a(lt)?$/i.test(s))t=!0;else if(/^(c|ctrl|control)$/i.test(s))r=!0;else{if(!/^s(hift)?$/i.test(s))throw new Error("Unrecognized modifier name: "+s);n=!0}}return t&&(e="Alt-"+e),r&&(e="Ctrl-"+e),i&&(e="Cmd-"+e),n&&(e="Shift-"+e),e}function jo(e){var t={};for(var r in e)if(e.hasOwnProperty(r)){var n=e[r];if(/^(name|fallthrough|(de|at)tach)$/.test(r))continue;if("..."==n){delete e[r];continue}for(var i=_(r.split(" "),Ko),o=0;o<i.length;o++){var l=void 0,s=void 0;o==i.length-1?(s=i.join(" "),l=n):(s=i.slice(0,o+1).join(" "),l="...");var a=t[s];if(a){if(a!=l)throw new Error("Inconsistent bindings for "+s)}else t[s]=l}delete e[r]}for(var u in t)e[u]=t[u];return e}function Xo(e,t,r,n){var i=(t=qo(t)).call?t.call(e,n):t[e];if(!1===i)return"nothing";if("..."===i)return"multi";if(null!=i&&r(i))return"handled";if(t.fallthrough){if("[object Array]"!=Object.prototype.toString.call(t.fallthrough))return Xo(e,t.fallthrough,r,n);for(var o=0;o<t.fallthrough.length;o++){var l=Xo(e,t.fallthrough[o],r,n);if(l)return l}}}function Yo(e){var t="string"==typeof e?e:zo[e.keyCode];return"Ctrl"==t||"Alt"==t||"Shift"==t||"Mod"==t}function $o(e,t,r){var n=e;return t.altKey&&"Alt"!=n&&(e="Alt-"+e),(C?t.metaKey:t.ctrlKey)&&"Ctrl"!=n&&(e="Ctrl-"+e),(C?t.ctrlKey:t.metaKey)&&"Mod"!=n&&(e="Cmd-"+e),!r&&t.shiftKey&&"Shift"!=n&&(e="Shift-"+e),e}function _o(e,t){if(h&&34==e.keyCode&&e.char)return!1;var r=zo[e.keyCode];return null!=r&&!e.altGraphKey&&(3==e.keyCode&&e.code&&(r=e.code),$o(r,e,t))}function qo(e){return"string"==typeof e?Vo[e]:e}function Zo(e,t){for(var r=e.doc.sel.ranges,n=[],i=0;i<r.length;i++){for(var o=t(r[i]);n.length&&tt(o.from,$(n).to)<=0;){var l=n.pop();if(tt(l.from,o.from)<0){o.from=l.from;break}}n.push(o)}Jn(e,(function(){for(var t=n.length-1;t>=0;t--)go(e.doc,"",n[t].from,n[t].to,"+delete");Dn(e)}))}function Qo(e,t,r){var n=ie(e.text,t+r,r);return n<0||n>e.text.length?null:n}function Jo(e,t,r){var n=Qo(e,t.ch,r);return null==n?null:new et(t.line,n,r<0?"after":"before")}function el(e,t,r,n,i){if(e){"rtl"==t.doc.direction&&(i=-i);var o=ue(r,t.doc.direction);if(o){var l,s=i<0?$(o):o[0],a=i<0==(1==s.level)?"after":"before";if(s.level>0||"rtl"==t.doc.direction){var u=Dr(t,r);l=i<0?r.text.length-1:0;var c=Wr(t,u,l).top;l=oe((function(e){return Wr(t,u,e).top==c}),i<0==(1==s.level)?s.from:s.to-1,l),"before"==a&&(l=Qo(r,l,1))}else l=i<0?s.to:s.from;return new et(n,l,a)}}return new et(n,i<0?r.text.length:0,i<0?"before":"after")}Vo.basic={Left:"goCharLeft",Right:"goCharRight",Up:"goLineUp",Down:"goLineDown",End:"goLineEnd",Home:"goLineStartSmart",PageUp:"goPageUp",PageDown:"goPageDown",Delete:"delCharAfter",Backspace:"delCharBefore","Shift-Backspace":"delCharBefore",Tab:"defaultTab","Shift-Tab":"indentAuto",Enter:"newlineAndIndent",Insert:"toggleOverwrite",Esc:"singleSelection"},Vo.pcDefault={"Ctrl-A":"selectAll","Ctrl-D":"deleteLine","Ctrl-Z":"undo","Shift-Ctrl-Z":"redo","Ctrl-Y":"redo","Ctrl-Home":"goDocStart","Ctrl-End":"goDocEnd","Ctrl-Up":"goLineUp","Ctrl-Down":"goLineDown","Ctrl-Left":"goGroupLeft","Ctrl-Right":"goGroupRight","Alt-Left":"goLineStart","Alt-Right":"goLineEnd","Ctrl-Backspace":"delGroupBefore","Ctrl-Delete":"delGroupAfter","Ctrl-S":"save","Ctrl-F":"find","Ctrl-G":"findNext","Shift-Ctrl-G":"findPrev","Shift-Ctrl-F":"replace","Shift-Ctrl-R":"replaceAll","Ctrl-[":"indentLess","Ctrl-]":"indentMore","Ctrl-U":"undoSelection","Shift-Ctrl-U":"redoSelection","Alt-U":"redoSelection",fallthrough:"basic"},Vo.emacsy={"Ctrl-F":"goCharRight","Ctrl-B":"goCharLeft","Ctrl-P":"goLineUp","Ctrl-N":"goLineDown","Ctrl-A":"goLineStart","Ctrl-E":"goLineEnd","Ctrl-V":"goPageDown","Shift-Ctrl-V":"goPageUp","Ctrl-D":"delCharAfter","Ctrl-H":"delCharBefore","Alt-Backspace":"delWordBefore","Ctrl-K":"killLine","Ctrl-T":"transposeChars","Ctrl-O":"openLine"},Vo.macDefault={"Cmd-A":"selectAll","Cmd-D":"deleteLine","Cmd-Z":"undo","Shift-Cmd-Z":"redo","Cmd-Y":"redo","Cmd-Home":"goDocStart","Cmd-Up":"goDocStart","Cmd-End":"goDocEnd","Cmd-Down":"goDocEnd","Alt-Left":"goGroupLeft","Alt-Right":"goGroupRight","Cmd-Left":"goLineLeft","Cmd-Right":"goLineRight","Alt-Backspace":"delGroupBefore","Ctrl-Alt-Backspace":"delGroupAfter","Alt-Delete":"delGroupAfter","Cmd-S":"save","Cmd-F":"find","Cmd-G":"findNext","Shift-Cmd-G":"findPrev","Cmd-Alt-F":"replace","Shift-Cmd-Alt-F":"replaceAll","Cmd-[":"indentLess","Cmd-]":"indentMore","Cmd-Backspace":"delWrappedLineLeft","Cmd-Delete":"delWrappedLineRight","Cmd-U":"undoSelection","Shift-Cmd-U":"redoSelection","Ctrl-Up":"goDocStart","Ctrl-Down":"goDocEnd",fallthrough:["basic","emacsy"]},Vo.default=y?Vo.macDefault:Vo.pcDefault;var tl={selectAll:so,singleSelection:function(e){return e.setSelection(e.getCursor("anchor"),e.getCursor("head"),U)},killLine:function(e){return Zo(e,(function(t){if(t.empty()){var r=Xe(e.doc,t.head.line).text.length;return t.head.ch==r&&t.head.line<e.lastLine()?{from:t.head,to:et(t.head.line+1,0)}:{from:t.head,to:et(t.head.line,r)}}return{from:t.from(),to:t.to()}}))},deleteLine:function(e){return Zo(e,(function(t){return{from:et(t.from().line,0),to:st(e.doc,et(t.to().line+1,0))}}))},delLineLeft:function(e){return Zo(e,(function(e){return{from:et(e.from().line,0),to:e.from()}}))},delWrappedLineLeft:function(e){return Zo(e,(function(t){var r=e.charCoords(t.head,"div").top+5;return{from:e.coordsChar({left:0,top:r},"div"),to:t.from()}}))},delWrappedLineRight:function(e){return Zo(e,(function(t){var r=e.charCoords(t.head,"div").top+5,n=e.coordsChar({left:e.display.lineDiv.offsetWidth+100,top:r},"div");return{from:t.from(),to:n}}))},undo:function(e){return e.undo()},redo:function(e){return e.redo()},undoSelection:function(e){return e.undoSelection()},redoSelection:function(e){return e.redoSelection()},goDocStart:function(e){return e.extendSelection(et(e.firstLine(),0))},goDocEnd:function(e){return e.extendSelection(et(e.lastLine()))},goLineStart:function(e){return e.extendSelectionsBy((function(t){return rl(e,t.head.line)}),{origin:"+move",bias:1})},goLineStartSmart:function(e){return e.extendSelectionsBy((function(t){return nl(e,t.head)}),{origin:"+move",bias:1})},goLineEnd:function(e){return e.extendSelectionsBy((function(t){return function(e,t){var r=Xe(e.doc,t),n=function(e){for(var t;t=Pt(e);)e=t.find(1,!0).line;return e}(r);n!=r&&(t=qe(n));return el(!0,e,r,t,-1)}(e,t.head.line)}),{origin:"+move",bias:-1})},goLineRight:function(e){return e.extendSelectionsBy((function(t){var r=e.cursorCoords(t.head,"div").top+5;return e.coordsChar({left:e.display.lineDiv.offsetWidth+100,top:r},"div")}),K)},goLineLeft:function(e){return e.extendSelectionsBy((function(t){var r=e.cursorCoords(t.head,"div").top+5;return e.coordsChar({left:0,top:r},"div")}),K)},goLineLeftSmart:function(e){return e.extendSelectionsBy((function(t){var r=e.cursorCoords(t.head,"div").top+5,n=e.coordsChar({left:0,top:r},"div");return n.ch<e.getLine(n.line).search(/\S/)?nl(e,t.head):n}),K)},goLineUp:function(e){return e.moveV(-1,"line")},goLineDown:function(e){return e.moveV(1,"line")},goPageUp:function(e){return e.moveV(-1,"page")},goPageDown:function(e){return e.moveV(1,"page")},goCharLeft:function(e){return e.moveH(-1,"char")},goCharRight:function(e){return e.moveH(1,"char")},goColumnLeft:function(e){return e.moveH(-1,"column")},goColumnRight:function(e){return e.moveH(1,"column")},goWordLeft:function(e){return e.moveH(-1,"word")},goGroupRight:function(e){return e.moveH(1,"group")},goGroupLeft:function(e){return e.moveH(-1,"group")},goWordRight:function(e){return e.moveH(1,"word")},delCharBefore:function(e){return e.deleteH(-1,"codepoint")},delCharAfter:function(e){return e.deleteH(1,"char")},delWordBefore:function(e){return e.deleteH(-1,"word")},delWordAfter:function(e){return e.deleteH(1,"word")},delGroupBefore:function(e){return e.deleteH(-1,"group")},delGroupAfter:function(e){return e.deleteH(1,"group")},indentAuto:function(e){return e.indentSelection("smart")},indentMore:function(e){return e.indentSelection("add")},indentLess:function(e){return e.indentSelection("subtract")},insertTab:function(e){return e.replaceSelection("\t")},insertSoftTab:function(e){for(var t=[],r=e.listSelections(),n=e.options.tabSize,i=0;i<r.length;i++){var o=r[i].from(),l=R(e.getLine(o.line),o.ch,n);t.push(Y(n-l%n))}e.replaceSelections(t)},defaultTab:function(e){e.somethingSelected()?e.indentSelection("add"):e.execCommand("insertTab")},transposeChars:function(e){return Jn(e,(function(){for(var t=e.listSelections(),r=[],n=0;n<t.length;n++)if(t[n].empty()){var i=t[n].head,o=Xe(e.doc,i.line).text;if(o)if(i.ch==o.length&&(i=new et(i.line,i.ch-1)),i.ch>0)i=new et(i.line,i.ch+1),e.replaceRange(o.charAt(i.ch-1)+o.charAt(i.ch-2),et(i.line,i.ch-2),i,"+transpose");else if(i.line>e.doc.first){var l=Xe(e.doc,i.line-1).text;l&&(i=new et(i.line,1),e.replaceRange(o.charAt(0)+e.doc.lineSeparator()+l.charAt(l.length-1),et(i.line-1,l.length-1),i,"+transpose"))}r.push(new Si(i,i))}e.setSelections(r)}))},newlineAndIndent:function(e){return Jn(e,(function(){for(var t=e.listSelections(),r=t.length-1;r>=0;r--)e.replaceRange(e.doc.lineSeparator(),t[r].anchor,t[r].head,"+input");t=e.listSelections();for(var n=0;n<t.length;n++)e.indentLine(t[n].from().line,null,!0);Dn(e)}))},openLine:function(e){return e.replaceSelection("\n","start")},toggleOverwrite:function(e){return e.toggleOverwrite()}};function rl(e,t){var r=Xe(e.doc,t),n=Rt(r);return n!=r&&(t=qe(n)),el(!0,e,n,t,1)}function nl(e,t){var r=rl(e,t.line),n=Xe(e.doc,r.line),i=ue(n,e.doc.direction);if(!i||0==i[0].level){var o=Math.max(r.ch,n.text.search(/\S/)),l=t.line==r.line&&t.ch<=o&&t.ch;return et(r.line,l?0:o,r.sticky)}return r}function il(e,t,r){if("string"==typeof t&&!(t=tl[t]))return!1;e.display.input.ensurePolled();var n=e.display.shift,i=!1;try{e.isReadOnly()&&(e.state.suppressEdits=!0),r&&(e.display.shift=!1),i=t(e)!=G}finally{e.display.shift=n,e.state.suppressEdits=!1}return i}var ol=new z;function ll(e,t,r,n){var i=e.state.keySeq;if(i){if(Yo(t))return"handled";if(/\'$/.test(t)?e.state.keySeq=null:ol.set(50,(function(){e.state.keySeq==i&&(e.state.keySeq=null,e.display.input.reset())})),sl(e,i+" "+t,r,n))return!0}return sl(e,t,r,n)}function sl(e,t,r,n){var i=function(e,t,r){for(var n=0;n<e.state.keyMaps.length;n++){var i=Xo(t,e.state.keyMaps[n],r,e);if(i)return i}return e.options.extraKeys&&Xo(t,e.options.extraKeys,r,e)||Xo(t,e.options.keyMap,r,e)}(e,t,n);return"multi"==i&&(e.state.keySeq=t),"handled"==i&&sr(e,"keyHandled",e,t,r),"handled"!=i&&"multi"!=i||(be(r),xn(e)),!!i}function al(e,t){var r=_o(t,!0);return!!r&&(t.shiftKey&&!e.state.keySeq?ll(e,"Shift-"+r,t,(function(t){return il(e,t,!0)}))||ll(e,r,t,(function(t){if("string"==typeof t?/^go[A-Z]/.test(t):t.motion)return il(e,t)})):ll(e,r,t,(function(t){return il(e,t)})))}var ul=null;function cl(e){var t=this;if(!(e.target&&e.target!=t.display.input.getField()||(t.curOp.focus=W(),ge(t,e)))){l&&s<11&&27==e.keyCode&&(e.returnValue=!1);var n=e.keyCode;t.display.shift=16==n||e.shiftKey;var i=al(t,e);h&&(ul=i?n:null,i||88!=n||He||!(y?e.metaKey:e.ctrlKey)||t.replaceSelection("",null,"cut")),r&&!y&&!i&&46==n&&e.shiftKey&&!e.ctrlKey&&document.execCommand&&document.execCommand("cut"),18!=n||/\bCodeMirror-crosshair\b/.test(t.display.lineDiv.className)||function(e){var t=e.display.lineDiv;function r(e){18!=e.keyCode&&e.altKey||(T(t,"CodeMirror-crosshair"),de(document,"keyup",r),de(document,"mouseover",r))}H(t,"CodeMirror-crosshair"),he(document,"keyup",r),he(document,"mouseover",r)}(t)}}function hl(e){16==e.keyCode&&(this.doc.sel.shift=!1),ge(this,e)}function fl(e){var t=this;if(!(e.target&&e.target!=t.display.input.getField()||xr(t.display,e)||ge(t,e)||e.ctrlKey&&!e.altKey||y&&e.metaKey)){var r=e.keyCode,n=e.charCode;if(h&&r==ul)return ul=null,void be(e);if(!h||e.which&&!(e.which<10)||!al(t,e)){var i=String.fromCharCode(null==n?r:n);"\b"!=i&&(function(e,t,r){return ll(e,"'"+r+"'",t,(function(t){return il(e,t,!0)}))}(t,e,i)||t.display.input.onKeyPress(e))}}}var dl,pl,gl=function(e,t,r){this.time=e,this.pos=t,this.button=r};function vl(e){var t=this,r=t.display;if(!(ge(t,e)||r.activeTouch&&r.input.supportsTouch()))if(r.input.ensurePolled(),r.shift=e.shiftKey,xr(r,e))a||(r.scroller.draggable=!1,setTimeout((function(){return r.scroller.draggable=!0}),100));else if(!bl(t,e)){var n=un(t,e),i=Le(e),o=n?function(e,t){var r=+new Date;return pl&&pl.compare(r,e,t)?(dl=pl=null,"triple"):dl&&dl.compare(r,e,t)?(pl=new gl(r,e,t),dl=null,"double"):(dl=new gl(r,e,t),pl=null,"single")}(n,i):"single";window.focus(),1==i&&t.state.selectingText&&t.state.selectingText(e),n&&function(e,t,r,n,i){var o="Click";"double"==n?o="Double"+o:"triple"==n&&(o="Triple"+o);return ll(e,$o(o=(1==t?"Left":2==t?"Middle":"Right")+o,i),i,(function(t){if("string"==typeof t&&(t=tl[t]),!t)return!1;var n=!1;try{e.isReadOnly()&&(e.state.suppressEdits=!0),n=t(e,r)!=G}finally{e.state.suppressEdits=!1}return n}))}(t,i,n,o,e)||(1==i?n?function(e,t,r,n){l?setTimeout(E(Cn,e),0):e.curOp.focus=W();var i,o=function(e,t,r){var n=e.getOption("configureMouse"),i=n?n(e,t,r):{};if(null==i.unit){var o=b?r.shiftKey&&r.metaKey:r.altKey;i.unit=o?"rectangle":"single"==t?"char":"double"==t?"word":"line"}(null==i.extend||e.doc.extend)&&(i.extend=e.doc.extend||r.shiftKey);null==i.addNew&&(i.addNew=y?r.metaKey:r.ctrlKey);null==i.moveOnDrag&&(i.moveOnDrag=!(y?r.altKey:r.ctrlKey));return i}(e,r,n),u=e.doc.sel;e.options.dragDrop&&Me&&!e.isReadOnly()&&"single"==r&&(i=u.contains(t))>-1&&(tt((i=u.ranges[i]).from(),t)<0||t.xRel>0)&&(tt(i.to(),t)>0||t.xRel<0)?function(e,t,r,n){var i=e.display,o=!1,u=ei(e,(function(t){a&&(i.scroller.draggable=!1),e.state.draggingText=!1,e.state.delayingBlurEvent&&(e.hasFocus()?e.state.delayingBlurEvent=!1:Sn(e)),de(i.wrapper.ownerDocument,"mouseup",u),de(i.wrapper.ownerDocument,"mousemove",c),de(i.scroller,"dragstart",h),de(i.scroller,"drop",u),o||(be(t),n.addNew||$i(e.doc,r,null,null,n.extend),a&&!f||l&&9==s?setTimeout((function(){i.wrapper.ownerDocument.body.focus({preventScroll:!0}),i.input.focus()}),20):i.input.focus())})),c=function(e){o=o||Math.abs(t.clientX-e.clientX)+Math.abs(t.clientY-e.clientY)>=10},h=function(){return o=!0};a&&(i.scroller.draggable=!0);e.state.draggingText=u,u.copy=!n.moveOnDrag,he(i.wrapper.ownerDocument,"mouseup",u),he(i.wrapper.ownerDocument,"mousemove",c),he(i.scroller,"dragstart",h),he(i.scroller,"drop",u),e.state.delayingBlurEvent=!0,setTimeout((function(){return i.input.focus()}),20),i.scroller.dragDrop&&i.scroller.dragDrop()}(e,n,t,o):function(e,t,r,n){l&&Sn(e);var i=e.display,o=e.doc;be(t);var s,a,u=o.sel,c=u.ranges;n.addNew&&!n.extend?(a=o.sel.contains(r),s=a>-1?c[a]:new Si(r,r)):(s=o.sel.primary(),a=o.sel.primIndex);if("rectangle"==n.unit)n.addNew||(s=new Si(r,r)),r=un(e,t,!0,!0),a=-1;else{var h=ml(e,r,n.unit);s=n.extend?Yi(s,h.anchor,h.head,n.extend):h}n.addNew?-1==a?(a=c.length,Ji(o,Li(e,c.concat([s]),a),{scroll:!1,origin:"*mouse"})):c.length>1&&c[a].empty()&&"char"==n.unit&&!n.extend?(Ji(o,Li(e,c.slice(0,a).concat(c.slice(a+1)),0),{scroll:!1,origin:"*mouse"}),u=o.sel):qi(o,a,s,V):(a=0,Ji(o,new Ci([s],0),V),u=o.sel);var f=r;function d(t){if(0!=tt(f,t))if(f=t,"rectangle"==n.unit){for(var i=[],l=e.options.tabSize,c=R(Xe(o,r.line).text,r.ch,l),h=R(Xe(o,t.line).text,t.ch,l),d=Math.min(c,h),p=Math.max(c,h),g=Math.min(r.line,t.line),v=Math.min(e.lastLine(),Math.max(r.line,t.line));g<=v;g++){var m=Xe(o,g).text,y=j(m,d,l);d==p?i.push(new Si(et(g,y),et(g,y))):m.length>y&&i.push(new Si(et(g,y),et(g,j(m,p,l))))}i.length||i.push(new Si(r,r)),Ji(o,Li(e,u.ranges.slice(0,a).concat(i),a),{origin:"*mouse",scroll:!1}),e.scrollIntoView(t)}else{var b,w=s,x=ml(e,t,n.unit),C=w.anchor;tt(x.anchor,C)>0?(b=x.head,C=ot(w.from(),x.anchor)):(b=x.anchor,C=it(w.to(),x.head));var S=u.ranges.slice(0);S[a]=function(e,t){var r=t.anchor,n=t.head,i=Xe(e.doc,r.line);if(0==tt(r,n)&&r.sticky==n.sticky)return t;var o=ue(i);if(!o)return t;var l=se(o,r.ch,r.sticky),s=o[l];if(s.from!=r.ch&&s.to!=r.ch)return t;var a,u=l+(s.from==r.ch==(1!=s.level)?0:1);if(0==u||u==o.length)return t;if(n.line!=r.line)a=(n.line-r.line)*("ltr"==e.doc.direction?1:-1)>0;else{var c=se(o,n.ch,n.sticky),h=c-l||(n.ch-r.ch)*(1==s.level?-1:1);a=c==u-1||c==u?h<0:h>0}var f=o[u+(a?-1:0)],d=a==(1==f.level),p=d?f.from:f.to,g=d?"after":"before";return r.ch==p&&r.sticky==g?t:new Si(new et(r.line,p,g),n)}(e,new Si(st(o,C),b)),Ji(o,Li(e,S,a),V)}}var p=i.wrapper.getBoundingClientRect(),g=0;function v(t){var r=++g,l=un(e,t,!0,"rectangle"==n.unit);if(l)if(0!=tt(l,f)){e.curOp.focus=W(),d(l);var s=Nn(i,o);(l.line>=s.to||l.line<s.from)&&setTimeout(ei(e,(function(){g==r&&v(t)})),150)}else{var a=t.clientY<p.top?-20:t.clientY>p.bottom?20:0;a&&setTimeout(ei(e,(function(){g==r&&(i.scroller.scrollTop+=a,v(t))})),50)}}function m(t){e.state.selectingText=!1,g=1/0,t&&(be(t),i.input.focus()),de(i.wrapper.ownerDocument,"mousemove",y),de(i.wrapper.ownerDocument,"mouseup",b),o.history.lastSelOrigin=null}var y=ei(e,(function(e){0!==e.buttons&&Le(e)?v(e):m(e)})),b=ei(e,m);e.state.selectingText=b,he(i.wrapper.ownerDocument,"mousemove",y),he(i.wrapper.ownerDocument,"mouseup",b)}(e,n,t,o)}(t,n,o,e):Se(e)==r.scroller&&be(e):2==i?(n&&$i(t.doc,n),setTimeout((function(){return r.input.focus()}),20)):3==i&&(S?t.display.input.onContextMenu(e):Sn(t)))}}function ml(e,t,r){if("char"==r)return new Si(t,t);if("word"==r)return e.findWordAt(t);if("line"==r)return new Si(et(t.line,0),st(e.doc,et(t.line+1,0)));var n=r(e,t);return new Si(n.from,n.to)}function yl(e,t,r,n){var i,o;if(t.touches)i=t.touches[0].clientX,o=t.touches[0].clientY;else try{i=t.clientX,o=t.clientY}catch(e){return!1}if(i>=Math.floor(e.display.gutters.getBoundingClientRect().right))return!1;n&&be(t);var l=e.display,s=l.lineDiv.getBoundingClientRect();if(o>s.bottom||!me(e,r))return xe(t);o-=s.top-l.viewOffset;for(var a=0;a<e.display.gutterSpecs.length;++a){var u=l.gutters.childNodes[a];if(u&&u.getBoundingClientRect().right>=i)return pe(e,r,e,Ze(e.doc,o),e.display.gutterSpecs[a].className,t),xe(t)}}function bl(e,t){return yl(e,t,"gutterClick",!0)}function wl(e,t){xr(e.display,t)||function(e,t){if(!me(e,"gutterContextMenu"))return!1;return yl(e,t,"gutterContextMenu",!1)}(e,t)||ge(e,t,"contextmenu")||S||e.display.input.onContextMenu(t)}function xl(e){e.display.wrapper.className=e.display.wrapper.className.replace(/\s*cm-s-\S+/g,"")+e.options.theme.replace(/(^|\s)\s*/g," cm-s-"),zr(e)}gl.prototype.compare=function(e,t,r){return this.time+400>e&&0==tt(t,this.pos)&&r==this.button};var Cl={toString:function(){return"CodeMirror.Init"}},Sl={},Ll={};function kl(e,t,r){if(!t!=!(r&&r!=Cl)){var n=e.display.dragFunctions,i=t?he:de;i(e.display.scroller,"dragstart",n.start),i(e.display.scroller,"dragenter",n.enter),i(e.display.scroller,"dragover",n.over),i(e.display.scroller,"dragleave",n.leave),i(e.display.scroller,"drop",n.drop)}}function Tl(e){e.options.lineWrapping?(H(e.display.wrapper,"CodeMirror-wrap"),e.display.sizer.style.minWidth="",e.display.sizerWidth=null):(T(e.display.wrapper,"CodeMirror-wrap"),jt(e)),an(e),hn(e),zr(e),setTimeout((function(){return Gn(e)}),100)}function Ml(e,t){var r=this;if(!(this instanceof Ml))return new Ml(e,t);this.options=t=t?I(t):{},I(Sl,t,!1);var n=t.value;"string"==typeof n?n=new Do(n,t.mode,null,t.lineSeparator,t.direction):t.mode&&(n.modeOption=t.mode),this.doc=n;var i=new Ml.inputStyles[t.inputStyle](this),o=this.display=new vi(e,n,i,t);for(var u in o.wrapper.CodeMirror=this,xl(this),t.lineWrapping&&(this.display.wrapper.className+=" CodeMirror-wrap"),Kn(this),this.state={keyMaps:[],overlays:[],modeGen:0,overwrite:!1,delayingBlurEvent:!1,focused:!1,suppressEdits:!1,pasteIncoming:-1,cutIncoming:-1,selectingText:!1,draggingText:!1,highlight:new z,keySeq:null,specialChars:null},t.autofocus&&!m&&o.input.focus(),l&&s<11&&setTimeout((function(){return r.display.input.reset(!0)}),20),function(e){var t=e.display;he(t.scroller,"mousedown",ei(e,vl)),he(t.scroller,"dblclick",l&&s<11?ei(e,(function(t){if(!ge(e,t)){var r=un(e,t);if(r&&!bl(e,t)&&!xr(e.display,t)){be(t);var n=e.findWordAt(r);$i(e.doc,n.anchor,n.head)}}})):function(t){return ge(e,t)||be(t)});he(t.scroller,"contextmenu",(function(t){return wl(e,t)})),he(t.input.getField(),"contextmenu",(function(r){t.scroller.contains(r.target)||wl(e,r)}));var r,n={end:0};function i(){t.activeTouch&&(r=setTimeout((function(){return t.activeTouch=null}),1e3),(n=t.activeTouch).end=+new Date)}function o(e){if(1!=e.touches.length)return!1;var t=e.touches[0];return t.radiusX<=1&&t.radiusY<=1}function a(e,t){if(null==t.left)return!0;var r=t.left-e.left,n=t.top-e.top;return r*r+n*n>400}he(t.scroller,"touchstart",(function(i){if(!ge(e,i)&&!o(i)&&!bl(e,i)){t.input.ensurePolled(),clearTimeout(r);var l=+new Date;t.activeTouch={start:l,moved:!1,prev:l-n.end<=300?n:null},1==i.touches.length&&(t.activeTouch.left=i.touches[0].pageX,t.activeTouch.top=i.touches[0].pageY)}})),he(t.scroller,"touchmove",(function(){t.activeTouch&&(t.activeTouch.moved=!0)})),he(t.scroller,"touchend",(function(r){var n=t.activeTouch;if(n&&!xr(t,r)&&null!=n.left&&!n.moved&&new Date-n.start<300){var o,l=e.coordsChar(t.activeTouch,"page");o=!n.prev||a(n,n.prev)?new Si(l,l):!n.prev.prev||a(n,n.prev.prev)?e.findWordAt(l):new Si(et(l.line,0),st(e.doc,et(l.line+1,0))),e.setSelection(o.anchor,o.head),e.focus(),be(r)}i()})),he(t.scroller,"touchcancel",i),he(t.scroller,"scroll",(function(){t.scroller.clientHeight&&(Pn(e,t.scroller.scrollTop),In(e,t.scroller.scrollLeft,!0),pe(e,"scroll",e))})),he(t.scroller,"mousewheel",(function(t){return xi(e,t)})),he(t.scroller,"DOMMouseScroll",(function(t){return xi(e,t)})),he(t.wrapper,"scroll",(function(){return t.wrapper.scrollTop=t.wrapper.scrollLeft=0})),t.dragFunctions={enter:function(t){ge(e,t)||Ce(t)},over:function(t){ge(e,t)||(!function(e,t){var r=un(e,t);if(r){var n=document.createDocumentFragment();yn(e,r,n),e.display.dragCursor||(e.display.dragCursor=O("div",null,"CodeMirror-cursors CodeMirror-dragcursors"),e.display.lineSpace.insertBefore(e.display.dragCursor,e.display.cursorDiv)),N(e.display.dragCursor,n)}}(e,t),Ce(t))},start:function(t){return function(e,t){if(l&&(!e.state.draggingText||+new Date-Wo<100))Ce(t);else if(!ge(e,t)&&!xr(e.display,t)&&(t.dataTransfer.setData("Text",e.getSelection()),t.dataTransfer.effectAllowed="copyMove",t.dataTransfer.setDragImage&&!f)){var r=O("img",null,null,"position: fixed; left: 0; top: 0;");r.src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",h&&(r.width=r.height=1,e.display.wrapper.appendChild(r),r._top=r.offsetTop),t.dataTransfer.setDragImage(r,0,0),h&&r.parentNode.removeChild(r)}}(e,t)},drop:ei(e,Ho),leave:function(t){ge(e,t)||Fo(e)}};var u=t.input.getField();he(u,"keyup",(function(t){return hl.call(e,t)})),he(u,"keydown",ei(e,cl)),he(u,"keypress",ei(e,fl)),he(u,"focus",(function(t){return Ln(e,t)})),he(u,"blur",(function(t){return kn(e,t)}))}(this),Io(),Xn(this),this.curOp.forceUpdate=!0,Pi(this,n),t.autofocus&&!m||this.hasFocus()?setTimeout((function(){r.hasFocus()&&!r.state.focused&&Ln(r)}),20):kn(this),Ll)Ll.hasOwnProperty(u)&&Ll[u](this,t[u],Cl);fi(this),t.finishInit&&t.finishInit(this);for(var c=0;c<Nl.length;++c)Nl[c](this);Yn(this),a&&t.lineWrapping&&"optimizelegibility"==getComputedStyle(o.lineDiv).textRendering&&(o.lineDiv.style.textRendering="auto")}Ml.defaults=Sl,Ml.optionHandlers=Ll;var Nl=[];function Ol(e,t,r,n){var i,o=e.doc;null==r&&(r="add"),"smart"==r&&(o.mode.indent?i=dt(e,t).state:r="prev");var l=e.options.tabSize,s=Xe(o,t),a=R(s.text,null,l);s.stateAfter&&(s.stateAfter=null);var u,c=s.text.match(/^\s*/)[0];if(n||/\S/.test(s.text)){if("smart"==r&&((u=o.mode.indent(i,s.text.slice(c.length),s.text))==G||u>150)){if(!n)return;r="prev"}}else u=0,r="not";"prev"==r?u=t>o.first?R(Xe(o,t-1).text,null,l):0:"add"==r?u=a+e.options.indentUnit:"subtract"==r?u=a-e.options.indentUnit:"number"==typeof r&&(u=a+r),u=Math.max(0,u);var h="",f=0;if(e.options.indentWithTabs)for(var d=Math.floor(u/l);d;--d)f+=l,h+="\t";if(f<u&&(h+=Y(u-f)),h!=c)return go(o,h,et(t,0),et(t,c.length),"+input"),s.stateAfter=null,!0;for(var p=0;p<o.sel.ranges.length;p++){var g=o.sel.ranges[p];if(g.head.line==t&&g.head.ch<c.length){var v=et(t,c.length);qi(o,p,new Si(v,v));break}}}Ml.defineInitHook=function(e){return Nl.push(e)};var Al=null;function Dl(e){Al=e}function Wl(e,t,r,n,i){var o=e.doc;e.display.shift=!1,n||(n=o.sel);var l=+new Date-200,s="paste"==i||e.state.pasteIncoming>l,a=De(t),u=null;if(s&&n.ranges.length>1)if(Al&&Al.text.join("\n")==t){if(n.ranges.length%Al.text.length==0){u=[];for(var c=0;c<Al.text.length;c++)u.push(o.splitLines(Al.text[c]))}}else a.length==n.ranges.length&&e.options.pasteLinesPerSelection&&(u=_(a,(function(e){return[e]})));for(var h=e.curOp.updateInput,f=n.ranges.length-1;f>=0;f--){var d=n.ranges[f],p=d.from(),g=d.to();d.empty()&&(r&&r>0?p=et(p.line,p.ch-r):e.state.overwrite&&!s?g=et(g.line,Math.min(Xe(o,g.line).text.length,g.ch+$(a).length)):s&&Al&&Al.lineWise&&Al.text.join("\n")==a.join("\n")&&(p=g=et(p.line,0)));var v={from:p,to:g,text:u?u[f%u.length]:a,origin:i||(s?"paste":e.state.cutIncoming>l?"cut":"+input")};uo(e.doc,v),sr(e,"inputRead",e,v)}t&&!s&&Fl(e,t),Dn(e),e.curOp.updateInput<2&&(e.curOp.updateInput=h),e.curOp.typing=!0,e.state.pasteIncoming=e.state.cutIncoming=-1}function Hl(e,t){var r=e.clipboardData&&e.clipboardData.getData("Text");if(r)return e.preventDefault(),t.isReadOnly()||t.options.disableInput||!t.hasFocus()||Jn(t,(function(){return Wl(t,r,0,null,"paste")})),!0}function Fl(e,t){if(e.options.electricChars&&e.options.smartIndent)for(var r=e.doc.sel,n=r.ranges.length-1;n>=0;n--){var i=r.ranges[n];if(!(i.head.ch>100||n&&r.ranges[n-1].head.line==i.head.line)){var o=e.getModeAt(i.head),l=!1;if(o.electricChars){for(var s=0;s<o.electricChars.length;s++)if(t.indexOf(o.electricChars.charAt(s))>-1){l=Ol(e,i.head.line,"smart");break}}else o.electricInput&&o.electricInput.test(Xe(e.doc,i.head.line).text.slice(0,i.head.ch))&&(l=Ol(e,i.head.line,"smart"));l&&sr(e,"electricInput",e,i.head.line)}}}function Pl(e){for(var t=[],r=[],n=0;n<e.doc.sel.ranges.length;n++){var i=e.doc.sel.ranges[n].head.line,o={anchor:et(i,0),head:et(i+1,0)};r.push(o),t.push(e.getRange(o.anchor,o.head))}return{text:t,ranges:r}}function El(e,t,r,n){e.setAttribute("autocorrect",r?"":"off"),e.setAttribute("autocapitalize",n?"":"off"),e.setAttribute("spellcheck",!!t)}function Il(){var e=O("textarea",null,null,"position: absolute; bottom: -1em; padding: 0; width: 1px; height: 1em; min-height: 1em; outline: none"),t=O("div",[e],null,"overflow: hidden; position: relative; width: 3px; height: 0px;");return a?e.style.width="1000px":e.setAttribute("wrap","off"),g&&(e.style.border="1px solid black"),El(e),t}function Rl(e,t,r,n,i){var o=t,l=r,s=Xe(e,t.line),a=i&&"rtl"==e.direction?-r:r;function u(o){var l,u;if("codepoint"==n){var c=s.text.charCodeAt(t.ch+(r>0?0:-1));if(isNaN(c))l=null;else{var h=r>0?c>=55296&&c<56320:c>=56320&&c<57343;l=new et(t.line,Math.max(0,Math.min(s.text.length,t.ch+r*(h?2:1))),-r)}}else l=i?function(e,t,r,n){var i=ue(t,e.doc.direction);if(!i)return Jo(t,r,n);r.ch>=t.text.length?(r.ch=t.text.length,r.sticky="before"):r.ch<=0&&(r.ch=0,r.sticky="after");var o=se(i,r.ch,r.sticky),l=i[o];if("ltr"==e.doc.direction&&l.level%2==0&&(n>0?l.to>r.ch:l.from<r.ch))return Jo(t,r,n);var s,a=function(e,r){return Qo(t,e instanceof et?e.ch:e,r)},u=function(r){return e.options.lineWrapping?(s=s||Dr(e,t),Zr(e,t,s,r)):{begin:0,end:t.text.length}},c=u("before"==r.sticky?a(r,-1):r.ch);if("rtl"==e.doc.direction||1==l.level){var h=1==l.level==n<0,f=a(r,h?1:-1);if(null!=f&&(h?f<=l.to&&f<=c.end:f>=l.from&&f>=c.begin)){var d=h?"before":"after";return new et(r.line,f,d)}}var p=function(e,t,n){for(var o=function(e,t){return t?new et(r.line,a(e,1),"before"):new et(r.line,e,"after")};e>=0&&e<i.length;e+=t){var l=i[e],s=t>0==(1!=l.level),u=s?n.begin:a(n.end,-1);if(l.from<=u&&u<l.to)return o(u,s);if(u=s?l.from:a(l.to,-1),n.begin<=u&&u<n.end)return o(u,s)}},g=p(o+n,n,c);if(g)return g;var v=n>0?c.end:a(c.begin,-1);return null==v||n>0&&v==t.text.length||!(g=p(n>0?0:i.length-1,n,u(v)))?null:g}(e.cm,s,t,r):Jo(s,t,r);if(null==l){if(o||(u=t.line+a)<e.first||u>=e.first+e.size||(t=new et(u,t.ch,t.sticky),!(s=Xe(e,u))))return!1;t=el(i,e.cm,s,t.line,a)}else t=l;return!0}if("char"==n||"codepoint"==n)u();else if("column"==n)u(!0);else if("word"==n||"group"==n)for(var c=null,h="group"==n,f=e.cm&&e.cm.getHelper(t,"wordChars"),d=!0;!(r<0)||u(!d);d=!1){var p=s.text.charAt(t.ch)||"\n",g=ee(p,f)?"w":h&&"\n"==p?"n":!h||/\s/.test(p)?null:"p";if(!h||d||g||(g="s"),c&&c!=g){r<0&&(r=1,u(),t.sticky="after");break}if(g&&(c=g),r>0&&!u(!d))break}var v=oo(e,t,o,l,!0);return rt(o,v)&&(v.hitSide=!0),v}function zl(e,t,r,n){var i,o,l=e.doc,s=t.left;if("page"==n){var a=Math.min(e.display.wrapper.clientHeight,window.innerHeight||document.documentElement.clientHeight),u=Math.max(a-.5*rn(e.display),3);i=(r>0?t.bottom:t.top)+r*u}else"line"==n&&(i=r>0?t.bottom+3:t.top-3);for(;(o=_r(e,s,i)).outside;){if(r<0?i<=0:i>=l.height){o.hitSide=!0;break}i+=5*r}return o}var Bl=function(e){this.cm=e,this.lastAnchorNode=this.lastAnchorOffset=this.lastFocusNode=this.lastFocusOffset=null,this.polling=new z,this.composing=null,this.gracePeriod=!1,this.readDOMTimeout=null};function Gl(e,t){var r=Ar(e,t.line);if(!r||r.hidden)return null;var n=Xe(e.doc,t.line),i=Nr(r,n,t.line),o=ue(n,e.doc.direction),l="left";o&&(l=se(o,t.ch)%2?"right":"left");var s=Pr(i.map,t.ch,l);return s.offset="right"==s.collapse?s.end:s.start,s}function Ul(e,t){return t&&(e.bad=!0),e}function Vl(e,t,r){var n;if(t==e.display.lineDiv){if(!(n=e.display.lineDiv.childNodes[r]))return Ul(e.clipPos(et(e.display.viewTo-1)),!0);t=null,r=0}else for(n=t;;n=n.parentNode){if(!n||n==e.display.lineDiv)return null;if(n.parentNode&&n.parentNode==e.display.lineDiv)break}for(var i=0;i<e.display.view.length;i++){var o=e.display.view[i];if(o.node==n)return Kl(o,t,r)}}function Kl(e,t,r){var n=e.text.firstChild,i=!1;if(!t||!D(n,t))return Ul(et(qe(e.line),0),!0);if(t==n&&(i=!0,t=n.childNodes[r],r=0,!t)){var o=e.rest?$(e.rest):e.line;return Ul(et(qe(o),o.text.length),i)}var l=3==t.nodeType?t:null,s=t;for(l||1!=t.childNodes.length||3!=t.firstChild.nodeType||(l=t.firstChild,r&&(r=l.nodeValue.length));s.parentNode!=n;)s=s.parentNode;var a=e.measure,u=a.maps;function c(t,r,n){for(var i=-1;i<(u?u.length:0);i++)for(var o=i<0?a.map:u[i],l=0;l<o.length;l+=3){var s=o[l+2];if(s==t||s==r){var c=qe(i<0?e.line:e.rest[i]),h=o[l]+n;return(n<0||s!=t)&&(h=o[l+(n?1:0)]),et(c,h)}}}var h=c(l,s,r);if(h)return Ul(h,i);for(var f=s.nextSibling,d=l?l.nodeValue.length-r:0;f;f=f.nextSibling){if(h=c(f,f.firstChild,0))return Ul(et(h.line,h.ch-d),i);d+=f.textContent.length}for(var p=s.previousSibling,g=r;p;p=p.previousSibling){if(h=c(p,p.firstChild,-1))return Ul(et(h.line,h.ch+g),i);g+=p.textContent.length}}Bl.prototype.init=function(e){var t=this,r=this,n=r.cm,i=r.div=e.lineDiv;function o(e){for(var t=e.target;t;t=t.parentNode){if(t==i)return!0;if(/\bCodeMirror-(?:line)?widget\b/.test(t.className))break}return!1}function l(e){if(o(e)&&!ge(n,e)){if(n.somethingSelected())Dl({lineWise:!1,text:n.getSelections()}),"cut"==e.type&&n.replaceSelection("",null,"cut");else{if(!n.options.lineWiseCopyCut)return;var t=Pl(n);Dl({lineWise:!0,text:t.text}),"cut"==e.type&&n.operation((function(){n.setSelections(t.ranges,0,U),n.replaceSelection("",null,"cut")}))}if(e.clipboardData){e.clipboardData.clearData();var l=Al.text.join("\n");if(e.clipboardData.setData("Text",l),e.clipboardData.getData("Text")==l)return void e.preventDefault()}var s=Il(),a=s.firstChild;n.display.lineSpace.insertBefore(s,n.display.lineSpace.firstChild),a.value=Al.text.join("\n");var u=W();P(a),setTimeout((function(){n.display.lineSpace.removeChild(s),u.focus(),u==i&&r.showPrimarySelection()}),50)}}i.contentEditable=!0,El(i,n.options.spellcheck,n.options.autocorrect,n.options.autocapitalize),he(i,"paste",(function(e){!o(e)||ge(n,e)||Hl(e,n)||s<=11&&setTimeout(ei(n,(function(){return t.updateFromDOM()})),20)})),he(i,"compositionstart",(function(e){t.composing={data:e.data,done:!1}})),he(i,"compositionupdate",(function(e){t.composing||(t.composing={data:e.data,done:!1})})),he(i,"compositionend",(function(e){t.composing&&(e.data!=t.composing.data&&t.readFromDOMSoon(),t.composing.done=!0)})),he(i,"touchstart",(function(){return r.forceCompositionEnd()})),he(i,"input",(function(){t.composing||t.readFromDOMSoon()})),he(i,"copy",l),he(i,"cut",l)},Bl.prototype.screenReaderLabelChanged=function(e){e?this.div.setAttribute("aria-label",e):this.div.removeAttribute("aria-label")},Bl.prototype.prepareSelection=function(){var e=mn(this.cm,!1);return e.focus=W()==this.div,e},Bl.prototype.showSelection=function(e,t){e&&this.cm.display.view.length&&((e.focus||t)&&this.showPrimarySelection(),this.showMultipleSelections(e))},Bl.prototype.getSelection=function(){return this.cm.display.wrapper.ownerDocument.getSelection()},Bl.prototype.showPrimarySelection=function(){var e=this.getSelection(),t=this.cm,n=t.doc.sel.primary(),i=n.from(),o=n.to();if(t.display.viewTo==t.display.viewFrom||i.line>=t.display.viewTo||o.line<t.display.viewFrom)e.removeAllRanges();else{var l=Vl(t,e.anchorNode,e.anchorOffset),s=Vl(t,e.focusNode,e.focusOffset);if(!l||l.bad||!s||s.bad||0!=tt(ot(l,s),i)||0!=tt(it(l,s),o)){var a=t.display.view,u=i.line>=t.display.viewFrom&&Gl(t,i)||{node:a[0].measure.map[2],offset:0},c=o.line<t.display.viewTo&&Gl(t,o);if(!c){var h=a[a.length-1].measure,f=h.maps?h.maps[h.maps.length-1]:h.map;c={node:f[f.length-1],offset:f[f.length-2]-f[f.length-3]}}if(u&&c){var d,p=e.rangeCount&&e.getRangeAt(0);try{d=k(u.node,u.offset,c.offset,c.node)}catch(e){}d&&(!r&&t.state.focused?(e.collapse(u.node,u.offset),d.collapsed||(e.removeAllRanges(),e.addRange(d))):(e.removeAllRanges(),e.addRange(d)),p&&null==e.anchorNode?e.addRange(p):r&&this.startGracePeriod()),this.rememberSelection()}else e.removeAllRanges()}}},Bl.prototype.startGracePeriod=function(){var e=this;clearTimeout(this.gracePeriod),this.gracePeriod=setTimeout((function(){e.gracePeriod=!1,e.selectionChanged()&&e.cm.operation((function(){return e.cm.curOp.selectionChanged=!0}))}),20)},Bl.prototype.showMultipleSelections=function(e){N(this.cm.display.cursorDiv,e.cursors),N(this.cm.display.selectionDiv,e.selection)},Bl.prototype.rememberSelection=function(){var e=this.getSelection();this.lastAnchorNode=e.anchorNode,this.lastAnchorOffset=e.anchorOffset,this.lastFocusNode=e.focusNode,this.lastFocusOffset=e.focusOffset},Bl.prototype.selectionInEditor=function(){var e=this.getSelection();if(!e.rangeCount)return!1;var t=e.getRangeAt(0).commonAncestorContainer;return D(this.div,t)},Bl.prototype.focus=function(){"nocursor"!=this.cm.options.readOnly&&(this.selectionInEditor()&&W()==this.div||this.showSelection(this.prepareSelection(),!0),this.div.focus())},Bl.prototype.blur=function(){this.div.blur()},Bl.prototype.getField=function(){return this.div},Bl.prototype.supportsTouch=function(){return!0},Bl.prototype.receivedFocus=function(){var e=this,t=this;this.selectionInEditor()?setTimeout((function(){return e.pollSelection()}),20):Jn(this.cm,(function(){return t.cm.curOp.selectionChanged=!0})),this.polling.set(this.cm.options.pollInterval,(function e(){t.cm.state.focused&&(t.pollSelection(),t.polling.set(t.cm.options.pollInterval,e))}))},Bl.prototype.selectionChanged=function(){var e=this.getSelection();return e.anchorNode!=this.lastAnchorNode||e.anchorOffset!=this.lastAnchorOffset||e.focusNode!=this.lastFocusNode||e.focusOffset!=this.lastFocusOffset},Bl.prototype.pollSelection=function(){if(null==this.readDOMTimeout&&!this.gracePeriod&&this.selectionChanged()){var e=this.getSelection(),t=this.cm;if(v&&c&&this.cm.display.gutterSpecs.length&&function(e){for(var t=e;t;t=t.parentNode)if(/CodeMirror-gutter-wrapper/.test(t.className))return!0;return!1}(e.anchorNode))return this.cm.triggerOnKeyDown({type:"keydown",keyCode:8,preventDefault:Math.abs}),this.blur(),void this.focus();if(!this.composing){this.rememberSelection();var r=Vl(t,e.anchorNode,e.anchorOffset),n=Vl(t,e.focusNode,e.focusOffset);r&&n&&Jn(t,(function(){Ji(t.doc,ki(r,n),U),(r.bad||n.bad)&&(t.curOp.selectionChanged=!0)}))}}},Bl.prototype.pollContent=function(){null!=this.readDOMTimeout&&(clearTimeout(this.readDOMTimeout),this.readDOMTimeout=null);var e,t,r,n=this.cm,i=n.display,o=n.doc.sel.primary(),l=o.from(),s=o.to();if(0==l.ch&&l.line>n.firstLine()&&(l=et(l.line-1,Xe(n.doc,l.line-1).length)),s.ch==Xe(n.doc,s.line).text.length&&s.line<n.lastLine()&&(s=et(s.line+1,0)),l.line<i.viewFrom||s.line>i.viewTo-1)return!1;l.line==i.viewFrom||0==(e=cn(n,l.line))?(t=qe(i.view[0].line),r=i.view[0].node):(t=qe(i.view[e].line),r=i.view[e-1].node.nextSibling);var a,u,c=cn(n,s.line);if(c==i.view.length-1?(a=i.viewTo-1,u=i.lineDiv.lastChild):(a=qe(i.view[c+1].line)-1,u=i.view[c+1].node.previousSibling),!r)return!1;for(var h=n.doc.splitLines(function(e,t,r,n,i){var o="",l=!1,s=e.doc.lineSeparator(),a=!1;function u(e){return function(t){return t.id==e}}function c(){l&&(o+=s,a&&(o+=s),l=a=!1)}function h(e){e&&(c(),o+=e)}function f(t){if(1==t.nodeType){var r=t.getAttribute("cm-text");if(r)return void h(r);var o,d=t.getAttribute("cm-marker");if(d){var p=e.findMarks(et(n,0),et(i+1,0),u(+d));return void(p.length&&(o=p[0].find(0))&&h(Ye(e.doc,o.from,o.to).join(s)))}if("false"==t.getAttribute("contenteditable"))return;var g=/^(pre|div|p|li|table|br)$/i.test(t.nodeName);if(!/^br$/i.test(t.nodeName)&&0==t.textContent.length)return;g&&c();for(var v=0;v<t.childNodes.length;v++)f(t.childNodes[v]);/^(pre|p)$/i.test(t.nodeName)&&(a=!0),g&&(l=!0)}else 3==t.nodeType&&h(t.nodeValue.replace(/\u200b/g,"").replace(/\u00a0/g," "))}for(;f(t),t!=r;)t=t.nextSibling,a=!1;return o}(n,r,u,t,a)),f=Ye(n.doc,et(t,0),et(a,Xe(n.doc,a).text.length));h.length>1&&f.length>1;)if($(h)==$(f))h.pop(),f.pop(),a--;else{if(h[0]!=f[0])break;h.shift(),f.shift(),t++}for(var d=0,p=0,g=h[0],v=f[0],m=Math.min(g.length,v.length);d<m&&g.charCodeAt(d)==v.charCodeAt(d);)++d;for(var y=$(h),b=$(f),w=Math.min(y.length-(1==h.length?d:0),b.length-(1==f.length?d:0));p<w&&y.charCodeAt(y.length-p-1)==b.charCodeAt(b.length-p-1);)++p;if(1==h.length&&1==f.length&&t==l.line)for(;d&&d>l.ch&&y.charCodeAt(y.length-p-1)==b.charCodeAt(b.length-p-1);)d--,p++;h[h.length-1]=y.slice(0,y.length-p).replace(/^\u200b+/,""),h[0]=h[0].slice(d).replace(/\u200b+$/,"");var x=et(t,d),C=et(a,f.length?$(f).length-p:0);return h.length>1||h[0]||tt(x,C)?(go(n.doc,h,x,C,"+input"),!0):void 0},Bl.prototype.ensurePolled=function(){this.forceCompositionEnd()},Bl.prototype.reset=function(){this.forceCompositionEnd()},Bl.prototype.forceCompositionEnd=function(){this.composing&&(clearTimeout(this.readDOMTimeout),this.composing=null,this.updateFromDOM(),this.div.blur(),this.div.focus())},Bl.prototype.readFromDOMSoon=function(){var e=this;null==this.readDOMTimeout&&(this.readDOMTimeout=setTimeout((function(){if(e.readDOMTimeout=null,e.composing){if(!e.composing.done)return;e.composing=null}e.updateFromDOM()}),80))},Bl.prototype.updateFromDOM=function(){var e=this;!this.cm.isReadOnly()&&this.pollContent()||Jn(this.cm,(function(){return hn(e.cm)}))},Bl.prototype.setUneditable=function(e){e.contentEditable="false"},Bl.prototype.onKeyPress=function(e){0==e.charCode||this.composing||(e.preventDefault(),this.cm.isReadOnly()||ei(this.cm,Wl)(this.cm,String.fromCharCode(null==e.charCode?e.keyCode:e.charCode),0))},Bl.prototype.readOnlyChanged=function(e){this.div.contentEditable=String("nocursor"!=e)},Bl.prototype.onContextMenu=function(){},Bl.prototype.resetPosition=function(){},Bl.prototype.needsContentAttribute=!0;var jl=function(e){this.cm=e,this.prevInput="",this.pollingFast=!1,this.polling=new z,this.hasSelection=!1,this.composing=null};jl.prototype.init=function(e){var t=this,r=this,n=this.cm;this.createField(e);var i=this.textarea;function o(e){if(!ge(n,e)){if(n.somethingSelected())Dl({lineWise:!1,text:n.getSelections()});else{if(!n.options.lineWiseCopyCut)return;var t=Pl(n);Dl({lineWise:!0,text:t.text}),"cut"==e.type?n.setSelections(t.ranges,null,U):(r.prevInput="",i.value=t.text.join("\n"),P(i))}"cut"==e.type&&(n.state.cutIncoming=+new Date)}}e.wrapper.insertBefore(this.wrapper,e.wrapper.firstChild),g&&(i.style.width="0px"),he(i,"input",(function(){l&&s>=9&&t.hasSelection&&(t.hasSelection=null),r.poll()})),he(i,"paste",(function(e){ge(n,e)||Hl(e,n)||(n.state.pasteIncoming=+new Date,r.fastPoll())})),he(i,"cut",o),he(i,"copy",o),he(e.scroller,"paste",(function(t){if(!xr(e,t)&&!ge(n,t)){if(!i.dispatchEvent)return n.state.pasteIncoming=+new Date,void r.focus();var o=new Event("paste");o.clipboardData=t.clipboardData,i.dispatchEvent(o)}})),he(e.lineSpace,"selectstart",(function(t){xr(e,t)||be(t)})),he(i,"compositionstart",(function(){var e=n.getCursor("from");r.composing&&r.composing.range.clear(),r.composing={start:e,range:n.markText(e,n.getCursor("to"),{className:"CodeMirror-composing"})}})),he(i,"compositionend",(function(){r.composing&&(r.poll(),r.composing.range.clear(),r.composing=null)}))},jl.prototype.createField=function(e){this.wrapper=Il(),this.textarea=this.wrapper.firstChild},jl.prototype.screenReaderLabelChanged=function(e){e?this.textarea.setAttribute("aria-label",e):this.textarea.removeAttribute("aria-label")},jl.prototype.prepareSelection=function(){var e=this.cm,t=e.display,r=e.doc,n=mn(e);if(e.options.moveInputWithCursor){var i=Xr(e,r.sel.primary().head,"div"),o=t.wrapper.getBoundingClientRect(),l=t.lineDiv.getBoundingClientRect();n.teTop=Math.max(0,Math.min(t.wrapper.clientHeight-10,i.top+l.top-o.top)),n.teLeft=Math.max(0,Math.min(t.wrapper.clientWidth-10,i.left+l.left-o.left))}return n},jl.prototype.showSelection=function(e){var t=this.cm.display;N(t.cursorDiv,e.cursors),N(t.selectionDiv,e.selection),null!=e.teTop&&(this.wrapper.style.top=e.teTop+"px",this.wrapper.style.left=e.teLeft+"px")},jl.prototype.reset=function(e){if(!this.contextMenuPending&&!this.composing){var t=this.cm;if(t.somethingSelected()){this.prevInput="";var r=t.getSelection();this.textarea.value=r,t.state.focused&&P(this.textarea),l&&s>=9&&(this.hasSelection=r)}else e||(this.prevInput=this.textarea.value="",l&&s>=9&&(this.hasSelection=null))}},jl.prototype.getField=function(){return this.textarea},jl.prototype.supportsTouch=function(){return!1},jl.prototype.focus=function(){if("nocursor"!=this.cm.options.readOnly&&(!m||W()!=this.textarea))try{this.textarea.focus()}catch(e){}},jl.prototype.blur=function(){this.textarea.blur()},jl.prototype.resetPosition=function(){this.wrapper.style.top=this.wrapper.style.left=0},jl.prototype.receivedFocus=function(){this.slowPoll()},jl.prototype.slowPoll=function(){var e=this;this.pollingFast||this.polling.set(this.cm.options.pollInterval,(function(){e.poll(),e.cm.state.focused&&e.slowPoll()}))},jl.prototype.fastPoll=function(){var e=!1,t=this;t.pollingFast=!0,t.polling.set(20,(function r(){t.poll()||e?(t.pollingFast=!1,t.slowPoll()):(e=!0,t.polling.set(60,r))}))},jl.prototype.poll=function(){var e=this,t=this.cm,r=this.textarea,n=this.prevInput;if(this.contextMenuPending||!t.state.focused||We(r)&&!n&&!this.composing||t.isReadOnly()||t.options.disableInput||t.state.keySeq)return!1;var i=r.value;if(i==n&&!t.somethingSelected())return!1;if(l&&s>=9&&this.hasSelection===i||y&&/[\uf700-\uf7ff]/.test(i))return t.display.input.reset(),!1;if(t.doc.sel==t.display.selForContextMenu){var o=i.charCodeAt(0);if(8203!=o||n||(n="​"),8666==o)return this.reset(),this.cm.execCommand("undo")}for(var a=0,u=Math.min(n.length,i.length);a<u&&n.charCodeAt(a)==i.charCodeAt(a);)++a;return Jn(t,(function(){Wl(t,i.slice(a),n.length-a,null,e.composing?"*compose":null),i.length>1e3||i.indexOf("\n")>-1?r.value=e.prevInput="":e.prevInput=i,e.composing&&(e.composing.range.clear(),e.composing.range=t.markText(e.composing.start,t.getCursor("to"),{className:"CodeMirror-composing"}))})),!0},jl.prototype.ensurePolled=function(){this.pollingFast&&this.poll()&&(this.pollingFast=!1)},jl.prototype.onKeyPress=function(){l&&s>=9&&(this.hasSelection=null),this.fastPoll()},jl.prototype.onContextMenu=function(e){var t=this,r=t.cm,n=r.display,i=t.textarea;t.contextMenuPending&&t.contextMenuPending();var o=un(r,e),u=n.scroller.scrollTop;if(o&&!h){r.options.resetSelectionOnContextMenu&&-1==r.doc.sel.contains(o)&&ei(r,Ji)(r.doc,ki(o),U);var c,f=i.style.cssText,d=t.wrapper.style.cssText,p=t.wrapper.offsetParent.getBoundingClientRect();if(t.wrapper.style.cssText="position: static",i.style.cssText="position: absolute; width: 30px; height: 30px;\n      top: "+(e.clientY-p.top-5)+"px; left: "+(e.clientX-p.left-5)+"px;\n      z-index: 1000; background: "+(l?"rgba(255, 255, 255, .05)":"transparent")+";\n      outline: none; border-width: 0; outline: none; overflow: hidden; opacity: .05; filter: alpha(opacity=5);",a&&(c=window.scrollY),n.input.focus(),a&&window.scrollTo(null,c),n.input.reset(),r.somethingSelected()||(i.value=t.prevInput=" "),t.contextMenuPending=m,n.selForContextMenu=r.doc.sel,clearTimeout(n.detectingSelectAll),l&&s>=9&&v(),S){Ce(e);var g=function(){de(window,"mouseup",g),setTimeout(m,20)};he(window,"mouseup",g)}else setTimeout(m,50)}function v(){if(null!=i.selectionStart){var e=r.somethingSelected(),o="​"+(e?i.value:"");i.value="⇚",i.value=o,t.prevInput=e?"":"​",i.selectionStart=1,i.selectionEnd=o.length,n.selForContextMenu=r.doc.sel}}function m(){if(t.contextMenuPending==m&&(t.contextMenuPending=!1,t.wrapper.style.cssText=d,i.style.cssText=f,l&&s<9&&n.scrollbars.setScrollTop(n.scroller.scrollTop=u),null!=i.selectionStart)){(!l||l&&s<9)&&v();var e=0,o=function(){n.selForContextMenu==r.doc.sel&&0==i.selectionStart&&i.selectionEnd>0&&"​"==t.prevInput?ei(r,so)(r):e++<10?n.detectingSelectAll=setTimeout(o,500):(n.selForContextMenu=null,n.input.reset())};n.detectingSelectAll=setTimeout(o,200)}}},jl.prototype.readOnlyChanged=function(e){e||this.reset(),this.textarea.disabled="nocursor"==e,this.textarea.readOnly=!!e},jl.prototype.setUneditable=function(){},jl.prototype.needsContentAttribute=!1,function(e){var t=e.optionHandlers;function r(r,n,i,o){e.defaults[r]=n,i&&(t[r]=o?function(e,t,r){r!=Cl&&i(e,t,r)}:i)}e.defineOption=r,e.Init=Cl,r("value","",(function(e,t){return e.setValue(t)}),!0),r("mode",null,(function(e,t){e.doc.modeOption=t,Ai(e)}),!0),r("indentUnit",2,Ai,!0),r("indentWithTabs",!1),r("smartIndent",!0),r("tabSize",4,(function(e){Di(e),zr(e),hn(e)}),!0),r("lineSeparator",null,(function(e,t){if(e.doc.lineSep=t,t){var r=[],n=e.doc.first;e.doc.iter((function(e){for(var i=0;;){var o=e.text.indexOf(t,i);if(-1==o)break;i=o+t.length,r.push(et(n,o))}n++}));for(var i=r.length-1;i>=0;i--)go(e.doc,t,r[i],et(r[i].line,r[i].ch+t.length))}})),r("specialChars",/[\u0000-\u001f\u007f-\u009f\u00ad\u061c\u200b\u200e\u200f\u2028\u2029\ufeff\ufff9-\ufffc]/g,(function(e,t,r){e.state.specialChars=new RegExp(t.source+(t.test("\t")?"":"|\t"),"g"),r!=Cl&&e.refresh()})),r("specialCharPlaceholder",Qt,(function(e){return e.refresh()}),!0),r("electricChars",!0),r("inputStyle",m?"contenteditable":"textarea",(function(){throw new Error("inputStyle can not (yet) be changed in a running editor")}),!0),r("spellcheck",!1,(function(e,t){return e.getInputField().spellcheck=t}),!0),r("autocorrect",!1,(function(e,t){return e.getInputField().autocorrect=t}),!0),r("autocapitalize",!1,(function(e,t){return e.getInputField().autocapitalize=t}),!0),r("rtlMoveVisually",!w),r("wholeLineUpdateBefore",!0),r("theme","default",(function(e){xl(e),gi(e)}),!0),r("keyMap","default",(function(e,t,r){var n=qo(t),i=r!=Cl&&qo(r);i&&i.detach&&i.detach(e,n),n.attach&&n.attach(e,i||null)})),r("extraKeys",null),r("configureMouse",null),r("lineWrapping",!1,Tl,!0),r("gutters",[],(function(e,t){e.display.gutterSpecs=di(t,e.options.lineNumbers),gi(e)}),!0),r("fixedGutter",!0,(function(e,t){e.display.gutters.style.left=t?ln(e.display)+"px":"0",e.refresh()}),!0),r("coverGutterNextToScrollbar",!1,(function(e){return Gn(e)}),!0),r("scrollbarStyle","native",(function(e){Kn(e),Gn(e),e.display.scrollbars.setScrollTop(e.doc.scrollTop),e.display.scrollbars.setScrollLeft(e.doc.scrollLeft)}),!0),r("lineNumbers",!1,(function(e,t){e.display.gutterSpecs=di(e.options.gutters,t),gi(e)}),!0),r("firstLineNumber",1,gi,!0),r("lineNumberFormatter",(function(e){return e}),gi,!0),r("showCursorWhenSelecting",!1,vn,!0),r("resetSelectionOnContextMenu",!0),r("lineWiseCopyCut",!0),r("pasteLinesPerSelection",!0),r("selectionsMayTouch",!1),r("readOnly",!1,(function(e,t){"nocursor"==t&&(kn(e),e.display.input.blur()),e.display.input.readOnlyChanged(t)})),r("screenReaderLabel",null,(function(e,t){t=""===t?null:t,e.display.input.screenReaderLabelChanged(t)})),r("disableInput",!1,(function(e,t){t||e.display.input.reset()}),!0),r("dragDrop",!0,kl),r("allowDropFileTypes",null),r("cursorBlinkRate",530),r("cursorScrollMargin",0),r("cursorHeight",1,vn,!0),r("singleCursorHeightPerLine",!0,vn,!0),r("workTime",100),r("workDelay",100),r("flattenSpans",!0,Di,!0),r("addModeClass",!1,Di,!0),r("pollInterval",100),r("undoDepth",200,(function(e,t){return e.doc.history.undoDepth=t})),r("historyEventDelay",1250),r("viewportMargin",10,(function(e){return e.refresh()}),!0),r("maxHighlightLength",1e4,Di,!0),r("moveInputWithCursor",!0,(function(e,t){t||e.display.input.resetPosition()})),r("tabindex",null,(function(e,t){return e.display.input.getField().tabIndex=t||""})),r("autofocus",null),r("direction","ltr",(function(e,t){return e.doc.setDirection(t)}),!0),r("phrases",null)}(Ml),function(e){var t=e.optionHandlers,r=e.helpers={};e.prototype={constructor:e,focus:function(){window.focus(),this.display.input.focus()},setOption:function(e,r){var n=this.options,i=n[e];n[e]==r&&"mode"!=e||(n[e]=r,t.hasOwnProperty(e)&&ei(this,t[e])(this,r,i),pe(this,"optionChange",this,e))},getOption:function(e){return this.options[e]},getDoc:function(){return this.doc},addKeyMap:function(e,t){this.state.keyMaps[t?"push":"unshift"](qo(e))},removeKeyMap:function(e){for(var t=this.state.keyMaps,r=0;r<t.length;++r)if(t[r]==e||t[r].name==e)return t.splice(r,1),!0},addOverlay:ti((function(t,r){var n=t.token?t:e.getMode(this.options,t);if(n.startState)throw new Error("Overlays may not be stateful.");!function(e,t,r){for(var n=0,i=r(t);n<e.length&&r(e[n])<=i;)n++;e.splice(n,0,t)}(this.state.overlays,{mode:n,modeSpec:t,opaque:r&&r.opaque,priority:r&&r.priority||0},(function(e){return e.priority})),this.state.modeGen++,hn(this)})),removeOverlay:ti((function(e){for(var t=this.state.overlays,r=0;r<t.length;++r){var n=t[r].modeSpec;if(n==e||"string"==typeof e&&n.name==e)return t.splice(r,1),this.state.modeGen++,void hn(this)}})),indentLine:ti((function(e,t,r){"string"!=typeof t&&"number"!=typeof t&&(t=null==t?this.options.smartIndent?"smart":"prev":t?"add":"subtract"),Qe(this.doc,e)&&Ol(this,e,t,r)})),indentSelection:ti((function(e){for(var t=this.doc.sel.ranges,r=-1,n=0;n<t.length;n++){var i=t[n];if(i.empty())i.head.line>r&&(Ol(this,i.head.line,e,!0),r=i.head.line,n==this.doc.sel.primIndex&&Dn(this));else{var o=i.from(),l=i.to(),s=Math.max(r,o.line);r=Math.min(this.lastLine(),l.line-(l.ch?0:1))+1;for(var a=s;a<r;++a)Ol(this,a,e);var u=this.doc.sel.ranges;0==o.ch&&t.length==u.length&&u[n].from().ch>0&&qi(this.doc,n,new Si(o,u[n].to()),U)}}})),getTokenAt:function(e,t){return yt(this,e,t)},getLineTokens:function(e,t){return yt(this,et(e),t,!0)},getTokenTypeAt:function(e){e=st(this.doc,e);var t,r=ft(this,Xe(this.doc,e.line)),n=0,i=(r.length-1)/2,o=e.ch;if(0==o)t=r[2];else for(;;){var l=n+i>>1;if((l?r[2*l-1]:0)>=o)i=l;else{if(!(r[2*l+1]<o)){t=r[2*l+2];break}n=l+1}}var s=t?t.indexOf("overlay "):-1;return s<0?t:0==s?null:t.slice(0,s-1)},getModeAt:function(t){var r=this.doc.mode;return r.innerMode?e.innerMode(r,this.getTokenAt(t).state).mode:r},getHelper:function(e,t){return this.getHelpers(e,t)[0]},getHelpers:function(e,t){var n=[];if(!r.hasOwnProperty(t))return n;var i=r[t],o=this.getModeAt(e);if("string"==typeof o[t])i[o[t]]&&n.push(i[o[t]]);else if(o[t])for(var l=0;l<o[t].length;l++){var s=i[o[t][l]];s&&n.push(s)}else o.helperType&&i[o.helperType]?n.push(i[o.helperType]):i[o.name]&&n.push(i[o.name]);for(var a=0;a<i._global.length;a++){var u=i._global[a];u.pred(o,this)&&-1==B(n,u.val)&&n.push(u.val)}return n},getStateAfter:function(e,t){var r=this.doc;return dt(this,(e=lt(r,null==e?r.first+r.size-1:e))+1,t).state},cursorCoords:function(e,t){var r=this.doc.sel.primary();return Xr(this,null==e?r.head:"object"==typeof e?st(this.doc,e):e?r.from():r.to(),t||"page")},charCoords:function(e,t){return jr(this,st(this.doc,e),t||"page")},coordsChar:function(e,t){return _r(this,(e=Kr(this,e,t||"page")).left,e.top)},lineAtHeight:function(e,t){return e=Kr(this,{top:e,left:0},t||"page").top,Ze(this.doc,e+this.display.viewOffset)},heightAtLine:function(e,t,r){var n,i=!1;if("number"==typeof e){var o=this.doc.first+this.doc.size-1;e<this.doc.first?e=this.doc.first:e>o&&(e=o,i=!0),n=Xe(this.doc,e)}else n=e;return Vr(this,n,{top:0,left:0},t||"page",r||i).top+(i?this.doc.height-Vt(n):0)},defaultTextHeight:function(){return rn(this.display)},defaultCharWidth:function(){return nn(this.display)},getViewport:function(){return{from:this.display.viewFrom,to:this.display.viewTo}},addWidget:function(e,t,r,n,i){var o,l,s,a=this.display,u=(e=Xr(this,st(this.doc,e))).bottom,c=e.left;if(t.style.position="absolute",t.setAttribute("cm-ignore-events","true"),this.display.input.setUneditable(t),a.sizer.appendChild(t),"over"==n)u=e.top;else if("above"==n||"near"==n){var h=Math.max(a.wrapper.clientHeight,this.doc.height),f=Math.max(a.sizer.clientWidth,a.lineSpace.clientWidth);("above"==n||e.bottom+t.offsetHeight>h)&&e.top>t.offsetHeight?u=e.top-t.offsetHeight:e.bottom+t.offsetHeight<=h&&(u=e.bottom),c+t.offsetWidth>f&&(c=f-t.offsetWidth)}t.style.top=u+"px",t.style.left=t.style.right="","right"==i?(c=a.sizer.clientWidth-t.offsetWidth,t.style.right="0px"):("left"==i?c=0:"middle"==i&&(c=(a.sizer.clientWidth-t.offsetWidth)/2),t.style.left=c+"px"),r&&(o=this,l={left:c,top:u,right:c+t.offsetWidth,bottom:u+t.offsetHeight},null!=(s=On(o,l)).scrollTop&&Pn(o,s.scrollTop),null!=s.scrollLeft&&In(o,s.scrollLeft))},triggerOnKeyDown:ti(cl),triggerOnKeyPress:ti(fl),triggerOnKeyUp:hl,triggerOnMouseDown:ti(vl),execCommand:function(e){if(tl.hasOwnProperty(e))return tl[e].call(null,this)},triggerElectric:ti((function(e){Fl(this,e)})),findPosH:function(e,t,r,n){var i=1;t<0&&(i=-1,t=-t);for(var o=st(this.doc,e),l=0;l<t&&!(o=Rl(this.doc,o,i,r,n)).hitSide;++l);return o},moveH:ti((function(e,t){var r=this;this.extendSelectionsBy((function(n){return r.display.shift||r.doc.extend||n.empty()?Rl(r.doc,n.head,e,t,r.options.rtlMoveVisually):e<0?n.from():n.to()}),K)})),deleteH:ti((function(e,t){var r=this.doc.sel,n=this.doc;r.somethingSelected()?n.replaceSelection("",null,"+delete"):Zo(this,(function(r){var i=Rl(n,r.head,e,t,!1);return e<0?{from:i,to:r.head}:{from:r.head,to:i}}))})),findPosV:function(e,t,r,n){var i=1,o=n;t<0&&(i=-1,t=-t);for(var l=st(this.doc,e),s=0;s<t;++s){var a=Xr(this,l,"div");if(null==o?o=a.left:a.left=o,(l=zl(this,a,i,r)).hitSide)break}return l},moveV:ti((function(e,t){var r=this,n=this.doc,i=[],o=!this.display.shift&&!n.extend&&n.sel.somethingSelected();if(n.extendSelectionsBy((function(l){if(o)return e<0?l.from():l.to();var s=Xr(r,l.head,"div");null!=l.goalColumn&&(s.left=l.goalColumn),i.push(s.left);var a=zl(r,s,e,t);return"page"==t&&l==n.sel.primary()&&An(r,jr(r,a,"div").top-s.top),a}),K),i.length)for(var l=0;l<n.sel.ranges.length;l++)n.sel.ranges[l].goalColumn=i[l]})),findWordAt:function(e){var t=Xe(this.doc,e.line).text,r=e.ch,n=e.ch;if(t){var i=this.getHelper(e,"wordChars");"before"!=e.sticky&&n!=t.length||!r?++n:--r;for(var o=t.charAt(r),l=ee(o,i)?function(e){return ee(e,i)}:/\s/.test(o)?function(e){return/\s/.test(e)}:function(e){return!/\s/.test(e)&&!ee(e)};r>0&&l(t.charAt(r-1));)--r;for(;n<t.length&&l(t.charAt(n));)++n}return new Si(et(e.line,r),et(e.line,n))},toggleOverwrite:function(e){null!=e&&e==this.state.overwrite||((this.state.overwrite=!this.state.overwrite)?H(this.display.cursorDiv,"CodeMirror-overwrite"):T(this.display.cursorDiv,"CodeMirror-overwrite"),pe(this,"overwriteToggle",this,this.state.overwrite))},hasFocus:function(){return this.display.input.getField()==W()},isReadOnly:function(){return!(!this.options.readOnly&&!this.doc.cantEdit)},scrollTo:ti((function(e,t){Wn(this,e,t)})),getScrollInfo:function(){var e=this.display.scroller;return{left:e.scrollLeft,top:e.scrollTop,height:e.scrollHeight-kr(this)-this.display.barHeight,width:e.scrollWidth-kr(this)-this.display.barWidth,clientHeight:Mr(this),clientWidth:Tr(this)}},scrollIntoView:ti((function(e,t){null==e?(e={from:this.doc.sel.primary().head,to:null},null==t&&(t=this.options.cursorScrollMargin)):"number"==typeof e?e={from:et(e,0),to:null}:null==e.from&&(e={from:e,to:null}),e.to||(e.to=e.from),e.margin=t||0,null!=e.from.line?function(e,t){Hn(e),e.curOp.scrollToPos=t}(this,e):Fn(this,e.from,e.to,e.margin)})),setSize:ti((function(e,t){var r=this,n=function(e){return"number"==typeof e||/^\d+$/.test(String(e))?e+"px":e};null!=e&&(this.display.wrapper.style.width=n(e)),null!=t&&(this.display.wrapper.style.height=n(t)),this.options.lineWrapping&&Rr(this);var i=this.display.viewFrom;this.doc.iter(i,this.display.viewTo,(function(e){if(e.widgets)for(var t=0;t<e.widgets.length;t++)if(e.widgets[t].noHScroll){fn(r,i,"widget");break}++i})),this.curOp.forceUpdate=!0,pe(this,"refresh",this)})),operation:function(e){return Jn(this,e)},startOperation:function(){return Xn(this)},endOperation:function(){return Yn(this)},refresh:ti((function(){var e=this.display.cachedTextHeight;hn(this),this.curOp.forceUpdate=!0,zr(this),Wn(this,this.doc.scrollLeft,this.doc.scrollTop),ui(this.display),(null==e||Math.abs(e-rn(this.display))>.5||this.options.lineWrapping)&&an(this),pe(this,"refresh",this)})),swapDoc:ti((function(e){var t=this.doc;return t.cm=null,this.state.selectingText&&this.state.selectingText(),Pi(this,e),zr(this),this.display.input.reset(),Wn(this,e.scrollLeft,e.scrollTop),this.curOp.forceScroll=!0,sr(this,"swapDoc",this,t),t})),phrase:function(e){var t=this.options.phrases;return t&&Object.prototype.hasOwnProperty.call(t,e)?t[e]:e},getInputField:function(){return this.display.input.getField()},getWrapperElement:function(){return this.display.wrapper},getScrollerElement:function(){return this.display.scroller},getGutterElement:function(){return this.display.gutters}},ye(e),e.registerHelper=function(t,n,i){r.hasOwnProperty(t)||(r[t]=e[t]={_global:[]}),r[t][n]=i},e.registerGlobalHelper=function(t,n,i,o){e.registerHelper(t,n,o),r[t]._global.push({pred:i,val:o})}}(Ml);var Xl="iter insert remove copy getEditor constructor".split(" ");for(var Yl in Do.prototype)Do.prototype.hasOwnProperty(Yl)&&B(Xl,Yl)<0&&(Ml.prototype[Yl]=function(e){return function(){return e.apply(this.doc,arguments)}}(Do.prototype[Yl]));return ye(Do),Ml.inputStyles={textarea:jl,contenteditable:Bl},Ml.defineMode=function(e){Ml.defaults.mode||"null"==e||(Ml.defaults.mode=e),Ie.apply(this,arguments)},Ml.defineMIME=function(e,t){Ee[e]=t},Ml.defineMode("null",(function(){return{token:function(e){return e.skipToEnd()}}})),Ml.defineMIME("text/plain","null"),Ml.defineExtension=function(e,t){Ml.prototype[e]=t},Ml.defineDocExtension=function(e,t){Do.prototype[e]=t},Ml.fromTextArea=function(e,t){if((t=t?I(t):{}).value=e.value,!t.tabindex&&e.tabIndex&&(t.tabindex=e.tabIndex),!t.placeholder&&e.placeholder&&(t.placeholder=e.placeholder),null==t.autofocus){var r=W();t.autofocus=r==e||null!=e.getAttribute("autofocus")&&r==document.body}function n(){e.value=s.getValue()}var i;if(e.form&&(he(e.form,"submit",n),!t.leaveSubmitMethodAlone)){var o=e.form;i=o.submit;try{var l=o.submit=function(){n(),o.submit=i,o.submit(),o.submit=l}}catch(e){}}t.finishInit=function(r){r.save=n,r.getTextArea=function(){return e},r.toTextArea=function(){r.toTextArea=isNaN,n(),e.parentNode.removeChild(r.getWrapperElement()),e.style.display="",e.form&&(de(e.form,"submit",n),t.leaveSubmitMethodAlone||"function"!=typeof e.form.submit||(e.form.submit=i))}},e.style.display="none";var s=Ml((function(t){return e.parentNode.insertBefore(t,e.nextSibling)}),t);return s},function(e){e.off=de,e.on=he,e.wheelEventPixels=wi,e.Doc=Do,e.splitLines=De,e.countColumn=R,e.findColumn=j,e.isWordChar=J,e.Pass=G,e.signal=pe,e.Line=Xt,e.changeEnd=Ti,e.scrollbarModel=Vn,e.Pos=et,e.cmpPos=tt,e.modes=Pe,e.mimeModes=Ee,e.resolveMode=Re,e.getMode=ze,e.modeExtensions=Be,e.extendMode=Ge,e.copyState=Ue,e.startState=Ke,e.innerMode=Ve,e.commands=tl,e.keyMap=Vo,e.keyName=_o,e.isModifierKey=Yo,e.lookupKey=Xo,e.normalizeKeyMap=jo,e.StringStream=je,e.SharedTextMarker=Mo,e.TextMarker=ko,e.LineWidget=Co,e.e_preventDefault=be,e.e_stopPropagation=we,e.e_stop=Ce,e.addClass=H,e.contains=D,e.rmClass=T,e.keyNames=zo}(Ml),Ml.version="5.65.4",Ml}));
//# sourceMappingURL=/sm/c52b3363fb1c07885cf7348cebdd42852a891fe1be5155fb96ce5befc7d5a207.map// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: https://codemirror.net/LICENSE

// Utility function that allows modes to be combined. The mode given
// as the base argument takes care of most of the normal mode
// functionality, but a second (typically simple) mode is used, which
// can override the style of text. Both modes get to parse all of the
// text, but when both assign a non-null style to a piece of code, the
// overlay wins, unless the combine argument was true and not overridden,
// or state.overlay.combineTokens was true, in which case the styles are
// combined.

(function(mod) {
  if (typeof exports == "object" && typeof module == "object") // CommonJS
    mod(require("../../lib/codemirror"));
  else if (typeof define == "function" && define.amd) // AMD
    define(["../../lib/codemirror"], mod);
  else // Plain browser env
    mod(CodeMirror);
})(function(CodeMirror) {
"use strict";

CodeMirror.overlayMode = function(base, overlay, combine) {
  return {
    startState: function() {
      return {
        base: CodeMirror.startState(base),
        overlay: CodeMirror.startState(overlay),
        basePos: 0, baseCur: null,
        overlayPos: 0, overlayCur: null,
        streamSeen: null
      };
    },
    copyState: function(state) {
      return {
        base: CodeMirror.copyState(base, state.base),
        overlay: CodeMirror.copyState(overlay, state.overlay),
        basePos: state.basePos, baseCur: null,
        overlayPos: state.overlayPos, overlayCur: null
      };
    },

    token: function(stream, state) {
      if (stream != state.streamSeen ||
          Math.min(state.basePos, state.overlayPos) < stream.start) {
        state.streamSeen = stream;
        state.basePos = state.overlayPos = stream.start;
      }

      if (stream.start == state.basePos) {
        state.baseCur = base.token(stream, state.base);
        state.basePos = stream.pos;
      }
      if (stream.start == state.overlayPos) {
        stream.pos = stream.start;
        state.overlayCur = overlay.token(stream, state.overlay);
        state.overlayPos = stream.pos;
      }
      stream.pos = Math.min(state.basePos, state.overlayPos);

      // state.overlay.combineTokens always takes precedence over combine,
      // unless set to null
      if (state.overlayCur == null) return state.baseCur;
      else if (state.baseCur != null &&
               state.overlay.combineTokens ||
               combine && state.overlay.combineTokens == null)
        return state.baseCur + " " + state.overlayCur;
      else return state.overlayCur;
    },

    indent: base.indent && function(state, textAfter, line) {
      return base.indent(state.base, textAfter, line);
    },
    electricChars: base.electricChars,

    innerMode: function(state) { return {state: state.base, mode: base}; },

    blankLine: function(state) {
      var baseToken, overlayToken;
      if (base.blankLine) baseToken = base.blankLine(state.base);
      if (overlay.blankLine) overlayToken = overlay.blankLine(state.overlay);

      return overlayToken == null ?
        baseToken :
        (combine && baseToken != null ? baseToken + " " + overlayToken : overlayToken);
    }
  };
};

});
// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: https://codemirror.net/LICENSE

(function(mod) {
  if (typeof exports == "object" && typeof module == "object") // CommonJS
    mod(require("../../lib/codemirror"), require("../xml/xml"), require("../meta"));
  else if (typeof define == "function" && define.amd) // AMD
    define(["../../lib/codemirror", "../xml/xml", "../meta"], mod);
  else // Plain browser env
    mod(CodeMirror);
})(function(CodeMirror) {
"use strict";

CodeMirror.defineMode("markdown", function(cmCfg, modeCfg) {

  var htmlMode = CodeMirror.getMode(cmCfg, "text/html");
  var htmlModeMissing = htmlMode.name == "null"

  function getMode(name) {
    if (CodeMirror.findModeByName) {
      var found = CodeMirror.findModeByName(name);
      if (found) name = found.mime || found.mimes[0];
    }
    var mode = CodeMirror.getMode(cmCfg, name);
    return mode.name == "null" ? null : mode;
  }

  // Should characters that affect highlighting be highlighted separate?
  // Does not include characters that will be output (such as `1.` and `-` for lists)
  if (modeCfg.highlightFormatting === undefined)
    modeCfg.highlightFormatting = false;

  // Maximum number of nested blockquotes. Set to 0 for infinite nesting.
  // Excess `>` will emit `error` token.
  if (modeCfg.maxBlockquoteDepth === undefined)
    modeCfg.maxBlockquoteDepth = 0;

  // Turn on task lists? ("- [ ] " and "- [x] ")
  if (modeCfg.taskLists === undefined) modeCfg.taskLists = false;

  // Turn on strikethrough syntax
  if (modeCfg.strikethrough === undefined)
    modeCfg.strikethrough = false;

  if (modeCfg.emoji === undefined)
    modeCfg.emoji = false;

  if (modeCfg.fencedCodeBlockHighlighting === undefined)
    modeCfg.fencedCodeBlockHighlighting = true;

  if (modeCfg.fencedCodeBlockDefaultMode === undefined)
    modeCfg.fencedCodeBlockDefaultMode = 'text/plain';

  if (modeCfg.xml === undefined)
    modeCfg.xml = true;

  // Allow token types to be overridden by user-provided token types.
  if (modeCfg.tokenTypeOverrides === undefined)
    modeCfg.tokenTypeOverrides = {};

  var tokenTypes = {
    header: "header",
    code: "comment",
    quote: "quote",
    list1: "variable-2",
    list2: "variable-3",
    list3: "keyword",
    hr: "hr",
    image: "image",
    imageAltText: "image-alt-text",
    imageMarker: "image-marker",
    formatting: "formatting",
    linkInline: "link",
    linkEmail: "link",
    linkText: "link",
    linkHref: "string",
    em: "em",
    strong: "strong",
    strikethrough: "strikethrough",
    emoji: "builtin"
  };

  for (var tokenType in tokenTypes) {
    if (tokenTypes.hasOwnProperty(tokenType) && modeCfg.tokenTypeOverrides[tokenType]) {
      tokenTypes[tokenType] = modeCfg.tokenTypeOverrides[tokenType];
    }
  }

  var hrRE = /^([*\-_])(?:\s*\1){2,}\s*$/
  ,   listRE = /^(?:[*\-+]|^[0-9]+([.)]))\s+/
  ,   taskListRE = /^\[(x| )\](?=\s)/i // Must follow listRE
  ,   atxHeaderRE = modeCfg.allowAtxHeaderWithoutSpace ? /^(#+)/ : /^(#+)(?: |$)/
  ,   setextHeaderRE = /^ {0,3}(?:\={1,}|-{2,})\s*$/
  ,   textRE = /^[^#!\[\]*_\\<>` "'(~:]+/
  ,   fencedCodeRE = /^(~~~+|```+)[ \t]*([\w\/+#-]*)[^\n`]*$/
  ,   linkDefRE = /^\s*\[[^\]]+?\]:.*$/ // naive link-definition
  ,   punctuation = /[!"#$%&'()*+,\-.\/:;<=>?@\[\\\]^_`{|}~\xA1\xA7\xAB\xB6\xB7\xBB\xBF\u037E\u0387\u055A-\u055F\u0589\u058A\u05BE\u05C0\u05C3\u05C6\u05F3\u05F4\u0609\u060A\u060C\u060D\u061B\u061E\u061F\u066A-\u066D\u06D4\u0700-\u070D\u07F7-\u07F9\u0830-\u083E\u085E\u0964\u0965\u0970\u0AF0\u0DF4\u0E4F\u0E5A\u0E5B\u0F04-\u0F12\u0F14\u0F3A-\u0F3D\u0F85\u0FD0-\u0FD4\u0FD9\u0FDA\u104A-\u104F\u10FB\u1360-\u1368\u1400\u166D\u166E\u169B\u169C\u16EB-\u16ED\u1735\u1736\u17D4-\u17D6\u17D8-\u17DA\u1800-\u180A\u1944\u1945\u1A1E\u1A1F\u1AA0-\u1AA6\u1AA8-\u1AAD\u1B5A-\u1B60\u1BFC-\u1BFF\u1C3B-\u1C3F\u1C7E\u1C7F\u1CC0-\u1CC7\u1CD3\u2010-\u2027\u2030-\u2043\u2045-\u2051\u2053-\u205E\u207D\u207E\u208D\u208E\u2308-\u230B\u2329\u232A\u2768-\u2775\u27C5\u27C6\u27E6-\u27EF\u2983-\u2998\u29D8-\u29DB\u29FC\u29FD\u2CF9-\u2CFC\u2CFE\u2CFF\u2D70\u2E00-\u2E2E\u2E30-\u2E42\u3001-\u3003\u3008-\u3011\u3014-\u301F\u3030\u303D\u30A0\u30FB\uA4FE\uA4FF\uA60D-\uA60F\uA673\uA67E\uA6F2-\uA6F7\uA874-\uA877\uA8CE\uA8CF\uA8F8-\uA8FA\uA8FC\uA92E\uA92F\uA95F\uA9C1-\uA9CD\uA9DE\uA9DF\uAA5C-\uAA5F\uAADE\uAADF\uAAF0\uAAF1\uABEB\uFD3E\uFD3F\uFE10-\uFE19\uFE30-\uFE52\uFE54-\uFE61\uFE63\uFE68\uFE6A\uFE6B\uFF01-\uFF03\uFF05-\uFF0A\uFF0C-\uFF0F\uFF1A\uFF1B\uFF1F\uFF20\uFF3B-\uFF3D\uFF3F\uFF5B\uFF5D\uFF5F-\uFF65]|\uD800[\uDD00-\uDD02\uDF9F\uDFD0]|\uD801\uDD6F|\uD802[\uDC57\uDD1F\uDD3F\uDE50-\uDE58\uDE7F\uDEF0-\uDEF6\uDF39-\uDF3F\uDF99-\uDF9C]|\uD804[\uDC47-\uDC4D\uDCBB\uDCBC\uDCBE-\uDCC1\uDD40-\uDD43\uDD74\uDD75\uDDC5-\uDDC9\uDDCD\uDDDB\uDDDD-\uDDDF\uDE38-\uDE3D\uDEA9]|\uD805[\uDCC6\uDDC1-\uDDD7\uDE41-\uDE43\uDF3C-\uDF3E]|\uD809[\uDC70-\uDC74]|\uD81A[\uDE6E\uDE6F\uDEF5\uDF37-\uDF3B\uDF44]|\uD82F\uDC9F|\uD836[\uDE87-\uDE8B]/
  ,   expandedTab = "    " // CommonMark specifies tab as 4 spaces

  function switchInline(stream, state, f) {
    state.f = state.inline = f;
    return f(stream, state);
  }

  function switchBlock(stream, state, f) {
    state.f = state.block = f;
    return f(stream, state);
  }

  function lineIsEmpty(line) {
    return !line || !/\S/.test(line.string)
  }

  // Blocks

  function blankLine(state) {
    // Reset linkTitle state
    state.linkTitle = false;
    state.linkHref = false;
    state.linkText = false;
    // Reset EM state
    state.em = false;
    // Reset STRONG state
    state.strong = false;
    // Reset strikethrough state
    state.strikethrough = false;
    // Reset state.quote
    state.quote = 0;
    // Reset state.indentedCode
    state.indentedCode = false;
    if (state.f == htmlBlock) {
      var exit = htmlModeMissing
      if (!exit) {
        var inner = CodeMirror.innerMode(htmlMode, state.htmlState)
        exit = inner.mode.name == "xml" && inner.state.tagStart === null &&
          (!inner.state.context && inner.state.tokenize.isInText)
      }
      if (exit) {
        state.f = inlineNormal;
        state.block = blockNormal;
        state.htmlState = null;
      }
    }
    // Reset state.trailingSpace
    state.trailingSpace = 0;
    state.trailingSpaceNewLine = false;
    // Mark this line as blank
    state.prevLine = state.thisLine
    state.thisLine = {stream: null}
    return null;
  }

  function blockNormal(stream, state) {
    var firstTokenOnLine = stream.column() === state.indentation;
    var prevLineLineIsEmpty = lineIsEmpty(state.prevLine.stream);
    var prevLineIsIndentedCode = state.indentedCode;
    var prevLineIsHr = state.prevLine.hr;
    var prevLineIsList = state.list !== false;
    var maxNonCodeIndentation = (state.listStack[state.listStack.length - 1] || 0) + 3;

    state.indentedCode = false;

    var lineIndentation = state.indentation;
    // compute once per line (on first token)
    if (state.indentationDiff === null) {
      state.indentationDiff = state.indentation;
      if (prevLineIsList) {
        state.list = null;
        // While this list item's marker's indentation is less than the deepest
        //  list item's content's indentation,pop the deepest list item
        //  indentation off the stack, and update block indentation state
        while (lineIndentation < state.listStack[state.listStack.length - 1]) {
          state.listStack.pop();
          if (state.listStack.length) {
            state.indentation = state.listStack[state.listStack.length - 1];
          // less than the first list's indent -> the line is no longer a list
          } else {
            state.list = false;
          }
        }
        if (state.list !== false) {
          state.indentationDiff = lineIndentation - state.listStack[state.listStack.length - 1]
        }
      }
    }

    // not comprehensive (currently only for setext detection purposes)
    var allowsInlineContinuation = (
        !prevLineLineIsEmpty && !prevLineIsHr && !state.prevLine.header &&
        (!prevLineIsList || !prevLineIsIndentedCode) &&
        !state.prevLine.fencedCodeEnd
    );

    var isHr = (state.list === false || prevLineIsHr || prevLineLineIsEmpty) &&
      state.indentation <= maxNonCodeIndentation && stream.match(hrRE);

    var match = null;
    if (state.indentationDiff >= 4 && (prevLineIsIndentedCode || state.prevLine.fencedCodeEnd ||
         state.prevLine.header || prevLineLineIsEmpty)) {
      stream.skipToEnd();
      state.indentedCode = true;
      return tokenTypes.code;
    } else if (stream.eatSpace()) {
      return null;
    } else if (firstTokenOnLine && state.indentation <= maxNonCodeIndentation && (match = stream.match(atxHeaderRE)) && match[1].length <= 6) {
      state.quote = 0;
      state.header = match[1].length;
      state.thisLine.header = true;
      if (modeCfg.highlightFormatting) state.formatting = "header";
      state.f = state.inline;
      return getType(state);
    } else if (state.indentation <= maxNonCodeIndentation && stream.eat('>')) {
      state.quote = firstTokenOnLine ? 1 : state.quote + 1;
      if (modeCfg.highlightFormatting) state.formatting = "quote";
      stream.eatSpace();
      return getType(state);
    } else if (!isHr && !state.setext && firstTokenOnLine && state.indentation <= maxNonCodeIndentation && (match = stream.match(listRE))) {
      var listType = match[1] ? "ol" : "ul";

      state.indentation = lineIndentation + stream.current().length;
      state.list = true;
      state.quote = 0;

      // Add this list item's content's indentation to the stack
      state.listStack.push(state.indentation);
      // Reset inline styles which shouldn't propagate across list items
      state.em = false;
      state.strong = false;
      state.code = false;
      state.strikethrough = false;

      if (modeCfg.taskLists && stream.match(taskListRE, false)) {
        state.taskList = true;
      }
      state.f = state.inline;
      if (modeCfg.highlightFormatting) state.formatting = ["list", "list-" + listType];
      return getType(state);
    } else if (firstTokenOnLine && state.indentation <= maxNonCodeIndentation && (match = stream.match(fencedCodeRE, true))) {
      state.quote = 0;
      state.fencedEndRE = new RegExp(match[1] + "+ *$");
      // try switching mode
      state.localMode = modeCfg.fencedCodeBlockHighlighting && getMode(match[2] || modeCfg.fencedCodeBlockDefaultMode );
      if (state.localMode) state.localState = CodeMirror.startState(state.localMode);
      state.f = state.block = local;
      if (modeCfg.highlightFormatting) state.formatting = "code-block";
      state.code = -1
      return getType(state);
    // SETEXT has lowest block-scope precedence after HR, so check it after
    //  the others (code, blockquote, list...)
    } else if (
      // if setext set, indicates line after ---/===
      state.setext || (
        // line before ---/===
        (!allowsInlineContinuation || !prevLineIsList) && !state.quote && state.list === false &&
        !state.code && !isHr && !linkDefRE.test(stream.string) &&
        (match = stream.lookAhead(1)) && (match = match.match(setextHeaderRE))
      )
    ) {
      if ( !state.setext ) {
        state.header = match[0].charAt(0) == '=' ? 1 : 2;
        state.setext = state.header;
      } else {
        state.header = state.setext;
        // has no effect on type so we can reset it now
        state.setext = 0;
        stream.skipToEnd();
        if (modeCfg.highlightFormatting) state.formatting = "header";
      }
      state.thisLine.header = true;
      state.f = state.inline;
      return getType(state);
    } else if (isHr) {
      stream.skipToEnd();
      state.hr = true;
      state.thisLine.hr = true;
      return tokenTypes.hr;
    } else if (stream.peek() === '[') {
      return switchInline(stream, state, footnoteLink);
    }

    return switchInline(stream, state, state.inline);
  }

  function htmlBlock(stream, state) {
    var style = htmlMode.token(stream, state.htmlState);
    if (!htmlModeMissing) {
      var inner = CodeMirror.innerMode(htmlMode, state.htmlState)
      if ((inner.mode.name == "xml" && inner.state.tagStart === null &&
           (!inner.state.context && inner.state.tokenize.isInText)) ||
          (state.md_inside && stream.current().indexOf(">") > -1)) {
        state.f = inlineNormal;
        state.block = blockNormal;
        state.htmlState = null;
      }
    }
    return style;
  }

  function local(stream, state) {
    var currListInd = state.listStack[state.listStack.length - 1] || 0;
    var hasExitedList = state.indentation < currListInd;
    var maxFencedEndInd = currListInd + 3;
    if (state.fencedEndRE && state.indentation <= maxFencedEndInd && (hasExitedList || stream.match(state.fencedEndRE))) {
      if (modeCfg.highlightFormatting) state.formatting = "code-block";
      var returnType;
      if (!hasExitedList) returnType = getType(state)
      state.localMode = state.localState = null;
      state.block = blockNormal;
      state.f = inlineNormal;
      state.fencedEndRE = null;
      state.code = 0
      state.thisLine.fencedCodeEnd = true;
      if (hasExitedList) return switchBlock(stream, state, state.block);
      return returnType;
    } else if (state.localMode) {
      return state.localMode.token(stream, state.localState);
    } else {
      stream.skipToEnd();
      return tokenTypes.code;
    }
  }

  // Inline
  function getType(state) {
    var styles = [];

    if (state.formatting) {
      styles.push(tokenTypes.formatting);

      if (typeof state.formatting === "string") state.formatting = [state.formatting];

      for (var i = 0; i < state.formatting.length; i++) {
        styles.push(tokenTypes.formatting + "-" + state.formatting[i]);

        if (state.formatting[i] === "header") {
          styles.push(tokenTypes.formatting + "-" + state.formatting[i] + "-" + state.header);
        }

        // Add `formatting-quote` and `formatting-quote-#` for blockquotes
        // Add `error` instead if the maximum blockquote nesting depth is passed
        if (state.formatting[i] === "quote") {
          if (!modeCfg.maxBlockquoteDepth || modeCfg.maxBlockquoteDepth >= state.quote) {
            styles.push(tokenTypes.formatting + "-" + state.formatting[i] + "-" + state.quote);
          } else {
            styles.push("error");
          }
        }
      }
    }

    if (state.taskOpen) {
      styles.push("meta");
      return styles.length ? styles.join(' ') : null;
    }
    if (state.taskClosed) {
      styles.push("property");
      return styles.length ? styles.join(' ') : null;
    }

    if (state.linkHref) {
      styles.push(tokenTypes.linkHref, "url");
    } else { // Only apply inline styles to non-url text
      if (state.strong) { styles.push(tokenTypes.strong); }
      if (state.em) { styles.push(tokenTypes.em); }
      if (state.strikethrough) { styles.push(tokenTypes.strikethrough); }
      if (state.emoji) { styles.push(tokenTypes.emoji); }
      if (state.linkText) { styles.push(tokenTypes.linkText); }
      if (state.code) { styles.push(tokenTypes.code); }
      if (state.image) { styles.push(tokenTypes.image); }
      if (state.imageAltText) { styles.push(tokenTypes.imageAltText, "link"); }
      if (state.imageMarker) { styles.push(tokenTypes.imageMarker); }
    }

    if (state.header) { styles.push(tokenTypes.header, tokenTypes.header + "-" + state.header); }

    if (state.quote) {
      styles.push(tokenTypes.quote);

      // Add `quote-#` where the maximum for `#` is modeCfg.maxBlockquoteDepth
      if (!modeCfg.maxBlockquoteDepth || modeCfg.maxBlockquoteDepth >= state.quote) {
        styles.push(tokenTypes.quote + "-" + state.quote);
      } else {
        styles.push(tokenTypes.quote + "-" + modeCfg.maxBlockquoteDepth);
      }
    }

    if (state.list !== false) {
      var listMod = (state.listStack.length - 1) % 3;
      if (!listMod) {
        styles.push(tokenTypes.list1);
      } else if (listMod === 1) {
        styles.push(tokenTypes.list2);
      } else {
        styles.push(tokenTypes.list3);
      }
    }

    if (state.trailingSpaceNewLine) {
      styles.push("trailing-space-new-line");
    } else if (state.trailingSpace) {
      styles.push("trailing-space-" + (state.trailingSpace % 2 ? "a" : "b"));
    }

    return styles.length ? styles.join(' ') : null;
  }

  function handleText(stream, state) {
    if (stream.match(textRE, true)) {
      return getType(state);
    }
    return undefined;
  }

  function inlineNormal(stream, state) {
    var style = state.text(stream, state);
    if (typeof style !== 'undefined')
      return style;

    if (state.list) { // List marker (*, +, -, 1., etc)
      state.list = null;
      return getType(state);
    }

    if (state.taskList) {
      var taskOpen = stream.match(taskListRE, true)[1] === " ";
      if (taskOpen) state.taskOpen = true;
      else state.taskClosed = true;
      if (modeCfg.highlightFormatting) state.formatting = "task";
      state.taskList = false;
      return getType(state);
    }

    state.taskOpen = false;
    state.taskClosed = false;

    if (state.header && stream.match(/^#+$/, true)) {
      if (modeCfg.highlightFormatting) state.formatting = "header";
      return getType(state);
    }

    var ch = stream.next();

    // Matches link titles present on next line
    if (state.linkTitle) {
      state.linkTitle = false;
      var matchCh = ch;
      if (ch === '(') {
        matchCh = ')';
      }
      matchCh = (matchCh+'').replace(/([.?*+^\[\]\\(){}|-])/g, "\\$1");
      var regex = '^\\s*(?:[^' + matchCh + '\\\\]+|\\\\\\\\|\\\\.)' + matchCh;
      if (stream.match(new RegExp(regex), true)) {
        return tokenTypes.linkHref;
      }
    }

    // If this block is changed, it may need to be updated in GFM mode
    if (ch === '`') {
      var previousFormatting = state.formatting;
      if (modeCfg.highlightFormatting) state.formatting = "code";
      stream.eatWhile('`');
      var count = stream.current().length
      if (state.code == 0 && (!state.quote || count == 1)) {
        state.code = count
        return getType(state)
      } else if (count == state.code) { // Must be exact
        var t = getType(state)
        state.code = 0
        return t
      } else {
        state.formatting = previousFormatting
        return getType(state)
      }
    } else if (state.code) {
      return getType(state);
    }

    if (ch === '\\') {
      stream.next();
      if (modeCfg.highlightFormatting) {
        var type = getType(state);
        var formattingEscape = tokenTypes.formatting + "-escape";
        return type ? type + " " + formattingEscape : formattingEscape;
      }
    }

    if (ch === '!' && stream.match(/\[[^\]]*\] ?(?:\(|\[)/, false)) {
      state.imageMarker = true;
      state.image = true;
      if (modeCfg.highlightFormatting) state.formatting = "image";
      return getType(state);
    }

    if (ch === '[' && state.imageMarker && stream.match(/[^\]]*\](\(.*?\)| ?\[.*?\])/, false)) {
      state.imageMarker = false;
      state.imageAltText = true
      if (modeCfg.highlightFormatting) state.formatting = "image";
      return getType(state);
    }

    if (ch === ']' && state.imageAltText) {
      if (modeCfg.highlightFormatting) state.formatting = "image";
      var type = getType(state);
      state.imageAltText = false;
      state.image = false;
      state.inline = state.f = linkHref;
      return type;
    }

    if (ch === '[' && !state.image) {
      if (state.linkText && stream.match(/^.*?\]/)) return getType(state)
      state.linkText = true;
      if (modeCfg.highlightFormatting) state.formatting = "link";
      return getType(state);
    }

    if (ch === ']' && state.linkText) {
      if (modeCfg.highlightFormatting) state.formatting = "link";
      var type = getType(state);
      state.linkText = false;
      state.inline = state.f = stream.match(/\(.*?\)| ?\[.*?\]/, false) ? linkHref : inlineNormal
      return type;
    }

    if (ch === '<' && stream.match(/^(https?|ftps?):\/\/(?:[^\\>]|\\.)+>/, false)) {
      state.f = state.inline = linkInline;
      if (modeCfg.highlightFormatting) state.formatting = "link";
      var type = getType(state);
      if (type){
        type += " ";
      } else {
        type = "";
      }
      return type + tokenTypes.linkInline;
    }

    if (ch === '<' && stream.match(/^[^> \\]+@(?:[^\\>]|\\.)+>/, false)) {
      state.f = state.inline = linkInline;
      if (modeCfg.highlightFormatting) state.formatting = "link";
      var type = getType(state);
      if (type){
        type += " ";
      } else {
        type = "";
      }
      return type + tokenTypes.linkEmail;
    }

    if (modeCfg.xml && ch === '<' && stream.match(/^(!--|\?|!\[CDATA\[|[a-z][a-z0-9-]*(?:\s+[a-z_:.\-]+(?:\s*=\s*[^>]+)?)*\s*(?:>|$))/i, false)) {
      var end = stream.string.indexOf(">", stream.pos);
      if (end != -1) {
        var atts = stream.string.substring(stream.start, end);
        if (/markdown\s*=\s*('|"){0,1}1('|"){0,1}/.test(atts)) state.md_inside = true;
      }
      stream.backUp(1);
      state.htmlState = CodeMirror.startState(htmlMode);
      return switchBlock(stream, state, htmlBlock);
    }

    if (modeCfg.xml && ch === '<' && stream.match(/^\/\w*?>/)) {
      state.md_inside = false;
      return "tag";
    } else if (ch === "*" || ch === "_") {
      var len = 1, before = stream.pos == 1 ? " " : stream.string.charAt(stream.pos - 2)
      while (len < 3 && stream.eat(ch)) len++
      var after = stream.peek() || " "
      // See http://spec.commonmark.org/0.27/#emphasis-and-strong-emphasis
      var leftFlanking = !/\s/.test(after) && (!punctuation.test(after) || /\s/.test(before) || punctuation.test(before))
      var rightFlanking = !/\s/.test(before) && (!punctuation.test(before) || /\s/.test(after) || punctuation.test(after))
      var setEm = null, setStrong = null
      if (len % 2) { // Em
        if (!state.em && leftFlanking && (ch === "*" || !rightFlanking || punctuation.test(before)))
          setEm = true
        else if (state.em == ch && rightFlanking && (ch === "*" || !leftFlanking || punctuation.test(after)))
          setEm = false
      }
      if (len > 1) { // Strong
        if (!state.strong && leftFlanking && (ch === "*" || !rightFlanking || punctuation.test(before)))
          setStrong = true
        else if (state.strong == ch && rightFlanking && (ch === "*" || !leftFlanking || punctuation.test(after)))
          setStrong = false
      }
      if (setStrong != null || setEm != null) {
        if (modeCfg.highlightFormatting) state.formatting = setEm == null ? "strong" : setStrong == null ? "em" : "strong em"
        if (setEm === true) state.em = ch
        if (setStrong === true) state.strong = ch
        var t = getType(state)
        if (setEm === false) state.em = false
        if (setStrong === false) state.strong = false
        return t
      }
    } else if (ch === ' ') {
      if (stream.eat('*') || stream.eat('_')) { // Probably surrounded by spaces
        if (stream.peek() === ' ') { // Surrounded by spaces, ignore
          return getType(state);
        } else { // Not surrounded by spaces, back up pointer
          stream.backUp(1);
        }
      }
    }

    if (modeCfg.strikethrough) {
      if (ch === '~' && stream.eatWhile(ch)) {
        if (state.strikethrough) {// Remove strikethrough
          if (modeCfg.highlightFormatting) state.formatting = "strikethrough";
          var t = getType(state);
          state.strikethrough = false;
          return t;
        } else if (stream.match(/^[^\s]/, false)) {// Add strikethrough
          state.strikethrough = true;
          if (modeCfg.highlightFormatting) state.formatting = "strikethrough";
          return getType(state);
        }
      } else if (ch === ' ') {
        if (stream.match('~~', true)) { // Probably surrounded by space
          if (stream.peek() === ' ') { // Surrounded by spaces, ignore
            return getType(state);
          } else { // Not surrounded by spaces, back up pointer
            stream.backUp(2);
          }
        }
      }
    }

    if (modeCfg.emoji && ch === ":" && stream.match(/^(?:[a-z_\d+][a-z_\d+-]*|\-[a-z_\d+][a-z_\d+-]*):/)) {
      state.emoji = true;
      if (modeCfg.highlightFormatting) state.formatting = "emoji";
      var retType = getType(state);
      state.emoji = false;
      return retType;
    }

    if (ch === ' ') {
      if (stream.match(/^ +$/, false)) {
        state.trailingSpace++;
      } else if (state.trailingSpace) {
        state.trailingSpaceNewLine = true;
      }
    }

    return getType(state);
  }

  function linkInline(stream, state) {
    var ch = stream.next();

    if (ch === ">") {
      state.f = state.inline = inlineNormal;
      if (modeCfg.highlightFormatting) state.formatting = "link";
      var type = getType(state);
      if (type){
        type += " ";
      } else {
        type = "";
      }
      return type + tokenTypes.linkInline;
    }

    stream.match(/^[^>]+/, true);

    return tokenTypes.linkInline;
  }

  function linkHref(stream, state) {
    // Check if space, and return NULL if so (to avoid marking the space)
    if(stream.eatSpace()){
      return null;
    }
    var ch = stream.next();
    if (ch === '(' || ch === '[') {
      state.f = state.inline = getLinkHrefInside(ch === "(" ? ")" : "]");
      if (modeCfg.highlightFormatting) state.formatting = "link-string";
      state.linkHref = true;
      return getType(state);
    }
    return 'error';
  }

  var linkRE = {
    ")": /^(?:[^\\\(\)]|\\.|\((?:[^\\\(\)]|\\.)*\))*?(?=\))/,
    "]": /^(?:[^\\\[\]]|\\.|\[(?:[^\\\[\]]|\\.)*\])*?(?=\])/
  }

  function getLinkHrefInside(endChar) {
    return function(stream, state) {
      var ch = stream.next();

      if (ch === endChar) {
        state.f = state.inline = inlineNormal;
        if (modeCfg.highlightFormatting) state.formatting = "link-string";
        var returnState = getType(state);
        state.linkHref = false;
        return returnState;
      }

      stream.match(linkRE[endChar])
      state.linkHref = true;
      return getType(state);
    };
  }

  function footnoteLink(stream, state) {
    if (stream.match(/^([^\]\\]|\\.)*\]:/, false)) {
      state.f = footnoteLinkInside;
      stream.next(); // Consume [
      if (modeCfg.highlightFormatting) state.formatting = "link";
      state.linkText = true;
      return getType(state);
    }
    return switchInline(stream, state, inlineNormal);
  }

  function footnoteLinkInside(stream, state) {
    if (stream.match(']:', true)) {
      state.f = state.inline = footnoteUrl;
      if (modeCfg.highlightFormatting) state.formatting = "link";
      var returnType = getType(state);
      state.linkText = false;
      return returnType;
    }

    stream.match(/^([^\]\\]|\\.)+/, true);

    return tokenTypes.linkText;
  }

  function footnoteUrl(stream, state) {
    // Check if space, and return NULL if so (to avoid marking the space)
    if(stream.eatSpace()){
      return null;
    }
    // Match URL
    stream.match(/^[^\s]+/, true);
    // Check for link title
    if (stream.peek() === undefined) { // End of line, set flag to check next line
      state.linkTitle = true;
    } else { // More content on line, check if link title
      stream.match(/^(?:\s+(?:"(?:[^"\\]|\\.)+"|'(?:[^'\\]|\\.)+'|\((?:[^)\\]|\\.)+\)))?/, true);
    }
    state.f = state.inline = inlineNormal;
    return tokenTypes.linkHref + " url";
  }

  var mode = {
    startState: function() {
      return {
        f: blockNormal,

        prevLine: {stream: null},
        thisLine: {stream: null},

        block: blockNormal,
        htmlState: null,
        indentation: 0,

        inline: inlineNormal,
        text: handleText,

        formatting: false,
        linkText: false,
        linkHref: false,
        linkTitle: false,
        code: 0,
        em: false,
        strong: false,
        header: 0,
        setext: 0,
        hr: false,
        taskList: false,
        list: false,
        listStack: [],
        quote: 0,
        trailingSpace: 0,
        trailingSpaceNewLine: false,
        strikethrough: false,
        emoji: false,
        fencedEndRE: null
      };
    },

    copyState: function(s) {
      return {
        f: s.f,

        prevLine: s.prevLine,
        thisLine: s.thisLine,

        block: s.block,
        htmlState: s.htmlState && CodeMirror.copyState(htmlMode, s.htmlState),
        indentation: s.indentation,

        localMode: s.localMode,
        localState: s.localMode ? CodeMirror.copyState(s.localMode, s.localState) : null,

        inline: s.inline,
        text: s.text,
        formatting: false,
        linkText: s.linkText,
        linkTitle: s.linkTitle,
        linkHref: s.linkHref,
        code: s.code,
        em: s.em,
        strong: s.strong,
        strikethrough: s.strikethrough,
        emoji: s.emoji,
        header: s.header,
        setext: s.setext,
        hr: s.hr,
        taskList: s.taskList,
        list: s.list,
        listStack: s.listStack.slice(0),
        quote: s.quote,
        indentedCode: s.indentedCode,
        trailingSpace: s.trailingSpace,
        trailingSpaceNewLine: s.trailingSpaceNewLine,
        md_inside: s.md_inside,
        fencedEndRE: s.fencedEndRE
      };
    },

    token: function(stream, state) {

      // Reset state.formatting
      state.formatting = false;

      if (stream != state.thisLine.stream) {
        state.header = 0;
        state.hr = false;

        if (stream.match(/^\s*$/, true)) {
          blankLine(state);
          return null;
        }

        state.prevLine = state.thisLine
        state.thisLine = {stream: stream}

        // Reset state.taskList
        state.taskList = false;

        // Reset state.trailingSpace
        state.trailingSpace = 0;
        state.trailingSpaceNewLine = false;

        if (!state.localState) {
          state.f = state.block;
          if (state.f != htmlBlock) {
            var indentation = stream.match(/^\s*/, true)[0].replace(/\t/g, expandedTab).length;
            state.indentation = indentation;
            state.indentationDiff = null;
            if (indentation > 0) return null;
          }
        }
      }
      return state.f(stream, state);
    },

    innerMode: function(state) {
      if (state.block == htmlBlock) return {state: state.htmlState, mode: htmlMode};
      if (state.localState) return {state: state.localState, mode: state.localMode};
      return {state: state, mode: mode};
    },

    indent: function(state, textAfter, line) {
      if (state.block == htmlBlock && htmlMode.indent) return htmlMode.indent(state.htmlState, textAfter, line)
      if (state.localState && state.localMode.indent) return state.localMode.indent(state.localState, textAfter, line)
      return CodeMirror.Pass
    },

    blankLine: blankLine,

    getType: getType,

    blockCommentStart: "<!--",
    blockCommentEnd: "-->",
    closeBrackets: "()[]{}''\"\"``",
    fold: "markdown"
  };
  return mode;
}, "xml");

CodeMirror.defineMIME("text/markdown", "markdown");

CodeMirror.defineMIME("text/x-markdown", "markdown");

});
// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: https://codemirror.net/LICENSE

(function(mod) {
  if (typeof exports == "object" && typeof module == "object") // CommonJS
    mod(require("../../lib/codemirror"));
  else if (typeof define == "function" && define.amd) // AMD
    define(["../../lib/codemirror"], mod);
  else // Plain browser env
    mod(CodeMirror);
})(function(CodeMirror) {
"use strict";

var htmlConfig = {
  autoSelfClosers: {'area': true, 'base': true, 'br': true, 'col': true, 'command': true,
                    'embed': true, 'frame': true, 'hr': true, 'img': true, 'input': true,
                    'keygen': true, 'link': true, 'meta': true, 'param': true, 'source': true,
                    'track': true, 'wbr': true, 'menuitem': true},
  implicitlyClosed: {'dd': true, 'li': true, 'optgroup': true, 'option': true, 'p': true,
                     'rp': true, 'rt': true, 'tbody': true, 'td': true, 'tfoot': true,
                     'th': true, 'tr': true},
  contextGrabbers: {
    'dd': {'dd': true, 'dt': true},
    'dt': {'dd': true, 'dt': true},
    'li': {'li': true},
    'option': {'option': true, 'optgroup': true},
    'optgroup': {'optgroup': true},
    'p': {'address': true, 'article': true, 'aside': true, 'blockquote': true, 'dir': true,
          'div': true, 'dl': true, 'fieldset': true, 'footer': true, 'form': true,
          'h1': true, 'h2': true, 'h3': true, 'h4': true, 'h5': true, 'h6': true,
          'header': true, 'hgroup': true, 'hr': true, 'menu': true, 'nav': true, 'ol': true,
          'p': true, 'pre': true, 'section': true, 'table': true, 'ul': true},
    'rp': {'rp': true, 'rt': true},
    'rt': {'rp': true, 'rt': true},
    'tbody': {'tbody': true, 'tfoot': true},
    'td': {'td': true, 'th': true},
    'tfoot': {'tbody': true},
    'th': {'td': true, 'th': true},
    'thead': {'tbody': true, 'tfoot': true},
    'tr': {'tr': true}
  },
  doNotIndent: {"pre": true},
  allowUnquoted: true,
  allowMissing: true,
  caseFold: true
}

var xmlConfig = {
  autoSelfClosers: {},
  implicitlyClosed: {},
  contextGrabbers: {},
  doNotIndent: {},
  allowUnquoted: false,
  allowMissing: false,
  allowMissingTagName: false,
  caseFold: false
}

CodeMirror.defineMode("xml", function(editorConf, config_) {
  var indentUnit = editorConf.indentUnit
  var config = {}
  var defaults = config_.htmlMode ? htmlConfig : xmlConfig
  for (var prop in defaults) config[prop] = defaults[prop]
  for (var prop in config_) config[prop] = config_[prop]

  // Return variables for tokenizers
  var type, setStyle;

  function inText(stream, state) {
    function chain(parser) {
      state.tokenize = parser;
      return parser(stream, state);
    }

    var ch = stream.next();
    if (ch == "<") {
      if (stream.eat("!")) {
        if (stream.eat("[")) {
          if (stream.match("CDATA[")) return chain(inBlock("atom", "]]>"));
          else return null;
        } else if (stream.match("--")) {
          return chain(inBlock("comment", "-->"));
        } else if (stream.match("DOCTYPE", true, true)) {
          stream.eatWhile(/[\w\._\-]/);
          return chain(doctype(1));
        } else {
          return null;
        }
      } else if (stream.eat("?")) {
        stream.eatWhile(/[\w\._\-]/);
        state.tokenize = inBlock("meta", "?>");
        return "meta";
      } else {
        type = stream.eat("/") ? "closeTag" : "openTag";
        state.tokenize = inTag;
        return "tag bracket";
      }
    } else if (ch == "&") {
      var ok;
      if (stream.eat("#")) {
        if (stream.eat("x")) {
          ok = stream.eatWhile(/[a-fA-F\d]/) && stream.eat(";");
        } else {
          ok = stream.eatWhile(/[\d]/) && stream.eat(";");
        }
      } else {
        ok = stream.eatWhile(/[\w\.\-:]/) && stream.eat(";");
      }
      return ok ? "atom" : "error";
    } else {
      stream.eatWhile(/[^&<]/);
      return null;
    }
  }
  inText.isInText = true;

  function inTag(stream, state) {
    var ch = stream.next();
    if (ch == ">" || (ch == "/" && stream.eat(">"))) {
      state.tokenize = inText;
      type = ch == ">" ? "endTag" : "selfcloseTag";
      return "tag bracket";
    } else if (ch == "=") {
      type = "equals";
      return null;
    } else if (ch == "<") {
      state.tokenize = inText;
      state.state = baseState;
      state.tagName = state.tagStart = null;
      var next = state.tokenize(stream, state);
      return next ? next + " tag error" : "tag error";
    } else if (/[\'\"]/.test(ch)) {
      state.tokenize = inAttribute(ch);
      state.stringStartCol = stream.column();
      return state.tokenize(stream, state);
    } else {
      stream.match(/^[^\s\u00a0=<>\"\']*[^\s\u00a0=<>\"\'\/]/);
      return "word";
    }
  }

  function inAttribute(quote) {
    var closure = function(stream, state) {
      while (!stream.eol()) {
        if (stream.next() == quote) {
          state.tokenize = inTag;
          break;
        }
      }
      return "string";
    };
    closure.isInAttribute = true;
    return closure;
  }

  function inBlock(style, terminator) {
    return function(stream, state) {
      while (!stream.eol()) {
        if (stream.match(terminator)) {
          state.tokenize = inText;
          break;
        }
        stream.next();
      }
      return style;
    }
  }

  function doctype(depth) {
    return function(stream, state) {
      var ch;
      while ((ch = stream.next()) != null) {
        if (ch == "<") {
          state.tokenize = doctype(depth + 1);
          return state.tokenize(stream, state);
        } else if (ch == ">") {
          if (depth == 1) {
            state.tokenize = inText;
            break;
          } else {
            state.tokenize = doctype(depth - 1);
            return state.tokenize(stream, state);
          }
        }
      }
      return "meta";
    };
  }

  function lower(tagName) {
    return tagName && tagName.toLowerCase();
  }

  function Context(state, tagName, startOfLine) {
    this.prev = state.context;
    this.tagName = tagName || "";
    this.indent = state.indented;
    this.startOfLine = startOfLine;
    if (config.doNotIndent.hasOwnProperty(tagName) || (state.context && state.context.noIndent))
      this.noIndent = true;
  }
  function popContext(state) {
    if (state.context) state.context = state.context.prev;
  }
  function maybePopContext(state, nextTagName) {
    var parentTagName;
    while (true) {
      if (!state.context) {
        return;
      }
      parentTagName = state.context.tagName;
      if (!config.contextGrabbers.hasOwnProperty(lower(parentTagName)) ||
          !config.contextGrabbers[lower(parentTagName)].hasOwnProperty(lower(nextTagName))) {
        return;
      }
      popContext(state);
    }
  }

  function baseState(type, stream, state) {
    if (type == "openTag") {
      state.tagStart = stream.column();
      return tagNameState;
    } else if (type == "closeTag") {
      return closeTagNameState;
    } else {
      return baseState;
    }
  }
  function tagNameState(type, stream, state) {
    if (type == "word") {
      state.tagName = stream.current();
      setStyle = "tag";
      return attrState;
    } else if (config.allowMissingTagName && type == "endTag") {
      setStyle = "tag bracket";
      return attrState(type, stream, state);
    } else {
      setStyle = "error";
      return tagNameState;
    }
  }
  function closeTagNameState(type, stream, state) {
    if (type == "word") {
      var tagName = stream.current();
      if (state.context && state.context.tagName != tagName &&
          config.implicitlyClosed.hasOwnProperty(lower(state.context.tagName)))
        popContext(state);
      if ((state.context && state.context.tagName == tagName) || config.matchClosing === false) {
        setStyle = "tag";
        return closeState;
      } else {
        setStyle = "tag error";
        return closeStateErr;
      }
    } else if (config.allowMissingTagName && type == "endTag") {
      setStyle = "tag bracket";
      return closeState(type, stream, state);
    } else {
      setStyle = "error";
      return closeStateErr;
    }
  }

  function closeState(type, _stream, state) {
    if (type != "endTag") {
      setStyle = "error";
      return closeState;
    }
    popContext(state);
    return baseState;
  }
  function closeStateErr(type, stream, state) {
    setStyle = "error";
    return closeState(type, stream, state);
  }

  function attrState(type, _stream, state) {
    if (type == "word") {
      setStyle = "attribute";
      return attrEqState;
    } else if (type == "endTag" || type == "selfcloseTag") {
      var tagName = state.tagName, tagStart = state.tagStart;
      state.tagName = state.tagStart = null;
      if (type == "selfcloseTag" ||
          config.autoSelfClosers.hasOwnProperty(lower(tagName))) {
        maybePopContext(state, tagName);
      } else {
        maybePopContext(state, tagName);
        state.context = new Context(state, tagName, tagStart == state.indented);
      }
      return baseState;
    }
    setStyle = "error";
    return attrState;
  }
  function attrEqState(type, stream, state) {
    if (type == "equals") return attrValueState;
    if (!config.allowMissing) setStyle = "error";
    return attrState(type, stream, state);
  }
  function attrValueState(type, stream, state) {
    if (type == "string") return attrContinuedState;
    if (type == "word" && config.allowUnquoted) {setStyle = "string"; return attrState;}
    setStyle = "error";
    return attrState(type, stream, state);
  }
  function attrContinuedState(type, stream, state) {
    if (type == "string") return attrContinuedState;
    return attrState(type, stream, state);
  }

  return {
    startState: function(baseIndent) {
      var state = {tokenize: inText,
                   state: baseState,
                   indented: baseIndent || 0,
                   tagName: null, tagStart: null,
                   context: null}
      if (baseIndent != null) state.baseIndent = baseIndent
      return state
    },

    token: function(stream, state) {
      if (!state.tagName && stream.sol())
        state.indented = stream.indentation();

      if (stream.eatSpace()) return null;
      type = null;
      var style = state.tokenize(stream, state);
      if ((style || type) && style != "comment") {
        setStyle = null;
        state.state = state.state(type || style, stream, state);
        if (setStyle)
          style = setStyle == "error" ? style + " error" : setStyle;
      }
      return style;
    },

    indent: function(state, textAfter, fullLine) {
      var context = state.context;
      // Indent multi-line strings (e.g. css).
      if (state.tokenize.isInAttribute) {
        if (state.tagStart == state.indented)
          return state.stringStartCol + 1;
        else
          return state.indented + indentUnit;
      }
      if (context && context.noIndent) return CodeMirror.Pass;
      if (state.tokenize != inTag && state.tokenize != inText)
        return fullLine ? fullLine.match(/^(\s*)/)[0].length : 0;
      // Indent the starts of attribute names.
      if (state.tagName) {
        if (config.multilineTagIndentPastTag !== false)
          return state.tagStart + state.tagName.length + 2;
        else
          return state.tagStart + indentUnit * (config.multilineTagIndentFactor || 1);
      }
      if (config.alignCDATA && /<!\[CDATA\[/.test(textAfter)) return 0;
      var tagAfter = textAfter && /^<(\/)?([\w_:\.-]*)/.exec(textAfter);
      if (tagAfter && tagAfter[1]) { // Closing tag spotted
        while (context) {
          if (context.tagName == tagAfter[2]) {
            context = context.prev;
            break;
          } else if (config.implicitlyClosed.hasOwnProperty(lower(context.tagName))) {
            context = context.prev;
          } else {
            break;
          }
        }
      } else if (tagAfter) { // Opening tag spotted
        while (context) {
          var grabbers = config.contextGrabbers[lower(context.tagName)];
          if (grabbers && grabbers.hasOwnProperty(lower(tagAfter[2])))
            context = context.prev;
          else
            break;
        }
      }
      while (context && context.prev && !context.startOfLine)
        context = context.prev;
      if (context) return context.indent + indentUnit;
      else return state.baseIndent || 0;
    },

    electricInput: /<\/[\s\w:]+>$/,
    blockCommentStart: "<!--",
    blockCommentEnd: "-->",

    configuration: config.htmlMode ? "html" : "xml",
    helperType: config.htmlMode ? "html" : "xml",

    skipAttribute: function(state) {
      if (state.state == attrValueState)
        state.state = attrState
    },

    xmlCurrentTag: function(state) {
      return state.tagName ? {name: state.tagName, close: state.type == "closeTag"} : null
    },

    xmlCurrentContext: function(state) {
      var context = []
      for (var cx = state.context; cx; cx = cx.prev)
        context.push(cx.tagName)
      return context.reverse()
    }
  };
});

CodeMirror.defineMIME("text/xml", "xml");
CodeMirror.defineMIME("application/xml", "xml");
if (!CodeMirror.mimeModes.hasOwnProperty("text/html"))
  CodeMirror.defineMIME("text/html", {name: "xml", htmlMode: true});

});
// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: https://codemirror.net/LICENSE

(function(mod) {
  if (typeof exports == "object" && typeof module == "object") // CommonJS
    mod(require("../../lib/codemirror"), require("../markdown/markdown"), require("../../addon/mode/overlay"));
  else if (typeof define == "function" && define.amd) // AMD
    define(["../../lib/codemirror", "../markdown/markdown", "../../addon/mode/overlay"], mod);
  else // Plain browser env
    mod(CodeMirror);
})(function(CodeMirror) {
"use strict";

var urlRE = /^((?:(?:aaas?|about|acap|adiumxtra|af[ps]|aim|apt|attachment|aw|beshare|bitcoin|bolo|callto|cap|chrome(?:-extension)?|cid|coap|com-eventbrite-attendee|content|crid|cvs|data|dav|dict|dlna-(?:playcontainer|playsingle)|dns|doi|dtn|dvb|ed2k|facetime|feed|file|finger|fish|ftp|geo|gg|git|gizmoproject|go|gopher|gtalk|h323|hcp|https?|iax|icap|icon|im|imap|info|ipn|ipp|irc[6s]?|iris(?:\.beep|\.lwz|\.xpc|\.xpcs)?|itms|jar|javascript|jms|keyparc|lastfm|ldaps?|magnet|mailto|maps|market|message|mid|mms|ms-help|msnim|msrps?|mtqp|mumble|mupdate|mvn|news|nfs|nih?|nntp|notes|oid|opaquelocktoken|palm|paparazzi|platform|pop|pres|proxy|psyc|query|res(?:ource)?|rmi|rsync|rtmp|rtsp|secondlife|service|session|sftp|sgn|shttp|sieve|sips?|skype|sm[bs]|snmp|soap\.beeps?|soldat|spotify|ssh|steam|svn|tag|teamspeak|tel(?:net)?|tftp|things|thismessage|tip|tn3270|tv|udp|unreal|urn|ut2004|vemmi|ventrilo|view-source|webcal|wss?|wtai|wyciwyg|xcon(?:-userid)?|xfire|xmlrpc\.beeps?|xmpp|xri|ymsgr|z39\.50[rs]?):(?:\/{1,3}|[a-z0-9%])|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}\/)(?:[^\s()<>]|\([^\s()<>]*\))+(?:\([^\s()<>]*\)|[^\s`*!()\[\]{};:'".,<>?«»“”‘’]))/i

CodeMirror.defineMode("gfm", function(config, modeConfig) {
  var codeDepth = 0;
  function blankLine(state) {
    state.code = false;
    return null;
  }
  var gfmOverlay = {
    startState: function() {
      return {
        code: false,
        codeBlock: false,
        ateSpace: false
      };
    },
    copyState: function(s) {
      return {
        code: s.code,
        codeBlock: s.codeBlock,
        ateSpace: s.ateSpace
      };
    },
    token: function(stream, state) {
      state.combineTokens = null;

      // Hack to prevent formatting override inside code blocks (block and inline)
      if (state.codeBlock) {
        if (stream.match(/^```+/)) {
          state.codeBlock = false;
          return null;
        }
        stream.skipToEnd();
        return null;
      }
      if (stream.sol()) {
        state.code = false;
      }
      if (stream.sol() && stream.match(/^```+/)) {
        stream.skipToEnd();
        state.codeBlock = true;
        return null;
      }
      // If this block is changed, it may need to be updated in Markdown mode
      if (stream.peek() === '`') {
        stream.next();
        var before = stream.pos;
        stream.eatWhile('`');
        var difference = 1 + stream.pos - before;
        if (!state.code) {
          codeDepth = difference;
          state.code = true;
        } else {
          if (difference === codeDepth) { // Must be exact
            state.code = false;
          }
        }
        return null;
      } else if (state.code) {
        stream.next();
        return null;
      }
      // Check if space. If so, links can be formatted later on
      if (stream.eatSpace()) {
        state.ateSpace = true;
        return null;
      }
      if (stream.sol() || state.ateSpace) {
        state.ateSpace = false;
        if (modeConfig.gitHubSpice !== false) {
          if(stream.match(/^(?:[a-zA-Z0-9\-_]+\/)?(?:[a-zA-Z0-9\-_]+@)?(?=.{0,6}\d)(?:[a-f0-9]{7,40}\b)/)) {
            // User/Project@SHA
            // User@SHA
            // SHA
            state.combineTokens = true;
            return "link";
          } else if (stream.match(/^(?:[a-zA-Z0-9\-_]+\/)?(?:[a-zA-Z0-9\-_]+)?#[0-9]+\b/)) {
            // User/Project#Num
            // User#Num
            // #Num
            state.combineTokens = true;
            return "link";
          }
        }
      }
      if (stream.match(urlRE) &&
          stream.string.slice(stream.start - 2, stream.start) != "](" &&
          (stream.start == 0 || /\W/.test(stream.string.charAt(stream.start - 1)))) {
        // URLs
        // Taken from http://daringfireball.net/2010/07/improved_regex_for_matching_urls
        // And then (issue #1160) simplified to make it not crash the Chrome Regexp engine
        // And then limited url schemes to the CommonMark list, so foo:bar isn't matched as a URL
        state.combineTokens = true;
        return "link";
      }
      stream.next();
      return null;
    },
    blankLine: blankLine
  };

  var markdownConfig = {
    taskLists: true,
    strikethrough: true,
    emoji: true
  };
  for (var attr in modeConfig) {
    markdownConfig[attr] = modeConfig[attr];
  }
  markdownConfig.name = "markdown";
  return CodeMirror.overlayMode(CodeMirror.getMode(config, markdownConfig), gfmOverlay);

}, "markdown");

  CodeMirror.defineMIME("text/x-gfm", "gfm");
});
// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: https://codemirror.net/LICENSE

(function(mod) {
  if (typeof exports == "object" && typeof module == "object") // CommonJS
    mod(require("../../lib/codemirror"));
  else if (typeof define == "function" && define.amd) // AMD
    define(["../../lib/codemirror"], mod);
  else // Plain browser env
    mod(CodeMirror);
})(function(CodeMirror) {
"use strict";

CodeMirror.defineMode("javascript", function(config, parserConfig) {
  var indentUnit = config.indentUnit;
  var statementIndent = parserConfig.statementIndent;
  var jsonldMode = parserConfig.jsonld;
  var jsonMode = parserConfig.json || jsonldMode;
  var trackScope = parserConfig.trackScope !== false
  var isTS = parserConfig.typescript;
  var wordRE = parserConfig.wordCharacters || /[\w$\xa1-\uffff]/;

  // Tokenizer

  var keywords = function(){
    function kw(type) {return {type: type, style: "keyword"};}
    var A = kw("keyword a"), B = kw("keyword b"), C = kw("keyword c"), D = kw("keyword d");
    var operator = kw("operator"), atom = {type: "atom", style: "atom"};

    return {
      "if": kw("if"), "while": A, "with": A, "else": B, "do": B, "try": B, "finally": B,
      "return": D, "break": D, "continue": D, "new": kw("new"), "delete": C, "void": C, "throw": C,
      "debugger": kw("debugger"), "var": kw("var"), "const": kw("var"), "let": kw("var"),
      "function": kw("function"), "catch": kw("catch"),
      "for": kw("for"), "switch": kw("switch"), "case": kw("case"), "default": kw("default"),
      "in": operator, "typeof": operator, "instanceof": operator,
      "true": atom, "false": atom, "null": atom, "undefined": atom, "NaN": atom, "Infinity": atom,
      "this": kw("this"), "class": kw("class"), "super": kw("atom"),
      "yield": C, "export": kw("export"), "import": kw("import"), "extends": C,
      "await": C
    };
  }();

  var isOperatorChar = /[+\-*&%=<>!?|~^@]/;
  var isJsonldKeyword = /^@(context|id|value|language|type|container|list|set|reverse|index|base|vocab|graph)"/;

  function readRegexp(stream) {
    var escaped = false, next, inSet = false;
    while ((next = stream.next()) != null) {
      if (!escaped) {
        if (next == "/" && !inSet) return;
        if (next == "[") inSet = true;
        else if (inSet && next == "]") inSet = false;
      }
      escaped = !escaped && next == "\\";
    }
  }

  // Used as scratch variables to communicate multiple values without
  // consing up tons of objects.
  var type, content;
  function ret(tp, style, cont) {
    type = tp; content = cont;
    return style;
  }
  function tokenBase(stream, state) {
    var ch = stream.next();
    if (ch == '"' || ch == "'") {
      state.tokenize = tokenString(ch);
      return state.tokenize(stream, state);
    } else if (ch == "." && stream.match(/^\d[\d_]*(?:[eE][+\-]?[\d_]+)?/)) {
      return ret("number", "number");
    } else if (ch == "." && stream.match("..")) {
      return ret("spread", "meta");
    } else if (/[\[\]{}\(\),;\:\.]/.test(ch)) {
      return ret(ch);
    } else if (ch == "=" && stream.eat(">")) {
      return ret("=>", "operator");
    } else if (ch == "0" && stream.match(/^(?:x[\dA-Fa-f_]+|o[0-7_]+|b[01_]+)n?/)) {
      return ret("number", "number");
    } else if (/\d/.test(ch)) {
      stream.match(/^[\d_]*(?:n|(?:\.[\d_]*)?(?:[eE][+\-]?[\d_]+)?)?/);
      return ret("number", "number");
    } else if (ch == "/") {
      if (stream.eat("*")) {
        state.tokenize = tokenComment;
        return tokenComment(stream, state);
      } else if (stream.eat("/")) {
        stream.skipToEnd();
        return ret("comment", "comment");
      } else if (expressionAllowed(stream, state, 1)) {
        readRegexp(stream);
        stream.match(/^\b(([gimyus])(?![gimyus]*\2))+\b/);
        return ret("regexp", "string-2");
      } else {
        stream.eat("=");
        return ret("operator", "operator", stream.current());
      }
    } else if (ch == "`") {
      state.tokenize = tokenQuasi;
      return tokenQuasi(stream, state);
    } else if (ch == "#" && stream.peek() == "!") {
      stream.skipToEnd();
      return ret("meta", "meta");
    } else if (ch == "#" && stream.eatWhile(wordRE)) {
      return ret("variable", "property")
    } else if (ch == "<" && stream.match("!--") ||
               (ch == "-" && stream.match("->") && !/\S/.test(stream.string.slice(0, stream.start)))) {
      stream.skipToEnd()
      return ret("comment", "comment")
    } else if (isOperatorChar.test(ch)) {
      if (ch != ">" || !state.lexical || state.lexical.type != ">") {
        if (stream.eat("=")) {
          if (ch == "!" || ch == "=") stream.eat("=")
        } else if (/[<>*+\-|&?]/.test(ch)) {
          stream.eat(ch)
          if (ch == ">") stream.eat(ch)
        }
      }
      if (ch == "?" && stream.eat(".")) return ret(".")
      return ret("operator", "operator", stream.current());
    } else if (wordRE.test(ch)) {
      stream.eatWhile(wordRE);
      var word = stream.current()
      if (state.lastType != ".") {
        if (keywords.propertyIsEnumerable(word)) {
          var kw = keywords[word]
          return ret(kw.type, kw.style, word)
        }
        if (word == "async" && stream.match(/^(\s|\/\*([^*]|\*(?!\/))*?\*\/)*[\[\(\w]/, false))
          return ret("async", "keyword", word)
      }
      return ret("variable", "variable", word)
    }
  }

  function tokenString(quote) {
    return function(stream, state) {
      var escaped = false, next;
      if (jsonldMode && stream.peek() == "@" && stream.match(isJsonldKeyword)){
        state.tokenize = tokenBase;
        return ret("jsonld-keyword", "meta");
      }
      while ((next = stream.next()) != null) {
        if (next == quote && !escaped) break;
        escaped = !escaped && next == "\\";
      }
      if (!escaped) state.tokenize = tokenBase;
      return ret("string", "string");
    };
  }

  function tokenComment(stream, state) {
    var maybeEnd = false, ch;
    while (ch = stream.next()) {
      if (ch == "/" && maybeEnd) {
        state.tokenize = tokenBase;
        break;
      }
      maybeEnd = (ch == "*");
    }
    return ret("comment", "comment");
  }

  function tokenQuasi(stream, state) {
    var escaped = false, next;
    while ((next = stream.next()) != null) {
      if (!escaped && (next == "`" || next == "$" && stream.eat("{"))) {
        state.tokenize = tokenBase;
        break;
      }
      escaped = !escaped && next == "\\";
    }
    return ret("quasi", "string-2", stream.current());
  }

  var brackets = "([{}])";
  // This is a crude lookahead trick to try and notice that we're
  // parsing the argument patterns for a fat-arrow function before we
  // actually hit the arrow token. It only works if the arrow is on
  // the same line as the arguments and there's no strange noise
  // (comments) in between. Fallback is to only notice when we hit the
  // arrow, and not declare the arguments as locals for the arrow
  // body.
  function findFatArrow(stream, state) {
    if (state.fatArrowAt) state.fatArrowAt = null;
    var arrow = stream.string.indexOf("=>", stream.start);
    if (arrow < 0) return;

    if (isTS) { // Try to skip TypeScript return type declarations after the arguments
      var m = /:\s*(?:\w+(?:<[^>]*>|\[\])?|\{[^}]*\})\s*$/.exec(stream.string.slice(stream.start, arrow))
      if (m) arrow = m.index
    }

    var depth = 0, sawSomething = false;
    for (var pos = arrow - 1; pos >= 0; --pos) {
      var ch = stream.string.charAt(pos);
      var bracket = brackets.indexOf(ch);
      if (bracket >= 0 && bracket < 3) {
        if (!depth) { ++pos; break; }
        if (--depth == 0) { if (ch == "(") sawSomething = true; break; }
      } else if (bracket >= 3 && bracket < 6) {
        ++depth;
      } else if (wordRE.test(ch)) {
        sawSomething = true;
      } else if (/["'\/`]/.test(ch)) {
        for (;; --pos) {
          if (pos == 0) return
          var next = stream.string.charAt(pos - 1)
          if (next == ch && stream.string.charAt(pos - 2) != "\\") { pos--; break }
        }
      } else if (sawSomething && !depth) {
        ++pos;
        break;
      }
    }
    if (sawSomething && !depth) state.fatArrowAt = pos;
  }

  // Parser

  var atomicTypes = {"atom": true, "number": true, "variable": true, "string": true,
                     "regexp": true, "this": true, "import": true, "jsonld-keyword": true};

  function JSLexical(indented, column, type, align, prev, info) {
    this.indented = indented;
    this.column = column;
    this.type = type;
    this.prev = prev;
    this.info = info;
    if (align != null) this.align = align;
  }

  function inScope(state, varname) {
    if (!trackScope) return false
    for (var v = state.localVars; v; v = v.next)
      if (v.name == varname) return true;
    for (var cx = state.context; cx; cx = cx.prev) {
      for (var v = cx.vars; v; v = v.next)
        if (v.name == varname) return true;
    }
  }

  function parseJS(state, style, type, content, stream) {
    var cc = state.cc;
    // Communicate our context to the combinators.
    // (Less wasteful than consing up a hundred closures on every call.)
    cx.state = state; cx.stream = stream; cx.marked = null, cx.cc = cc; cx.style = style;

    if (!state.lexical.hasOwnProperty("align"))
      state.lexical.align = true;

    while(true) {
      var combinator = cc.length ? cc.pop() : jsonMode ? expression : statement;
      if (combinator(type, content)) {
        while(cc.length && cc[cc.length - 1].lex)
          cc.pop()();
        if (cx.marked) return cx.marked;
        if (type == "variable" && inScope(state, content)) return "variable-2";
        return style;
      }
    }
  }

  // Combinator utils

  var cx = {state: null, column: null, marked: null, cc: null};
  function pass() {
    for (var i = arguments.length - 1; i >= 0; i--) cx.cc.push(arguments[i]);
  }
  function cont() {
    pass.apply(null, arguments);
    return true;
  }
  function inList(name, list) {
    for (var v = list; v; v = v.next) if (v.name == name) return true
    return false;
  }
  function register(varname) {
    var state = cx.state;
    cx.marked = "def";
    if (!trackScope) return
    if (state.context) {
      if (state.lexical.info == "var" && state.context && state.context.block) {
        // FIXME function decls are also not block scoped
        var newContext = registerVarScoped(varname, state.context)
        if (newContext != null) {
          state.context = newContext
          return
        }
      } else if (!inList(varname, state.localVars)) {
        state.localVars = new Var(varname, state.localVars)
        return
      }
    }
    // Fall through means this is global
    if (parserConfig.globalVars && !inList(varname, state.globalVars))
      state.globalVars = new Var(varname, state.globalVars)
  }
  function registerVarScoped(varname, context) {
    if (!context) {
      return null
    } else if (context.block) {
      var inner = registerVarScoped(varname, context.prev)
      if (!inner) return null
      if (inner == context.prev) return context
      return new Context(inner, context.vars, true)
    } else if (inList(varname, context.vars)) {
      return context
    } else {
      return new Context(context.prev, new Var(varname, context.vars), false)
    }
  }

  function isModifier(name) {
    return name == "public" || name == "private" || name == "protected" || name == "abstract" || name == "readonly"
  }

  // Combinators

  function Context(prev, vars, block) { this.prev = prev; this.vars = vars; this.block = block }
  function Var(name, next) { this.name = name; this.next = next }

  var defaultVars = new Var("this", new Var("arguments", null))
  function pushcontext() {
    cx.state.context = new Context(cx.state.context, cx.state.localVars, false)
    cx.state.localVars = defaultVars
  }
  function pushblockcontext() {
    cx.state.context = new Context(cx.state.context, cx.state.localVars, true)
    cx.state.localVars = null
  }
  pushcontext.lex = pushblockcontext.lex = true
  function popcontext() {
    cx.state.localVars = cx.state.context.vars
    cx.state.context = cx.state.context.prev
  }
  popcontext.lex = true
  function pushlex(type, info) {
    var result = function() {
      var state = cx.state, indent = state.indented;
      if (state.lexical.type == "stat") indent = state.lexical.indented;
      else for (var outer = state.lexical; outer && outer.type == ")" && outer.align; outer = outer.prev)
        indent = outer.indented;
      state.lexical = new JSLexical(indent, cx.stream.column(), type, null, state.lexical, info);
    };
    result.lex = true;
    return result;
  }
  function poplex() {
    var state = cx.state;
    if (state.lexical.prev) {
      if (state.lexical.type == ")")
        state.indented = state.lexical.indented;
      state.lexical = state.lexical.prev;
    }
  }
  poplex.lex = true;

  function expect(wanted) {
    function exp(type) {
      if (type == wanted) return cont();
      else if (wanted == ";" || type == "}" || type == ")" || type == "]") return pass();
      else return cont(exp);
    };
    return exp;
  }

  function statement(type, value) {
    if (type == "var") return cont(pushlex("vardef", value), vardef, expect(";"), poplex);
    if (type == "keyword a") return cont(pushlex("form"), parenExpr, statement, poplex);
    if (type == "keyword b") return cont(pushlex("form"), statement, poplex);
    if (type == "keyword d") return cx.stream.match(/^\s*$/, false) ? cont() : cont(pushlex("stat"), maybeexpression, expect(";"), poplex);
    if (type == "debugger") return cont(expect(";"));
    if (type == "{") return cont(pushlex("}"), pushblockcontext, block, poplex, popcontext);
    if (type == ";") return cont();
    if (type == "if") {
      if (cx.state.lexical.info == "else" && cx.state.cc[cx.state.cc.length - 1] == poplex)
        cx.state.cc.pop()();
      return cont(pushlex("form"), parenExpr, statement, poplex, maybeelse);
    }
    if (type == "function") return cont(functiondef);
    if (type == "for") return cont(pushlex("form"), pushblockcontext, forspec, statement, popcontext, poplex);
    if (type == "class" || (isTS && value == "interface")) {
      cx.marked = "keyword"
      return cont(pushlex("form", type == "class" ? type : value), className, poplex)
    }
    if (type == "variable") {
      if (isTS && value == "declare") {
        cx.marked = "keyword"
        return cont(statement)
      } else if (isTS && (value == "module" || value == "enum" || value == "type") && cx.stream.match(/^\s*\w/, false)) {
        cx.marked = "keyword"
        if (value == "enum") return cont(enumdef);
        else if (value == "type") return cont(typename, expect("operator"), typeexpr, expect(";"));
        else return cont(pushlex("form"), pattern, expect("{"), pushlex("}"), block, poplex, poplex)
      } else if (isTS && value == "namespace") {
        cx.marked = "keyword"
        return cont(pushlex("form"), expression, statement, poplex)
      } else if (isTS && value == "abstract") {
        cx.marked = "keyword"
        return cont(statement)
      } else {
        return cont(pushlex("stat"), maybelabel);
      }
    }
    if (type == "switch") return cont(pushlex("form"), parenExpr, expect("{"), pushlex("}", "switch"), pushblockcontext,
                                      block, poplex, poplex, popcontext);
    if (type == "case") return cont(expression, expect(":"));
    if (type == "default") return cont(expect(":"));
    if (type == "catch") return cont(pushlex("form"), pushcontext, maybeCatchBinding, statement, poplex, popcontext);
    if (type == "export") return cont(pushlex("stat"), afterExport, poplex);
    if (type == "import") return cont(pushlex("stat"), afterImport, poplex);
    if (type == "async") return cont(statement)
    if (value == "@") return cont(expression, statement)
    return pass(pushlex("stat"), expression, expect(";"), poplex);
  }
  function maybeCatchBinding(type) {
    if (type == "(") return cont(funarg, expect(")"))
  }
  function expression(type, value) {
    return expressionInner(type, value, false);
  }
  function expressionNoComma(type, value) {
    return expressionInner(type, value, true);
  }
  function parenExpr(type) {
    if (type != "(") return pass()
    return cont(pushlex(")"), maybeexpression, expect(")"), poplex)
  }
  function expressionInner(type, value, noComma) {
    if (cx.state.fatArrowAt == cx.stream.start) {
      var body = noComma ? arrowBodyNoComma : arrowBody;
      if (type == "(") return cont(pushcontext, pushlex(")"), commasep(funarg, ")"), poplex, expect("=>"), body, popcontext);
      else if (type == "variable") return pass(pushcontext, pattern, expect("=>"), body, popcontext);
    }

    var maybeop = noComma ? maybeoperatorNoComma : maybeoperatorComma;
    if (atomicTypes.hasOwnProperty(type)) return cont(maybeop);
    if (type == "function") return cont(functiondef, maybeop);
    if (type == "class" || (isTS && value == "interface")) { cx.marked = "keyword"; return cont(pushlex("form"), classExpression, poplex); }
    if (type == "keyword c" || type == "async") return cont(noComma ? expressionNoComma : expression);
    if (type == "(") return cont(pushlex(")"), maybeexpression, expect(")"), poplex, maybeop);
    if (type == "operator" || type == "spread") return cont(noComma ? expressionNoComma : expression);
    if (type == "[") return cont(pushlex("]"), arrayLiteral, poplex, maybeop);
    if (type == "{") return contCommasep(objprop, "}", null, maybeop);
    if (type == "quasi") return pass(quasi, maybeop);
    if (type == "new") return cont(maybeTarget(noComma));
    return cont();
  }
  function maybeexpression(type) {
    if (type.match(/[;\}\)\],]/)) return pass();
    return pass(expression);
  }

  function maybeoperatorComma(type, value) {
    if (type == ",") return cont(maybeexpression);
    return maybeoperatorNoComma(type, value, false);
  }
  function maybeoperatorNoComma(type, value, noComma) {
    var me = noComma == false ? maybeoperatorComma : maybeoperatorNoComma;
    var expr = noComma == false ? expression : expressionNoComma;
    if (type == "=>") return cont(pushcontext, noComma ? arrowBodyNoComma : arrowBody, popcontext);
    if (type == "operator") {
      if (/\+\+|--/.test(value) || isTS && value == "!") return cont(me);
      if (isTS && value == "<" && cx.stream.match(/^([^<>]|<[^<>]*>)*>\s*\(/, false))
        return cont(pushlex(">"), commasep(typeexpr, ">"), poplex, me);
      if (value == "?") return cont(expression, expect(":"), expr);
      return cont(expr);
    }
    if (type == "quasi") { return pass(quasi, me); }
    if (type == ";") return;
    if (type == "(") return contCommasep(expressionNoComma, ")", "call", me);
    if (type == ".") return cont(property, me);
    if (type == "[") return cont(pushlex("]"), maybeexpression, expect("]"), poplex, me);
    if (isTS && value == "as") { cx.marked = "keyword"; return cont(typeexpr, me) }
    if (type == "regexp") {
      cx.state.lastType = cx.marked = "operator"
      cx.stream.backUp(cx.stream.pos - cx.stream.start - 1)
      return cont(expr)
    }
  }
  function quasi(type, value) {
    if (type != "quasi") return pass();
    if (value.slice(value.length - 2) != "${") return cont(quasi);
    return cont(maybeexpression, continueQuasi);
  }
  function continueQuasi(type) {
    if (type == "}") {
      cx.marked = "string-2";
      cx.state.tokenize = tokenQuasi;
      return cont(quasi);
    }
  }
  function arrowBody(type) {
    findFatArrow(cx.stream, cx.state);
    return pass(type == "{" ? statement : expression);
  }
  function arrowBodyNoComma(type) {
    findFatArrow(cx.stream, cx.state);
    return pass(type == "{" ? statement : expressionNoComma);
  }
  function maybeTarget(noComma) {
    return function(type) {
      if (type == ".") return cont(noComma ? targetNoComma : target);
      else if (type == "variable" && isTS) return cont(maybeTypeArgs, noComma ? maybeoperatorNoComma : maybeoperatorComma)
      else return pass(noComma ? expressionNoComma : expression);
    };
  }
  function target(_, value) {
    if (value == "target") { cx.marked = "keyword"; return cont(maybeoperatorComma); }
  }
  function targetNoComma(_, value) {
    if (value == "target") { cx.marked = "keyword"; return cont(maybeoperatorNoComma); }
  }
  function maybelabel(type) {
    if (type == ":") return cont(poplex, statement);
    return pass(maybeoperatorComma, expect(";"), poplex);
  }
  function property(type) {
    if (type == "variable") {cx.marked = "property"; return cont();}
  }
  function objprop(type, value) {
    if (type == "async") {
      cx.marked = "property";
      return cont(objprop);
    } else if (type == "variable" || cx.style == "keyword") {
      cx.marked = "property";
      if (value == "get" || value == "set") return cont(getterSetter);
      var m // Work around fat-arrow-detection complication for detecting typescript typed arrow params
      if (isTS && cx.state.fatArrowAt == cx.stream.start && (m = cx.stream.match(/^\s*:\s*/, false)))
        cx.state.fatArrowAt = cx.stream.pos + m[0].length
      return cont(afterprop);
    } else if (type == "number" || type == "string") {
      cx.marked = jsonldMode ? "property" : (cx.style + " property");
      return cont(afterprop);
    } else if (type == "jsonld-keyword") {
      return cont(afterprop);
    } else if (isTS && isModifier(value)) {
      cx.marked = "keyword"
      return cont(objprop)
    } else if (type == "[") {
      return cont(expression, maybetype, expect("]"), afterprop);
    } else if (type == "spread") {
      return cont(expressionNoComma, afterprop);
    } else if (value == "*") {
      cx.marked = "keyword";
      return cont(objprop);
    } else if (type == ":") {
      return pass(afterprop)
    }
  }
  function getterSetter(type) {
    if (type != "variable") return pass(afterprop);
    cx.marked = "property";
    return cont(functiondef);
  }
  function afterprop(type) {
    if (type == ":") return cont(expressionNoComma);
    if (type == "(") return pass(functiondef);
  }
  function commasep(what, end, sep) {
    function proceed(type, value) {
      if (sep ? sep.indexOf(type) > -1 : type == ",") {
        var lex = cx.state.lexical;
        if (lex.info == "call") lex.pos = (lex.pos || 0) + 1;
        return cont(function(type, value) {
          if (type == end || value == end) return pass()
          return pass(what)
        }, proceed);
      }
      if (type == end || value == end) return cont();
      if (sep && sep.indexOf(";") > -1) return pass(what)
      return cont(expect(end));
    }
    return function(type, value) {
      if (type == end || value == end) return cont();
      return pass(what, proceed);
    };
  }
  function contCommasep(what, end, info) {
    for (var i = 3; i < arguments.length; i++)
      cx.cc.push(arguments[i]);
    return cont(pushlex(end, info), commasep(what, end), poplex);
  }
  function block(type) {
    if (type == "}") return cont();
    return pass(statement, block);
  }
  function maybetype(type, value) {
    if (isTS) {
      if (type == ":") return cont(typeexpr);
      if (value == "?") return cont(maybetype);
    }
  }
  function maybetypeOrIn(type, value) {
    if (isTS && (type == ":" || value == "in")) return cont(typeexpr)
  }
  function mayberettype(type) {
    if (isTS && type == ":") {
      if (cx.stream.match(/^\s*\w+\s+is\b/, false)) return cont(expression, isKW, typeexpr)
      else return cont(typeexpr)
    }
  }
  function isKW(_, value) {
    if (value == "is") {
      cx.marked = "keyword"
      return cont()
    }
  }
  function typeexpr(type, value) {
    if (value == "keyof" || value == "typeof" || value == "infer" || value == "readonly") {
      cx.marked = "keyword"
      return cont(value == "typeof" ? expressionNoComma : typeexpr)
    }
    if (type == "variable" || value == "void") {
      cx.marked = "type"
      return cont(afterType)
    }
    if (value == "|" || value == "&") return cont(typeexpr)
    if (type == "string" || type == "number" || type == "atom") return cont(afterType);
    if (type == "[") return cont(pushlex("]"), commasep(typeexpr, "]", ","), poplex, afterType)
    if (type == "{") return cont(pushlex("}"), typeprops, poplex, afterType)
    if (type == "(") return cont(commasep(typearg, ")"), maybeReturnType, afterType)
    if (type == "<") return cont(commasep(typeexpr, ">"), typeexpr)
    if (type == "quasi") { return pass(quasiType, afterType); }
  }
  function maybeReturnType(type) {
    if (type == "=>") return cont(typeexpr)
  }
  function typeprops(type) {
    if (type.match(/[\}\)\]]/)) return cont()
    if (type == "," || type == ";") return cont(typeprops)
    return pass(typeprop, typeprops)
  }
  function typeprop(type, value) {
    if (type == "variable" || cx.style == "keyword") {
      cx.marked = "property"
      return cont(typeprop)
    } else if (value == "?" || type == "number" || type == "string") {
      return cont(typeprop)
    } else if (type == ":") {
      return cont(typeexpr)
    } else if (type == "[") {
      return cont(expect("variable"), maybetypeOrIn, expect("]"), typeprop)
    } else if (type == "(") {
      return pass(functiondecl, typeprop)
    } else if (!type.match(/[;\}\)\],]/)) {
      return cont()
    }
  }
  function quasiType(type, value) {
    if (type != "quasi") return pass();
    if (value.slice(value.length - 2) != "${") return cont(quasiType);
    return cont(typeexpr, continueQuasiType);
  }
  function continueQuasiType(type) {
    if (type == "}") {
      cx.marked = "string-2";
      cx.state.tokenize = tokenQuasi;
      return cont(quasiType);
    }
  }
  function typearg(type, value) {
    if (type == "variable" && cx.stream.match(/^\s*[?:]/, false) || value == "?") return cont(typearg)
    if (type == ":") return cont(typeexpr)
    if (type == "spread") return cont(typearg)
    return pass(typeexpr)
  }
  function afterType(type, value) {
    if (value == "<") return cont(pushlex(">"), commasep(typeexpr, ">"), poplex, afterType)
    if (value == "|" || type == "." || value == "&") return cont(typeexpr)
    if (type == "[") return cont(typeexpr, expect("]"), afterType)
    if (value == "extends" || value == "implements") { cx.marked = "keyword"; return cont(typeexpr) }
    if (value == "?") return cont(typeexpr, expect(":"), typeexpr)
  }
  function maybeTypeArgs(_, value) {
    if (value == "<") return cont(pushlex(">"), commasep(typeexpr, ">"), poplex, afterType)
  }
  function typeparam() {
    return pass(typeexpr, maybeTypeDefault)
  }
  function maybeTypeDefault(_, value) {
    if (value == "=") return cont(typeexpr)
  }
  function vardef(_, value) {
    if (value == "enum") {cx.marked = "keyword"; return cont(enumdef)}
    return pass(pattern, maybetype, maybeAssign, vardefCont);
  }
  function pattern(type, value) {
    if (isTS && isModifier(value)) { cx.marked = "keyword"; return cont(pattern) }
    if (type == "variable") { register(value); return cont(); }
    if (type == "spread") return cont(pattern);
    if (type == "[") return contCommasep(eltpattern, "]");
    if (type == "{") return contCommasep(proppattern, "}");
  }
  function proppattern(type, value) {
    if (type == "variable" && !cx.stream.match(/^\s*:/, false)) {
      register(value);
      return cont(maybeAssign);
    }
    if (type == "variable") cx.marked = "property";
    if (type == "spread") return cont(pattern);
    if (type == "}") return pass();
    if (type == "[") return cont(expression, expect(']'), expect(':'), proppattern);
    return cont(expect(":"), pattern, maybeAssign);
  }
  function eltpattern() {
    return pass(pattern, maybeAssign)
  }
  function maybeAssign(_type, value) {
    if (value == "=") return cont(expressionNoComma);
  }
  function vardefCont(type) {
    if (type == ",") return cont(vardef);
  }
  function maybeelse(type, value) {
    if (type == "keyword b" && value == "else") return cont(pushlex("form", "else"), statement, poplex);
  }
  function forspec(type, value) {
    if (value == "await") return cont(forspec);
    if (type == "(") return cont(pushlex(")"), forspec1, poplex);
  }
  function forspec1(type) {
    if (type == "var") return cont(vardef, forspec2);
    if (type == "variable") return cont(forspec2);
    return pass(forspec2)
  }
  function forspec2(type, value) {
    if (type == ")") return cont()
    if (type == ";") return cont(forspec2)
    if (value == "in" || value == "of") { cx.marked = "keyword"; return cont(expression, forspec2) }
    return pass(expression, forspec2)
  }
  function functiondef(type, value) {
    if (value == "*") {cx.marked = "keyword"; return cont(functiondef);}
    if (type == "variable") {register(value); return cont(functiondef);}
    if (type == "(") return cont(pushcontext, pushlex(")"), commasep(funarg, ")"), poplex, mayberettype, statement, popcontext);
    if (isTS && value == "<") return cont(pushlex(">"), commasep(typeparam, ">"), poplex, functiondef)
  }
  function functiondecl(type, value) {
    if (value == "*") {cx.marked = "keyword"; return cont(functiondecl);}
    if (type == "variable") {register(value); return cont(functiondecl);}
    if (type == "(") return cont(pushcontext, pushlex(")"), commasep(funarg, ")"), poplex, mayberettype, popcontext);
    if (isTS && value == "<") return cont(pushlex(">"), commasep(typeparam, ">"), poplex, functiondecl)
  }
  function typename(type, value) {
    if (type == "keyword" || type == "variable") {
      cx.marked = "type"
      return cont(typename)
    } else if (value == "<") {
      return cont(pushlex(">"), commasep(typeparam, ">"), poplex)
    }
  }
  function funarg(type, value) {
    if (value == "@") cont(expression, funarg)
    if (type == "spread") return cont(funarg);
    if (isTS && isModifier(value)) { cx.marked = "keyword"; return cont(funarg); }
    if (isTS && type == "this") return cont(maybetype, maybeAssign)
    return pass(pattern, maybetype, maybeAssign);
  }
  function classExpression(type, value) {
    // Class expressions may have an optional name.
    if (type == "variable") return className(type, value);
    return classNameAfter(type, value);
  }
  function className(type, value) {
    if (type == "variable") {register(value); return cont(classNameAfter);}
  }
  function classNameAfter(type, value) {
    if (value == "<") return cont(pushlex(">"), commasep(typeparam, ">"), poplex, classNameAfter)
    if (value == "extends" || value == "implements" || (isTS && type == ",")) {
      if (value == "implements") cx.marked = "keyword";
      return cont(isTS ? typeexpr : expression, classNameAfter);
    }
    if (type == "{") return cont(pushlex("}"), classBody, poplex);
  }
  function classBody(type, value) {
    if (type == "async" ||
        (type == "variable" &&
         (value == "static" || value == "get" || value == "set" || (isTS && isModifier(value))) &&
         cx.stream.match(/^\s+[\w$\xa1-\uffff]/, false))) {
      cx.marked = "keyword";
      return cont(classBody);
    }
    if (type == "variable" || cx.style == "keyword") {
      cx.marked = "property";
      return cont(classfield, classBody);
    }
    if (type == "number" || type == "string") return cont(classfield, classBody);
    if (type == "[")
      return cont(expression, maybetype, expect("]"), classfield, classBody)
    if (value == "*") {
      cx.marked = "keyword";
      return cont(classBody);
    }
    if (isTS && type == "(") return pass(functiondecl, classBody)
    if (type == ";" || type == ",") return cont(classBody);
    if (type == "}") return cont();
    if (value == "@") return cont(expression, classBody)
  }
  function classfield(type, value) {
    if (value == "!") return cont(classfield)
    if (value == "?") return cont(classfield)
    if (type == ":") return cont(typeexpr, maybeAssign)
    if (value == "=") return cont(expressionNoComma)
    var context = cx.state.lexical.prev, isInterface = context && context.info == "interface"
    return pass(isInterface ? functiondecl : functiondef)
  }
  function afterExport(type, value) {
    if (value == "*") { cx.marked = "keyword"; return cont(maybeFrom, expect(";")); }
    if (value == "default") { cx.marked = "keyword"; return cont(expression, expect(";")); }
    if (type == "{") return cont(commasep(exportField, "}"), maybeFrom, expect(";"));
    return pass(statement);
  }
  function exportField(type, value) {
    if (value == "as") { cx.marked = "keyword"; return cont(expect("variable")); }
    if (type == "variable") return pass(expressionNoComma, exportField);
  }
  function afterImport(type) {
    if (type == "string") return cont();
    if (type == "(") return pass(expression);
    if (type == ".") return pass(maybeoperatorComma);
    return pass(importSpec, maybeMoreImports, maybeFrom);
  }
  function importSpec(type, value) {
    if (type == "{") return contCommasep(importSpec, "}");
    if (type == "variable") register(value);
    if (value == "*") cx.marked = "keyword";
    return cont(maybeAs);
  }
  function maybeMoreImports(type) {
    if (type == ",") return cont(importSpec, maybeMoreImports)
  }
  function maybeAs(_type, value) {
    if (value == "as") { cx.marked = "keyword"; return cont(importSpec); }
  }
  function maybeFrom(_type, value) {
    if (value == "from") { cx.marked = "keyword"; return cont(expression); }
  }
  function arrayLiteral(type) {
    if (type == "]") return cont();
    return pass(commasep(expressionNoComma, "]"));
  }
  function enumdef() {
    return pass(pushlex("form"), pattern, expect("{"), pushlex("}"), commasep(enummember, "}"), poplex, poplex)
  }
  function enummember() {
    return pass(pattern, maybeAssign);
  }

  function isContinuedStatement(state, textAfter) {
    return state.lastType == "operator" || state.lastType == "," ||
      isOperatorChar.test(textAfter.charAt(0)) ||
      /[,.]/.test(textAfter.charAt(0));
  }

  function expressionAllowed(stream, state, backUp) {
    return state.tokenize == tokenBase &&
      /^(?:operator|sof|keyword [bcd]|case|new|export|default|spread|[\[{}\(,;:]|=>)$/.test(state.lastType) ||
      (state.lastType == "quasi" && /\{\s*$/.test(stream.string.slice(0, stream.pos - (backUp || 0))))
  }

  // Interface

  return {
    startState: function(basecolumn) {
      var state = {
        tokenize: tokenBase,
        lastType: "sof",
        cc: [],
        lexical: new JSLexical((basecolumn || 0) - indentUnit, 0, "block", false),
        localVars: parserConfig.localVars,
        context: parserConfig.localVars && new Context(null, null, false),
        indented: basecolumn || 0
      };
      if (parserConfig.globalVars && typeof parserConfig.globalVars == "object")
        state.globalVars = parserConfig.globalVars;
      return state;
    },

    token: function(stream, state) {
      if (stream.sol()) {
        if (!state.lexical.hasOwnProperty("align"))
          state.lexical.align = false;
        state.indented = stream.indentation();
        findFatArrow(stream, state);
      }
      if (state.tokenize != tokenComment && stream.eatSpace()) return null;
      var style = state.tokenize(stream, state);
      if (type == "comment") return style;
      state.lastType = type == "operator" && (content == "++" || content == "--") ? "incdec" : type;
      return parseJS(state, style, type, content, stream);
    },

    indent: function(state, textAfter) {
      if (state.tokenize == tokenComment || state.tokenize == tokenQuasi) return CodeMirror.Pass;
      if (state.tokenize != tokenBase) return 0;
      var firstChar = textAfter && textAfter.charAt(0), lexical = state.lexical, top
      // Kludge to prevent 'maybelse' from blocking lexical scope pops
      if (!/^\s*else\b/.test(textAfter)) for (var i = state.cc.length - 1; i >= 0; --i) {
        var c = state.cc[i];
        if (c == poplex) lexical = lexical.prev;
        else if (c != maybeelse && c != popcontext) break;
      }
      while ((lexical.type == "stat" || lexical.type == "form") &&
             (firstChar == "}" || ((top = state.cc[state.cc.length - 1]) &&
                                   (top == maybeoperatorComma || top == maybeoperatorNoComma) &&
                                   !/^[,\.=+\-*:?[\(]/.test(textAfter))))
        lexical = lexical.prev;
      if (statementIndent && lexical.type == ")" && lexical.prev.type == "stat")
        lexical = lexical.prev;
      var type = lexical.type, closing = firstChar == type;

      if (type == "vardef") return lexical.indented + (state.lastType == "operator" || state.lastType == "," ? lexical.info.length + 1 : 0);
      else if (type == "form" && firstChar == "{") return lexical.indented;
      else if (type == "form") return lexical.indented + indentUnit;
      else if (type == "stat")
        return lexical.indented + (isContinuedStatement(state, textAfter) ? statementIndent || indentUnit : 0);
      else if (lexical.info == "switch" && !closing && parserConfig.doubleIndentSwitch != false)
        return lexical.indented + (/^(?:case|default)\b/.test(textAfter) ? indentUnit : 2 * indentUnit);
      else if (lexical.align) return lexical.column + (closing ? 0 : 1);
      else return lexical.indented + (closing ? 0 : indentUnit);
    },

    electricInput: /^\s*(?:case .*?:|default:|\{|\})$/,
    blockCommentStart: jsonMode ? null : "/*",
    blockCommentEnd: jsonMode ? null : "*/",
    blockCommentContinue: jsonMode ? null : " * ",
    lineComment: jsonMode ? null : "//",
    fold: "brace",
    closeBrackets: "()[]{}''\"\"``",

    helperType: jsonMode ? "json" : "javascript",
    jsonldMode: jsonldMode,
    jsonMode: jsonMode,

    expressionAllowed: expressionAllowed,

    skipExpression: function(state) {
      parseJS(state, "atom", "atom", "true", new CodeMirror.StringStream("", 2, null))
    }
  };
});

CodeMirror.registerHelper("wordChars", "javascript", /[\w$]/);

CodeMirror.defineMIME("text/javascript", "javascript");
CodeMirror.defineMIME("text/ecmascript", "javascript");
CodeMirror.defineMIME("application/javascript", "javascript");
CodeMirror.defineMIME("application/x-javascript", "javascript");
CodeMirror.defineMIME("application/ecmascript", "javascript");
CodeMirror.defineMIME("application/json", { name: "javascript", json: true });
CodeMirror.defineMIME("application/x-json", { name: "javascript", json: true });
CodeMirror.defineMIME("application/manifest+json", { name: "javascript", json: true })
CodeMirror.defineMIME("application/ld+json", { name: "javascript", jsonld: true });
CodeMirror.defineMIME("text/typescript", { name: "javascript", typescript: true });
CodeMirror.defineMIME("application/typescript", { name: "javascript", typescript: true });

});
// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: https://codemirror.net/LICENSE

(function(mod) {
  if (typeof exports == "object" && typeof module == "object") // CommonJS
    mod(require("../../lib/codemirror"));
  else if (typeof define == "function" && define.amd) // AMD
    define(["../../lib/codemirror"], mod);
  else // Plain browser env
    mod(CodeMirror);
})(function(CodeMirror) {
"use strict";

CodeMirror.defineMode("go", function(config) {
  var indentUnit = config.indentUnit;

  var keywords = {
    "break":true, "case":true, "chan":true, "const":true, "continue":true,
    "default":true, "defer":true, "else":true, "fallthrough":true, "for":true,
    "func":true, "go":true, "goto":true, "if":true, "import":true,
    "interface":true, "map":true, "package":true, "range":true, "return":true,
    "select":true, "struct":true, "switch":true, "type":true, "var":true,
    "bool":true, "byte":true, "complex64":true, "complex128":true,
    "float32":true, "float64":true, "int8":true, "int16":true, "int32":true,
    "int64":true, "string":true, "uint8":true, "uint16":true, "uint32":true,
    "uint64":true, "int":true, "uint":true, "uintptr":true, "error": true,
    "rune":true
  };

  var atoms = {
    "true":true, "false":true, "iota":true, "nil":true, "append":true,
    "cap":true, "close":true, "complex":true, "copy":true, "delete":true, "imag":true,
    "len":true, "make":true, "new":true, "panic":true, "print":true,
    "println":true, "real":true, "recover":true
  };

  var isOperatorChar = /[+\-*&^%:=<>!|\/]/;

  var curPunc;

  function tokenBase(stream, state) {
    var ch = stream.next();
    if (ch == '"' || ch == "'" || ch == "`") {
      state.tokenize = tokenString(ch);
      return state.tokenize(stream, state);
    }
    if (/[\d\.]/.test(ch)) {
      if (ch == ".") {
        stream.match(/^[0-9]+([eE][\-+]?[0-9]+)?/);
      } else if (ch == "0") {
        stream.match(/^[xX][0-9a-fA-F]+/) || stream.match(/^0[0-7]+/);
      } else {
        stream.match(/^[0-9]*\.?[0-9]*([eE][\-+]?[0-9]+)?/);
      }
      return "number";
    }
    if (/[\[\]{}\(\),;\:\.]/.test(ch)) {
      curPunc = ch;
      return null;
    }
    if (ch == "/") {
      if (stream.eat("*")) {
        state.tokenize = tokenComment;
        return tokenComment(stream, state);
      }
      if (stream.eat("/")) {
        stream.skipToEnd();
        return "comment";
      }
    }
    if (isOperatorChar.test(ch)) {
      stream.eatWhile(isOperatorChar);
      return "operator";
    }
    stream.eatWhile(/[\w\$_\xa1-\uffff]/);
    var cur = stream.current();
    if (keywords.propertyIsEnumerable(cur)) {
      if (cur == "case" || cur == "default") curPunc = "case";
      return "keyword";
    }
    if (atoms.propertyIsEnumerable(cur)) return "atom";
    return "variable";
  }

  function tokenString(quote) {
    return function(stream, state) {
      var escaped = false, next, end = false;
      while ((next = stream.next()) != null) {
        if (next == quote && !escaped) {end = true; break;}
        escaped = !escaped && quote != "`" && next == "\\";
      }
      if (end || !(escaped || quote == "`"))
        state.tokenize = tokenBase;
      return "string";
    };
  }

  function tokenComment(stream, state) {
    var maybeEnd = false, ch;
    while (ch = stream.next()) {
      if (ch == "/" && maybeEnd) {
        state.tokenize = tokenBase;
        break;
      }
      maybeEnd = (ch == "*");
    }
    return "comment";
  }

  function Context(indented, column, type, align, prev) {
    this.indented = indented;
    this.column = column;
    this.type = type;
    this.align = align;
    this.prev = prev;
  }
  function pushContext(state, col, type) {
    return state.context = new Context(state.indented, col, type, null, state.context);
  }
  function popContext(state) {
    if (!state.context.prev) return;
    var t = state.context.type;
    if (t == ")" || t == "]" || t == "}")
      state.indented = state.context.indented;
    return state.context = state.context.prev;
  }

  // Interface

  return {
    startState: function(basecolumn) {
      return {
        tokenize: null,
        context: new Context((basecolumn || 0) - indentUnit, 0, "top", false),
        indented: 0,
        startOfLine: true
      };
    },

    token: function(stream, state) {
      var ctx = state.context;
      if (stream.sol()) {
        if (ctx.align == null) ctx.align = false;
        state.indented = stream.indentation();
        state.startOfLine = true;
        if (ctx.type == "case") ctx.type = "}";
      }
      if (stream.eatSpace()) return null;
      curPunc = null;
      var style = (state.tokenize || tokenBase)(stream, state);
      if (style == "comment") return style;
      if (ctx.align == null) ctx.align = true;

      if (curPunc == "{") pushContext(state, stream.column(), "}");
      else if (curPunc == "[") pushContext(state, stream.column(), "]");
      else if (curPunc == "(") pushContext(state, stream.column(), ")");
      else if (curPunc == "case") ctx.type = "case";
      else if (curPunc == "}" && ctx.type == "}") popContext(state);
      else if (curPunc == ctx.type) popContext(state);
      state.startOfLine = false;
      return style;
    },

    indent: function(state, textAfter) {
      if (state.tokenize != tokenBase && state.tokenize != null) return CodeMirror.Pass;
      var ctx = state.context, firstChar = textAfter && textAfter.charAt(0);
      if (ctx.type == "case" && /^(?:case|default)\b/.test(textAfter)) {
        state.context.type = "}";
        return ctx.indented;
      }
      var closing = firstChar == ctx.type;
      if (ctx.align) return ctx.column + (closing ? 0 : 1);
      else return ctx.indented + (closing ? 0 : indentUnit);
    },

    electricChars: "{}):",
    closeBrackets: "()[]{}''\"\"``",
    fold: "brace",
    blockCommentStart: "/*",
    blockCommentEnd: "*/",
    lineComment: "//"
  };
});

CodeMirror.defineMIME("text/x-go", "go");

});
// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: https://codemirror.net/5/LICENSE

// declare global: DOMRect

(function(mod) {
  if (typeof exports == "object" && typeof module == "object") // CommonJS
    mod(require("../../lib/codemirror"));
  else if (typeof define == "function" && define.amd) // AMD
    define(["../../lib/codemirror"], mod);
  else // Plain browser env
    mod(CodeMirror);
})(function(CodeMirror) {
  "use strict";

  var HINT_ELEMENT_CLASS        = "CodeMirror-hint";
  var ACTIVE_HINT_ELEMENT_CLASS = "CodeMirror-hint-active";

  // This is the old interface, kept around for now to stay
  // backwards-compatible.
  CodeMirror.showHint = function(cm, getHints, options) {
    if (!getHints) return cm.showHint(options);
    if (options && options.async) getHints.async = true;
    var newOpts = {hint: getHints};
    if (options) for (var prop in options) newOpts[prop] = options[prop];
    return cm.showHint(newOpts);
  };

  CodeMirror.defineExtension("showHint", function(options) {
    options = parseOptions(this, this.getCursor("start"), options);
    var selections = this.listSelections()
    if (selections.length > 1) return;
    // By default, don't allow completion when something is selected.
    // A hint function can have a `supportsSelection` property to
    // indicate that it can handle selections.
    if (this.somethingSelected()) {
      if (!options.hint.supportsSelection) return;
      // Don't try with cross-line selections
      for (var i = 0; i < selections.length; i++)
        if (selections[i].head.line != selections[i].anchor.line) return;
    }

    if (this.state.completionActive) this.state.completionActive.close();
    var completion = this.state.completionActive = new Completion(this, options);
    if (!completion.options.hint) return;

    CodeMirror.signal(this, "startCompletion", this);
    completion.update(true);
  });

  CodeMirror.defineExtension("closeHint", function() {
    if (this.state.completionActive) this.state.completionActive.close()
  })

  function Completion(cm, options) {
    this.cm = cm;
    this.options = options;
    this.widget = null;
    this.debounce = 0;
    this.tick = 0;
    this.startPos = this.cm.getCursor("start");
    this.startLen = this.cm.getLine(this.startPos.line).length - this.cm.getSelection().length;

    if (this.options.updateOnCursorActivity) {
      var self = this;
      cm.on("cursorActivity", this.activityFunc = function() { self.cursorActivity(); });
    }
  }

  var requestAnimationFrame = window.requestAnimationFrame || function(fn) {
    return setTimeout(fn, 1000/60);
  };
  var cancelAnimationFrame = window.cancelAnimationFrame || clearTimeout;

  Completion.prototype = {
    close: function() {
      if (!this.active()) return;
      this.cm.state.completionActive = null;
      this.tick = null;
      if (this.options.updateOnCursorActivity) {
        this.cm.off("cursorActivity", this.activityFunc);
      }

      if (this.widget && this.data) CodeMirror.signal(this.data, "close");
      if (this.widget) this.widget.close();
      CodeMirror.signal(this.cm, "endCompletion", this.cm);
    },

    active: function() {
      return this.cm.state.completionActive == this;
    },

    pick: function(data, i) {
      var completion = data.list[i], self = this;
      this.cm.operation(function() {
        if (completion.hint)
          completion.hint(self.cm, data, completion);
        else
          self.cm.replaceRange(getText(completion), completion.from || data.from,
                               completion.to || data.to, "complete");
        CodeMirror.signal(data, "pick", completion);
        self.cm.scrollIntoView();
      });
      if (this.options.closeOnPick) {
        this.close();
      }
    },

    cursorActivity: function() {
      if (this.debounce) {
        cancelAnimationFrame(this.debounce);
        this.debounce = 0;
      }

      var identStart = this.startPos;
      if(this.data) {
        identStart = this.data.from;
      }

      var pos = this.cm.getCursor(), line = this.cm.getLine(pos.line);
      if (pos.line != this.startPos.line || line.length - pos.ch != this.startLen - this.startPos.ch ||
          pos.ch < identStart.ch || this.cm.somethingSelected() ||
          (!pos.ch || this.options.closeCharacters.test(line.charAt(pos.ch - 1)))) {
        this.close();
      } else {
        var self = this;
        this.debounce = requestAnimationFrame(function() {self.update();});
        if (this.widget) this.widget.disable();
      }
    },

    update: function(first) {
      if (this.tick == null) return
      var self = this, myTick = ++this.tick
      fetchHints(this.options.hint, this.cm, this.options, function(data) {
        if (self.tick == myTick) self.finishUpdate(data, first)
      })
    },

    finishUpdate: function(data, first) {
      if (this.data) CodeMirror.signal(this.data, "update");

      var picked = (this.widget && this.widget.picked) || (first && this.options.completeSingle);
      if (this.widget) this.widget.close();

      this.data = data;

      if (data && data.list.length) {
        if (picked && data.list.length == 1) {
          this.pick(data, 0);
        } else {
          this.widget = new Widget(this, data);
          CodeMirror.signal(data, "shown");
        }
      }
    }
  };

  function parseOptions(cm, pos, options) {
    var editor = cm.options.hintOptions;
    var out = {};
    for (var prop in defaultOptions) out[prop] = defaultOptions[prop];
    if (editor) for (var prop in editor)
      if (editor[prop] !== undefined) out[prop] = editor[prop];
    if (options) for (var prop in options)
      if (options[prop] !== undefined) out[prop] = options[prop];
    if (out.hint.resolve) out.hint = out.hint.resolve(cm, pos)
    return out;
  }

  function getText(completion) {
    if (typeof completion == "string") return completion;
    else return completion.text;
  }

  function buildKeyMap(completion, handle) {
    var baseMap = {
      Up: function() {handle.moveFocus(-1);},
      Down: function() {handle.moveFocus(1);},
      PageUp: function() {handle.moveFocus(-handle.menuSize() + 1, true);},
      PageDown: function() {handle.moveFocus(handle.menuSize() - 1, true);},
      Home: function() {handle.setFocus(0);},
      End: function() {handle.setFocus(handle.length - 1);},
      Enter: handle.pick,
      Tab: handle.pick,
      Esc: handle.close
    };

    var mac = /Mac/.test(navigator.platform);

    if (mac) {
      baseMap["Ctrl-P"] = function() {handle.moveFocus(-1);};
      baseMap["Ctrl-N"] = function() {handle.moveFocus(1);};
    }

    var custom = completion.options.customKeys;
    var ourMap = custom ? {} : baseMap;
    function addBinding(key, val) {
      var bound;
      if (typeof val != "string")
        bound = function(cm) { return val(cm, handle); };
      // This mechanism is deprecated
      else if (baseMap.hasOwnProperty(val))
        bound = baseMap[val];
      else
        bound = val;
      ourMap[key] = bound;
    }
    if (custom)
      for (var key in custom) if (custom.hasOwnProperty(key))
        addBinding(key, custom[key]);
    var extra = completion.options.extraKeys;
    if (extra)
      for (var key in extra) if (extra.hasOwnProperty(key))
        addBinding(key, extra[key]);
    return ourMap;
  }

  function getHintElement(hintsElement, el) {
    while (el && el != hintsElement) {
      if (el.nodeName.toUpperCase() === "LI" && el.parentNode == hintsElement) return el;
      el = el.parentNode;
    }
  }

  function Widget(completion, data) {
    this.id = "cm-complete-" + Math.floor(Math.random(1e6))
    this.completion = completion;
    this.data = data;
    this.picked = false;
    var widget = this, cm = completion.cm;
    var ownerDocument = cm.getInputField().ownerDocument;
    var parentWindow = ownerDocument.defaultView || ownerDocument.parentWindow;

    var hints = this.hints = ownerDocument.createElement("ul");
    hints.setAttribute("role", "listbox")
    hints.setAttribute("aria-expanded", "true")
    hints.id = this.id
    var theme = completion.cm.options.theme;
    hints.className = "CodeMirror-hints " + theme;
    this.selectedHint = data.selectedHint || 0;

    var completions = data.list;
    for (var i = 0; i < completions.length; ++i) {
      var elt = hints.appendChild(ownerDocument.createElement("li")), cur = completions[i];
      var className = HINT_ELEMENT_CLASS + (i != this.selectedHint ? "" : " " + ACTIVE_HINT_ELEMENT_CLASS);
      if (cur.className != null) className = cur.className + " " + className;
      elt.className = className;
      if (i == this.selectedHint) elt.setAttribute("aria-selected", "true")
      elt.id = this.id + "-" + i
      elt.setAttribute("role", "option")
      if (cur.render) cur.render(elt, data, cur);
      else elt.appendChild(ownerDocument.createTextNode(cur.displayText || getText(cur)));
      elt.hintId = i;
    }

    var container = completion.options.container || ownerDocument.body;
    var pos = cm.cursorCoords(completion.options.alignWithWord ? data.from : null);
    var left = pos.left, top = pos.bottom, below = true;
    var offsetLeft = 0, offsetTop = 0;
    if (container !== ownerDocument.body) {
      // We offset the cursor position because left and top are relative to the offsetParent's top left corner.
      var isContainerPositioned = ['absolute', 'relative', 'fixed'].indexOf(parentWindow.getComputedStyle(container).position) !== -1;
      var offsetParent = isContainerPositioned ? container : container.offsetParent;
      var offsetParentPosition = offsetParent.getBoundingClientRect();
      var bodyPosition = ownerDocument.body.getBoundingClientRect();
      offsetLeft = (offsetParentPosition.left - bodyPosition.left - offsetParent.scrollLeft);
      offsetTop = (offsetParentPosition.top - bodyPosition.top - offsetParent.scrollTop);
    }
    hints.style.left = (left - offsetLeft) + "px";
    hints.style.top = (top - offsetTop) + "px";

    // If we're at the edge of the screen, then we want the menu to appear on the left of the cursor.
    var winW = parentWindow.innerWidth || Math.max(ownerDocument.body.offsetWidth, ownerDocument.documentElement.offsetWidth);
    var winH = parentWindow.innerHeight || Math.max(ownerDocument.body.offsetHeight, ownerDocument.documentElement.offsetHeight);
    container.appendChild(hints);
    cm.getInputField().setAttribute("aria-autocomplete", "list")
    cm.getInputField().setAttribute("aria-owns", this.id)
    cm.getInputField().setAttribute("aria-activedescendant", this.id + "-" + this.selectedHint)

    var box = completion.options.moveOnOverlap ? hints.getBoundingClientRect() : new DOMRect();
    var scrolls = completion.options.paddingForScrollbar ? hints.scrollHeight > hints.clientHeight + 1 : false;

    // Compute in the timeout to avoid reflow on init
    var startScroll;
    setTimeout(function() { startScroll = cm.getScrollInfo(); });

    var overlapY = box.bottom - winH;
    if (overlapY > 0) {
      var height = box.bottom - box.top, curTop = pos.top - (pos.bottom - box.top);
      if (curTop - height > 0) { // Fits above cursor
        hints.style.top = (top = pos.top - height - offsetTop) + "px";
        below = false;
      } else if (height > winH) {
        hints.style.height = (winH - 5) + "px";
        hints.style.top = (top = pos.bottom - box.top - offsetTop) + "px";
        var cursor = cm.getCursor();
        if (data.from.ch != cursor.ch) {
          pos = cm.cursorCoords(cursor);
          hints.style.left = (left = pos.left - offsetLeft) + "px";
          box = hints.getBoundingClientRect();
        }
      }
    }
    var overlapX = box.right - winW;
    if (scrolls) overlapX += cm.display.nativeBarWidth;
    if (overlapX > 0) {
      if (box.right - box.left > winW) {
        hints.style.width = (winW - 5) + "px";
        overlapX -= (box.right - box.left) - winW;
      }
      hints.style.left = (left = Math.max(pos.left - overlapX - offsetLeft, 0)) + "px";
    }
    if (scrolls) for (var node = hints.firstChild; node; node = node.nextSibling)
      node.style.paddingRight = cm.display.nativeBarWidth + "px"

    cm.addKeyMap(this.keyMap = buildKeyMap(completion, {
      moveFocus: function(n, avoidWrap) { widget.changeActive(widget.selectedHint + n, avoidWrap); },
      setFocus: function(n) { widget.changeActive(n); },
      menuSize: function() { return widget.screenAmount(); },
      length: completions.length,
      close: function() { completion.close(); },
      pick: function() { widget.pick(); },
      data: data
    }));

    if (completion.options.closeOnUnfocus) {
      var closingOnBlur;
      cm.on("blur", this.onBlur = function() { closingOnBlur = setTimeout(function() { completion.close(); }, 100); });
      cm.on("focus", this.onFocus = function() { clearTimeout(closingOnBlur); });
    }

    cm.on("scroll", this.onScroll = function() {
      var curScroll = cm.getScrollInfo(), editor = cm.getWrapperElement().getBoundingClientRect();
      if (!startScroll) startScroll = cm.getScrollInfo();
      var newTop = top + startScroll.top - curScroll.top;
      var point = newTop - (parentWindow.pageYOffset || (ownerDocument.documentElement || ownerDocument.body).scrollTop);
      if (!below) point += hints.offsetHeight;
      if (point <= editor.top || point >= editor.bottom) return completion.close();
      hints.style.top = newTop + "px";
      hints.style.left = (left + startScroll.left - curScroll.left) + "px";
    });

    CodeMirror.on(hints, "dblclick", function(e) {
      var t = getHintElement(hints, e.target || e.srcElement);
      if (t && t.hintId != null) {widget.changeActive(t.hintId); widget.pick();}
    });

    CodeMirror.on(hints, "click", function(e) {
      var t = getHintElement(hints, e.target || e.srcElement);
      if (t && t.hintId != null) {
        widget.changeActive(t.hintId);
        if (completion.options.completeOnSingleClick) widget.pick();
      }
    });

    CodeMirror.on(hints, "mousedown", function() {
      setTimeout(function(){cm.focus();}, 20);
    });

    // The first hint doesn't need to be scrolled to on init
    var selectedHintRange = this.getSelectedHintRange();
    if (selectedHintRange.from !== 0 || selectedHintRange.to !== 0) {
      this.scrollToActive();
    }

    CodeMirror.signal(data, "select", completions[this.selectedHint], hints.childNodes[this.selectedHint]);
    return true;
  }

  Widget.prototype = {
    close: function() {
      if (this.completion.widget != this) return;
      this.completion.widget = null;
      if (this.hints.parentNode) this.hints.parentNode.removeChild(this.hints);
      this.completion.cm.removeKeyMap(this.keyMap);
      var input = this.completion.cm.getInputField()
      input.removeAttribute("aria-activedescendant")
      input.removeAttribute("aria-owns")

      var cm = this.completion.cm;
      if (this.completion.options.closeOnUnfocus) {
        cm.off("blur", this.onBlur);
        cm.off("focus", this.onFocus);
      }
      cm.off("scroll", this.onScroll);
    },

    disable: function() {
      this.completion.cm.removeKeyMap(this.keyMap);
      var widget = this;
      this.keyMap = {Enter: function() { widget.picked = true; }};
      this.completion.cm.addKeyMap(this.keyMap);
    },

    pick: function() {
      this.completion.pick(this.data, this.selectedHint);
    },

    changeActive: function(i, avoidWrap) {
      if (i >= this.data.list.length)
        i = avoidWrap ? this.data.list.length - 1 : 0;
      else if (i < 0)
        i = avoidWrap ? 0  : this.data.list.length - 1;
      if (this.selectedHint == i) return;
      var node = this.hints.childNodes[this.selectedHint];
      if (node) {
        node.className = node.className.replace(" " + ACTIVE_HINT_ELEMENT_CLASS, "");
        node.removeAttribute("aria-selected")
      }
      node = this.hints.childNodes[this.selectedHint = i];
      node.className += " " + ACTIVE_HINT_ELEMENT_CLASS;
      node.setAttribute("aria-selected", "true")
      this.completion.cm.getInputField().setAttribute("aria-activedescendant", node.id)
      this.scrollToActive()
      CodeMirror.signal(this.data, "select", this.data.list[this.selectedHint], node);
    },

    scrollToActive: function() {
      var selectedHintRange = this.getSelectedHintRange();
      var node1 = this.hints.childNodes[selectedHintRange.from];
      var node2 = this.hints.childNodes[selectedHintRange.to];
      var firstNode = this.hints.firstChild;
      if (node1.offsetTop < this.hints.scrollTop)
        this.hints.scrollTop = node1.offsetTop - firstNode.offsetTop;
      else if (node2.offsetTop + node2.offsetHeight > this.hints.scrollTop + this.hints.clientHeight)
        this.hints.scrollTop = node2.offsetTop + node2.offsetHeight - this.hints.clientHeight + firstNode.offsetTop;
    },

    screenAmount: function() {
      return Math.floor(this.hints.clientHeight / this.hints.firstChild.offsetHeight) || 1;
    },

    getSelectedHintRange: function() {
      var margin = this.completion.options.scrollMargin || 0;
      return {
        from: Math.max(0, this.selectedHint - margin),
        to: Math.min(this.data.list.length - 1, this.selectedHint + margin),
      };
    }
  };

  function applicableHelpers(cm, helpers) {
    if (!cm.somethingSelected()) return helpers
    var result = []
    for (var i = 0; i < helpers.length; i++)
      if (helpers[i].supportsSelection) result.push(helpers[i])
    return result
  }

  function fetchHints(hint, cm, options, callback) {
    if (hint.async) {
      hint(cm, callback, options)
    } else {
      var result = hint(cm, options)
      if (result && result.then) result.then(callback)
      else callback(result)
    }
  }

  function resolveAutoHints(cm, pos) {
    var helpers = cm.getHelpers(pos, "hint"), words
    if (helpers.length) {
      var resolved = function(cm, callback, options) {
        var app = applicableHelpers(cm, helpers);
        function run(i) {
          if (i == app.length) return callback(null)
          fetchHints(app[i], cm, options, function(result) {
            if (result && result.list.length > 0) callback(result)
            else run(i + 1)
          })
        }
        run(0)
      }
      resolved.async = true
      resolved.supportsSelection = true
      return resolved
    } else if (words = cm.getHelper(cm.getCursor(), "hintWords")) {
      return function(cm) { return CodeMirror.hint.fromList(cm, {words: words}) }
    } else if (CodeMirror.hint.anyword) {
      return function(cm, options) { return CodeMirror.hint.anyword(cm, options) }
    } else {
      return function() {}
    }
  }

  CodeMirror.registerHelper("hint", "auto", {
    resolve: resolveAutoHints
  });

  CodeMirror.registerHelper("hint", "fromList", function(cm, options) {
    var cur = cm.getCursor(), token = cm.getTokenAt(cur)
    var term, from = CodeMirror.Pos(cur.line, token.start), to = cur
    if (token.start < cur.ch && /\w/.test(token.string.charAt(cur.ch - token.start - 1))) {
      term = token.string.substr(0, cur.ch - token.start)
    } else {
      term = ""
      from = cur
    }
    var found = [];
    for (var i = 0; i < options.words.length; i++) {
      var word = options.words[i];
      if (word.slice(0, term.length) == term)
        found.push(word);
    }

    if (found.length) return {list: found, from: from, to: to};
  });

  CodeMirror.commands.autocomplete = CodeMirror.showHint;

  var defaultOptions = {
    hint: CodeMirror.hint.auto,
    completeSingle: true,
    alignWithWord: true,
    closeCharacters: /[\s()\[\]{};:>,]/,
    closeOnPick: true,
    closeOnUnfocus: true,
    updateOnCursorActivity: true,
    completeOnSingleClick: true,
    container: null,
    customKeys: null,
    extraKeys: null,
    paddingForScrollbar: true,
    moveOnOverlap: true,
  };

  CodeMirror.defineOption("hintOptions", null);
});
