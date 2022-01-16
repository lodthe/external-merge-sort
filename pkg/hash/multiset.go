package hash

// MultiSet calculates hash for multiset of integers.
// hash(a1, a2, a3) = (x^a1 + x^a2 + x^a3) % mod...
type MultiSet struct {
	x   int64
	mod int64

	hash int64
}

func NewMultiset(x, mod int64) *MultiSet {
	const seed = 3912929

	return &MultiSet{
		x:    x,
		mod:  mod,
		hash: seed % mod,
	}
}

func (s *MultiSet) Hash() int64 {
	return s.hash
}

func (s *MultiSet) Add(value int64) {
	s.hash = (s.hash + s.binpow(s.x, value)) % s.mod
}

// binpow calculates a^power using binary exponentiation.
func (s *MultiSet) binpow(a, power int64) int64 {
	if power == 0 {
		return 1
	}

	tmp := s.binpow(a, power/2)
	tmp = tmp * tmp % s.mod

	if power%2 == 1 {
		return tmp * a % s.mod
	}

	return tmp
}
