package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/marktlinn/zipper/zip"
)

func main() {
	// Define the flags
	filePath := flag.String("f", "", "Path to the input file")
	outputPath := flag.String("o", "", "Path to the output file")
	compress := flag.Bool("c", false, "Performs a compression operation")
	decompress := flag.Bool("d", false, "Performs a decompression operation")

	// Parse the flags
	flag.Parse()

	noOperation := !*compress && !*decompress
	if *filePath == "" || noOperation || *outputPath == "" {
		log.Fatalf("Please provide valid file path, operation, and output path.")
	}

	// Open the input file
	inputFile, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer inputFile.Close()

	// Create the output file
	outFile, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Perform the operation

	if *compress {
		err = zip.Compress(inputFile, outFile, *filePath)
		if err != nil {
			slog.Error("Error compressing file", "error", err)
			log.Fatalf("Failed to compress file: %v", err)
		}
		slog.Info("Compression successful")

		return
	} else if *decompress {
		err = zip.Decompress(inputFile, outFile)
		if err != nil {
			slog.Error("Error decompressing file", "error", err)
			log.Fatalf("Failed to decompress file: %v", err)
		}
		slog.Info("Decompression successful")

		return
	}
	log.Fatal("Invalid operation: Use '-c' flag for compress or '-d' for decompress.")
}
