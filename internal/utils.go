package utils

import (
	"fmt"
	"os"
)

// CheckError is a utility function to handle errors gracefully.
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
