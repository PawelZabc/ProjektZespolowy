// tools/assetgen/main.go
package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"
)

//go:embed template.go.tmpl
var templateFS embed.FS

type AssetCategory struct {
	Name      string // Models, Textures, etc.
	Directory string // models, graphics, etc.
	Assets    []Asset
}

type Asset struct {
	ConstName string // MODEL_PLAYER
	Filename  string // player.glb
	Comment   string // Player model
}

type TemplateData struct {
	Timestamp  string
	Categories []AssetCategory
}

func main() {
	if err := generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Generated assets_gen.go successfully")
}

// findAssetsPath searches for the client/assets directory
func findAssetsPath() string {
	// Try different relative paths
	possiblePaths := []string{
		"../../client/assets", // From tools/assetgen/
		"client/assets",       // From project root
		"assets",              // From client/
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			abs, _ := filepath.Abs(path)
			return abs
		}
	}

	return ""
}

func generate() error {
	// Find assets directory - try relative paths from tools/assetgen/
	assetsPath := findAssetsPath()
	if assetsPath == "" {
		return fmt.Errorf("could not find client/assets directory")
	}

	fmt.Printf("Found assets directory: %s\n", assetsPath)

	// Define asset categories and their directories
	categories := []struct {
		name string
		dir  string
	}{
		{"Models", "models"},
		{"Textures", "graphics"},
		{"Sounds", "audio"},
		{"Fonts", "fonts"},
		{"Shaders", "shaders"},
	}

	var assetCategories []AssetCategory

	// Scan each directory
	for _, cat := range categories {
		dirPath := filepath.Join(assetsPath, cat.dir)

		// Check if directory exists
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Printf("Directory not found: %s (skipping)\n", dirPath)
			continue
		}

		assets, err := scanDirectory(dirPath, cat.name)
		if err != nil {
			return fmt.Errorf("failed to scan %s: %w", cat.dir, err)
		}

		if len(assets) > 0 {
			assetCategories = append(assetCategories, AssetCategory{
				Name:      cat.name,
				Directory: cat.dir,
				Assets:    assets,
			})
			fmt.Printf("Found %d %s\n", len(assets), cat.name)
		}
	}

	// Load template
	tmplContent, err := templateFS.ReadFile("template.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New("assets").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	data := TemplateData{
		Timestamp:  time.Now().Format(time.RFC3339),
		Categories: assetCategories,
	}

	// Create output file
	outputPath := filepath.Join(assetsPath, "assets_gen.go")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Execute template
	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func scanDirectory(dirPath, categoryName string) ([]Asset, error) {
	var assets []Asset

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		// Skip hidden files and non-asset files
		if strings.HasPrefix(filename, ".") {
			continue
		}

		constName := generateConstName(categoryName, filename)
		comment := generateComment(filename)

		assets = append(assets, Asset{
			ConstName: constName,
			Filename:  filename,
			Comment:   comment,
		})
	}

	return assets, nil
}

func generateConstName(category, filename string) string {
	// Remove extension
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Convert to PascalCase
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || r == '.'
	})

	var result strings.Builder
	result.WriteString(strings.TrimSuffix(category, "s")) // Remove plural: Models -> Model

	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(strings.ToLower(part[1:]))
			}
		}
	}

	return result.String()
}

func generateComment(filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Convert to readable format
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	// Capitalize first letter
	if len(name) > 0 {
		runes := []rune(name)
		runes[0] = unicode.ToUpper(runes[0])
		name = string(runes)
	}

	return name
}
