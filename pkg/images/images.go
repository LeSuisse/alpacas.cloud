package images

import (
	"errors"
	"github.com/disintegration/imaging"
	"image"
	"io/ioutil"
	"path"
	"path/filepath"
)

type Images []string

func New(imagesFolderPath string) (Images, error) {
	files, err := ioutil.ReadDir(imagesFolderPath)
	if err != nil {
		return nil, err
	}

	var im Images

	for _, file := range files {
		if file.Mode().IsRegular() && filepath.Ext(file.Name()) == ".jpg" {
			im = append(im, path.Join(imagesFolderPath, file.Name()))
		}
	}

	if len(im) <= 0 {
		return nil, errors.New("Can't find images")
	}

	return im, nil
}

func (images Images) Get() (image.Image, error) {
	return imaging.Open(images[randomSource.Intn(len(images))], imaging.AutoOrientation(true))
}

func (images Images) GetWithWidth(width int) (image.Image, error) {
	for i := 0; i < 10; i++ {
		im, err := images.Get()
		if err != nil {
			return nil, err
		}
		if im.Bounds().Max.X >= width {
			return imaging.Resize(im, width, 0, imaging.Lanczos), nil
		}
	}
	return nil, errors.New("Can't find an image with the right width")
}
