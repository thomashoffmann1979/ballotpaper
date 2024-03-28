package main

import (
	"fmt"
	"image/color"
	"gocv.io/x/gocv"
)

func cameras() {
	fmt.Println("cameras noting don yet here.")
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println("Error opening capture device: ", 0)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("cameras")
	img := gocv.NewMat()
	defer img.Close()
	for {
		webcam.Read(&img)
		rotated := gocv.NewMat()
		gocv.Rotate(img, &rotated, gocv.Rotate90CounterClockwise)
		contour := findPaperContour(rotated)
		fmt.Println("contours: ", contour.Size())

		cornerPoints := getCornerPoints(contour)
		fmt.Println("cornerPoints: ", cornerPoints)

		drawContours := gocv.NewPointsVector()
		defer drawContours.Close()
		drawContours.Append(contour)
		gocv.DrawContours(&rotated, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)

		window.IMShow(rotated)
		window.WaitKey(1)
	}
}
