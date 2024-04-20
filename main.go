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

	strSystemUrl string
	strSystemLogin string
	strSystemPassword string
	
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

	
	flag.StringVar(&tesseractPrefix, "tessdata", "" , "path to your tessdata directory")
	

	flag.StringVar(&strSystemUrl, "url", "http://localhost/wm/" , "system url")
	flag.StringVar(&strSystemLogin, "login", "max.muster@tualo.io" , "system login")
	flag.StringVar(&strSystemPassword, "password", "none" , "system password")

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
		case "app":
			fmt.Println("app")
			appwindow()
		case "camera":
			fmt.Println("camera",			dontScanBarcode		)
			fmt.Println("not implemented");
			// cameras( )
		case "detect":
			/*
			image := gocv.IMRead(strInputFile, gocv.IMReadColor)
			lastTesseract := TesseractReturnType{}
			process(image,lastTesseract);
			*/
			fmt.Println("not implemented");
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