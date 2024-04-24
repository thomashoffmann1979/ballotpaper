package main

import (
	"image"
	"gocv.io/x/gocv"
)

type TesseractReturnType struct {
	Point    image.Point
	BoxBarcode   string
	StackBarcode   string
	Barcode   string
	Title   string
	Id string
	IsCorrect bool
	Marks   []bool
	PageRois []DocumentConfigurationPageRoi
	Pagesize DocumentConfigurationPageSize
	CircleSize int
	CircleMinDistance int
}

type RoisChannelStruct struct {
	tesseractReturn TesseractReturnType
	mat gocv.Mat
}

type CheckMarkList struct {
	Count int
	Sum int
	AVG float64
	Checked bool
}

type CheckMarks struct {
	Mean float64
    X       int 
	Y       int
	Radius	   int
}


type ReturnType struct {
	Point    image.Point
	FCBarcode   string
	Barcode   string
	Title   string
	IsCorrect bool
	Marks   []bool
	PageRois []DocumentConfigurationPageRoi
	Pagesize DocumentConfigurationPageSize
	CircleSize int
	CircleMinDistance int
}

type CameraList struct {
	Width int
	Height int
	Index int
	Title string
}