package main

import (
	"fmt"
	"image"
	"sort"
	"image/color"
	"gocv.io/x/gocv"
)

type CheckMarks struct {
	Mean float64
    X       int 
	Y       int
	Radius	   int
}


type CheckMarkList struct {
	Count int
	Sum int
	AVG float64
	Checked bool
}

func cameras(checkMarkList []CheckMarkList, lastTitle string) {
	
	webcam, err := gocv.VideoCaptureDevice(intCamera)
	if err != nil {
		fmt.Println("Error opening capture device: ", 0)
		return
	}
	defer webcam.Close()

	maxCircles := 0
	window := gocv.NewWindow("cameras")
	img := gocv.NewMat()
	defer img.Close()
	for {
		webcam.Read(&img)
		rotated := gocv.NewMat()
		gocv.Rotate(img, &rotated, gocv.Rotate90CounterClockwise)
		contour := findPaperContour(rotated)
		// fmt.Println("contours: ", contour.Size())

		cornerPoints := getCornerPoints(contour)
		// fmt.Println("cornerPoints: ", cornerPoints)

		// gocv.Rectangle(&rook, image.Rect(cornerPoints.topLeftCorner, bottomRightCorner), color.RGBA{255, 255, 0, 0}, -1)

		paper := extractPaper(rotated, contour, 500, 700, cornerPoints)

		if paper.Empty() {
			continue
		}

		imgGray := gocv.NewMat()
		gocv.CvtColor(paper, &imgGray, gocv.ColorBGRToGray)


		point,ballotPaperTitle:=tesseract(imgGray) // to do offeset
		if point.X!=0 && point.Y!=0{
			rect:=image.Rect(point.X+55-30, point.Y+50 +40-10, point.X+55+50-30, point.Y+50+500-10)
			x := paper.Clone() 
			croppedMat := x.Region(rect)
			gocv.Rectangle(&paper, rect, color.RGBA{255, 255, 0, 0}, -1)
			if !croppedMat.Empty() {
				

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

				blue := color.RGBA{0, 0, 255, 0}
				red := color.RGBA{255, 0, 0, 0}

			    // xw := gocv.NewWindow("spot")

				checkMarks := []CheckMarks{}

				checkMarksList := []bool{}
				
				if circles.Cols() != maxCircles {
					checkMarkList = []CheckMarkList{}
				}

				if lastTitle!=ballotPaperTitle { 
					checkMarkList = []CheckMarkList{}
				}

				maxCircles = circles.Cols()
				for i := 0; i < circles.Cols(); i++ {
					v := circles.GetVecfAt(0, i)
					// if circles are found
					if len(v) > 2 {
						x := int(v[0])
						y := int(v[1])
						r := int(v[2])


						if r-4 < 0 {
							continue
						}
						if y-r +1 < 0 || x-r +1 <0 || y+r-2 > croppedMat.Rows() || x+r-2 > croppedMat.Cols() {	
							continue
						}


						rect_circle:=image.Rect(x-r +1, y-r +1, x+r-1, y+r-1)
						// fmt.Println(  rect_circle )
						rect_circleMat := croppedMat.Region(rect_circle)

						imgGray := gocv.NewMat()

						/*
						imgGrayCroppedMat := gocv.NewMat()
						defer imgGrayCroppedMat.Close()
						gocv.CvtColor(croppedMat.Clone(), &imgGrayCroppedMat, gocv.ColorBGRToGray)
						gocv.Threshold(imgGrayCroppedMat, &imgGrayCroppedMat, 127, 200, gocv.ThresholdBinary+gocv.ThresholdOtsu)
						xw.IMShow(imgGrayCroppedMat)
						*/
						
						gocv.CvtColor(rect_circleMat, &imgGray, gocv.ColorBGRToGray)
						gocv.Threshold(imgGray, &imgGray, 30, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)

						mean := imgGray.Mean()
						
						defer rect_circleMat.Close()
						defer imgGray.Close()
						
						/*
						if mean.Val1>180 {
							fmt.Println( fmt.Sprintf(" Circle A: %d \t X,Y,R %d,%d,%d " , i, x, y, r  ))
							}else{
							fmt.Println( fmt.Sprintf(" Circle B: %d \t X,Y,R %d,%d,%d " , i, x, y, r  ))
						}
						*/
						checkMarks = append(checkMarks, CheckMarks{mean.Val1, x, y, r})


						//mask := gocv.Zeros(rect_circleMat.Rows(), rect_circleMat.Cols(), rect_circleMat.Type())
						
//						gocv.Circle(&mask, image.Pt(x, y), r, color.RGBA{255, 255, 255, 0}, -1)
//						gocv.IMWrite("mask.jpg", mask)
						//gocv.Subtract(mask,rect_circleMat,&rect_circleMat)

						//gocv.IMWrite("mask.jpg", mask)
						//name := fmt.Sprintf("circle_%d.jpg", i)
						//gocv.IMWrite(name, rect_circleMat)
						/*
						
						name2 := fmt.Sprintf("circle_%d_mask.jpg", i)
						gocv.IMWrite(name2, mask)
						*/
						
						//xw.IMShow(rect_circleMat)


						gocv.Circle(&croppedMat, image.Pt(x, y), r, blue, 2)

						gocv.Circle(&croppedMat, image.Pt(x, y), 2, red, 3)
					}


				}
				

				sort.Slice(checkMarks[:], func(i, j int) bool {
					return checkMarks[i].Y < checkMarks[j].Y
				})

				
				for i := 0; i < len(checkMarks); i++ {
					if i >= len(checkMarkList) {
						checkMarkList = append(checkMarkList, CheckMarkList{})
					}
					if checkMarks[i].Mean>150 {
						checkMarksList = append(checkMarksList, false)

						checkMarkList[i].Sum += 0
					}else{
						checkMarksList = append(checkMarksList, true)
						checkMarkList[i].Sum += 1
					}
					checkMarkList[i].Count++
					checkMarkList[i].AVG = float64(checkMarkList[i].Sum) / float64(checkMarkList[i].Count)
					checkMarkList[i].Checked = checkMarkList[i].AVG > 0.6
				}

				outList:=[]string{}
				for i := 0; i < len(checkMarkList); i++ {
					 
					if checkMarkList[i].Checked {
						outList = append(outList, "😎")
					} else {
						outList = append(outList, "🥶")
					}
				}						

				fmt.Println( ballotPaperTitle, outList , checkMarkList[0].Count );

				lastTitle = ballotPaperTitle

				roiwindow := gocv.NewWindow("roi")
				roiwindow.IMShow(croppedMat)
				
			}
		}

		paperwindow := gocv.NewWindow("paper")
		paperwindow.IMShow(paper)

		drawContours := gocv.NewPointsVector()
		defer drawContours.Close()
		drawContours.Append(contour)
		gocv.DrawContours(&rotated, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)

		if rotated.Empty() {
			continue
		}
		window.IMShow(rotated)
		window.WaitKey(1)
	}
}
