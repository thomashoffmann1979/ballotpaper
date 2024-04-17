package main

import (
	// "fmt"
	"log"
	//"image/color"
	"gocv.io/x/gocv"
)



func processTesseractChannelImage() {
	var tesseractReturn TesseractReturnType
	for range grabVideoCameraTicker.C {	
		img,ok := <-tesseractChannelImage
		if false {
			log.Println("got image *",ok,img.Size(),img.Empty(),img.Channels(),img.Rows(),img.Cols())
		}
		if ok {
			if !img.Empty() {
				imgGray := gocv.NewMat()
				gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)
				tesseractReturn = tesseract(imgGray)
				imgGray.Close()
				 
				channelData := RoisChannelStruct{
					tesseractReturn: tesseractReturn,
					mat: img.Clone(),
				}
				if false {
					log.Println(tesseractReturn);
				}
				if true {
					if len(tesseractReturnChannel)==cap(tesseractReturnChannel) {
						cD,_ := <-tesseractReturnChannel
						cD.mat.Close()
					}
					tesseractReturnChannel <- channelData
				}
			}
		}
		img.Close()
	}
}