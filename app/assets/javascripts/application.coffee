#= require jquery
#= require jquery_ujs
#= require foundation
#= require turbolinks
#= require_tree .

$ ->
  onPageUpdate = ->
    $(document).foundation() 
    $('textarea[data-expanding]').autosize()

document.addEventListener 'page:update', onPageUpdate

# cache pages 
Turbolinks.enableTransitionCache()