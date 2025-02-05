package assert

import "fmt"

func Nil(val any) {
	if val != nil {
		panic(val)
	}
}

func Eq[T comparable](a, b T) {
	if a != b {
		panic(fmt.Sprintf("not equal, %+v, %+v", a, b))
	}
}
