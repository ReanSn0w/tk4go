package tools

import (
	"context"
	"os"

	"golang.org/x/term"
)

// AnyKeyToExit - реализует выход из приложения после нажатия любой клавици пользователем
//
// https://stackoverflow.com/questions/15159118/read-a-character-from-standard-input-in-go-without-pressing-enter
func AnyKeyToExit(log Logger, cancel context.CancelFunc) {
	log.Logf("[INFO] press any key to exit")
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	b := make([]byte, 1)
	os.Stdin.Read(b)
	cancel()
}
