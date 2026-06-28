package images

import (
	"errors"
	"github.com/davidbyttow/govips/v2/vips"
	"math"
	"os"
	"path"
	"path/filepath"
)

func init() {
	vips.LoggingSettings(nil, vips.LogLevelError)
}

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
	files, err := os.ReadDir(imagesFolderPath)
	if err != nil {
		return nil, err
	}

	var im Images

	for _, file := range files {
		if file.Type().IsRegular() && filepath.Ext(file.Name()) == ".jpg" {
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
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err := vips.NewImageFromReader(file)
	if err != nil {
		return nil, err
	}
	defer image.Close()
	err = image.AutoRotate()
	if err != nil {
		return nil, err
	}

	imageWidth := image.Width()
	imageHeight := image.Height()

	maxWidth := requestedOpts.MaxWidth
	if maxWidth == 0 {
		maxWidth = imageWidth
	}
	maxHeight := requestedOpts.MaxHeight
	if maxHeight == 0 {
		maxHeight = imageHeight
	}

	scaleWidth := float64(maxWidth) / float64(imageWidth)
	scaleHeight := float64(maxHeight) / float64(imageHeight)

	if scaleWidth > 1 || scaleHeight > 1 {
		return nil, &RequestedSizeTooBigError{}
	}
	scale := math.Min(scaleWidth, scaleHeight)

	if scale != 1 {
		err = image.Resize(scale, vips.KernelLanczos3)
		if err != nil {
			return nil, err
		}
	}

	imageData, _, err := image.ExportNative()
	if err != nil {
		return nil, err
	}
	return &OutputImage{
		Name: filepath.Base(imagePath),
		Data: imageData,
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
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err := vips.NewImageFromReader(file)
	if err != nil {
		return nil, err
	}
	defer image.Close()
	err = image.AutoRotate()
	if err != nil {
		return nil, err
	}

	imageWidth := image.Width()
	imageHeight := image.Height()

	width := requestedOpts.MaxWidth
	if width == 0 {
		width = imageWidth
	}
	height := requestedOpts.MaxHeight
	if height == 0 {
		height = imageHeight
	}

	if width > imageWidth || height > imageHeight {
		return nil, &RequestedSizeTooBigError{}
	}

	err = image.SmartCrop(width, height, vips.InterestingCentre)
	if err != nil {
		return nil, err
	}
	outputImg, _, err := image.ExportNative()
	if err != nil {
		return nil, err
	}

	return &OutputImage{
		Name: filepath.Base(imagePath),
		Data: outputImg,
	}, nil
}
