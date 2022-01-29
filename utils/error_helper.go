package utils

import (
	"fmt"
	"strings"
)

func HandleError(err error, field string) {
	if strings.Contains(err.Error(), "not ok with status") {
		fmt.Println("error from api while working on " + field)
	} else {
		fmt.Println(err.Error())
	}
}
