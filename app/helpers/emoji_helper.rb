module EmojiHelper
 def emojify(content)
    content.gsub(/:([a-z0-9\+\-_]+):/) do |match|
      if Emoji.names.include?($1)
        "<img title=\"#{$1}\" alt=\"#{$1}\" height=\"20\" src=\"/images/emoji/#{$1}.png\" style=\"vertical-align:middle\" width=\"20\" />"
      else
        match
      end
    end if content.present?
  end
end