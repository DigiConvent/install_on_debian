package user

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/DigiConvent/install_on_debian/utils"
)

const sudoersFolder = "/etc/sudoers.d"

func CreateOrGetUser(name string) (*OsUserAccount, error) {
	_, err := user.Lookup(name)
	userAccount := &OsUserAccount{
		name: &name,
	}
	if err != nil {
		output, err := exec.Command("sudo", "useradd", "--create-home", name).Output()
		if err != nil {
			return nil, errors.New(err.Error() + string(output))
		}

		if utils.FileExists(userAccount.sudoersFile()) {
			fmt.Println("Sudoersfile already exists, deleting")
			err := os.Remove(userAccount.sudoersFile())
			fmt.Println(err)
		}

		sctl := "/bin/systemctl"
		cmds := []string{sctl + " daemon-reload"}       // can restart daemon
		cmds = append(cmds, sctl+" * "+name+".service") // can start, restart, stop, show, enable, etc...

		file, err := os.OpenFile(userAccount.sudoersFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0x440)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		varName := strings.ToUpper(name) + "_EXCEPTIONS"
		_, err = file.WriteString("Cmnd_Alias " + varName + " = " + strings.Join(cmds, ", ") + "\n\n" + name + " ALL=(ALL) NOPASSWD: " + varName + "\n")
		if err != nil {
			return nil, err
		}
	}

	return userAccount, nil
}

type OsUserAccount struct {
	name *string
}

func (u *OsUserAccount) sudoersFile() string {
	return path.Join(sudoersFolder, *u.name)
}

func (u *OsUserAccount) Delete() error {
	output, err := exec.Command("sudo", "userdel", "-r", *u.name).Output()
	if err != nil {
		return errors.New("could not userdel -r " + *u.name + " " + err.Error() + string(output))
	}

	err = os.Remove(u.sudoersFile())
	if err != nil {
		return errors.New("could not delete sudoersfile")
	}
	return nil
}
