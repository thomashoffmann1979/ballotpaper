package main

import (
	// "fmt"
	// "log"
	"image"
	"gocv.io/x/gocv"
	"github.com/bieber/barcode"
)

type BarcodeSymbol struct {
	Type string
	Data string
	Quality int
	Boundary []image.Point
}

var scannerChannelImage = make(chan gocv.Mat, 1)
var scannerChannelBarcodes = make(chan []BarcodeSymbol, 10)

// used only in grabcamera.go
var currentPageBarcode = make(chan string, 1)
var currentStackBarcode = make(chan string, 1)
var currentBoxBarcode = make(chan string, 1)


func processBarcodeSymbol(barcodeSymbol *barcode.Symbol ) {
	//blog.Println("got barcode",barcodeSymbol.Type,barcodeSymbol.Data)
	if barcode.Code128 == barcodeSymbol.Type{
		if len(currentPageBarcode) == cap(currentPageBarcode) {
			<-currentPageBarcode
		}
		currentPageBarcode <- barcodeSymbol.Data
	}
	if barcode.Code39 == barcodeSymbol.Type{
		if len(barcodeSymbol.Data)>3 {
			if barcodeSymbol.Data[0:3]=="FC4" {
				if len(currentBoxBarcode) == cap(currentBoxBarcode) {
					<-currentBoxBarcode
				}
				currentBoxBarcode <- barcodeSymbol.Data
			}
			if barcodeSymbol.Data[0:3]=="FC3" {
				if len(currentStackBarcode) == cap(currentStackBarcode) {
					<-currentStackBarcode
				}
				currentStackBarcode <- barcodeSymbol.Data
			}
		}
	}
}

func scanBarcodeChannel() {
	scanner := barcode.NewScanner()
	scanner.SetEnabledAll(false)
	scanner.SetEnabledSymbology(barcode.Code39,true)
	scanner.SetEnabledSymbology(barcode.Code128,true)
	for range grabVideoCameraTicker.C {	
		img,ok := <-scannerChannelImage
		// fmt.Println("got image",paper,ok,paper.Size())
		// log.Println("got image",ok,img.Size())
		if ok {
			if !img.Empty() {
				gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
				// mean := img.Mean()
				//if (mean.Val1+mean.Val2+mean.Val3)/3 > 100 {

					symbols, err := scanner.ScanMat(&img)
					if err != nil {
						panic(err)
					}
					syms := []BarcodeSymbol{}
					for _, s := range symbols {
						syms = append(syms,BarcodeSymbol{Type:s.Type.Name(),Data:s.Data,Quality:s.Quality,Boundary:s.Boundary})
						processBarcodeSymbol(s)
					}
					if len(scannerChannelBarcodes) == cap(scannerChannelBarcodes) {
						<-scannerChannelBarcodes
					}
					scannerChannelBarcodes <- syms
				//}
				img.Close()
			}
		}
		// paper.Close()
	}
}
