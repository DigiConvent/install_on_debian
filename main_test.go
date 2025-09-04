package install_on_debian_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DigiConvent/install_on_debian"
)

func Main() {
	fmt.Println("Start of the program")
	time.Sleep(10 * time.Second)
	fmt.Println("End of the program")
	os.Exit(0)
}

func TestInstallOnDebian(t *testing.T) {
	err := install_on_debian.InstallThisBinary("main_test")
	if err != nil {
		t.Fatal(err)
	}
}
