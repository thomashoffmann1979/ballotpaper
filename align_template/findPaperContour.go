package main

import (
	// "fmt"
	"image"
	"gocv.io/x/gocv"
)

func findPaperContour(img gocv.Mat) gocv.PointVector {
	imgGray := gocv.NewMat()
	gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)

	imgBlur := gocv.NewMat()
	gocv.GaussianBlur(imgGray, &imgBlur, image.Point{5, 5}, 0, 0, gocv.BorderDefault)

	imgThresh := gocv.NewMat()
	gocv.Threshold(imgBlur, &imgThresh, 90, 120, gocv.ThresholdBinary+gocv.ThresholdOtsu)
	
	showImage("imgThresh", imgThresh, 0)

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
	maxContour := contours.At(maxContourIndex)

	imgGray.Close()
	imgBlur.Close()
	imgThresh.Close()
	return maxContour
}
/*
findPaperContour: function(img) {
	const imgGray = new cv.Mat();
	cv.cvtColor(img, imgGray, cv.COLOR_RGBA2GRAY);

	const imgBlur = new cv.Mat();
	cv.GaussianBlur(
	  imgGray,
	  imgBlur,
	  new cv.Size(5, 5),
	  0,
	  0,
	  cv.BORDER_DEFAULT
	);

	

	const imgThresh = new cv.Mat();
	cv.threshold(
	  imgBlur,
	  imgThresh,
	  0,
	  255,
	  cv.THRESH_BINARY + cv.THRESH_OTSU
	);

	let contours = new cv.MatVector();
	let hierarchy = new cv.Mat();

	
	cv.findContours(
	  imgThresh,
	  contours,
	  hierarchy,
	  cv.RETR_CCOMP,
	  cv.CHAIN_APPROX_SIMPLE
	);
	let maxArea = 0;
	let maxContourIndex = -1;
	for (let i = 0; i < contours.size(); ++i) {
	  let contourArea = cv.contourArea(contours.get(i));
	  if (contourArea > maxArea) {
		maxArea = contourArea;
		maxContourIndex = i;
	  }
	}
	if (maxContourIndex === -1) {
		imgGray.delete();
		imgBlur.delete();
		imgThresh.delete();
		contours.delete();
		hierarchy.delete();
		return null;
	}
	const maxContour = contours.get(maxContourIndex);

	imgGray.delete();
	imgBlur.delete();
	imgThresh.delete();
	contours.delete();
	hierarchy.delete();
	return maxContour;
},
*/