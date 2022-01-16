package hash

func Polynomial(data []byte) int64 {
	const p = 327
	const mod = 1e9 + 123
	const seed = 1932932

	h := int64(seed ^ len(data))
	for _, b := range data {
		h = (h*p + int64(b)) % mod
	}

	return h
}
