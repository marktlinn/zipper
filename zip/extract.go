package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log/slog"
)

const chunkSize = 2 * 1024 * 1024 // 2MB

// Decompress processes a ZIP file from a reader and writes the contents to the writer.
func Decompress(r io.Reader, w io.Writer) error {
	for {
		buf, n, err := readChunk(r)
		if err != nil {
			if err != io.EOF {
				slog.Error("error reading from input", "error", err)
				return fmt.Errorf("failed to read from input: %w", err)
			}
			if n == 0 { // End of file and no data read
				break
			}
		}

		if err := processChunk(buf[:n], w); err != nil {
			return err
		}

		if err == io.EOF && n == 0 { // End of file with no more data to process
			break
		}
	}

	return nil
}

// readChunk reads a chunk of data from the reader.
func readChunk(r io.Reader) ([]byte, int, error) {
	buf := make([]byte, chunkSize)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		slog.Error("error reading from input", "error", err)
		return nil, 0, err
	}
	return buf, n, err
}

// processChunk processes a chunk of data and writes its contents to the writer.
func processChunk(chunk []byte, w io.Writer) error {
	zipReader, err := zip.NewReader(bytes.NewReader(chunk), int64(len(chunk)))
	if err != nil {
		slog.Error("error creating zip reader", "error", err)
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	for _, f := range zipReader.File {
		if err := processFile(f, w); err != nil {
			return err
		}
	}

	return nil
}

// processFile processes a single file from the ZIP archive.
func processFile(f *zip.File, w io.Writer) error {
	file, err := f.Open()
	fmt.Printf("file name: %s\n", f.Name)
	if err != nil {
		slog.Error("error opening file in zip archive", "error", err)
		return fmt.Errorf("failed to open file in zip archive: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		slog.Error("error writing decompressed data", "error", err)
		return fmt.Errorf("failed to write decompressed data: %w", err)
	}

	return nil
}
