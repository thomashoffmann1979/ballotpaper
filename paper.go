package main

import (
	// "fmt"
	"log"
	"image/color"
	"gocv.io/x/gocv"
)

var pixelScale	float64 = 1
var pixelScaleY	float64 = 1


func processPaperChannelImage() {

	for range grabVideoCameraTicker.C {	
		img,ok := <-paperChannelImage
		// fmt.Println("got image",paper,ok,paper.Size())
		// log.Println("got image",ok,img.Size())
		if ok {
			if !img.Empty() {

				contour := findPaperContour(img)
				cornerPoints := getCornerPoints(contour)
				topLeftCorner := cornerPoints["topLeftCorner"]
				bottomRightCorner := cornerPoints["bottomRightCorner"]
				if false {
					log.Printf("template: %d %d",  bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y )
				}

				paper := extractPaper(img, contour, bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y, cornerPoints)
				mean := paper.Mean()
				if (mean.Val1+mean.Val2+mean.Val3)/3 > 150 {

					// Barcodescanner
					if len(scannerChannelImage)==cap(scannerChannelImage) {
						mat,_ := <-scannerChannelImage
						mat.Close()
					}
					scannerCloned := paper.Clone()
					scannerChannelImage <- scannerCloned
					
					// log.Println("mean",mean.Val1,mean.Val2,mean.Val3)

					if !paper.Empty() {
						if len(tesseractChannelImage)==cap(tesseractChannelImage) {
							mat,_:=<-tesseractChannelImage
							mat.Close()
						}
						imgGray := paper.Clone()

						tesseractChannelImage <- imgGray
					}
					paper.Close()
					



					
					drawContours := gocv.NewPointsVector()
					drawContours.Append(contour)
					if readyToSave {
						gocv.DrawContours(&img, drawContours, -1, color.RGBA{0, 255, 0, 120}, int(8.0*pixelScale))
					} else {
						gocv.DrawContours(&img, drawContours, -1, color.RGBA{255, 0, 0, 120}, int(8.0*pixelScale))
					}
					drawContours.Close()

				}
					
				if len(imageChannelPaper)==cap(imageChannelPaper) {
					mat,_:=<-imageChannelPaper
					mat.Close()
				}
				cloned := img.Clone()
				imageChannelPaper <- cloned

				img.Close()
			}
		}
		// paper.Close()
	}
}
