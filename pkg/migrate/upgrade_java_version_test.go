package migrate

import (
	"context"
	"strings"
	"testing"

	"rewrite-migrate-java/pkg/recipe"
)

func TestUpgradeJavaVersionMaven(t *testing.T) {
	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <properties>
        <maven.compiler.source>8</maven.compiler.source>
        <maven.compiler.target>8</maven.compiler.target>
    </properties>
</project>`

	expectedContent := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
    </properties>
</project>`

	recipe := NewUpgradeJavaVersion(17)
	visitor := recipe.GetVisitor()

	sourceFile := &mockSourceFile{
		path:    "pom.xml",
		content: pomContent,
	}

	ctx := &recipe.ExecutionContext{
		Context:    context.Background(),
		Properties: make(map[string]interface{}),
	}

	result, err := visitor.Visit(sourceFile, ctx)
	if err != nil {
		t.Fatalf("Visit failed: %v", err)
	}

	if result.GetContent() != expectedContent {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expectedContent, result.GetContent())
	}
}

func TestUpgradeJavaVersionGradle(t *testing.T) {
	gradleContent := `plugins {
    id 'java'
}

sourceCompatibility = JavaVersion.VERSION_8
targetCompatibility = JavaVersion.VERSION_8`

	expectedContent := `plugins {
    id 'java'
}

sourceCompatibility = JavaVersion.VERSION_17
targetCompatibility = JavaVersion.VERSION_17`

	recipe := NewUpgradeJavaVersion(17)
	visitor := recipe.GetVisitor()

	sourceFile := &mockSourceFile{
		path:    "build.gradle",
		content: gradleContent,
	}

	ctx := &recipe.ExecutionContext{
		Context:    context.Background(),
		Properties: make(map[string]interface{}),
	}

	result, err := visitor.Visit(sourceFile, ctx)
	if err != nil {
		t.Fatalf("Visit failed: %v", err)
	}

	if result.GetContent() != expectedContent {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expectedContent, result.GetContent())
	}
}

func TestUseJavaUtilBase64(t *testing.T) {
	javaContent := `import sun.misc.BASE64Encoder;
import sun.misc.BASE64Decoder;

public class Example {
    public void encode() {
        BASE64Encoder encoder = new BASE64Encoder();
        String result = encoder.encode(data);
    }
    
    public void decode() {
        BASE64Decoder decoder = new BASE64Decoder();
        byte[] result = decoder.decodeBuffer(data);
    }
}`

	recipe := NewUseJavaUtilBase64("sun.misc", false)
	visitor := recipe.GetVisitor()

	sourceFile := &mockSourceFile{
		path:    "Example.java",
		content: javaContent,
	}

	ctx := &recipe.ExecutionContext{
		Context:    context.Background(),
		Properties: make(map[string]interface{}),
	}

	result, err := visitor.Visit(sourceFile, ctx)
	if err != nil {
		t.Fatalf("Visit failed: %v", err)
	}

	resultContent := result.GetContent()

	// Check that imports were updated
	if !strings.Contains(resultContent, "import java.util.Base64;") {
		t.Error("Expected java.util.Base64 import")
	}

	// Check that sun.misc imports were removed
	if strings.Contains(resultContent, "sun.misc.BASE64") {
		t.Error("sun.misc.BASE64 imports should be removed")
	}

	// Check that Base64 encoder/decoder instantiation was updated
	if strings.Contains(resultContent, "new BASE64Encoder()") ||
		strings.Contains(resultContent, "new BASE64Decoder()") {
		t.Error("BASE64Encoder/Decoder instantiation should be replaced")
	}
}

// mockSourceFile implements recipe.SourceFile for testing
type mockSourceFile struct {
	path    string
	content string
}

func (m *mockSourceFile) GetPath() string {
	return m.path
}

func (m *mockSourceFile) GetContent() string {
	return m.content
}

func (m *mockSourceFile) GetClasses() []recipe.ClassDeclaration {
	return nil
}

func (m *mockSourceFile) GetImports() []recipe.ImportDeclaration {
	return nil
}

func (m *mockSourceFile) GetPackage() string {
	return ""
}

func (m *mockSourceFile) WithContent(content string) recipe.SourceFile {
	return &mockSourceFile{
		path:    m.path,
		content: content,
	}
}
