package diff

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/sol1du2/rdetective/cmd/rdetective/common"
	"github.com/sol1du2/rdetective/rdiff"
)

func CommandDiff() *cobra.Command {
	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Computes difference",
		Run: func(_ *cobra.Command, _ []string) {
			if err := diff(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	common.SetDefaults(diffCmd)

	return diffCmd
}

func diff() error {
	if err := common.ApplyConfiguration(); err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}

	logger, err := newLogger(!common.LogTimestamp, common.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Debugln("chunk size ", common.ChunkSize)
	logger.Debugln("original file ", common.OriginalFilePath)
	logger.Debugln("updated file ", common.UpdatedFilePath)
	logger.Debugln("diff start")

	rd, err := rdiff.New(&rdiff.Config{
		Logger:    logger,
		ChunkSize: common.ChunkSize,

		OriginalSource: FileSource{fileName: common.OriginalFilePath},
		UpdatedSource:  FileSource{fileName: common.UpdatedFilePath},
	})
	if err != nil {
		return fmt.Errorf("failed to create rolling diff: %w", err)
	}

	signature, err := rd.GenerateSignature()
	if err != nil {
		return fmt.Errorf("failed to generate signature: %w", err)
	}

	logger.Info("\n---signature---")
	logger.Debugln(signature)
	for i, s := range signature.Chunks {
		logger.Info("chunk ", i, ", hash ", s.Adler32, ", bytes ", s.Window)
	}

	delta, err := rd.GenerateDelta()
	if err != nil {
		return fmt.Errorf("failed to generate delta: %w", err)
	}

	// Sort it by position for easy printing.
	sort.SliceStable(delta.Changes, func(i, j int) bool {
		return delta.Changes[i].Position < delta.Changes[j].Position
	})

	logger.Info("\n---delta---")
	logger.Debugln(delta)
	for _, v := range delta.Changes {
		if len(v.NewBytes) > 0 {
			if v.ChunkIndex < len(signature.Chunks) {
				logger.Info("chunk ", v.ChunkIndex, " is at position ", v.Position, " and has new bytes prepended: ", string(v.NewBytes))
			} else {
				logger.Info("new chuck at end of file (position ", v.Position, ") with bytes: ", string(v.NewBytes))
			}
		} else {
			logger.Info("chunk ", v.ChunkIndex, " is at position ", v.Position)
		}
	}

	for _, index := range delta.MissingChunks {
		logger.Info("chunk ", index, " is missing")
	}

	return nil
}
