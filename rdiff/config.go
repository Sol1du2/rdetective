package rdiff

import "github.com/sirupsen/logrus"

type Config struct {
	Logger    logrus.FieldLogger
	ChunkSize int

	OriginalSource DataSource
	UpdatedSource  DataSource
}
