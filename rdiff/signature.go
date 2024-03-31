package rdiff

import "github.com/sol1du2/rdetective/rdiff/rhash"

// SignatureChunk represents a part of a file, along with its hashed value.
// Note(Sol1du2): For real use we should also use a strong hashing like sha1 to
// properly take care of collisions. And the weaker alder32 just for efficient
// lookups. However, for exercise purposes and to simplify the code we are only
// using adler32 here. Adding a stronger hash to the algorithm should be fairly
// trivial.
type SignatureChunk struct {
	Adler32 uint32
	Window  []byte
}

// Signature represents a file consisting of several chunks.
// indexMap represents the index (position) of each chunk with the hash value as
// the key. This is to help find matching chunks.
type Signature struct {
	Chunks   []SignatureChunk
	indexMap map[uint32][]int
}

func (s *Signature) AddChunk(chunk []byte) {
	a, w := generateHash(chunk)
	sc := SignatureChunk{
		Adler32: a,
		Window:  w,
	}

	s.Chunks = append(s.Chunks, sc)

	index := len(s.Chunks) - 1

	if existingValue, ok := s.indexMap[sc.Adler32]; ok {
		s.indexMap[sc.Adler32] = append(existingValue, index)
	} else {
		s.indexMap[sc.Adler32] = []int{index}
	}
}

func (s *Signature) MatchChunk(hash uint32) int {
	if indexes, ok := s.indexMap[hash]; ok {
		index := indexes[0]

		// Remove Matched Chunk
		if len(indexes) == 1 {
			delete(s.indexMap, hash)
		} else {
			indexes = indexes[1:]
			s.indexMap[hash] = indexes
		}

		return index
	}

	return -1
}

func generateHash(data []byte) (hash uint32, window []byte) {
	r := rhash.New()

	for _, b := range data {
		r.Update(b)
	}

	return r.Sum(), r.Window
}
