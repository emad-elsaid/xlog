$ ->
	$(document).on 'page:before-change', ->
		$('.loader').fadeIn()

	$(document).on 'page:change', ->
		$('.loader').hide()