package main

/*
#include <string.h>
#include <stdlib.h>

typedef int reader;
typedef int writer;
*/
import "C"

import (
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/visenze/recordio"
)

var nullPtr = unsafe.Pointer(uintptr(0))

type writer struct {
	w *recordio.Writer
	f *os.File
}

type reader struct {
	scanner *recordio.Scanner
}

//export create_recordio_writer
func create_recordio_writer(path *C.char, maxChunkSize C.int, compressor C.int) C.writer {
	p := C.GoString(path)
	f, err := os.Create(p)
	if err != nil {
		log.Println(err)
		return -1
	}

	w := recordio.NewWriter(f, int(maxChunkSize), int(compressor))
	writer := &writer{f: f, w: w}
	return addWriter(writer)
}

//export recordio_write
func recordio_write(writer C.writer, buf *C.uchar, size C.int) C.int {
	w := getWriter(writer)

	// Make a copy of the C buffer rather than create a slice
	// backed by the C buffer. This is because RecordIO caches the
	// slice in memory until the max chunk size is reached and
	// then dump the slice to disk. At which point the C buffer is
	// no longer valid.
	b := make([]byte, int(size))
	for i := 0; i < int(size); i++ {
		ptr := (*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(buf)) + uintptr(i)))
		b[i] = byte(*ptr)
	}

	c, err := w.w.Write(b)
	if err != nil {
		log.Println(err)
		return -1
	}
	return C.int(c)
}

//export release_recordio_writer
func release_recordio_writer(writer C.writer) {
	w := removeWriter(writer)
	w.w.Close()
	w.f.Close()
}

//export create_recordio_reader
func create_recordio_reader(path *C.char) C.reader {
	p := C.GoString(path)
	s, err := recordio.NewScanner(strings.Split(p, ",")...)
	if err != nil {
		log.Println(err)
		return -1
	}

	r := &reader{scanner: s}
	return addReader(r)
}

//export recordio_read
func recordio_read(reader C.reader, record **C.uchar) C.int {
	r := getReader(reader)
	if r.scanner.Scan() {
		buf := r.scanner.Record()
		if len(buf) == 0 {
			*record = (*C.uchar)(nullPtr)
			return 0
		}

		size := C.int(len(buf))
		*record = (*C.uchar)(C.malloc(C.size_t(len(buf))))
		C.memcpy(unsafe.Pointer(*record), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
		return size
	}

	return -1
}

//export release_recordio_reader
func release_recordio_reader(reader C.reader) {
	r := removeReader(reader)
	r.scanner.Close()
}

//export mem_free
func mem_free(p unsafe.Pointer) {
	// "free" may be a better name for this function, but doing so
	// will cause calling any function of this library from Python
	// ctypes hanging.
	C.free(p)
}

func main() {} // Required but ignored
