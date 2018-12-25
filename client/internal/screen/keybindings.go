package screen

import (
	"github.com/jroimartin/gocui"
)

type keybindings []*keybinding

func emptyKeybindings() *keybindings {
	k := make(keybindings, 0, 128)
	return &k
}

func (k *keybindings) Add(key *keybinding) {
	*k = append(*k, key)
}

func (k keybindings) Bind(g *gocui.Gui) error {
	for _, key := range k {
		if key.bound {
			continue
		}
		err := g.SetKeybinding(key.Namespace, key.Key, key.Modifier, key.Exec)
		if err != nil {
			return err
		}
		key.bound = true
	}
	return nil
}

func (k keybindings) Unbind(g *gocui.Gui) error {
	for _, key := range k {
		if !key.bound {
			continue
		}
		err := g.DeleteKeybinding(key.Namespace, key.Key, key.Modifier)
		if err != nil {
			return err
		}
		key.bound = false
	}
	return nil
}

type keybinding struct {
	Namespace string
	Key       interface{}
	Modifier  gocui.Modifier
	Exec      func(*gocui.Gui, *gocui.View) error

	bound bool
}

func newKey(
	namespace string,
	key interface{},
	modifier gocui.Modifier,
	exec func(*gocui.Gui, *gocui.View) error,
) *keybinding {
	return &keybinding{
		Namespace: namespace,
		Key:       key,
		Modifier:  modifier,
		Exec:      exec,
	}
}
