/*

main() which holds the MonkSeal gui

*/

package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"strconv"
	"fmt"
)

// Slide is a function which returns the slide's main primitive and its title.
// It receives a "nextSlide" function which can be called to advance the
// presentation to the next slide.
type Slide func(nextSlide func()) (title string, content tview.Primitive)

// Instantiate the app
var app = tview.NewApplication()
// Widget switching
var pages = tview.NewPages() 

func main() {
	
	// The presentation slides.
	slides := []Slide{
		inputMessage,
		selectChannel,
	}

	// The bottom row has some info on where we are.
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	// Create the pages for all slides.
	currentSlide := 0
	info.Highlight(strconv.Itoa(currentSlide))
	pages := tview.NewPages()
	
	nextSlide := func() {
		currentSlide = (currentSlide + 1) % len(slides)
		info.Highlight(strconv.Itoa(currentSlide)).
			ScrollToHighlight()
		pages.SwitchToPage(strconv.Itoa(currentSlide))
	}
	for index, slide := range slides {
		title, primitive := slide(nextSlide)
		pages.AddPage(strconv.Itoa(index), primitive, true, index == currentSlide)
		fmt.Fprintf(info, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, title)
	}

	// Shortcuts to navigate 
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			nextSlide()
		} 
		return event
	})


	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
