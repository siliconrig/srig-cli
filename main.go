package main

import (
	"fmt"
	"os"
)

const usage = `fbay — Flashbay CLI

Usage:
  fbay <command> [options]

Commands:
  flash    Flash firmware to a remote board
  serial   Open serial terminal to a remote board
  session  Manage sessions (create, end)
  status   Show board availability and session status

Use "fbay <command> --help" for more information.`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "flash":
		cmdFlash(os.Args[2:])
	case "serial":
		cmdSerial(os.Args[2:])
	case "session":
		cmdSession(os.Args[2:])
	case "status":
		cmdStatus(os.Args[2:])
	case "--help", "-h", "help":
		fmt.Println(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		fmt.Println(usage)
		os.Exit(1)
	}
}
