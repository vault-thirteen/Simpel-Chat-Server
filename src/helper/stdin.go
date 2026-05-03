package helper

import (
	"fmt"
	"os"
	"testing"

	"golang.org/x/term"
)

const StdinPasswordForTests = "test"

func GetPasswordFromStdin(object string) (pwd string, err error) {
	msg := fmt.Sprintf("Enter password for %s:", object)
	fmt.Println(msg)

	// Be aware that 'stdin' does not work in Go's tests !
	if testing.Testing() {
		return StdinPasswordForTests, nil
	}

	var buf []byte
	buf, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
