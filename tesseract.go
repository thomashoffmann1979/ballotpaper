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
func tesseract(img gocv.Mat) (image.Point,string,[]DocumentConfigurationPageRoi,DocumentConfigurationPageSize,int) {
	result:=image.Point{0,	0}
	title:=""
	pageRois:=[]DocumentConfigurationPageRoi{}
	pagesize:=DocumentConfigurationPageSize{}
	circleSize:=1
	client := gosseract.NewClient()
	defer client.Close()


	fmt.Println("tesseract",documentConfigurations, img.Cols(), img.Rows())
	for i := 0; i < len(documentConfigurations); i++ {
		
		pagesize=documentConfigurations[i].Pagesize
		circleSize=documentConfigurations[i].CircleSize
		X := documentConfigurations[i].TitleRegion.X * img.Cols() / pagesize.Width
		Y := documentConfigurations[i].TitleRegion.Y * img.Rows() / pagesize.Height
		W := documentConfigurations[i].TitleRegion.Width * img.Cols() / pagesize.Width
		H := documentConfigurations[i].TitleRegion.Height * img.Rows() / pagesize.Height

		croppedMat := img.Region(image.Rect(X, Y, W+X, H+Y))

		pageRois:=documentConfigurations[i].Rois

		if croppedMat.Empty() {
			return result,title,pageRois,pagesize,circleSize
		}



		imgColor := gocv.NewMat()
		gocv.CvtColor(croppedMat, &imgColor, gocv.ColorGrayToBGR)
		client.SetWhitelist("EinzelhandelEnergie")
		seterror := client.SetImageFromBytes(fileformatBytes(croppedMat))
		if seterror != nil {
			fmt.Println(seterror)
			return result,title,pageRois,pagesize,circleSize
		}
		out, herr := client.GetBoundingBoxes(3)
		if herr != nil {
			fmt.Println(herr)
			return result,title,pageRois,pagesize,circleSize
		}else{
			fmt.Println(out[0].Word)

			// documentConfigurations[i].Titles
			for j := 0; j < len(documentConfigurations[i].Titles); j++ {
				distance := levenshtein.ComputeDistance(out[0].Word, documentConfigurations[i].Titles[j])
				fmt.Printf("The distance between %s and %s is %d %d.\n", out[0].Word, documentConfigurations[i].Titles[j], len( documentConfigurations[i].Titles[j]), distance)
				if distance < 3 {
					title = documentConfigurations[i].Titles[j]

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
					result = out[0].Box.Min;
			
					return result,title,pageRois,pagesize,circleSize
				}

			}

			croppedwindow := gocv.NewWindow("cropped")
			croppedwindow.IMShow(croppedMat)
		
		}
	}
	
	

	return result,title,pageRois,pagesize,circleSize

}