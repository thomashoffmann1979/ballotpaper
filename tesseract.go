package main

import (
	"fmt"
	"time"
	"image"
	"os"
	"image/color"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
	"github.com/agnivade/levenshtein"
)
func fileformatBytes(img gocv.Mat) []byte {
	buffer, err :=gocv.IMEncodeWithParams(gocv.PNGFileExt, img, []int{gocv.IMWriteJpegQuality, 100})
	if err != nil {
		return nil
	}
	return buffer.GetBytes(	)
}



func tesseract(img gocv.Mat) (TesseractReturnType) {

	start := time.Now()



	result:=TesseractReturnType{}
	result.Point=image.Point{0,	0}
	result.Title=""
	result.IsCorrect=false
	result.PageRois=[]DocumentConfigurationPageRoi{}
	result.Pagesize=DocumentConfigurationPageSize{}
	result.CircleSize=1
	result.CircleMinDistance=100
	result.Marks=[]bool{}
	
	client := gosseract.NewClient()
	defer client.Close()

	if tesseractPrefix != "" {
		client.SetTessdataPrefix(tesseractPrefix)
	}


	if false {
		fmt.Println("tesseract",documentConfigurations, img.Cols(), img.Rows())
	}
	for i := 0; i < len(documentConfigurations); i++ {
		

		result.CircleSize=documentConfigurations[i].CircleSize
		result.CircleMinDistance=documentConfigurations[i].CircleMinDistance
		result.Pagesize=documentConfigurations[i].Pagesize
		


		X := documentConfigurations[i].TitleRegion.X * img.Cols() / result.Pagesize.Width
		Y := documentConfigurations[i].TitleRegion.Y * img.Rows() / result.Pagesize.Height
		W := documentConfigurations[i].TitleRegion.Width * img.Cols() / result.Pagesize.Width
		H := documentConfigurations[i].TitleRegion.Height * img.Rows() / result.Pagesize.Height

		croppedMat := img.Region(image.Rect(X, Y, W+X, H+Y))

		if croppedMat.Empty() {
			croppedMat.Close()
			return result
		}



		// imgColor := gocv.NewMat()
		// gocv.CvtColor(croppedMat, &imgColor, gocv.ColorGrayToBGR)
		// client.SetWhitelist("EinzelhandelEnergie")
		smaller := ResizeMat(croppedMat.Clone(), croppedMat.Cols()/tesseractScale, croppedMat.Rows()/tesseractScale)

		seterror := client.SetImageFromBytes(fileformatBytes(smaller))
		if seterror != nil {
			fmt.Println(seterror)
			return result
		}
		out, herr := client.GetBoundingBoxes(3)
		if herr != nil {
			fmt.Println(herr)
			croppedMat.Close()
			return result
		}else{
			if false {
				for j := 0; j < len(out); j++ {
					fmt.Println(i,out[j].Word)
				}
			}
			for j := 0; j < len(documentConfigurations[i].Titles); j++ {
				distance := levenshtein.ComputeDistance(out[0].Word, documentConfigurations[i].Titles[j])
				if false {
					fmt.Printf("The distance between %s and %s is %d %d.\n", out[0].Word, documentConfigurations[i].Titles[j], len( documentConfigurations[i].Titles[j]), distance)
				}
				if distance < 3 {
					result.Title=documentConfigurations[i].Titles[j]
					//title = out[0].Word
					drawContours := gocv.NewPointsVector()
					contour:= gocv.NewPointVectorFromPoints([]image.Point{
						out[0].Box.Min,
						image.Point{out[0].Box.Max.X, out[0].Box.Min.Y},
						out[0].Box.Max,
						image.Point{out[0].Box.Min.X, out[0].Box.Max.Y} 			})
					drawContours.Append(contour)
					gocv.DrawContours(&croppedMat, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)
					result.Point = image.Point{documentConfigurations[i].TitleRegion.X, documentConfigurations[i].TitleRegion.Y}

					if false {
						debug( fmt.Sprintf("ocr %s %d %d %d",time.Since(start),croppedMat.Cols(),croppedMat.Rows(), os.Getpid() ) )
					}
					result.PageRois=documentConfigurations[i].Rois
					croppedMat.Close()
					drawContours.Close()
					smaller.Close()
					return result
				}

			}

			
		
		}
		croppedMat.Close()
		smaller.Close()
	}
	
	//debug( fmt.Sprintf("tesseract failed %s ",time.Since(start)) )
	readyToSave = false

	return result

}