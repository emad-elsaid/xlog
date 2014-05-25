//= require jquery
//= require jquery_ujs
//= require foundation
//= require turbolinks
//= require_tree .

var onPageUpdate = function(){
	$(document).foundation(); 
	$('textarea[data-expanding]').autosize();
}
$(onPageUpdate);
document.addEventListener("page:update", onPageUpdate);

// cache pages 
Turbolinks.enableTransitionCache();
