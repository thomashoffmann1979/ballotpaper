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
  //defer warpedDst.Close()

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

		M := gocv.GetPerspectiveTransform( gocv.NewPointVectorFromPoints(origImg)  , gocv.NewPointVectorFromPoints(newImg))
    defer M.Close()
		gocv.WarpPerspective(img, &warpedDst, M, dsize)
    // , gocv.InterpolationLinear, gocv.BorderConstant, gocv.NewScalar())
	}

	
	return warpedDst
}
/*
extractPaper: function(img,maxContour, resultWidth, resultHeight, cornerPoints) {


      
            const {
              topLeftCorner,
              topRightCorner,
              bottomLeftCorner,
              bottomRightCorner,
            } = cornerPoints || this.getCornerPoints(maxContour, img);
            let warpedDst = new cv.Mat();
     
            if (
              
              topLeftCorner && 
              topRightCorner && 
              bottomLeftCorner && 
              bottomRightCorner
              
              ){
                let dsize = new cv.Size(resultWidth, resultHeight);
                let srcTri = cv.matFromArray(4, 1, cv.CV_32FC2, [
                topLeftCorner.x,
                topLeftCorner.y,
                topRightCorner.x,
                topRightCorner.y,
                bottomLeftCorner.x,
                bottomLeftCorner.y,
                bottomRightCorner.x,
                bottomRightCorner.y,
                ]);
        
                let dstTri = cv.matFromArray(4, 1, cv.CV_32FC2, [
                0,
                0,
                resultWidth,
                0,
                0,
                resultHeight,
                resultWidth,
                resultHeight,
                ]);
        
                let M = cv.getPerspectiveTransform(srcTri, dstTri);
                cv.warpPerspective(
                img,
                warpedDst,
                M,
                dsize,
                cv.INTER_LINEAR,
                cv.BORDER_CONSTANT,
                new cv.Scalar()
                );
            }
            //cv.imshow(canvas, warpedDst);
      
            //img.delete()
            //warpedDst.delete()
            return warpedDst;
          },
*/