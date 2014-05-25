keyboard_handler = (e)->
  ctrlKey = if navigator.platform.match("Mac") then e.metaKey else e.ctrlKey
  if ctrlKey
    character = String.fromCharCode(e.keyCode).toLowerCase()
    handler = 'ctrl+' + character
    elements = $('[data-key="' + handler + '"]')
    if elements.size() > 0
      e.preventDefault()

      # execute elements
      elements.each ->
        
        $(@).focus() if $(@).is('input[type=text]')

        $(@).trigger('click') if  $(@).is('input[type=submit]') || 
                                  $(@).is('input[type=button]') ||
                                  $(@).is('input[type=reset]') ||
                                  $(@).is('button') ||
                                  $(@).is('a')
          
        if $(@).is('a')
          $(@).trigger('click') 
          document.location.href = $(@).attr('href')

document.addEventListener 'keydown', keyboard_handler, false