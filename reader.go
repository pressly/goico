package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"runtime"

	_ "image/png"
)

var magicString = string([]byte{0010})

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
	r     reader
	num   uint16
	dir   []entry
	image []image.Image
	tmp   [1024]byte
}

type entry struct {
	Width   uint8
	Height  uint8
	Palette uint8
	_       uint8 // Reserved byte
	Plane   uint16
	Bits    uint16
	Size    uint32
	Offset  uint32
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
	if err := d.readImageDir(); err != nil {
		return err
	}

	d.image = make([]image.Image, d.num)
	for i, entry := range d.dir {
		img, err := d.parseImage(entry)
		if err != nil {
			return err
		}
		d.image[i] = img
		runtime.GC()
	}
	return nil
}

func (d *decoder) readHeader() error {
	var first, second uint16
	binary.Read(d.r, binary.LittleEndian, &first)
	binary.Read(d.r, binary.LittleEndian, &second)
	binary.Read(d.r, binary.LittleEndian, &d.num)
	if first != 0 {
		return fmt.Errorf("First byte is %d instead of 0", first)
	}
	if second != 1 {
		return fmt.Errorf("Second byte is %d instead of 1, this is not an ICO file", second)
	}
	return nil
}

func (d *decoder) readImageDir() error {
	for i := 0; i < int(d.num); i++ {
		var e entry
		err := binary.Read(d.r, binary.LittleEndian, &e)
		if err != nil {
			return err
		}
		d.dir = append(d.dir, e)
	}
	fmt.Println(d.dir)
	return nil
}

func (d *decoder) parseImage(e entry) (image.Image, error) {
	off := int(e.Offset) - int((d.num*16)+6)
	b := make([]byte, e.Size)
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
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
		Num:   int(d.num),
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
