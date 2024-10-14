package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Constants for configuration
const (
	mdExtension     = ".md"
	htmlExtension   = ".html"
	dirPermissions  = 0755
	filePermissions = 0644
	indexPageName   = "index.html"
)

// mdToHTML converts markdown content to HTML
func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// readFile reads the content of a file
func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// writeFile writes data to a file
func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, filePermissions)
}

// generateHTMLPage creates a complete HTML page from markdown content
func generateHTMLPage(cssPath string, content string, backLink string) string {
	template, err := readFile("template.html")
	if err != nil {
		log.Fatalf("Failed to read HTML template: %v", err)
	}
	return fmt.Sprintf(string(template), cssPath, content, backLink)
}

// createDirectory creates a directory if it doesn't exist
func createDirectory(path string) error {
	if err := os.MkdirAll(path, dirPermissions); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", path, err)
	}
	return nil
}

// isMarkdownFile checks if a file has a .md extension
func isMarkdownFile(filename string) bool {
	return strings.HasSuffix(filename, mdExtension)
}

// processFile handles the conversion of a Markdown file to HTML
func processFile(srcPath, destPath string) error {
	mdData, err := readFile(srcPath)
	if err != nil {
		return fmt.Errorf("error reading markdown file %s: %v", srcPath, err)
	}

	htmlData := mdToHTML(mdData)

	// Calculate the relative path to the stylesheet
	depth := strings.Count(filepath.Dir(destPath), string(os.PathSeparator))
	cssPath := strings.Repeat("../", depth) + "style.css"
	backLink := strings.Repeat("../", depth) + "index.html"

	htmlContent := generateHTMLPage(cssPath, string(htmlData), backLink)

	destPath = strings.TrimSuffix(destPath, mdExtension) + htmlExtension
	if err := writeFile(destPath, []byte(htmlContent)); err != nil {
		return fmt.Errorf("error writing HTML file %s: %v", destPath, err)
	}

	log.Printf("Converted %s to %s", srcPath, destPath)
	return nil
}

func generateIndexPage(destDir string) error {
	var links []string

	// Traverse the destination directory to find all HTML files
	err := filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), htmlExtension) {
			// Skip the index.html file
			if info.Name() == indexPageName {
				return nil
			}
			// Create a relative link to the HTML file
			relPath, _ := filepath.Rel(destDir, path)
			links = append(links, relPath)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Generate the index HTML content without a title and bullets
	var indexContent strings.Builder
	indexContent.WriteString("<!DOCTYPE html>\n<html>\n<head>\n <meta charset=\"UTF-8\">\n<title>Notes</title\n<link rel=\"preconnect\" href=\"https://fonts.googleapis.com\">\n<link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin>\n<link href=\"https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap\" rel=\"stylesheet\">\n<link rel=\"apple-touch-icon\" sizes=\"180x180\" href=\"/apple-touch-icon.png\">\n<link rel=\"icon\" type=\"image/png\" sizes=\"32x32\" href=\"/favicon-32x32.png\">\n<link rel=\"icon\" type=\"image/png\" sizes=\"16x16\" href=\"/favicon-16x16.png\">\n<link rel=\"manifest\" href=\"/site.webmanifest\">\n<link rel=\"stylesheet\" type=\"text/css\" href=\"style.css\">\n</head>\n<body>\n<div class=\"wrapper\">\n")
	for _, link := range links {
		// Remove the .html extension for display
		nameWithoutExtension := strings.TrimSuffix(link, ".html")
		indexContent.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a><br>\n", link, nameWithoutExtension))
	}
	indexContent.WriteString("</div></body>\n</html>")

	// Write the index file
	indexPath := filepath.Join(destDir, indexPageName)
	return writeFile(indexPath, []byte(indexContent.String()))
}

// processDirectory processes the source directory and mirrors the structure to the destination directory
func processDirectory(srcDir, destDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}

		// Determine the relative path and create the new path
		relPath, _ := filepath.Rel(srcDir, path)
		newPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create directory in the destination
			if err := createDirectory(newPath); err != nil {
				return err
			}
			log.Printf("Created directory: %s", newPath)
			return nil
		}

		// Process Markdown files
		if isMarkdownFile(info.Name()) {
			return processFile(path, newPath)
		}

		log.Printf("Skipped non-Markdown file: %s", path)
		return nil
	})
}

func copyFileToPublic(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", path, err)
	}

	newPath := "public/" + path

	err = writeFile(newPath, file)
	if err != nil {
		return fmt.Errorf("error writing to public/%s: %v", path, err)
	}

	return nil
}

func main() {
	srcDir := "md"      // Source directory containing .md files
	destDir := "public" // Destination directory for HTML files

	if err := processDirectory(srcDir, destDir); err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}

	if err := generateIndexPage(destDir); err != nil {
		log.Fatalf("Error generating index page: %v", err)
	}

	if err := copyFileToPublic("style.css"); err != nil {
		log.Fatalf("Error copying stylesheet: %v", err)
	}

	if err := copyFileToPublic("site.webmanifest"); err != nil {
		log.Fatalf("Error copying site.webmanifest: %v", err)
	}

	if err := copyFileToPublic("favicon.ico"); err != nil {
		log.Fatalf("Error copying favicon: %v", err)
	}

	if err := copyFileToPublic("favicon-32x32.png"); err != nil {
		log.Fatalf("Error copying favicon: %v", err)
	}

	if err := copyFileToPublic("favicon-16x16.png"); err != nil {
		log.Fatalf("Error copying favicon: %v", err)
	}

	if err := copyFileToPublic("apple-touch-icon.png"); err != nil {
		log.Fatalf("Error copying favicon: %v", err)
	}

	if err := copyFileToPublic("apple-touch-icon.png"); err != nil {
		log.Fatalf("Error copying favicon: %v", err)
	}

	if err := copyFileToPublic("android-chrome-192x192.png"); err != nil {
		log.Fatalf("Error copying favicon: %v", err)
	}
	fmt.Println("Markdown files converted to HTML successfully, and index page created.")
}
