package utils

func ReturnError(err *error) {
	if r := recover(); r != nil {
		e := r.(error)
		*err = e
	}
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
