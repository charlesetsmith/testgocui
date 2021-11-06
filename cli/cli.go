// Test aplication and example of cli interface with a command and message split pane window.

package cli

import (
	"fmt"
	"strings"

	"github.com/charlesetsmith/testgocui/screen"
	"github.com/jroimartin/gocui"
)

// Cmdinfo - Previos Command history, prompt info
type Cmdinfo struct {
	Commands []string
	Prompt   string
	Ppad     int // Number of pad characters around prompt e.g. prompt[99]: would be 3
	Curline  int // What is the current command line # we are on
	// Numlines int
}

type cmdfunc func(*gocui.Gui, []string, Cmdinfo)

// Cmdhandler - Commands and function pointers to handle them
var Cmdhandler = map[string]cmdfunc{
	"ca":    cmda,
	"cb":    cmdb,
	"cc":    cmdc,
	"ls":    ls,
	"quit":  exit,
	"exit":  exit,
	"help":  usage,
	"usage": usage,
}

// The different command line input handlers

// cmda [args]...
func cmda(g *gocui.Gui, args []string, cmds Cmdinfo) {
	screen.Fprintln(g, "msg", "green_black", "Command A", args)
}

// cmdb [args]...
func cmdb(g *gocui.Gui, args []string, cmds Cmdinfo) {
	screen.Fprintln(g, "msg", "green_black", "Command B", args)
}

// cmdc [args]...
func cmdc(g *gocui.Gui, args []string, cmds Cmdinfo) {
	screen.Fprintln(g, "msg", "green_black", "Command B", args)
}

// ls - list the history of commands to the msg window
func ls(g *gocui.Gui, args []string, cmds Cmdinfo) {

	var s string
	for c := range cmds.Commands {
		s += fmt.Sprintf("%d=%s\n", c, cmds.Commands[c])
	}
	screen.Fprintln(g, "msg", "cyan_black", s)
}

// Quit saratoga
func exit(g *gocui.Gui, args []string, cmds Cmdinfo) {
	if len(args) == 1 { // exit 0
		screen.Fprintln(g, "msg", "green_black", "Gocui Good Bye!")
		return
	}
}

// usage - list usage of available commands
func usage(g *gocui.Gui, args []string, cmds Cmdinfo) {
	// Usage and description map of each command
	var help = map[string]string{
		"ca [arg]...": "Command Example a",
		"cb [arg]...": "Command Example b",
		"cc [arg]...": "Command Example c",
		"ls":          "History of commands entered",
		"quit":        "Bye!",
		"exit":        "Bye!",
		"help":        "List of available commands",
		"usage":       "List of available commands",
	}
	var s string
	for idx := range help {
		s += fmt.Sprintf("%s - %s\n", idx, help[idx])
	}
	screen.Fprintln(g, "msg", "cyan_black", s)
}

/* ************************************************************************** */

// Docmd -- Execute the command entered
func Docmd(g *gocui.Gui, s string, cmds Cmdinfo) {
	if s == "" { // Handle just return
		return
	}

	// Get rid of leading and trailing whitespace
	s = strings.TrimSpace(s)
	vals := strings.Fields(s)
	// Lookup the command and execute it
	for c := range Cmdhandler {
		if c == vals[0] {
			fn, ok := Cmdhandler[c]
			if ok {
				fn(g, vals, cmds)
				return
			}
			screen.Fprintln(g, "msg", "bwhite_red", "Cannot execute:", vals[0])
			return
		}
	}
	screen.Fprintln(g, "msg", "bwhite_red", "Invalid command:", vals[0])
	return
}
