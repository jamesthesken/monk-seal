/*

main() which holds the MonkSeal gui

*/

package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"fmt"
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

	test := tview.NewTextView().
		SetTextColor(tcell.ColorYellow).
		SetScrollable(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})
	
// The form is getting harder to use for handling text output to the main message area

/*
	form := tview.NewForm().
		AddInputField("> ", "", 100, nil, func(text string) {
  			myMessage := text
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			k := event.Key()
			if k == tcell.KeyEnter {
				fmt.Fprintf(test, myMessage)
			}
			return event
		}).
		AddDropDown("Channels", []string{"James", "Bob"}, 0, nil).
		AddInputField("New Channel", "", 20, nil, nil).
		AddButton("Connect", func() {
			go rendezvousChat(textView, test)
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetHorizontal(true);
*/

	form := tview.NewInputField().
			SetLabel("> ").
			SetFieldWidth(100)
			
	// When user presses enter, the text clears and prints properly
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			k := event.Key()
			myMsg := form.GetText()
			if k == tcell.KeyEnter {
				fmt.Fprintf(test,"(You) > 	")
				fmt.Fprintf(test, myMsg)
				fmt.Fprintf(test,"\n")
				form.SetText("")
			}
			return event
		})

	
	//TODO: Add a new primitive as a pointer to the handle stream funcitons of chat.go

	// Grid container
	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false)

	//sideBar := newPrimitive("Messages")	
		
	// Layout for screens less than 100 cells.
	grid.AddItem(textView, 1, 1, 1, 1, 0, 100, false).
		AddItem(form, 2, 0, 1, 3, 0, 0, true)
	// Layout for screens wider than 100 cells.
	grid.AddItem(textView, 1, 2, 1, 1, 0, 100, false).
		AddItem(test, 1, 0, 1, 2, 0, 100, false).
		AddItem(form, 2, 0, 1, 3, 0, 100, true)
	

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
