package main

import (
	"fmt"
	"flag"
	"encoding/json"
	"os"
	"gocv.io/x/gocv"
)

var (
	strInputFile string
	strCompareFile string
	strType string
	intCamera int
	boolVerbose bool
	documentConfigurations DocumentConfigurations
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
}

type DocumentConfigurations []struct {
	Titles       []string `json:"titles"`
	CircleSize   int `json:"circleSize"`
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
	flag.Parse()

	template := gocv.NewMat()

	dat, _ := os.ReadFile("config.json")
	json.Unmarshal([]byte(dat), &documentConfigurations)
	
	if strCompareFile != "" {
		template = gocv.IMRead(strCompareFile, gocv.IMReadColor)
	}

	switch strType {
		case "camera":
			cameras(template)
		case "detect":
			image := gocv.IMRead(strInputFile, gocv.IMReadColor)
			checkMarkList := []CheckMarkList{}
			lastTitle := ""
		
			process(image,template,checkMarkList,lastTitle);
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