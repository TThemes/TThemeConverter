package converter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TThemes/TThemeConverter/readers"
	"github.com/TThemes/TThemeConverter/utils"
	"github.com/TThemes/TThemeConverter/writers"
)

func readOverridesMapFile(filename string) map[string]string {
	overridemap := make(map[string]string)
	fullname := ""
	for _, name := range strings.Split(filename, ".") {
		fullname += name + "."
		fmt.Println("Looking for overrides at " + fullname + "map")
		for key, value := range readMapFile(fullname + "map") {
			overridemap[key] = value
		}
	}
	return overridemap
}

func readMapFile(filename string) map[string]string {
	filemap := make(map[string]string)
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return filemap
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rule := strings.Split(scanner.Text(), "=")
		key := rule[0]
		value := rule[1]
		if filemap[value] != "" {
			value = filemap[value]
		}
		filemap[key] = value
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return filemap
	}

	return filemap
}

func readConvMapFile(osin, osout, filein string) map[string]string {
	defaultMap := readMapFile("map/" + osin + "_" + osout + ".map")
	return defaultMap
}

func applyConvMap(thememap, convmap, transmap map[string]string, filein string) map[string]string {
	fname := strings.TrimSuffix(filein, filepath.Ext(filein))
	outMap := make(map[string]string)
	for key, value := range convmap {
		outMap[key] = thememap[value]
	}
	overrideMap := readOverridesMapFile("in/" + fname)
	for key, value := range overrideMap {
		outMap[key] = value
	}
	for key := range transmap {
		if len(outMap[key]) == 6 {
			outMap[key] += transmap[key]
		}
	}
	return outMap
}

// Convert converts filein to fileout using mapfile
func Convert(filein *string, fileout *string, isoffline *bool) {
	var isTiledBg bool
	osIn := utils.GuessOS(filein)
	osOut := utils.GuessOS(fileout)
	fmt.Println("Converting " + osIn + " theme to " + osOut + "...")
	if osIn != osOut {
		themeMap := readers.ReadThemeFile(filein)
		if !*isoffline {
			utils.GetLatestMaps(osIn, osOut)
		}
		if themeMap["chat_wallpaper"] != "" {
			utils.MakeTiledBg(themeMap["chat_wallpaper"])
			isTiledBg = true
		}
		transMap := readMapFile("map/" + osIn + "_" + osOut + "_trans.map")
		convMap := readConvMapFile(osIn, osOut, *filein)
		outMap := applyConvMap(themeMap, convMap, transMap, *filein)
		writers.WriteThemeFile(fileout, outMap, isTiledBg)
		os.RemoveAll("tmp/")
	} else {
		fmt.Println("You are trying convert to same OS!")
		os.RemoveAll("tmp/")
	}
}
