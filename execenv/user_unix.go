// +build !windows

package execenv

import (
	"fmt"
	"strconv"

	user "github.com/dnephin/go-os-user"
)

func valueFromUser(name string) (string, error) {
	currentUser, err := user.CurrentUser()
	if err != nil {
		return "", err
	}
	switch name {
	case "name":
		return currentUser.Name, nil
	case "uid":
		return strconv.Itoa(currentUser.Uid), nil
	case "gid":
		return strconv.Itoa(currentUser.Gid), nil
	case "home":
		return currentUser.Home, nil
	case "group":
		group, err := user.LookupGid(currentUser.Gid)
		return group.Name, err
	default:
		return "", fmt.Errorf("Unknown variable \"user.%s\"", name)
	}
}

func getUserName() (string, error) {
	currentUser, err := user.CurrentUser()
	return currentUser.Name, err
}
