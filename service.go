package install_on_debian

import (
	"errors"
	"os"
	"os/exec"
	"os/user"
)

func CreateUser(name string) (*userAccount, error) {
	u, _ := user.Lookup(name)
	if u == nil {
		output, err := exec.Command("sudo", "useradd", "--create-home", name).Output()
		if err != nil {
			return nil, errors.New(err.Error() + string(output))
		}
	}

	return &userAccount{name: &name}, nil
}

type userAccount struct {
	name *string
}

func (u *userAccount) Delete() error {
	output, err := exec.Command("sudo", "userdel", "-r", *u.name).Output()
	if err != nil {
		return errors.New(err.Error() + string(output))
	}
	return nil
}

func servicePath(name string) string {
	return `/etc/systemd/system/` + name + `.service`
}
func serviceFileContents(name string) string {
	return `
	[Unit]
	Description=` + name + `
	After=network.target
	
	[Service]
	Type=simple
	ExecStart=/home/` + name + `/main --run
	Restart=always
	User=` + name + `
	Group=` + name + `
	StandardOutput=journal
	StandardError=journal
	AmbientCapabilities=CAP_NET_BIND_SERVICE
	
	[Install]
	WantedBy=multi-user.target`
}
func (u *userAccount) InstallService() (*installedService, error) {
	// make sure that a service file is put into place
	if _, err := os.Stat(servicePath(*u.name)); err != nil {
		err := os.WriteFile(servicePath(*u.name), []byte(serviceFileContents(*u.name)), 0644)
		if err != nil {
			return nil, err
		}
	}
	return &installedService{name: u.name}, nil
}

type installedService struct {
	name *string
}

func (i *installedService) Uninstall() (*userAccount, error) {
	err := os.Remove(servicePath(*i.name))
	if err != nil {
		return nil, err
	}
	return &userAccount{name: i.name}, nil
}

func (i *installedService) Start() (*startedService, error) {
	for _, command := range []string{
		"systemctl daemon-reload",
		"systemctl enable" + *i.name,
		"systemctl start" + *i.name} {
		err := exec.Command("sudo", command).Run()
		if err != nil {
			return nil, err
		}
	}
	return &startedService{name: i.name}, nil
}

type startedService struct {
	name *string
}

func (s *startedService) Stop() (*installedService, error) {
	err := exec.Command("sudo", "systemctl stop"+*s.name).Run()
	if err != nil {
		return nil, err
	}
	return &installedService{
		name: s.name,
	}, nil
}
