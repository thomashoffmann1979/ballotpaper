package main

import (
	"fmt"
	//"os"
    "log"
	"time"
	"image"
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


func ResizeMat(img gocv.Mat,width int, height int) gocv.Mat {
	resizeMat := gocv.NewMat()

	if !img.Empty() {
		if img.Cols() >= width && img.Rows() >= height {
			if height>0 && width>0 {
				fmt.Println("ResizeMat",img.Cols(),img.Rows(),width,height)
				gocv.Resize(img, &resizeMat, image.Point{width, height}, 0, 0, gocv.InterpolationArea)
				img.Close()
			}
		}
	}
	if resizeMat.Empty() {
		return img
	}
	return resizeMat
}

func grabcamera( ) {
	webcam, err := gocv.VideoCaptureDeviceWithAPI(intCamera,0)
	/*
	webcam.Set(gocv.VideoCaptureFrameWidth, 1280*3)
	webcam.Set(gocv.VideoCaptureFrameHeight, 720*3)
	*/
	
//	webcam.Set(gocv.VideoCaptureFrameWidth, 1280*2)
//	webcam.Set(gocv.VideoCaptureFrameHeight, 720*2)

//	webcam.Set(gocv.VideoCaptureFrameWidth, 1280*1)
//	webcam.Set(gocv.VideoCaptureFrameHeight, 720*1)
//	4608 3456
	debug(fmt.Sprintf("Start grab camera %d ",forcedCameraWidth))

	if forcedCameraWidth > 0 {
		webcam.Set(gocv.VideoCaptureFrameWidth, float64(forcedCameraWidth))
	}
	if forcedCameraHeight > 0 {
		webcam.Set(gocv.VideoCaptureFrameHeight, float64(forcedCameraHeight))
	}

	debug(fmt.Sprintf("Start grab camera: %v %v",runVideo,intCamera))

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
		start:=time.Now()
		//debug("grabcamera")
		webcam.Read(&img)
		rotated := gocv.NewMat()

		gocv.Rotate(img, &rotated, gocv.Rotate90Clockwise)

		log.Println("grabcamera >>>>>>>>>>>>>>>>>>>>",rotated.Cols(),rotated.Rows(),time.Since(start))
		
		//debug( fmt.Sprintf("grab %s %d %d %d",time.Since(start),rotated.Cols(),rotated.Rows() , os.Getpid() ) )

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