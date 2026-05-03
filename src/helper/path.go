package helper

import (
	"fmt"
	"log"
	"os"
)

func GetCurrentDirectory() string {
	cd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return cd
}

func PrintCurrentDirectoryOrDie() {
	fmt.Println(fmt.Sprintf("Current directory: %s", GetCurrentDirectory()))
}
