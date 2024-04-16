package errors

import (
	"fmt"
)

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Print(LoxErrorFmt(line, where, message))
}

func LoxErrorFmt(line int, where string, message string) string {
	return fmt.Sprintf("[line %d] Error%s: %s\n", line, where, message)
}
