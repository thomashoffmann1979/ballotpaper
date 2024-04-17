package main

import (
	"gocv.io/x/gocv"
)

// var window *gocv.Window

var processing bool = false
var pcount int = 0


func showImage( name string,  img gocv.Mat ) {
	if strType != "app" {
		window := gocv.NewWindow(name)
		window.IMShow(img)
		window.WaitKey(1)
	}else{
		if name == "output" {
			if len(cameraChannelImage)==cap(cameraChannelImage) {
				<-cameraChannelImage
			}
			//fmt.Println("showImage",name)
			cloned := img.Clone()
			cameraChannelImage <- cloned
		}
		if name == "paper" {
			if len(imageChannelPaper)==cap(imageChannelPaper) {
				<-imageChannelPaper
			}
			//fmt.Println("showImage",name,len(imageChannelPaper),cap(imageChannelPaper))
			cloned := img.Clone()

			imageChannelPaper <- cloned
		}
		if name == "scannerImageWindow" {
			if len(imageChannelCircle)==cap(imageChannelCircle) {
				<-imageChannelCircle
			}
			//fmt.Println("showImage",name,len(imageChannelPaper),cap(imageChannelPaper))
			cloned := img.Clone()

			imageChannelCircle <- cloned
		}
	}
}

func showImageBlocked( name string,  img gocv.Mat, waitKey int ) {
	if strType != "app" {
		window := gocv.NewWindow(name)
		window.IMShow(img)
		window.WaitKey(1)
	}else{
		if name == "camera" && processing == false{
			processing = true
			if pcount==3 {
				image := matToImage(img)
				outputImage.Image = image
				outputImage.Refresh()
			}
			pcount++
			if pcount > 3 {
				pcount = 0
				// fmt.Println("showImage",name)
			}
			processing = false
		}
		// fmt.Println("showImage",name)
	}
}
