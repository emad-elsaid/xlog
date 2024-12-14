Xlog works in 2 modes:

* Read/Write
* ReadOnly

By default xlog server works in Read/Write mode where you can edit and delete files. this mode is not for production use. it's meant for local personal use. and this is meant for the first usecase: taking personal notes, local digital gardening.

ReadOnly mode which can be specified using `--readonly=true` flag. This flag is checked by xlog and extensions to turn of any code that writes to the filesystem. 

/alert don't run xlog server on production server neither in read/write nor in readonly. as it's meant for personal local use.

Generate static website process will turn on readonly mode automatically.

Any extensions that writes or modify the filesystem is responsible for checking if `Config.Readonly` global variable is true and make sure that part is not executed.
