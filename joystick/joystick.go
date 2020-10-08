package joystick

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/fyne-io/examples/img/icon"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	maxElements = 20
)

type joyUI struct {
	panel     fyne.CanvasObject
	labels    *widget.Box
	values    *widget.Box
	nameLabel *widget.Label
	buttons   [maxElements]*widget.Check
	axes      [maxElements]*widget.Label
	hats      [maxElements]*widget.Label
}

func newJoyUI() *joyUI {
	// We prebuild the interface and reuse it. Generating a new one each time leaks memory.
	var ui joyUI

	ui.labels = widget.NewVBox()
	ui.values = widget.NewVBox()
	ui.panel = widget.NewHBox(ui.labels, ui.values)

	ui.labels.Append(widget.NewLabelWithStyle("Name", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}))
	ui.nameLabel = widget.NewLabel("")
	ui.values.Append(ui.nameLabel)

	for i := 0; i < maxElements; i++ {
		ui.buttons[i] = widget.NewCheck("", nil)
		ui.labels.Append(widget.NewLabelWithStyle(fmt.Sprintf("Button %v", i), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}))
		ui.values.Append(ui.buttons[i])
	}
	for i := 0; i < maxElements; i++ {
		ui.axes[i] = widget.NewLabel("")
		ui.labels.Append(widget.NewLabelWithStyle(fmt.Sprintf("Axe %v", i), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}))
		ui.values.Append(ui.axes[i])
	}
	for i := 0; i < maxElements; i++ {
		ui.hats[i] = widget.NewLabel("")
		ui.labels.Append(widget.NewLabelWithStyle(fmt.Sprintf("Hat %v", i), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}))
		ui.values.Append(ui.hats[i])
	}

	return &ui
}

func (ui *joyUI) update(joy glfw.Joystick) {
	if !joy.Present() {
		ui.nameLabel.Text = fmt.Sprintf("Joystick %v not available", joy)
		for i := 1; i < len(ui.labels.Children); i++ {
			ui.labels.Children[i].Hide()
			ui.values.Children[i].Hide()
		}
		return
	}

	ui.nameLabel.Text = joy.GetName()
	buttons := joy.GetButtons()
	axes := joy.GetAxes()
	hats := joy.GetHats()

	w := 1
	for e := 0; e < maxElements; e++ {
		if e < len(buttons) {
			ui.buttons[e].SetChecked(buttons[e] == glfw.Press)
			ui.labels.Children[w].Show()
			ui.values.Children[w].Show()
		} else {
			ui.labels.Children[w].Hide()
			ui.values.Children[w].Hide()
		}
		w++
	}

	for e := 0; e < maxElements; e++ {
		if e < len(axes) {
			ui.axes[e].Text = fmt.Sprintf("%v", axes[e])
			ui.labels.Children[w].Show()
			ui.values.Children[w].Show()
		} else {
			ui.labels.Children[w].Hide()
			ui.values.Children[w].Hide()
		}
		w++
	}

	for e := 0; e < maxElements; e++ {
		if e < len(hats) {
			ui.hats[e].Text = fmt.Sprintf("%v", hats[e])
			ui.labels.Children[w].Show()
			ui.values.Children[w].Show()
		} else {
			ui.labels.Children[w].Hide()
			ui.values.Children[w].Hide()
		}
		w++
	}
}

// Show loads a joystick test window for the specified app context
func Show(app fyne.App) {
	window := app.NewWindow("Joystick")
	window.SetIcon(icon.BugBitmap)

	content := widget.NewHBox()
	window.SetContent(content)

	joyUI1 := newJoyUI()
	joyUI2 := newJoyUI()
	content = widget.NewHBox(
		joyUI1.panel,
		layout.NewSpacer(),
		canvas.NewLine(color.Black),
		layout.NewSpacer(),
		joyUI2.panel,
	)
	window.SetContent(content)

	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				joyUI1.update(glfw.Joystick1)
				joyUI2.update(glfw.Joystick2)
				content.Refresh()
			}
		}
	}()

	window.SetOnClosed(func() {
		done <- true
	})

	window.Show()
}
