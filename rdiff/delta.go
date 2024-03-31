package rdiff

// DeltaChunk represents a part of the file that maybe or may not have changed.
// ChunkIndex represents the index of this chunk relative to the Signature.
// NewBytes represents new bytes that are prepended to this chunk relative to
// the same indexed Signature chunk.
// Position represents the file position this chunk starts. It differs from the
// ChunkIndex as this is the actual file position and not relative to the
// Signature.
type DeltaChunk struct {
	ChunkIndex int
	NewBytes   []byte
	Position   int
}

// Delta represents all changed chunks relative to the Signature.
// If len(Changes) > len(Signature) that means new chunks were added at the end.
// MissingChunks contains all chunks that are not detectable anymore from the
// Signature of the original file.
type Delta struct {
	Changes       []DeltaChunk
	MissingChunks []int
}
