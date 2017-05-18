
package main

import (
	"bytes"
	"code.google.com/p/portaudio-go/portaudio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
        if len(os.Args) <= 1 {
                fmt.Println("Usage:\n \t./ffmpeg audio_file")
                os.Exit(0)
        }
        inputfile := os.Args[1]
	chk := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	framePerBuffer := 2048
	ff := newFfmpeg(inputfile)
	defer ff.Close()
	stream, err := portaudio.OpenDefaultStream(0, 2, 44100, framePerBuffer, ff)
	chk(err)
	defer stream.Close()
	chk(stream.Start())
	if err := ff.cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	chk(stream.Stop())
}

type ffmpeg struct {
	in  io.ReadCloser
	cmd *exec.Cmd
}

func newFfmpeg(filename string) *ffmpeg {
	cmd := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return &ffmpeg{stdout, cmd}
}

func (f *ffmpeg) Close() error {
	return f.in.Close()
}

func (f *ffmpeg) ProcessAudio(_, out [][]int16) {
	// int16 takes 2 bytes
	bufferSize := len(out[0]) * 4
	var pack = make([]byte, bufferSize)
	if _, err := f.in.Read(pack); err != nil {
		log.Fatal(err)
	}
	n := make([]int16, len(out[0])*2)
	for i := range n {
		var x int16
		buf := bytes.NewBuffer(pack[2*i : 2*(i+1)])
		binary.Read(buf, binary.LittleEndian, &x)
		n[i] = x
	}

	for i := range out[0] {
		out[0][i] = n[2*i]
		out[1][i] = n[2*i+1]
	}
}
