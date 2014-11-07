package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	bmp "github.com/jsummers/gobmp"
)

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
		fmt.Println(d.dir)
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
	return nil
}

type DIB struct {
	HeaderSize uint32
	Width      uint32
	Height     uint32
	Planes     uint16
	BPP        uint16
	_          uint32
	Size       uint32
	_          uint32
	_          uint32
	NumColors  uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
	_          uint32
}

func (d *decoder) parseImage(e entry) (image.Image, error) {
	b := make([]byte, e.Size)
	n, err := io.ReadFull(d.r, b)
	if n != int(e.Size) {
		return nil, fmt.Errorf("Only %d of %d bytes read.", n, e.Size)
	}
	if err != nil {
		return nil, err
	}

	fileHeader := make([]byte, 14)
	copy(fileHeader[0:2], "\x42\x4D")
	binary.LittleEndian.PutUint32(fileHeader[2:6], e.Size+14)
	var iSize uint32
	binary.Read(bytes.NewReader(b[24:28]), binary.LittleEndian, &iSize)
	binary.LittleEndian.PutUint32(fileHeader[10:14], uint32(e.Size+14)-iSize)

	bb := append(fileHeader, b...)
	fmt.Println(len(bb))
	fmt.Println(len(bb))
	fmt.Println(len(bb))

	err = ioutil.WriteFile("/vagrant_data/noheader.png", b, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile("/vagrant_data/header.png", bb, os.ModePerm)
	if err != nil {
		return nil, err
	}
	img, err := bmp.Decode(bytes.NewReader(bb))
	return img, err
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
	image.RegisterFormat("ico", "", Decode, DecodeConfig)
}
