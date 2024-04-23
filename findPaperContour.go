package main

import (
	// "fmt"
	"image"
	"gocv.io/x/gocv"
)

func findPaperContour(img gocv.Mat) gocv.PointVector {
	imgGray := gocv.NewMat()
	gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)


	// gocv.Resize(imgGray, &imgGray, image.Point{imgGray.Cols() / 8, imgGray.Rows() / 8}, 0, 0, gocv.InterpolationArea)


	imgBlur := gocv.NewMat()
	gocv.GaussianBlur(imgGray, &imgBlur, image.Point{5, 5}, 0, 0, gocv.BorderDefault)
	imgThresh := gocv.NewMat()
	gocv.Threshold(imgBlur, &imgThresh, 90, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)
	contours := gocv.FindContours(imgThresh, gocv.RetrievalCComp, gocv.ChainApproxSimple)
	  
	maxArea := 0.0
	maxContourIndex := -1
	for i := 0; i < contours.Size(); i++ {
		contourArea := gocv.ContourArea(contours.At(i))
		if contourArea > maxArea {
			maxArea = contourArea
			maxContourIndex = i
		}
	}
	if maxContourIndex == -1 {
		imgGray.Close()
		imgBlur.Close()
		imgThresh.Close()
		return gocv.NewPointVector()
	}
	maxContour := gocv.NewPointVectorFromPoints( contours.At(maxContourIndex).ToPoints() )
	imgGray.Close()
	imgBlur.Close()
	imgThresh.Close()
	contours.Close()
	return maxContour
}
