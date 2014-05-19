# == Schema Information
#
# Table name: settings
#
#  id         :integer          not null, primary key
#  name       :string(255)
#  value      :text
#  created_at :datetime
#  updated_at :datetime
#

class Setting < ActiveRecord::Base
	validates_uniqueness_of :name
  validates_presence_of :name
end
