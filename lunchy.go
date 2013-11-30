package main

import(
  "fmt"
  "os"
)

const(
  LUNCHY_VERSION = "0.1.0"
)

func printUsage() {
  fmt.Printf("Lunchy %s, the friendly launchctl wrapper\n", LUNCHY_VERSION)
  fmt.Println("Usage: lunchy [start|stop|restart|list|status|install|show|edit] [options]")
}

func main() {
  args := os.Args

  if (len(args) == 1) {
    printUsage()
    os.Exit(1)
  }
}