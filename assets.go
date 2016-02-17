package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	data []byte
	once sync.Once
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/assets/css/style.min.css": {
		local:   "assets/css/style.min.css",
		size:    3463,
		modtime: 1455731476,
		compressed: `
H4sIAAAJbogA/7RWTW+rOBTdjzT/Ab2qeo0UIwKh0xi91Yw0msVsRm8xWwMGrIIvsp2mLep/H9t8hBBI
0pHaSBXY9xzfLx9uoaqyKSjLC4U3nnf/EUP6Nl6IYpI85wL2PEUJlCDwne/741VBa0oU5tA9jfdqkEwx
4DihXFEx3pLsneIEXvRqBlyhjFSsfMP/QAwK1pJwiSQVLGt3D61HvudFnRdZln2Q5vgSKfqqUEoTEMQe
yYHTD4ILc0Qz3dQeUFEya5FBspcXLNx4rxRwHT5XRC8Ip1+JS0ieMcl0aGtX+62gGllJmhgmx9VO4owJ
qVBSsDLFMc1A0GsIrorW/iFYdZDGmOpM4m/for5G9esoqxhV8I6M10SgXJCUaeuHkmZq7Yg8Jg/e2v5W
jne/du52gfk5oXk53TalX50QH2j8zNSI1B6ydgy5o6DWDMah9tGWBUn9/HDGvTrZDkeOnO4YF6bQWZe+
JFz4GtpKfgnvlFNBW43/xZuxUjc0rgXkLMV//PtXRXL6U+gbqXuwcv9miQAJmXKH46QiQv1uSieV+PH9
zuv+vq8dytPZjT877M+3mv7YOKsoZbIuyRu2d+rDHS5FMwpz7K37uIrGOnVgqSrax0F1BC31ZX6ZucLN
YENiCeVe0Uj3HEZPnr5PpiI4HDi3nlm0+kBKlg9iZoXJytjGDWkVVUTkjCOLRr4FvSOmZeRVu3VNRZoO
bbzww/r1NtW5HdUkeyG1VNbArPMtEHtO4M3DBM0ElQXq4a3QPvlxuN3M2ZP0pbN1SWKS3kPSdBsHuzmI
OjClfZkcEYYkobpkmpDwhOqPCOG0bPr+sLLu6oBHxRyCIXsF0XlpW6HtCmnr4436ZcLWq3DTWuweDWjh
hGODzVLoLMKhOensvks6l3yb/WWw476QcsJw7Dvf9l1N0pTxvKfcXKK0DMhqw8IV8CK7i4NLNIpVtOnO
tTE43kVbx81Ic/Q7ON6W9rBNOOTY9rJ9199txRJSdrdOr189A8k9R3CcCtKdd90x/c2EESrY7p7S+ArK
VOUYz5ZWi/Y64XagaI76MhGOq1CbvbN8TWGFP9yDk5mpNP8ncrU5P9XloFj2maY/lezzNpxpoPaQfprp
EV1Mi/YOuTidTaao5lypFxVhEy4owtJotpyfZYwd5zpgENwPaWPcBIBOs9d3/xW6mXmymSjLDSRzI+an
WS5L1OOSQF3hLEmsRd+WUvWjB05IzZSu7PtM2U3/39D2E/1FNsAMQM3OBCOPO1k0YhcNX7Z0RO+5v2n6
UfNZQMd9sYV//eW/AAAA//8jh5imhw0AAA==
`,
	},

	"/assets/index.html": {
		local:   "assets/index.html",
		size:    5463,
		modtime: 1455731453,
		compressed: `
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
`,
	},

	"/assets/script/base.js": {
		local:   "assets/script/base.js",
		size:    2286,
		modtime: 1455702272,
		compressed: `
H4sIAAAJbogA/5RWUW/iRhB+R+I/jFAkTOSDXFX14ao+OMYJVo2NbHNpHtf2grc1Xnd3fVFU3X/vzBpI
aLlTayHBzn7zzXyzs2MWt+PRtjPiwF/hA9TGdPrTYrEXpu6LeSkPC234jrVFuegH1Hg0HkWi5K3mFfRt
xRWYmsM6zKEZzITwZfeqxL424JQz+OHu44+QWSLwpT4wgmy4OgithWxBaKi54sUr7BVrDa9c2CnOQe6g
rJnacxeMBNa+QseVRgdZGCZa0e6BQYmhxiOEmhp5tNyZF6Y4oitgWstSMCSESpb9gbeGGQq4Ew3X4FDi
k+zoMZnZKBVnzXgkWqvqtAcvWBDZG1BcGyVKInFBtGXTV5TFabsRB3EMQe62Ano8QtpeowhK1YWDrMSO
vrlV1vVFI3TtQiWIu+gNGjUZbTVdUrKQCjRvMDGkEJi6lfuWnwVR9h0V1RzLpMnyUsvDpRaBGe161WJQ
bp0qiWWzMX/npSEL4XeyaeQLqStlWwkSpT/RweW4yQr5hVs9wym30mC6QxZ0DN3b4R63dM2aBgp+rBpG
xhqz95IUZaANNoBgDXRS2ZD/lDq3KawCyJKH/MlLAwgz2KTJ53AZLGHiZbieuPAU5qtkmwMiUi/OnyF5
AC9+hl/DeOlC8NsmDbIMknQ8CtebKAzQGMZ+tF2G8SPco2Oc5BCF2NbImidAEY9cYZAR2zpI/RUuvfsw
CvNndzx6CPOYWB+SFDzYeGke+tvIS2GzTTdJFmACS+SNw/ghxTDBOojzOYZFGwSfcQHZyosiijUeeVsU
kFKK4Ceb5zR8XOWwSqJlgMb7AHPz7qNgiIW6/MgL1y4svbX3GFivBGlQHuGGBOFpFZCNInr48fMwiUmJ
n8R5iksXhab52fcpzAIXvDTMqCYPabJGjVRTdEksCzrGwUBD9YaLY0EIrbdZcGaEZeBFSJaRs1V5QuOp
3i7oZBcLeKr5cH06tsc2FLZRNTSS0V0bj26c02WezRVn1auz61t7JZ0Z/DUeAT6LW5w3TBlsMmHbibyB
QtjdBXhVZWOwVhyGC1swjR3d4MwYQDfOdG5k9wG7n4YNV9PZnFWVTwhnOvjxajr7+Q1eSGPk4f94fGHN
f6Dtkbb9Pg4lPXJjJeHOng9WnEldb5ypNV3DNrK06i/hJ+s1j+E1cIkfbNfQZAenwdWAoQFgh4tQmqbp
nz0OVJw+aGSDh+bD2Oi7CgVqaDmnYYGFQigOJWU5Z5cJXIa/haCt/n3y513P4LHjvuEKaWnAfIQDTqMX
nNgcXzy6fqd0mGn049Q9mpvw6Hul875Tl29m/NWFj3d3d3ALP92R7Zypjy+BP4DZEIAfzsoahnY4p3Pj
TI4dMpnNCfCNpG4cmsyzOdonJfHilLyOpOcLU8dA3hD9lzMBM0Y5Ezwd9mHIbPJeID2D2Xnv/x7y9ayb
fnwdFGPT+A1n6twA9A5lBd5ISX8pbhx8FVXyZTYv8NuZFhxbgfctne70mpBjDtOSSKfHQH8HAAD//37d
MsvuCAAA
`,
	},

	"/assets/script/functions.js": {
		local:   "assets/script/functions.js",
		size:    10676,
		modtime: 1455710940,
		compressed: `
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
`,
	},

	"/assets/script/require.config.js": {
		local:   "assets/script/require.config.js",
		size:    655,
		modtime: 1455727060,
		compressed: `
H4sIAAAJbogA/4ySS2rDMBCG94HcwTu1YKQkpRTcU5RSuihZKPbYkdHD0cikpfjulSwnftBAZmNZ//zf
PJCFUyss1Ehzo0tRPfyuV4mPA0f4sDJLCEcEhwxzKxrHSBr1hrsjZsmQHaJ+a8H+eANjvObftDKmksAb
EdCqv2NSHJDVp5DItsOBKqEv1BDKKNCu5+SFDo1J0xal5BYWoJhJa2Q7+kI3l/8F7wy+VbD3AIeGUKhG
wmf0sSe6obt/pWWlstW5E0b7vZDreZoQluq18CHxthtUPAo12yaJw5As+SJxs2Q/IZFhqpv6WH+GDVFA
gxNbeq2Vjtj9aOmm2L71O4hj+SlpPvOZC/cO/tkV3rx9Xq+6x9e/AAAA//8913egjwIAAA==
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},

	"/assets": {
		isDir: true,
		local: "/assets",
	},

	"/assets/css": {
		isDir: true,
		local: "/assets/css",
	},

	"/assets/script": {
		isDir: true,
		local: "/assets/script",
	},
}
