package main

import "io"

type myWriter struct {
}

/*
	func (w myWriter) Write(p []byte) (n int, err error) {
		return
	}
*/
// func (w myWriter) Write(p []byte) (n int, err error) {
// 	return
// }
func main() {
	// 检查 *myWriter 类型是否实现了 io.Writer 接口
	var _ io.Writer = (*myWriter)(nil)

	// 检查 myWriter 类型是否实现了 io.Writer 接口
	var _ io.Writer = myWriter{}
}

/*
# command-line-arguments
interface/check.go:18:20: cannot use (*myWriter)(nil) (value of type *myWriter) as type io.Writer in variable declaration:
	*myWriter does not implement io.Writer (missing Write method)
interface/check.go:21:20: cannot use myWriter{} (value of type myWriter) as type io.Writer in variable declaration:
	myWriter does not implement io.Writer (missing Write method)
*/
