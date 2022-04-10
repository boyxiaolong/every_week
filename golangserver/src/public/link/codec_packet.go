package link

import (
	"bufio"
	"encoding/binary"
	"errors"
	//"fmt"
	"io"
	"sync"
	"time"
)

type ByteOrder binary.ByteOrder

var (
	BigEndian    = binary.BigEndian
	LittleEndian = binary.LittleEndian
)

var (
	ErrPacketUnsupported = errors.New("funny/link: unsupported packet type")
	ErrPacketTooLarge    = errors.New("funny/link: too large packet")
	ErrPacketNoReadAll   = errors.New("funny/link: no read all content from packet")
	ErrPacketReadFial    = errors.New("funny/link: read fail")
)

func Packet(n, maxPacketSize, readBufferSize int, byteOrder ByteOrder, base CodecType) CodecType {
	if n != 1 && n != 2 && n != 4 && n != 8 {
		panic(ErrPacketUnsupported)
	}
	return &packetCodecType{
		n:              n,
		base:           base,
		maxPacketSize:  maxPacketSize,
		readBufferSize: readBufferSize,
		byteOrder:      byteOrder,
	}
}

type packetCodecType struct {
	n              int
	maxPacketSize  int
	readBufferSize int
	base           CodecType
	encoderPool    sync.Pool
	decoderPool    sync.Pool
	byteOrder      binary.ByteOrder
}

func (codecType *packetCodecType) NewEncoder(w io.Writer) Encoder {
	encoder, ok := codecType.encoderPool.Get().(*packetEncoder)
	if ok {
		encoder.writer = w
	} else {
		encoder = &packetEncoder{
			writer: w,
			parent: codecType,
		}
		encoder.buffer.data = make([]byte, 1024)
		encoder.buffer.n = codecType.n
		encoder.buffer.max = codecType.n + codecType.maxPacketSize
		encoder.encodeHead = codecType.encodeHead4
	}

	encoder.base = codecType.base.NewEncoder(&encoder.buffer)
	return encoder
}

func (codecType *packetCodecType) encodeHead4(b []byte, id uint16) ([]byte, error) {
	//mt.Printf("%v2222", b)
	if n := len(b) - 4; n <= codecType.maxPacketSize {
		size := len(b) + 4
		bufSize := Int16ToBytes(uint16(size))
		msgId := Int16ToBytes(id)
		bufSize = append(bufSize, msgId...)
		b = append(bufSize, b...)
		//fmt.Printf("%v", b)
		return b, nil
	}

	print("encoder too large packet:", time.Now().Format("2015-10-12 00:00:00"))
	return b, ErrPacketTooLarge
}

func Int16ToBytes(i uint16) []byte {
	var buf = make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(i))
	return buf
}

func (codecType *packetCodecType) NewDecoder(r io.Reader) Decoder {
	decoder, ok := codecType.decoderPool.Get().(*packetDecoder)
	if ok {
		decoder.reader.R.(*bufio.Reader).Reset(r)
	} else {
		decoder = &packetDecoder{
			n:      codecType.n,
			parent: codecType,
		}
		decoder.reader.R = bufio.NewReaderSize(r, codecType.readBufferSize)
		decoder.decodeHead = codecType.decodeHead2
	}
	decoder.base = codecType.base.NewDecoder(&decoder.reader)
	return decoder
}

func (codecType *packetCodecType) decodeHead1(b []byte) int {
	if n := int(b[0]); n <= 254 && n <= codecType.maxPacketSize {
		return n
	}
	print("decode too large packet:", time.Now().Format("2015-10-12 00:00:00"))
	return -1
}

func (codecType *packetCodecType) decodeHead2(b []byte) int {
	if n := int(codecType.byteOrder.Uint16(b)); n > 0 && n <= 65534 && n <= codecType.maxPacketSize {
		return n
	}
	//panic(ErrPacketTooLarge)

	print("decode too large packet:", time.Now().Format("2015-10-12 00:00:00"))
	return -1
}

func (codecType *packetCodecType) decodeHead4(b []byte) int {
	if n := int(codecType.byteOrder.Uint32(b)); n > 0 && n <= codecType.maxPacketSize {
		return n
	}
	print("decode too large packet:", time.Now().Format("2015-10-12 00:00:00"))
	return -1
}

func (codecType *packetCodecType) decodeHead8(b []byte) int {
	if n := int(codecType.byteOrder.Uint64(b)); n > 0 && n <= codecType.maxPacketSize {
		return n
	}
	print("decode too large packet:", time.Now().Format("2015-10-12 00:00:00"))
	return -1
}

type packetEncoder struct {
	base       Encoder
	buffer     PacketBuffer
	writer     io.Writer
	parent     *packetCodecType
	encodeHead func([]byte, uint16) ([]byte, error)
}

type packetDecoder struct {
	n          int
	base       Decoder
	head       [8]byte
	reader     io.LimitedReader
	parent     *packetCodecType
	decodeHead func([]byte) int
}

func (encoder *packetEncoder) Encode(msg interface{}, args ...interface{}) (err error) {
	if len(args) == 0 {
		panic("encoder args empty")
		return
	}
	//encoder.buffer.reset()
	//if err = encoder.base.Encode(msg); err != nil {
	//	return
	//}

	//b := encoder.buffer.bytes()
	newb, err := encoder.encodeHead(msg.([]byte), uint16(args[0].(int)))

	if err != nil {
		return
	}

	_, err = encoder.writer.Write(newb)

	return
}

func (decoder *packetDecoder) Decode(msg interface{}, msg_type *uint16) (err error) {
	head := decoder.head[:decoder.n]
	if _, err = io.ReadFull(decoder.reader.R, head); err != nil {
		return
	}

	head_size := head[:2]
	head_id := head[2:4]
	decoder.reader.N = int64(decoder.decodeHead(head_size)) - int64(decoder.n)
	*msg_type = uint16(decoder.decodeHead(head_id))

	if decoder.reader.N == -1 {
		return ErrPacketTooLarge
	}

	if decoder.base == nil {
		//println("read package fail.......................................................")
		return ErrPacketReadFial
	}

	if err = decoder.base.Decode(msg, msg_type); err != nil {
		return
	}

	//GStdout.Error("msg:%v", msg)

	if decoder.reader.N != 0 {
		println("packet no read all:", decoder.reader.N)
		return ErrPacketNoReadAll
	}
	return
}

func (encoder *packetEncoder) Dispose() {
	if d, ok := encoder.base.(Disposeable); ok {
		d.Dispose()
	}
	encoder.base = nil
	encoder.parent.encoderPool.Put(encoder)
}

func (decoder *packetDecoder) Dispose() {
	if d, ok := decoder.base.(Disposeable); ok {
		d.Dispose()
	}
	decoder.base = nil
	decoder.parent.decoderPool.Put(decoder)
}

type PacketBuffer struct {
	data []byte
	n    int
	max  int
	wpos int
}

func (pb *PacketBuffer) bytes() []byte {
	return pb.data[:pb.wpos]
}

func (pb *PacketBuffer) reset() {
	pb.wpos = pb.n
}

func (pb *PacketBuffer) gorws(n int) {
	if newLen := pb.wpos + n; newLen > len(pb.data) {
		newData := make([]byte, newLen, newLen+512)
		copy(newData, pb.data)
		pb.data = newData
	}
}

func (pb *PacketBuffer) Next(n int) (b []byte) {
	pb.gorws(n)
	n += pb.wpos
	if n > pb.max {
		panic(ErrPacketTooLarge)
	}
	b = pb.data[pb.wpos:n]
	pb.wpos = n
	return
}

func (pb *PacketBuffer) Write(b []byte) (int, error) {
	n := len(b)
	copy(pb.Next(n), b)
	return n, nil
}
