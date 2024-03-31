package rdiff

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

type StringSource struct {
	Data string
}

func (s StringSource) GetReader() (io.Reader, error) {
	return strings.NewReader(s.Data), nil
}

func getRDiff(original, updated string) (*RollingDiff, error) {
	lvl, _ := logrus.ParseLevel("debug")
	rh, err := New(&Config{
		ChunkSize:      2,
		OriginalSource: StringSource{Data: original},
		UpdatedSource:  StringSource{Data: updated},
		Logger: &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				DisableTimestamp: false,
			},
			Level: lvl,
		},
	})

	if err != nil {
		return nil, err
	}

	return rh, nil
}

func compareSignatures(sig1, sig2 Signature, t *testing.T) {
	if len(sig1.Chunks) != len(sig2.Chunks) {
		t.Errorf("unexpected length of chunks, got %d, expected %d", len(sig1.Chunks), len(sig2.Chunks))
	}

	if len(sig1.indexMap) != len(sig2.indexMap) {
		t.Errorf("unexpected length of indexed map, got %d, expected %d", len(sig1.indexMap), len(sig2.indexMap))
	}

	for i := range sig1.Chunks {
		chunk1 := sig1.Chunks[i]
		chunk2 := sig2.Chunks[i]

		if chunk1.Adler32 != chunk2.Adler32 {
			t.Errorf("hash value of chunk %d different than expected, got %d, expected %d", i, chunk1.Adler32, chunk2.Adler32)
		}

		if !bytes.Equal(chunk1.Window, chunk2.Window) {
			t.Errorf("hash value of chunk %d different than expected, got %s, expected %s",
				i, string(chunk1.Window), string(chunk2.Window))
		}

		indexed1, ok := sig1.indexMap[chunk1.Adler32]
		if !ok {
			t.Errorf("indexed chunk not found, expected %d", chunk1.Adler32)
		}

		indexed2, ok := sig2.indexMap[chunk1.Adler32]
		if !ok {
			t.Errorf("unexpected chunk %d", chunk1.Adler32)
		}

		if len(indexed1) != len(indexed2) {
			t.Errorf("unexpected length of indexed chunk, got %d, expected %d", len(indexed1), len(indexed2))
		}

		for j := range indexed1 {
			if indexed1[j] != indexed2[j] {
				t.Errorf("unexpected index, got %d, expected %d", indexed1[j], indexed2[j])
			}
		}
	}
}

func compareDeltas(delta1, delta2 Delta, t *testing.T) {
	if len(delta1.Changes) != len(delta2.Changes) {
		t.Errorf("unexpected length of changes, got %d, expected %d", len(delta1.Changes), len(delta2.Changes))
	}

	if len(delta1.MissingChunks) != len(delta2.MissingChunks) {
		t.Errorf("unexpected length of missing chunks, got %d, expected %d", len(delta1.MissingChunks), len(delta2.MissingChunks))
	}

	for i := range delta1.Changes {
		change1 := delta1.Changes[i]
		change2 := delta2.Changes[i]

		if change1.ChunkIndex != change2.ChunkIndex {
			t.Errorf("unexpected chunk index, got %d, expected %d", change1.ChunkIndex, change2.ChunkIndex)
		}

		if !bytes.Equal(change1.NewBytes, change2.NewBytes) {
			t.Errorf("unexpected new bytes, got %s, expected %s", string(change1.NewBytes), string(change2.NewBytes))
		}

		if change1.Position != change2.Position {
			t.Errorf("unexpected position, got %d, expected %d", change1.Position, change2.Position)
		}
	}

	for i := range delta1.MissingChunks {
		if delta1.MissingChunks[i] != delta2.MissingChunks[i] {
			t.Errorf("unexpected missing chunk, got %d, expected %d", delta1.MissingChunks[i], delta2.MissingChunks[i])
		}
	}
}

func TestSignature(t *testing.T) {
	original := "hello"
	updated := "not part of test"
	expectedS := Signature{
		Chunks: []SignatureChunk{
			{
				Adler32: 20381902,
				Window:  []byte("he"),
			},
			{
				Adler32: 21364953,
				Window:  []byte("ll"),
			},
			{
				Adler32: 7340144,
				Window:  []byte("o"),
			},
		},
		indexMap: map[uint32][]int{
			20381902: {
				0,
			},
			21364953: {
				1,
			},
			7340144: {
				2,
			},
		},
	}

	rh, err := getRDiff(original, updated)
	if err != nil {
		t.Errorf("error creating rdiff %s", err.Error())
	}

	s, err := rh.GenerateSignature()
	if err != nil {
		t.Errorf("error generating signature: %s", err.Error())
	}

	compareSignatures(s, expectedS, t)
}

func TestChunkChanged(t *testing.T) {
	original := "hello"
	updated := "heeello"
	expectedDelta := Delta{
		Changes: []DeltaChunk{
			{
				ChunkIndex: 0,
				NewBytes:   []byte{},
				Position:   0,
			},
			{
				ChunkIndex: 1,
				NewBytes:   []byte{'e', 'e'},
				Position:   2,
			},
			{
				ChunkIndex: 2,
				NewBytes:   []byte{},
				Position:   6,
			},
		},
		MissingChunks: []int{},
	}

	rh, err := getRDiff(original, updated)
	if err != nil {
		t.Errorf("error creating rdiff %s", err.Error())
	}

	_, err = rh.GenerateSignature()
	if err != nil {
		t.Errorf("error generating signature: %s", err.Error())
	}

	delta, err := rh.GenerateDelta()
	if err != nil {
		t.Errorf("error generating delta: %s", err.Error())
	}

	compareDeltas(delta, expectedDelta, t)
}

func TestChunkMoved(t *testing.T) {
	original := "hello"
	updated := "llohe"
	expectedDelta := Delta{
		Changes: []DeltaChunk{
			{
				ChunkIndex: 1,
				NewBytes:   []byte{},
				Position:   0,
			},
			{
				ChunkIndex: 0,
				NewBytes:   []byte{'o'},
				Position:   2,
			},
		},
		MissingChunks: []int{2},
	}

	rh, err := getRDiff(original, updated)
	if err != nil {
		t.Errorf("error creating rdiff %s", err.Error())
	}

	_, err = rh.GenerateSignature()
	if err != nil {
		t.Errorf("error generating signature: %s", err.Error())
	}

	delta, err := rh.GenerateDelta()
	if err != nil {
		t.Errorf("error generating delta: %s", err.Error())
	}

	compareDeltas(delta, expectedDelta, t)
}

func TestEqualChunks(t *testing.T) {
	original := "hellllllo"
	updated := "hellllo"
	expectedDelta := Delta{
		Changes: []DeltaChunk{
			{
				ChunkIndex: 0,
				NewBytes:   []byte{},
				Position:   0,
			},
			{
				ChunkIndex: 1,
				NewBytes:   []byte{},
				Position:   2,
			},
			{
				ChunkIndex: 2,
				NewBytes:   []byte{},
				Position:   4,
			},
			{
				ChunkIndex: 4,
				NewBytes:   []byte{},
				Position:   6,
			},
		},
		MissingChunks: []int{3},
	}

	rh, err := getRDiff(original, updated)
	if err != nil {
		t.Errorf("error creating rdiff %s", err.Error())
	}

	_, err = rh.GenerateSignature()
	if err != nil {
		t.Errorf("error generating signature: %s", err.Error())
	}

	delta, err := rh.GenerateDelta()
	if err != nil {
		t.Errorf("error generating delta: %s", err.Error())
	}

	compareDeltas(delta, expectedDelta, t)
}

func TestRemovedChunks(t *testing.T) {
	original := "hello"
	updated := "heo"
	expectedDelta := Delta{
		Changes: []DeltaChunk{
			{
				ChunkIndex: 0,
				NewBytes:   []byte{},
				Position:   0,
			},
			{
				ChunkIndex: 2,
				NewBytes:   []byte{},
				Position:   2,
			},
		},
		MissingChunks: []int{1},
	}

	rh, err := getRDiff(original, updated)
	if err != nil {
		t.Errorf("error creating rdiff %s", err.Error())
	}

	_, err = rh.GenerateSignature()
	if err != nil {
		t.Errorf("error generating signature: %s", err.Error())
	}

	delta, err := rh.GenerateDelta()
	if err != nil {
		t.Errorf("error generating delta: %s", err.Error())
	}

	compareDeltas(delta, expectedDelta, t)
}

func TestAddChunksEnd(t *testing.T) {
	original := "hello"
	updated := "heeello world"
	expectedDelta := Delta{
		Changes: []DeltaChunk{
			{
				ChunkIndex: 0,
				NewBytes:   []byte{},
				Position:   0,
			},
			{
				ChunkIndex: 1,
				NewBytes:   []byte{'e', 'e'},
				Position:   2,
			},
			{
				ChunkIndex: 3,
				NewBytes:   []byte{'o', ' ', 'w', 'o', 'r', 'l', 'd'},
				Position:   6,
			},
		},
		MissingChunks: []int{2},
	}

	rh, err := getRDiff(original, updated)
	if err != nil {
		t.Errorf("error creating rdiff %s", err.Error())
	}

	_, err = rh.GenerateSignature()
	if err != nil {
		t.Errorf("error generating signature: %s", err.Error())
	}

	delta, err := rh.GenerateDelta()
	if err != nil {
		t.Errorf("error generating delta: %s", err.Error())
	}

	compareDeltas(delta, expectedDelta, t)
}
