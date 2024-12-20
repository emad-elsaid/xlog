Xlog serves any files under current directory with exception of markdown files being accessed without `.md` extension and converted to HTML.

Besides that it serves files from embded files in the program from the core package or extensions.

# Overriding asset files 

Embded files are the last resort when looking up a file so to override an asset file you just need to put it in the same path in the current directory. that's all. that simple.

## CSS

- Xlog used to have a Go script to compile CSS/SASS to `public/style.css`.
- That changed to depend on Webpack in this commit 38c8171
- So chdir to `cmd/assets` and either build with `yarn build` or watch changes `yarn watch`
