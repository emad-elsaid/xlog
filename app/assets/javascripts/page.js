import $ from "jquery"
import autosize from "autosize"

let onPageUpdate = function() {
  $(document).foundation()
  autosize($('textarea[data-expanding]'))
}

$(onPageUpdate)
document.addEventListener('turbolinks:render', onPageUpdate)
