package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"unicode"
)

func extractScope(targetName string, file string) string {
	pattern := fmt.Sprintf(`(?sU).*%s.*}`, targetName)
	regexpScope := regexp.MustCompile(pattern)

	return regexpScope.FindString(file)
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
		panic("you can pass less than 2 arguments")
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
