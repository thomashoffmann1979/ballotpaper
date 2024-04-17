package main

import (
	"fmt"
	// "log"
	// "image"
	// "sort"
	// "image/color"
	"gocv.io/x/gocv"
)



func getCameraList() []CameraList {
	cameraList := []CameraList{}
	for i := 0; i < 5; i++ {
		webcam, err := gocv.VideoCaptureDeviceWithAPI(i,0)
		if err != nil {
			return cameraList
		}
		fmt.Println("Cam: ",i, webcam.Get(gocv.VideoCaptureFrameWidth), webcam.Get(gocv.VideoCaptureFrameHeight))
		cameraList = append(cameraList, CameraList{int(webcam.Get(gocv.VideoCaptureFrameWidth)), int(webcam.Get(gocv.VideoCaptureFrameHeight)), i, fmt.Sprintf("Camera %d",i)})
		webcam.Close()
	}
	return cameraList
}

func grabcamera( ) {
	webcam, err := gocv.VideoCaptureDeviceWithAPI(intCamera,0)
	/*
	webcam.Set(gocv.VideoCaptureFrameWidth, 1280*3)
	webcam.Set(gocv.VideoCaptureFrameHeight, 720*3)
	
	webcam.Set(gocv.VideoCaptureFrameWidth, 1280*2)
	webcam.Set(gocv.VideoCaptureFrameHeight, 720*2)
	*/
	
	if err != nil {
		fmt.Println("Error opening capture device: ", 0)
		return
	}
	defer webcam.Close()

	img := gocv.NewMat()
	defer img.Close()
	/*
	checkMarkList := []CheckMarkList{}
	lastReturnType := ReturnType{}
	*/
	for runVideo {
		webcam.Read(&img)
		rotated := gocv.NewMat()

		gocv.Rotate(img, &rotated, gocv.Rotate90Clockwise)


		// Videooutput
		if len(cameraChannelImage)==cap(cameraChannelImage) {
			mat,_ := <-cameraChannelImage
			mat.Close()
		}
		cameraCloned := rotated.Clone()
		cameraChannelImage <- cameraCloned

		// Paper
		if len(paperChannelImage)==cap(paperChannelImage) {
			mat,_ := <-paperChannelImage
			mat.Close()
		}
		paperCloned := rotated.Clone()
		paperChannelImage <- paperCloned

		rotated.Close()
	}
}