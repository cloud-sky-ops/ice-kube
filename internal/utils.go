package utils

import (
	"log"
	"os"

	"github.com/TwiN/go-color"
)

// PrintError is a utility function to handle errors gracefully.
func PrintError(errorMessage string, err error) {
	if err != nil {
		log.Println(errorMessage, color.Ize(color.Red, err)) // using log package to add timestamp with the error
		os.Exit(1)
	}
}
