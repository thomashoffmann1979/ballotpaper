package main

import (
	"log"
	"gocv.io/x/gocv"
	"github.com/bieber/barcode"
)



func findBarcodes(scanner *barcode.ImageScanner, img gocv.Mat)[]BarcodeSymbol{
	//smaller:=ResizeMat(img, img.Cols() / barcodeScale, img.Rows() /barcodeScale)
	smaller:=gocv.NewMat()
	gocv.CvtColor(img, &smaller, gocv.ColorBGRToGray)
	log.Println("findBarcodes",barcodeScale)
	symbols, err := scanner.ScanMat(&smaller)
	if err != nil {
		panic(err)
	}
	log.Println("findBarcodes",len(symbols))
	syms := []BarcodeSymbol{}
	for _, s := range symbols {
		syms = append(syms,BarcodeSymbol{Type:s.Type.Name(),Data:s.Data,Quality:s.Quality,Boundary:s.Boundary})
		log.Println("BarcodeSymbol",s.Type.Name(),s.Data,s.Quality,s.Boundary)
	}
	return syms
}

func processImage(){
	scanner := barcode.NewScanner()
	scanner.SetEnabledAll(false)
	scanner.SetEnabledSymbology(barcode.Code39,true)
	scanner.SetEnabledSymbology(barcode.Code128,true)
	log.Println("processImage starting ")
	for {
		if !runVideo {
			break
		}
		log.Println("processImage",runVideo)
		//for range grabVideoCameraTicker.C {	
		img,ok := <-paperChannelImage
		if ok {
			log.Println("got image",ok,img.Size())

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
					area := float64(paper.Size()[0]) * float64(paper.Size()[1]) / float64(img.Size()[0]) / float64(img.Size()[1])
					log.Println("mean",mean.Val1,mean.Val2,mean.Val3,area,paper.Size())
					if area > 0.1 {
						findBarcodes(scanner,paper)
					}
				}
			}
		}
		log.Println("processImage done",runVideo)
	}
	log.Println("processImage exit",runVideo)
}