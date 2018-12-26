package views

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var errNoSpace = fmt.Errorf("max amount of views reached")
var errNotFound = fmt.Errorf("unknown view")

type viewManager struct {
	sync.Mutex
	aliases map[int8]string
	spots   [10]bool

	recentView int8

	addCb func()
}

func createViewManager(addCb func()) *viewManager {
	return &viewManager{
		aliases:    make(map[int8]string),
		addCb:      addCb,
		recentView: -1,
	}
}

func (vm *viewManager) AllViews() []string {
	var res []string
	for i := range vm.spots {
		res = append(res, strconv.Itoa(i))
	}
	return res
}

func (vm *viewManager) AddView(name string) (string, error) {
	vm.Lock()
	defer vm.Unlock()

	var num int8
	var found bool
	for i, ex := range vm.spots {
		if !ex {
			num, found = int8(i), true
			break
		}
	}
	if !found {
		return "", errNoSpace
	}
	vm.aliases[num] = name
	vm.spots[num] = true
	vm.recentView = num
	vm.addCb()

	return strconv.Itoa(int(num)), nil
}

func (vm *viewManager) RemoveView(name string) (string, error) {
	vm.Lock()
	defer vm.Unlock()

	num, err := vm.findForName(name)
	if err != nil {
		return "", err
	}
	delete(vm.aliases, num)
	vm.spots[num] = false
	return strconv.Itoa(int(num)), nil
}

func (vm *viewManager) GetRecent() string {
	return strconv.Itoa(int(vm.recentView))
}

func (vm *viewManager) SetRecent(name string) error {
	num, err := vm.findForName(name)
	if err != nil {
		return err
	}
	vm.recentView = num
	return nil
}

func (vm *viewManager) findForName(name string) (int8, error) {
	var num int8
	var found bool
	for i, n := range vm.aliases {
		if n == name {
			num, found = int8(i), true
			break
		}
	}
	if !found {
		return 0, errNotFound
	}
	return num, nil
}

func (vm *viewManager) String() string {
	var res []string
	for i, b := range vm.spots {
		if b {
			num := int8(i)
			if vm.recentView == num {
				res = append(res, fmt.Sprintf("*%d:%s", i, vm.aliases[num]))
			} else {
				res = append(res, fmt.Sprintf("%d:%s", i, vm.aliases[num]))
			}
		}
	}
	return strings.Join(res, " ")
}
