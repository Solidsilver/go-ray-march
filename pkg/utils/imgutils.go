package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/rs/zerolog/log"
)

func EncodePNGToPath(imgPath string, img image.Image) error {
	out, err := os.Create(imgPath)
	if err != nil {
		log.Err(err).Msgf("Could not create output file: %v", imgPath)
		// log.Error().Msgf("Could not create output file: %v", imgPath)
		return err
	}
	defer out.Close()
	err = png.Encode(out, img)
	if err != nil {
		log.Err(err).Msg("Could not encode output image")
		// log.Error().Msg("Could not encode output image")
	}
	return err
}

const (
	IMG_PNG  = "png"
	IMG_JPEG = "jpeg"
)

func EncodeImageToBytes(img image.Image, format string) ([]byte, error) {
	var err error
	buf := new(bytes.Buffer)

	switch format {
	case IMG_PNG:
		err = png.Encode(buf, img)
	case IMG_JPEG:
		err = jpeg.Encode(buf, img, nil)
	}
	return buf.Bytes(), err

}
