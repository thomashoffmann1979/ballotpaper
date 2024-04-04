package main

import (
	"fmt"
	"image"
	"gocv.io/x/gocv"
	"github.com/bieber/barcode"
)

func scanBarcode(paper gocv.Mat) string {

	result:=""
	imgGray := gocv.NewMat()
	defer imgGray.Close()

	gocv.CvtColor(paper, &imgGray, gocv.ColorBGRToGray)
	scanner := barcode.NewScanner().SetEnabledAll(false)
	//scanner.SetEnabledSymbology(barcode.Code39,true)
	scanner.SetEnabledSymbology(barcode.Code128,true)

	rect := image.Rect(paper.Cols()/2, 0, paper.Cols(), paper.Rows()/10)
	if boolVerbose {
		fmt.Println("scanBarcode",paper.Cols(),paper.Rows())
	}
	if (paper.Rows()<400){
		scanner.SetEnabledSymbology(barcode.Code39,true)
		rect = image.Rect(0 , 0, paper.Cols(), paper.Rows() )
	}

	//rect := image.Rect(0, 0, paper.Cols(), paper.Rows()/10)
	if rect.Empty() {
		return result
	}
	scannerImage := paper.Region(rect)
	if scannerImage.Empty() {
		return result
	}
	if scannerImage.Rows()/ scannerImageSmallShrink < 36 || scannerImage.Cols() / scannerImageSmallShrink < 36 {
		return result
	}

	scannerImageSmall := gocv.NewMat()
	defer scannerImageSmall.Close()
	gocv.Resize(scannerImage, &scannerImageSmall, image.Point{scannerImage.Cols()/ scannerImageSmallShrink, scannerImage.Rows()/ scannerImageSmallShrink}, 0, 0, gocv.InterpolationLinear)

	if (paper.Rows()<400){
		scannerImageSmall = scannerImage.Clone()
	}
	
	if showScannerImage {
		scannerImageWindow := gocv.NewWindow("scannerImageWindow")
		scannerImageWindow.IMShow(scannerImageSmall)
	}

	symbols, _ := scanner.ScanMat(&scannerImageSmall)
	for _, s := range symbols {
		result = s.Data
		//fmt.Println(result)
	}
	defer scannerImage.Close()
	return result
}