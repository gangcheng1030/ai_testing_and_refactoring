package globalvariable

var Global GlobalVariable

type GlobalVariable struct {
	a int
}

func Add(a int) int {
	Global.a += a
	return Global.a
}

func Minus(a int) int {
	Global.a -= a
	return Global.a
}

func GetA() int {
	return Global.a
}

func SetA(a int) {
	Global.a = a
}
