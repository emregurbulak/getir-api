package errorTypes

import (
	"fmt"
	"os"
)

func HandleCustomError(message string, err error) {
	if err != nil {
		fmt.Println(message + ", somthing went wrong. ErorDetails: " + err.Error())
	} else {
		fmt.Println(message + ", somthing went wrong.")
	}
	os.Exit(1)
}
