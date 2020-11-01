package images

import (
	"errors"
	"github.com/h2non/bimg"
	"io/ioutil"
	"math"
	"path"
	"path/filepath"
)

const (
	maxAttempts = 10
)

type Images []string
type ImageOpts struct {
	MaxWidth  int
	MaxHeight int
}
type OutputImage struct {
	Name string
	Data []byte
}

type RequestedSizeTooBigError struct {
}

func (e *RequestedSizeTooBigError) Error() string {
	return "can't find an image with an appropriate size"
}

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
		return nil, errors.New("can't find images")
	}

	return im, nil
}

func (images Images) Get(requestedOpts ImageOpts) (*OutputImage, error) {
	var errRequestedSizeTooBig *RequestedSizeTooBigError
	for i := 0; i < maxAttempts; i++ {
		img, err := images.getResizedImage(requestedOpts)
		if err == nil {
			return img, nil
		}
		if !errors.As(err, &errRequestedSizeTooBig) {
			return nil, err
		}
	}

	return nil, errRequestedSizeTooBig
}

func (images Images) getResizedImage(requestedOpts ImageOpts) (*OutputImage, error) {
	imagePath := images[randomSource.Intn(len(images))]
	inputBuf, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}

	image := bimg.NewImage(inputBuf)
	_, err = image.AutoRotate()
	if err != nil {
		return nil, err
	}

	imageSize, err := image.Size()
	if err != nil {
		return nil, err
	}

	maxWidth := requestedOpts.MaxWidth
	if maxWidth == 0 {
		maxWidth = imageSize.Width
	}
	maxHeight := requestedOpts.MaxHeight
	if maxHeight == 0 {
		maxHeight = imageSize.Height
	}

	scaleWidth := float64(maxWidth) / float64(imageSize.Width)
	scaleHeight := float64(maxHeight) / float64(imageSize.Height)

	if scaleWidth > 1 || scaleHeight > 1 {
		return nil, &RequestedSizeTooBigError{}
	}
	scale := math.Min(scaleWidth, scaleHeight)

	if scale != 1 {
		_, err = image.Resize(int(float64(imageSize.Width) * scale), int(float64(imageSize.Height) * scale))
		if err != nil {
			return nil, err
		}
	}

	return &OutputImage{
		Name: filepath.Base(imagePath),
		Data: image.Image(),
	}, nil
}

func (images Images) GetPlaceHolder(requestedOpts ImageOpts) (*OutputImage, error) {
	var errRequestedSizeTooBig *RequestedSizeTooBigError
	for i := 0; i < maxAttempts; i++ {
		img, err := images.getPlaceHolderImage(requestedOpts)
		if err == nil {
			return img, nil
		}
		if !errors.As(err, &errRequestedSizeTooBig) {
			return nil, err
		}
	}

	return nil, errRequestedSizeTooBig
}

func (images Images) getPlaceHolderImage(requestedOpts ImageOpts) (*OutputImage, error) {
	imagePath := images[randomSource.Intn(len(images))]
	inputBuf, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}

	image := bimg.NewImage(inputBuf)
	_, err = image.AutoRotate()
	if err != nil {
		return nil, err
	}

	imageSize, err := image.Size()
	if err != nil {
		return nil, err
	}

	width := requestedOpts.MaxWidth
	if width == 0 {
		width = imageSize.Width
	}
	height := requestedOpts.MaxHeight
	if height == 0 {
		height = imageSize.Height
	}

	if width > imageSize.Width || height > imageSize.Height {
		return nil, &RequestedSizeTooBigError{}
	}

	outputImg, err := image.SmartCrop(width, height)
	if err != nil {
		return nil, err
	}

	return &OutputImage{
		Name: filepath.Base(imagePath),
		Data: outputImg,
	}, nil
}
