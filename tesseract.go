package main

import (
	"fmt"
	
	"image"
	"image/color"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)
func fileformatBytes(img gocv.Mat) []byte {
	buffer, err :=gocv.IMEncodeWithParams(gocv.PNGFileExt, img, []int{gocv.IMWriteJpegQuality, 100})
	if err != nil {
		return nil
	}
	return buffer.GetBytes(	)
}
func tesseract(img gocv.Mat) (image.Point,string) {
	result:=image.Point{0,	0}
	title:=""
	client := gosseract.NewClient()
	defer client.Close()

	croppedMat := img.Region(image.Rect(55, 50, 200, 90))

	if croppedMat.Empty() {
		return result,title
	}



	imgColor := gocv.NewMat()
	gocv.CvtColor(croppedMat, &imgColor, gocv.ColorGrayToBGR)
	client.SetWhitelist("EinzelhandelEnergie")
	seterror := client.SetImageFromBytes(fileformatBytes(croppedMat))
	if seterror != nil {
		fmt.Println(seterror)
		return result,title
	}
 	//client.SetImage("myimg.png")
	/*
	 _, err := client.Text()
	 if err != nil {
		 fmt.Println(err)
		 return result
	 }else{
		// fmt.Println(text)
	 }
	 */
	 //gosseract.PageIteratorLevel
	 out, herr := client.GetBoundingBoxes(3)
	 if herr != nil {
		fmt.Println(herr)
		return result,title
	}else{
		//fmt.Println(out[0].Word)


		if out[0].Word == "Einzelhandel" || out[0].Word == "Energie" {
			title = out[0].Word
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
	
		}



		croppedwindow := gocv.NewWindow("cropped")
		croppedwindow.IMShow(croppedMat)
	
	}
	return result,title

}