Tagline := "Nami - Aplicación CLI para la documentación de Go"

Green  := $(shell tput -Txterm setaf 2)
Yellow := $(shell tput -Txterm setaf 3)
White  := $(shell tput -Txterm setaf 7)
Red    := $(shell tput -Txterm setaf 1)
Reset  := $(shell tput -Txterm sgr0)

# Nombre del binario
BinName = nami
# Ubicación del archivo main
Main = *.go

# Vars

Version := $(git describe --abbrev=0 --tags)
CommitHash := $(shell git rev-parse --short HEAD)
BuildTimestamp := $(shell date '+%Y-%m-%dT%H:%M:%S')

TargetMaxCharNum := 20

.DEFAULT_GOAL := help

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: confirm
confirm:
	@echo -n '¿Seguro? [s/N] ' && read ans && [ $${ans:-N} = s ]

.PHONY: validate
## Valida si los programas obligatorios se encuentran instalados en el sistema
validate:
	@command -v go >/dev/null 2>&1 || { echo "${Red}Requiero el programa 'go' pero no está instalado.${Reset}" >&2; exit 1; }

.PHONY: install
## Instala las aplicaciones auxiliares para el desarrollo de la aplicación
install: validate
	go install go.uber.org/nilaway/cmd/nilaway@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/segmentio/golines@latest
	go install github.com/fatih/gomodifytags@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install github.com/cosmtrek/air@latest

.PHONY: help
## Muestra este mensaje de ayuda
help:
	@echo ${Tagline}
	@echo ''
	@echo 'Modo de uso:'
	@echo '  ${Yellow}make${Reset} ${Green}<target>${Reset}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9\/]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${Yellow}%-$(TargetMaxCharNum)s${Reset} ${Green}%s${Reset}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: tidy
## Formatea el código usando golines y gofumpt, despues actualiza las dependencias usando go mod tidy.
tidy:
	golines -m 120 --base-formatter gofumpt ./...
	go mod tidy -v

.PHONY: audit
# Realiza comprobaciones de control de calidad en el código base, incluida la verificación de las dependencias, la
# ejecución de análisis estáticos, la comprobación de vulnerabilidades, análisis estático para evitar pánicos nil en
# producción detectándolos en tiempo de compilación en lugar de en tiempo de ejecución y la ejecución de todas las pruebas.
## Realiza comprobaciones de control de calidad en el código
audit:
	go mod verify
	go vet ./...
	nilaway ./...
	staticcheck -checks=all,-ST1000,-U1000 -f stylish ./...
	govulncheck ./...
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: no-dirty
## Comprueba que no hay cambios no comprometidos en los archivos 'tracked' en el repositorio git actual.
no-dirty:
	git diff --exit-code

.PHONY: build
## Compila la aplicación y genera un binario en ./dist
build: validate
	@mkdir -p ./dist
	go build -ldflags="-X main.Version='$(Version)' -X main.CommitHash='$(CommitHash)' -X main.BuildTimestamp='$(BuildTimestamp)'" -o ./dist/$(BinName) $(Main)

.PHONY: run/go
## Inicia la aplicación desde el código fuente
run/go: validate
	@go run $(Main) $(Args)

.PHONY: run/bin
## Ejecuta el binario generado
run/bin: build
	@./dist/$(BinName) $(Args)
