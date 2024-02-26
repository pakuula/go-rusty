package io_r

import (
	"io"

	"github.com/pakuula/go-rusty/result"
)

type Reader interface {
	Read(p []byte) result.Result[int]
}

type wrappedReader struct {
	io.Reader
}

func WrapReader(r io.Reader) Reader {
	return wrappedReader{Reader: r}
}

func (self wrappedReader) Read(p []byte) result.Result[int] {
	return result.Wrap(self.Reader.Read(p))
}

func Read(r io.Reader, p []byte) result.Result[int] {
	return result.Wrap(r.Read(p))
}

func ReadAt(r io.ReaderAt, p []byte, off int64) result.Result[int] {
	return result.Wrap(r.ReadAt(p, off))
}

type Writer interface {
	Write(p []byte) result.Result[int]
}

type wrappedWriter struct {
	io.Writer
}

func WrapWriter(r io.Writer) Writer {
	return wrappedWriter{Writer: r}
}

func (self wrappedWriter) Write(p []byte) result.Result[int] {
	return result.Wrap(self.Writer.Write(p))
}

func Write(w io.Writer, p []byte) result.Result[int] {
	return result.Wrap(w.Write(p))
}

func WriteAt(w io.WriterAt, p []byte, off int64) result.Result[int] {
	return result.Wrap(w.WriteAt(p, off))
}

func Copy(dst io.Writer, src io.Reader) result.Result[int64] {
	return result.Wrap(io.Copy(dst, src))
}

func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) result.Result[int64] {
	return result.Wrap(io.CopyBuffer(dst, src, buf))
}

func CopyN(dst io.Writer, src io.Reader, n int64) result.Result[int64] {
	return result.Wrap(io.CopyN(dst, src, n))
}

func ReadAll(r io.Reader) result.Result[[]byte] {
	return result.Wrap(io.ReadAll(r))
}

func ReadAtLeast(r io.Reader, buf []byte, min int) result.Result[int] {
	return result.Wrap(io.ReadAtLeast(r, buf, min))
}

func ReadFull(r io.Reader, buf []byte) result.Result[int] {
	return result.Wrap(io.ReadFull(r, buf))
}

func WriteString(w io.Writer, s string) result.Result[int] {
	return result.Wrap(io.WriteString(w, s))
}
