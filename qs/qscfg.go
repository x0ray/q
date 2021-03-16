// Package qs - q scripting language
package qs

import (
	"os"
	"runtime"
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
)

func init() {
	// assert q command not embedded, running in its own process
	QsEmbedded = true

	// TODO fix location of QsLDir
	if QsOS == "linux" || QsOS == "darwin" { // unix-like
		QsLDir = "/usr/local/share/q"
		QsPathDefault = "./?.q;" + QsLDir + "/?.q;" + QsLDir + "/?/init.q"
		lgd("detectedOs", "os", QsOS, "libdir", QsLDir, "path", QsPathDefault)
	} else if QsOS == "windows" {
		QsLDir = "!\\q"
		QsPathDefault = ".\\?.q;" + QsLDir + "\\?.q;" + QsLDir + "\\?\\init.q"
		lgd("detectedOs", "os", QsOS, "libdir", QsLDir, "path", QsPathDefault)
	} else {
		lge("osNotSupported", "os", QsOS)
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
