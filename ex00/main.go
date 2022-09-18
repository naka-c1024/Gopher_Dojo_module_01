package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var exitStatus int

func readWrite(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		io.WriteString(w, scanner.Text()+"\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func trimSpaceLeft(err error) string {
	str := err.Error()
	spaceIndex := strings.Index(str, " ")
	if spaceIndex == -1 {
		return str
	}
	arr := strings.SplitAfter(str[spaceIndex+1:], " ")
	arr[1] = strings.Title(arr[1])
	rtn := strings.Join(arr, "")
	return rtn
}

func openFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	err = readWrite(f, os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func readFile(argc int) {
	for i := 1; i < argc; i++ {
		if os.Args[i] == "-" {
			readStdin()
			continue
		}
		err := openFile(os.Args[i])
		if err != nil {
			os.Stderr.WriteString("ft_cat: " + trimSpaceLeft(err) + "\n")
			exitStatus = 1
		}
	}
}

func readStdin() {
	err := readWrite(os.Stdin, os.Stdout)
	if err != nil {
		os.Stderr.WriteString("ft_cat: " + trimSpaceLeft(err) + "\n")
		os.Exit(1)
	}
}

func main() {
	argc := len(os.Args)
	if argc == 1 {
		readStdin()
	} else {
		readFile(argc)
	}
	os.Exit(exitStatus)
}
