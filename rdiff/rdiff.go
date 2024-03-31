package rdiff

import (
	"bufio"
	"io"

	"github.com/sol1du2/rdetective/rdiff/rhash"
)

type DataSource interface {
	GetReader() (io.Reader, error)
}

type RollingDiff struct {
	config *Config

	originalBuffer *bufio.Reader
	updatedBuffer  *bufio.Reader

	signature Signature
}

func New(config *Config) (*RollingDiff, error) {
	rd := RollingDiff{
		config: config,
	}

	if err := rd.InitDataReaders(); err != nil {
		return nil, err
	}

	return &rd, nil
}

func (rd *RollingDiff) InitDataReaders() (err error) {
	originalBuffer, err := rd.config.OriginalSource.GetReader()
	if err != nil {
		return err
	}

	updatedBuffer, err := rd.config.UpdatedSource.GetReader()
	if err != nil {
		return err
	}

	rd.originalBuffer = bufio.NewReader(originalBuffer)
	rd.updatedBuffer = bufio.NewReader(updatedBuffer)

	return nil
}

func (rd *RollingDiff) GenerateSignature() (Signature, error) {
	chunkSize := rd.config.ChunkSize
	reader := rd.originalBuffer

	chunkData := make([]byte, chunkSize)

	rd.signature = Signature{
		Chunks:   make([]SignatureChunk, 0),
		indexMap: make(map[uint32][]int),
	}

	for {
		read, err := reader.Read(chunkData)

		if read == 0 || err == io.EOF {
			break
		}

		if err != nil {
			return rd.signature, err
		}

		if read < chunkSize {
			chunkData = chunkData[:read]
		}

		rd.signature.AddChunk(chunkData)
	}

	return rd.signature, nil
}

func (rd *RollingDiff) GenerateDelta() (Delta, error) {
	chunkSize := rd.config.ChunkSize
	reader := rd.updatedBuffer
	sig := rd.signature

	delta := Delta{
		Changes: []DeltaChunk{},
	}

	adler32 := rhash.New()

	var newBytes []byte
	newBytesLen := 0
	i := 0
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		if err != nil {
			return delta, err
		}

		adler32.Update(b)

		if adler32.Size < chunkSize {
			continue // Continue until chunk is full or we reached EOF.
		}

		if adler32.Size > chunkSize {
			removed, rollErr := adler32.Roll()
			if rollErr != nil {
				break
			}

			newBytes = append(newBytes, removed)
		}

		// Check match with signature.
		index := sig.MatchChunk(adler32.Sum())
		if index >= 0 {
			delta.Changes = append(delta.Changes, DeltaChunk{
				ChunkIndex: index,
				NewBytes:   newBytes,
				Position:   i*chunkSize + newBytesLen,
			})

			newBytesLen += len(newBytes)

			newBytes = []byte{}
			adler32.Reset()

			i++
		}
	}

	if adler32.Size < chunkSize { // Try last chunk if it's smaller than size.
		index := sig.MatchChunk(adler32.Sum())
		if index >= 0 {
			delta.Changes = append(delta.Changes, DeltaChunk{
				ChunkIndex: index,
				NewBytes:   newBytes,
				Position:   i*chunkSize + newBytesLen,
			})
		}
	} else if len(newBytes) > 0 { // Add data that is detected at the end of the file.
		delta.Changes = append(delta.Changes, DeltaChunk{
			ChunkIndex: len(sig.Chunks), // New index
			NewBytes:   append(newBytes, adler32.Window...),
			Position:   i*chunkSize + newBytesLen,
		})
	}

	// Store missing chunks.
	// Note(sol1du2): We could potentially just compare the delta with the
	// signature for the missing chunks. But this makes the result a bit nicer
	// to parse.
	for _, indexes := range sig.indexMap {
		delta.MissingChunks = append(delta.MissingChunks, indexes...)
	}

	return delta, nil
}
