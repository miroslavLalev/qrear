package screen

import (
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/miroslavLalev/qrear/client/internal/screen/views"
)

type Controller struct {
	g *gocui.Gui

	ctrl    *views.Ctrl
	courier *views.Courier

	initCh chan<- time.Duration
	logger *log.Logger

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

	x, y := g.Size()

	c.ctrl = views.NewCtrl(func(
		name string,
		key interface{},
		mod gocui.Modifier,
		press func(g *gocui.Gui, v *gocui.View) error,
	) error {
		c.skbs.Add(newKey(name, key, mod, press))
		return c.skbs.Bind(g)
	}, func() error {
		return c.skbs.Unbind(g)
	}, func() error {
		return c.skbs.Bind(g)
	}, -1, y-2, x, y)

	c.courier = views.NewCourier(func(
		name string,
		key interface{},
		mod gocui.Modifier,
		press func(g *gocui.Gui, v *gocui.View) error,
	) error {
		c.skbs.Add(newKey(name, key, mod, press))
		return c.skbs.Bind(g)
	}, func() error {
		return c.skbs.Unbind(g)
	}, func() error {
		return c.skbs.Bind(g)
	}, func() {
		c.ctrl.Update(g, c.courier.String())
	}, -1, -1, x, y-2)

	g.SetManager(c.ctrl, c.courier)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := c.ctrl.SetKeybindings(g); err != nil {
		return err
	}

	if err := c.courier.SetKeybindings(g); err != nil {
		return err
	}

	c.initCh <- 500 * time.Millisecond
	if err := g.MainLoop(); err != nil {
		return err
	}
	return nil
}

func (c *Controller) AddTab(name string) (func([]byte), error) {
	v, err := c.courier.AddLayer(c.g, name)
	if err != nil {
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
