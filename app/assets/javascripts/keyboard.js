import $ from 'jquery'

let keyboard_handler = function(e) {
  let ctrlKey

  if ( navigator.platform.match("Mac") ) {
    ctrlKey = e.metaKey
  } else {
    ctrlKey = e.ctrlKey
  }

  if ( ctrlKey ) {
    let character = String.fromCharCode(e.keyCode).toLowerCase()
    let handler = 'ctrl+' + character
    let elements = $('[data-key="' + handler + '"]')

    if ( elements.length > 0 ) {
      e.preventDefault()

      // execute elements
      elements.each( function() {
        if ($(this).is('input[type=text]')) {
          $(this).focus()
        }

        if ( $(this).is('input[type=submit]') || $(this).is('input[type=button]') || $(this).is('input[type=reset]') || $(this).is('button') || $(this).is('a') ) {
          $(this).trigger('click')
        }

        if ( $(this).is('a') ) {
          $(this).trigger('click')
          document.location.href = $(this).attr('href')
        }
      })
    }
  }
}

document.addEventListener('keydown', keyboard_handler, false)
