module ApplicationHelper

  # 
  # page body tag class attribute values helpers
  # 
  attr_accessor :body_classes

  # add new class to page body
  def add_body_class(*class_names)
    @body_classes ||= []
    @body_classes += class_names
  end

  # delete a class from page body
  def delete_body_class(*class_names)
    @body_classes ||= []
    @body_classes -= class_names
  end

  # get body class attribute total value
  def body_class
    @body_classes ||= []
    @body_classes.map(&:to_s).join ' '
  end

end
