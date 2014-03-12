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
## Available tasks

A Task to copy Gmoji gem images (smile faces) to public, this should sync them with gem latest
```bash
$ rake emoji
``` 

## Project Status

we can say that we're still in the alpha phase, it could be used for blogging but
we didn't test it enough to make sure it is stable enough for production use.

## Current Features

* simple theme
* responsive and mobile ready theme
* Login/Logout/Forgot password powered by devise
* Code highlight
* posts written in github Flavored Markdown
* Ruby on rails based, so extending won't be a problem
* Permenant links for posts and SEO friendly URL `domain.com/post-title-here`
* emoticons like github [sheatsheet](http://www.emoji-cheat-sheet.com/)
* replace Facebook/Twitter linksby embeded version of the post
* replaces youtube video link by a video player
* replaces Github Gist link by an embeded version of the gist


## TODO

* replace links with embeded versions for 
	* scribed
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
* catch 404 exceptions with 404 page ( make normal post using seeds and create a settings variable and make it the default value)
* add support for MathML
* enhance embeded links style a little bit further
* replace any link by an embeded preview like facebook links preview (this is gonna be hard but we can inspect discourse for it)

## license

MIT license personal and commercial usage.

## Credits

Emad Elsaid. <mailto:blazeeboy@gmail.com>