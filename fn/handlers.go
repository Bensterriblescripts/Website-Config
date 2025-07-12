package fn

import "fmt"

func CheckPanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
func CheckPrintErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func PanicErr[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
func PrintErr[T any](value T, err error) T {
	if err != nil {
		fmt.Println(err)
		return value
	}
	return value
}
