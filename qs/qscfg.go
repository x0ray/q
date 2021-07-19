// Package qs - q scripting language
package qs

import (
	"os"
	"runtime"
	"time"

	"github.com/x0ray/q/logwriter"

	"github.com/rs/zerolog"
)

const LNumberBit = 64
const LNumberScanFormat = "%f"

type LNumber float64

var (
	CompatVarArg     bool   = true
	FieldsPerFlush   int    = 50
	RegistrySize     int    = 256 * 20
	CallStackSize    int    = 256
	MaxOAListGetLoop int    = 100
	MaxArrayIndex    int    = 67108864
	QsPath           string = "QS_PATH"
	QsLDir           string
	QsPathDefault    string
	QsOS             string   = runtime.GOOS
	oaArgs           []string // q args before -- arg
	scrArgs          []string // script args after -- arg
	QsEmbedded       bool     // used by embedded proc, default is embedded, q sets to false

	log   zerolog.Logger
	debug bool
)

func init() {
	logWriter := logwriter.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	log = zerolog.New(logWriter).With().Caller().Timestamp().Logger()

	// assert q command embedded
	QsEmbedded = true

	// TODO fix location of QsLDir
	if QsOS == "linux" || QsOS == "darwin" { // unix-like
		QsLDir = "/usr/local/share/q"
		QsPathDefault = "./?.q;" + QsLDir + "/?.q;" + QsLDir + "/?/init.q"
		if debug {
			log.Debug().Msgf("Detected OS %s", QsOS)
		}
	} else if QsOS == "windows" {
		QsLDir = "!\\q"
		QsPathDefault = ".\\?.q;" + QsLDir + "\\?.q;" + QsLDir + "\\?\\init.q"
		if debug {
			log.Debug().Msgf("Detected OS %s", QsOS)
		}
	} else {
		log.Error().Msgf("OS %s not supported", QsOS)
	}

	// seperate Args before -- and after --
	scrArgFlg := false
	for _, v := range os.Args {
		if scrArgFlg {
			scrArgs = append(scrArgs, v)
		} else {
			if v == "--" {
				scrArgFlg = true
			} else {
				oaArgs = append(oaArgs, v)
			}
		}
	}
}

func GetDebug() bool {
	return debug
}

func SetDebug(d bool) {
	debug = d
}
