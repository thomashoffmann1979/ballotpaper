package main

import (
	"fmt"
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

// ,checkMarkList []CheckMarkList, lastTitle string
func cameras(template gocv.Mat) {
	
	webcam, err := gocv.VideoCaptureDeviceWithAPI(intCamera,0)
	if err != nil {
		fmt.Println("Error opening capture device: ", 0)
		return
	}
	defer webcam.Close()
	webcam.Set(gocv.VideoCaptureFrameWidth, 1280*3)
	webcam.Set(gocv.VideoCaptureFrameHeight, 720*3)
	window := gocv.NewWindow("cameras")
	img := gocv.NewMat()
	defer img.Close()
	checkMarkList := []CheckMarkList{}
	lastTitle := ""
	lastTesseract := TesseractReturnType{}

	for {
		webcam.Read(&img)
		checkMarkList,lastTitle,lastTesseract = process(img, template, checkMarkList, lastTitle, lastTesseract)
		window.IMShow(img)
		window.WaitKey(1)
	
	}
}
