package main

import (
	"flag"
	"fmt"
	"main/internal/apiserver"
	"main/internal/platform/dbmanager"

	"github.com/pterm/pterm"
)

func main() {
	banner()
	portPtr := flag.String("port", "8000", "specify port to run on")
	dbNamePtr := flag.String("name", "presence", "db name")
	flag.Parse()
	dbmanager.Start(*dbNamePtr)

	apiserver.Start(fmt.Sprintf(":%s", *portPtr))
}

// banner prints product banner
func banner() {
	pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Sprint("Presence"))
	pterm.Info.Println("Easy to use attendance manager by Akhil Datla")

}
