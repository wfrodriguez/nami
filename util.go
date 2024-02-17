package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gitlab.com/tozd/go/errors"
)

// ReadFile lee un archivo de texto
func ReadFile(path string) (string, errors.E) { // {{{
	content, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "Error al leer el archivo")
	}

	return string(content), nil
} // }}}

// Exec ejecuta un comando y devuelve la salida, en caso de un error devuelve un valor vacío
func Exec(name string, args ...string) string { //{{{
	cmd, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(cmd))
} //}}}

// GoPath devuelve el directorio GOPATH
func GoPath() string { // {{{
	gp := Exec("go", "env", "GOPATH")

	return Join(gp, "pkg", "mod")
} // }}}

// IF es un alias del condicional ternario: cond ? rt : rf
func IF[T any](cond bool, rt, rf T) T { // {{{
	if cond {
		return rt
	}
	return rf
} // }}}

// Join es un alias de filepath.Join
func Join(p ...string) string { // {{{
	return filepath.Join(p...)
} // }}}

// Exists valida si un archivo `path` existe
func Exists(path string) bool { // {{{
	_, err := os.Stat(path)
	return err == nil
} // }}}

// GetCurrentDir devuelve el directorio actual de trabajo
func GetCurrentDir() string { // {{{
	currentDir, err := os.Getwd()
	if err != nil {
		return "."
	}

	return currentDir
} // }}}
