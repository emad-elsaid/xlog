Xlog::Application.routes.draw do
  resources :posts

  devise_for :users
  root 'posts#index'
end
