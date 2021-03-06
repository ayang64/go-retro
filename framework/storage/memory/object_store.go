package memory

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/retro-framework/go-retro/framework/packing"
	"github.com/retro-framework/go-retro/framework/types"
)

var (
	ErrNoSuchObject          = errors.New("no such object in object database")
	ErrUnableToInflateObject = errors.New("error running zlib inflate")
)

type ObjectStore struct {
	o map[string][]byte
}

func (os *ObjectStore) ListObjects() {
	for k, _ := range os.o {
		fmt.Println(k)
	}
}

func (os *ObjectStore) WritePacked(p types.HashedObject) (int, error) {

	if os.o == nil {
		os.o = make(map[string][]byte)
	}

	var (
		b bytes.Buffer
		k = p.Hash().String()
	)

	w := zlib.NewWriter(&b)
	w.Write(p.Contents())
	w.Close()

	if _, ok := os.o[k]; !ok {
		os.o[k] = b.Bytes()
		return len(b.Bytes()), nil
	}

	return 0, nil
}

// TODO: should also parse the aglo out of the string and set the PO Hash
// algo/etc to the right values., the new PackedObject could be kept and
// maybe simply take an AlgoName in the second position?
func (os *ObjectStore) RetrievePacked(s string) (types.HashedObject, error) {
	if poB, ok := os.o[s]; ok {

		b := bytes.NewReader(poB)
		r, err := zlib.NewReader(b)
		if err != nil {
			return nil, ErrUnableToInflateObject
		}
		r.Close()

		orig, _ := ioutil.ReadAll(r)

		return packing.NewPackedObject(string(orig)), nil
	}
	return nil, ErrNoSuchObject
}
