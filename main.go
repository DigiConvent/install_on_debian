package install_on_debian

import (
	"github.com/DigiConvent/install_on_debian/binary"
	"github.com/DigiConvent/install_on_debian/systemctl"
	user "github.com/DigiConvent/install_on_debian/user"
)

func InstallThisBinary(name string) error {
	u, err := user.CreateOrGetUser(name)
	if err != nil {
		return err
	}

	sysCtl, err := systemctl.Get(name)
	if err != nil {
		return err
	}

	bin := binary.New(name)
	err = bin.HardLinkToHome()
	if err != nil {
		return err
	}

	sysCtl.User = u
	if sysCtl.IsInstalled() {
		return nil
	}
	_, err = sysCtl.Install("")
	if err != nil {
		return err
	}
	return nil
}

func UninstallThisBinary(name string) error {
	sysCtl, err := systemctl.Get(name)
	if err != nil {
		return err
	}

	if sysCtl.IsRunning() {
		_, err := sysCtl.Stop()
		if err != nil {
			return err
		}
	}

	_, err = sysCtl.Uninstall()
	if err != nil {
		return err
	}

	sysCtl.User.Delete()

	return nil
}
