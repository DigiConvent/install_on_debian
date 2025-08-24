package install_on_debian

import (
	"os"
	"os/exec"
)

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

func (u *userAccount) InstallService(name string) (*installedService, error) {
	// make sure that a service file is put into place
	if _, err := os.Stat(servicePath(name)); err != nil {
		err := os.WriteFile(servicePath(name), []byte(serviceFileContents(name)), 0644)
		if err != nil {
			return nil, err
		}
	}
	return &installedService{name: name}, nil
}

type installedService struct {
	name string
}

func (i *installedService) Start(name string) (*startedService, error) {
	for _, command := range []string{
		"systemctl daemon-reload",
		"systemctl enable" + name,
		"systemctl start" + name} {
		err := exec.Command("sudo", command).Run()
		if err != nil {
			return nil, err
		}
	}
	return &startedService{name: name}, nil
}
func (i *installedService) UninstallService(name string) (*userAccount, error) {
	err := os.Remove(servicePath(name))
	if err != nil {
		return nil, err
	}
	return &userAccount{name: name}, nil
}

type startedService struct {
	name string
}

func (s *startedService) StopService() error {
	return exec.Command("sudo", "systemctl stop"+s.name).Run()
}
