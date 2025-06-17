package java

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"rewrite-migrate-java/pkg/recipe"
)

// JavaSourceFile implements recipe.SourceFile for Java source files
type JavaSourceFile struct {
	path    string
	content string
	pkg     string
	imports []recipe.ImportDeclaration
	classes []recipe.ClassDeclaration
}

// NewJavaSourceFile creates a new JavaSourceFile from content
func NewJavaSourceFile(path, content string) (*JavaSourceFile, error) {
	jsf := &JavaSourceFile{
		path:    path,
		content: content,
	}

	if err := jsf.parse(); err != nil {
		return nil, fmt.Errorf("failed to parse Java source file %s: %w", path, err)
	}

	return jsf, nil
}

func (jsf *JavaSourceFile) GetPath() string {
	return jsf.path
}

func (jsf *JavaSourceFile) GetContent() string {
	return jsf.content
}

func (jsf *JavaSourceFile) GetClasses() []recipe.ClassDeclaration {
	return jsf.classes
}

func (jsf *JavaSourceFile) GetImports() []recipe.ImportDeclaration {
	return jsf.imports
}

func (jsf *JavaSourceFile) GetPackage() string {
	return jsf.pkg
}

func (jsf *JavaSourceFile) WithContent(content string) recipe.SourceFile {
	newFile := &JavaSourceFile{
		path:    jsf.path,
		content: content,
	}
	newFile.parse() // Re-parse with new content
	return newFile
}

// parse performs basic parsing of the Java source file
func (jsf *JavaSourceFile) parse() error {
	scanner := bufio.NewScanner(strings.NewReader(jsf.content))

	// Parse package declaration
	packageRegex := regexp.MustCompile(`^\s*package\s+([a-zA-Z_][a-zA-Z0-9_.]*)\s*;`)

	// Parse import statements
	importRegex := regexp.MustCompile(`^\s*import\s+(static\s+)?([a-zA-Z_][a-zA-Z0-9_.*]*)\s*;`)

	// Parse class declarations
	classRegex := regexp.MustCompile(`^\s*(public\s+|private\s+|protected\s+)?(abstract\s+|final\s+)?(class|interface|enum)\s+([a-zA-Z_][a-zA-Z0-9_]*)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "//") || strings.HasPrefix(strings.TrimSpace(line), "/*") {
			continue
		}

		// Parse package
		if matches := packageRegex.FindStringSubmatch(line); matches != nil {
			jsf.pkg = matches[1]
			continue
		}

		// Parse imports
		if matches := importRegex.FindStringSubmatch(line); matches != nil {
			isStatic := strings.TrimSpace(matches[1]) == "static"
			packageName := matches[2]
			isWildcard := strings.HasSuffix(packageName, ".*")

			jsf.imports = append(jsf.imports, &JavaImportDeclaration{
				packageName: packageName,
				isStatic:    isStatic,
				isWildcard:  isWildcard,
			})
			continue
		}

		// Parse class declarations
		if matches := classRegex.FindStringSubmatch(line); matches != nil {
			className := matches[4]
			fullyQualifiedName := className
			if jsf.pkg != "" {
				fullyQualifiedName = jsf.pkg + "." + className
			}

			jsf.classes = append(jsf.classes, &JavaClassDeclaration{
				simpleName:         className,
				fullyQualifiedName: fullyQualifiedName,
				methods:            []recipe.MethodDeclaration{},
				fields:             []recipe.FieldDeclaration{},
			})
		}
	}

	return scanner.Err()
}

// JavaImportDeclaration implements recipe.ImportDeclaration
type JavaImportDeclaration struct {
	packageName string
	isStatic    bool
	isWildcard  bool
}

func (i *JavaImportDeclaration) GetPackageName() string {
	return i.packageName
}

func (i *JavaImportDeclaration) IsStatic() bool {
	return i.isStatic
}

func (i *JavaImportDeclaration) IsWildcard() bool {
	return i.isWildcard
}

// JavaClassDeclaration implements recipe.ClassDeclaration
type JavaClassDeclaration struct {
	simpleName         string
	fullyQualifiedName string
	methods            []recipe.MethodDeclaration
	fields             []recipe.FieldDeclaration
}

func (c *JavaClassDeclaration) GetSimpleName() string {
	return c.simpleName
}

func (c *JavaClassDeclaration) GetFullyQualifiedName() string {
	return c.fullyQualifiedName
}

func (c *JavaClassDeclaration) GetMethods() []recipe.MethodDeclaration {
	return c.methods
}

func (c *JavaClassDeclaration) GetFields() []recipe.FieldDeclaration {
	return c.fields
}

// JavaMethodDeclaration implements recipe.MethodDeclaration
type JavaMethodDeclaration struct {
	name       string
	returnType string
	parameters []recipe.Parameter
	body       string
}

func (m *JavaMethodDeclaration) GetName() string {
	return m.name
}

func (m *JavaMethodDeclaration) GetReturnType() string {
	return m.returnType
}

func (m *JavaMethodDeclaration) GetParameters() []recipe.Parameter {
	return m.parameters
}

func (m *JavaMethodDeclaration) GetBody() string {
	return m.body
}

// JavaFieldDeclaration implements recipe.FieldDeclaration
type JavaFieldDeclaration struct {
	name      string
	fieldType string
	modifiers []string
}

func (f *JavaFieldDeclaration) GetName() string {
	return f.name
}

func (f *JavaFieldDeclaration) GetType() string {
	return f.fieldType
}

func (f *JavaFieldDeclaration) GetModifiers() []string {
	return f.modifiers
}

// JavaParameter implements recipe.Parameter
type JavaParameter struct {
	name      string
	paramType string
}

func (p *JavaParameter) GetName() string {
	return p.name
}

func (p *JavaParameter) GetType() string {
	return p.paramType
}
