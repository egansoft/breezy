package config

import (
	"fmt"
)

var (
	Port        = 8080
	DebugMode   = false
	Shell       = "sh"
	ShellPrefix = ""
	ShellArg    = "-c"
	ShellBind   = ShBind
)

func SetShell(shell string) {
	switch shell {
	case "sh":
		ShellPrefix = ""
		ShellArg = "-c"
		ShellBind = ShBind
	case "bash":
		ShellPrefix = ""
		ShellArg = "-c"
		ShellBind = ShBind
	case "zsh":
		ShellPrefix = ""
		ShellArg = "-c"
		ShellBind = ShBind
	case "python":
		ShellPrefix = "import sys\n"
		ShellArg = "-c"
		ShellBind = PythonBind
	case "python3":
		ShellPrefix = "import sys\n"
		ShellArg = "-c"
		ShellBind = PythonBind
	case "node":
		ShellPrefix = ""
		ShellArg = "-e"
		ShellBind = NodeBind
	case "ruby":
		ShellPrefix = ""
		ShellArg = "-e"
		ShellBind = RubyBind
	}

	Shell = shell
}

func ShBind(idx int) string {
	return fmt.Sprintf("$%v", idx)
}

func PythonBind(idx int) string {
	return fmt.Sprintf("sys.argv[%v]", idx+1)
}

func NodeBind(idx int) string {
	return fmt.Sprintf("process.argv[%v]", idx+1)
}

func RubyBind(idx int) string {
	return fmt.Sprintf("ARGV[%v]", idx)
}
