package example

import (
	"fmt"
	"math/rand"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
)

func init() {
	http.HandleFunc("/", handler)
}

func randomString() string {
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	return string(letters[:rand.Intn(len(letters))])
}

func writeToCloudStorage(r *http.Request) error {
	fileName := "file.txt"

	ctx := appengine.NewContext(r)

	// determine default bucket name
	bucketName, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	wc := bucket.Object(fileName).NewWriter(ctx)
	wc.ContentType = "text/plain"

	if _, err := wc.Write([]byte("abcde\n")); err != nil {
		log.Errorf(ctx, "createFile: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
		return err
	}

	if err := wc.Close(); err != nil {
		log.Errorf(ctx, "createFile: unable to close bucket %q, file %q: %v", bucket, fileName, err)
		return err
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello, sHAR!\n")
	fmt.Fprintf(w, randomString())

	writeToCloudStorage(r)

}
