// Test aplication and example of cli interface with a command and message split pane window.

package cli

import (
	"fmt"
	"strings"

	"github.com/charlesetsmith/testgocui/screen"
	"github.com/jroimartin/gocui"
)

// Viewinfo - Command history
type Viewinfo struct {
	Commands []string
	Prompt   string
	Ppad     int // Number of pad characters around prompt e.g. prompt[99]: would be 3
	Curline  int
	Numlines int
}

// Cinfo -- Used in the "ls" command - WHICH CURRENTLY DUMPS CORE!!!! FIX IT!!!!!
// var Cinfo *Viewinfo

// All of the different command line input handlers

func cmda(g *gocui.Gui, args []string, cmds Viewinfo) {
	screen.Fprintln(g, "msg", "green_black", "Command A", args)
}

func cmdb(g *gocui.Gui, args []string, cmds Viewinfo) {
	screen.Fprintln(g, "msg", "green_black", "Command B", args)
}

func cmdc(g *gocui.Gui, args []string, cmds Viewinfo) {
	screen.Fprintln(g, "msg", "green_black", "Command B", args)
}

func ls(g *gocui.Gui, args []string, cmds Viewinfo) {

	s := fmt.Sprintf("Number Commands=%d\n", len(cmds.Commands))
	for c := range cmds.Commands {
		s += fmt.Sprintf("Commands=%s\n", cmds.Commands[c])
	}
	screen.Fprintln(g, "msg", "cyan_black", s)
}

// Quit saratoga
func exit(g *gocui.Gui, args []string, cmds Viewinfo) {
	if len(args) == 1 { // exit 0
		screen.Fprintln(g, "msg", "green_black", "Gocui Good Bye!")
		return
	}
}

/* ************************************************************************** */

type cmdfunc func(*gocui.Gui, []string, Viewinfo)

// Commands and function pointers to handle them
var cmdhandler = map[string]cmdfunc{
	"ca":   cmda,
	"cb":   cmdb,
	"cc":   cmdc,
	"ls":   ls, // THIS CURRENTLY DUMPS CORE!!!!! FIX IT!!!!
	"quit": exit,
}

// Docmd -- Execute the command entered
func Docmd(g *gocui.Gui, s string, cmds Viewinfo) {
	if s == "" { // Handle just return
		return
	}

	// Get rid of leading and trailing whitespace
	s = strings.TrimSpace(s)
	vals := strings.Fields(s)
	// Lookup the command and execute it
	for c := range cmdhandler {
		if c == vals[0] {
			fn, ok := cmdhandler[c]
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
