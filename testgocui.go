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
var Cinfo cli.Cmdhist

// *******************************************************************

// Place a "packet" on top of the others if it is to be shown.
func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	_, err := g.SetViewOnTop(name)
	if showpacket {
		g.SetViewOnTop("packet")
	}
	return nil, err
}

// Rotate through the views - CtrlSpace
// only show packet view if it is enabled
func switchView(g *gocui.Gui, v *gocui.View) error {
	var err error
	var view string

	switch v.Name() {
	case "cmd":
		view = "msg"
	case "msg":
		if showpacket {
			view = "packet"
		} else {
			view = "err"
		}
	case "err":
		view = "cmd"
	case "packet":
		view = "err"
	}
	if _, err = setCurrentViewOnTop(g, view); err != nil {
		return err
	}
	return nil
}

// Backspace or Delete -- All good
func backSpace(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case "cmd":
		cx, _ := v.Cursor()
		if cx <= promptlen(Cinfo) { // Dont move we are at the prompt
			return nil
		}
		// Delete rune backwards
		v.EditDelete(true)
	case "msg", "packet", "err":
		return nil
	}
	return nil
}

// Handle Left Arrow Move -- All good
func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case "cmd":
		cx, cy := v.Cursor()
		if cx <= promptlen(Cinfo) { // Dont move we are at the prompt
			return nil
		}
		// Move back a character
		if err := v.SetCursor(cx-1, cy); err != nil {
			screen.ErrPrintln(g, "white_black", v.Name(), "LeftArrow:", "cx=", cx, "cy=", cy, "error=", err)
		}
	case "msg", "packet":
		return nil
	}
	return nil
}

// Handle Right Arrow Move - All good
func cursorRight(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case "cmd":
		cx, cy := v.Cursor()
		line, _ := v.Line(cy)
		if cx >= len(line)-1 { // We are at the end of line do nothing
			v.SetCursor(len(line), cy)
			return nil
		}
		// Move forward a character
		if err := v.SetCursor(cx+1, cy); err != nil {
			screen.ErrPrintln(g, "white_red", "RightArrow:", "cx=", cx, "cy=", cy, "error=", err)
		}
	case "msg", "packet":
		return nil
	}
	return nil
}

// Handle down cursor -- All good!
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	// Don't move down if we are at the last line in current views Bufferlines
	if oy+cy == len(v.BufferLines())-1 {
		screen.ErrPrintf(g, "white_black", "%s Down oy=%d cy=%d lines=%d\n",
			v.Name(), oy, cy, len(v.BufferLines()))
		return nil
	}
	if err := v.SetCursor(cx, cy+1); err != nil {
		screen.ErrPrintf(g, "magenta_black", "%s Down oy=%d cy=%d lines=%d err=%s\n",
			v.Name(), oy, cy, len(v.BufferLines()), err.Error())
		// ox, oy = v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			screen.ErrPrintf(g, "cyan_black", "%s Down oy=%d cy=%d lines=%d err=%s\n",
				v.Name(), oy, cy, len(v.BufferLines()), err.Error())
			return err
		}
	}
	screen.ErrPrintf(g, "green_black", "%s Down oy=%d cy=%d lines=%d\n",
		v.Name(), oy, cy, len(v.BufferLines()))
	return nil
}

// Why don't we scroll up!!!
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	screen.ErrPrintf(g, "green_black", "%s Up ox=%d oy=%d cx=%d cy=%d lines=%d\n",
		v.Name(), ox, oy, cx, cy, len(v.BufferLines()))
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		screen.ErrPrintf(g, "magenta_black", "%s SetCur Up oy=%d cy=%d lines=%d err=%s\n",
			v.Name(), oy, cy, len(v.BufferLines()), err.Error())
		if err := v.SetOrigin(ox, oy-1); err != nil {
			screen.ErrPrintf(g, "cyan_black", "%s SetOri Up oy=%d cy=%d lines=%d err=%s\n",
				v.Name(), oy-1, cy, len(v.BufferLines()), err.Error())
			return err
		} else {
			screen.ErrPrintf(g, "red_black", "%s SetOri Up oy=%d cy=%d lines=%d\n",
				v.Name(), oy, cy, len(v.BufferLines()))
			return nil
		}
	}
	_, cy = v.Cursor()
	screen.ErrPrintf(g, "green_black", "%s Up oy=%d cy=%d lines=%d\n",
		v.Name(), oy, cy, len(v.BufferLines()))
	return nil
}

// This is where we process command line inputs after a CR entered
func getLine(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		log.Fatal("getLine - g or v is nil")
	}
	switch v.Name() {
	case "cmd":
		// c := &Cinfo
		// Find out where we are
		_, cy := v.Cursor()
		// Get the line
		line, _ := v.Line(cy)

		command := strings.SplitN(line, ":", 2)
		if command[1] == "" { // We have just hit enter - do nothing
			return nil
		}
		// Save the command into history
		Cinfo.Commands = append(Cinfo.Commands, command[1])

		// Spawn a go to run the command
		go func(*gocui.Gui, string) {
			// defer Sarwg.Done()
			cli.Docmd(g, command[1], Cinfo)
		}(g, command[1])

		if command[1] == "exit" || command[1] == "quit" {
			// Sarwg.Wait()
			err := quit(g, v)
			// THIS IS A KLUDGE FIX IT WITH A CHANNEL
			log.Fatal("\nGocui Exit. Bye!\n", err)
		}
		prompt(g, v)
	case "msg", "packet", "err":
		return cursorDown(g, v)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// ShowPacket - Show Packet trace info
var showpacket bool = false

// Turn on/off the Packet View
func showPacket(g *gocui.Gui, v *gocui.View) error {
	var err error

	if g == nil || v == nil {
		log.Fatal("showPacket g is nil")
	}
	showpacket = !showpacket
	if showpacket {
		_, err = g.SetViewOnTop("packet")
	} else {
		_, err = g.SetViewOnTop("msg")
	}
	return err
}

// Bind keys to function handlers
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlSpace, gocui.ModNone, switchView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, cursorLeft); err != nil {
		return nil
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, cursorRight); err != nil {
		return nil
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, showPacket); err != nil {
		return nil
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyBackspace, gocui.ModNone, backSpace); err != nil {
		return nil
	}
	if err := g.SetKeybinding("", gocui.KeyBackspace2, gocui.ModNone, backSpace); err != nil {
		return nil
	}
	if err := g.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, backSpace); err != nil {
		return nil
	}
	return nil
}

// Return the length of the prompt
func promptlen(v cli.Cmdhist) int {
	return len(v.Prompt) + len(strconv.Itoa(v.Curline)) + v.Ppad
}

// Display the prompt
func prompt(g *gocui.Gui, v *gocui.View) {
	if g == nil || v == nil || v.Name() != "cmd" {
		log.Fatal("prompt must be in cmd view")
	}
	_, oy := v.Origin()
	_, cy := v.Cursor()
	// Only display it if it is on the next new line
	if oy+cy == Cinfo.Curline {
		if FirstPass { // Just the prompt no precedin \n as we are the first line
			screen.CmdPrintf(g, "yellow_black", "%s[%d]:", Cinfo.Prompt, Cinfo.Curline)
			v.SetCursor(promptlen(Cinfo), cy)
		} else { // End the last command by going to new lin \n then put up the new prompt
			Cinfo.Curline++
			screen.CmdPrintf(g, "yellow_black", "\n%s[%d]:", Cinfo.Prompt, Cinfo.Curline)
			_, cy := v.Cursor()
			v.SetCursor(promptlen(Cinfo), cy)
			if err := cursorDown(g, v); err != nil {
				screen.MsgPrintln(g, "red_black", "Cannot move to next line")
			}
			_, cy = v.Cursor()
			v.SetCursor(promptlen(Cinfo), cy+1)
		}
	}
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
	var packet *gocui.View

	ratio := 4 // Ratio of cmd to msg views

	// Maximum size of x and y
	maxx, maxy := g.Size()
	// This is the command line input view -- cli inputs and return messages go here
	if cmd, err = g.SetView("cmd", 0, maxy-(maxy/ratio)+1, maxx/2-1, maxy-1); err != nil {
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
		cmd.Autoscroll = false // This (false) enables vertical scrolling!
	}
	// This is the error msg view -- mic errors go here
	if cmd, err = g.SetView("err", maxx/2, maxy-(maxy/ratio)+1, maxx-1, maxy-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		cmd.Title = "Errors"
		cmd.Highlight = false
		cmd.BgColor = gocui.ColorBlack
		cmd.FgColor = gocui.ColorGreen
		cmd.Editable = false
		cmd.Overwrite = false
		cmd.Wrap = true
		cmd.Autoscroll = false // This (false) enables vertical scrolling!
	}
	// This is the packet trace window - packet trace history goes here
	// Toggles on/off with CtrlP
	if packet, err = g.SetView("packet", maxx-maxx/4, 1, maxx-2, maxy-maxy/ratio-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		packet.Title = "Packets"
		packet.Highlight = false
		packet.BgColor = gocui.ColorBlack
		packet.FgColor = gocui.ColorMagenta
		packet.Editable = false
		packet.Wrap = true
		packet.Overwrite = false
		packet.Autoscroll = false // This (false) enables vertical scrolling!
	}

	// This is the message view window - Status & error messages go here
	if msg, err = g.SetView("msg", 0, 0, maxx-1, maxy-maxy/ratio); err != nil {
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
		msg.Autoscroll = false // This (false) enables vertical scrolling!
	}

	// Display the prompt without the \n first time around
	if FirstPass {
		// All inputs happen via the cmd view and go there to start
		if cmd, err = g.SetCurrentView("cmd"); err != nil {
			return err
		}
		g.Cursor = true
		g.Highlight = true
		g.SelFgColor = gocui.ColorRed
		g.SelBgColor = gocui.ColorWhite
		cmd.SetCursor(0, 0)
		Cinfo.Curline = 0
		prompt(g, cmd)
		FirstPass = false
		screen.MsgPrintln(g, "white_black", "CtrlSpace - Rotate between views")
		screen.MsgPrintln(g, "white_black", "CtrlP - Show/Hide Packet view")
		screen.MsgPrintln(g, "white_black", "? for help")
	}
	return nil
}

// Main
func main() {

	// The prompt for the command view
	Cinfo.Prompt = "testgocui"
	Cinfo.Ppad = 3 // len("[]:") // For []: in chars in the prompt e.g. "gocui[5]:"

	// Set up the gocui interface and start the mainloop
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Printf("Cannot run gocui user interface")
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	// The Base calling functions for testgocui live in cli.go so look there first!
	errflag := make(chan error, 1)
	go mainloop(g, errflag)

	err = <-errflag
	fmt.Println("Mainloop has quit with error", err.Error())
	/*
		select {
		case err := <-errflag:
			fmt.Println("Mainloop has quit:", err.Error())
		}
	*/
}

// Go routine for command line loop
func mainloop(g *gocui.Gui, done chan error) {
	var err error

	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Printf("%s", err.Error())
	}
	done <- err
}
