package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

var args struct {
	NeedConfirm    bool     `arg:"-n, --needConfirm" default:"false" help:"confirm before remove every single entry"`
	Recursion      bool     `arg:"-r, --recursion" default:"false"`
	Except         bool     `arg:"-e, --except" default:"false"`
	RemoveEmptyDir bool     `arg:"-d, --removeEmptyDir" default:"true"`
	DirOrFile      string   `arg:"positional" help:"target file or dir"`
	Pattern        []string `arg:"positional" help:"partern to match"`
}

const seperator = string(os.PathSeparator)

func matchPattern(name string) bool {
	if len(args.Pattern) == 0 {
		return true
	}
	for _, item := range args.Pattern {
		if item == name {
			return true
		}
	}
	return false
}
func confirm(element string, name string) bool {
	fmt.Printf("remove %s %s:(y/N)", element, name)
	var c rune
	_, err := fmt.Scanf("%c", &c)
	if err != nil {
		panic(err)
	}
	return c == 'y'
}
func main() {
	arg.MustParse(&args)
	info, err := os.Stat(args.DirOrFile)
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		dir()
		return
	}
	// is normal file
	file()
}
func file() {
	if args.NeedConfirm && !confirm("file", args.DirOrFile) {
		return
	}
	os.Remove(args.DirOrFile)
}
func dir() {
	if !args.Recursion {
		panic(fmt.Errorf("target %s is a dir", args.DirOrFile))
	}
	if strings.HasSuffix(args.DirOrFile, seperator) {
		args.DirOrFile = args.DirOrFile[:len(args.DirOrFile)-2]
	}
	if len(args.Pattern) == 0 {
		// remove all
		os.RemoveAll(args.DirOrFile)
		return
	}
	// start to recursion
	WalkDir(args.DirOrFile)
}
func WalkDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		name := strings.Join([]string{dir, entry.Name()}, seperator)
		if !entry.IsDir() { // is a file
			if matchPattern(entry.Name()) {
				if args.NeedConfirm && !confirm("file", name) {
					continue
				}
				err := os.Remove(name)
				if err != nil {
					panic(err)
				}
			}
			continue
		}
		if args.Recursion {
			WalkDir(name)
		}
		if args.RemoveEmptyDir {
			if args.NeedConfirm && !confirm("dir", name) {
				continue
			}
			err = os.Remove(name) // only the empty dir would be successfully removed
			if err != nil {
				panic(err)
			}
		}
	}
}
