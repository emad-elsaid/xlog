* Xlog doesn't require custom structure for your markdown files
* Default main file name is `index.md` and can be overriden with `--index` flag

# Create empty directory

Create a new empty directory and `cd` into it or simple navigate to existing directory which has markdown files

```shell
mkdir myblog
cd myblog
```

# Run Xlog

Assuming you already went through one of the Installation methods. `xlog` should be in your **PATH**. 

```shell
xlog
```

# Running on different port

The previous command starts a server on port **3000** if you want to specify the port you can do so using `--bind` flag

```shell
xlog --bind 127.0.0.1:4000
```

This will run the server on port **4000** instead of **3000** 

# Using different index page

Xlog assumes the main page is **index.md** if you're working in an existing github repository for example you may need to specify **README.md** as you index page as follows

```shell
xlog --index README
```

Notice that specifying the index page doesn't need the extension `.md`.

# Open your new site

Now you can navigate to http://localhost:3000 in your browser to start browsing the markdown files. if it's a new directory you'll be redirected to the editor to write your first page. 

Note that "Ctrl+S" will save the page if you're in the edit page. and also navigate to the edit page if you're viewing a page. so you can switch between edit and view mode with `Ctrl+S` and on MacOS `⌘+S`

# Generating a static site

You can generate HTML files from your markdown files using `--build` flag

```shell
xlog --build .
```

Which will convert all of your markdown files to HTML files in the current directory. 

You can specify a destination for the HTML output.

```shell
xlog --build /destination/directory/path
```

# Integration with Github pages

If your markdown is hosted as Gituhub repository. You can use github workflows to download and execute xlog to generate HTML pages and host it with github pages. an examples can be found here:

- [Emad Elsaid Blog](https://github.com/emad-elsaid/emad-elsaid.github.io/blob/master/.github/workflows/xlog.yml)
- [Xlog documentation](https://github.com/emad-elsaid/xlog/blob/master/.github/workflows/xlog.yml)