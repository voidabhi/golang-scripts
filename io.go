
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"bytes"
)

type Io struct {
	reader    *bufio.Reader
	writer    *bufio.Writer
	tokens    []string
	nextToken int
}

func NewIo() *Io {
	return &Io{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

func (io *Io) Flush() {
	io.writer.Flush()
}

func (io *Io) NextLine() string {
	var buf []byte
	for {
		line, isPrefix, _ := io.reader.ReadLine()
		buf = append(buf, line...)
		if !isPrefix {
			break
		}
	}
	return string(buf)
}

func (io *Io) Next() string {
	if io.nextToken >= len(io.tokens) {
		line := io.NextLine()
		io.tokens = strings.Fields(line)
		io.nextToken = 0
	}
	r := io.tokens[io.nextToken]
	io.nextToken++
	return r
}

func (io *Io) NextInt() int {
	i, _ := strconv.Atoi(io.Next())
	return i
}

func (io *Io) NextLong() int64 {
	i, _ := strconv.ParseInt(io.Next(), 10, 64)
	return i
}

func (io *Io) NextDouble() float64 {
	i, _ := strconv.ParseFloat(io.Next(), 64)
	return i
}

func (io *Io) Println(a ...interface{}) {
	var buffer bytes.Buffer
	for i := range a {
		if i > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString("%v")
	}
	io.Printfln(buffer.String(), a...)
}

func (io *Io) Printfln(format string, a ...interface{}) {
	fmt.Fprintf(io.writer, format + "\n", a...)
}

func Add(a, b int) int {
	return a + b
}

func main() {
	io := NewIo()
	defer io.Flush()
	a := io.NextInt()
	b := io.Next()
	c := io.NextDouble()
	d := io.NextLine()
	io.Println(a, b, c, d)
	io.Printfln("%.3f", c)
}
