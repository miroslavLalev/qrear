package views

import (
	"github.com/jroimartin/gocui"
)

type Courier struct {
	bindKey  BindFn
	insMode  ModeFn
	viewMode ModeFn

	vm *viewManager

	x0, y0 int
	x, y   int
}

func NewCourier(bindKey BindFn, insMode, viewMode ModeFn, addCb func(), x0, y0, x, y int) *Courier {
	return &Courier{
		bindKey:  bindKey,
		insMode:  insMode,
		viewMode: viewMode,
		vm:       createViewManager(addCb),

		x0: x0, y0: y0,
		x: x, y: y,
	}
}

func (c *Courier) SetKeybindings(g *gocui.Gui) error {
	for _, layer := range c.vm.AllViews() {
		if err := g.SetKeybinding(layer, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			v.Editable = false
			g.Cursor = false
			return c.viewMode()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (c *Courier) Layout(g *gocui.Gui) error {
	for _, layer := range c.vm.AllViews() {
		if _, err := g.SetView(
			layer, c.x0, c.y0, c.x, c.y,
		); err != nil && err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

func (c *Courier) String() string {
	return c.vm.String()
}

func (c *Courier) AddLayer(g *gocui.Gui, name string) (*gocui.View, error) {
	numStr, err := c.vm.AddView(name)
	if err != nil {
		return nil, err
	}
	v, _ := g.SetViewOnTop(numStr)
	g.SetCurrentView(numStr)
	c.bindKey("", getKeyForName(numStr), gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		c.vm.SetRecent(name)
		g.SetViewOnTop(numStr)
		g.SetCurrentView(numStr)

		// TODO: This does not belong here!
		c.vm.addCb()
		return nil
	})

	c.bindKey(numStr, 'a', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.Editable = true
		g.Cursor = true
		return c.insMode()
	})

	return v, nil
}

func getKeyForName(name string) rune {
	return []rune(name)[0] // quick'n'dirty
}
