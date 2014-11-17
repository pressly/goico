package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"

	bmp "github.com/jsummers/gobmp"
)

// A FormatError reports that the input is not a valid ICO.
type FormatError string

func (e FormatError) Error() string { return "invalid ICO format: " + string(e) }

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
	cfg   image.Config
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
	if err := d.readImageDir(configOnly); err != nil {
		return err
	}
	if configOnly {
		cfg, err := d.parseConfig(d.dir[0])
		if err != nil {
			return err
		}
		d.cfg = cfg
	} else {
		d.image = make([]image.Image, d.num)
		for i, entry := range d.dir {
			img, err := d.parseImage(entry)
			if err != nil {
				return err
			}
			d.image[i] = img
		}
	}
	return nil
}

func (d *decoder) readHeader() error {
	var first, second uint16
	binary.Read(d.r, binary.LittleEndian, &first)
	binary.Read(d.r, binary.LittleEndian, &second)
	binary.Read(d.r, binary.LittleEndian, &d.num)
	if first != 0 {
		return FormatError(fmt.Sprintf("first byte is %d instead of 0", first))
	}
	if second != 1 {
		return FormatError(fmt.Sprintf("second byte is %d instead of 1", second))
	}
	return nil
}

func (d *decoder) readImageDir(configOnly bool) error {
	n := int(d.num)
	if configOnly {
		n = 1
	}
	for i := 0; i < n; i++ {
		var e entry
		err := binary.Read(d.r, binary.LittleEndian, &e)
		if err != nil {
			return err
		}
		d.dir = append(d.dir, e)
	}
	return nil
}

func (d *decoder) parseImage(e entry) (image.Image, error) {
	tmp := make([]byte, 14+e.Size)
	n, err := io.ReadFull(d.r, tmp[14:])
	if n != int(e.Size) {
		return nil, FormatError(fmt.Sprintf("only %d of %d bytes read.", n, e.Size))
	}
	if err != nil {
		return nil, err
	}

	_, err = png.DecodeConfig(bytes.NewReader(tmp[14:]))
	if err == nil {
		return png.Decode(bytes.NewReader(tmp[14:]))
	} else {
		tmp = d.setupBMP(e, tmp)
		return bmp.Decode(bytes.NewReader(tmp))
	}
}

func (d *decoder) parseConfig(e entry) (cfg image.Config, err error) {
	tmp := make([]byte, 14+e.Size)
	n, err := io.ReadFull(d.r, tmp[14:])
	if n != int(e.Size) {
		return cfg, fmt.Errorf("Only %d of %d bytes read.", n, e.Size)
	}
	if err != nil {
		return cfg, err
	}

	cfg, err = png.DecodeConfig(bytes.NewReader(tmp[14:]))
	if err != nil {
		tmp = d.setupBMP(e, tmp)
		cfg, err = bmp.DecodeConfig(bytes.NewReader(tmp))
	}
	return cfg, err
}

func (d *decoder) setupBMP(e entry, tmp []byte) []byte {
	var dibSize, w, h uint32
	binary.Read(bytes.NewReader(tmp[14:14+4]), binary.LittleEndian, &dibSize)
	binary.Read(bytes.NewReader(tmp[14+4:14+8]), binary.LittleEndian, &w)
	binary.Read(bytes.NewReader(tmp[14+8:14+12]), binary.LittleEndian, &h)

	if h > w {
		binary.LittleEndian.PutUint32(tmp[14+8:14+12], h/2)
	}

	copy(tmp[0:2], "\x42\x4D")
	binary.LittleEndian.PutUint32(tmp[2:6], e.Size+14)

	var numColors uint32
	binary.Read(bytes.NewReader(tmp[14+32:14+36]), binary.LittleEndian, &numColors)
	var offset uint32
	offset = 14 + dibSize + (numColors * 4)
	if dibSize > 40 {
		var iccSize uint32
		binary.Read(bytes.NewReader(tmp[14+dibSize-8:14+dibSize-4]), binary.LittleEndian, &iccSize)
		offset += iccSize
	}
	binary.LittleEndian.PutUint32(tmp[10:14], offset)

	return tmp
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
	if err := d.decode(r, true); err != nil {
		return image.Config{}, err
	}
	return d.cfg, nil
}

func init() {
	image.RegisterFormat("ico", "", Decode, DecodeConfig)
}
