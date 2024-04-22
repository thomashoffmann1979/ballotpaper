package main

import (
	"fmt"
	"log"
	"image"
	//"image/color"
	//"gocv.io/x/gocv"
)

var readyToSave bool = false
func sumMarks(checkMarkList []CheckMarkList, processResult TesseractReturnType) []CheckMarkList{
	if processResult.IsCorrect {
		if len(checkMarkList)==0 {
			
		}else{
			if len(checkMarkList) != len(processResult.Marks) {
				checkMarkList = []CheckMarkList{}
				log.Println("reset checkMarkList")
			}
		}

		
		// log.Printf("Box: %s, Stack: %s, Barcode: %s, Title: %s, Marks: %v",processResult.BoxBarcode,processResult.StackBarcode, processResult.Barcode , processResult.Title, processResult.Marks)

		for i := 0; i < len(processResult.Marks); i++ {
			if i >= len(checkMarkList) {
				checkMarkList = append(checkMarkList, CheckMarkList{})
			}
			if processResult.Marks[i] {
				checkMarkList[i].Sum += 1
			}
			checkMarkList[i].Count++
			checkMarkList[i].AVG = float64(checkMarkList[i].Sum) / float64(checkMarkList[i].Count)
			checkMarkList[i].Checked = checkMarkList[i].AVG > sumMarksAVG
		}


		outList:=[]string{}
		for i := 0; i < len(checkMarkList); i++ {
			
			if checkMarkList[i].Checked {
				outList = append(outList, "😎")
			} else {
				outList = append(outList, "🥶")
			}
		}	


		if len(checkMarkList)>0 && checkMarkList[0].Count>15 {
			log.Printf("Box: %s, Stack: %s, Barcode: %s, Title: %s, Marks: %s",processResult.BoxBarcode,processResult.StackBarcode, processResult.Barcode , processResult.Title, outList)
			readyToSave = true
		}else{
			readyToSave = false
		}
		// checkMarkList = []CheckMarkList{}
	}
	return checkMarkList
}

func processRoisChannel() {
	var circleSize int=100
	var minDist float64=1000
	var checkMarkList = []CheckMarkList{}
	lastTesseract := TesseractReturnType{}

	for range grabVideoCameraTicker.C {	
		roisReturn,ok := <-tesseractReturnChannel
		if ok {
			for pRoiIndex := 0; pRoiIndex < len(roisReturn.tesseractReturn.PageRois); pRoiIndex++ {

				if false {
					log.Printf("Title: %s", roisReturn.tesseractReturn.Title)
				}
				if (IndexOf(roisReturn.tesseractReturn.PageRois[pRoiIndex].Titles, roisReturn.tesseractReturn.Title)>-1) {
						
					pixelScale =  float64(roisReturn.mat.Cols()) /  float64(roisReturn.tesseractReturn.Pagesize.Width)
					pixelScaleY =  float64(roisReturn.mat.Rows()) /  float64(roisReturn.tesseractReturn.Pagesize.Height)

					if pixelScale==0 {
						pixelScale=1
					}
					if pixelScaleY==0 {
						pixelScaleY=1
					}


					log.Printf("pixelScale: %.1f , %.1f  ",pixelScale, pixelScaleY);
					debug(fmt.Sprintf("pixelScale: %.1f , %.1f  ",pixelScale, pixelScaleY))


					X := int(float64(roisReturn.tesseractReturn.PageRois[pRoiIndex].X) * pixelScale)
					Y := int(float64(roisReturn.tesseractReturn.PageRois[pRoiIndex].Y) * pixelScaleY)
					W := int(float64(roisReturn.tesseractReturn.PageRois[pRoiIndex].Width) * pixelScale)
					H := int(float64(roisReturn.tesseractReturn.PageRois[pRoiIndex].Height) * pixelScaleY)
					//rect:=image.Rect(point.X+X, point.Y+Y, point.X+X+W, point.Y+Y+H)

					circleSize = int(float64(roisReturn.tesseractReturn.CircleSize) * pixelScale)
					minDist =float64(roisReturn.tesseractReturn.CircleMinDistance) * pixelScale
					

					if false {
						log.Printf("circleSize: %d px %d mm, minDist: %d ", circleSize, roisReturn.tesseractReturn.CircleSize, minDist )
					}
					rect:=image.Rect( X, Y, X+W, Y+H)
					
					croppedMat := roisReturn.mat.Region(rect)
					if !croppedMat.Empty() {


						marks:=findCircles(croppedMat, circleSize,minDist )
						roisReturn.tesseractReturn.Marks=marks
						roisReturn.tesseractReturn.BoxBarcode= boxLabelWidget.Text
						roisReturn.tesseractReturn.StackBarcode= stackLabelWidget.Text
						roisReturn.tesseractReturn.Barcode= ballotLabelWidget.Text
						// log.Println("marks: ", marks)
						
						if roisReturn.tesseractReturn.PageRois[pRoiIndex].ExcpectedMarks==len(marks) {
							roisReturn.tesseractReturn.IsCorrect=true
						}

						if lastTesseract.Title != roisReturn.tesseractReturn.Title || lastTesseract.Barcode != roisReturn.tesseractReturn.Barcode {
							checkMarkList = []CheckMarkList{}
						}
						checkMarkList = sumMarks(checkMarkList, roisReturn.tesseractReturn)



						ret := roisReturn;
						ret.mat = roisReturn.mat.Clone()
						ret.tesseractReturn.BoxBarcode= boxLabelWidget.Text
						ret.tesseractReturn.StackBarcode= stackLabelWidget.Text
						ret.tesseractReturn.Barcode= ballotLabelWidget.Text
						if len(roisReturnChannel)==cap(roisReturnChannel) {
							cD,_ := <-roisReturnChannel
							cD.mat.Close()
						}
						roisReturnChannel <- ret

						
						lastTesseract	= roisReturn.tesseractReturn
						/*
						if (roisReturn.tesseractReturn.Barcode!="") {
							if len(marks)== roisReturn.tesseractReturn.PageRois[pRoiIndex].ExcpectedMarks {
								log.Printf("Barcode: %s, Title: %s, Marks: %v", roisReturn.tesseractReturn.Barcode , roisReturn.tesseractReturn.Title, marks)

								roisReturn.tesseractReturn.Marks=marks
								roisReturn.tesseractReturn.IsCorrect=true
							}else{
								// log.Printf("Barcode: %s, Title: %s, Marks: %v, Expected: %d", tesseractReturn.Barcode , tesseractReturn.Title, marks, tesseractReturn.PageRois[pRoiIndex].ExcpectedMarks)
							}
						}else{
							// log.Printf("Title: %s, Marks: %v", tesseractReturn.Title, marks)
						}
						*/
					}
					croppedMat.Close()

					/*
					gocv.Rectangle(&paper, rect, color.RGBA{255, 255, 0, 0}, 4)

					drawContours := gocv.NewPointsVector()
					defer drawContours.Close()

					drawContours.Append(contour)
					gocv.DrawContours(&rotated, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)
					*/
				}
			}
			roisReturn.mat.Close()
		}
	}
}
