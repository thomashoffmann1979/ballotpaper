package main

import (
	"log"
	//"gocv.io/x/gocv"
)

func processImage(){
	
	log.Println("processImage starting ")
	for {
		if !runVideo {
			break
		}
		log.Println("processImage",runVideo)
		//for range grabVideoCameraTicker.C {	
		img,ok := <-paperChannelImage
		if ok {
			log.Println("got image",ok,img.Size())
		}
		log.Println("processImage done",runVideo)
	}
	log.Println("processImage exit",runVideo)
}