package main

import l4g "log4go"

func main() {
	// Load the configuration (isn't this easy?)
	l4g.LoadConfiguration("example.xml")
	l4g.Debug("hello world")
	l4g.Trace("hello world")

	// And now we're ready!
	l4g.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	l4g.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	l4g.Info("About that time, eh chaps?")
}
