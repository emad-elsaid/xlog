Xlog::Application.routes.draw do
  resources :posts

  devise_for :users
  get ':id' => 'posts#show', as: 'post_link'
  root 'posts#index'
end
