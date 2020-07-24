package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
)

const (
	LUNCHY_VERSION = "0.2.1"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func fileCopy(src string, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}

	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	return d.Close()
}

// Plist holds launch agent name and plist path
type Plist struct {
	Name string
	Path string
}

type option struct {
	name  string
	value interface{}
}

func findPlists(dir string, options ...option) []Plist {
	plists := []Plist{}
	args := []string{"-L", dir, "-name", "*.plist", "-type", "f"}
	for _, option := range options {
		args = append(args, option.name, fmt.Sprintf("%v", option.value))
	}
	output, err := exec.Command("find", args...).Output()
	if err != nil {
		return plists
	}

	for _, plistPath := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		name := strings.Replace(path.Base(plistPath), ".plist", "", 1)
		plists = append(plists, Plist{name, plistPath})
	}

	sort.SliceStable(plists, func(i, j int) bool {
		return plists[i].Name < plists[j].Name
	})

	return plists
}

func getPlists() []Plist {
	dirs := []string{"/Library/LaunchAgents", path.Join(os.Getenv("HOME"), "/Library/LaunchAgents")}

	isRoot := os.Geteuid() == 0
	if isRoot {
		dirs = append(dirs, "/Library/LaunchDaemons", "/System/Library/LaunchDaemons")
	}

	plists := []Plist{}
	for _, dir := range dirs {
		plists = append(plists, findPlists(dir)...)
	}

	sort.SliceStable(plists, func(i, j int) bool {
		return plists[i].Name < plists[j].Name
	})

	return plists
}

func getPlist(name string) Plist {
	plists := []Plist{}
	for _, plist := range getPlists() {
		if strings.Index(plist.Name, name) != -1 {
			plists = append(plists, plist)
		}
	}

	if len(plists) == 0 {
		fatal("no launch agent found matching:", name)
	}

	if len(plists) > 1 {
		var matches strings.Builder
		for _, plist := range plists {
			matches.WriteString("\n")
			matches.WriteString(plist.Name)
		}
		fatal("multiple launch agents found matching:", name, "\n\nmatches found are:", matches.String())
	}

	return plists[0]
}

func sliceIncludes(slice []string, match string) bool {
	for _, val := range slice {
		if val == match {
			return true
		}
	}

	return false
}

func printUsage() {
	fmt.Printf("Lunchy %s, the friendly launchctl wrapper\n", LUNCHY_VERSION)
	fmt.Println("Usage: lunchy [start|stop|restart|list|status|install|show|edit|remove|scan] [options]")
}

func printList(args []string) {
	pattern := ""

	if len(args) > 0 {
		pattern = args[0]
	}

	for _, plist := range getPlists() {
		if strings.Index(plist.Name, pattern) != -1 {
			fmt.Println(plist.Name)
		}
	}
}

func printStatus(args []string) {
	out, err := exec.Command("launchctl", "list").Output()

	if err != nil {
		fatal("failed to get process list")
	}

	pattern := ""

	if len(args) > 0 {
		pattern = args[0]
	}

	installed := []string{}
	for _, plist := range getPlists() {
		installed = append(installed, plist.Name)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	for _, line := range lines {
		chunks := strings.Split(line, "\t")

		if chunks[2] == "Label" {
			fmt.Println(line)
			continue
		}

		if len(pattern) > 0 {
			if strings.Index(chunks[2], pattern) != -1 {
				if sliceIncludes(installed, chunks[2]) {
					fmt.Println(line)
				}
			}
		} else {
			if sliceIncludes(installed, chunks[2]) {
				fmt.Println(line)
			}
		}
	}
}

func exitWithInvalidArgs(args []string, msg string) {
	if len(args) < 1 {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func startDaemons(args []string) {
	// Check if name pattern is not given and try profiles
	if len(args) == 0 {
		if profileExists() {
			startProfile()
			return
		}
		exitWithInvalidArgs(args, "name required")
	}

	name := args[0]
	startDaemon(getPlist(name))
}

func startDaemon(plist Plist) {
	_, err := exec.Command("launchctl", "load", plist.Path).Output()

	if err != nil {
		fmt.Println("failed to start", plist.Name)
		return
	}

	fmt.Println("started", plist.Name)
}

func stopDaemons(args []string) {
	// Check if name pattern is not given and try profiles
	if len(args) == 0 {
		if profileExists() {
			stopProfile()
			return
		}
		exitWithInvalidArgs(args, "name required")
	}

	name := args[0]
	stopDaemon(getPlist(name))
}

func stopDaemon(plist Plist) {
	_, err := exec.Command("launchctl", "unload", plist.Path).Output()

	if err != nil {
		fmt.Println("failed to stop", plist.Name)
		return
	}

	fmt.Println("stopped", plist.Name)
}

func restartDaemons(args []string) {
	// Check if name pattern is not given and try profiles
	if len(args) == 0 {
		if profileExists() {
			restartProfile()
			return
		}
		exitWithInvalidArgs(args, "name required")
	}

	name := args[0]
	plist := getPlist(name)
	stopDaemon(plist)
	startDaemon(plist)
}

func showPlist(args []string) {
	exitWithInvalidArgs(args, "name required")

	name := args[0]
	printPlistContent(getPlist(name))
}

func printPlistContent(plist Plist) {
	contents, err := ioutil.ReadFile(plist.Path)

	if err != nil {
		fatal("unable to read plist")
	}

	fmt.Printf(string(contents))
}

func editPlist(args []string) {
	exitWithInvalidArgs(args, "name required")

	name := args[0]
	editPlistContent(getPlist(name))
}

func editPlistContent(plist Plist) {
	editor := os.Getenv("EDITOR")

	if len(editor) == 0 {
		fatal("EDITOR environment variable is not set")
	}

	cmd := exec.Command(editor, plist.Path)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.Start()
	cmd.Wait()
}

func installPlist(args []string) {
	exitWithInvalidArgs(args, "path required")

	srcPath := args[0]

	if !fileExists(srcPath) {
		fatal("source file does not exist")
	}

	info, _ := os.Stat(srcPath)
	dirs := []string{
		path.Join(os.Getenv("HOME"), "/Library/LaunchAgents"),
		"/Library/LaunchAgents",
	}

	for _, dir := range dirs {
		if !fileExists(dir) {
			continue
		}

		destPath := path.Join(dir, info.Name())

		if fileExists(destPath) && os.Remove(destPath) != nil {
			fatal("unable to delete existing plist")
		}

		if fileCopy(srcPath, destPath) != nil {
			fatal("failed to copy file")
		}

		fmt.Println(srcPath, "installed to", dir)
		return
	}

}

func removePlist(args []string) {
	exitWithInvalidArgs(args, "name required")

	name := args[0]
	plist := getPlist(name)
	if os.Remove(plist.Path) == nil {
		fmt.Println("removed", plist.Path)
	} else {
		fmt.Println("failed to remove", plist.Path)
	}
}

func scanPath(args []string) {
	exitWithInvalidArgs(args, "path required")

	options := []option{}
	dir := path.Join(os.Getenv("HOME"), "/Library/LaunchAgents")

	// This is a handy override to find all homebrew-based lists
	if dir == "homebrew" || dir == "Homebrew" {
		prefix := "/usr/local"
		output, _ := exec.Command("brew", "--prefix").Output()
		if output != nil {
			prefix = strings.TrimSpace(string(output))
		}
		dir = path.Join(prefix, "/Cellar")
		options = append(options, option{"-maxdepth", 3})
	}

	for _, plist := range findPlists(dir, options...) {
		fmt.Println(plist.Name)
	}
}

// Get full path to lunchy profile file
func profilePath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path.Join(dir, "/.lunchy")
}

// Check if profile file exists
func profileExists() bool {
	return fileExists(profilePath())
}

// Get daemon names specified in lunchy profile
func readProfile() []string {
	path := profilePath()
	if path == "" {
		return []string{}
	}

	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}
	}

	result := []string{}
	lines := strings.Split(strings.TrimSpace(string(buff)), "\n")

	for _, l := range lines {
		line := strings.TrimSpace(l)

		// Skip comments (starts with #)
		if line[0] == 35 {
			continue
		}

		result = append(result, line)
	}

	return result
}

func plistsAction(names []string, action string) {
	plists := getPlists()

	for _, name := range names {
		for _, plist := range plists {
			if strings.Index(plist.Name, name) != -1 {
				switch action {
				case "start":
					startDaemon(plist)
				case "stop":
					stopDaemon(plist)
				case "restart":
					stopDaemon(plist)
					startDaemon(plist)
				}
			}
		}
	}
}

func startProfile() {
	fmt.Println("Starting daemons in profile:", profilePath())
	plistsAction(readProfile(), "start")
}

func stopProfile() {
	fmt.Println("Stopping daemons in profile:", profilePath())
	plistsAction(readProfile(), "stop")
}

func restartProfile() {
	fmt.Println("Restarting daemons in profile:", profilePath())
	plistsAction(readProfile(), "restart")
}

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func main() {
	args := os.Args

	if len(args) == 1 {
		printUsage()
		os.Exit(1)
	}

	switch args[1] {
	default:
		printUsage()
		os.Exit(1)
	case "help":
		printUsage()
		return
	case "list", "ls":
		printList(args[2:])
		return
	case "status", "ps":
		printStatus(args[2:])
		return
	case "start":
		startDaemons(args[2:])
		return
	case "stop":
		stopDaemons(args[2:])
		return
	case "restart":
		restartDaemons(args[2:])
		return
	case "show":
		showPlist(args[2:])
		return
	case "edit":
		editPlist(args[2:])
		return
	case "install", "add":
		installPlist(args[2:])
		return
	case "remove", "rm", "uninstall":
		removePlist(args[2:])
		return
	case "scan":
		scanPath(args[2:])
		return
	}
}
