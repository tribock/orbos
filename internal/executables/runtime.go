//go:generate goderive .

package executables

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/caos/orbos/internal/helpers"
	"github.com/pkg/errors"
)

var executables map[string][]byte

var populate = func() {}

func Populate() {
	populate()
}

func PreBuilt(name string) ([]byte, error) {
	executable, ok := executables[name]
	if !ok {
		return nil, errors.Errorf("%s was not prebuilt", name)
	}
	return executable, nil
}

func PreBuild(packables <-chan PackableTuple) (err error) {
	sp := selfPath()
	tmpFile := filepath.Join(sp, "prebuilt.tmp")
	outFile := filepath.Join(sp, "prebuilt.go")
	prebuilt, err := os.Create(tmpFile)
	if err != nil {
		return errors.Wrapf(err, "creating %s failed", tmpFile)
	}
	defer func() {
		if err != nil {
			os.Remove(tmpFile)
		}
		err = os.Rename(tmpFile, outFile)
	}()

	if _, err = prebuilt.WriteString(`package executables

func init() {
	populate = func(){
		executables = map[string][]byte{`); err != nil {
		return err
	}

	for pt := range deriveFmapPack(pack, packables) {
		packable, packed, packErr := pt()
		err = helpers.Concat(err, packErr)
		if packErr != nil {
			continue
		}

		_, packErr = prebuilt.WriteString(fmt.Sprintf(`
		"%s": unpack("%s"),`, packable.key, *packed))
		err = helpers.Concat(err, packErr)
	}

	if err != nil {
		os.Remove(prebuilt.Name())
		return err
	}

	_, err = prebuilt.WriteString(`
		}
	}
}
`)
	return err
}

func packedTupleFunc(packable *packable) func(*string, error) packedTuple {
	return func(packed *string, err error) packedTuple {
		return deriveTuplePacked(packable, packed, err)
	}
}

type packable struct {
	key  string
	data io.ReadCloser
}

type PackableTuple func() (*packable, error)

func NewPackableTuple(key string, data io.ReadCloser) PackableTuple {
	return deriveTuplePackable(&packable{
		key:  key,
		data: data,
	}, nil)
}

type packedTuple func() (*packable, *string, error)

func pack(packableTuple PackableTuple) packedTuple {

	packable, err := packableTuple()
	defer func() {
		if packable != nil {
			packable.data.Close()
		}
	}()

	packedTuple := packedTupleFunc(packable)
	if err != nil {
		return packedTuple(nil, err)
	}

	gzipBuffer := new(bytes.Buffer)
	defer gzipBuffer.Reset()

	gzipWriter := gzip.NewWriter(gzipBuffer)
	_, err = io.Copy(gzipWriter, packable.data)
	if err != nil {
		return packedTuple(nil, errors.Wrap(err, "gzipping failed"))
	}

	if err := packable.data.Close(); err != nil {
		return packedTuple(nil, errors.Wrap(err, "closing data failed"))
	}

	if err := gzipWriter.Close(); err != nil {
		return packedTuple(nil, errors.Wrap(err, "closing gzip writer failed"))
	}

	packed := base64.StdEncoding.EncodeToString(gzipBuffer.Bytes())

	return packedTuple(&packed, nil)
}

func unpack(executable string) []byte {
	gzipNodeAgent, err := base64.StdEncoding.DecodeString(executable)
	if err != nil {
		panic(errors.Wrap(err, "decoding node agent from base64 failed"))
	}
	bytesReader := bytes.NewReader(gzipNodeAgent)

	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		panic(errors.Wrap(err, "ungzipping node agent failed"))
	}
	defer gzipReader.Close()

	unpacked, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		panic(errors.Wrap(err, "reading unpacked node agent failed"))
	}
	return unpacked
}
