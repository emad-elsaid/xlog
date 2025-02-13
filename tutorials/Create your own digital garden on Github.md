![](/docs/public/sprout.png)

#tutorial

# Create a Repository

Assuming your user name on github is `user-name`:

1. Create a repository in your account called **user-name.github.io** : https://github.com/new
2. Go to repository **Settings** (last tab on top of the repository)
3. Go to **Pages** from the sidebar under **Code and automation**
4. Under **Build and deployment** > **Source** choose **Github Actions**

# Create A Home page

1. Go to your repository **Code** page (first tab at the top of the repository)
2. Click **Add file** > **Create new file**
3. Name your file: **index.md**
4. Add any content to your file. for example: **Hello world!**
5. Click **Commit new file**

# Add Github workflow

1. Go to your repository **Code** page
2. Click **Add file** > **Create new file**
3. Name your file **.github/workflows/xlog.yml**
4. Add The following content to your file
```yaml
name: Xlog

on:
  push:
    branches: [ "master" ]

  # Allows you to run this workflow manually from the Actions tab
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Allow non-ASCII character
        run: git config core.quotepath false

      - name: restore timestamps
        uses: chetan/git-restore-mtime-action@v1

      - name: Install xlog
        env:
          XLOG_VERSION: v1.6.6
        run: curl --location https://github.com/emad-elsaid/xlog/releases/download/${XLOG_VERSION}/xlog-${XLOG_VERSION}-linux-amd64.tar.gz | tar -xz -C ..

      - name: Build
        run: |
          ../xlog \
          --build . \
          --sitename "user-name"
          rm *.md
          chmod -R 0777 .

      - name: Upload GitHub Pages artifact
        uses: actions/upload-pages-artifact@v3.0.1
        with:
          path: .

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```
5. if your main branch name is different than `master` please change it in the previous file.
6. Make sure you replace any occurrence of **user-name** with your user name.
7. Click **Commit new file**
8. Go to **Actions** tab. you should find a run in progress. After it completes move to the next step.

# Visiting your digital garden

* Now your digital garden should be available under **user-name.github.io**
* Visit the previous URL to make sure everything works
* If the page doesn't appear with **Hello world!** content (or the content you wrote) then something is wrong with the previous steps.

# Writing new pages

1. When you want to add a new page to your garden add the file in your repository and make sure it's in **.md** (Markdown) format.
2. After Github finishes building the page should be served as HTML with the same name without **.md** extension
3. For example **about.md** will be served as **/about**
4. keep your files names meaningful. as Xlog will autolink pages together by the file name.
5. so in any page when you have the word **about** in the text, Xlog will convert it to a link to the **about** page.
