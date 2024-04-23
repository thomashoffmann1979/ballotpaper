package main

import (
	"flag"
	"time"
	"log"
	"gocv.io/x/gocv"
	"os"
)
var intTestCamera int
var forcedTestCameraWidth int
var forcedTestCameraHeight int


//var wnd gocv.Window
func main(){

	flag.IntVar(&intTestCamera, "camera", 0 , "camera")
	flag.IntVar(&forcedTestCameraWidth, "width", 0 , "width")
	flag.IntVar(&forcedTestCameraHeight, "height", 0 , "height")
	flag.Parse()

	

	webcam, _ := gocv.VideoCaptureDeviceWithAPI(intTestCamera,0)
	if forcedTestCameraWidth > 0 {
		webcam.Set(gocv.VideoCaptureFrameWidth, float64(forcedTestCameraWidth))
	}
	if forcedTestCameraHeight > 0 {
		webcam.Set(gocv.VideoCaptureFrameHeight, float64(forcedTestCameraHeight))
	}
	run := true
	img := gocv.NewMat()
	defer img.Close()

	for	run {
		start := time.Now()
		webcam.Read(&img)
		rotated := gocv.NewMat()

		gocv.Rotate(img, &rotated, gocv.Rotate90Clockwise)
		log.Println("grabvideo_test",time.Since(start),rotated.Cols(),rotated.Rows(),os.Getpid())

	}
}