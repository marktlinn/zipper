package zip

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// Compress takes an io.Reader as input, an io.Writer as output, and a filename string.
// It reads the data from the io.Reader in chunks, compresses it using the Deflate
// method, and writes the compressed data to the io.Writer. The compressed data is
// written to a file within the zip archive with the specified filename.
//
// This function efficiently handles large files by streaming data without reading
// the entire file into memory.
func Compress(r io.Reader, w io.Writer, filename string) error {
	zipWriter := zip.NewWriter(w)

	defer zipWriter.Close()

	header := &zip.FileHeader{
		Name:     filename,
		Method:   zip.Deflate,
		Modified: time.Now(),
	}

	f, err := zipWriter.CreateHeader(header)
	if err != nil {
		slog.Error("error creating header in Zip Compress")
		return fmt.Errorf("failed to createheader: %w", err)
	}

	srcBuf := bufio.NewReader(r)
	_, err = io.CopyBuffer(f, srcBuf, make([]byte, 2048))
	if err != nil {
		slog.Error("error writing buffer from src to dest")
		return fmt.Errorf("failed to copy from buffer src to destination: %w", err)
	}

	return nil
}
