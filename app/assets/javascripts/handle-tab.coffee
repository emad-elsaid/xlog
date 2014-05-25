# credit to : 
# http://stackoverflow.com/questions/6140632/how-to-handle-tab-in-textarea
$ ->
  $(document).on 'keydown', 'textarea', (e)->
    if e.keyCode == 9 # tab was pressed
      # get caret position/selection
      start = this.selectionStart
      end = this.selectionEnd
      value = $(@).val()

      # set textarea value to: text before caret + tab + text after caret
      $(@).val(value.substring(0, start) + "\t" + value.substring(end))

      # put caret at right position again (add one for the tab)
      this.selectionStart = this.selectionEnd = start + 1;

      # prevent the focus lose
      e.preventDefault()
