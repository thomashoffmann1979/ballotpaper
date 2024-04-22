package main

func debug(str string) {
	//log.Printf("pixelScale: %.1f , %.1f  ",pixelScale, pixelScaleY);
	if len(debugChannel)==cap(debugChannel) {
		<-debugChannel
	}
	debugChannel <- str
}

