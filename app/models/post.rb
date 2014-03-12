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

class Post < ActiveRecord::Base
  has_permalink
  belongs_to :user

  validates :title, presence: true, uniqueness: true, length: { minimum: 3 }
  validates :body, presence: true, length: { minimum: 20 }

  # convert the body written by user to HTML
  def body_html
  	renderer = Redcarpet::Markdown.new(PygmentizeHTML, 
                            fenced_code_blocks: true,
                            disable_indented_code_blocks: true,
                            strikethrough: true,
                            superscript: true,
                            underline: true,
                            autolink: true )
  	renderer.render(body)
  end

end

# i think putting two classes in one file
# is a bad practice, but till we find
# a suitable place for it just keep it here
class PygmentizeHTML < Redcarpet::Render::HTML
  def block_code(code, language)
    Pygmentize.process(code, language)
  end
end