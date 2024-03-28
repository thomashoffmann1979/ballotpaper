package main

import (
	"fmt"
	"flag"
)

var (
	strInputFile string
	strCompareFile string
	strType string
	intCamera int
	boolVerbose bool
)

func main() {
    fmt.Println("Hello, world.")

	flag.StringVar(&strType, "type", "" , "detect type of image")
	flag.StringVar(&strInputFile, "input", "" , "input file")
	flag.IntVar(&intCamera, "camera", 0 , "camera")
	flag.StringVar(&strCompareFile, "template", "" , "template file")
	flag.BoolVar(&boolVerbose, "verbose", false, "verbose output")
	flag.Parse()

	switch strType {
		case "camera":
			cameras()
		case "detect":
			fmt.Println("detecting image")
		case "compare":
			fmt.Println("comparing image")
		case "help":
			help()
		case "tesseract":
			tesseract()
		default:

			help()
	}
}