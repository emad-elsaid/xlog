![](/docs/public/website.png)

#tutorial

* Xlog doesn't require custom structure for your markdown files
* Default main file name is `index.md` and can be overriden with `--index` flag

# Create empty directory

Create a new empty directory and `cd` into it or simple navigate to existing directory which has markdown files

```shell
mkdir myblog
cd myblog
```

# Run Xlog

Assuming you already went through one of the Installation methods. `xlog` should be in your **PATH**. Simply executing it in current directory starts an HTTP server on port 3000

```shell
xlog
```

# Running on a different port

The previous command starts a server on port **3000** if you want to specify the port you can do so using `--bind` flag

```shell
xlog --bind 127.0.0.1:4000
```

This will run the server on port **4000** instead of **3000** 

# Using a different index page

Xlog assumes the main page is **index.md** if you're working in an existing github repository for example you may need to specify **README.md** as your index page as follows

```shell
xlog --index README
```

Notice that specifying the index page doesn't need the extension `.md`.

# Open your new site

Now you can navigate to [http://localhost:3000](http://localhost:3000) in your browser to start browsing the markdown files. if it's a new directory your editor will open to write your first page. 

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

If your markdown is hosted as a Github repository, you can use the [Xlog GitHub Action](Create your own digital garden on Github) to build the site and deploy it to GitHub Pages. The action checks out your repository, downloads the xlog binary, builds your site, and deploys it in a single workflow step.

Full tutorial: [Create your own digital garden on Github](Create your own digital garden on Github)

Examples of the workflow file:
- [Emad Elsaid Blog](https://github.com/emad-elsaid/emad-elsaid.github.io/blob/master/.github/workflows/xlog.yml)
- [Xlog documentation](https://github.com/emad-elsaid/xlog/blob/master/.github/workflows/xlog.yml)
