package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// MakeTiledBg creates 1x1 converted.jpg
func MakeTiledBg(bgcolor string) {
	back := image.NewRGBA(image.Rect(0, 0, 1, 1))
	alpha, _ := strconv.ParseUint("ff", 16, 32)
	red, _ := strconv.ParseUint(bgcolor[0:2], 16, 32)
	green, _ := strconv.ParseUint(bgcolor[2:4], 16, 32)
	blue, _ := strconv.ParseUint(bgcolor[4:6], 16, 32)
	if len(bgcolor) == 8 {
		alpha, _ = strconv.ParseUint(bgcolor[6:8], 16, 32)
	}

	back.SetRGBA(0, 0, color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(alpha)})

	toimg, _ := os.Create("tmp/tiled.jpg")
	defer toimg.Close()
	err := jpeg.Encode(toimg, back, &jpeg.Options{Quality: 87})
	if err != nil {
		log.Fatal(err)
	}

}

// FileExists is checks if file exist
// https://stackoverflow.com/a/12527546
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// ConvertBg converts background from png / jpeg to jpeg in case it's tiled fills
// 1920x1920 square with tiles
func ConvertBg() {
	var fname, fext string
	var isTiled, isPng bool
	var img image.Image
	filenames := []string{"tiled", "background"}
	fileexts := []string{".jpg", ".png"}
	for _, filename := range filenames {
		for _, fileext := range fileexts {
			if FileExists("tmp/" + filename + fileext) {
				fname = filename
				fext = fileext
			}
		}
	}
	if fname == "tiled" {
		isTiled = true
	}
	if fext == ".png" {
		isPng = true
	}

	if isPng {
		fimg, _ := os.Open("tmp/" + fname + fext)
		defer fimg.Close()
		img, _ = png.Decode(fimg)
	} else {
		fimg, _ := os.Open("tmp/" + fname + fext)
		defer fimg.Close()
		img, _ = jpeg.Decode(fimg)
	}

	imgsizex := 1920
	imgsizey := 1920

	back := image.NewRGBA(image.Rect(0, 0, imgsizex, imgsizey))

	if isTiled {
		tilesizex := img.Bounds().Max.X
		tilesizey := img.Bounds().Max.Y
		stepx := imgsizex / tilesizex
		stepy := imgsizey / tilesizey
		for x := 0; x <= stepx; x++ {
			for y := 0; y <= stepy; y++ {
				draw.Draw(back, back.Bounds(), img, image.Point{-x * tilesizex, -y * tilesizey}, draw.Src)
			}
		}
	}

	toimg, _ := os.Create("tmp/converted.jpg")
	defer toimg.Close()
	if isTiled {
		err := jpeg.Encode(toimg, back, &jpeg.Options{Quality: 87})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := jpeg.Encode(toimg, img, &jpeg.Options{Quality: 87})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// PrepareFolders creates folders for convertor
func PrepareFolders() {
	dirs := []string{"in", "out", "map", "tmp"}
	for _, dir := range dirs {
		if !FileExists(dir) {
			os.Mkdir(dir, 0777)
		}
	}
}

// GetLatestMaps downloads latest maps for conversion from repos_url
// https://github.com/TThemes/TThemeMap/
func GetLatestMaps(osin, osout string) {
	baseurl := "https://raw.githubusercontent.com/TThemes/TThemeMap/master/"
	downloadFromURL(baseurl + osin + "_" + osout + ".map")
	downloadFromURL(baseurl + osin + "_" + osout + "_trans.map")
}

func downloadFromURL(url string) {
	tokens := strings.Split(url, "/")
	fileName := "map/" + tokens[len(tokens)-1]

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
}

// GuessOS returns OS depending on theme file extencion
func GuessOS(filename *string) string {
	filenameParts := strings.Split(*filename, ".")
	fileExt := filenameParts[len(filenameParts)-1]
	switch fileExt {
	case "tdesktop-theme":
		{
			return "desktop"
		}
	case "attheme":
		{
			return "android"
		}
	case "iostheme":
		{
			return "ios"
		}
	default:
		{
			return "other"
		}
	}
}
