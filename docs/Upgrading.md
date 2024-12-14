# Upgrading to V2

If you were a user for v1 and would like to upgrade to v2 please take the following steps:

* `--sidebar` command-line argument is removed. it had no effect for a long time already and was kept for backward compatibility
* The `book` extension is renamed to `blocks` and has the same `book` shortcode. if you're importing it manually you need to change the import path to import `blocks` instead
* For extensions development 
  * You're now required to `xlog.RegisterExtension` your extension in your `init()` function then `Register*` the rest of your components in the extension `.Init()` function instead of the global one. this allow for future development to enable/disable extensions by the user.
  * Your extension can now check for `xlog.Config.Readonly` instead of `xlog.READONLY` during initialization (extension `Init()`)
