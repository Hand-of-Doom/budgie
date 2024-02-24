package main

import (
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func extractScope(targetName string, file string) string {
	before, after, ok := strings.Cut(file, targetName)
	if !ok {
		return before
	}
	charsetAfter := []rune(after)
	var next bool
	for i, char := range charsetAfter {
		if char == '}' {
			if next == false {
				charsetAfter = charsetAfter[:i+1]
				break
			}
			next = false
			continue
		}
		if char == '{' {
			next = true
		}
	}

	scope := before + targetName + string(charsetAfter)
	return scope
}

func runTarget(name string, file string) error {
	scope := extractScope(name, file)
	scope += "\n"+name

	cmd := exec.Command("sh", "-c", scope)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}


func runApplication() error {
	args := os.Args[1:]
	if len(args) < 1 {
		panic("budgie requires a target name as the first argument")
	}
	if len(args) > 2 {
		panic("you can only pass less than 2 arguments")
	}

	var fname string
	if len(args) == 2 {
		fname = args[1]
	} else {
		fname = "tasks.sh"
	}

	targetName := args[0]
	charset := []rune(targetName)
	charset[0] = unicode.ToUpper(charset[0])
	targetName = string(charset)

	fbytes, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	file := string(fbytes)
	return runTarget(targetName, file)
}

func main() {
	if err := runApplication(); err != nil {
		panic(err)
	}
}
