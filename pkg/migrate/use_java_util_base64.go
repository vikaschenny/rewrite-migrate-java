package migrate

import (
	"regexp"
	"strings"

	"rewrite-migrate-java/pkg/recipe"
)

// UseJavaUtilBase64 replaces sun.misc Base64 classes with java.util.Base64
type UseJavaUtilBase64 struct {
	*recipe.BaseRecipe
	SunPackage   string
	UseMimeCoder bool
}

// NewUseJavaUtilBase64 creates a new UseJavaUtilBase64 recipe
func NewUseJavaUtilBase64(sunPackage string, useMimeCoder bool) *UseJavaUtilBase64 {
	if sunPackage == "" {
		sunPackage = "sun.misc"
	}

	return &UseJavaUtilBase64{
		BaseRecipe: &recipe.BaseRecipe{
			DisplayName: "Prefer `java.util.Base64` instead of `sun.misc`",
			Description: "Prefer `java.util.Base64` instead of using `sun.misc` in Java 8 or higher. " +
				"`sun.misc` is not exported by the Java module system and accessing this class will " +
				"result in a warning in Java 11 and an error in Java 17.",
		},
		SunPackage:   sunPackage,
		UseMimeCoder: useMimeCoder,
	}
}

func (u *UseJavaUtilBase64) GetVisitor() recipe.TreeVisitor {
	return &UseJavaUtilBase64Visitor{
		sunPackage:   u.SunPackage,
		useMimeCoder: u.UseMimeCoder,
	}
}

func (u *UseJavaUtilBase64) ApplicabilityTest() recipe.Precondition {
	return &UsesTypePrecondition{
		types: []string{
			u.SunPackage + ".BASE64Encoder",
			u.SunPackage + ".BASE64Decoder",
		},
	}
}

// UseJavaUtilBase64Visitor handles the actual transformation
type UseJavaUtilBase64Visitor struct {
	sunPackage   string
	useMimeCoder bool
}

func (v *UseJavaUtilBase64Visitor) Visit(node recipe.SourceFile, ctx *recipe.ExecutionContext) (recipe.SourceFile, error) {
	content := node.GetContent()

	// Check if already using incompatible Base64
	if v.alreadyUsingIncompatibleBase64(node) {
		// Return with warning - manual intervention required
		return node, nil
	}

	// Replace imports
	content = v.replaceImports(content)

	// Replace class instantiations
	content = v.replaceClassInstantiations(content)

	// Replace method calls
	content = v.replaceMethodCalls(content)

	if content != node.GetContent() {
		return node.WithContent(content), nil
	}

	return node, nil
}

func (v *UseJavaUtilBase64Visitor) replaceImports(content string) string {
	patterns := []struct {
		regex       *regexp.Regexp
		replacement string
	}{
		{
			regex:       regexp.MustCompile(`import\s+` + regexp.QuoteMeta(v.sunPackage) + `\.BASE64Encoder;`),
			replacement: "import java.util.Base64;",
		},
		{
			regex:       regexp.MustCompile(`import\s+` + regexp.QuoteMeta(v.sunPackage) + `\.BASE64Decoder;`),
			replacement: "import java.util.Base64;",
		},
	}

	for _, pattern := range patterns {
		content = pattern.regex.ReplaceAllString(content, pattern.replacement)
	}

	return content
}

func (v *UseJavaUtilBase64Visitor) replaceClassInstantiations(content string) string {
	patterns := []struct {
		regex       *regexp.Regexp
		replacement string
	}{
		{
			regex: regexp.MustCompile(`new\s+BASE64Encoder\s*\(\s*\)`),
			replacement: func() string {
				if v.useMimeCoder {
					return "Base64.getMimeEncoder()"
				}
				return "Base64.getEncoder()"
			}(),
		},
		{
			regex: regexp.MustCompile(`new\s+BASE64Decoder\s*\(\s*\)`),
			replacement: func() string {
				if v.useMimeCoder {
					return "Base64.getMimeDecoder()"
				}
				return "Base64.getDecoder()"
			}(),
		},
	}

	for _, pattern := range patterns {
		content = pattern.regex.ReplaceAllString(content, pattern.replacement)
	}

	return content
}

func (v *UseJavaUtilBase64Visitor) replaceMethodCalls(content string) string {
	// Replace encode methods
	encodePatterns := []struct {
		regex       *regexp.Regexp
		replacement string
	}{
		{
			regex: regexp.MustCompile(`(\w+)\.encode\s*\(\s*([^)]+)\s*\)`),
			replacement: func() string {
				if v.useMimeCoder {
					return "Base64.getMimeEncoder().encodeToString($2)"
				}
				return "Base64.getEncoder().encodeToString($2)"
			}(),
		},
		{
			regex: regexp.MustCompile(`(\w+)\.encodeBuffer\s*\(\s*([^)]+)\s*\)`),
			replacement: func() string {
				if v.useMimeCoder {
					return "Base64.getMimeEncoder().encodeToString($2)"
				}
				return "Base64.getEncoder().encodeToString($2)"
			}(),
		},
	}

	// Replace decode methods
	decodePatterns := []struct {
		regex       *regexp.Regexp
		replacement string
	}{
		{
			regex: regexp.MustCompile(`(\w+)\.decodeBuffer\s*\(\s*([^)]+)\s*\)`),
			replacement: func() string {
				if v.useMimeCoder {
					return "Base64.getMimeDecoder().decode($2)"
				}
				return "Base64.getDecoder().decode($2)"
			}(),
		},
	}

	allPatterns := append(encodePatterns, decodePatterns...)

	for _, pattern := range allPatterns {
		content = pattern.regex.ReplaceAllString(content, pattern.replacement)
	}

	return content
}

func (v *UseJavaUtilBase64Visitor) alreadyUsingIncompatibleBase64(node recipe.SourceFile) bool {
	// Check if there's already a Base64 class that's not java.util.Base64
	for _, class := range node.GetClasses() {
		if class.GetSimpleName() == "Base64" {
			return true
		}
	}

	// Check imports for incompatible Base64 classes
	for _, imp := range node.GetImports() {
		packageName := imp.GetPackageName()
		if strings.HasSuffix(packageName, ".Base64") && packageName != "java.util.Base64" {
			return true
		}
	}

	return false
}

// UsesTypePrecondition checks if the source file uses specific types
type UsesTypePrecondition struct {
	types []string
}

func (p *UsesTypePrecondition) Check(sourceFile recipe.SourceFile) bool {
	content := sourceFile.GetContent()

	for _, typeName := range p.types {
		// Simple check for type usage in content
		if strings.Contains(content, typeName) {
			return true
		}
	}

	return false
}
