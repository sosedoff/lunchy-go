package main

import(
  "fmt"
  "os"
  "io/ioutil"
  "path/filepath"
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
    if (filepath.Ext(file.Name())) == ".plist" {
      result = append(result, file.Name())
    }
  }

  return result
}

func printList() {
  path  := fmt.Sprintf("%s/Library/LaunchAgents", os.Getenv("HOME"))
  files := findPlists(path)

  for _, file := range files {
    fmt.Println(file)
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
  }
}