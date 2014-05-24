document.addEventListener("keydown", function(e) {
  if (navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey){
    var character = String.fromCharCode(e.keyCode).toLowerCase();
    var handler = "ctrl+"+character;
    var elements = $('[data-key="'+handler+'"]');
    if( elements.size() > 0){
      e.preventDefault();

      // execute elements
      elements.each(function(){
        var $this = $(this);
        
        if($this.is('input[type=text]')){
          $this.focus();

        }else if( $this.is('input[type=submit]') || 
                  $this.is('input[type=button]') ||
                  $this.is('input[type=reset]') ||
                  $this.is('button')){
          $this.trigger('click');

        }else if( $this.is('a') ){
          $this.trigger('click');
          document.location.href = $this.attr('href');
        }
      });

    }
  }
}, false);