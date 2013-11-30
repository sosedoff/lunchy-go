package main

import(
  "fmt"
  "os"
  "io/ioutil"
  "path/filepath"
  "os/exec"
  "strings"
)

const(
  LUNCHY_VERSION = "0.1.0"
)

func printUsage() {
  fmt.Printf("Lunchy %s, the friendly launchctl wrapper\n", LUNCHY_VERSION)
  fmt.Println("Usage: lunchy [start|stop|restart|list|status|install|show|edit] [options]")
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

func sliceIncludes(slice []string, match string) bool {
  for _, val := range slice {
    if val == match {
      return true
    }
  }

  return false
}

func printList() {
  for _, file := range getPlists() {
    fmt.Println(file)
  }
}

func printStatus(args []string) {
  out, err := exec.Command("launchctl", "list").Output()

  if err != nil {
    fmt.Println("Failed to execute", err)
    os.Exit(1)
  }

  pattern := ""

  if len(args) == 3 {
    pattern = args[2]
  }

  installed := getPlists()
  lines := strings.Split(strings.TrimSpace(string(out)), "\n")
  
  for _, line := range lines {
    chunks := strings.Split(line, "\t")

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
  case "list":
    printList()
    return
  case "status":
    printStatus(args)
    return
  }
}