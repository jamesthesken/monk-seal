package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)


// HelloWorld shows a simple "Hello world" example.
func selectChannel(nextSlide func()) (title string, content tview.Primitive) {

	// Aligns the text to the center of the grid section
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	
	// Formatting for the rendezvous printouts
	textView := tview.NewTextView().
		SetTextColor(tcell.ColorYellow).
		SetScrollable(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})

	// List
	list := tview.NewList().
		AddItem("Chat", "Some explanatory text", 'a', func() {
			go rendezvousChat(textView)
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	
	// Input text
	inputField := tview.NewInputField().
		SetLabel("Enter a number: ").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)
	
	// Grid container
	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false)

	sideBar := newPrimitive("Side Bar")		

	// Layout for screens wider than 100 cells.
	grid.AddItem(list, 1, 0, 1, 1, 0, 100, true).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false).
		AddItem(textView, 1, 1, 1, 1, 0, 100, false).
		AddItem(inputField, 2, 0, 1, 3, 0, 100, false)

	return "selectChannel", grid
}
