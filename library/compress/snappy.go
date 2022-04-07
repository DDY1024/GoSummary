package compress

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"

	"github.com/golang/snappy"
)

var (
	snappyReaderPool = &sync.Pool{
		New: func() interface{} {
			return snappy.NewReader(nil)
		},
	}

	snappyWriterPool = &sync.Pool{
		New: func() interface{} {
			return snappy.NewBufferedWriter(nil)
		},
	}
)

type snappyCodec struct{}

func (s *snappyCodec) getReader(r io.Reader) *snappy.Reader {
	rd := snappyReaderPool.Get().(*snappy.Reader)
	rd.Reset(r)
	return rd
}

func (s *snappyCodec) getWriter(b *bytes.Buffer) *snappy.Writer {
	w := snappyWriterPool.Get().(*snappy.Writer)
	w.Reset(b)
	return w
}

func (s *snappyCodec) putWriter(w *snappy.Writer) {
	snappyWriterPool.Put(w)
}

func (s *snappyCodec) putReader(r *snappy.Reader) {
	snappyReaderPool.Put(r)
}

func (s *snappyCodec) Decode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	reader := s.getReader(bytes.NewReader(data))
	defer s.putReader(reader)
	return ioutil.ReadAll(reader)
}

func (s *snappyCodec) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Tips: 此处预估下内存占用大小，暂且用 0xff 代替
	buf := bytes.NewBuffer(make([]byte, 0, 0xff))
	writer := s.getWriter(buf)
	defer s.putWriter(writer)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
