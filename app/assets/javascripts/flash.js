import $ from 'jquery'

$(function(){
  $(document).on('click', '.flash, #error_explanation', function(){
    $(this).fadeOut()
  })
})
