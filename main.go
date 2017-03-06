package main

import (
	"flag"
	"fmt"

	"github.com/TThemes/TThemeConverter/converter"
	"github.com/TThemes/TThemeConverter/utils"
)

func main() {
	fileIn := flag.String("in", "", "Use that file as input")
	// fileIn := flag.String("in", "arc_dark.attheme", "Use that file as input")
	fileOut := flag.String("out", "", "Use that file as output")
	// fileOut := flag.String("out", "arc_dark.tdesktop-theme", "Use that file as output")
	isOffline := flag.Bool("offline", false, "Allows you to run without internet connection if you provide it with map files for conversion")
	flag.Parse()
	if *isOffline {
		fmt.Println("Running in offline mode, map files coluld be outdated!")
	}
	utils.PrepareFolders()
	converter.Convert(fileIn, fileOut, isOffline)
}
