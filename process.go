package main

import (
	"fmt"
	"image"
	// "sort"
	"image/color"
	"gocv.io/x/gocv"
	"github.com/bieber/barcode"
)


func process(img gocv.Mat, template gocv.Mat, checkMarkList []CheckMarkList, lastTitle string, lastTesseract TesseractReturnType) ([]CheckMarkList,string,TesseractReturnType) {
//	maxCircles := 0

	rotated := gocv.NewMat()
	gocv.Rotate(img, &rotated, gocv.Rotate90CounterClockwise)
	var tesseractReturn TesseractReturnType

	if rotated.Empty() {
		return checkMarkList, lastTitle,lastTesseract
	}

	contour := findPaperContour(rotated)
	cornerPoints := getCornerPoints(contour)
	paper := extractPaper(rotated, contour, template.Cols()*2, template.Rows()*2, cornerPoints)

	if paper.Empty() {
		return checkMarkList, lastTitle,lastTesseract
	}

	imgGray := gocv.NewMat()
	gocv.CvtColor(paper, &imgGray, gocv.ColorBGRToGray)


	// read barcode with zbar from the frame
	scanner := barcode.NewScanner().SetEnabledAll(false)
	//scanner.SetEnabledSymbology(barcode.Code39,true)
	scanner.SetEnabledSymbology(barcode.Code128,true)

	rect := image.Rect(paper.Cols()/2, 0, paper.Cols(), paper.Rows()/10)
	scannerImage := paper.Region(rect)

	scannerImageWindow := gocv.NewWindow("scannerImageWindow")
	scannerImageWindow.IMShow(scannerImage)


	/*
	imgObj, _ := scannerImage.ToImage()
	src := barcode.NewImage(imgObj)
	*/
	symbols, _ := scanner.ScanMat(&scannerImage)
	for _, s := range symbols {
		data := s.Data
		fmt.Println("data",data)
	}
	defer scannerImage.Close()

	if true {
		if lastTesseract.Title!="" {
			tesseractReturn = lastTesseract
		}else{
			tesseractReturn = tesseract(imgGray)
		}
		point := tesseractReturn.Point
		circleSize := tesseractReturn.CircleSize
		if boolVerbose {
			fmt.Println("Ballot Paper Title: ", tesseractReturn.Title,tesseractReturn.PageRois)
		}
		if point.X!=0 && point.Y!=0{

			X := tesseractReturn.PageRois[0].X * paper.Cols() / tesseractReturn.Pagesize.Width
			Y := tesseractReturn.PageRois[0].Y * paper.Rows() / tesseractReturn.Pagesize.Height
			W := tesseractReturn.PageRois[0].Width * paper.Cols() / tesseractReturn.Pagesize.Width
			H := tesseractReturn.PageRois[0].Height * paper.Rows() / tesseractReturn.Pagesize.Height

			//rect:=image.Rect(point.X+55-30, point.Y+50 +40-10, point.X+55+50-30, point.Y+50+500-10)
			//fmt.Println("paper",paper.Cols(),paper.Rows() )
			//fmt.Println("SCALE",paper.Rows() / pagesize.Height)
			//fmt.Println("pageRois[0]",pageRois[0])
			//fmt.Println("point",point)
			rect:=image.Rect(point.X+X, point.Y+Y, point.X+X+W, point.Y+Y+H)
			//fmt.Println("rect",rect)


			x := paper.Clone() 
			croppedMat := x.Region(rect)
			gocv.Rectangle(&paper, rect, color.RGBA{255, 255, 0, 0}, 4)
			if !croppedMat.Empty() {
				
				findCircles(croppedMat, circleSize)

				/*
				croppedMatGray := gocv.NewMat()
				// gocv.IMWrite("croppedMat.jpg", croppedMat)
				gocv.CvtColor(croppedMat, &croppedMatGray, gocv.ColorBGRToGray)
				//gocv.MedianBlur(croppedMatGray, &croppedMatGray, 5)
				circles := gocv.NewMat()
				defer circles.Close()

				gocv.HoughCirclesWithParams(
					croppedMatGray,
					&circles,
					gocv.HoughGradient,
					1,                     // dp
					20, //float64(croppedMatGray.Rows()/50), // minDist
					90,                    // param1
					10,                    // param2
					4,                    // minRadius
					10,                     // maxRadius
				)

				// blue := color.RGBA{0, 0, 255, 0}
				// red := color.RGBA{255, 0, 0, 0}

				if false {

				}
				*/

			}
		}

		paperwindow := gocv.NewWindow("paper")
		paperwindow.IMShow(paper)

		drawContours := gocv.NewPointsVector()
		defer drawContours.Close()
		drawContours.Append(contour)
		gocv.DrawContours(&rotated, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)

	}
	window := gocv.NewWindow("output")
	window.IMShow(rotated)
	window.WaitKey(1)


	if rotated.Empty() {
		return checkMarkList, lastTitle,tesseractReturn
	}

	window.IMShow(rotated)
	window.WaitKey(1)

	return checkMarkList, lastTitle,tesseractReturn
	
	
}