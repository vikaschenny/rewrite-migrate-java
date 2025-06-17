package recipe

import (
	"context"
	"time"
)

// ExecutionContext holds context information for recipe execution
type ExecutionContext struct {
	context.Context
	Properties map[string]interface{}
}

// Recipe represents a migration recipe that can be applied to Java source code
type Recipe interface {
	GetDisplayName() string
	GetDescription() string
	GetEstimatedEffortPerOccurrence() time.Duration
	GetVisitor() TreeVisitor
	GetRecipeList() []Recipe
	ApplicabilityTest() Precondition
}

// TreeVisitor visits and potentially modifies nodes in a source tree
type TreeVisitor interface {
	Visit(node SourceFile, ctx *ExecutionContext) (SourceFile, error)
}

// Precondition represents a condition that must be met for a recipe to apply
type Precondition interface {
	Check(sourceFile SourceFile) bool
}

// SourceFile represents a Java source file with its AST and metadata
type SourceFile interface {
	GetPath() string
	GetContent() string
	GetClasses() []ClassDeclaration
	GetImports() []ImportDeclaration
	GetPackage() string
	WithContent(content string) SourceFile
}

// ClassDeclaration represents a Java class declaration
type ClassDeclaration interface {
	GetSimpleName() string
	GetFullyQualifiedName() string
	GetMethods() []MethodDeclaration
	GetFields() []FieldDeclaration
}

// MethodDeclaration represents a Java method declaration
type MethodDeclaration interface {
	GetName() string
	GetReturnType() string
	GetParameters() []Parameter
	GetBody() string
}

// FieldDeclaration represents a Java field declaration
type FieldDeclaration interface {
	GetName() string
	GetType() string
	GetModifiers() []string
}

// ImportDeclaration represents a Java import statement
type ImportDeclaration interface {
	GetPackageName() string
	IsStatic() bool
	IsWildcard() bool
}

// Parameter represents a method parameter
type Parameter interface {
	GetName() string
	GetType() string
}

// BaseRecipe provides a basic implementation of Recipe
type BaseRecipe struct {
	DisplayName     string
	Description     string
	EstimatedEffort time.Duration
}

func (r *BaseRecipe) GetDisplayName() string {
	return r.DisplayName
}

func (r *BaseRecipe) GetDescription() string {
	return r.Description
}

func (r *BaseRecipe) GetEstimatedEffortPerOccurrence() time.Duration {
	return r.EstimatedEffort
}

func (r *BaseRecipe) GetRecipeList() []Recipe {
	return []Recipe{} // Override in composite recipes
}

func (r *BaseRecipe) ApplicabilityTest() Precondition {
	return nil // Override if needed
}

// CompositeRecipe combines multiple recipes
type CompositeRecipe struct {
	*BaseRecipe
	Recipes []Recipe
}

func (c *CompositeRecipe) GetRecipeList() []Recipe {
	return c.Recipes
}

func (c *CompositeRecipe) GetVisitor() TreeVisitor {
	return &CompositeVisitor{Visitors: c.getVisitors()}
}

func (c *CompositeRecipe) getVisitors() []TreeVisitor {
	var visitors []TreeVisitor
	for _, recipe := range c.Recipes {
		if visitor := recipe.GetVisitor(); visitor != nil {
			visitors = append(visitors, visitor)
		}
	}
	return visitors
}

// CompositeVisitor applies multiple visitors in sequence
type CompositeVisitor struct {
	Visitors []TreeVisitor
}

func (v *CompositeVisitor) Visit(node SourceFile, ctx *ExecutionContext) (SourceFile, error) {
	current := node
	for _, visitor := range v.Visitors {
		var err error
		current, err = visitor.Visit(current, ctx)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}
