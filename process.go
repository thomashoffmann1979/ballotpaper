package main

import (
	"fmt"
	"image"
	// "sort"
	"image/color"
	"gocv.io/x/gocv"
)

func process(img gocv.Mat, template gocv.Mat, checkMarkList []CheckMarkList, lastTitle string) ([]CheckMarkList,string) {
//	maxCircles := 0

	rotated := gocv.NewMat()
	gocv.Rotate(img, &rotated, gocv.Rotate90CounterClockwise)
	contour := findPaperContour(rotated)
	cornerPoints := getCornerPoints(contour)
	paper := extractPaper(rotated, contour, template.Cols(), template.Rows(), cornerPoints)

	if paper.Empty() {
		return checkMarkList, lastTitle
	}

	imgGray := gocv.NewMat()
	gocv.CvtColor(paper, &imgGray, gocv.ColorBGRToGray)


	point,ballotPaperTitle,pageRois,pagesize,circleSize := tesseract(imgGray)
	fmt.Println("Ballot Paper Title: ", ballotPaperTitle,pageRois)
	if point.X!=0 && point.Y!=0{

		X := pageRois[0].X * paper.Cols() / pagesize.Width
		Y := pageRois[0].Y * paper.Rows() / pagesize.Height
		W := pageRois[0].Width * paper.Cols() / pagesize.Width
		H := pageRois[0].Height * paper.Rows() / pagesize.Height

		//rect:=image.Rect(point.X+55-30, point.Y+50 +40-10, point.X+55+50-30, point.Y+50+500-10)
		fmt.Println("paper",paper.Cols(),paper.Rows() )
		fmt.Println("SCALE",paper.Rows() / pagesize.Height)
		fmt.Println("pageRois[0]",pageRois[0])
		fmt.Println("point",point)
		rect:=image.Rect(point.X+X, point.Y+Y, point.X+X+W, point.Y+Y+H)
		fmt.Println("rect",rect)


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

	window := gocv.NewWindow("output")
	window.IMShow(rotated)
	window.WaitKey(0)


	if rotated.Empty() {
		return checkMarkList, lastTitle
	}
	return checkMarkList, lastTitle
	/*
	window.IMShow(rotated)
	window.WaitKey(1)
	*/
}