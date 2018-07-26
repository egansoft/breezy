package actions

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/egansoft/silly/config"
	"github.com/egansoft/silly/utils"
)

type Action interface {
	Handle(http.ResponseWriter, []string, []string)
}

type Cmd struct {
	script string
}

type Fs struct {
	root string
}

func NewCmd(urlPath []string, line string) (Action, error) {
	varToUrlIndex := make(map[string]int)
	for i, token := range urlPath {
		if _, exists := varToUrlIndex[token]; exists {
			return nil, fmt.Errorf("Duplicate var defined in %s", line)
		}
		if utils.TokenIsVar(token) {
			varToUrlIndex[token] = i
		}
	}

	cmdVars := utils.VarsInCmd(line)
	script := line
	for i, cmdVar := range cmdVars {
		_, exists := varToUrlIndex[cmdVar]
		if !exists {
			return nil, fmt.Errorf("Var %s used in %s but not defined in %v", cmdVar, line, urlPath)
		}

		varArg := fmt.Sprintf("$%v", i)
		script = strings.Replace(script, cmdVar, varArg, -1)
	}

	cmd := &Cmd{
		script: script,
	}
	return cmd, nil
}

func (c *Cmd) Handle(w http.ResponseWriter, args []string, residual []string) {
	bashArgs := []string{"-c", c.script}
	allArgs := append(bashArgs, args...)

	cmd := exec.Command(config.Shell, allArgs...)
	cmd.Stdout = w
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func NewFs(root string) (Action, error) {
	fs := &Fs{
		root: root,
	}
	return fs, nil
}

func (f *Fs) Handle(w http.ResponseWriter, args []string, residual []string) {
	pathEnd := strings.Join(residual, "/")
	path := f.root + "/" + pathEnd
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		io.WriteString(w, "uh oh")
	}
	defer file.Close()

	ctype := mime.TypeByExtension(path)
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	io.Copy(w, file)
}
