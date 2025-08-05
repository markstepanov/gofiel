package comprassion

type ComprassionInfo struct {
	ComprassionAlgorithm string
	UncompressedSize     int
	CompressedSize       int
}

var ComressionGzip = "Gzip"
var ComressionZstd = "Zstd"
var ComressionNone = "non-compressed"

var AvailableCompressions = [3]string{ComressionGzip, ComressionZstd, ComressionNone}
