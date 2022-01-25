package primitiveUtil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Mode int

const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

func Tranform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	in, err := tempFile("in_", ext)
	if err != nil {
		return nil, errors.New("primitiveUtils: Failed to create temp input file")
	}
	defer in.Close()
	defer os.Remove(in.Name())
	out, err := tempFile("out_", ext)
	if err != nil {
		return nil, errors.New("primitiveUtils: Failed to create temp output file")
	}
	fmt.Println(out.Name())
	defer out.Close()
	defer os.Remove(out.Name())

	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("primitiveUtils: Failed to copyimgae into input file")
	}
	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, ModeCombo)
	if err != nil {
		return nil, fmt.Errorf("primitiveUtils: Failed to run primitive command, stdCombo = %s", stdCombo)
	}
	fmt.Println(stdCombo)
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, errors.New("primitiveUtils: Failed to copy outfut file to buffer")
	}

	return b, nil
}

func primitive(inputFile,
	outputFile string,
	numShapes int, mode Mode) (string, error) {
	argsString := fmt.Sprintf("-i %s -o %s -n %d -m %d ", inputFile, outputFile, numShapes, mode)
	cmd := exec.Command("primitive", strings.Fields(argsString)...)
	// doesnt gice reat debuuign change
	b, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return string(b), err
}

func tempFile(prefix, ext string) (*os.File, error) {
	// wtf is the need fo temp?
	in, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, errors.New("primitiveUtils: Failed to create temp input file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
