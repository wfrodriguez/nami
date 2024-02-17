package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	fzf "github.com/ktr0731/go-fuzzyfinder"
	"gitlab.com/tozd/go/errors"
)

const GoExt = ".go"

var Version = "dev"
var CommitHash = "dev"
var CommitDate = "-now-"

var reDel = regexp.MustCompile(`/[^/]+$`)
var dbg *bool

// check Muestra la información relacionada con un error
func check(err error) { // {{{
	if err != nil {
		fmt.Println("¡¡¡Ha ocurrido un error!!!")
		fmt.Println(err)
		fmt.Println("Causa:", errors.Cause(err))
		d := errors.Details(err)
		if len(d) > 0 {
			fmt.Println("Detalles:")
			for k, v := range d {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}
		if *dbg {
			fmt.Println("Stacktrace:")
			fmt.Printf("  %+v\n", err)
		}
		os.Exit(1)
	}
} // }}}

// main Funcion principal
func main() { // {{{
	dbg = flag.Bool("d", false, "Habilitar modo debug")
	ver := flag.Bool("v", false, "Mostrar la versión")
	flag.Parse()

	if *ver {
		fmt.Printf("%s (build %s @ %s)", Version, CommitHash, CommitDate)
		os.Exit(0)
	}

	if *dbg {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	pkgs, e := GetPackages()
	check(e)

	for _, pkg := range pkgs {
		pkgpath := Join(GoPath(), pkg.String())
		slog.Debug(fmt.Sprint(" Paquete ", pkg))
		check(FindGoFiles(pkgpath, pkg))
		// fmt.Println("\n" + pkg.Describe())
	}

	idx, err := fzf.Find(
		GoItems,
		func(i int) string {
			return GoItems[i]
		},
		fzf.WithHeader("Documentación de Go"))
	check(err)

	PrintDef(Exec("go", "doc", GoItems[idx]))
} // }}}
