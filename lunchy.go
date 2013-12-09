package main

import(
  "fmt"
  "os"
  "io"
  "io/ioutil"
  "path/filepath"
  "os/exec"
  "strings"
)

const(
  LUNCHY_VERSION = "0.1.2"
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

func findPlists(path string) []string {
  result := []string{}
  files, err := ioutil.ReadDir(path)

  if err != nil {
    return result
  }

  for _, file := range files {
    if !file.IsDir() {
      if (filepath.Ext(file.Name())) == ".plist" {
        name := strings.Replace(file.Name(), ".plist", "", -1)
        result = append(result, name)
      }
    }
  }

  return result
}

func getPlists() []string {
  path := fmt.Sprintf("%s/Library/LaunchAgents", os.Getenv("HOME")) 
  files := findPlists(path)

  return files
}

func getPlist(name string) string {
  for _, plist := range getPlists() {
    if strings.Index(plist, name) != -1 {
      return plist
    }
  }

  return ""
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
  fmt.Println("Usage: lunchy [start|stop|restart|list|status|install|show|edit|remove] [options]")
}

func printList() {
  for _, file := range getPlists() {
    fmt.Println(file)
  }
}

func printStatus(args []string) {
  out, err := exec.Command("launchctl", "list").Output()

  if err != nil {
    fatal("failed to get process list")
  }

  pattern := ""

  if len(args) == 3 {
    pattern = args[2]
  }

  installed := getPlists()
  lines := strings.Split(strings.TrimSpace(string(out)), "\n")

  for _, line := range lines {
    chunks := strings.Split(line, "\t")
    clean_line := strings.Replace(line, "\t", " ", -1)

    if len(pattern) > 0 {
      if strings.Index(chunks[2], pattern) != -1 {
        if sliceIncludes(installed, chunks[2]) {
          fmt.Println(clean_line)
        }
      }
    } else {
      if sliceIncludes(installed, chunks[2]) {
        fmt.Println(clean_line)
      }
    }
  }
}

func exitWithInvalidArgs(args []string, msg string) {
  if len(args) < 3 {
    fmt.Println(msg)
    os.Exit(1)
  }
}

func startDaemons(args []string) {
  exitWithInvalidArgs(args, "name required")

  name := args[2]

  for _, plist := range getPlists() {
    if strings.Index(plist, name) != -1 {
      startDaemon(plist)
    }
  }
}

func startDaemon(name string) {
  path := fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), name)
  _, err := exec.Command("launchctl", "load", path).Output()

  if err != nil {
    fmt.Println("failed to start", name)
    return
  }

  fmt.Println("started", name)
}

func stopDaemons(args []string) {
  exitWithInvalidArgs(args, "name required")

  name := args[2]

  for _, plist := range getPlists() {
    if strings.Index(plist, name) != -1 {
      stopDaemon(plist)
    }
  }
}

func stopDaemon(name string) {
  path := fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), name)
  _, err := exec.Command("launchctl", "unload", path).Output()

  if err != nil {
    fmt.Println("failed to stop", name)
    return
  }

  fmt.Println("stopped", name)
}

func restartDaemons(args []string) {
  exitWithInvalidArgs(args, "name required")

  name := args[2]

  for _, plist := range getPlists() {
    if strings.Index(plist, name) != -1 {
      stopDaemon(plist)
      startDaemon(plist)
    }
  }
}

func showPlist(args []string) {
  exitWithInvalidArgs(args, "name required")

  name := args[2]

  for _, plist := range getPlists() {
    if strings.Index(plist, name) != -1 {
      printPlistContent(plist)
      return
    }
  }
}

func printPlistContent(name string) {
  path := fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), name)
  contents, err := ioutil.ReadFile(path)

  if err != nil {
    fatal("unable to read plist")
  }

  fmt.Printf(string(contents))
}

func editPlist(args []string) {
  exitWithInvalidArgs(args, "name required")

  name := args[2]

  for _, plist := range getPlists() {
    if strings.Index(plist, name) != -1 {
      editPlistContent(plist)
      return
    }
  }
}

func editPlistContent(name string) {
  path := fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), name)
  editor := os.Getenv("EDITOR")

  if len(editor) == 0 {
    fatal("EDITOR environment variable is not set")
  }

  cmd := exec.Command(editor, path)
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout

  cmd.Start()
  cmd.Wait()
}

func installPlist(args []string) {
  exitWithInvalidArgs(args, "path required")

  path := args[2]

  if !fileExists(path) {
    fatal("source file does not exist")
  }

  info, _ := os.Stat(path)
  base_path := fmt.Sprintf("%s/%s", os.Getenv("HOME"), "Library/LaunchAgents")
  new_path := fmt.Sprintf("%s/%s", base_path, info.Name())

  if fileExists(new_path) && os.Remove(new_path) != nil {
    fatal("unable to delete existing plist")
  }

  if fileCopy(path, new_path) != nil {
    fatal("failed to copy file")
  }

  fmt.Println(path, "installed to", base_path)
}

func removePlist(args []string) {
  exitWithInvalidArgs(args, "name required")

  name := args[2]
  base_path := fmt.Sprintf("%s/%s", os.Getenv("HOME"), "Library/LaunchAgents")

  for _, plist := range getPlists() {
    if strings.Index(plist, name) != -1 {
      path := fmt.Sprintf("%s/%s.plist", base_path, plist)

      if os.Remove(path) == nil {
        fmt.Println("removed", path)
      } else {
        fmt.Println("failed to remove", path)
      }
    }
  }
}

func fatal(message string) {
  fmt.Println(message)
  os.Exit(1)
}

func main() {
  args := os.Args

  if (len(args) == 1) {
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
    printList()
    return
  case "status":
    printStatus(args)
    return
  case "start":
    startDaemons(args)
    return
  case "stop":
    stopDaemons(args)
    return
  case "restart":
    restartDaemons(args)
    return
  case "show":
    showPlist(args)
    return
  case "edit":
    editPlist(args)
    return
  case "install":
    installPlist(args)
    return
  case "remove":
    removePlist(args)
    return
  }
}
