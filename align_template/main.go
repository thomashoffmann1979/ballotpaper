package main

import (
	"fmt"
	"flag"

	"image/color"
	"image"
	"sort"
	"gocv.io/x/gocv"
)

var (
	strInputFile string
	strCompareFile string
	points = []image.Point{}
)
// https://pyimagesearch.com/2020/08/31/image-alignment-and-registration-with-opencv/

/*
image: Our input photo/scan of a form (such as the IRS W-4). The form itself, from an arbitrary viewpoint, should be identical to the template image but with form data present.
template: The template form image.
maxFeatures: Places an upper bound on the number of candidate keypoint regions to consider.
keepPercent: Designates the percentage of keypoint matches to keep, effectively allowing us to eliminate noisy keypoint matching results
debug: A flag indicating whether to display the matched keypoints. By default, keypoints are not displayed; however I recommend setting this value to True for debugging purposes.	
*/

type Point struct {
	X, Y int
}

func align_images(use_image gocv.Mat, template gocv.Mat, maxFeatures  int  , keepPercent float64 , debug bool ) {

		
			// convert both the input image and template to grayscale
	imageGray := gocv.NewMat()
	templateGray := gocv.NewMat()
	gocv.CvtColor(use_image, &imageGray, gocv.ColorBGRToGray)
	gocv.CvtColor(template, &templateGray, gocv.ColorBGRToGray)


	// detect ORB keypoints and descriptors in the grayscale images
	// orb := gocv.NewORB()
	orb := gocv.NewORBWithParams(maxFeatures, 1.2, 8, 31, 0, 2, 0, 31, 20)
	defer orb.Close()
	kpsA, descsA := orb.DetectAndCompute(imageGray, gocv.NewMat())
	kpsB, descsB := orb.DetectAndCompute(templateGray, gocv.NewMat())

	if debug {
		fmt.Println("keypointsA: ", kpsA)
		fmt.Println("keypointsB: ", kpsB)
		fmt.Println("descsA: ", descsA)
		fmt.Println("descsB: ", descsB)
	}

	matcher:= gocv.NewBFMatcherWithParams(gocv.NormHamming, false)
	defer matcher.Close()
	matches := matcher.Match(descsA, descsB)
	//.KnnMatch(descsA, descsB, 2)

	compareDistance := func (a int, b int) bool {
		return matches[a].Distance < matches[b].Distance
	}

	keep:= int( float64(len(matches)) * keepPercent)

	sort.Slice(matches, compareDistance)

	//keep=2
	matches = matches[:keep]
	sort.Slice(matches[:], func(i, j int) bool {
		return i<keep
	})
	
	fmt.Println("matches: ", matches)
	fmt.Println("len(matches): ", len(matches))
	fmt.Println("keep: ", keep)

	output := gocv.NewMat()
	gocv.DrawMatches(use_image, kpsA, template, kpsB, matches, &output, color.RGBA{0, 255, 0, 0},color.RGBA{0, 0, 255, 0}, nil, gocv.DrawDefault)


	gocv.IMWrite("draw_matches.jpg", output)


	// initialize the list of actual matches
	//matches := []
	// Zeros := gocv.NewMatWithSize(1, 2, gocv.MatTypeCV32F)

	ptsA := gocv.NewMatWithSize(len(matches), 2, gocv.MatTypeCV32F)
	ptsB := gocv.NewMatWithSize(len(matches), 2, gocv.MatTypeCV32F)
	for i := 0; i < len(matches); i++ {
		
		ptsA.SetFloatAt(i, 0, float32(kpsA[matches[i].QueryIdx].X))
		ptsA.SetFloatAt(i, 1, float32(kpsA[matches[i].QueryIdx].Y))
		ptsB.SetFloatAt(i, 0, float32(kpsB[matches[i].TrainIdx].X))
		ptsB.SetFloatAt(i, 1, float32(kpsB[matches[i].TrainIdx].Y))
	}

	if debug {
		fmt.Println("ptsA: ", ptsA)
		fmt.Println("ptsB: ", ptsB)
	}

	// compute the homography between the two sets of points
	//func FindHomography(ptsA, dstPoints *Mat, method HomographyMethod, ransacReprojThreshold float64, mask *Mat, maxIters int, confidence float64) Mat

	mask := gocv.NewMat()
	h := gocv.FindHomography(ptsB, &ptsA, gocv.HomograpyMethodRANSAC, 3, &mask, 2000, 0.995)

	if debug {
		fmt.Println("mask: ", mask)
		fmt.Println("h: ", h)
	}


	warpedDst := gocv.NewMat()
	defer warpedDst.Close()

	dsize := image.Point{template.Cols(), template.Rows()}
	gocv.WarpPerspective(use_image, &warpedDst, h, dsize )

	gocv.IMWrite("warped.jpg", warpedDst)
	fmt.Println("warpedDst: ", warpedDst)
	fmt.Println("dsize: ", dsize)

}


func main() {

	flag.StringVar(&strInputFile, "input", "" , "input file")
	flag.StringVar(&strCompareFile, "template", "" , "template file")

	flag.Parse()

	image := gocv.IMRead(strInputFile, gocv.IMReadColor)
	template := gocv.IMRead(strCompareFile, gocv.IMReadColor)
	rotated := gocv.NewMat()
	gocv.Rotate(image, &rotated, gocv.Rotate90CounterClockwise)
	contour := findPaperContour(rotated)
	cornerPoints := getCornerPoints(contour)
	paper := extractPaper(rotated, contour, 500, 700, cornerPoints)

	align_images(paper, template, 500, 0.2, false)

	fmt.Println( "done!")
}