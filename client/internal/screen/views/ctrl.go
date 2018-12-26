package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const ctrlViewName = "ctrl"
const ctrlByte byte = ':'

// Ctrl is the control view
type Ctrl struct {
	bindKey  BindFn
	insMode  ModeFn
	viewMode ModeFn

	x0, y0 int
	x, y   int

	// internal state
	oldView string
	oldBuf  string
}

func NewCtrl(bind BindFn, insMode, viewMode ModeFn, x0, y0, x, y int) *Ctrl {
	return &Ctrl{
		bindKey:  bind,
		insMode:  insMode,
		viewMode: viewMode,

		x0: x0, y0: y0,
		x: x, y: y,
	}
}

func (c *Ctrl) SetKeybindings(g *gocui.Gui) error {
	if err := c.bindKey("", ':', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err := c.insMode()
		if err != nil {
			return err
		}
		return c.focus(g)
	}); err != nil {
		return err
	}

	return g.SetKeybinding(ctrlViewName, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err := c.defocus(g)
		if err != nil {
			return err
		}
		return c.viewMode()
	})
}

func (c *Ctrl) Layout(g *gocui.Gui) error {
	_, err := g.SetView(ctrlViewName, c.x0, c.y0, c.x, c.y)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	return nil
}

func (c *Ctrl) Update(g *gocui.Gui, content string) {
	v, _ := g.View(ctrlViewName)
	v.Clear()
	fmt.Fprint(v, content)
}

func (c *Ctrl) focus(g *gocui.Gui) error {
	v := g.CurrentView()
	if v != nil {
		c.oldView = v.Name()
	}
	v, err := g.SetCurrentView(ctrlViewName)
	if err != nil {
		return err
	}
	c.oldBuf = v.ViewBuffer()
	v.Clear()
	n, err := v.Write([]byte{ctrlByte})
	if err != nil {
		return err
	}
	v.Editable = true
	g.Cursor = true

	return v.SetCursor(n, 0)
}

func (c *Ctrl) defocus(g *gocui.Gui) error {
	v := g.CurrentView()
	if v == nil || v.Name() != ctrlViewName {
		return nil
	}
	g.Cursor = false
	v.Editable = false
	v.Clear()
	v.Write([]byte(c.oldBuf))
	g.SetCurrentView(c.oldView)
	return nil
}
