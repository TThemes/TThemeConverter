package readers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/TThemes/TThemeConverter/types"
	"github.com/TThemes/TThemeConverter/utils"
	"github.com/mholt/archiver"
)

func prepareDesktopFile(filename *string) {
	err := archiver.Zip.Open("in/"+*filename, "tmp/")
	if err != nil {
		log.Fatal(err)
	}
	utils.ConvertBg()
}

func prepareAndroidFile(filename *string) {
	f, err := os.Open("in/" + *filename)
	if err != nil {
		log.Fatal(err)
	}

	fo, err := os.OpenFile(
		"tmp/background.jpg",
		os.O_WRONLY|os.O_CREATE,
		0755,
	)
	if err != nil {
		log.Panic(err)
	}

	fto, err := os.OpenFile(
		"tmp/theme.attheme",
		os.O_WRONLY|os.O_CREATE,
		0755,
	)
	if err != nil {
		log.Panic(err)
	}

	var fpos int64
	var fpose int64
	var fposend int64

	fposend, _ = f.Seek(0, 2)
	fpos = fposend

	for i := 0; i < int(fposend)-2; i++ {
		_, err = f.Seek(int64(i), 0)
		if err != nil {
			log.Panic(err)
		}
		b := make([]byte, 3)
		_, err = io.ReadAtLeast(f, b, 3)
		if err != nil {
			log.Panic(err)
		}
		if string(b) == "WPS" {
			fpos = int64(i)
		}
		if string(b) == "WPE" {
			fpose = int64(i)
		}
	}

	if fposend != fpos {
		f.Seek(fpos+4, 0)                    // skip "WPS\n"
		imgbuf := make([]byte, fpose-5-fpos) // get rid of "\nWPE\n"
		_, err = io.ReadAtLeast(f, imgbuf, int(fpose-5-fpos))
		if err != nil {
			log.Panic(err)
		}
		fo.Write(imgbuf)
	}

	f.Seek(0, 0)
	themebuf := make([]byte, fpos)
	_, err = io.ReadAtLeast(f, themebuf, int(fpos))
	if err != nil {
		log.Panic(err)
	}
	fto.Write(themebuf)
}

func prepareIOSFile(filename *string) {
	fmt.Println(*filename)
	fmt.Println("Looks like we cant read ios themes yet")
}

func themeReader(osname string) map[string]string {
	themeMap := make(map[string]string)
	readerParams := types.GetThemeParams(osname)
	file, err := os.Open("tmp/" + readerParams.FileName)
	if err != nil {
		log.Fatal(err)
		return themeMap
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		themeline := scanner.Text()
		rule := strings.Split(themeline, readerParams.Divider)
		if len(rule) >= 2 {
			rulekey := rule[0]
			rulevalue := strings.Split(rule[1], readerParams.Eol)[0]
			switch osname {
			case "desktop":
				{
					if !strings.HasPrefix(rulevalue, "#") {
						rulevalue = themeMap[rulevalue]
					}
					rulevalue = strings.TrimPrefix(rulevalue, "#")
				}
			case "android":
				{
					inval, err := strconv.Atoi(rulevalue)
					if err != nil {
						rulevalue = strings.TrimPrefix(rulevalue, "#")
						if len(rulevalue) == 8 {
							rulevalue = rulevalue[2:8] + rulevalue[0:2]
						}
					} else {
						tohex := ((inval << 8) & 0xffffff00) | ((inval >> 24) & 0xff)
						rulevalue = fmt.Sprintf("%08x", tohex)
					}
				}
			}
			themeMap[rulekey] = rulevalue
		}
	}
	themeMap["whatever"] = "ff00ff"
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return themeMap
	}
	return themeMap
}

// ReadThemeFile is used for reading theme file to map
func ReadThemeFile(filename *string) map[string]string {
	var osname string
	switch osname = utils.GuessOS(filename); osname {
	case "android":
		{
			prepareAndroidFile(filename)
		}
	case "desktop":
		{
			prepareDesktopFile(filename)
		}
	case "ios":
		{
			prepareIOSFile(filename)
		}
	}
	themeMap := themeReader(osname)
	return themeMap
}
