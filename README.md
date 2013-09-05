fileutil
========

A collection of simple file related utilities which are not provided by the standard library.

Provides:

```go

   // Will fully copy all the files underneath the src/ directory into
   // the dst/ directory.
   fileutil.CopyDirectory("/path/to/destination/", "/path/to/directory/")
   ok, err := fileutil.IsSymLink("/path/to/possible/symlink/")
   if err != nil {
       log.Fatal(err)
   }
   if ok {
       os.Remove("/path/to/possible/symlink/")
   }
   
   plus, minus, err := fileutil.DiffDirectories("/path/to/dir/", "/path/to/dir2/")
   if err != nil {
       log.Fatal(err)
   }
   if plus != nil || minus != nil {
       log.Println("The directories do not match.")
   }
```
