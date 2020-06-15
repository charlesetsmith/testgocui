// Test aplication and example of cli interface with a command and message split pane window.

package cli

import (
	"fmt"
	"strings"

	"github.com/charlesetsmith/testgocui/screen"
	"github.com/jroimartin/gocui"
)

// Cinfo -- Used in the "ls" command
var Cinfo *screen.Viewinfo

// All of the different command line input handlers

func cmda(g *gocui.Gui, args []string) {
	screen.Fprintln(g, "msg", "green_black", "Command A", args)
}

func cmdb(g *gocui.Gui, args []string) {
	screen.Fprintln(g, "msg", "green_black", "Command B", args)
}

func cmdc(g *gocui.Gui, args []string) {
	screen.Fprintln(g, "msg", "green_black", "Command B", args)
}

func ls(g *gocui.Gui, args []string) {

	s := fmt.Sprintf("Number Commands=%d\n", len(Cinfo.Commands))
	for c := range Cinfo.Commands {
		s += fmt.Sprintf("Commands=%s\n", Cinfo.Commands[c])
	}
	screen.Fprintln(g, "msg", "cyan_black", s)
}

// Quit saratoga
func exit(g *gocui.Gui, args []string) {
	if len(args) == 1 { // exit 0
		screen.Fprintln(g, "msg", "green_black", "Gocui Good Bye!")
		return
	}
}

/* ************************************************************************** */

type cmdfunc func(*gocui.Gui, []string)

// Commands and function pointers to handle them
var cmdhandler = map[string]cmdfunc{
	"ca":   cmda,
	"cb":   cmdb,
	"cc":   cmdc,
	"ls":   ls,
	"quit": exit,
}

// Docmd -- Execute the command entered
func Docmd(g *gocui.Gui, s string) {
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
				fn(g, vals)
				return
			}
			screen.Fprintln(g, "msg", "bwhite_red", "Cannot execute:", vals[0])
			return
		}
	}
	screen.Fprintln(g, "msg", "bwhite_red", "Invalid command:", vals[0])
	return
}
