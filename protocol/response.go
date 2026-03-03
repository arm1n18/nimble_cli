package protocol

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Underline = "\033[4m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[95m"
	Cyan      = "\033[36m"
	Gray      = "\033[37m"
	White     = "\033[97m"
)

func Error(format string, a ...any) {
	fmt.Printf(Red+"Error: "+format+Reset+"", a...)
}

func FatalError(format string, a ...any) {
	log.Printf(Red+"Error: "+format+Reset+"", a...)
	os.Exit(1)
}

func Help(format string, a ...any) {
	fmt.Printf(Yellow+format+Reset+"\n", a...)
}

func BoldText(format string, a ...any) {
	fmt.Printf(Bold+format+Reset+"\n", a...)
}

func Commad(cmd, desc, format string, lenCmd, lenDesc int) {
	fmt.Printf(Bold+Green+"%-*s"+Reset+"%-*s"+Magenta+"%s"+Reset+"\n", lenCmd, cmd, lenDesc, desc, format)
}

func ClearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
