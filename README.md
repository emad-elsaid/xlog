# Xlog

XLog is a simple blogging system for hackers, it will utilize rdoc syntax for posts content, simple file manager, and medium.com like style of interaction, also i'll experiment some new UX ideas. 

## why ?

nothing fancy, i just wanted to build a system in my mind, simple, uses markdown, rdoc ..etc, and some ux ideas i got in my mind, so why not, lets build it :)

## install

```bash
  bundle install
  rake db:setup
  rails s
```
## install on ubuntu

if you want to install on [koding](http://www.koding.com), open terminal and follow the following steps

```bash
git clone https://github.com/blazeeboy/xlog.git
cd xlog
sudo apt-get install libmagickwand-dev imagemagick
gem install minitest -v '4.7.5'
gem install json -v '1.8.1'
gem install rmagick -v '2.13.2'
export PKG_CONFIG_PATH="/opt/local/lib/pkgconfig:$PKG_CONFIG_PATH"
bundle install
rake db:setup
rake assets:precompile RAILS_ENV=production
rails server -d -e production
```
now it is ready to run your server of choice.

## Default credentials

```
email: admin@example.com
password: password
```

## Project Status

project still under development, DO NOT USE IT IN PRODUCTION
i still didn't add previlages for production. 

## Current Features

* simple theme
* responsive and mobile ready theme
* Login/Logout/Forgot password powered by devise
* Code highlight
* posts written in github Flavored Markdown
* Ruby on rails based, so extending won't be a problem
* Permenant links for posts and SEO friendly URL `domain.com/post-title-here`


## TODO

* replace links with embeded versions for 
	* FB
	* Twitter
	* Youtube
	* scribed
	* gist
	* soundcloud
	* ...etc
* blocks system to add more complex post contents such as:
	* Fiddles/Codepen style
* get rid of some gems that is not so important
* File manager
* Actions panel like sublime/atom panel and action/extending system
* keyboard shortcuts (ctrl+p for actions panel for example)
* enhance design typography
* add caching solution for views rendering (i think ryan bates made many screen casts about this)
* try to save posts to folder instead of sqlite db
* change views to haml
* emoticons like github

## license

MIT license personal and commercial usage.

## Credits

Emad Elsaid. <mailto:blazeeboy@gmail.com>