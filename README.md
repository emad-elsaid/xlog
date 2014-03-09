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

## TODO

* code highlight
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

## license

MIT license personal and commercial usage.

## Credits

Emad Elsaid. <mailto:blazeeboy@gmail.com>