# Upgrading to V2

If you were a user for v1 and would like to upgrade to v2 please take the following steps:

* `--sidebar` command-line argument is removed. it had no effect for a long time already and was kept for backward compatibility
* The `book` extension is renamed to `blocks` and has the same `book` shortcode. if you're importing it manually you need to change the import path to import `blocks` instead
* For extensions development 
  * You're now required to `xlog.RegisterExtension` your extension in your `init()` function then `Register*` the rest of your components in the extension `.Init()` function instead of the global one. this allow for future development to enable/disable extensions by the user.
  * Your extension can now check for `xlog.Config.Readonly` instead of `xlog.READONLY` during initialization (extension `Init()`)
  * `Get/Post/Delete/..etc` doesn't accept middlewares parameters anymore
* The version extension has been removed as it was incomplete. if you are running xlog on the same directory that has `.versions` subdirectories you can remove them by running `rm -rf *.versions`
* `github.repo` and `github.branch` are removed in favor of `github.url` which is the full URL of the editing. so it should work with other git online editors
* `custom_css` was removed as its functionality can be achieved by using `custom.head` 
* `custom_head/before_view/after_view` name changed to `custom.` replacing the `_` with `.` for consistency with other flags
