/*

main() which holds the MonkSeal gui

*/

package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

// Instantiate the app
var app = tview.NewApplication()

func main() {
	
	// Aligns the text to the center of the grid section
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	
	// Formatting for the rendezvous printouts
	textView := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetScrollable(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})

	
	form := tview.NewForm().
		AddInputField("> ", "", 100, nil, nil).
		AddDropDown("Channels", []string{"James", "Bob"}, 0, nil).
		AddInputField("New Channel", "", 20, nil, nil).
		AddButton("Connect", func() {
			go rendezvousChat(textView)
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetHorizontal(true);
	
	// Grid container
	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false)

	sideBar := newPrimitive("Messages")	
		
	// Layout for screens less than 100 cells.
	grid.AddItem(textView, 1, 1, 1, 1, 0, 100, false).
		AddItem(form, 2, 0, 1, 3, 0, 0, true)

	// Layout for screens wider than 100 cells.
	grid.AddItem(textView, 1, 2, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 0, 1, 2, 0, 100, false).
		AddItem(form, 2, 0, 1, 3, 0, 100, true)
	

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
