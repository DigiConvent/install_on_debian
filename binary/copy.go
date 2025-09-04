package binary

import (
	"fmt"
	"os"
	"path"

	"github.com/DigiConvent/install_on_debian/utils"
)

// this packages makes sure that the binary is in place and handles "it"

type BinaryOperations interface {
	uri() string
	HardLinkToHome() error
}

type Binary struct {
	name string
}

func (b *Binary) uri() string {
	u, err := os.Executable()
	if err != nil {
		return ""
	}
	return u
}

func (b *Binary) HardLinkToHome() error {
	target := path.Join("/home", b.name, "main")
	if utils.FileExists(target) {
		original, err := os.Stat(b.uri())
		if err != nil {
			return err
		}

		hardLink, err := os.Stat(target)
		if err != nil {
			return err
		}
		isSameFile := os.SameFile(original, hardLink)
		if isSameFile {
			err = os.Remove(target)
			if err != nil {
				fmt.Println(err)
			}
		}
		return b.HardLinkToHome()
	}

	return os.Link(b.uri(), target)
}

func New(name string) BinaryOperations {
	return &Binary{name: name}
}
