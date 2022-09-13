// imgconv は自作パッケージです。
package imgconv

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type ErrMsg string

func IsPng(path string) bool {
	if ext := filepath.Ext(path); ext == ".png" {
		return true
	} else {
		return false
	}
}

func TrimSpaceLeft(err error) string {
	str := err.Error()
	spaceIndex := strings.Index(str, " ")
	if spaceIndex == -1 {
		return str
	}
	return str[spaceIndex+1:]
}

var OsExit = os.Exit
var OsStderr = os.Stderr

func CheckError(err error, msg ErrMsg, path string) {
	if err != nil {
		fmt.Fprintf(OsStderr, "error: %s: %v: %v\n", path, msg, TrimSpaceLeft(err))
		OsExit(1)
	}
}

func JpgToPng(path string) {
	file, err := os.Open(path)
	CheckError(err, "open", path)
	defer file.Close()

	img, _, err := image.Decode(file)
	CheckError(err, "decode", path)

	png_file := strings.Replace(path, "jpg", "png", -1)
	out, err := os.Create(png_file)
	CheckError(err, "create", path)
	defer out.Close()

	png.Encode(out, img)
}

func DirExists(dirname string) {
	if _, err := os.Stat(dirname); err != nil {
		fmt.Fprintf(OsStderr, "error: %s\n", TrimSpaceLeft(err))
		OsExit(1)
	}
}

func WalkMainFunc(path string, info os.FileInfo) {
	if filepath.Ext(path) == ".jpg" {
		JpgToPng(path)
	} else if info.IsDir() == false && IsPng(path) == false {
		fmt.Fprintf(os.Stderr, "error: %s is not a valid file\n", path)
	}
}

func ConvertMain(dirname string) {
	DirExists(dirname)
	err := filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			WalkMainFunc(path, info)
			return nil
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", TrimSpaceLeft(err))
		OsExit(1)
	}
}
