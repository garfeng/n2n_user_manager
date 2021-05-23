package main

import (
	"fmt"

	"github.com/robfig/cron"

	"github.com/garfeng/n2n_user_manager/client"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// TODO: UI will use walk

type N2NUI struct {
	username *walk.LineEdit
	password *walk.LineEdit
	ip       *walk.LineEdit
	mask     *walk.LineEdit

	mainWindow *MainWindow

	controller *client.Controller
	job        *cron.Cron
}

func NewUI() *N2NUI {

	ui := &N2NUI{
		username:   &walk.LineEdit{},
		password:   &walk.LineEdit{},
		ip:         &walk.LineEdit{},
		mask:       &walk.LineEdit{},
		mainWindow: nil,
		controller: client.NewController("config.toml"),
	}

	ui.CreateWindow()
	err := ui.controller.ReadConfig()
	if err != nil {
		// TODO: alert error
		panic(err)
	}

	return ui
}

func (ui *N2NUI) Run() {
	ui.mainWindow.Run()
}

func (ui *N2NUI) CreateWindow() {
	ui.mainWindow = &MainWindow{
		Title:  "N2N Client",
		Size:   Size{300, 400},
		Layout: VBox{},
		Children: []Widget{
			LineEdit{
				AssignTo: &ui.username,
			},
			LineEdit{AssignTo: &ui.password},
			LineEdit{AssignTo: &ui.ip},
			LineEdit{AssignTo: &ui.mask},
			PushButton{
				Text:      "Connect",
				OnClicked: ui.Connect,
			},
		},
	}
}

func (ui *N2NUI) Connect() {
	username := ui.username.Text()
	password := ui.password.Text()

	if username == "" || password == "" {
		fmt.Println("empty username or password")
		//TODO: alert error
		return
	}

	ipAndMask := []string{}
	if ui.ip.Text() != "" {
		ipAndMask = append(ipAndMask, ui.ip.Text())

		if ui.mask.Text() != "" {
			ipAndMask = append(ipAndMask, ui.mask.Text())
		}
	}

	ui.controller.InitUserInfo(username, password, ipAndMask...)

	err := ui.controller.Reconnect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer close(ui.controller.ErrChan)

	ui.job = cron.New()

	// Every 2:02:00 am
	ui.job.AddFunc("0 2 2 * * ?", func() {
		fmt.Println("reconnect")
		ui.controller.Reconnect()
	})
	ui.job.Start()
}

func (ui *N2NUI) Disconnect() {
	ui.job.Stop()
	ui.controller.Disconnect()
}
