package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.com/tozd/go/errors"
	"golang.org/x/mod/modfile"
)

type Package struct {
	Pkg       string
	Version   string
	Functions []string
	Methods   []string
	Types     []string
}

// NewPackage crea una nueva instancia de Package
func NewPackage(pkg, ver string) *Package { // {{{
	return &Package{
		Pkg:     pkg,
		Version: ver,
	}
} // }}}

// String convierte Package a string
func (p *Package) String() string { //{{{
	return fmt.Sprintf("%s@%s", p.Pkg, p.Version)
} //}}}

// GetPackages devuelve todos los paquetes listados en el archivo `go.mod`
func GetPackages() ([]*Package, errors.E) { // {{{
	cwd := GetCurrentDir()
	goModPath := Join(cwd, "go.mod")
	if !Exists(goModPath) {
		return nil, errors.Errorf("Este directorio (`%s`) no contiene el archivo `go.mod`", cwd)
	}
	modPath := Join(GetCurrentDir(), "go.mod")
	modContent, err := os.ReadFile(modPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error al leer el archivo `go.mod`")
	}

	// Parsea el contenido del archivo go.mod
	modFile, err := modfile.Parse("go.mod", modContent, nil)
	if err != nil {
		log.Fatal(err)
	}

	pkgs := make([]*Package, 0)
	for _, req := range modFile.Require {
		if !req.Indirect {
			pkgs = append(pkgs, NewPackage(req.Mod.Path, req.Mod.Version))
		}
	}

	return pkgs, nil
} // }}}

// FindGoFiles busca todos los archivos .go en el directorio `path` y almacena los nombres de funciones, metodos y tipos en `pkg`
func FindGoFiles(path string, pkg *Package) errors.E { // {{{
	slog.Debug(fmt.Sprintf(" 󰥨 Buscando archivos .go en %s", path))
	err := filepath.Walk(path, func(p string, fi os.FileInfo, er error) error {
		if er != nil {
			return er
		}

		if !fi.IsDir() && filepath.Ext(p) == GoExt && !strings.HasSuffix(p, "_test.go") &&
			!strings.Contains(p, "example") && !strings.Contains(p, "vendor") &&
			!strings.Contains(p, "internal") && !strings.Contains(p, "builtin") {

			slog.Debug(fmt.Sprintf("  󰟓  %s\n", p))
			code, e := ReadFile(p)
			if e != nil {
				return errors.WithDetails(e, "file", p)
			}

			prefix := strings.TrimPrefix(p, path)
			prefix = reDel.ReplaceAllString(prefix, "")

			slog.Debug(fmt.Sprintf("tozd: %q", prefix))
			pkg.Types = append(pkg.Types, getTypes(pkg.Pkg, prefix, code)...)
			pkg.Methods = append(pkg.Methods, getMethods(pkg.Pkg, prefix, code)...)
			pkg.Functions = append(pkg.Functions, getFunctions(pkg.Pkg, prefix, code)...)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "Error al buscar archivos .go")
	}

	GoItems = append(GoItems, pkg.Types...)
	GoItems = append(GoItems, pkg.Methods...)
	GoItems = append(GoItems, pkg.Functions...)

	return nil
} // }}}

// getFunctions extrae las funciones de un archivo .go
func getFunctions(pkg, prefix, text string) []string { // {{{
	re := regexp.MustCompile(`(?m)^func ([A-Z][^\(]+)`)
	matches := re.FindAllStringSubmatch(text, -1)
	fns := make([]string, 0, len(matches))
	// slog.Debug(fmt.Sprintf("Package: %q, Prefix: %q, Matches: %q", pkg, prefix, matches))
	for _, l := range matches {
		f := pkg + prefix + "." + l[1]
		slog.Debug(fmt.Sprintf("     󰊕 %s", f))
		fns = append(fns, f)
	}

	return fns
} // }}}

// getMethods extrae los métodos de un archivo .go
func getMethods(pkg, prefix, text string) []string { // {{{
	re := regexp.MustCompile(`(?m)^func \([a-z]+ \*?([A-Z][^\)]+)\) ([A-Z][^\(]+)`)
	matches := re.FindAllStringSubmatch(text, -1)
	meths := make([]string, 0, len(matches))
	for _, l := range matches {
		m := fmt.Sprintf("%.s.%s.%s", prefix, l[1], l[2])
		slog.Debug(fmt.Sprintf("     󰡱 %s", pkg+m))
		meths = append(meths, pkg+m)
	}

	return meths
} // }}}

// getTypes extrae los tipos de un archivo .go
func getTypes(pkg, prefix, text string) []string { // {{{
	re := regexp.MustCompile(`(?m)^type ([A-Z]\w+)\b`)
	matches := re.FindAllStringSubmatch(text, -1)
	types := make([]string, 0, len(matches))
	prefix = pkg + prefix
	for _, l := range matches {
		t := prefix + "." + l[1]
		slog.Debug(fmt.Sprintf("      %s", t))
		types = append(types, t)
	}

	return types
} // }}}
