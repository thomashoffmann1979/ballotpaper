package main

import (
	"log"
	"gocv.io/x/gocv"
	"image"
	"time"
	"github.com/bieber/barcode"
)



func findBarcodes(scanner *barcode.ImageScanner, img gocv.Mat)[]BarcodeSymbol{
	syms := []BarcodeSymbol{}
	if img.Empty() {
		return syms
	}
	smaller:=gocv.NewMat()
	gocv.CvtColor(img, &smaller, gocv.ColorBGRToGray)
	if smaller.Cols() > 800 {
		gocv.GaussianBlur(smaller, &smaller, image.Point{5, 5}, 0, 0, gocv.BorderDefault)
		gocv.Resize(smaller, &smaller, image.Point{smaller.Cols() / barcodeScale, smaller.Rows() / barcodeScale}, 0, 0, gocv.InterpolationArea)
	}
	log.Println("barcodeScale",barcodeScale)
	symbols, err := scanner.ScanMat(&smaller)
	if err != nil {
		panic(err)
	}
	
	/*
	log.Println("findBarcodes",len(symbols))
	if len(symbols) == 0 {
		gocv.IMWrite("noBarcode.png",img)
	}else{
		gocv.IMWrite("barcode.png",img)
	
	}
	*/
	
	for _, s := range symbols {
		syms = append(syms,BarcodeSymbol{Type:s.Type.Name(),Data:s.Data,Quality:s.Quality,Boundary:s.Boundary})
		log.Println("BarcodeSymbol",s.Type.Name(),s.Data,s.Quality,s.Boundary)
	}
	smaller.Close()
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
		start:=time.Now()
		log.Println("processImage ************")
		//for range grabVideoCameraTicker.C {	
		img,ok := <-paperChannelImage
		if ok {
			log.Println("got image",ok,img.Size(),len(paperChannelImage))

			if !img.Empty() {
				contour := findPaperContour(img)
				log.Println("findPaperContour done %s %v",time.Since(start),contour)

				cornerPoints := getCornerPoints(contour)
				topLeftCorner := cornerPoints["topLeftCorner"]
				bottomRightCorner := cornerPoints["bottomRightCorner"]
				if false {
					log.Printf("template: %d %d",  bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y )
				}

				paper := extractPaper(img, contour, bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y, cornerPoints)
				
				if paper.Empty() {
					contour.Close()
					img.Close()
					continue
				}
				// mean := paper.Mean()
				// if (mean.Val1+mean.Val2+mean.Val3)/3 > 150 {
					area := float64(paper.Size()[0]) * float64(paper.Size()[1]) / float64(img.Size()[0]) / float64(img.Size()[1])
					// log.Println("mean",mean.Val1,mean.Val2,mean.Val3,area,paper.Size(),time.Since(start))
					if area > 0.1 {
						findBarcodes(scanner,paper)
					}
				// }
				contour.Close()
				paper.Close()
			}
			img.Close()
		}
		log.Println("processImage done %s",time.Since(start))
	}
	log.Println("processImage exit",runVideo)
}