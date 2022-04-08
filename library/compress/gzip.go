package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

const (
	bufferChanSize = 4096 * 2
	initBytesSize  = 40 * 1024
)

type ICodec interface {
	// 标准压缩、解压接口
	Decode(data []byte) ([]byte, error)
	Encode(data []byte) ([]byte, error)
	Prefix() byte
}

var (
	bufferChan chan *bytes.Buffer
)

func init() {
	bufferChan = make(chan *bytes.Buffer, bufferChanSize)
}

func getBuffer() *bytes.Buffer {
	select {
	case buff := <-bufferChan:
		buff.Reset()
		return buff
	default:
		data := make([]byte, 0, initBytesSize)
		return bytes.NewBuffer(data)
	}
}

func backBuffer(buff *bytes.Buffer) {
	buff.Reset()
	select {
	case bufferChan <- buff:
		return
	default:
		return
	}
}

const (
	gzipWriterChanSize = 256
	gzipReaderChanSize = 4096
)

var gzipWriterChan chan *gzip.Writer
var gzipReaderChan chan *gzip.Reader

func init() {
	gzipWriterChan = make(chan *gzip.Writer, gzipWriterChanSize)
	gzipReaderChan = make(chan *gzip.Reader, gzipReaderChanSize)
}

type gzipCodec struct{}

func (g *gzipCodec) Decode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	reader, e := g.getReader(bytes.NewReader(data))
	if e != nil {
		return nil, e
	}
	defer g.backReader(reader)

	return ioutil.ReadAll(reader)
}

func (g *gzipCodec) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	buff := getBuffer()
	defer backBuffer(buff)

	writer := g.getWriter(buff)
	defer g.backWriter(writer)

	_, e := writer.Write(data)
	if e != nil {
		return nil, e
	}
	e = writer.Close()
	if e != nil {
		return nil, e
	}

	ret := make([]byte, len(buff.Bytes()))
	copy(ret, buff.Bytes())
	return ret, nil
}

func (g *gzipCodec) getWriter(buff *bytes.Buffer) *gzip.Writer {

	select {
	case writer := <-gzipWriterChan:
		writer.Reset(buff)
		return writer
	default:
		return gzip.NewWriter(buff)
	}
}

func (g *gzipCodec) backWriter(w *gzip.Writer) {
	w.Reset(nil)
	select {
	case gzipWriterChan <- w:
		return
	default:
	}
}

func (g *gzipCodec) getReader(r io.Reader) (*gzip.Reader, error) {
	select {
	case reader := <-gzipReaderChan:
		e := reader.Reset(r)
		return reader, e
	default:
		return gzip.NewReader(r)
	}
}

func (g *gzipCodec) backReader(reader *gzip.Reader) {
	select {
	case gzipReaderChan <- reader:
		return
	default:
	}
}
