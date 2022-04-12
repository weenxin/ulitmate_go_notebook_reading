## go test 命令

`go help test`

`go test` 会自动运行引入的包。并产生如下的summary信息：

```
 ok   archive/tar   0.011s
 FAIL archive/zip   0.022s
 ok   compress/gzip 0.033s
```

'Go test' 会重新编译所有有`*_test.go`文件的包。这些文件可以包含功能测试，benchmark测试，模糊测试。**所有以_和.开头的文件都被忽略，包括_test.go**。

在*_test*开头的包内声明的测试文件，会作为一个单独的包测试，并会以主测试二进制编译，链接与运行。可以在每个测试包中加入`testdata`文件夹保存用于测试的文件。

`go test`前会运行`go vet`测试文件的签名问题。如果在`go vert`过程中发现问题，将停止运行`go test`。

`go test`两种模式：
- `go test` 编译运行当前文件测试文件。侧重模式下不`caching`
- `go test .` , `go test math` , `go test ./...` 会编译指定包






