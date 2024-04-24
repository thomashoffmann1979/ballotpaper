package main

import (
	"log"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"time"
	"github.com/bieber/barcode"
	"fmt"
	api "tualo.de/ballotpaper/api"
	"strings"
	"encoding/json"

	"encoding/base64"
)



func findBarcodes(scanner *barcode.ImageScanner, img gocv.Mat)[]BarcodeSymbol{
	syms := []BarcodeSymbol{}
	if img.Empty() {
		return syms
	}
	smaller:=gocv.NewMat()
	gocv.CvtColor(img, &smaller, gocv.ColorBGRToGray)
	if smaller.Cols() > 800 {
		gocv.GaussianBlur(smaller, &smaller, image.Point{5, 5}, 0, 0, gocv.BorderDefault)
		gocv.Resize(smaller, &smaller, image.Point{smaller.Cols() / barcodeScale, smaller.Rows() / barcodeScale}, 0, 0, gocv.InterpolationArea)
	}
	if false {
		log.Println("barcodeScale",barcodeScale)
	}
	symbols, err := scanner.ScanMat(&smaller)
	if err != nil {
		panic(err)
	}
	
	/*
	log.Println("findBarcodes",len(symbols))
	if len(symbols) == 0 {
		gocv.IMWrite("noBarcode.png",img)
	}else{
		gocv.IMWrite("barcode.png",img)
	
	}
	*/
	
	for _, s := range symbols {
		syms = append(syms,BarcodeSymbol{Type:s.Type.Name(),Data:s.Data,Quality:s.Quality,Boundary:s.Boundary})
		if false {
			log.Println("BarcodeSymbol",s.Type.Name(),s.Data,s.Quality,s.Boundary)
		}
	}
	smaller.Close()
	return syms
}



func processRegionsOfInterest(tr TesseractReturnType,img gocv.Mat, useRoi int) TesseractReturnType{
	

						
	pixelScale =  float64(img.Cols()) /  float64(tr.Pagesize.Width)
	pixelScaleY =  float64(img.Rows()) /  float64(tr.Pagesize.Height)

	if pixelScale==0 {
		pixelScale=1
	}
	if pixelScaleY==0 {
		pixelScaleY=1
	}
	circleSize := int(float64(tr.CircleSize) * pixelScale)
	minDist :=float64(tr.CircleMinDistance) * pixelScale

	if useRoi<len(tr.PageRois) {
		pRoiIndex := useRoi
		// for pRoiIndex := 0; pRoiIndex < len(tr.PageRois); pRoiIndex++ {
		X := int(float64(tr.PageRois[pRoiIndex].X) * pixelScale)
		Y := int(float64(tr.PageRois[pRoiIndex].Y) * pixelScaleY)
		W := int(float64(tr.PageRois[pRoiIndex].Width) * pixelScale)
		H := int(float64(tr.PageRois[pRoiIndex].Height) * pixelScaleY)

		rect:=image.Rect( X, Y, X+W, Y+H)
		croppedMat := img.Region(rect)
		if !croppedMat.Empty() {
			marks:=findCircles(croppedMat, circleSize,minDist )
			tr.Marks=marks
			/*
			tr.BoxBarcode= boxLabelWidget.Text
			tr.StackBarcode= stackLabelWidget.Text
			tr.Barcode= ballotLabelWidget.Text
			*/
			
			
			if tr.PageRois[pRoiIndex].ExcpectedMarks==len(marks) {
				tr.IsCorrect=true

				if false {
					log.Println("marks: ", marks)
				}
			}

			/*
			if lastTesseract.Title != rois.tesseractReturn.Title || lastTesseract.Barcode != rois.tesseractReturn.Barcode {
				checkMarkList = []CheckMarkList{}
			}
			checkMarkList = sumMarks(checkMarkList, rois.tesseractReturn)
			*/


			/*
			ret := rois;
			ret.mat = rois.mat.Clone()
			ret.tesseractReturn.BoxBarcode= boxLabelWidget.Text
			ret.tesseractReturn.StackBarcode= stackLabelWidget.Text
			ret.tesseractReturn.Barcode= ballotLabelWidget.Text
			
			if len(roisReturnChannel)==cap(roisReturnChannel) {
				cD,_ := <-roisReturnChannel
				cD.mat.Close()
			}
			roisReturnChannel <- ret
			*/

			
			// lastTesseract	= rois.tesseractReturn
		}
		croppedMat.Close()
	}
	return tr
	

}

func processImage(){
	scanner := barcode.NewScanner()
	scanner.SetEnabledAll(false)
	scanner.SetEnabledSymbology(barcode.Code39,true)
	scanner.SetEnabledSymbology(barcode.Code128,true)
	log.Println("processImage starting ")
	tesseractNeeded := true
	lastTesseractResult := TesseractReturnType{}
	lastBarcode := "wlekfjwuqezgzw"
	doFindCircles := false
	checkMarkList := []CheckMarkList{}
	green := 0
	red := 0
	blue := 0
	for {
		if !runVideo {
			break
		}
		start:=time.Now()
		if false {
			log.Println("processImage ************")
		}
		//for range grabVideoCameraTicker.C {	
		img,ok := <-paperChannelImage
		if ok {
			if false {
				log.Println("got image",ok,img.Size(),len(paperChannelImage))
			}

			green = 0
			red = 0
			blue = 0

			if !img.Empty() {
				contour := findPaperContour(img)

				if false {
					log.Println("findPaperContour done %s %v",time.Since(start),contour)
				}

				cornerPoints := getCornerPoints(contour)
				topLeftCorner := cornerPoints["topLeftCorner"]
				bottomRightCorner := cornerPoints["bottomRightCorner"]
				if false {
					log.Printf("template: %d %d",  bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y )
				}

				paper := extractPaper(img, contour, bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y, cornerPoints)
				
				if paper.Empty() {

					contour.Close()
					img.Close()
					continue
				}

				

				// mean := paper.Mean()
				// if (mean.Val1+mean.Val2+mean.Val3)/3 > 150 {
					area := float64(paper.Size()[0]) * float64(paper.Size()[1]) / float64(img.Size()[0]) / float64(img.Size()[1])
					// log.Println("mean",mean.Val1,mean.Val2,mean.Val3,area,paper.Size(),time.Since(start))
					if area > 0.1 {
						codes := findBarcodes(scanner,paper)
						if len(codes) > 0 {
							for _, code := range codes {
								if code.Type == "CODE-128" {

									if code.Data != lastBarcode {
										lastBarcode = code.Data
										log.Println("code",code)
										tesseractNeeded = true
										doFindCircles = false
										checkMarkList = []CheckMarkList{}

										green = 0
										red = 0
										blue = 255
									}

									if tesseractNeeded {
										result := tesseract(paper)
										if len(result.PageRois)>0 {
											tesseractNeeded = false
											lastTesseractResult = result
											doFindCircles = true
											checkMarkList = []CheckMarkList{}
											fmt.Println("lastTesseractResult",lastTesseractResult.Title)
											green = 255
											red = 255
											blue = 255

										}else{
											green = 100
											red = 255
											blue = 100
											fmt.Println("tesseract no bp found")
										}
									}

									if doFindCircles {


										for pRoiIndex := 0; pRoiIndex < len(lastTesseractResult.PageRois); pRoiIndex++ {
											titles := []string{}
				
											for i := 0; i < len(lastTesseractResult.PageRois[pRoiIndex].Types); i++ {
												titles = append(titles, lastTesseractResult.PageRois[pRoiIndex].Types[i].Title)
											}
											foundIndex := IndexOf(titles, lastTesseractResult.Title)
											if (foundIndex>-1) {

												res := processRegionsOfInterest(lastTesseractResult,paper,pRoiIndex)

												log.Println("res.Marks",res.Marks)
												log.Println("res.Id",lastTesseractResult.PageRois[pRoiIndex].Types[foundIndex].Id)
												green = 255
												red = 0
												blue = 255


												if res.IsCorrect {
													// log.Println("IsCorrect",res)
													//lastTesseractResult=res
													for i := 0; i < len(res.Marks); i++ {
														if i >= len(checkMarkList) {
															checkMarkList = append(checkMarkList, CheckMarkList{})
														}
														if res.Marks[i] {
															checkMarkList[i].Sum += 1
														}
														checkMarkList[i].Count++
														checkMarkList[i].AVG = float64(checkMarkList[i].Sum) / float64(checkMarkList[i].Count)
														checkMarkList[i].Checked = checkMarkList[i].AVG > sumMarksAVG
													}

													if len(checkMarkList)>0 && checkMarkList[0].Count>10 {
														//

														outList:=[]string{}
														for i := 0; i < len(checkMarkList); i++ {
															
															if checkMarkList[i].Checked {
																outList = append(outList, "X")
															} else {
																outList = append(outList, "O")
															}
														}	
														res.Barcode = lastBarcode
														fmt.Printf("Box: %s, Stack: %s, Barcode: %s, Title: %s, List: %v \n",res.BoxBarcode,res.StackBarcode, res.Barcode , lastTesseractResult.Title, outList)
														//checkMarkList = sumMarks(checkMarkList, res)
														b := new(strings.Builder)
														json.NewEncoder(b).Encode(outList)

														image_bytes, _ := gocv.IMEncode(gocv.JPEGFileExt, paper)
														image_base64 := base64.StdEncoding.EncodeToString(image_bytes.GetBytes())
														//fmt.Println("pic",image_base64[0:100])
														image_bytes.Close()

														res,err := api.SendReading(
															res.BoxBarcode,
															res.StackBarcode,
															res.Barcode,
															lastTesseractResult.PageRois[pRoiIndex].Types[foundIndex].Id,
															b.String(),
															"data:image/jpeg;base64,"+image_base64,
														)
														
														if err != nil {
															log.Println("SendReading ERROR",err)
															green = 0
															red = 255
															blue = 0
														}else{
															// log.Println(">>>>",res.Msg)
															if res.Success {
																doFindCircles = false


																green = 150
																red = 0
																blue = 50
															}else{
																log.Println("SendReading ERROR",res.Msg)
															}
														}

													}
												}
											}
										}
				

										

										

										/*
										marks:=findCircles(croppedMat, circleSize,minDist )

										circles := findCircles(paper,lastTesseractResult)
										if len(circles) > 0 {
											for _, circle := range circles {
												log.Println("circle",circle)
											}
										}
										*/
									}else{
										green = 150
										red = 0
										blue = 50
									}
									//log.Println("code use tesseract",code.Data,tesseractNeeded,lastTesseractResult)
								}
							}
							// gocv.IMWrite("paper.png",paper)
						}

					}else{
						green = 0
						red = 0
						blue = 0
				 	}

								

				drawContours := gocv.NewPointsVector()
				drawContours.Append(contour)
				gocv.DrawContours(&img, drawContours, -1, color.RGBA{uint8(red), uint8(green), uint8(blue), 120}, int(8.0*pixelScale))
				drawContours.Close()

				if len(imageChannelPaper)==cap(imageChannelPaper) {
					mat,_:=<-imageChannelPaper
					mat.Close()
				}
				cloned := img.Clone()
				imageChannelPaper <- cloned

				contour.Close()
				paper.Close()
			}
			img.Close()
		}
		//log.Println("processImage done %s",time.Since(start))
	}
	//log.Println("processImage exit",runVideo)
}