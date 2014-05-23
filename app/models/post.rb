# == Schema Information
#
# Table name: posts
#
#  id         :integer          not null, primary key
#  title      :string(255)
#  body       :text
#  user_id    :integer
#  created_at :datetime
#  updated_at :datetime
#  permalink  :string(255)
#

require 'xloghtmlrenderer'

class Post < ActiveRecord::Base
  include EmojiHelper
  
  has_permalink
  self.per_page = 20
  belongs_to :user

  validates :title, presence: true, uniqueness: true, length: { minimum: 3 }
  validates :body, presence: true, length: { minimum: 20 }

  # convert the body written by user to HTML
  def body_html
  	renderer = Redcarpet::Markdown.new(XlogHTMLRenderer, 
                            fenced_code_blocks: true,
                            disable_indented_code_blocks: true,
                            strikethrough: true,
                            superscript: true,
                            underline: true,
                            autolink: true )
  	emojify renderer.render(body)
  end

end