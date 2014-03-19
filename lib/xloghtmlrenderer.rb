class XlogHTMLRenderer < Redcarpet::Render::HTML

  def embed_youtube url, path, query
    v = query['v'].first || path[1..-1]
    return nil if v.blank?
    "<div class=\"embed_youtube flex-video\"><iframe src=\"//www.youtube.com/embed/#{v}\" frameborder=\"0\" allowfullscreen></iframe></div>"
  end

  def embed_facebook url, path, query
    '<div class="embed_facebook"><div id="fb-root"></div> <script>(function(d, s, id) { var js, fjs = d.getElementsByTagName(s)[0]; if (d.getElementById(id)) return; js = d.createElement(s); js.id = id; js.src = "//connect.facebook.net/en_US/all.js#xfbml=1"; fjs.parentNode.insertBefore(js, fjs); }(document, \'script\', \'facebook-jssdk\'));</script><div class="fb-post" data-href="'+url+'" ></div></div>'
  end

  def embed_twitter url, path, query
    '<div class="embed_twitter"><blockquote class="twitter-tweet" lang="en"><a href="'+url+'"></a></blockquote> <script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script></div>'
  end

  def embed_gist url, path, query
    '<div class="embed_gist"><script src="'+url+'.js"></script></div>'
  end

  def embed_image url, path, query
    '<div class="embed_image"><img src="'+url+'"/></div>'
  end

  def preprocess(full_document)
    full_document.gsub(/^http(s)?:\/\/\S*/) do |url|
      u = URI.parse(url)
      query = u.query.nil? ? {} : CGI.parse(u.query)
      path = u.path
      host = u.host

      if host.end_with? 'youtube.com'
        embed_youtube url, path, query
      elsif host.end_with? 'facebook.com' and (path.include?('posts') or path.include?('photo'))
        embed_facebook url, path, query
      elsif host.end_with? 'twitter.com'
        embed_twitter url, path, query
      elsif host.end_with? 'gist.github.com'
        embed_gist url, path, query
      elsif url.end_with? '.jpg' or url.end_with? '.png' or url.end_with? '.gif'
        embed_image url, path, query
      else
        url
      end
    end
  end

  def block_code(code, language)
    Pygmentize.process(code, language)
  end
end