package images

import (
	"errors"
	"github.com/discordapp/lilliput"
	"io/ioutil"
	"math"
	"path"
	"path/filepath"
)

const (
	maxAttempts        = 10
	maxOutputImageSize = 5 * 1024 * 1024
)

type Images []string
type ImageOpts struct {
	MaxWidth  int
	MaxHeight int
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

func (images Images) Get(requestedOpts ImageOpts) ([]byte, error) {
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

func (images Images) getResizedImage(requestedOpts ImageOpts) ([]byte, error) {
	inputBuf, err := ioutil.ReadFile(images[randomSource.Intn(len(images))])
	if err != nil {
		return nil, err
	}

	decoder, err := lilliput.NewDecoder(inputBuf)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	header, err := decoder.Header()
	if err != nil {
		return nil, err
	}

	maxWidth := requestedOpts.MaxWidth
	if maxWidth == 0 {
		maxWidth = header.Width()
	}
	maxHeight := requestedOpts.MaxHeight
	if maxHeight == 0 {
		maxHeight = header.Height()
	}

	scaleWidth := float64(maxWidth) / float64(header.Width())
	scaleHeight := float64(maxHeight) / float64(header.Height())

	if scaleWidth > 1 || scaleHeight > 1 {
		return nil, &RequestedSizeTooBigError{}
	}
	scale := math.Min(scaleWidth, scaleHeight)

	maxSize := header.Height()
	if header.Width() > maxSize {
		maxSize = header.Width()
	}

	ops := lilliput.NewImageOps(maxSize)
	defer ops.Close()

	opts := &lilliput.ImageOptions{
		FileType:             ".jpeg",
		NormalizeOrientation: true,
		ResizeMethod:         lilliput.ImageOpsNoResize,
		EncodeOptions:        map[int]int{lilliput.JpegQuality: 85},
	}

	if scale != 1 {
		opts.ResizeMethod = lilliput.ImageOpsFit
		opts.Width = int(float64(header.Width()) * scale)
		opts.Height = int(float64(header.Height()) * scale)
	}

	outputImg := make([]byte, maxOutputImageSize)

	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		return nil, err
	}

	return outputImg, nil
}

func (images Images) GetPlaceHolder(requestedOpts ImageOpts) ([]byte, error) {
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

func (images Images) getPlaceHolderImage(requestedOpts ImageOpts) ([]byte, error) {
	inputBuf, err := ioutil.ReadFile(images[randomSource.Intn(len(images))])
	if err != nil {
		return nil, err
	}

	decoder, err := lilliput.NewDecoder(inputBuf)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	header, err := decoder.Header()
	if err != nil {
		return nil, err
	}

	width := requestedOpts.MaxWidth
	if width == 0 {
		width = header.Width()
	}
	height := requestedOpts.MaxHeight
	if height == 0 {
		height = header.Height()
	}

	if width > header.Width() || height > header.Height() {
		return nil, &RequestedSizeTooBigError{}
	}

	maxSize := header.Height()
	if header.Width() > maxSize {
		maxSize = header.Width()
	}

	ops := lilliput.NewImageOps(maxSize)
	defer ops.Close()

	opts := &lilliput.ImageOptions{
		FileType:             ".jpeg",
		NormalizeOrientation: true,
		ResizeMethod:         lilliput.ImageOpsFit,
		EncodeOptions:        map[int]int{lilliput.JpegQuality: 85},
		Width:                width,
		Height:               height,
	}

	outputImg := make([]byte, maxOutputImageSize)

	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		return nil, err
	}

	return outputImg, nil
}
