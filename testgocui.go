// Saratoga Interactive Client
// Test aplication and example of cli interface with a command and message split pane window.

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charlesetsmith/testgocui/cli"
	"github.com/charlesetsmith/testgocui/screen"
	"github.com/jroimartin/gocui"
)

// Cinfo - Information held on the cmd view
var Cinfo screen.Viewinfo

// Minfo - Information held on the msg view
var Minfo screen.Viewinfo

// *******************************************************************

// Return the length of the prompt
func promptlen(v screen.Viewinfo) int {
	return len(v.Prompt) + len(strconv.Itoa(v.Curline)) + v.Ppad
}

// Change view between "cmd" and "msg" window
func switchView(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		log.Fatal("switchView - g or v is nil")
	}
	var err error

	if v.Name() == "cmd" {
		v, err = g.SetCurrentView("msg")
	} else {
		v, err = g.SetCurrentView("cmd")
	}
	return err
}

// Backspace or Delete
func backSpace(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		log.Fatal("backSpace - g or v is nil")
	}
	cx, _ := v.Cursor()
	if cx <= promptlen(Cinfo) { // Dont move we are at the prompt
		return nil
	}
	// Delete rune backwards
	v.EditDelete(true)
	return nil
}

// Handle Left Arrow Move -- All good
func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		log.Fatal("cursorLeft - g or v is nil")
	}
	cx, cy := v.Cursor()
	if cx <= promptlen(Cinfo) { // Dont move
		return nil
	}
	// Move back a character
	if err := v.SetCursor(cx-1, cy); err != nil {
		screen.Fprintln(g, "msg", "bwhite_black", "LeftArrow:", "cx=", cx, "cy=", cy, "error=", err)
	}
	return nil
}

// Handle Right Arrow Move - All good
func cursorRight(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		log.Fatal("cursorRight - g or v is nil")
	}
	cx, cy := v.Cursor()
	line, _ := v.Line(cy)
	if cx >= len(line)-1 { // We are at the end of line do nothing
		v.SetCursor(len(line), cy)
		return nil
	}
	// Move forward a character
	if err := v.SetCursor(cx+1, cy); err != nil {
		screen.Fprintln(g, "msg", "bwhite_red", "RightArrow:", "cx=", cx, "cy=", cy, "error=", err)
	}
	return nil
}

// Handle down cursor -- All good!
// well not quite, still issue if we scroll down before hitting return
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		log.Fatal("cursorDown - g or v is nil")
	}
	_, oy := v.Origin()
	cx, cy := v.Cursor()

	// Don't move down if we are at the last line in current views Bufferlines
	if oy+cy >= len(v.BufferLines())-1 {
		return nil
	}
	err := v.SetCursor(cx, cy+1)
	if err != nil { // Reset the origin
		if err := v.SetOrigin(0, oy+1); err != nil { // changed ox to 0
			screen.Fprintf(g, "msg", "bwhite_red", "SetOrigin error=%s", err)
			return err
		}

	}
	// Move the cursor to the end of the current line
	_, cy = v.Cursor()
	if line, err := v.Line(cy); err == nil {
		v.SetCursor(len(line), cy)
	}
	return nil
}

// Handle up cursor -- All good!
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		log.Fatal("cursorUp - g or v is nil")
	}
	_, oy := v.Origin()
	cx, cy := v.Cursor()
	err := v.SetCursor(cx, cy-1)
	if err != nil && oy > 0 { // Reset the origin
		if err := v.SetOrigin(0, oy-1); err != nil { // changed ox to 0
			screen.Fprintf(g, "msg", "bwhite_red", "SetOrigin error=%s", err)
			return err
		}
	}
	// Move the cursor to the end of the current line
	_, cy = v.Cursor()
	if line, err := v.Line(cy); err == nil {
		v.SetCursor(len(line), cy)
	}
	return nil
}

// This is where we process command line inputs after a CR entered
func getLine(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		log.Fatal("getLine - g or v is nil")
	}
	// Find out where we are
	_, cy := v.Cursor()
	// Get the line
	line, _ := v.Line(cy)
	// screen.Fprintf(g, "msg", "red_black", "cx=%d cy=%d lines=%d line=%s\n",
	//	len(v.BufferLines()), cx, cy, line)
	command := strings.SplitN(line, ":", 2)
	if command[1] == "" { // We have just hit enter - do nothing
		return nil
	}
	Cinfo.Commands = append(Cinfo.Commands, command[1])

	go func(*gocui.Gui, string) {
		// defer Sarwg.Done()
		cli.Docmd(g, command[1])
	}(g, command[1])

	if command[1] == "exit" || command[1] == "quit" {
		// Sarwg.Wait()
		err := quit(g, v)
		// THIS IS A KLUDGE FIX IT WITH A CHANNEL
		log.Fatal("\nGocui Exit. Bye!", err)
	}

	// Bump the number of the current line
	Cinfo.Curline++
	// Our new x position will always be after the prompt + 3 for []: chars
	xpos := promptlen(Cinfo)
	// Have we scrolled past the length of v, if so reset the origin

	if err := v.SetCursor(xpos, cy+1); err != nil {
		// screen.Fprintln(g, "msg", "red_black", "We Scrolled past length of v", err)
		_, oy := v.Origin()
		// screen.Fprintf(g, "msg", "red_black", "Origin reset ox=%d oy=%d\n", ox, oy)
		if err := v.SetOrigin(0, oy+1); err != nil { // changed xpos to 0
			// screen.Fprintln(g, "msg", "red_black", "SetOrigin Error:", err)
			return err
		}
		// Set the cursor to last line in v
		if verr := v.SetCursor(xpos, cy); verr != nil {
			screen.Fprintln(g, "msg", "bwite_red", "Setcursor out of bounds:", verr)
		}
		// cx, cy := v.Cursor()
		// screen.Fprintf(g, "msg", "red_black", "cx=%d cy=%d line=%s\n", cx, cy, line)
	}
	// Put up the new prompt on the next line
	screen.Fprintf(g, "cmd", "yellow_black", "\n%s[%d]:", Cinfo.Prompt, Cinfo.Curline)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("cmd", gocui.KeyCtrlSpace, gocui.ModNone, switchView); err != nil {
		return err
	}

	if err := g.SetKeybinding("cmd", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyArrowLeft, gocui.ModNone, cursorLeft); err != nil {
		return nil
	}
	if err := g.SetKeybinding("cmd", gocui.KeyArrowRight, gocui.ModNone, cursorRight); err != nil {
		return nil
	}
	if err := g.SetKeybinding("cmd", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyBackspace, gocui.ModNone, backSpace); err != nil {
		return nil
	}
	if err := g.SetKeybinding("cmd", gocui.KeyBackspace2, gocui.ModNone, backSpace); err != nil {
		return nil
	}
	if err := g.SetKeybinding("cmd", gocui.KeyDelete, gocui.ModNone, backSpace); err != nil {
		return nil
	}
	return nil
}

// FirstPass -- First time around layout we don;t put \n at end of prompt
var FirstPass = true

// For working out screen positions in cli i/o

// CmdLines - Number of lines in Cmd View
var CmdLines int

// MaxX - Maximum screen X Value
var MaxX int

// MaxY - Maximum screen Y Value
var MaxY int

func layout(g *gocui.Gui) error {

	var err error
	var cmd *gocui.View
	var msg *gocui.View

	ratio := 4            // Ratio of cmd to err views
	MaxX, MaxY = g.Size() // Set the MaxX and MaxY to current size
	// This is the command line input view -- cli inputs and return messages go here
	if cmd, err = g.SetView("cmd", 0, MaxY-(MaxY/ratio)+1, MaxX-1, MaxY-1); err != nil {
		CmdLines = (MaxY / ratio) - 3 // Number of input lines in cmd view
		if err != gocui.ErrUnknownView {
			return err
		}
		cmd.Title = "Command Line"
		cmd.Highlight = false
		cmd.BgColor = gocui.ColorBlack
		cmd.FgColor = gocui.ColorGreen
		cmd.Editable = true
		cmd.Overwrite = true
		cmd.Wrap = true
	}
	// This is the message view window - All sorts of status & error messages go here
	if msg, err = g.SetView("msg", 0, 0, MaxX-1, MaxY-MaxY/ratio); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		msg.Title = "Messages"
		msg.Highlight = false
		msg.BgColor = gocui.ColorBlack
		msg.FgColor = gocui.ColorYellow
		msg.Editable = false
		msg.Wrap = true
		msg.Overwrite = false
		msg.Autoscroll = true
	}

	// All inputs happen via the cmd view
	if cmd, err = g.SetCurrentView("cmd"); err != nil {
		return err
	}

	// Display the prompt without the \n first time around
	if FirstPass {
		cmd.SetCursor(0, 0)
		xpos := promptlen(Cinfo)
		Cinfo.Curline = 0
		screen.Fprintf(g, "cmd", "yellow_black", "%s[%d]:", Cinfo.Prompt, Cinfo.Curline)
		cmd.SetCursor(xpos, 0)
		FirstPass = false
	}
	return nil
}

// Main
func main() {

	// cli.Cinfo = &Cinfo

	Cinfo.Prompt = "testgocui"
	Cinfo.Ppad = 3 // len("[]:") // For []: in chars in the prompt e.g. "test[5]:"

	// Set up the gocui interface and start the mainloop
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Printf("Cannot run gocui user interface")
		log.Fatal(err)
	}
	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	// The Base calling functions for testgocui live in cli.go so look there first!
	errflag := make(chan error, 1)
	go mainloop(g, errflag)

	select {
	case err := <-errflag:
		fmt.Println("Mainloop has quit:", err.Error())
	}
	return
}

// Go routine for command line loop
func mainloop(g *gocui.Gui, done chan error) {
	var err error

	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Printf("%s", err.Error())
	}
	done <- err
}
