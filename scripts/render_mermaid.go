package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/render_mermaid.go <markdown-file> [output-type: png|svg]")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outType := "png"
	if len(os.Args) > 2 {
		outType = os.Args[2]
	}

	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var mermaidBlocks []string
	var currentBlock strings.Builder
	inBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "```mermaid") {
			inBlock = true
			continue
		}
		if inBlock && strings.HasPrefix(strings.TrimSpace(line), "```") {
			inBlock = false
			mermaidBlocks = append(mermaidBlocks, currentBlock.String())
			currentBlock.Reset()
			continue
		}
		if inBlock {
			currentBlock.WriteString(line + "\n")
		}
	}

	if len(mermaidBlocks) == 0 {
		fmt.Println("No mermaid blocks found in file.")
		return
	}

	for i, block := range mermaidBlocks {
		tmpFile := fmt.Sprintf("temp_diagram_%d.mmd", i)
		outName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
		if len(mermaidBlocks) > 1 {
			outName = fmt.Sprintf("%s_%d.%s", outName, i, outType)
		} else {
			outName = fmt.Sprintf("%s.%s", outName, outType)
		}

		err := os.WriteFile(tmpFile, []byte(block), 0644)
		if err != nil {
			fmt.Printf("Error writing temp file: %v\n", err)
			continue
		}

		fmt.Printf("Rendering %s -> %s...\n", inputPath, outName)
		cmd := exec.Command("npx", "-y", "@mermaid-js/mermaid-cli", "-i", tmpFile, "-o", outName, "-b", "transparent")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error rendering diagram %d: %v\nStderr: %s\n", i, err, stderr.String())
		} else {
			fmt.Printf("✅ Successfully saved: %s\n", outName)
		}

		os.Remove(tmpFile)
	}
}
