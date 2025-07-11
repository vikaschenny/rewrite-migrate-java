<p align="center">
  <a href="https://docs.openrewrite.org">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://github.com/openrewrite/rewrite/raw/main/doc/logo-oss-dark.svg">
      <source media="(prefers-color-scheme: light)" srcset="https://github.com/openrewrite/rewrite/raw/main/doc/logo-oss-light.svg">
      <img alt="OpenRewrite Logo" src="https://github.com/openrewrite/rewrite/raw/main/doc/logo-oss-light.svg" width='600px'>
    </picture>
  </a>
</p>

<div align="center">
  <h1>rewrite-migrate-java</h1>
</div>

<div align="center">

<!-- Keep the gap above this line, otherwise they won't render correctly! -->
[![ci](https://github.com/openrewrite/rewrite-migrate-java/actions/workflows/ci.yml/badge.svg)](https://github.com/openrewrite/rewrite-migrate-java/actions/workflows/ci.yml)
[![Maven Central](https://img.shields.io/maven-central/v/org.openrewrite.recipe/rewrite-migrate-java.svg)](https://mvnrepository.com/artifact/org.openrewrite.recipe/rewrite-migrate-java)
[![Revved up by Develocity](https://img.shields.io/badge/Revved%20up%20by-Develocity-06A0CE?logo=Gradle&labelColor=02303A)](https://ge.openrewrite.org/scans)
[![Contributing Guide](https://img.shields.io/badge/Contributing-Guide-informational)](https://github.com/openrewrite/.github/blob/main/CONTRIBUTING.md)
</div>

### What is this?

This project implements a [Rewrite module](https://github.com/openrewrite/rewrite) that performs common tasks when
migrating to a new version of either Java and/or J2EE.

Browse [a selection of recipes available through this module in the recipe catalog](https://docs.openrewrite.org/recipes/java/migrate).

# Java

In releases of Java prior to Java 11, it was not uncommon for there to be more than three years between major releases
of the platform. This changed in June 2018 when a new, six-month release cadence was adopted by the OpenJDK community.
The new model allows features to be released within any six-month window allowing features to be incrementally
introduced when they are ready. Additionally, there are Java LTS (Long term support) releases on which there exists
enterprise support offered through several vendors that provide builds of the JVM, compiler, and standard libraries. The
current LTS versions of the Java platform (Java 8, 11, 17, and 21) are the most common versions in use within the Java
ecosystem.

# Java EE/Jakarta EE

The Java Platform, Enterprise Edition (Java EE) consists of a set of specifications that extend Java Standard Edition to
enable the development of distributed applications and web services. Examples of the most commonly used parts of Java EE
include JAXB, JAX-WS, and the activation framework. These APIs and their associated reference implementations were
bundled with the Java standard library in JDK 6 through JDK 8 and deprecated in JDK 9. Starting with JDK 11, the
libraries were removed from the standard library to reduce the footprint of the Java standard
library ([See JEP 320 for details](https://openjdk.org/jeps/320)).

**Any projects that continue to use the JAXB framework (on JDK 11+) must now explicitly add the JAXB API and a runtime
implementation to their builds.**

To muddy the waters further, the governance of the Java Platform, Enterprise Edition, was transferred to the Eclipse
Foundation and was renamed to Jakarta EE. The Jakarta EE 8 release (the first under the Jakarta name) maintains
the `javax.xml.bind` package namespace, whereas Jakarta EE 9 is the first release where the package namespace was changed
to `jakarta.xml.bind`:

## Java Architecture for XML Binding (JAXB)

Java Architecture for XML Binding (JAXB) provides a framework for mapping XML documents to/from a Java representation of
those documents. The specification/implementation of this library that is bundled with older versions of the JDK was
part of the Java EE specification before it was moved to the Jakarta project. It can be confusing because Java EE 8
and Jakarta EE 8 provides exactly the same specification (they use the same `javax.xml.bind` namespace), and there are
two different reference implementations for the specification.

| Jakarta EE Version | XML Binding Artifact                        | Package Namespace | Description                   |
|--------------------|---------------------------------------------|-------------------|-------------------------------|
| Java EE 8          | javax.xml.bind:jaxb-api:2.3.x               | javax.xml.bind    | JAXB API                      |
| Jakarta EE 8       | com.sun.xml.bind:jaxb-impl:2.3.x            | javax.xml.bind    | JAXB Reference Implementation |
| Jakarta EE 8       | jakarta.xml.bind:jakarta.xml.bind-api:2.3.x | javax.xml.bind    | JAXB API                      |
| Jakarta EE 8       | org.glassfish.jaxb:jaxb-runtime:2.3.x       | javax.xml.bind    | JAXB Reference Implementation |
| Jakarta EE 9       | jakarta.xml.bind:jakarta.xml.bind-api:3.x   | jakarta.xml.bind  | JAXB API                      |
| Jakarta EE 9       | org.glassfish.jaxb:jaxb-runtime:3.x         | jakarta.xml.bind  | JAXB Reference Implementation |

## Java API for XML Web Services (JAX-WS)

Java API for XML Web Services (JAX-WS) provides a framework for building SOAP-based XML web services in Java. This
framework was originally part of the Java Platform, Enterprise Edition (J2EE), and both the API and the reference
implementation was governed as part of the J2EE specification.

| Jakarta EE Version | XML Web Services Artifact               | Package Namespace | Description                     |
|--------------------|-----------------------------------------|-------------------|---------------------------------|
| Java EE 8          | javax.xml.ws:jaxws-api:2.3.1            | javax.jws         | JAX-WS API                      |
| Jakarta EE 8       | jakarta.xml.ws:jakarta.xml.ws-api:2.3.x | javax.jws         | JAX-WS API                      |
| Jakarta EE 8       | com.sun.xml.ws:jaxws-rt:2.3.x           | javax.jws         | JAX-WS Reference Implementation |
| Jakarta EE 9       | jakarta.xml.ws:jakarta.xml.ws-api:2.3.x | jakarta.jws       | JAX-WS API                      |
| Jakarta EE 9       | com.sun.xml.ws:jaxws-rt:2.3.x           | jakarta.jws       | JAX-WS Reference Implementation |

# Java Migration Recipes

OpenRewrite provides a set of recipes that will help developers migrate to either Java 11, Java 17, or Java 21. These LTS
releases are the most common targets for organizations that are looking to modernize their applications.

## Java 11 Migrations

OpenRewrite provides a set of recipes that will help developers migrate to Java 11 when their existing application
workloads are on Java 8 through 10. The biggest obstacles for the move to Java 11 are the introduction of the module
system (in Java 9) and the removal of J2EE libraries that were previously packaged with the core JDK.

The composite recipe for migrating to Java 11 `org.openrewrite.java.migrate.Java8toJava11` will allow developers to
migrate applications that were previously running on Java 8 through 10. This recipe covers the following themes:

- Applications that use any of the Java EE specifications will have those dependencies migrated to Jakarta EE 8.
  Additionally, the migration to Jakarta EE 8 will also add explicit runtime dependencies on those projects that have
  transitive dependencies on the Jakarta EE APIs. **Currently, only Maven-based build files are supported.**
- Applications that use maven plugins for generating source code from XSDs and WSDLs will have their plugins
  updated to use a version of the plugin that is compatible with Java 11.
- Any deprecated APIs in the earlier versions of Java that have a well-defined migration path will be automatically
  applied to an application's sources. The remediations included with this recipe were originally identified using
  a build plugin called [`Jdeprscan`](https://docs.oracle.com/javase/9/tools/jdeprscan.htm).
- Illegal Reflective Access warnings will be logged when an application attempts to use an API that has not been
  publicly exported via the module system. This recipe will upgrade well-known, third-party libraries if they provide
  a version that is compliant with the Java module system. See [Illegal Reflective Access](#IllegalReflectiveAccess) for
  more information.

## Java 17 Migrations

OpenRewrite provides a set of recipes that will help developers migrate to Java 17 when their existing application
workloads are on Java 11 through 16. The composite recipe `org.openrewrite.java.migrate.UpgradeToJava17` will cover the
following themes:

- Any deprecated APIs in the earlier versions of Java that have a well-defined migration path will be automatically
  applied to an application's sources. The remediations included with this recipe were originally identified using
  a build plugin called [`Jdeprscan`](https://docs.oracle.com/javase/9/tools/jdeprscan.htm).
- Illegal Reflective Access errors are fatal in Java 17 and will result in the application terminating when an
  application or a third-party library attempts to access an API that has not been publicly exported via the module
  system. This recipe will upgrade well-known, third-party libraries if they provide a version that is compliant with
  the Java module system.

## Java 21 Migrations

OpenRewrite provides a set of recipes that will help developers migrate to Java 21 when their existing application
workloads are on Java 11 through 20. The composite recipe `org.openrewrite.java.migrate.UpgradeToJava21` will cover the
following themes:

- everything covered by the Java 17 Migration
- initial support for the migration to Sequenced collections

## Illegal Reflective Access<a name="IllegalReflectiveAccess"></a>

The Java module system was introduced in Java 9 and provides a higher-level abstraction for grouping a set of Java
packages and resources along with additional metadata. The metadata is used to identify what services the module
offers, what dependencies the module requires, and provides a mechanism for explicitly defining which module classes are
"visible" to Java classes that are external to the module.

The module system provides strong encapsulation and the core Java libraries, starting with Java 9, have been designed
to use the module specification. The rules of the module system, if strictly enforced, introduce breaking changes to
downstream projects that have not yet adopted the module system. In fact, it is very common for a typical Java
application to have a mix of module-compliant code along with code that is not aware of modules.

Even as Java has reached Java 15, there are a large number of applications and libraries that are not compliant with
the rules defined by the Java module system. Rather than breaking those libraries, the Java runtime has been configured
to allow mixed-use applications. If an application makes an illegal, reflective call to a module's unpublished resource,
a warning will be logged.

The default behavior, starting with Java 11, is to log a warning the first time an illegal access call is made. All
subsequent calls will not be logged, and the warning looks similar to the following:

```log
WARNING: An illegal reflective access operation has occurred
WARNING: Illegal reflective access by com.thoughtworks.xstream.core.util.Fields (file.....)
WARNING: Please consider reporting this to the maintainers of com.thoughtworks.xstream.core.util.Fields
WARNING: Use --illegal-access=warn to enable warnings of further illegal reflective access operations
WARNING: All illegal access operations will be denied in a future release
```

This warning, while valid, produces noise in an organization's logging infrastructure. In Java 17, these types of issues
are now fatal and the application will terminate if such illegal access occurs.

### Suppressing Illegal Reflective Access Exceptions.

In situations where a third-party library does not provide a version that is compliant with the Java module
system, it is possible to suppress these warnings/errors. An application may add `Add-Opens` declarations to its
top-level JAR's manifest:

```xml

<Add-Opens>
    java.base/java.lang java.base/java.util java.base/java.lang.reflect java.base/java.text java.desktop/java.awt.font
</Add-Opens>
```

This solution will suppress the warnings and errors in the deployed artifacts while still surfacing the warning when
developers run the application from their development environments.

**NOTE: You cannot add these directives to a library that is transitively included by an application. The only place the
Java runtime will enforce the suppressions when they are applied to the top-level, executable Jar.**

There are currently no recipes that will automatically apply "<Add-Opens>" directives to Jar manifests.

## Helpful tools

- http://ibm.biz/WAMT4AppBinaries

## Contributing

We appreciate all types of contributions. See the [contributing guide](https://github.com/openrewrite/.github/blob/main/CONTRIBUTING.md) for detailed instructions on how to get started.

### Licensing

For more information about licensing, please visit our [licensing page](https://docs.openrewrite.org/licensing/openrewrite-licensing).

# Java Migration Tool - Go Implementation

A Go-based tool for migrating Java applications between versions, inspired by OpenRewrite's rewrite-migrate-java.

## Overview

This tool provides automated migration capabilities for Java applications, helping developers upgrade from older Java versions to newer ones (Java 8 → 11, 11 → 17, 17 → 21). It handles common migration tasks including:

- **Build Configuration Updates**: Updates Maven and Gradle build files with new Java versions
- **Dependency Migration**: Migrates Java EE dependencies to Jakarta EE
- **API Replacements**: Replaces deprecated APIs with modern equivalents
- **Code Transformations**: Handles package namespace changes and method signature updates

## Features

### Supported Migration Paths
- **Java 8 to 11**: Complete migration including Java EE to Jakarta EE transition
- **Java 11 to 17**: Handles deprecated APIs and reflective access issues
- **Java 17 to 21**: Includes sequenced collections and latest features
- **Custom Versions**: Support for any target Java version

### Key Transformations
- Replace `sun.misc.BASE64Encoder/Decoder` with `java.util.Base64`
- Migrate Java EE packages to Jakarta EE equivalents
- Update Maven compiler plugin configurations
- Update Gradle Java compatibility settings
- Handle illegal reflective access warnings

## Installation

### Prerequisites
- Go 1.21 or later
- Java project with Maven or Gradle build system

### Build from Source
```bash
git clone <repository>
cd rewrite-migrate-java
go build -o java-migrate cmd/migrate/main.go
```

## Usage

### Basic Usage
```bash
# Migrate to Java 17
./java-migrate /path/to/java/project

# Migrate to specific version
./java-migrate -version=11 /path/to/java/project

# Dry run (show changes without applying)
./java-migrate -dry-run /path/to/java/project

# Custom source directory
./java-migrate -src=src/main/java /path/to/java/project
```

### Command Line Options
- `-version`: Target Java version (default: 17)
- `-src`: Source directory to scan (default: "src/main/java")
- `-dry-run`: Show what would be changed without applying changes

### Examples

#### Migrate Java 8 project to Java 11
```bash
./java-migrate -version=11 ./my-java8-project
```

#### Preview changes for Java 17 migration
```bash
./java-migrate -version=17 -dry-run ./my-java-project
```

## Architecture

The tool follows the recipe pattern used by OpenRewrite, with Go-specific adaptations:

```
pkg/
├── recipe/          # Core recipe interfaces and base types
├── java/            # Java source file parsing and representation
├── migrate/         # Migration recipe implementations
└── cmd/migrate/     # CLI application
```

### Key Components

1. **Recipe Interface**: Defines migration operations
2. **TreeVisitor**: Traverses and modifies source files
3. **SourceFile**: Represents Java source files and build files
4. **Migration Recipes**: Specific transformations (e.g., Java version upgrade, API replacement)

## Supported Transformations

### Build File Updates
- **Maven**: Updates `maven.compiler.source`, `maven.compiler.target`, and plugin configurations
- **Gradle**: Updates `sourceCompatibility`, `targetCompatibility`, and toolchain settings

### Code Transformations

#### Base64 Migration
```java
// Before (Java 8)
import sun.misc.BASE64Encoder;
BASE64Encoder encoder = new BASE64Encoder();
String encoded = encoder.encode(data);

// After (Java 8+)
import java.util.Base64;
String encoded = Base64.getEncoder().encodeToString(data);
```

#### Java EE to Jakarta EE
```java
// Before
import javax.persistence.Entity;
import javax.servlet.http.HttpServlet;

// After
import jakarta.persistence.Entity;
import jakarta.servlet.http.HttpServlet;
```

## Recipe System

### Built-in Recipes

#### Version Upgrade Recipes
- `UpgradeJavaVersion`: Updates build files to target specific Java version
- `Java8ToJava11`: Comprehensive Java 8 to 11 migration
- `UpgradeToJava17`: Migration to Java 17 with LTS features
- `UpgradeToJava21`: Latest Java 21 migration

#### API Migration Recipes
- `UseJavaUtilBase64`: Replaces sun.misc Base64 with java.util.Base64
- `JavaEEToJakartaEE`: Migrates Java EE to Jakarta EE packages
- `RemoveDeprecatedAPIs`: Removes usage of deprecated APIs

### Creating Custom Recipes

```go
type MyCustomRecipe struct {
    *recipe.BaseRecipe
}

func (r *MyCustomRecipe) GetVisitor() recipe.TreeVisitor {
    return &MyCustomVisitor{}
}

type MyCustomVisitor struct{}

func (v *MyCustomVisitor) Visit(node recipe.SourceFile, ctx *recipe.ExecutionContext) (recipe.SourceFile, error) {
    // Implement custom transformation logic
    content := node.GetContent()
    
    // Modify content as needed
    modifiedContent := transformContent(content)
    
    if modifiedContent != content {
        return node.WithContent(modifiedContent), nil
    }
    
    return node, nil
}
```

## Comparison with Original Java Implementation

| Feature | Java (OpenRewrite) | Go Implementation |
|---------|-------------------|-------------------|
| AST Parsing | Full Java AST | Regex-based parsing |
| Template System | JavaTemplate | String replacements |
| Type Safety | Full type checking | Pattern matching |
| Performance | JVM overhead | Native binary |
| Extensibility | Rich plugin system | Go interfaces |
| Complexity | High | Simplified |

## Limitations

- **Parsing**: Uses regex-based parsing instead of full AST
- **Type Checking**: Limited type safety compared to OpenRewrite
- **Template System**: Simple string replacements vs. sophisticated templates
- **Coverage**: Subset of OpenRewrite's transformation capabilities

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests for new functionality
5. Submit a pull request

### Adding New Recipes

1. Create recipe struct implementing `recipe.Recipe`
2. Implement `TreeVisitor` for transformations
3. Add to appropriate migration composite recipe
4. Add tests and documentation

## Testing

```bash
# Run tests
go test ./...

# Test with a sample project
go run cmd/migrate/main.go -dry-run ./testdata/sample-project
```

## License

Licensed under the Moderne Source Available License (same as original).

## Acknowledgments

This Go implementation is inspired by and based on the [OpenRewrite](https://github.com/openrewrite/rewrite) project and specifically [rewrite-migrate-java](https://github.com/openrewrite/rewrite-migrate-java). The original Java implementation provides the conceptual foundation and migration patterns implemented here in Go.
