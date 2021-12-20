package main

import (
	"flag"
	"fmt"
	"main/internal/apiserver"
	"main/internal/platform/dbmanager"
)

func main() {
	banner()
	portPtr := flag.Int("port", 8000, "specify port to run on")
	dbNamePtr := flag.String("name", "presence", "db name")
	flag.Parse()
	dbmanager.Start(*dbNamePtr)
	apiserver.Start(fmt.Sprintf(":%d", *portPtr))
}

// banner prints product banner
func banner() {
	fmt.Print(`
Presence - An easy-to-use attendance manager
(c)2021 - Akhil Datla

`)

}
