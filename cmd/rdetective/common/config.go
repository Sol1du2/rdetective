package common

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Generic settings.
	LogTimestamp bool
	LogLevel     string

	OriginalFilePath string
	UpdatedFilePath  string

	ChunkSize int
)

func SetDefaults(cmd *cobra.Command) {
	// Defaults
	viper.SetDefault("LOG_TIMESTAMP", true)
	viper.SetDefault("LOG_LEVEL", "info")

	// Command line flags
	cmd.Flags().Bool("log-timestamp", true, "prefix each log line with timestamp")
	cmd.Flags().String("log-level", "info", "log level (one of panic, fatal, error, warn, info or debug)")

	cmd.Flags().String("original", "", "original file")
	cmd.Flags().String("updated", "", "updated file")

	cmd.Flags().Int("chunk-size", 2, "the size of each hashed chunk (window)")

	_ = viper.BindPFlag("LOG_TIMESTAMP", cmd.Flags().Lookup("log-timestamp"))
	_ = viper.BindPFlag("LOG_LEVEL", cmd.Flags().Lookup("log-level"))

	_ = viper.BindPFlag("ORIGINAL", cmd.Flags().Lookup("original"))
	_ = viper.BindPFlag("UPDATED", cmd.Flags().Lookup("updated"))

	_ = viper.BindPFlag("CHUNK_SIZE", cmd.Flags().Lookup("chunk-size"))

	// Setup env.
	viper.SetEnvPrefix("rdetective")
	viper.AutomaticEnv()
}

func ApplyConfiguration() error {
	LogTimestamp = viper.GetBool("LOG_TIMESTAMP")
	LogLevel = viper.GetString("LOG_LEVEL")

	OriginalFilePath = viper.GetString("ORIGINAL")
	UpdatedFilePath = viper.GetString("UPDATED")

	ChunkSize = viper.GetInt("CHUNK_SIZE")

	return nil
}
