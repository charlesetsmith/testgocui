// Test aplication and example of cli interface with a command and message split pane window.

package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charlesetsmith/testgocui/screen"
	"github.com/jroimartin/gocui"
)

// Cmdinfo - Previos Command history, prompt info
type Cmdhist struct {
	Commands []string
	Prompt   string
	Ppad     int // Number of pad characters around prompt e.g. prompt[99]: would be 3
	Curline  int // What is the current command line # we are on
}

type Cmdfunc func(*gocui.Gui, []string, Cmdhist)

type Cmd struct {
	Usage string
	Help  string
}

var Commands = map[string]Cmd{
	"ca":    {Usage: "ca [arg]...", Help: "Command Example a"},
	"cb":    {Usage: "cb [arg]...", Help: "Command Example b"},
	"cc":    {Usage: "cc [arg]...", Help: "Command Example c"},
	"buf":   {Usage: "buf", Help: "Show Buffer"},
	"ls":    {Usage: "ls", Help: "History of commands entered"},
	"quit":  {Usage: "quit", Help: "Bye!"},
	"exit":  {Usage: "exit", Help: "Bye!"},
	"help":  {Usage: "help", Help: "List of available commands"},
	"usage": {Usage: "usage", Help: "List of available commands"},
	"?":     {Usage: "?", Help: "List of available commands"},
}

var Commandfuncs = map[string]Cmdfunc{
	"ca":    cmda,
	"cb":    cmdb,
	"cc":    cmdc,
	"buf":   cmdbuf,
	"ls":    ls,
	"quit":  exit,
	"exit":  exit,
	"help":  usage,
	"usage": usage,
	"?":     usage,
}

// The different command line input handlers

// cmda [args]...
func cmda(g *gocui.Gui, args []string, cmds Cmdhist) {
	screen.MsgPrintln(g, "green_black", "Command A", args)
}

// cmdb [args]...
func cmdb(g *gocui.Gui, args []string, cmds Cmdhist) {
	screen.MsgPrintln(g, "green_black", "Command B", args)
}

// cmdc [args]...
func cmdc(g *gocui.Gui, args []string, cmds Cmdhist) {
	screen.MsgPrintln(g, "green_black", "Command B", args)
}

// ls - list the history of commands to the msg window
func ls(g *gocui.Gui, args []string, cmds Cmdhist) {
	var s string

	for i := 0; i < len(cmds.Commands); i++ {
		s += fmt.Sprintf("%d=%s\n", i, cmds.Commands[i])
	}
	screen.MsgPrintln(g, "cyan_black", s)
}

func cmdbuf(g *gocui.Gui, args []string, cmds Cmdhist) {
	// var b string
	v, _ := g.View("cmd")
	s := v.Buffer()
	screen.MsgPrintln(g, "", s)
}

// Quit saratoga
func exit(g *gocui.Gui, args []string, cmds Cmdhist) {
	if len(args) == 1 { // exit 0
		screen.MsgPrintln(g, "green_black", "Gocui Good Bye!")
		return
	}
}

// usage - sort list usage of available commands and help
func usage(g *gocui.Gui, args []string, cmds Cmdhist) {
	s := "CtrlSpace - Rotate Between Views\nCtrlP - Show/Hode Packet View\n\n"
	keys := make([]string, 0, len(Commands))
	for k := range Commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s += fmt.Sprintf("%s: %s\n", Commands[k].Usage, Commands[k].Help)
	}
	screen.MsgPrintln(g, "cyan_black", s)
}

/* ************************************************************************** */

// Docmd -- Execute the command entered
func Docmd(g *gocui.Gui, s string, cmds Cmdhist) {
	if s == "" { // Handle just return
		return
	}
	s = strings.TrimSpace(s)  // Get rid of leading and trailing whitespace
	vals := strings.Fields(s) // Split each field into a slice of strings
	// Lookup the command and execute it if it is a valid command!
	if Commands[vals[0]].Help != "" {
		fn := Commandfuncs[vals[0]]
		fn(g, vals, cmds)
		return
	}
	screen.MsgPrintln(g, "red_black", "Invalid command:", vals[0])
}
