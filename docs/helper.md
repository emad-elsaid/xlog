A helper function is a function that is used in an html template. the output of the function gets printed to the template. a helper function can take one or more parameters.

Helper functions are Go `html/template` concept it's not introduced by xlog. an example can be found in html/template [documentation](https://pkg.go.dev/html/template#example-Template-Helpers). 

Extensions can define their own helpers to be used by any template using [`RegisterHelper`](https://pkg.go.dev/github.com/emad-elsaid/xlog#RegisterHelper) function. Registering a new helper has to be in the extension `init` function to be find at the time of parsing the templates.