Xlog is built to be small core that offers small set of features. And focus on offering a developer friendly public API to allow extending it with more features.

# Extension points

- Add any HTTP route with your handler function
- Add Goldmark extension for parsing or rendering
- Add a Preprocessor function to process page content before converting to HTML
- Listen to Page events such as Write or Delete.
- Define a helper function for templates
- Add a directory to be parsed as a template
- Add widgets in selected areas of the view page such as before or after rendered HTML
- Add a command to the list of commands triggered with `Ctrl+K` which can execute arbitrary Javascript.
- Add a route to be exported in case of building static site
- Add arbitrary link to pages or any URL
- Add quick command to appear on top of the view page

# Overview

An extension is a 

* Go module/package that imports xlog package
* Can be hosted anywhere
* Implements `xlog.Extension` interface
* Has an `Init` function to register all of its components using `Register*`
* Uses `RegisterExtension` functions in the `init` function of the package to register the extension
* Adds or improves a feature in xlog using one or more of the extension points.
* Imported by a the `main` package of your knowledgebase along with all other extensions and Xlog itself. an example can be found in Xlog CLI

# Creating extensions

* Hello world extension
