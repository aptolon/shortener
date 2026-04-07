package codec

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	base     = uint64(len(alphabet))
	codeLen  = 10
)

func Encode(id uint64) string {
	buf := make([]byte, codeLen)

	for i := codeLen - 1; i >= 0; i-- {
		buf[i] = alphabet[id%base]
		id /= base
	}

	return string(buf)
}

func Decode(code string) (uint64, error) {
	return 0, nil
}
