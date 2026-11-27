// Topic: readers-writers | Concept: io | Stage: 03-advanced
// Run:   go run ./03-advanced/io/readers-writers
//
// Explore:
// * io.Reader (Read(p []byte)) and io.Writer (Write(p []byte)) are tiny interfaces
// * Files, network connections, buffers and strings all implement them
// * Compose them: io.Copy, io.MultiWriter, io.TeeReader, io.LimitReader
// * bufio.Scanner reads line by line; bufio.Writer batches writes (Flush!)
// * Writing your own Reader/Writer is just one method
// * os.CreateTemp + defer os.Remove for scratch files
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// upperWriter wraps another writer and upper-cases everything
type upperWriter struct {
	w io.Writer
}

func (u upperWriter) Write(p []byte) (int, error) {
	return u.w.Write(bytes.ToUpper(p))
}

// countingReader counts letters as data passes through
type countingReader struct {
	r       io.Reader
	letters int
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	for _, r := range string(p[:n]) {
		if unicode.IsLetter(r) {
			c.letters++
		}
	}
	return n, err
}

func main() {
	fmt.Println("Topic: io.Reader and io.Writer")

	text := "line one\nline two\nline three\n"

	// Scanner over a strings.Reader
	sc := bufio.NewScanner(strings.NewReader(text))
	n := 0
	for sc.Scan() {
		n++
		fmt.Printf("  scanned %d: %s\n", n, sc.Text())
	}

	// Custom writer + io.Copy
	fmt.Print("  upper: ")
	_, _ = io.Copy(upperWriter{os.Stdout}, strings.NewReader("shout this\n"))

	// Custom reader
	cr := &countingReader{r: strings.NewReader("Go 1.26 rocks!")}
	data, _ := io.ReadAll(cr)
	fmt.Printf("  read %q, letters counted = %d\n", data, cr.letters)

	// MultiWriter: one write, two destinations
	var buf bytes.Buffer
	mw := io.MultiWriter(&buf, upperWriter{io.Discard})
	fmt.Fprintf(mw, "written to %d writers", 2)
	fmt.Println("  buffer got:", buf.String())

	// LimitReader
	limited, _ := io.ReadAll(io.LimitReader(strings.NewReader("0123456789"), 4))
	fmt.Println("  LimitReader(4):", string(limited))

	// Files are Readers and Writers too
	f, err := os.CreateTemp("", "go-io-lesson-*.txt")
	if err != nil {
		fmt.Println("temp file error:", err)
		return
	}
	defer os.Remove(f.Name())
	bw := bufio.NewWriter(f)
	for i := 1; i <= 3; i++ {
		fmt.Fprintf(bw, "row %d\n", i)
	}
	bw.Flush() // without Flush the data stays in memory
	f.Close()

	content, _ := os.ReadFile(f.Name())
	fmt.Printf("  temp file content: %q\n", content)
}
