/*

main() which holds the MonkSeal gui

*/

package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"fmt"
	"strings"
)


const logo = `
   __  ___          __     ____         __  
  /  |/  /__  ___  / /__  / __/__ ___ _/ /  
 / /|_/ / _ \/ _ \/  '_/ _\ \/ -_) _  / /   
/_/  /_/\___/_//_/_/\_\ /___/\__/\_,_/_/    
                                           
`                                                                            

// Instantiate the app
var app = tview.NewApplication()

func main() {
	
	// What's the size of the logo?
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)
	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}

	logoBox := tview.NewTextView().
		SetTextColor(tcell.ColorYellow)
		
	fmt.Fprint(logoBox, logo)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 7, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logoBox, logoWidth, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true)

	// Formatting for the rendezvous printouts
	textView := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetScrollable(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})

	// Message box
	msgField := tview.NewTextView().
		SetTextColor(tcell.ColorYellow).
		SetScrollable(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})

	myMessage := tview.NewInputField().
			SetLabel("> ").
			SetFieldWidth(100)
			
/*
	// When user presses enter, the text clears and prints properly
	myMessage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			k := event.Key()
			myMsg := myMessage.GetText()
			if k == tcell.KeyEnter {
				fmt.Fprintf(msgField,"(You) > 	")
				fmt.Fprintf(msgField, myMsg)
				fmt.Fprintf(msgField,"\n")
				myMessage.SetText("")
			}
			return event
		})

*/

	// Settings area for selecting which chat room the user is in
	mySettings := tview.NewForm().
		AddDropDown("Channels", []string{"James", "Bob"}, 0, nil).
		AddInputField("New Channel", "", 20, nil, nil).
		AddButton("Connect", func() {
			go rendezvousChat(textView, msgField, myMessage) // Add argument based on selected dropdown
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetHorizontal(true);

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			j := event.Key()

			if j == tcell.KeyRight {
				app.SetFocus(mySettings)
			}else if j == tcell.KeyLeft{
				app.SetFocus(myMessage)
			}
			return event
		})

	//TODO: Add a new primitive as a pointer to the handle stream funcitons of chat.go

	// Grid container
	grid := tview.NewGrid().
		SetRows(4, 0, 4).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(flex, 0, 0, 1, 3, 0, 0, false)
		
	// Layout for screens less than 100 cells.
	grid.AddItem(textView, 1, 1, 1, 1, 0, 100, false).
		AddItem(myMessage, 2, 0, 1, 3, 0, 0, false).
		AddItem(mySettings, 3, 0, 1, 3, 0, 0, false).
		AddItem(msgField, 1, 0, 1, 3, 0, 0, false)
		

	// Layout for screens wider than 100 cells.
	grid.AddItem(textView, 1, 2, 1, 1, 0, 100, false).
		AddItem(msgField, 1, 0, 1, 2, 0, 100, false).
		AddItem(myMessage, 2, 0, 1, 3, 0, 100, false).
		AddItem(mySettings, 3, 0, 1, 3, 0, 100, false)
	

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
