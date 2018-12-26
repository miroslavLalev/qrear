package views

import "github.com/jroimartin/gocui"

type BindFn func(
	name string,
	key interface{},
	mod gocui.Modifier,
	press func(g *gocui.Gui, v *gocui.View) error,
) error

type ModeFn func() error
