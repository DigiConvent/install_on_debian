package install_on_debian

import (
	"errors"
	"os/exec"
	"os/user"
)

func (u *userAccount) Delete() error {
	output, err := exec.Command("sudo", "userdel", "-r", u.name).Output()
	if err != nil {
		return errors.New(err.Error() + string(output))
	}
	return nil
}

type userAccount struct {
	name string
}

// if the user already exists, nil is returned
func CreateUser(name string) (*userAccount, error) {
	u, _ := user.Lookup(name)
	if u == nil {
		output, err := exec.Command("sudo", "useradd", "--create-home", name).Output()
		if err != nil {
			return nil, errors.New(err.Error() + string(output))
		}
	}

	return &userAccount{name: name}, nil
}
