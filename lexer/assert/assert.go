package assert

func Nil(val any) {
	if val != nil {
		panic(val)
	}
}
