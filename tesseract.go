package main

import (
	"fmt"
	
	"image"
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

type TesseractReturnType struct {
	Point    image.Point
	Title   string
	PageRois []DocumentConfigurationPageRoi
	Pagesize DocumentConfigurationPageSize
	CircleSize int
}

func tesseract(img gocv.Mat) (TesseractReturnType) {
	result:=TesseractReturnType{}
	result.Point=image.Point{0,	0}
	result.Title=""
	result.PageRois=[]DocumentConfigurationPageRoi{}
	result.Pagesize=DocumentConfigurationPageSize{}
	result.CircleSize=1
	
	client := gosseract.NewClient()
	defer client.Close()


	if boolVerbose {
		fmt.Println("tesseract",documentConfigurations, img.Cols(), img.Rows())
	}
	for i := 0; i < len(documentConfigurations); i++ {
		

		result.CircleSize=documentConfigurations[i].CircleSize
		result.Pagesize=documentConfigurations[i].Pagesize
		result.PageRois=documentConfigurations[i].Rois


		X := documentConfigurations[i].TitleRegion.X * img.Cols() / result.Pagesize.Width
		Y := documentConfigurations[i].TitleRegion.Y * img.Rows() / result.Pagesize.Height
		W := documentConfigurations[i].TitleRegion.Width * img.Cols() / result.Pagesize.Width
		H := documentConfigurations[i].TitleRegion.Height * img.Rows() / result.Pagesize.Height

		croppedMat := img.Region(image.Rect(X, Y, W+X, H+Y))


		if croppedMat.Empty() {
			return result
		}



		imgColor := gocv.NewMat()
		gocv.CvtColor(croppedMat, &imgColor, gocv.ColorGrayToBGR)
		client.SetWhitelist("EinzelhandelEnergie")
		seterror := client.SetImageFromBytes(fileformatBytes(croppedMat))
		if seterror != nil {
			fmt.Println(seterror)
			return result
		}
		out, herr := client.GetBoundingBoxes(3)
		if herr != nil {
			fmt.Println(herr)
			return result
		}else{
			if boolVerbose {
				fmt.Println(out[0].Word)
			}

			// documentConfigurations[i].Titles
			for j := 0; j < len(documentConfigurations[i].Titles); j++ {
				distance := levenshtein.ComputeDistance(out[0].Word, documentConfigurations[i].Titles[j])
				if boolVerbose {
					fmt.Printf("The distance between %s and %s is %d %d.\n", out[0].Word, documentConfigurations[i].Titles[j], len( documentConfigurations[i].Titles[j]), distance)
				}
				if distance < 3 {

					result.Title=documentConfigurations[i].Titles[j]
					//title = out[0].Word
					drawContours := gocv.NewPointsVector()
					defer drawContours.Close()
					contour:= gocv.NewPointVectorFromPoints([]image.Point{
						out[0].Box.Min,
						image.Point{out[0].Box.Max.X, out[0].Box.Min.Y},
						out[0].Box.Max,
						image.Point{out[0].Box.Min.X, out[0].Box.Max.Y} 			})
					drawContours.Append(contour)
					gocv.DrawContours(&croppedMat, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)

					result.Point = out[0].Box.Min;
			
					return result
				}

			}

			croppedwindow := gocv.NewWindow("cropped")
			croppedwindow.IMShow(croppedMat)
		
		}
	}
	
	

	return result

}