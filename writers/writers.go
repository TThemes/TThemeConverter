package writers

import (
	"io"
	"log"
	"os"

	"github.com/TThemes/TThemeConverter/types"
	"github.com/TThemes/TThemeConverter/utils"
	"github.com/mholt/archiver"
)

// WriteThemeFile saves theme file
func WriteThemeFile(fileout *string, thememap map[string]string, isTiledBg bool) {
	osname := utils.GuessOS(fileout)
	var fname string
	writerParams := types.GetThemeParams(osname)
	if osname == "desktop" {
		fname = writerParams.FileName
	} else {
		fname = *fileout
	}
	file, err := os.OpenFile(
		"out/"+fname,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0755,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for key, value := range thememap {
		if writerParams.Argb {
			if len(value) == 8 {
				value = value[6:] + value[:6]
			}
		}
		byteSlice := []byte(key + writerParams.Divider + "#" + value + writerParams.Eols)
		_, err := file.Write(byteSlice)
		if err != nil {
			log.Fatal(err)
		}
	}

	switch osname {
	case "desktop":
		{
			if isTiledBg {
				err := archiver.Zip.Make("out/"+*fileout, []string{"out/colors.tdesktop-theme", "tmp/tiled.jpg"})
				if err != nil {
					log.Fatal(err)
				}
			} else {
				err := archiver.Zip.Make("out/"+*fileout, []string{"out/colors.tdesktop-theme", "tmp/background.jpg"})
				if err != nil {
					log.Fatal(err)
				}
			}
			os.Remove("out/colors.tdesktop-theme")
		}
	case "android":
		{
			if thememap["chat_wallpaper"] == "" {
				fi, err := os.Open("tmp/converted.jpg")
				if err != nil {
					panic(err)
				}
				_, err = file.Write([]byte("WPS\n"))
				if err != nil {
					log.Fatal(err)
				}

				buf := make([]byte, 1024)
				for {
					n, err := fi.Read(buf)
					if err != nil && err != io.EOF {
						log.Fatal(err)
					}
					if n == 0 {
						break
					}
					if _, err := file.Write(buf[:n]); err != nil {
						log.Fatal(err)
					}
				}

				_, err = file.Write([]byte("\nWPE"))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	case "ios":
		{
		}
	}
}
