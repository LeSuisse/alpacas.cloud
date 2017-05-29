package images

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
)

type Images []string

func LoadImages(imagesFolderPath string) (Images, error) {
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

func (images Images) Random() string {
	return images[randomSource.Intn(len(images))]
}
