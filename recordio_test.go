package recordio_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/PaddlePaddle/recordio"
)

func ExampleWriter_Write() {
	var buf bytes.Buffer
	w := recordio.NewWriter(&buf, -1, -1)
	w.Write([]byte("Hello"))
	w.Write([]byte("World!"))
	w.Close()
}

func ExampleScanner_Scan() {
	f, err := os.Create("/tmp/example_recordio_0")
	if err != nil {
		panic(err)
	}

	w := recordio.NewWriter(f, -1, -1)
	w.Write([]byte("Hello"))
	w.Close()
	f.Close()

	f, err = os.Create("/tmp/example_recordio_1")
	if err != nil {
		panic(err)
	}

	w = recordio.NewWriter(f, -1, -1)
	w.Write([]byte("World!"))
	w.Close()
	f.Close()

	s, err := recordio.NewScanner("/tmp/example_recordio_*")
	if err != nil {
		panic(err)
	}

	for s.Scan() {
		fmt.Println(string(s.Record()))
	}
	// Output:
	// Hello
	// World!
}

func ExampleRangeScanner_Scan() {
	const path = "/tmp/example_recordio_0"
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	w := recordio.NewWriter(f, -1, -1)
	w.Write([]byte("Hello"))
	w.Write([]byte("World!"))
	w.Close()
	f.Close()

	f, err = os.Open(path)
	if err != nil {
		panic(err)
	}

	idx, err := recordio.LoadIndex(f)
	if err != nil {
		panic(err)
	}

	f.Seek(0, 0)
	s := recordio.NewRangeScanner(f, idx, 1, -1)
	if err != nil {
		panic(err)
	}

	for s.Scan() {
		fmt.Println(string(s.Record()))
	}
	// Output: World!
}

func TestWriteRead(t *testing.T) {
	const total = 1000
	var buf bytes.Buffer
	w := recordio.NewWriter(&buf, 0, -1)
	for i := 0; i < total; i++ {
		_, err := w.Write(make([]byte, i))
		if err != nil {
			t.Fatal(err)
		}
	}
	w.Close()

	idx, err := recordio.LoadIndex(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	if idx.NumRecords != total {
		t.Fatal("num record does not match:", idx.NumRecords, total)
	}

	s := recordio.NewRangeScanner(bytes.NewReader(buf.Bytes()), idx, -1, -1)
	i := 0
	for s.Scan() {
		if !reflect.DeepEqual(s.Record(), make([]byte, i)) {
			t.Fatal("not equal:", len(s.Record()), len(make([]byte, i)))
		}
		i++
	}

	if i != total {
		t.Fatal("total count not match:", i, total)
	}
}

func TestChunkIndex(t *testing.T) {
	const total = 1000
	var buf bytes.Buffer
	w := recordio.NewWriter(&buf, 0, -1)
	for i := 0; i < total; i++ {
		_, err := w.Write(make([]byte, i))
		if err != nil {
			t.Fatal(err)
		}
	}
	w.Close()

	idx, err := recordio.LoadIndex(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	if idx.NumChunks() != total {
		t.Fatal("unexpected chunk num:", idx.NumChunks(), total)
	}

	for i := 0; i < total; i++ {
		newIdx := idx.ChunkIndex(i)
		s := recordio.NewRangeScanner(bytes.NewReader(buf.Bytes()), newIdx, -1, -1)
		j := 0
		for s.Scan() {
			if !reflect.DeepEqual(s.Record(), make([]byte, i)) {
				t.Fatal("not equal:", len(s.Record()), len(make([]byte, i)))
			}
			j++
		}
		if j != 1 {
			t.Fatal("unexpected record per chunk:", j)
		}
	}
}
