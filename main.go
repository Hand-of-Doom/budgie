package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func removeComments(file string) string {
	pattern := regexp.MustCompile(`#.*`)
	return pattern.ReplaceAllString(file, "")
}

type target struct {
	name string
	src  string
	deps []*target
}

type targetsSet map[string]*target

func extractTargetScope(target string, file string) string {
	beforeCurTarget := strings.Split(file, target)[0]
	regexpTargets := regexp.MustCompile(`target[ \S\n\t]+end`)
	return regexpTargets.ReplaceAllString(beforeCurTarget, "")
}

func splitIntoTargetSet(file string) targetsSet {
	pattern := regexp.MustCompile(`(?U)target (?P<name>\S+)(\((?P<deps>[\S ]+)?\))?:(?P<src>[ \S\n\t]+)?\bend\b`)
	namePos := pattern.SubexpIndex("name")
	depsPos := pattern.SubexpIndex("deps")
	srcPos := pattern.SubexpIndex("src")

	targets := targetsSet{}

	groups := pattern.FindAllStringSubmatch(file, -1)
	matches := pattern.FindAllString(file, -1)

	var insertDeps func(donor *target, depNames []string)
	insertDeps = func(donor *target, depNames []string) {
		if depNames[0] == "" {
			return
		}
		for _, depName := range depNames {
			if targets[depName] == nil {
				panic("no target named " + depName)
			}
			for i, group := range groups {
				if depName != group[namePos] {
					continue
				}
				name := group[namePos]
				var dep *target
				if _, ok := targets[name]; ok {
					dep = targets[name]
				} else {
					dep = new(target)
					dep.name = name
					scope := extractTargetScope(matches[i], file)
					dep.src = scope + group[srcPos]

					deps := strings.Split(group[depsPos], " ")
					insertDeps(dep, deps)
					targets[name] = dep
				}
				donor.deps = append(donor.deps, dep)
			}
		}
	}

	for i, group := range groups {
		name := group[namePos]
		if _, ok := targets[name]; !ok {
			t := new(target)
			t.name = group[namePos]
			scope := extractTargetScope(matches[i], file)
			t.src = scope + group[srcPos]

			deps := strings.Split(group[depsPos], " ")
			insertDeps(t, deps)

			targets[name] = t
		}
	}

	return targets
}

func runTarget(targetName string, targets targetsSet) error {
	target, ok := targets[targetName]
	if !ok {
		panic("no target named " + targetName)
	}

	for _, dep := range target.deps {
		err := runTarget(dep.name, targets)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command("/bin/sh", "-c", target.src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

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
	const ext = ".budgie"
	if len(args) == 2 {
		fname = args[1]
		if filepath.Ext(fname) != ext {
			panic("the file extension must be specified as " + ext)
		}
	} else {
		fname = "tasks" + ext
	}

	targetName := args[0]

	fbytes, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	file := string(fbytes)
	file = removeComments(file)
	targets := splitIntoTargetSet(file)

	return runTarget(targetName, targets)
}

func main() {
	if err := runApplication(); err != nil {
		panic(err)
	}
}
