// imgconv は自作パッケージです。
package imgconv

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// MyError はユーザー定義型です。
type MyError string

// Error はMyError型のメソッドです。
func (e MyError) Error() string {
	return string(e)
}

var ExitStatus int

// IsPng は.pngファイルかどうかbool値で返す関数です。
func IsPng(path string) bool {
	return filepath.Ext(path) == ".png"
}

// TrimSpaceLeft はエラーメッセージにおいて不要なスペースから左部分を除く関数です。
func TrimSpaceLeft(err error) string {
	str := err.Error()
	spaceIndex := strings.Index(str, " ")
	if spaceIndex == -1 {
		return str
	}
	return str[spaceIndex+1:]
}

// JPGtoPng はJPGファイルから.pngファイルに変換する関数です。
func JPGtoPng(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return MyError(err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return MyError(err.Error())
	}

	var pngFile string
	switch filepath.Ext(path) {
	case ".jpg":
		pngFile = strings.TrimSuffix(path, ".jpg") + ".png"
	case ".jpeg":
		pngFile = strings.TrimSuffix(path, ".jpeg") + ".png"
	}
	out, err := os.Create(pngFile)
	if err != nil {
		return MyError(err.Error())
	}

	err = png.Encode(out, img)
	if err != nil {
		return MyError(err.Error())
	}

	return nil
}

var OsStderr = os.Stderr

// FindJPG はJPGファイルを探す関数です。
func FindJPG(dirname string) {
	err := filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg" {
				err := JPGtoPng(path)
				if err != nil {
					fmt.Fprintf(OsStderr, "error: %s: %s\n", path, err.Error())
					ExitStatus = 1
				}
			} else if info.IsDir() == false && IsPng(path) == false {
				fmt.Fprintf(OsStderr, "error: %s is not a valid file\n", path)
				ExitStatus = 1
			}
			return nil
		})
	if err != nil {
		fmt.Fprintf(OsStderr, "%s\n", err.Error())
		ExitStatus = 1
	}
}

// Convert はmypkgのmainとなる関数です。
func Convert() int {
	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Fprintf(OsStderr, "error: invalid argument\n")
		return 1
	}
	for i := 0; flag.Arg(i) != ""; i++ {
		if _, err := os.Stat(flag.Arg(i)); err != nil {
			fmt.Fprintf(OsStderr, "error: %s\n", TrimSpaceLeft(err))
			ExitStatus = 1
			continue
		}
		FindJPG(flag.Arg(i))
	}
	return ExitStatus
}
