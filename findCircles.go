package main

import (
	// "fmt"
	"image"
 	"sort"
	"math"
	"image/color"
	"gocv.io/x/gocv"
)

func DrawCircles(img *gocv.Mat, circles *gocv.Mat,  innerOverdraw int, outerOverdraw int, marks []CheckMarks ) {
	var _color color.RGBA = color.RGBA{255, 255, 255, 0}
	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])
			if r-innerOverdraw/10> 0 {
				if len(marks) > i {
					_color = color.RGBA{220, 220, 220, 0}
					/*
					if math.Round(marks[i].Mean) > meanFindCircles  {
						_color = color.RGBA{0, 255, 0, 0}
					}else{
						_color = color.RGBA{220, 220, 220, 0}
					}
					*/
				}
				gocv.Circle(img, image.Pt(x, y), r-innerOverdraw/10, _color, outerOverdraw/10)
			}
		}
	}
}

func findCircles(croppedMat gocv.Mat , circleSize int,minDist float64) []bool {
	croppedMatGray := gocv.NewMat()
	gocv.CvtColor(croppedMat, &croppedMatGray, gocv.ColorBGRToGray)
	circles := gocv.NewMat()
	gocv.HoughCirclesWithParams(
		croppedMatGray,
		&circles,
		gocv.HoughGradient,
		dpHoughCircles,                     // dp
		minDist, //float64(croppedMatGray.Rows()/50), // minDist
		thresholdHoughCircles,                    // param1
		accumulatorThresholdHoughCircles,                    // param2
		circleSize,                    // minRadius
		circleSize,                     // maxRadius
	)


	imgRGray := gocv.NewMat()
	imgGray := gocv.NewMat()
	imgBlur := gocv.NewMat()
	gocv.CvtColor(croppedMat, &imgGray, gocv.ColorBGRToGray)
	if gaussianBlurFindCircles % 2!=1 {
		gaussianBlurFindCircles++
	}

	if adaptiveThresholdBlockSize % 2!=1 {
		adaptiveThresholdBlockSize++
	}

	gocv.GaussianBlur(imgGray, &imgBlur, image.Point{gaussianBlurFindCircles, gaussianBlurFindCircles}, 0, 0, gocv.BorderDefault)
	gocv.AdaptiveThreshold(imgBlur, &imgRGray, 255.0, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, adaptiveThresholdBlockSize, adaptiveThresholdSubtractMean)
	imgBlur.Close()
	imgGray.Close()

	checkMarks := []CheckMarks{}
	checkMarksList := []bool{}

	DrawCircles(&imgRGray, &circles,  innerOverdrawDrawCircles*int(pixelScale), outerOverdrawDrawCircles*int(pixelScale), checkMarks)

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])
			rect_circle:=image.Rect(x-r , y-r  , x+r , y+r )
			if rect_circle.Min.X < 0 || rect_circle.Min.Y < 0 || rect_circle.Max.X > imgRGray.Cols() || rect_circle.Max.Y > imgRGray.Rows() {
				continue
			}else{
				rect_circleMat := imgRGray.Region(rect_circle)
				mean := rect_circleMat.Mean()
				rect_circleMat.Close()
				checkMarks = append(checkMarks, CheckMarks{mean.Val1, x, y, r})
			}
		}
	}


	sort.Slice(checkMarks[:], func(i, j int) bool {
		return checkMarks[i].Y < checkMarks[j].Y
	})
	for i := 0; i < len(checkMarks); i++ {
		if math.Round(checkMarks[i].Mean) > meanFindCircles  {
			checkMarksList = append(checkMarksList, false)
		}else{
			checkMarksList = append(checkMarksList, true)
		}
	}


	imgCol := gocv.NewMat()
	gocv.CvtColor(imgRGray, &imgCol, gocv.ColorGrayToBGR)
	DrawCircles(&imgCol, &circles, innerOverdrawDrawCircles*int(pixelScale), outerOverdrawDrawCircles*int(pixelScale), checkMarks)
	if len(imageChannelCircle)==cap(imageChannelCircle) {
		mat,_ := <-imageChannelCircle
		mat.Close()
	}
	imageChannelCircle <- imgCol
	circles.Close()
	imgRGray.Close()

	return checkMarksList
}