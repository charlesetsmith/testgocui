# testgocui
Example of gocui user interface

This shows an example of two views with the bottom view handling command line input and the top view handling messages.

Each command line input runs a go routine to carry out the task in the background with any output's showing in the 
message window.

The command line view can be scrolled up and down.
Next stage will be to implemnt "shell like" copy and paste !! (to run the last command) !xx where xx is the number
of the command line to run again. "ctrl c" to copy "ctrl v" to paste.

Is a reasonably flexible construct where all you have to do is edit "cli.go" to add in your own command functions, fiddle
the help and command funtion pointer map's.

Enjoy.
