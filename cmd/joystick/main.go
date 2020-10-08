// Package main launches the joystick example directly
package main

import (
	"fyne.io/fyne/app"
	"github.com/fyne-io/examples/img/icon"
	"github.com/fyne-io/examples/joystick"
)

func main() {
	app := app.New()
	app.SetIcon(icon.BugBitmap)

	joystick.Show(app)
	app.Run()
}
