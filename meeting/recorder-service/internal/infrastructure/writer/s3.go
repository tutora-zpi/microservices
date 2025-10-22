package writer

import "github.com/pion/rtp"

type s3Writer struct{}

// Close implements Writer.
func (s *s3Writer) Close() error {
	panic("unimplemented")
}

// Write implements Writer.
func (s *s3Writer) Write(packet *rtp.Packet) error {
	panic("unimplemented")
}

func NewS3Writer() Writer {
	return &s3Writer{}
}

func (s *s3Writer) GetPath() string {
	// S3 does not have a local path, return an empty string or a placeholder.
	return ""
}
