# testgocui
Example of gocui user interface with multiple views and colours.

This shows an example of 4 views with the bottom left "cmd" view handling command line input;
 the top "msg" view showing messages; the bottom right "err" view showing error info and the optional "packet" view that can be tunred on and off (brought to the front).

Each command line input runs as an independant go routine to carry out the task in the background with output's showing in the "msg", "packet" and "err" views.

I a am working on getting vertical scrolling going properly. Currently it does not scroll above,
or below the cursor view. Although the "bufer" is retained.

Is a reasonably flexible construct where all you have to do is edit "cli.go" to add in your own command functions and change the associated help and command funtion pointer map's.

Enjoy.
