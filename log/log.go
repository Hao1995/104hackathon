package log

import "fmt"

//Errorf ...
func Errorf(str string) {
	fmt.Println("[ERROR] ", str)
}
