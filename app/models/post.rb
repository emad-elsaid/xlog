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

  validates :title, presence: true, length: { minimum: 3 }
  validates :body, presence: true, length: { minimum: 20 }

end
