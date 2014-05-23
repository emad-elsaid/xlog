json.array!(@posts) do |post|
  json.extract! post, :id, :title, :user_id, :created_at, :updated_at
  json.url post_link_url(post, format: :json)
end
