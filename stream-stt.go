package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	speech "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
)

func main() {
	ctx := context.Background()
	conn, err := transport.DialGRPC(ctx,
		option.WithEndpoint("speech.googleapis.com:443"),
		option.WithScopes("https://www.googleapis.com/auth/cloud-platform"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	stream, err := speech.NewSpeechClient(conn).StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// send the initial configuration message.
	if err := stream.Send(&speech.StreamingRecognizeRequest{
		StreamingRequest: &speech.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speech.StreamingRecognitionConfig{
				Config: &speech.RecognitionConfig{
					Encoding:   speech.RecognitionConfig_LINEAR16,
					SampleRate: 16000,
				},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		// pipe stdin to the API
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err == io.EOF {
				return // nothing else to pipe, kill this goroutine
			}
			if err != nil {
				log.Printf("reading stdin error: %v", err)
				continue
			}
			if err = stream.Send(&speech.StreamingRecognizeRequest{
				StreamingRequest: &speech.StreamingRecognizeRequest_AudioContent{
					AudioContent: buf[:n],
				},
			}); err != nil {
				log.Printf("sending audio error: %v", err)
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// TODO: handle error
			continue
		}
		if resp.Error != nil {
			// TODO: handle error
			continue
		}
		for _, result := range resp.Results {
			fmt.Printf("result: %+v\n", result)
		}
	}
}
