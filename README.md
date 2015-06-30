# goexpose

For when you can't use full go-gettable names.

## Usage


### First, install it

```
go get github.com/MOZGIII/goexpose
```

### Then use it

```
$ cd /some/path/myproject
$ ls
package1  package2  package3  file1.go  file2.go
$ goexpose .
Exposing "/some/path/myproject" as "myproject" at "/home/user/gopath"
```

This creates a symlink at `/home/user/gopath/src/myproject` that points to `/some/path/myproject`.

And now you can use `myproject` in `import`s:

```go
import (
  "myproject"
  "myproject/package1"
  "myproject/package2"
)
```

### Custom project name

This example shows that you can pass your own name for the project.

```
$ cd /some/path/myotherproject
$ ls
package1  package2  package3  file1.go  file2.go
$ goexpose . awesomeproject
Exposing "/some/path/myotherproject" as "awesomeproject" at "/home/user/gopath"
```

This would create symlink at `/home/user/gopath/src/awesomeproject` that points to `/some/path/myotherproject`.

```go
import (
  "awesomeproject"
  "awesomeproject/package1"
  "awesomeproject/package2"
)
```

### Project name guessing

Noticed in the first example project name is omitted?
This is because this command can guess the project name from it's path.

It can even guess if you have a `/src` dir for your code.

For the project that has all code stored in `src` dir:

```
$ cd /some/path/to/project/root/mysuperproject
$ ls
examples  src  README.md
$ cd src
$ ls
package1  package2  package3  file1.go  file2.go
$ goexpose .
Exposing "/some/path/to/project/root/mysuperproject/src" as "mysuperproject" at "/home/user/gopath"
```

The symlink would be `/home/user/gopath/src/mysuperproject` that points to `/some/path/to/project/root/mysuperproject/src`.

As you can see, it takes not `src` but `mysuperproject` as a project name.
You may find it useful. Like I do. I made this thing that way, of course it's useful for me.

Guessing uses simple hardcoded heuristics.
Don't forget to file a pull request if you implement one that fits your needs.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
