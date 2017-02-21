package web

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/goamz/goamz/s3"
	"gopkg.in/h2non/bimg.v1"
)

//UploadBase64Image resizes image and puts it into s3 bucket
func UploadBase64Image(bucket *s3.Bucket, image string, imagename string, width int) (string, error) {
	byt, err := base64.StdEncoding.DecodeString(strings.Split(image, "base64,")[1])
	if err != nil {
		log.Println(err)
		return "", err
	}

	meta := strings.Split(image, "base64,")[0]
	newmeta := strings.Replace(strings.Replace(meta, "data:", "", -1), ";", "", -1)

	opt := bimg.Options{
		Width: width,
	}

	newImage, err := bimg.Resize(byt, opt)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = bucket.Put(imagename, newImage, newmeta, s3.PublicReadWrite, s3.Options{})
	if err != nil {
		log.Println(err)
		return "", err
	}

	return bucket.URL(imagename), nil
}

//ResizeToThumbnailSize resizes image to a 300px width
func ResizeToThumbnailSize(byt []byte) ([]byte, error) {
	opt := bimg.Options{
		Width: 300,
	}

	newImage, err := bimg.Resize(byt, opt)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	//bimg.Write("smaller.jpg", newImage)
	return newImage, nil
}
