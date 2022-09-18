package imgconv_test

import (
	"convert/imgconv"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMyError(t *testing.T) {
	t.Parallel()
	expect := "hoge"
	actual := imgconv.MyError(expect).Error()
	if actual != expect {
		t.Errorf(`expect="%s" actual="%s"`, expect, actual)
	}
}

func AssertEquals(t *testing.T, expected bool, actual bool) {
	t.Helper()
	if expected != actual {
		t.Errorf("Unexpected bool\nexpected:%v, actual:%v", expected, actual)
	}
}

func TestIsPng(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    string
		expected bool
	}{
		{input: "../testdata/isPng/bar.jpg", expected: false},
		{input: "../testdata/isPng/baz.png.jpg", expected: false},
	}
	for _, c := range cases {
		AssertEquals(t, c.expected, imgconv.IsPng(c.input))
	}
}

func TestTrimSpaceLeft(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    error
		expected string
	}{
		{input: errors.New("open aaa: no such file or directory"), expected: "aaa: no such file or directory"},
		{input: errors.New("stat bbb: no such file or directory"), expected: "bbb: no such file or directory"},
		{input: errors.New("image: unknown format"), expected: "unknown format"},
		{input: errors.New("abc"), expected: "abc"},
	}
	for _, c := range cases {
		if actual := imgconv.TrimSpaceLeft(c.input); actual != c.expected {
			t.Errorf("want TrimSpaceLeft(%v) = %s, but actual = %s", c.input, c.expected, actual)
		}
	}
}

func errorJPGtoPng(t *testing.T) {
	cases := []struct {
		input    string
		expected error
	}{
		{input: "open_hoge", expected: errors.New("open open_hoge: no such file or directory")},
		{input: "../testdata/JPGtoPng/decode.txt", expected: errors.New("image: unknown format")},
	}
	for _, c := range cases {
		if actual := imgconv.JPGtoPng(c.input); actual.Error() != c.expected.Error() {
			t.Errorf("want JPGtoPng(%s) = %s, but actual = %s", c.input, c.expected, actual)
		}
	}
}

func TestJPGtoPng(t *testing.T) {
	// t.Parallel() // kore
	cases := []struct {
		input string
	}{
		{input: "../testdata/Go-BB_cover.jpg"},
		{input: "../testdata/Go-BB_spread1.jpg"},
		{input: "../testdata/Go-BB_spread2.jpeg"},
		{input: "../testdata/logos.jpg"},
	}
	for _, c := range cases {
		imgconv.JPGtoPng(c.input)
		var pngFile string
		switch filepath.Ext(c.input) {
		case ".jpg":
			pngFile = strings.TrimSuffix(c.input, ".jpg") + ".png"
		case ".jpeg":
			pngFile = strings.TrimSuffix(c.input, ".jpeg") + ".png"
		}
		if _, err := os.Stat(pngFile); err != nil {
			t.Errorf("%v do not create", pngFile)
		}
		exec.Command("rm", pngFile).Run()
	}
	errorJPGtoPng(t)
}

func errorFindJPG(t *testing.T) {
	oldOsStderr := imgconv.OsStderr
	defer func() { imgconv.OsStderr = oldOsStderr }()
	imgconv.OsStderr = nil

	cases := []struct {
		input    string
		expected int
	}{
		{input: "open_hoge", expected: 1},
		{input: "../testdata/JPGtoPng/decode.txt", expected: 1},
		{input: "../testdata/FindJPG/aaa.jpeg", expected: 1},
	}
	for _, c := range cases {
		imgconv.FindJPG(c.input)
		if imgconv.ExitStatus != c.expected {
			t.Errorf("want FindJPG(%s) = %d, but actual = %d", c.input, c.expected, imgconv.ExitStatus)
		}
		imgconv.ExitStatus = 0
	}
}

func TestFindJPG(t *testing.T) {
	// t.Parallel() // kore
	cases := []struct {
		input string
	}{
		{input: "../testdata/Go-BB_cover.jpg"},
		{input: "../testdata/Go-BB_spread1.jpg"},
		{input: "../testdata/Go-BB_spread2.jpeg"},
		{input: "../testdata/logos.jpg"},
	}
	for _, c := range cases {
		imgconv.FindJPG(c.input)
		var pngFile string
		switch filepath.Ext(c.input) {
		case ".jpg":
			pngFile = strings.TrimSuffix(c.input, ".jpg") + ".png"
		case ".jpeg":
			pngFile = strings.TrimSuffix(c.input, ".jpeg") + ".png"
		}
		if _, err := os.Stat(pngFile); err != nil {
			t.Errorf("%v do not create", pngFile)
		}
		exec.Command("rm", pngFile).Run()
	}
	errorFindJPG(t)
}

func TestFlag(t *testing.T) {
	t.Parallel()
	oldOsStderr := imgconv.OsStderr
	defer func() { imgconv.OsStderr = oldOsStderr }()
	imgconv.OsStderr = nil

	cases := []struct {
		input    []string
		expected int
	}{
		{input: nil, expected: 1},
		{input: []string{"hoge"}, expected: 1},
		{input: []string{"foo", "bar"}, expected: 1},
	}
	for _, c := range cases {
		actual := imgconv.Flag(c.input)
		if actual != c.expected {
			t.Errorf("Fail assert equal. Expect: %d Actual: %d", 1, actual)
		}
	}
}
