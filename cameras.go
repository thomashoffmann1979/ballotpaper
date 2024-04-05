package main

import (
	"fmt"
	"log"
	// "image"
	// "sort"
	// "image/color"
	"gocv.io/x/gocv"
)

type CheckMarks struct {
	Mean float64
    X       int 
	Y       int
	Radius	   int
}


type CheckMarkList struct {
	Count int
	Sum int
	AVG float64
	Checked bool
}

var runVideo bool = true

// ,checkMarkList []CheckMarkList, lastTitle string
func cameras( ) {
	
	webcam, err := gocv.VideoCaptureDeviceWithAPI(intCamera,0)
	if err != nil {
		fmt.Println("Error opening capture device: ", 0)
		return
	}
	defer webcam.Close()

	webcam.Set(gocv.VideoCaptureFrameWidth, 1280*3)
	webcam.Set(gocv.VideoCaptureFrameHeight, 720*3)
	
	
	img := gocv.NewMat()
	defer img.Close()
	checkMarkList := []CheckMarkList{}
	//lastTitle := ""
	lastTesseract := TesseractReturnType{}

	currentBox:=""
	currentStack:=""

	for runVideo {
		webcam.Read(&img)

		if runProcess {
			processResult := process(img , lastTesseract)

			if processResult.FCBarcode!="" {
				log.Printf("using BOX: %s ", processResult.FCBarcode)
				if processResult.FCBarcode[0:3]=="FC4" {
					currentBox=processResult.FCBarcode
					currentStack = processResult.FCBarcode[0:2]+"3"+processResult.FCBarcode[3:]
				}
				if processResult.FCBarcode[0:3]=="FC3" {
					currentStack=processResult.FCBarcode
				}
				//log.Printf("Saving BOX: %s ", processResult.FCBarcode)
			}
			if processResult.IsCorrect {
				if processResult.Title != lastTesseract.Title || processResult.Barcode != lastTesseract.Barcode {
					checkMarkList = []CheckMarkList{}
				}

				


				for i := 0; i < len(lastTesseract.Marks); i++ {
					if i >= len(checkMarkList) {
						checkMarkList = append(checkMarkList, CheckMarkList{})
					}
					if lastTesseract.Marks[i] {
						checkMarkList[i].Sum += 1
					}
					checkMarkList[i].Count++
					checkMarkList[i].AVG = float64(checkMarkList[i].Sum) / float64(checkMarkList[i].Count)
					checkMarkList[i].Checked = checkMarkList[i].AVG > 0.6
				}


				outList:=[]string{}
				for i := 0; i < len(checkMarkList); i++ {
					
					if checkMarkList[i].Checked {
						outList = append(outList, "😎")
					} else {
						outList = append(outList, "🥶")
					}
				}	

				if len(checkMarkList)>0 && checkMarkList[0].Count>5 {
					if currentBox!="" {
						log.Printf("Box: %s, Stack: %s, Barcode: %s, Title: %s, Marks: %s",currentBox,currentStack, processResult.Barcode , processResult.Title, outList)
					//	processResult=TesseractReturnType{}
					}
					//checkMarkList = []CheckMarkList{}
				}
			}
			lastTesseract = processResult
		}
		showImage("camera", img, 1)
	}
}
