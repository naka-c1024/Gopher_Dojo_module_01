package imgconv_test

import (
	"convert/imgconv"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func AssertEquals(t *testing.T, expected bool, actual bool) {
	t.Helper()
	if expected != actual {
		t.Errorf("Unexpected bool\nexpected:%v, actual:%v", expected, actual)
	}
}

func TestIsPng(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    string // 関数に渡すもの
		expected bool   // 期待するもの
	}{
		{input: "../testdata/isPng/foo.png", expected: true},
		{input: "../testdata/isPng/bar.jpg", expected: false},
		{input: "../testdata/isPng/baz.png.jpg", expected: false},
		{input: "../testdata/isPng/hoge.jpg.png", expected: true},
	}
	for _, c := range cases {
		AssertEquals(t, c.expected, imgconv.IsPng(c.input))
	}
}

func TestTrimSpaceLeft(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    error  // 関数に渡すもの
		expected string // 期待するもの
	}{
		{input: errors.New("open aaa: no such file or directory"), expected: "aaa: no such file or directory"},
		{input: errors.New("stat bbb: no such file or directory"), expected: "bbb: no such file or directory"},
		{input: errors.New("image: unknown format"), expected: "unknown format"},
		{input: errors.New("abc"), expected: "abc"},
	}
	for _, c := range cases {
		if actual := imgconv.TrimSpaceLeft(c.input); actual != c.expected {
			t.Errorf("want IsEven(%v) = %s, but actual = %s", c.input, c.expected, actual)
		}
	}
}

func TestCheckError(t *testing.T) {
	// OsExit のバックアップと defer でリカバー
	oldOsExit := imgconv.OsExit
	defer func() { imgconv.OsExit = oldOsExit }()
	// あとで OsExit 内で終了ステータスをキャプチャするための変数
	var capture int
	imgconv.OsExit = func(code int) { capture = code }

	oldOsStderr := imgconv.OsStderr
	defer func() { imgconv.OsStderr = oldOsStderr }()
	imgconv.OsStderr = nil

	cases := []struct {
		input    error // 関数に渡すもの
		expected int
	}{
		{input: nil, expected: 0},
		{input: errors.New("open aaa: no such file or directory"), expected: 1},
	}

	for _, c := range cases {
		imgconv.CheckError(c.input, imgconv.ErrMsg("foo"), "bar")
		actual := capture
		if actual != c.expected {
			t.Errorf("Fail assert equal. Expect: %v Actual: %v", c.expected, actual)
		}
	}
}

func TestJpgToPng(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string // 関数に渡すもの
	}{
		{input: "../testdata/Go-BB_cover.jpg"},
		{input: "../testdata/Go-BB_spread1.jpg"},
		{input: "../testdata/Go-BB_spread2.jpg"},
		{input: "../testdata/logos.jpg"},
	}
	for _, c := range cases {
		imgconv.JpgToPng(c.input)
		png_file := strings.Replace(c.input, "jpg", "png", -1)
		if _, err := os.Stat(png_file); err != nil {
			t.Errorf("%v do not create", png_file)
		}
		exec.Command("rm", png_file).Run()
	}
}

func TestDirExists(t *testing.T) {
	// OsExit のバックアップと defer でリカバー
	oldOsExit := imgconv.OsExit
	defer func() { imgconv.OsExit = oldOsExit }()
	// あとで OsExit 内で終了ステータスをキャプチャするための変数
	var capture int
	imgconv.OsExit = func(code int) { capture = code }

	oldOsStderr := imgconv.OsStderr
	defer func() { imgconv.OsStderr = oldOsStderr }()
	imgconv.OsStderr = nil

	cases := []struct {
		input    string
		expected int
	}{
		{input: "../testdata", expected: 0},
		{input: "hoge", expected: 1},
	}

	for _, c := range cases {
		imgconv.DirExists(c.input)
		actual := capture
		if actual != c.expected {
			t.Errorf("Fail assert equal. Expect: %v Actual: %v", c.expected, actual)
		}
	}
}

func TestWalkMainFunc(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string // 関数に渡すもの
	}{
		{input: "../testdata/Go-BB_cover.jpg"},
		{input: "../testdata/Go-BB_spread1.jpg"},
		{input: "../testdata/Go-BB_spread2.jpg"},
		{input: "../testdata/logos.jpg"},
	}
	for _, c := range cases {
		imgconv.WalkMainFunc(c.input, nil)
		png_file := strings.Replace(c.input, "jpg", "png", -1)
		if _, err := os.Stat(png_file); err != nil {
			t.Errorf("%v do not create", png_file)
		}
		exec.Command("rm", png_file).Run()
	}
}

func TestConvertMain(t *testing.T) {
	// OsExit のバックアップと defer でリカバー
	oldOsExit := imgconv.OsExit
	defer func() { imgconv.OsExit = oldOsExit }()
	// あとで OsExit 内で終了ステータスをキャプチャするための変数
	var capture int
	imgconv.OsExit = func(code int) { capture = code }

	oldOsStderr := imgconv.OsStderr
	defer func() { imgconv.OsStderr = oldOsStderr }()
	imgconv.OsStderr = nil

	cases := []struct {
		input    string
		expected int
	}{
		{input: "hogehoge", expected: 1},
	}

	for _, c := range cases {
		imgconv.ConvertMain(c.input)
		actual := capture
		if actual != c.expected {
			t.Errorf("Fail assert equal. Expect: %v Actual: %v", c.expected, actual)
		}
	}
}
