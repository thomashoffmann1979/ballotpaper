package main

import (
	"fmt"
	"flag"
	"time"
	"encoding/json"
	"os"
	"gocv.io/x/gocv"
)

var (
	strInputFile string
	strCompareFile string
	strType string

	showTesseractCropped bool
	showScannerImage bool

	dontScanBarcode bool
	intCamera int
	boolVerbose bool
	documentConfigurations DocumentConfigurations

	runProcess bool
	runContour bool
	runTesseract bool
	runScanner bool
	runMarkdetection bool

	start time.Time
	scannerImageSmallShrink int = 3
)

type DocumentConfigurationPageSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DocumentConfigurationPageRoi struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
	ExcpectedMarks int `json:"excpectedMarks"`
	Titles       []string `json:"titles"`
}

type DocumentConfigurations []struct {
	Titles       []string `json:"titles"`
	CircleSize   int `json:"circleSize"`
	CircleMinDistance int `json:"circleMinDistance"`
	TitleRegion struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"titleRegion"`
	Pagesize DocumentConfigurationPageSize `json:"pagesize"`
	Rois []DocumentConfigurationPageRoi `json:"rois"`
}

func main() {

	flag.StringVar(&strType, "type", "" , "detect type of image")
	flag.StringVar(&strInputFile, "input", "" , "input file")
	flag.IntVar(&intCamera, "camera", 0 , "camera")
	flag.StringVar(&strCompareFile, "template", "" , "template file")
	flag.BoolVar(&boolVerbose, "verbose", false, "verbose output")
	flag.BoolVar(&dontScanBarcode, "nobarcode", true, "do not scan barcode")
	flag.BoolVar(&showTesseractCropped, "showtesseractcropped", false, "show tesseract cropped")
	flag.BoolVar(&showScannerImage, "showscannerimage", false, "show scanner image")
	flag.Parse()

	start = time.Now()

	runContour=true
	runProcess=true
	runTesseract=true
	runMarkdetection=true
	runScanner=true

	template := gocv.NewMat()
	defer template.Close()

	dat, _ := os.ReadFile("config.json")
	json.Unmarshal([]byte(dat), &documentConfigurations)
	
	if strCompareFile != "" {
		template = gocv.IMRead(strCompareFile, gocv.IMReadColor)
	}

	switch strType {
		case "camera":
			fmt.Println("camera",			dontScanBarcode		)
			cameras(template)
		case "detect":
			image := gocv.IMRead(strInputFile, gocv.IMReadColor)

			lastTesseract := TesseractReturnType{}
		
			process(image,template,lastTesseract);
		case "compare":
			fmt.Println("comparing image")
		case "help":
			help()
		//case "tesseract":
		//	tesseract()
		default:

			help()
	}
}