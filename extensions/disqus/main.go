package disqus

import (
	"flag"
	"fmt"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

const tmpl = `
<div id="disqus_thread"></div>
<script>
    /**
    *  RECOMMENDED CONFIGURATION VARIABLES: EDIT AND UNCOMMENT THE SECTION BELOW TO INSERT DYNAMIC VALUES FROM YOUR PLATFORM OR CMS.
    *  LEARN WHY DEFINING THESE VARIABLES IS IMPORTANT: https://disqus.com/admin/universalcode/#configuration-variables    */
    /*
    var disqus_config = function () {
	  this.page.identifier = "%s";
    };
    */
    (function() { // DON'T EDIT BELOW THIS LINE
    var d = document, s = d.createElement('script');
    s.src = 'https://%s/embed.js';
    s.setAttribute('data-timestamp', +new Date());
    (d.head || d.body).appendChild(s);
    })();
</script>`

var domain string

func init() {
	flag.StringVar(&domain, "disqus", "", "Disqus domain name for example: xlog-emadelsaid.disqus.com")
	RegisterExtension(Disqus{})
}

type Disqus struct{}

func (Disqus) Name() string { return "disqus" }
func (Disqus) Init() {
	RegisterWidget(WidgetAfterView, 2, widget)
}

func widget(p Page) template.HTML {
	if domain == "" {
		return ""
	}

	script := fmt.Sprintf(tmpl, template.JSEscapeString(p.Name()), domain)
	return template.HTML(script)
}
