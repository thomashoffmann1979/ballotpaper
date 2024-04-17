package main

import (
	// "fmt"
	"image"
	"gocv.io/x/gocv"
)

func extractPaper(img gocv.Mat, maxContour gocv.PointVector, resultWidth int, resultHeight int, cornerPoints map[string]image.Point) gocv.Mat {
	topLeftCorner := cornerPoints["topLeftCorner"]
	topRightCorner := cornerPoints["topRightCorner"]
	bottomLeftCorner := cornerPoints["bottomLeftCorner"]
	bottomRightCorner := cornerPoints["bottomRightCorner"]
	warpedDst := gocv.NewMat()
	if topLeftCorner != (image.Point{}) && topRightCorner != (image.Point{}) && bottomLeftCorner != (image.Point{}) && bottomRightCorner != (image.Point{}) {
		dsize := image.Point{resultWidth, resultHeight}
    newImg := []image.Point{
      image.Point{0, 0},
      image.Point{0, resultHeight},
      image.Point{resultWidth, resultHeight},
      image.Point{resultWidth, 0},
    }
    origImg := []image.Point{
      topLeftCorner, // top-left
      bottomLeftCorner, // bottom-left
      bottomRightCorner, // bottom-right
      topRightCorner,  // top-right
    }
    origV := gocv.NewPointVectorFromPoints(origImg)
    newV := gocv.NewPointVectorFromPoints(newImg)

		M := gocv.GetPerspectiveTransform( origV  , newV)
		gocv.WarpPerspective(img, &warpedDst, M, dsize)
    origV.Close()
    newV.Close()
    M.Close()
	}
	return warpedDst
}