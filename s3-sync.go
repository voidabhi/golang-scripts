package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"os"
	"path/filepath"
)

const debug = true

var (
	bucketName string
	localPath  string
)

func init() {
	flag.StringVar(&bucketName, "b", "", "Bucket Name")
	flag.StringVar(&localPath, "p", "", "Local Path")
}

func HashFile(reader io.Reader) []byte {
	hasher := md5.New()
	io.Copy(hasher, reader)
	return hasher.Sum(nil)
}

func HashLocal(fname string) []byte {
	file, err := os.Open(fname)
	if err != nil {
		panic(err.Error())
	}
	return HashFile(file)
}

func HashRemote(bucket *s3.Bucket, path string) []byte {
	file, err := bucket.GetReader(path)
	if err != nil {
		panic(err.Error())
	}
	return HashFile(file)
}

func UploadWalker(path string, info os.FileInfo, err error) error {
	fmt.Println(path)
	fmt.Println(info.Name())
	fmt.Println(info.IsDir())
	return nil
}

func main() {
	flag.Parse()
	if localPath == "" || bucketName == "" {
		flag.PrintDefaults()
		return
	}
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}

	s3Conn := s3.New(auth, aws.USWest)
	bucket := s3Conn.Bucket(bucketName)
	contents, err := bucket.List("/", "", "", 1000)
	if err != nil {
		fmt.Println(contents)
		panic(err.Error())
	}
	fmt.Println(contents)
	filepath.Walk(localPath, UploadWalker)
	return
}
