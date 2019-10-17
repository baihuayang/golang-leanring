package game

type (
	ErrorCode struct {
		id      int
		message string
	}
	newec map[string]int
)

var (
	errorCodes map[int]string
)

func add(errno int, msg string) int {
	if errorCodes == nil {
		errorCodes = make(map[int]string)
	}
	errorCodes[errno] = msg
	return errno
}

func init() {
	// errorCodes = make(map[int]string)
}

func PrintErrno(errno int) {
	println(errorCodes[errno])
}

var (
	SystemError = newec{
		"success": add(0x0000, "成功"),
		"forward": add(0x0001, "重定向"),
	}
)
