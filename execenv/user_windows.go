package execenv

import (
	"fmt"
	"os/user"
)

func valueFromUser(name string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	switch name {
	case "name":
		return currentUser.Username, nil
	case "uid":
		return currentUser.Uid, nil
	case "gid":
		return currentUser.Gid, nil
	case "home":
		return currentUser.HomeDir, nil
	case "group":
		group, err := user.LookupGroupId(currentUser.Gid)
		return group.Name, err
	default:
		return "", fmt.Errorf("Unknown variable \"user.%s\"", name)
	}
}
