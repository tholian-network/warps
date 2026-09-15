package http

type Encoding string

const (
	EncodingIdentity Encoding = "identity"
	EncodingBzip2    Encoding = "bzip2"
	EncodingDeflate  Encoding = "deflate"  // RFC 1951
	EncodingGzip     Encoding = "gzip"     // RFC 1952
	EncodingZstd     Encoding = "zstd"     // RFC 8878
)

func (encoding Encoding) String() string {

	switch encoding {
	case EncodingIdentity:
		return "identity"
	case EncodingBzip2:
		return "bzip2"
	case EncodingDeflate:
		return "deflate"
	case EncodingGzip:
		return "gzip"
	case EncodingZstd:
		return "zstd"
	default:
		return ""
	}

}
