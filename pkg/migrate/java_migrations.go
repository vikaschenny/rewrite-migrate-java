package migrate

import (
	"regexp"
	"strings"
	"time"

	"rewrite-migrate-java/pkg/recipe"
)

// replacePackageReferences replaces all references to oldPackage with newPackage
func replacePackageReferences(content, oldPackage, newPackage string) string {
	// Replace import statements
	importPattern := regexp.MustCompile(`import\s+` + regexp.QuoteMeta(oldPackage) + `(\.[^;]*)?;`)
	content = importPattern.ReplaceAllStringFunc(content, func(match string) string {
		return strings.Replace(match, oldPackage, newPackage, 1)
	})

	// Replace package references in code
	// This is a simple replacement - more sophisticated parsing would be better
	content = strings.ReplaceAll(content, oldPackage+".", newPackage+".")

	return content
}

// Java8ToJava11 provides a composite recipe for migrating from Java 8 to Java 11
type Java8ToJava11 struct {
	*recipe.CompositeRecipe
}

// NewJava8ToJava11 creates a new Java 8 to 11 migration recipe
func NewJava8ToJava11() *Java8ToJava11 {
	return &Java8ToJava11{
		CompositeRecipe: &recipe.CompositeRecipe{
			BaseRecipe: &recipe.BaseRecipe{
				DisplayName: "Migrate from Java 8 to Java 11",
				Description: "Migrates Java 8 applications to Java 11, including updating dependencies, " +
					"replacing deprecated APIs, and handling Java EE to Jakarta EE transitions.",
				EstimatedEffort: 30 * time.Minute,
			},
			Recipes: []recipe.Recipe{
				NewUpgradeJavaVersion(11),
				NewUseJavaUtilBase64("sun.misc", false),
				NewJavaEEToJakartaEE(),
				NewRemoveDeprecatedAPIs(),
			},
		},
	}
}

// UpgradeToJava17 provides a composite recipe for upgrading to Java 17
type UpgradeToJava17 struct {
	*recipe.CompositeRecipe
}

// NewUpgradeToJava17 creates a new Java 17 upgrade recipe
func NewUpgradeToJava17() *UpgradeToJava17 {
	return &UpgradeToJava17{
		CompositeRecipe: &recipe.CompositeRecipe{
			BaseRecipe: &recipe.BaseRecipe{
				DisplayName: "Upgrade to Java 17",
				Description: "Upgrades applications to Java 17, handling deprecated APIs and " +
					"illegal reflective access issues.",
				EstimatedEffort: 20 * time.Minute,
			},
			Recipes: []recipe.Recipe{
				NewUpgradeJavaVersion(17),
				NewUseJavaUtilBase64("sun.misc", false),
				NewFixReflectiveAccess(),
			},
		},
	}
}

// UpgradeToJava21 provides a composite recipe for upgrading to Java 21
type UpgradeToJava21 struct {
	*recipe.CompositeRecipe
}

// NewUpgradeToJava21 creates a new Java 21 upgrade recipe
func NewUpgradeToJava21() *UpgradeToJava21 {
	return &UpgradeToJava21{
		CompositeRecipe: &recipe.CompositeRecipe{
			BaseRecipe: &recipe.BaseRecipe{
				DisplayName: "Upgrade to Java 21",
				Description: "Upgrades applications to Java 21, including all Java 17 changes plus " +
					"support for sequenced collections and other Java 21 features.",
				EstimatedEffort: 15 * time.Minute,
			},
			Recipes: []recipe.Recipe{
				NewUpgradeJavaVersion(21),
				NewUseJavaUtilBase64("sun.misc", false),
				NewFixReflectiveAccess(),
				NewSequencedCollectionsMigration(),
			},
		},
	}
}

// JavaEEToJakartaEE migrates Java EE dependencies to Jakarta EE
type JavaEEToJakartaEE struct {
	*recipe.BaseRecipe
}

// NewJavaEEToJakartaEE creates a new Java EE to Jakarta EE migration recipe
func NewJavaEEToJakartaEE() *JavaEEToJakartaEE {
	return &JavaEEToJakartaEE{
		BaseRecipe: &recipe.BaseRecipe{
			DisplayName:     "Migrate Java EE to Jakarta EE",
			Description:     "Migrates Java EE dependencies and imports to Jakarta EE equivalents",
			EstimatedEffort: 10 * time.Minute,
		},
	}
}

func (j *JavaEEToJakartaEE) GetVisitor() recipe.TreeVisitor {
	return &JavaEEToJakartaEEVisitor{}
}

// JavaEEToJakartaEEVisitor handles Java EE to Jakarta EE transformations
type JavaEEToJakartaEEVisitor struct{}

func (v *JavaEEToJakartaEEVisitor) Visit(node recipe.SourceFile, ctx *recipe.ExecutionContext) (recipe.SourceFile, error) {
	content := node.GetContent()

	// Common Java EE to Jakarta EE package migrations
	replacements := map[string]string{
		"javax.persistence": "jakarta.persistence",
		"javax.servlet":     "jakarta.servlet",
		"javax.ejb":         "jakarta.ejb",
		"javax.jms":         "jakarta.jms",
		"javax.mail":        "jakarta.mail",
		"javax.xml.bind":    "jakarta.xml.bind",
		"javax.xml.ws":      "jakarta.xml.ws",
		"javax.annotation":  "jakarta.annotation",
		"javax.enterprise":  "jakarta.enterprise",
		"javax.inject":      "jakarta.inject",
		"javax.interceptor": "jakarta.interceptor",
		"javax.validation":  "jakarta.validation",
		"javax.ws.rs":       "jakarta.ws.rs",
		"javax.json":        "jakarta.json",
		"javax.websocket":   "jakarta.websocket",
	}

	for javaEE, jakartaEE := range replacements {
		content = replacePackageReferences(content, javaEE, jakartaEE)
	}

	if content != node.GetContent() {
		return node.WithContent(content), nil
	}

	return node, nil
}

// RemoveDeprecatedAPIs removes usage of deprecated APIs
type RemoveDeprecatedAPIs struct {
	*recipe.BaseRecipe
}

// NewRemoveDeprecatedAPIs creates a recipe to remove deprecated API usage
func NewRemoveDeprecatedAPIs() *RemoveDeprecatedAPIs {
	return &RemoveDeprecatedAPIs{
		BaseRecipe: &recipe.BaseRecipe{
			DisplayName:     "Remove Deprecated API Usage",
			Description:     "Removes usage of APIs that have been deprecated and provides modern alternatives",
			EstimatedEffort: 15 * time.Minute,
		},
	}
}

func (r *RemoveDeprecatedAPIs) GetVisitor() recipe.TreeVisitor {
	return &RemoveDeprecatedAPIsVisitor{}
}

// RemoveDeprecatedAPIsVisitor handles deprecated API transformations
type RemoveDeprecatedAPIsVisitor struct{}

func (v *RemoveDeprecatedAPIsVisitor) Visit(node recipe.SourceFile, ctx *recipe.ExecutionContext) (recipe.SourceFile, error) {
	// This would contain specific deprecated API replacements
	// For now, return unchanged
	return node, nil
}

// FixReflectiveAccess handles illegal reflective access issues
type FixReflectiveAccess struct {
	*recipe.BaseRecipe
}

// NewFixReflectiveAccess creates a recipe to fix reflective access issues
func NewFixReflectiveAccess() *FixReflectiveAccess {
	return &FixReflectiveAccess{
		BaseRecipe: &recipe.BaseRecipe{
			DisplayName:     "Fix Illegal Reflective Access",
			Description:     "Updates code to fix illegal reflective access issues in Java 17+",
			EstimatedEffort: 20 * time.Minute,
		},
	}
}

func (f *FixReflectiveAccess) GetVisitor() recipe.TreeVisitor {
	return &FixReflectiveAccessVisitor{}
}

// FixReflectiveAccessVisitor handles reflective access fixes
type FixReflectiveAccessVisitor struct{}

func (v *FixReflectiveAccessVisitor) Visit(node recipe.SourceFile, ctx *recipe.ExecutionContext) (recipe.SourceFile, error) {
	// This would contain specific reflective access fixes
	// For now, return unchanged
	return node, nil
}

// SequencedCollectionsMigration handles Java 21 sequenced collections
type SequencedCollectionsMigration struct {
	*recipe.BaseRecipe
}

// NewSequencedCollectionsMigration creates a recipe for sequenced collections
func NewSequencedCollectionsMigration() *SequencedCollectionsMigration {
	return &SequencedCollectionsMigration{
		BaseRecipe: &recipe.BaseRecipe{
			DisplayName:     "Migrate to Sequenced Collections",
			Description:     "Updates code to use Java 21 sequenced collections features",
			EstimatedEffort: 10 * time.Minute,
		},
	}
}

func (s *SequencedCollectionsMigration) GetVisitor() recipe.TreeVisitor {
	return &SequencedCollectionsMigrationVisitor{}
}

// SequencedCollectionsMigrationVisitor handles sequenced collections migration
type SequencedCollectionsMigrationVisitor struct{}

func (v *SequencedCollectionsMigrationVisitor) Visit(node recipe.SourceFile, ctx *recipe.ExecutionContext) (recipe.SourceFile, error) {
	// This would contain specific sequenced collections updates
	// For now, return unchanged
	return node, nil
}
