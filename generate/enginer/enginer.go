package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	typeName               = flag.String("type", "", "type name; must be set")
	output                 = flag.String("output", "", "output file name; default srcdir/<type>_enginer.go")
	trimprefix             = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	buildTags              = flag.String("tags", "", "comma-separated list of build tags to apply")
	packageName            = flag.String("packagename", "", "name of the package for generated code; default current package")
	interfacesImport       = flag.String("interfacesimport", "github.com/hearchco/agent/src/search/scraper", "source of the interface import, which is prefixed to interfaces; default github.com/hearchco/agent/src/search/scraper")
	interfacesPackage      = flag.String("interfacespackage", "scraper", "name of the package for the interfaces; default scraper")
	interfaceEnginer       = flag.String("interfaceenginer", "Enginer", "name of the nginer interface; default scraper.Enginer")
	interfaceWebSearcher   = flag.String("interfacewebsearcher", "WebSearcher", "name of the web searcher interface; default scraper.WebSearcher")
	interfaceImageSearcher = flag.String("interfaceimagesearcher", "ImageSearcher", "name of the image searcher interface; default scraper.ImageSearcher")
	interfaceSuggester     = flag.String("interfacesuggester", "Suggester", "name of the suggester interface; default scraper.Suggester")
	enginesImport          = flag.String("enginesimport", "github.com/hearchco/agent/src/search/engines", "source of the engines import, which is prefixed to imports for engines; default github.com/hearchco/agent/src/search/engines")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of enginer:\n")
	fmt.Fprintf(os.Stderr, "\tenginer [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tenginer [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("enginer: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	/* ----------------------------------
	//! Should be comma seperated list of type names, currently is only the first type name
	   ---------------------------------- */
	types := strings.Split(*typeName, ",")
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	g := Generator{
		trimPrefix: *trimprefix,
	}

	if len(args) == 1 && isDirectoryFatal(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
			// ^FATAL
		}
		dir = path.Dir(args[0])
	}

	g.parsePackage(args, tags)

	// Print the header and package clause.
	g.Printf("// Code generated by \"enginer %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	var pkgName string
	if *packageName == "" {
		pkgName = g.pkg.name
	} else {
		pkgName = *packageName
	}
	g.Printf("package %s", pkgName)
	g.Printf("\n")
	g.Printf("import \"%s\"\n", *interfacesImport)
	g.Printf("import \"%s\"\n", *enginesImport)

	// Run generate for each type.
	for _, typeName := range types {
		g.generate(typeName)
	}

	// Format the output.
	src := g.format()

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_enginer.go", types[0])
		outputName = path.Join(dir, strings.ToLower(baseName))
	}
	err := os.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
		// ^FATAL
	}
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
		Logf:       g.logf,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
		// ^FATAL
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages matching %v", len(pkgs), strings.Join(patterns, " "))
		// ^FATAL
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file:       file,
			pkg:        g.pkg,
			trimPrefix: g.trimPrefix,
		}
	}
}

// generate produces imports and the NewEngineStarter method for the named type.
func (g *Generator) generate(typeName string) {
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
		// ^FATAL
	}

	// Generate code for importing engines
	for _, v := range values {
		if validConst(v) && validInterfacer(v, *interfaceEnginer) {
			g.Printf("import \"%s/%s\"\n", *enginesImport, strings.ToLower(v.name))
		}
	}

	// Generate code that will fail if the constants change value.
	g.Printf("func _() {\n")
	g.Printf("\t// An \"invalid array index\" compiler error signifies that the constant values have changed.\n")
	g.Printf("\t// Re-run the enginer command to generate them again.\n")
	g.Printf("\tvar x [1]struct{}\n")
	for _, v := range values {
		origName := v.originalName
		if *packageName != "" {
			origName = fmt.Sprintf("%s.%s", g.pkg.name, v.originalName)
		}
		g.Printf("\t_ = x[%s - (%s)]\n", origName, v.str)
	}
	g.Printf("}\n")

	g.printEnginerLen(values)
	g.printInterfaces(values, *interfaceEnginer)
	g.printInterfaces(values, *interfaceWebSearcher)
	g.printInterfaces(values, *interfaceImageSearcher)
	g.printInterfaces(values, *interfaceSuggester)
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

func (v *Value) String() string {
	return v.str
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// We only care about const declarations.
		return true
	}
	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		if vspec.Type == nil && len(vspec.Values) > 0 {
			// "X = 1". With no type but a value. If the constant is untyped,
			// skip this vspec and reset the remembered type.
			typ = ""

			// If this is a simple type conversion, remember the type.
			// We don't mind if this is actually a call; a qualified call won't
			// be matched (that will be SelectorExpr, not Ident), and only unusual
			// situations will result in a function call that appears to be
			// a type conversion.
			ce, ok := vspec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			id, ok := ce.Fun.(*ast.Ident)
			if !ok {
				continue
			}
			typ = id.Name
		}
		if vspec.Type != nil {
			// "X T". We have a type. Remember it.
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != f.typeName {
			// This is not the type we're looking for.
			continue
		}
		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}
			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the name, find its
			// types.Const, and extract its value.
			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
				// ^FATAL
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				log.Fatalf("can't handle non-integer constant type %s", typ)
				// ^FATAL
			}
			value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
			if value.Kind() != constant.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
				// ^FATAL
			}
			i64, isInt := constant.Int64Val(value)
			u64, isUint := constant.Uint64Val(value)
			if !isInt && !isUint {
				log.Fatalf("internal error: value of %s is not an integer: %s", name, value.String())
				// ^FATAL
			}
			if !isInt {
				u64 = uint64(i64)
			}
			v := Value{
				originalName: name.Name,
				value:        u64,
				signed:       info&types.IsUnsigned == 0,
				str:          value.String(),
			}
			v.name = strings.TrimPrefix(v.originalName, f.trimPrefix)
			if c := vspec.Comment; c != nil && len(c.List) == 1 {
				v.interfaces = strings.Split(strings.TrimSpace(c.Text()), ",")
			}
			f.values = append(f.values, v)
		}
	}
	return false
}

func (g *Generator) printEnginerLen(values []Value) {
	g.Printf("\n")
	g.Printf("\nconst enginerLen = %d", len(values))
	g.Printf("\n")
}

func (g *Generator) printInterfaces(values []Value, interfaceName string) {
	g.Printf("\n")
	g.Printf("\nfunc %sArray() [enginerLen]%s.%s {", toLowerFirstChar(interfaceName), *interfacesPackage, interfaceName)
	g.Printf("\n\tvar engineArray [enginerLen]%s.%s", *interfacesPackage, interfaceName)
	for _, v := range values {
		if validConst(v) && validInterfacer(v, interfaceName) {
			g.Printf("\n\tengineArray[%s.%s] = %s.New()", g.pkg.name, v.name, strings.ToLower(v.name))
		}
	}
	g.Printf("\n\treturn engineArray")
	g.Printf("\n}")
}

func toLowerFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
