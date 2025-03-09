package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	mdExtension     = ".md"
	htmlExtension   = ".html"
	dirPermissions  = 0755
	filePermissions = 0644
	indexPageName   = "index.html"
	publicDir       = "public/"
)

type NavItem struct {
	Name     string
	Link     string
	IsDir    bool
	Children []*NavItem
}

func formatNavName(name string) string {
	name = strings.ReplaceAll(name, "_", " ")

	words := strings.Fields(name)

	for i := range words {
		if len(words[i]) > 0 {
			isAcronym := true
			for _, char := range words[i] {
				if !unicode.IsUpper(char) && !unicode.IsDigit(char) && char != '.' {
					isAcronym = false
					break
				}
			}

			if !isAcronym {
				r := []rune(words[i])
				r[0] = unicode.ToUpper(r[0])
				for j := 1; j < len(r); j++ {
					r[j] = unicode.ToLower(r[j])
				}
				words[i] = string(r)
			}
		}
	}

	return strings.Join(words, " ")
}

func extractCategoryFromPath(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))

	if len(parts) >= 2 {
		categoryParts := parts[1 : len(parts)-1]

		if len(categoryParts) == 0 {
			return ""
		}

		formattedParts := make([]string, len(categoryParts))
		for i, part := range categoryParts {
			if i == 0 {
				formattedParts[i] = strings.ToUpper(formatNavName(part))
			} else {
				formattedParts[i] = formatNavName(part)
			}
		}

		return strings.Join(formattedParts, " > ")
	}
	return ""
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, filePermissions)
}

func createDirectory(path string) error {
	if err := os.MkdirAll(path, dirPermissions); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", path, err)
	}
	return nil
}

func isMarkdownFile(filename string) bool {
	return strings.HasSuffix(filename, mdExtension)
}

func buildNavTree(srcDir string) (*NavItem, error) {
	root := &NavItem{
		Name:     filepath.Base(srcDir),
		Link:     "",
		IsDir:    true,
		Children: []*NavItem{},
	}

	categories := make(map[string]*NavItem)

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == srcDir {
			return nil
		}

		relPath, _ := filepath.Rel(srcDir, path)
		parts := strings.Split(relPath, string(os.PathSeparator))

		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() && !isMarkdownFile(path) {
			return nil
		}

		category := parts[0]
		var categoryItem *NavItem
		if item, exists := categories[category]; exists {
			categoryItem = item
		} else {
			categoryItem = &NavItem{
				Name:     category,
				Link:     "",
				IsDir:    true,
				Children: []*NavItem{},
			}
			categories[category] = categoryItem
			root.Children = append(root.Children, categoryItem)
		}

		if len(parts) == 1 && !info.IsDir() {
			htmlLink := strings.TrimSuffix(relPath, mdExtension) + htmlExtension
			fileItem := &NavItem{
				Name:  strings.TrimSuffix(filepath.Base(path), mdExtension),
				Link:  htmlLink,
				IsDir: false,
			}
			categoryItem.Children = append(categoryItem.Children, fileItem)
			return nil
		}

		if len(parts) > 1 {
			currentLevel := categoryItem

			for i := 1; i < len(parts); i++ {
				part := parts[i]

				if i == len(parts)-1 && !info.IsDir() {
					htmlLink := strings.TrimSuffix(relPath, mdExtension) + htmlExtension
					fileItem := &NavItem{
						Name:  strings.TrimSuffix(part, mdExtension),
						Link:  htmlLink,
						IsDir: false,
					}
					currentLevel.Children = append(currentLevel.Children, fileItem)
					break
				}

				if info.IsDir() || i < len(parts)-1 {
					found := false
					for _, child := range currentLevel.Children {
						if child.Name == part && child.IsDir {
							currentLevel = child
							found = true
							break
						}
					}

					if !found {
						newDir := &NavItem{
							Name:     part,
							Link:     "",
							IsDir:    true,
							Children: []*NavItem{},
						}
						currentLevel.Children = append(currentLevel.Children, newDir)
						currentLevel = newDir
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return root, nil
}

func renderNavTree(root *NavItem) string {
	var sb strings.Builder

	for _, category := range root.Children {
		formattedCategoryName := formatNavName(category.Name)
		sb.WriteString(fmt.Sprintf("<div class=\"category\">%s\n", formattedCategoryName))
		sb.WriteString("    <ul>\n")

		for _, item := range category.Children {
			if !item.IsDir {
				formattedItemName := formatNavName(item.Name)
				sb.WriteString(fmt.Sprintf("        <li><a href=\"/%s\">%s</a></li>\n", item.Link, formattedItemName))
			}
		}

		for _, subcategory := range category.Children {
			if !subcategory.IsDir {
				continue
			}

			formattedSubcategoryName := formatNavName(subcategory.Name)
			expandableClass := ""
			if len(subcategory.Children) > 0 {
				expandableClass = " class=\"expandable\""
			}
			sb.WriteString(fmt.Sprintf("        <li%s>%s\n", expandableClass, formattedSubcategoryName))

			if len(subcategory.Children) > 0 {
				sb.WriteString("            <ul>\n")

				for _, item := range subcategory.Children {
					if !item.IsDir {
						formattedItemName := formatNavName(item.Name)
						sb.WriteString(fmt.Sprintf("                <li><a href=\"/%s\">%s</a></li>\n", item.Link, formattedItemName))
					} else {
						formattedItemName := formatNavName(item.Name)
						expandableClass := ""
						if len(item.Children) > 0 {
							expandableClass = " class=\"expandable\""
						}
						sb.WriteString(fmt.Sprintf("                <li%s>%s\n", expandableClass, formattedItemName))

						if len(item.Children) > 0 {
							sb.WriteString("                    <ul>\n")
							for _, page := range item.Children {
								if page.IsDir {
									formattedPageName := formatNavName(page.Name)
									sb.WriteString(fmt.Sprintf("                        <li>%s</li>\n", formattedPageName))
								} else {
									formattedPageName := formatNavName(page.Name)
									sb.WriteString(fmt.Sprintf("                        <li><a href=\"/%s\">%s</a></li>\n", page.Link, formattedPageName))
								}
							}
							sb.WriteString("                    </ul>\n")
						}
						sb.WriteString("                </li>\n")
					}
				}

				sb.WriteString("            </ul>\n")
			}

			sb.WriteString("        </li>\n")
		}

		sb.WriteString("    </ul>\n")
		sb.WriteString("</div>\n")
	}

	return sb.String()
}

func processFile(srcPath, destPath, srcDir string, navContent string) error {
	mdData, err := readFile(srcPath)
	if err != nil {
		return fmt.Errorf("error reading markdown file %s: %v", srcPath, err)
	}

	htmlData := mdToHTML(mdData)

	category := extractCategoryFromPath(srcPath)

	categoryHTML := ""
	if category != "" {
		parts := strings.Split(category, " > ")
		var formattedParts []string

		for _, part := range parts {
			formattedParts = append(formattedParts, fmt.Sprintf(`<span class="category-part">%s</span>`, part))
		}

		categoryPath := strings.Join(formattedParts, `<span class="separator">›</span>`)
		categoryHTML = fmt.Sprintf(`<div class="content-category">%s</div>`, categoryPath)
	}

	htmlWithCategory := string(htmlData)
	if categoryHTML != "" {
		h1Regex := regexp.MustCompile(`(?i)(<h1[^>]*>)`)
		if h1Regex.MatchString(htmlWithCategory) {
			htmlWithCategory = h1Regex.ReplaceAllString(htmlWithCategory, categoryHTML+"$1")
		} else {
			htmlWithCategory = categoryHTML + htmlWithCategory
		}
	}

	htmlWithCategory = addHeaderIDs(htmlWithCategory)

	depth := strings.Count(filepath.Dir(destPath), string(os.PathSeparator))
	relativePathPrefix := strings.Repeat("../", depth)

	apptouchPath := relativePathPrefix + "apple-touch-icon.png"
	favicon32Path := relativePathPrefix + "favicon-32x32.png"
	favicon16Path := relativePathPrefix + "favicon-16x16.png"
	webmanifestPath := relativePathPrefix + "site.webmanifest"
	cssPath := relativePathPrefix + "style.css"
	logoPath := relativePathPrefix + "favicon-32x32.png" // Fix for logo path
	backLink := relativePathPrefix + "index.html"

	relPath, _ := filepath.Rel(srcDir, srcPath)
	currentPage := strings.TrimSuffix(relPath, mdExtension) + htmlExtension

	tocContent := generateTableOfContents(htmlWithCategory)

	htmlContent := generateHTMLPage(apptouchPath, favicon32Path, favicon16Path, webmanifestPath, cssPath, logoPath,
		htmlWithCategory, backLink, navContent, tocContent, currentPage)

	destPath = strings.TrimSuffix(destPath, mdExtension) + htmlExtension
	if err := writeFile(destPath, []byte(htmlContent)); err != nil {
		return fmt.Errorf("error writing HTML file %s: %v", destPath, err)
	}

	log.Printf("Converted %s to %s", srcPath, destPath)
	return nil
}

func generateHTMLPage(apptouchPath, favicon32Path, favicon16Path, webmanifesstPath, cssPath, logoPath string,
	content, backLink, navContent, tocContent, currentPage string) string {
	template, err := readFile("template.html")
	if err != nil {
		log.Fatalf("Failed to read HTML template: %v", err)
	}

	htmlContent := string(template)
	htmlContent = strings.Replace(htmlContent, "<body>", fmt.Sprintf("<body data-current-page=\"/%s\">", currentPage), 1)

	containerEndIndex := strings.LastIndex(htmlContent, "</div>")
	if containerEndIndex != -1 {
		beforeEnd := htmlContent[:containerEndIndex]
		afterEnd := htmlContent[containerEndIndex:]
		htmlContent = beforeEnd + tocContent + afterEnd
	}

	return fmt.Sprintf(htmlContent,
		apptouchPath, favicon32Path, favicon16Path, webmanifesstPath,
		cssPath, logoPath, navContent, content, backLink)
}

func generateTableOfContents(htmlContent string) string {
	h2Regex := regexp.MustCompile(`(?i)<h2[^>]*id=["']([^"']+)["'][^>]*>(.*?)</h2>`)
	matches := h2Regex.FindAllStringSubmatch(htmlContent, -1)

	if len(matches) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<aside class="toc-sidebar">
    <div class="toc-container">
        <ul id="toc-list">`)

	for _, match := range matches {
		id := match[1]
		title := stripHTMLTags(match[2])
		sb.WriteString(fmt.Sprintf(`
            <li><a href="#%s" data-target="%s">%s</a></li>`, id, id, title))
	}

	sb.WriteString(`
        </ul>
    </div>
</aside>`)

	return sb.String()
}

func stripHTMLTags(input string) string {
	tagRegex := regexp.MustCompile("<[^>]*>")
	return tagRegex.ReplaceAllString(input, "")
}

func addHeaderIDs(htmlContent string) string {
	h2Regex := regexp.MustCompile(`(?i)(<h2)([^>]*)>([^<]*)(</h2>)`)

	idMap := make(map[string]int)

	slugify := func(text string) string {
		text = stripHTMLTags(text)
		text = strings.ToLower(text)
		text = strings.ReplaceAll(text, " ", "-")
		reg := regexp.MustCompile(`[^a-z0-9\-]`)
		text = reg.ReplaceAllString(text, "")
		text = regexp.MustCompile(`\-+`).ReplaceAllString(text, "-")
		text = strings.Trim(text, "-")

		return text
	}

	result := h2Regex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		textMatch := h2Regex.FindStringSubmatch(match)
		if len(textMatch) < 4 {
			return match
		}

		if strings.Contains(strings.ToLower(match), "id=") {
			return match
		}

		headingText := textMatch[3]
		baseID := slugify(headingText)

		if baseID == "" {
			baseID = "heading"
		}

		finalID := baseID
		if count, exists := idMap[baseID]; exists {
			count++
			idMap[baseID] = count
			finalID = fmt.Sprintf("%s-%d", baseID, count)
		} else {
			idMap[baseID] = 1
		}

		return fmt.Sprintf("%s%s id=\"%s\">%s%s", textMatch[1], textMatch[2], finalID, textMatch[3], textMatch[4])
	})

	return result
}

func generateIndexPage(destDir string, navContent string) error {
	var links []string

	err := filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			if info.Name() == "index.html" {
				return nil
			}
			relPath, _ := filepath.Rel(destDir, path)
			links = append(links, relPath)
		}
		return nil
	})

	if err != nil {
		return err
	}

	var linksHTML strings.Builder
	for _, link := range links {
		nameWithoutExtension := strings.TrimSuffix(filepath.Base(link), ".html")
		formattedName := formatNavName(nameWithoutExtension)
		linksHTML.WriteString(fmt.Sprintf("<div class=\"notes-link\"><a href=\"/%s\">%s</a></div>\n", link, formattedName))
	}

	templateContent, err := readFile("index_template.html")

	if err != nil {
		log.Fatalf("Failed to read index HTML template: %v", err)
	}
	log.Printf("Successfully read index HTML template")

	finalContent := strings.Replace(string(templateContent), "{{ links }}", linksHTML.String(), 1)
	finalContent = strings.Replace(finalContent, "{{ navigation }}", navContent, 1)

	finalContent = strings.Replace(finalContent, `<img src="favicon-32x32.png" alt="Logo">`,
		`<img src="/favicon-32x32.png" alt="Logo">`, 1)

	finalContent = strings.Replace(finalContent, "<body>", "<body data-current-page=\"/index.html\">", 1)

	indexPath := filepath.Join(destDir, "index.html")
	return writeFile(indexPath, []byte(finalContent))
}

func processDirectory(srcDir, destDir string, navContent string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}

		relPath, _ := filepath.Rel(srcDir, path)
		newPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			if err := createDirectory(newPath); err != nil {
				return err
			}
			log.Printf("Created directory: %s", newPath)
			return nil
		}

		if isMarkdownFile(info.Name()) {
			return processFile(path, newPath, srcDir, navContent)
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

	log.Printf("Copy %s to %s", path, newPath)
	return nil
}

func main() {
	srcDir := "md"
	destDir := "public"

	navRoot, err := buildNavTree(srcDir)
	if err != nil {
		log.Fatalf("Error building navigation tree: %v", err)
	}

	navContent := renderNavTree(navRoot)

	if err := processDirectory(srcDir, destDir, navContent); err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}

	if err := generateIndexPage(destDir, navContent); err != nil {
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

	if err := copyFileToPublic("android-chrome-192x192.png"); err != nil {
		log.Fatalf("Error copying favicon: %v", err)
	}
	fmt.Println("Markdown files converted to HTML successfully, and index page created.")
	fmt.Println("To view the site locally, run a web server in the project root directory.")
	fmt.Println("If you have Python installed, you can use:")
	fmt.Println("  python -m http.server")
	fmt.Println("Then access the site at: http://localhost:8000/public/")
	fmt.Println("Or to serve from the public directory directly:")
	fmt.Println("  cd public && python -m http.server")
	fmt.Println("Then access the site at: http://localhost:8000/")
}
