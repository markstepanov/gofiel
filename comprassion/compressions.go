package comprassion

import (
	"bytes"
	"compress/gzip"
	"errors"

	"github.com/klauspost/compress/zstd"
)




func CompressBytes(data *[]byte, comressionType string) (*[]byte, error) {
    switch comressionType {
        case ComressionGzip: {
            return compressGzip(data)
        }
        case ComressionZstd: {
            return compressZstd(data)
        }
    }

    return nil, errors.New("unknown compression algorithm")
}



func compressGzip(data *[]byte) (*[]byte, error) {
    var buf bytes.Buffer
    gw := gzip.NewWriter(&buf)

    defer gw.Close()

    _, err := gw.Write(*data)
    if err != nil {
        return nil, err
    }

    err = gw.Close() 
    if err != nil {
        return nil, err
    }

    resp := buf.Bytes()
    return &resp, nil
}


func compressZstd(data *[]byte) (*[]byte, error) {
    var buff []byte = make([]byte, 0, len(*data))

	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}
	defer encoder.Close()

	compressed := encoder.EncodeAll(*data, buff)
	return &compressed, nil

}


// func decompressZstd(data *[]byte)  (*[] byte, error){
//     // todo should pass 
//     decoder , err := zstd.NewReader(nil)

//     if err != nil {
//         return  nil, err
//     }
//     defer decoder.Close()

//     decoder.DecodeAll(data, )





// }
