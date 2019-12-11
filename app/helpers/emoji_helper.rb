module EmojiHelper
  def emojify(content)
    content.to_str.gsub(/:([\w+-]+):/) do |match|
      if emoji = Emoji.find_by_alias($1)
        emoji.raw
      else
        match
      end
    end.html_safe if content.present?
  end
end
