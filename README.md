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
```

## Demo

you can check the demo on my koding vm from here 
http://blazeeboy.kd.io:3000/

```
email: admin@example.com
password: password
```

now it is ready to run your server of choice.

## Project Status

project still under development, DO NOT USE IT IN PRODUCTION
i still didn't add previlages for production. 

## license

MIT license personal and commercial usage.

## Credits

Emad Elsaid. <mailto:blazeeboy@gmail.com>