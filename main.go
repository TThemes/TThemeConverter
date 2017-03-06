package main

import (
	"flag"
	"fmt"

	"github.com/TThemes/TThemeConverter/converter"
	"github.com/TThemes/TThemeConverter/utils"
)

//
//
// func makeAttheme(mapfile, tdesktopmap, transmap, overridemap map[string]string, filename string) {
// 	var color string
// 	atthememap := make(map[string]string)
// 	file, err := os.OpenFile(
// 		"./atthemes/"+filename+".attheme",
// 		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
// 		0666,
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
// 	for key, value := range mapfile {
// 		var desktopvalue string
// 		if overridemap[key] != "" {
// 			desktopvalue = overridemap[key]
// 			delete(overridemap, key)
// 		} else {
// 			desktopvalue = tdesktopmap[value]
// 		}
// 		if len(desktopvalue) == 6 {
// 			if transmap[key] != "" {
// 				color = transmap[key] + desktopvalue
// 			} else {
// 				color = desktopvalue
// 			}
// 		} else if len(desktopvalue) == 8 {
// 			color = desktopvalue[6:] + desktopvalue[:6]
// 		} else {
// 			fmt.Println("Key " + mapfile[key] + " missing from .tdesktop-theme")
// 			color = "00ff00"
// 		}
// 		atthememap[key] = strings.ToUpper(color)
// 	}
// 	if len(overridemap) != 0 {
// 		for key, value := range overridemap {
// 			trans := "ff"
// 			if len(value) == 6 {
// 				if transmap[key] != "" {
// 					color = transmap[key] + value
// 				} else {
// 					color = trans + value
// 				}
// 			} else {
// 				color = value[6:] + value[:6]
// 			}
// 			atthememap[key] = color
// 		}
// 	}
// 	for key, value := range atthememap {
// 		byteSlice := []byte(key + "=#" + value + "\n")
// 		_, err := file.Write(byteSlice)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	if atthememap["chat_wallpaper"] == "" {
// 		fi, err := os.Open("./wip/" + filename + "/converted.jpg")
// 		if err != nil {
// 			panic(err)
// 		}
// 		_, err = file.Write([]byte("WPS\n"))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		buf := make([]byte, 1024)
// 		for {
// 			n, err := fi.Read(buf)
// 			if err != nil && err != io.EOF {
// 				log.Fatal(err)
// 			}
// 			if n == 0 {
// 				break
// 			}
// 			if _, err := file.Write(buf[:n]); err != nil {
// 				log.Fatal(err)
// 			}
// 		}
//
// 		_, err = file.Write([]byte("\nWPE"))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

func main() {
	fileIn := flag.String("in", "", "Use that file as input")
	// fileIn := flag.String("in", "arc_dark.attheme", "Use that file as input")
	fileOut := flag.String("out", "", "Use that file as output")
	// fileOut := flag.String("out", "arc_dark.tdesktop-theme", "Use that file as output")
	isOffline := flag.Bool("offline", false, "Allows you to run without internet connection if you provide it with map files for convertion")
	flag.Parse()
	if *isOffline {
		fmt.Println("Running in offline mode, map files coluld be outdated!")
	}
	utils.PrepareFolders()
	converter.Convert(fileIn, fileOut, isOffline)
}
