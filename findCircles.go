package main

import (
	"fmt"
	"image"
	// "sort"
	"image/color"
	"gocv.io/x/gocv"
)
func findCircles(croppedMat gocv.Mat , circleSize int) {
	croppedMatGray := gocv.NewMat()
	gocv.CvtColor(croppedMat, &croppedMatGray, gocv.ColorBGRToGray)
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


	imgRGray := gocv.NewMat()
	imgGray := gocv.NewMat()
	imgBlur := gocv.NewMat()
	defer imgGray.Close()
	gocv.CvtColor(croppedMat, &imgGray, gocv.ColorBGRToGray)

	mSize := (circleSize%2) +  circleSize
	fmt.Println("mSize: ", mSize)
	gocv.GaussianBlur(imgGray, &imgBlur, image.Point{5, 5}, 0, 0, gocv.BorderDefault)
	gocv.AdaptiveThreshold(imgBlur, &imgRGray, 255.0, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 7, 3.0)
	

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])
			fmt.Println("x,y,r: ", x, y, r)
			if r-4 < 0 {
				continue
			}
			if y-r +1 < 0 || x-r +1 <0 || y+r-2 > croppedMat.Rows() || x+r-2 > croppedMat.Cols() {	
				continue
			}
			// rect_circle:=image.Rect(x-r +1, y-r +1, x+r-1, y+r-1)
			// rect_circleMat := croppedMat.Region(rect_circle)
			// mean := imgGray.Mean()
			// defer rect_circleMat.Close()


			color := color.RGBA{255, 255, 255, 0}
			gocv.Circle(&imgRGray, image.Pt(x, y), r-2, color, 4)


			r+=2
			gocv.Circle(&imgRGray, image.Pt(x, y), r, color, 2) 
			r+=2
			gocv.Circle(&imgRGray, image.Pt(x, y), r, color, 2) 
			r+=2
			gocv.Circle(&imgRGray, image.Pt(x, y), r, color, 2) 
			r+=2
			gocv.Circle(&imgRGray, image.Pt(x, y), r, color, 2) 
			r+=2
			gocv.Circle(&imgRGray, image.Pt(x, y), r, color, 2) 
			
			

			rect_circle:=image.Rect(x-r +1, y-r +1, x+r-1, y+r-1)
			rect_circleMat := imgRGray.Region(rect_circle)
			/*imgCGray := gocv.NewMat()
			defer imgCGray.Close()
			gocv.CvtColor(croppedMat, &imgCGray, gocv.ColorBGRToGray)
			*/
			mean := rect_circleMat.Mean()
			defer rect_circleMat.Close()
			fmt.Println("mean: ", mean)


		}


	}

	
	// gocv.Threshold(imgGray, &imgRGray, 40, 255, gocv.ThresholdBinary + gocv.ThresholdOtsu)
			

	findCirclesWindow := gocv.NewWindow("findCircles")
	findCirclesWindow.IMShow(imgRGray)

}