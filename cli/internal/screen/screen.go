package screen

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

type Controller struct {
	g *gocui.Gui

	ctrl *gocui.View

	initCh chan<- time.Duration
	logger *log.Logger
	vm     *viewManager

	skbs *keybindings
}

func New(logger *log.Logger, initCh chan<- time.Duration) *Controller {
	return &Controller{
		initCh: initCh,
		logger: logger,
		skbs:   emptyKeybindings(),
	}
}

func (c *Controller) Init() error {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return err
	}
	defer g.Close()

	g.InputEsc = true

	c.g = g
	c.vm = createViewManager(c.refreshCtrl)

	g.SetManagerFunc(c.mainLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err := c.skbs.Unbind(g)
		if err != nil {
			return err
		}
		err = c.focusCtrl()
		if err != nil {
			return err
		}
		c.ctrl.Clear()
		n, err := fmt.Fprint(c.ctrl, ":")
		if err != nil {
			return err
		}
		err = c.ctrl.SetCursor(n, 0)
		if err != nil {
			return err
		}
		c.ctrl.Editable = true
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("ctrl", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		buf := make([]byte, 1024)
		n, err := v.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		buf = buf[:n]
		c.logger.Printf("sending command: '%s'", buf)

		c.refreshCtrl()
		return c.skbs.Bind(g)
	}); err != nil {
		return err
	}

	c.initCh <- 200 * time.Millisecond
	if err := g.MainLoop(); err != nil {
		return err
	}
	return nil
}

func (c *Controller) AddTab(name string) (func([]byte), error) {
	numStr, err := c.vm.AddView(name)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}
	err = c.vm.SetRecent(name)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}
	v, err := c.g.SetViewOnTop(numStr)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}
	_, err = c.g.SetCurrentView(numStr)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	c.skbs.Add(newKey("", []rune(numStr)[0], gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetViewOnTop(numStr)
		if err != nil {
			return err
		}
		_, err = g.SetCurrentView(numStr)
		if err != nil {
			return err
		}
		err = c.vm.SetRecent(name)
		if err != nil {
			return err
		}
		c.refreshCtrl()
		return nil
	}))
	c.skbs.Bind(c.g)

	if err := c.g.SetKeybinding(numStr, gocui.KeyCtrlP, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.Clear()
		return nil
	}); err != nil {
		return nil, err
	}

	return func(b []byte) {
		c.g.Update(func(g *gocui.Gui) error {
			_, err = v.Write(b)
			if err != nil {
				c.logger.Println(err)
				return err
			}
			return nil
		})
	}, nil
}

func (c *Controller) mainLayout(g *gocui.Gui) error {
	x, y := g.Size()
	ctrl, err := g.SetView("ctrl", -1, y-2, x, y)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	c.ctrl = ctrl

	for _, name := range c.vm.AllViews() {
		_, err := g.SetView(name, -1, -1, x, y-2)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
	}

	return nil
}

func (c *Controller) refreshCtrl() {
	c.g.Update(func(g *gocui.Gui) error {
		c.ctrl.Clear()
		fmt.Fprint(c.ctrl, c.vm)
		return nil
	})
}

func (c *Controller) focusCtrl() error {
	_, err := c.g.SetCurrentView(c.ctrl.Name())
	if err != nil {
		c.logger.Println(err)
		return err
	}
	_, err = c.g.SetViewOnTop(c.ctrl.Name())
	if err != nil {
		c.logger.Println(err)
	}
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
