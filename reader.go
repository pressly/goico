package ico

import (
	"bufio"
	"image"
	"io"
)

const magicString = "?0?0?1?0"

// A FormatError reports that the input is not a valid ICO.
type FormatError string

func (e FormatError) Error() string { return "invalid ICO format: " + string(e) }

// An UnsupportedError reports that the input uses a valid but unimplemented ICO feature.
type UnsupportedError string

func (e UnsupportedError) Error() string { return "unsupported ICO feature: " + string(e) }

// If the io.Reader does not also have ReadByte, then decode will introduce its own buffering.
type reader interface {
	io.Reader
	io.ByteReader
}

type decoder struct {
	r     io.Reader
	num   int
	image []image.Image
	tmp   [1024]byte
}

func (d *decoder) decode(r io.Reader, configOnly bool) error {
	// Add buffering if r does not provide ReadByte.
	if rr, ok := r.(reader); ok {
		d.r = rr
	} else {
		d.r = bufio.NewReader(r)
	}

	if err := d.readHeader(); err != nil {
		return err
	}
	if configOnly {
		return nil
	}

	for {
		// TODO: Write decode loop here
	}
	return nil
}

func (d *decoder) readHeader() error {
	_, err := io.ReadFull(d.r, d.tmp[0:6])
	if err != nil {
		return err
	}
	if d.tmp[0] != 0 {
		// TODO: Error out here
	}
	if d.tmp[1] != 1 {
		// TODO: Error out here
	}
	d.num = int(d.tmp[2]) // The number of images
	return nil
}

func (d *decoder) parseImage(b []byte) (image.Image, error) {
	var img image.Image
	return img, nil
}

func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	if err := d.decode(r, false); err != nil {
		return nil, err
	}
	return d.image[0], nil
}

func DecodeAll(r io.Reader) (*ICO, error) {
	var d decoder
	if err := d.decode(r, false); err != nil {
		return nil, err
	}
	ico := &ICO{
		Num:   len(d.image),
		Image: d.image,
	}
	return ico, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var d decoder
	if err := d.decode(r, false); err != nil {
		return image.Config{}, err
	}
	return image.Config{
	// TODO: Need to fill this in with ???
	}, nil
}

func init() {
	image.RegisterFormat("ico", magicString, Decode, DecodeConfig)
}
