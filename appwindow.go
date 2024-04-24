package main

import (
	"os"
	"fmt"
	"time"
	// "math"
	"image"
	"image/color"
	"log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/cmd/fyne_demo/data"
	//"fyne.io/fyne/v2/cmd/fyne_demo/tutorials"
	//"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/driver/desktop"
	// "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"


	"tualo.de/ballotpaper/data"
	api "tualo.de/ballotpaper/api"
	"gocv.io/x/gocv"



)


const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

var loginContainer *fyne.Container
var mainAppContainer fyne.CanvasObject
var pingResponse api.PingResponse
var kandidatenResponse api.KandidatenResponse

var fullNameWidget *widget.Label
var boxLabelWidget *widget.Label
var stackLabelWidget *widget.Label
var ballotLabelWidget *widget.Label

var cameraSelectWidget *widget.Select
var thresholdHoughCirclesWidget *widget.Slider
var meanFindCirclesWidget *widget.Slider
var dpHoughCirclesWidget *widget.Slider
var gaussianBlurFindCirclesWidget *widget.Slider
var adaptiveThresholdBlockSizeWidget *widget.Slider
var adaptiveThresholdSubtractMeanWidget *widget.Slider

var thresholdHoughCirclesWidgetLabel *widget.Label
var meanFindCirclesWidgetLabel *widget.Label
var dpHoughCirclesWidgetLabel *widget.Label
var gaussianBlurFindCirclesWidgetLabel *widget.Label
var adaptiveThresholdBlockSizeWidgetLabel *widget.Label
var adaptiveThresholdSubtractMeanWidgetLabel *widget.Label

var debugListWidget *widget.List

var outputImage *canvas.Image
var paperImage *canvas.Image
var circleImage *canvas.Image

var videoIsRunning bool = false
var grabVideoCameraTicker *time.Ticker
var cameraChannelImage chan gocv.Mat

var paperChannelImage = make(chan gocv.Mat, 1)
var readyToSaveChannelImage = make(chan gocv.Mat, 1)
var tesseractChannelImage = make(chan gocv.Mat, 1)
var tesseractReturnChannel = make(chan RoisChannelStruct, 1)


var debugChannel = make(chan string, 10)
var debugData = []string{"", "", "", "", "", "", "", "", "", ""}


var roisReturnChannel = make(chan RoisChannelStruct, 1)


var imageChannelPaper chan gocv.Mat
var imageChannelCircle chan gocv.Mat

var currentBoxChannel chan string
var currentStackChannel chan string
var currentBarcodeChannel chan string

func matToImage(mat gocv.Mat) image.Image {
	img, _ := mat.ToImage()
	return img
}

var window *gocv.Window


func grabVideoImage() {
	for range grabVideoCameraTicker.C {
		mat,ok := <-cameraChannelImage
		if ok {


			start := time.Now()
			if false {
				image := matToImage(mat)
				outputImage.Image = image
				outputImage.Refresh()
			}
			if false {
				fmt.Println("grabVideoImage time",time.Since(start))
			}
			
			mat.Close()
		}
    }
}


func grabChannelBarcodes() {
	for range grabVideoCameraTicker.C {
		syms,ok := <-scannerChannelBarcodes
		if ok {
			for _,sym := range syms {
				// log.Println("got barcode",sym.Type,sym.Data)
				if sym.Type == "CODE-128" {
					ballotLabelWidget.SetText("Stimmzettel: "+sym.Data)

					/*
					a := ballotLabelWidget.NewColorRGBAAnimation(theme.PrimaryColorNamed(theme.ColorBlue), theme.PrimaryColorNamed(theme.ColorGreen),
						time.Second*1, func(c color.Color) {
							ballotLabelWidget.TextStyle.Color = c
							ballotLabelWidget.Refresh()
						})
					a.RepeatCount = fyne.AnimationRepeatForever
					a.AutoReverse = true
					a.Start()
					*/
					
				}
				if sym.Type == "CODE-39" {
					data := sym.Data
					if len(data)>3 {
						if data[0:3]=="FC4" {
							boxLabelWidget.SetText("Kiste: "+data)
						}
						if data[0:3]=="FC3" {
							stackLabelWidget.SetText("Stapel: "+data)
						}
					}
				}
			}
		}
    }
}

func grabReadyToSaveImage() {
	for range grabVideoCameraTicker.C {
		mat,ok := <-readyToSaveChannelImage
		if ok {
			if readyToSave {
				//fileName := fmt.Sprintf("outimages/%s.jpg", ballotLabelWidget.Text)
				// ballotLabelWidget.GetText()
				//gocv.IMWrite(fileName, mat)
				log.Println("readyToSave") // MatProfile %d",gocv.MatProfile.Count())
			}
			mat.Close()
		}
	}
}

func grabPaperImage() {
	for range grabVideoCameraTicker.C {
		mat,ok := <-imageChannelPaper
		if ok {

			image := matToImage(mat)
			paperImage.Image = image
			paperImage.Refresh()
			mat.Close()


		}
    }
}

func grabCircleImage() {
	for range grabVideoCameraTicker.C {
		mat,ok := <-imageChannelCircle
		if ok {
			image := matToImage(mat)
			circleImage.Image = image
			circleImage.Refresh()
			mat.Close()
		}
    }
}

func grabCurrentBox() {
	for range grabVideoCameraTicker.C {
		data,ok := <-currentBoxChannel
		if ok {
			boxLabelWidget.SetText("Kiste: "+data)
		}
    }
}

func grabCurrentStack() {
	for range grabVideoCameraTicker.C {
		data,ok := <-currentStackChannel
		if ok {
			stackLabelWidget.SetText("Stapel: " +data)
		}
    }
}


func grabDebugs() {
	for range grabVideoCameraTicker.C {
		data,ok := <-debugChannel
		if ok {
			for i:=len(debugData)-1;i>0;i-- {
				debugData[i] = debugData[i-1]
			}
			debugData[0] = data
			debugListWidget.Refresh()
			// debugListWidget.Add(data)
		}
	}
}


func updateDebugData(i widget.ListItemID, o fyne.CanvasObject) {
	o.(*widget.Label).SetText(debugData[i])
}

func makeDebugList() *widget.List {
	return widget.NewList(
		func() int {
			return len(debugData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		updateDebugData,
	)
}

func makeVideoGrid() fyne.CanvasObject {

	outputImage = canvas.NewImageFromImage(matToImage(gocv.NewMatWithSize(640, 480, gocv.MatTypeCV8UC3)))
	outputImage.FillMode = canvas.ImageFillContain
	if !showOutputImage {
		outputImage.Hide()
	}

	paperImage = canvas.NewImageFromImage(matToImage(gocv.NewMatWithSize(640, 480, gocv.MatTypeCV8UC3)))
	paperImage.FillMode = canvas.ImageFillContain
	if !showPaperImage {
		paperImage.Hide()
	}

	circleImage = canvas.NewImageFromImage(matToImage(gocv.NewMatWithSize(640, 480, gocv.MatTypeCV8UC3)))
	circleImage.FillMode = canvas.ImageFillContain
	if !showCirlceImage {
		circleImage.Hide()
	}


	debugListWidget = makeDebugList()
	if !showDebugList {
		debugListWidget.Hide()
	}
	
	return container.New(
		layout.NewGridLayout(3), 
		outputImage, 
		paperImage, 
		circleImage,
		debugListWidget,
	)
}

func onCameraSelectWidget(item string) { 
	fmt.Println("Select",item,IndexOf(cameraSelectWidget.Options,item)) 
	intCamera = IndexOf(cameraSelectWidget.Options,item)
	
}

func makeSettingsForm() fyne.CanvasObject {

	cameraList := getCameraList()
	fmt.Println("maxcameranum",len(cameraList))	
	//"Camera 1", "Camera 2", "Camera 3", "Camera 4"
	cameraSelectWidget = widget.NewSelect([]string{},onCameraSelectWidget)
	for i:=0;i<len(cameraList);i++ {
		cameraSelectWidget.Options = append(cameraSelectWidget.Options,fmt.Sprintf("Camera %d (%dx%d)",(i+1),cameraList[i].Width,cameraList[i].Height))
	}
	cameraSelectWidget.PlaceHolder = "Bitte wählen Sie eine Kamera aus"
	if intCamera>len(cameraList) {
		intCamera = 0
	}
	if len(cameraList)>0 {
		cameraSelectWidget.SetSelected(cameraSelectWidget.Options[intCamera])
	}


	thresholdHoughCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", thresholdHoughCircles))
	meanFindCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", meanFindCircles))
	dpHoughCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", dpHoughCircles))
	gaussianBlurFindCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%d", gaussianBlurFindCircles))
	adaptiveThresholdBlockSizeWidgetLabel = widget.NewLabel(fmt.Sprintf("%d", adaptiveThresholdBlockSize))
	adaptiveThresholdSubtractMeanWidgetLabel = widget.NewLabel(fmt.Sprintf("%.1f", adaptiveThresholdSubtractMean))



	thresholdHoughCirclesWidget = widget.NewSlider(0, 255)
	thresholdHoughCirclesWidget.Value = thresholdHoughCircles
	thresholdHoughCirclesWidget.OnChangeEnded = func(value float64) {
		thresholdHoughCircles = value
		thresholdHoughCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}


	meanFindCirclesWidget = widget.NewSlider(0, 255)
	meanFindCirclesWidget.Value = meanFindCircles
	meanFindCirclesWidget.OnChangeEnded = func(value float64) {
		meanFindCircles = value
		meanFindCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}

	dpHoughCirclesWidget = widget.NewSlider(0, 3)
	dpHoughCirclesWidget.Value = dpHoughCircles
	dpHoughCirclesWidget.OnChangeEnded = func(value float64) {
		dpHoughCircles = value
		dpHoughCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}

	gaussianBlurFindCirclesWidget = widget.NewSlider(0, 255)
	gaussianBlurFindCirclesWidget.Value = float64(gaussianBlurFindCircles)
	gaussianBlurFindCirclesWidget.OnChangeEnded = func(value float64) {

		gaussianBlurFindCircles = int(value)
		gaussianBlurFindCirclesWidgetLabel.SetText(fmt.Sprintf("%d", int(value)))

	}

	adaptiveThresholdBlockSizeWidget = widget.NewSlider(0, 255)
	adaptiveThresholdBlockSizeWidget.Value = float64(adaptiveThresholdBlockSize)
	adaptiveThresholdBlockSizeWidget.OnChangeEnded = func(value float64) {
		adaptiveThresholdBlockSize = int(value)
		adaptiveThresholdBlockSizeWidgetLabel.SetText(fmt.Sprintf("%d", int(value)))

	}

	adaptiveThresholdSubtractMeanWidget = widget.NewSlider(-10, 10)
	adaptiveThresholdSubtractMeanWidget.Value = float64(adaptiveThresholdSubtractMean)
	adaptiveThresholdSubtractMeanWidget.OnChangeEnded = func(value float64) {
		adaptiveThresholdSubtractMean = float32(value)
		adaptiveThresholdSubtractMeanWidgetLabel.SetText(fmt.Sprintf("%.1f", value))

	}

	

	paperImageCheck := widget.NewCheck("Anzeigen", func(c bool) {
		if c {
			paperImage.Show()
		} else {
			paperImage.Hide()
		}
	})
	paperImageCheck.SetChecked(showPaperImage)

	circleImageCheck := widget.NewCheck("Anzeigen", func(c bool) {
		if c {
			circleImage.Show()
		} else {
			circleImage.Hide()
		}
	})
	circleImageCheck.SetChecked(showCirlceImage)

	outputImageCheck := widget.NewCheck("Anzeigen", func(c bool) {
		if c {
			outputImage.Show()
		} else {
			outputImage.Hide()
		}
		
	})
	outputImageCheck.SetChecked(showOutputImage)

	debugListCheck := widget.NewCheck("Anzeigen", func(c bool) {
		if c {
			debugListWidget.Show()
		} else {
			debugListWidget.Hide()
		}
		
	})
	debugListCheck.SetChecked(showDebugList)

	//func(s string) { fmt.Println("selected", s) })
	// container.NewBorder(nil, nil, nil, nil,

	return container.New(layout.NewVBoxLayout(), 
		widget.NewLabel("Camera"),

		

		cameraSelectWidget,
		widget.NewAccordion(
			&widget.AccordionItem{
				Title:  "Kreisdetetion",
				Detail: container.New(
					layout.NewGridLayout(1), 
					widget.NewLabel("Mean Find Circles"),
					container.NewBorder( nil, nil, nil, meanFindCirclesWidgetLabel,meanFindCirclesWidget ),

					widget.NewLabel("Hough Circles Threshold"),
					container.NewBorder( nil, nil, nil, thresholdHoughCirclesWidgetLabel,thresholdHoughCirclesWidget ),

					widget.NewLabel("Inverse ratio of the accumulator"),
					container.NewBorder( nil, nil, nil, dpHoughCirclesWidgetLabel,dpHoughCirclesWidget ),

					widget.NewLabel("Blursize"),
					container.NewBorder( nil, nil, nil, gaussianBlurFindCirclesWidgetLabel,gaussianBlurFindCirclesWidget ),


					widget.NewLabel("Adaptive Threshold Block Size"),
					container.NewBorder( nil, nil, nil, adaptiveThresholdBlockSizeWidgetLabel,adaptiveThresholdBlockSizeWidget ),

					widget.NewLabel("Adaptive Threshold Subtract Mean"),
					container.NewBorder( nil, nil, nil, adaptiveThresholdSubtractMeanWidgetLabel,adaptiveThresholdSubtractMeanWidget ) ) }, 
			&widget.AccordionItem{
				Title:  "Ausgaben",
				Detail: container.New(
					layout.NewGridLayout(2), 
					widget.NewLabel("Kamerabild"),outputImageCheck,
					widget.NewLabel("Papier"),paperImageCheck,
					widget.NewLabel("Findcircle"),circleImageCheck,
					widget.NewLabel("Debug"),debugListCheck,) },
					
		),
	)
	/*

	var outputImage *canvas.Image
var paperImage *canvas.Image
var circleImage *canvas.Image

*/
		// widget.NewButton("Save", ))
}

func makeTopBar() fyne.CanvasObject {
	fullNameWidget = widget.NewLabel("Fullname")
	boxLabelWidget = widget.NewLabel("Kiste: UNBEKANNT")
	stackLabelWidget = widget.NewLabel("Stapel: UNBEKANNT")
	ballotLabelWidget = widget.NewLabel("Stimmzettel: UNBEKANNT")
	return container.New(layout.NewHBoxLayout(), 
		boxLabelWidget,
		stackLabelWidget,
		ballotLabelWidget,
		layout.NewSpacer(), fullNameWidget)
}

func makeOuterBorder() fyne.CanvasObject {
	// top := canvas.NewText("top bar", color.White)
	// left := canvas.NewText("left", color.White)
	bottom := widget.NewButton("Start/Stop", func() { 
		cameraChannelImage = make(chan gocv.Mat,1)

		imageChannelPaper = make(chan gocv.Mat,1)
		imageChannelCircle = make(chan gocv.Mat,1)

		currentBoxChannel = make(chan string,1)
		currentStackChannel = make(chan string,1)
		currentBarcodeChannel = make(chan string,1)
		
		
		if videoIsRunning {
			grabVideoCameraTicker.Stop()
			videoIsRunning = false
			runVideo = false
		} else {
			fmt.Println("id",os.Getpid())
			runVideo = true
			videoIsRunning = true
			//go grabcamera() 
			grabVideoCameraTicker = time.NewTicker(1 * time.Millisecond)
			
			// go grabDebugs()
			go grabcamera()  // Kamerabild abrufen
			go grabVideoImage() // kamera bild anzeigen

			go processImage() // Bild verarbeiten


			if false {
				go grabPaperImage()
				go grabChannelBarcodes()
				go grabCircleImage()
				go grabReadyToSaveImage()


				go scanBarcodeChannel()
				go processPaperChannelImage()
				go processTesseractChannelImage()
				go processRoisChannel()
				
				
				go grabCurrentBox()
				go grabCurrentStack()
			}
		}
		
	})
	// middle := canvas.NewText("content", color.White)
	return container.NewBorder(
		makeTopBar(), 
		bottom, nil, nil,  makeSplitTab())
}
func makeSplitTab() fyne.CanvasObject {
	/*
	left := widget.NewMultiLineEntry()
	left.Wrapping = fyne.TextWrapWord
	left.SetText("Long text is looooooooooooooong")
	*/
	
	/*
	outputImage = canvas.NewImageFromImage(matToImage(gocv.IMRead("Logo-large.png", gocv.IMReadColor)))
	outputImage.FillMode = canvas.ImageFillContain
	*/

	right := makeVideoGrid()
	/* container.NewVSplit(
		makeVideoGrid(),
		//nil,
	)
	*/
	
	// left := container.NewVScroll(canvas.NewText("Hello", color.White))
	left := container.NewVScroll(makeSettingsForm())
	
	// left.Width = 200
	c:=container.NewHSplit(left, right)
	c.Offset = 0.15
	return c
}


func makeMainPanel() fyne.CanvasObject {


	
	/*
	devices, err := avfoundation.Devices(avfoundation.Video)
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		fmt.Println(device.Name)
	}
	*/
	label := canvas.NewText("Anmelden", color.White)
	label.TextSize = 20
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle= fyne.TextStyle{Bold: true}



	ok := canvas.NewText("OK!", color.White)
	ok.TextSize = 20
	ok.Alignment = fyne.TextAlignCenter
	ok.TextStyle= fyne.TextStyle{Bold: true}

	nameText := canvas.NewText("Name!", color.White)
	nameText.TextSize = 20
	nameText.Alignment = fyne.TextAlignCenter
	nameText.TextStyle= fyne.TextStyle{Bold: true}
	
	image := canvas.NewImageFromResource(theme.FyneLogo())
	// image := canvas.NewImageFromURI(uri)
	// image := canvas.NewImageFromImage(src)
	// image := canvas.NewImageFromReader(reader, name)
	// image := canvas.NewImageFromFile(fileName)
	image.FillMode = canvas.ImageFillOriginal
	

	
	loginContainer = container.New(
		layout.NewVBoxLayout(), 
		layout.NewSpacer(),
		label,
		makeLoginFormTab(),
		layout.NewSpacer(),
	)

	mainAppContainer = makeOuterBorder()
	mainAppContainer.Hide()


	content := container.New(
		layout.NewStackLayout(), 
		loginContainer,
		mainAppContainer,
	)
	loginContainer.Hide()
					mainAppContainer.Show()
	return content
}

func makeLoginFormTab() fyne.CanvasObject {
	url := widget.NewEntry()
	url.SetPlaceHolder(strSystemUrl)
	url.SetText(strSystemUrl)

	login := widget.NewEntry()
	login.SetPlaceHolder(strSystemLogin)
	login.SetText(strSystemLogin)
	// email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	password.SetText(strSystemPassword)

	/*
	disabled := widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(string) {})
	disabled.Horizontal = true
	disabled.Disable()
	largeText := widget.NewMultiLineEntry()
	*/

	form := &widget.Form{
		SubmitText: "Anmelden",
		CancelText: "Abbrechen",
		Items: []*widget.FormItem{
			{Text: "URL", Widget: url, HintText: "Bitte gib die vollständige URL ein."},
			{Text: "Benutzername", Widget: login, HintText: "Bitte gib deinen Benutzernamen ein."},
			{Text: "Passwort", Widget: password, HintText: "Bitte gib dein Passwort ein."},
		},
		OnCancel: func() {
			// fmt.Println("Cancelled")
			os.Exit(0)
		},
		OnSubmit: func() {
			// fmt.Println("Form submitted")
			strUrl:= url.Text
			strLogin:= login.Text
			strPassword:= password.Text
			if strUrl == "" {
				strUrl = strSystemUrl
			}
			if strLogin == "" {
				strLogin = strSystemLogin
			}
			if strPassword == "" {
				strPassword = strSystemPassword
			}
			loginResponse, err := api.Login(strUrl, strLogin, strPassword)
			if err != nil {
				/*
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Login failed",
					Content: err.Error(),
				})
				*/
			} else {
				if loginResponse.Success {
					/*
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   "Login successful",
						Content: "Welcome " + loginResponse.Fullname,
					})
					*/

					pingResponse, _ = api.Ping(strUrl)
					fullNameWidget.SetText(loginResponse.Fullname)

					kandidatenResponse, _ = api.GetKandidaten(strUrl)
					fmt.Println(kandidatenResponse)
					loginContainer.Hide()
					mainAppContainer.Show()
				} else {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   "Login failed",
						Content: loginResponse.Msg,
					})
				}
			}
		},
	}
	// form.Append("Password", password)
	
	/*
	form.Resize(fyne.NewSize(1000, 300))
	content := container.New(layout.NewMaxLayout(), form)
	fmt.Println(content.Size())
	content.Resize(fyne.NewSize(800, 600))
	*/
	//form.NewSize(500, 300)
	// NewFormLayout(
	// form.Append("Disabled", disabled)
	//form.Append("Message", largeText)
	return form
}


func appwindow() {

	window = gocv.NewWindow("IMG")
		

	a := app.NewWithID("io.tualo.ballotscanner")
	
	
	a.SetIcon(data.FyneLogo)
	/*
	makeTray(a)
	logLifecycle(a)
	*/
	w := a.NewWindow("tualo ballot scanner")
	topWindow = w

	// w.SetMainMenu(makeMenu(a, w))
	
	w.SetMaster()
	w.SetContent(makeMainPanel())
	
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()

}