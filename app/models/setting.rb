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

  # check if value of a Setting value is truthy 
  def self.yes?(name) 
    setting = self.find_or_initialize_by(name: name)
    ['yes','y', 't', true].include? setting.value
  end

  # check if value of a setting is falsy
  def self.no?(name)
    !self.yes? name
  end

  def self.value(name)
    self.find_or_initialize_by(name: name).value
  end

  # gets Hash of {name: value, ...} and update them in database
  def self.set(params)
    params.each do |name, value|
      self.find_or_initialize_by(name: name).update(value: value)
    end
  end

end
