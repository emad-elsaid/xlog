$(function(){
	$(document).on('page:before-change', function(){
		$('.loader').fadeIn();
	});

	$(document).on('page:change', function(){
		$('.loader').hide();
	});
})