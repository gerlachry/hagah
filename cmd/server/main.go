// The main entry point for the server.
// Serves APIs for various tasks.
package main

import (
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
}

func main() {
	log.Print("Lets serve some data...")
}
