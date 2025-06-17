package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"rewrite-migrate-java/pkg/java"
	"rewrite-migrate-java/pkg/migrate"
	"rewrite-migrate-java/pkg/recipe"
)

var (
	version = flag.Int("version", 17, "Target Java version (8, 11, 17, 21)")
	srcDir  = flag.String("src", "src/main/java", "Source directory to scan")
	dryRun  = flag.Bool("dry-run", false, "Show what would be changed without applying changes")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <project-path>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	projectPath := flag.Arg(0)

	// Initialize migration recipe based on target version
	var migrationRecipe recipe.Recipe
	switch *version {
	case 11:
		migrationRecipe = migrate.NewJava8ToJava11()
	case 17:
		migrationRecipe = migrate.NewUpgradeToJava17()
	case 21:
		migrationRecipe = migrate.NewUpgradeToJava21()
	default:
		migrationRecipe = migrate.NewUpgradeJavaVersion(*version)
	}

	fmt.Printf("Starting Java migration to version %d\n", *version)
	fmt.Printf("Recipe: %s\n", migrationRecipe.GetDisplayName())
	fmt.Printf("Description: %s\n", migrationRecipe.GetDescription())
	fmt.Printf("Estimated effort: %v\n\n", migrationRecipe.GetEstimatedEffortPerOccurrence())

	// Find and process files
	err := processProject(projectPath, migrationRecipe)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed successfully!")
}

func processProject(projectPath string, migrationRecipe recipe.Recipe) error {
	ctx := &recipe.ExecutionContext{
		Context:    context.Background(),
		Properties: make(map[string]interface{}),
	}

	visitor := migrationRecipe.GetVisitor()
	if visitor == nil {
		return fmt.Errorf("no visitor available for recipe")
	}

	// Process Java source files
	sourceDir := filepath.Join(projectPath, *srcDir)
	err := processSourceFiles(sourceDir, visitor, ctx)
	if err != nil {
		return fmt.Errorf("failed to process source files: %w", err)
	}

	// Process build files
	err = processBuildFiles(projectPath, visitor, ctx)
	if err != nil {
		return fmt.Errorf("failed to process build files: %w", err)
	}

	return nil
}

func processSourceFiles(sourceDir string, visitor recipe.TreeVisitor, ctx *recipe.ExecutionContext) error {
	return filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil // Skip if source directory doesn't exist
			}
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".java") {
			return nil
		}

		return processFile(path, visitor, ctx)
	})
}

func processBuildFiles(projectPath string, visitor recipe.TreeVisitor, ctx *recipe.ExecutionContext) error {
	buildFiles := []string{
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
	}

	for _, buildFile := range buildFiles {
		buildPath := filepath.Join(projectPath, buildFile)
		if _, err := os.Stat(buildPath); err == nil {
			err = processFile(buildPath, visitor, ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func processFile(path string, visitor recipe.TreeVisitor, ctx *recipe.ExecutionContext) error {
	fmt.Printf("Processing: %s\n", path)

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var sourceFile recipe.SourceFile
	if strings.HasSuffix(path, ".java") {
		sourceFile, err = java.NewJavaSourceFile(path, string(content))
		if err != nil {
			return fmt.Errorf("failed to parse Java file %s: %w", path, err)
		}
	} else {
		// For build files, create a simple wrapper
		sourceFile = &SimpleSourceFile{
			path:    path,
			content: string(content),
		}
	}

	// Apply transformations
	transformedFile, err := visitor.Visit(sourceFile, ctx)
	if err != nil {
		return fmt.Errorf("failed to transform file %s: %w", path, err)
	}

	// Check if file was modified
	if transformedFile.GetContent() != sourceFile.GetContent() {
		fmt.Printf("  → Modified\n")

		if *dryRun {
			fmt.Printf("  → Would write changes (dry-run mode)\n")
		} else {
			err = os.WriteFile(path, []byte(transformedFile.GetContent()), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", path, err)
			}
			fmt.Printf("  → Changes written\n")
		}
	}

	return nil
}

// SimpleSourceFile implements recipe.SourceFile for non-Java files
type SimpleSourceFile struct {
	path    string
	content string
}

func (s *SimpleSourceFile) GetPath() string {
	return s.path
}

func (s *SimpleSourceFile) GetContent() string {
	return s.content
}

func (s *SimpleSourceFile) GetClasses() []recipe.ClassDeclaration {
	return nil
}

func (s *SimpleSourceFile) GetImports() []recipe.ImportDeclaration {
	return nil
}

func (s *SimpleSourceFile) GetPackage() string {
	return ""
}

func (s *SimpleSourceFile) WithContent(content string) recipe.SourceFile {
	return &SimpleSourceFile{
		path:    s.path,
		content: content,
	}
}
