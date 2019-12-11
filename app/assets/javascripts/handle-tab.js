import $ from 'jquery'

// credit to :
// http://stackoverflow.com/questions/6140632/how-to-handle-tab-in-textarea
$(function(){
  $(document).on('keydown', 'textarea', function(e) {
    if( e.keyCode == 9 ) { // tab was pressed
      // get caret position/selection
      let start = this.selectionStart
      let end = this.selectionEnd
      let value = $(this).val()

      // set textarea value to: text before caret + tab + text after caret
      $(this).val(value.substring(0, start) + "\t" + value.substring(end))

      // put caret at right position again (add one for the tab)
      this.selectionStart = this.selectionEnd = start + 1

      // prevent the focus lose
      e.preventDefault()
    }
  })
})
