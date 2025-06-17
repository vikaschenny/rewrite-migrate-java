package migrate

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"rewrite-migrate-java/pkg/recipe"
)

// UpgradeJavaVersion upgrades Java version in build files and source markers
type UpgradeJavaVersion struct {
	*recipe.BaseRecipe
	Version int
}

// NewUpgradeJavaVersion creates a new UpgradeJavaVersion recipe
func NewUpgradeJavaVersion(version int) *UpgradeJavaVersion {
	return &UpgradeJavaVersion{
		BaseRecipe: &recipe.BaseRecipe{
			DisplayName: "Upgrade Java version",
			Description: fmt.Sprintf("Upgrade build plugin configuration to use Java %d. "+
				"This recipe changes java.toolchain.languageVersion in build.gradle(.kts) of gradle projects, "+
				"or maven-compiler-plugin target version and related settings. "+
				"Will not downgrade if the version is newer than the specified version.", version),
			EstimatedEffort: time.Duration(0), // No manual effort required
		},
		Version: version,
	}
}

func (u *UpgradeJavaVersion) GetVisitor() recipe.TreeVisitor {
	return &UpgradeJavaVersionVisitor{
		targetVersion: u.Version,
	}
}

func (u *UpgradeJavaVersion) GetRecipeList() []recipe.Recipe {
	// Composite recipe that includes multiple sub-recipes
	return []recipe.Recipe{
		NewUpdateMavenCompilerPlugin(u.Version),
		NewUpdateGradleJavaCompatibility(u.Version),
	}
}

// UpgradeJavaVersionVisitor handles the actual transformation
type UpgradeJavaVersionVisitor struct {
	targetVersion int
}

func (v *UpgradeJavaVersionVisitor) Visit(node recipe.SourceFile, ctx *recipe.ExecutionContext) (recipe.SourceFile, error) {
	content := node.GetContent()
	path := node.GetPath()

	// Handle different build file types
	if strings.HasSuffix(path, "pom.xml") {
		return v.updateMavenPom(node, content)
	}

	if strings.HasSuffix(path, "build.gradle") || strings.HasSuffix(path, "build.gradle.kts") {
		return v.updateGradleBuild(node, content)
	}

	// Return unchanged for non-build files
	return node, nil
}

func (v *UpgradeJavaVersionVisitor) updateMavenPom(node recipe.SourceFile, content string) (recipe.SourceFile, error) {
	// Update Maven compiler plugin configuration
	patterns := []struct {
		regex       *regexp.Regexp
		replacement string
	}{
		{
			regex:       regexp.MustCompile(`(<maven\.compiler\.source>)\d+(<\/maven\.compiler\.source>)`),
			replacement: fmt.Sprintf("${1}%d${2}", v.targetVersion),
		},
		{
			regex:       regexp.MustCompile(`(<maven\.compiler\.target>)\d+(<\/maven\.compiler\.target>)`),
			replacement: fmt.Sprintf("${1}%d${2}", v.targetVersion),
		},
		{
			regex:       regexp.MustCompile(`(<source>)\d+(<\/source>)`),
			replacement: fmt.Sprintf("${1}%d${2}", v.targetVersion),
		},
		{
			regex:       regexp.MustCompile(`(<target>)\d+(<\/target>)`),
			replacement: fmt.Sprintf("${1}%d${2}", v.targetVersion),
		},
		{
			regex:       regexp.MustCompile(`(<release>)\d+(<\/release>)`),
			replacement: fmt.Sprintf("${1}%d${2}", v.targetVersion),
		},
	}

	updatedContent := content
	for _, pattern := range patterns {
		updatedContent = pattern.regex.ReplaceAllString(updatedContent, pattern.replacement)
	}

	if updatedContent != content {
		return node.WithContent(updatedContent), nil
	}

	return node, nil
}

func (v *UpgradeJavaVersionVisitor) updateGradleBuild(node recipe.SourceFile, content string) (recipe.SourceFile, error) {
	// Update Gradle build configuration
	patterns := []struct {
		regex       *regexp.Regexp
		replacement string
	}{
		{
			regex:       regexp.MustCompile(`(sourceCompatibility\s*=\s*JavaVersion\.VERSION_)\d+`),
			replacement: fmt.Sprintf("${1}%d", v.targetVersion),
		},
		{
			regex:       regexp.MustCompile(`(targetCompatibility\s*=\s*JavaVersion\.VERSION_)\d+`),
			replacement: fmt.Sprintf("${1}%d", v.targetVersion),
		},
		{
			regex:       regexp.MustCompile(`(sourceCompatibility\s*=\s*)\d+`),
			replacement: fmt.Sprintf("${1}%d", v.targetVersion),
		},
		{
			regex:       regexp.MustCompile(`(targetCompatibility\s*=\s*)\d+`),
			replacement: fmt.Sprintf("${1}%d", v.targetVersion),
		},
		{
			regex:       regexp.MustCompile(`(languageVersion\s*=\s*JavaLanguageVersion\.of\()\d+(\))`),
			replacement: fmt.Sprintf("${1}%d${2}", v.targetVersion),
		},
	}

	updatedContent := content
	for _, pattern := range patterns {
		updatedContent = pattern.regex.ReplaceAllString(updatedContent, pattern.replacement)
	}

	if updatedContent != content {
		return node.WithContent(updatedContent), nil
	}

	return node, nil
}

// UpdateMavenCompilerPlugin updates Maven compiler plugin configuration
type UpdateMavenCompilerPlugin struct {
	*recipe.BaseRecipe
	Version int
}

func NewUpdateMavenCompilerPlugin(version int) *UpdateMavenCompilerPlugin {
	return &UpdateMavenCompilerPlugin{
		BaseRecipe: &recipe.BaseRecipe{
			DisplayName: "Update Maven Compiler Plugin",
			Description: "Update Maven compiler plugin to use specified Java version",
		},
		Version: version,
	}
}

func (u *UpdateMavenCompilerPlugin) GetVisitor() recipe.TreeVisitor {
	return &UpgradeJavaVersionVisitor{targetVersion: u.Version}
}

// UpdateGradleJavaCompatibility updates Gradle Java compatibility settings
type UpdateGradleJavaCompatibility struct {
	*recipe.BaseRecipe
	Version int
}

func NewUpdateGradleJavaCompatibility(version int) *UpdateGradleJavaCompatibility {
	return &UpdateGradleJavaCompatibility{
		BaseRecipe: &recipe.BaseRecipe{
			DisplayName: "Update Gradle Java Compatibility",
			Description: "Update Gradle Java compatibility settings to specified version",
		},
		Version: version,
	}
}

func (u *UpdateGradleJavaCompatibility) GetVisitor() recipe.TreeVisitor {
	return &UpgradeJavaVersionVisitor{targetVersion: u.Version}
}
