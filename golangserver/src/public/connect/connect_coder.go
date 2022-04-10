package connect

import (
	//. "game/common"
	"io"
	"io/ioutil"
	"public/link"
)

type ConnectCodec struct {
}

type ConnectEncoder struct {
	w io.Writer
}

type ConnectDecoder struct {
	r io.Reader
}

func (codec ConnectCodec) NewEncoder(w io.Writer) link.Encoder {
	return &ConnectEncoder{w}
}

func (codec ConnectCodec) NewDecoder(r io.Reader) link.Decoder {
	return &ConnectDecoder{r}
}

func (encoder *ConnectEncoder) Encode(msg interface{}, args ...interface{}) error {
	_, err := encoder.w.Write(msg.([]byte))

	return err
}

func (decoder *ConnectDecoder) Decode(msg interface{}, msg_type *uint16) error {
	// We use ReadAll() here because we know the reader is a buffer object not a real net.Conn
	d, err := ioutil.ReadAll(decoder.r)
	if err != nil {
		return err
	}
	*(msg.(*[]byte)) = d
	return nil
}
