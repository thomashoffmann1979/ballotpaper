package main

import (
	"fmt"
	"image"
	"log"
	"time"
	// "sort"
	"image/color"
	"gocv.io/x/gocv"
)


func IndexOf[T comparable](collection []T, el T) int {
    for i, x := range collection {
        if x == el {
            return i
        }
    }
    return -1
}

func process(img gocv.Mat, lastTesseract TesseractReturnType) ( TesseractReturnType) {
//	maxCircles := 0

	lastTesseract.IsCorrect=false
	lastTesseract.FCBarcode=""

	rotated := gocv.NewMat()
	defer rotated.Close()
	imgGray := gocv.NewMat()
	defer imgGray.Close()


	gocv.Rotate(img, &rotated, gocv.Rotate90CounterClockwise)
	var tesseractReturn TesseractReturnType
	var strBarcode string=""
	var minDist float64=1000


	if rotated.Empty() {
		return lastTesseract
	}

	if runContour {

		contour := findPaperContour(rotated)
		elapsed := time.Since(start)
		if boolVerbose {
			log.Printf("A Query took %s countourSize %d Title: %s", elapsed, contour.Size(), lastTesseract.Title)
		}
		if true {
			cornerPoints := getCornerPoints(contour)

			
			topLeftCorner := cornerPoints["topLeftCorner"]
			bottomRightCorner := cornerPoints["bottomRightCorner"]
			
			// 1190 1682
			// 1987 2525
			if boolVerbose {
				log.Printf("template: %d %d",  bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y )
			}
			// paper := extractPaper(rotated, contour, template.Cols()*2, template.Rows()*2, cornerPoints)
			paper := extractPaper(rotated, contour, bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y, cornerPoints)
			defer paper.Close()
			

			if paper.Empty() {
				return  lastTesseract
			}

			showImage("paper", paper, 0)

			if runScanner {
				strBarcode = scanBarcode(paper)
				tesseractReturn.FCBarcode=""
				if strBarcode == "" {
					strBarcode = lastTesseract.Barcode
				}
				if len(strBarcode) > 3 && strBarcode[0:3]=="FC4" {
					tesseractReturn.FCBarcode=strBarcode
				}
				if len(strBarcode) > 3 && strBarcode[0:3]=="FC3" {
					tesseractReturn.FCBarcode=strBarcode
				}
				if boolVerbose {
					if strBarcode != "" {
						log.Printf("Barcode: %s", strBarcode)
					}
				}
			}

			if tesseractReturn.FCBarcode==""{
				if runTesseract {
					if lastTesseract.Barcode==strBarcode && lastTesseract.Title!="" {
						tesseractReturn = lastTesseract
					}else{
						gocv.CvtColor(paper, &imgGray, gocv.ColorBGRToGray)
						tesseractReturn = tesseract(imgGray)
						tesseractReturn.Barcode = strBarcode
						if boolVerbose {
							fmt.Println(tesseractReturn);
						}
					}
					//point := tesseractReturn.Point
					circleSize := tesseractReturn.CircleSize
					if boolVerbose {
						fmt.Println("Ballot Paper Title: ", tesseractReturn.Title,tesseractReturn.PageRois, tesseractReturn.Point)
					}

					//fmt.Println("PageRois: ", tesseractReturn.PageRois)

					
					// if point.X!=0 && point.Y!=0{

						for pRoiIndex := 0; pRoiIndex < len(tesseractReturn.PageRois); pRoiIndex++ {

							if boolVerbose {
							log.Printf("Title: %s", tesseractReturn.Title)
							}
							if (IndexOf(tesseractReturn.PageRois[pRoiIndex].Titles, tesseractReturn.Title)>-1) {
									
								X := tesseractReturn.PageRois[pRoiIndex].X * paper.Cols() / tesseractReturn.Pagesize.Width
								Y := tesseractReturn.PageRois[pRoiIndex].Y * paper.Rows() / tesseractReturn.Pagesize.Height
								W := tesseractReturn.PageRois[pRoiIndex].Width * paper.Cols() / tesseractReturn.Pagesize.Width
								H := tesseractReturn.PageRois[pRoiIndex].Height * paper.Rows() / tesseractReturn.Pagesize.Height
								//rect:=image.Rect(point.X+X, point.Y+Y, point.X+X+W, point.Y+Y+H)

								circleSize = tesseractReturn.CircleSize * paper.Cols() / tesseractReturn.Pagesize.Width
								minDist =float64(tesseractReturn.CircleMinDistance) * float64(paper.Cols()) / float64(tesseractReturn.Pagesize.Width)
								
								if boolVerbose {
									log.Printf("circleSize: %d px %d mm, minDist: %d ", circleSize, tesseractReturn.CircleSize, minDist )
								}
								rect:=image.Rect( X, Y, X+W, Y+H)
								
								x := paper.Clone() 
								defer x.Close()
								croppedMat := x.Region(rect)
								defer croppedMat.Close()

								gocv.Rectangle(&paper, rect, color.RGBA{255, 255, 0, 0}, 4)

								if runMarkdetection {
									if !croppedMat.Empty() {
										marks:=findCircles(croppedMat, circleSize,minDist)
										//fmt.Println("checkMarksList: ", checkMarksList)
										if (tesseractReturn.Barcode!="") {
											if len(marks)== tesseractReturn.PageRois[pRoiIndex].ExcpectedMarks {
												// log.Printf("Barcode: %s, Title: %s, Marks: %v", tesseractReturn.Barcode , tesseractReturn.Title, marks)

												tesseractReturn.Marks=marks
												tesseractReturn.IsCorrect=true
											}else{
												// log.Printf("Barcode: %s, Title: %s, Marks: %v, Expected: %d", tesseractReturn.Barcode , tesseractReturn.Title, marks, tesseractReturn.PageRois[pRoiIndex].ExcpectedMarks)
											}
										}else{
											// log.Printf("Title: %s, Marks: %v", tesseractReturn.Title, marks)
										}
									}
								}
									

								

								drawContours := gocv.NewPointsVector()
								defer drawContours.Close()

								drawContours.Append(contour)
								gocv.DrawContours(&rotated, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)
							}
						}
					// }
				}
				
				showImage("rotated", rotated, 1)


				if rotated.Empty() {
					return  tesseractReturn
				}
			}
		}
	}

	showImage("output", rotated, 1)

	return  tesseractReturn
	
	
}