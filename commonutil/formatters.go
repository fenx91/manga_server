package commonutil

import "fmt"

func TwoDigitInt(number int) string {
	return fmt.Sprintf("%02d", number)
}
