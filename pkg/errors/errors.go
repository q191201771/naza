// package errors 错误处理相关
package errors

func PanicIfErrorOccur(err error) {
	if err != nil {
		panic(err)
	}
}
