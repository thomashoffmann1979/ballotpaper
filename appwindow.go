package main

import (
	"os"
	"fmt"

	"image"
	"image/color"

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


	// "github.com/pion/mediadevices"
	// "github.com/pion/mediadevices/pkg/prop"

	// This is required to register camera adapter
	_ "github.com/pion/mediadevices/pkg/driver/camera" 
	"github.com/pion/mediadevices/pkg/avfoundation"
)


const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

var loginContainer *fyne.Container
var mainAppContainer fyne.CanvasObject
var pingResponse api.PingResponse
var kandidatenResponse api.KandidatenResponse
var fullNameWidget *widget.Label
var outputImage *canvas.Image
var videoIsRunning bool = false

func matToImage(mat gocv.Mat) image.Image {
	img, _ := mat.ToImage()
	return img
}

func makeSplitTab() fyne.CanvasObject {
	/*
	left := widget.NewMultiLineEntry()
	left.Wrapping = fyne.TextWrapWord
	left.SetText("Long text is looooooooooooooong")
	*/
	fullNameWidget = widget.NewLabel("Fullname")
	outputImage = canvas.NewImageFromImage(matToImage(gocv.IMRead("Logo-large.png", gocv.IMReadColor)))
	outputImage.FillMode = canvas.ImageFillContain
	right := container.NewVSplit(
		outputImage,
		widget.NewButton("Button", func() { 
			fmt.Println("Button clicked",runVideo,videoIsRunning)
			if videoIsRunning {
				videoIsRunning = false
				runVideo = false
			} else {
				runVideo = true
				videoIsRunning = true
				go cameras() 
			}
		}),
	)
	
	return container.NewHSplit(container.NewVScroll(fullNameWidget), right)
}


func makeMainPanel() fyne.CanvasObject {


	

	devices, err := avfoundation.Devices(avfoundation.Video)
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		fmt.Println(device.Name)
	}

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

	mainAppContainer = makeSplitTab()
	mainAppContainer.Hide()


	content := container.New(
		layout.NewStackLayout(), 
		loginContainer,
		mainAppContainer,
	)
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
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Login failed",
					Content: err.Error(),
				})
			} else {
				if loginResponse.Success {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   "Login successful",
						Content: "Welcome " + loginResponse.Fullname,
					})

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