// package qm stand alone command shell for Q language interpreter
package qm

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/x0ray/q/qs"
	"github.com/x0ray/q/qs/qsp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PGM  = "q"
	VER  = "2.0.0"
	NAME = "Q"
	EXTN = ".q"
)

// codes use by os.Exit()
const (
	RCOK    = 0
	RCWARN  = 1
	RCERROR = 2
	RCFATAL = 3
	NOTEST  = -1.0
	MEMDFLT = 0

	LOQUACITY = 2 // LOQUACITY - INFO msg control, 9=most 0=least
)

var VERDATE string = "20Mar2021"

const (
	nmDebug  = "debug"
	msgDebug = "Display debugging information during OA program execution"

	nmExec  = "exec"
	msgExec = "String of OA language statements to execute directly"

	nmH     = "h"
	nmHelp  = "help"
	msgHelp = "Display help information and exit"

	nmInter  = "inter"
	msgInter = "File name to use for profile data"

	nmLib  = "lib"
	msgLib = "OA library file name to access"

	nmLimit  = "limit"
	msgLimit = "Sets a memory limit for the executing OA program in MB. If 0 is used, no limit is set."

	nmLog  = "log"
	msgLog = "File name to use for logging"

	nmLoquacity  = "loquacity"
	msgLoquacity = "1-9 value, how much information is written to the log. 9 is more 1 is less."

	nmName  = "name"
	msgName = "Name used to name script for logging, and internal reference"

	nmPgm  = "pgm"
	msgPgm = "Name of file containing OA program"

	nmProfile  = "profile"
	msgProfile = "File name to use for profile data"

	nmQuiet  = "quiet"
	msgQuiet = "Hide program output during OA program execution"

	nmV        = "v"
	nmVersion  = "version"
	msgVersion = "Display version information and exit"

	nmVerbose  = "verbose"
	msgVerbose = "Displays more information for many functions that have display output"

	nmSyntax  = "Xsyntax" // X... == extended / hidden option
	msgSyntax = "Debugging tool, displays syntax tree information"

	nmIcode  = "Xicode" // X... == extended / hidden option
	msgIcode = "Debugging tool, displays internal instruction code list"

	nmXHelp  = "Xhelp" // X... == extended / hidden option
	msgXHelp = "Display extended help information and exit"
)

var u struct {
	// pgm - file name containing OA program
	pgm string

	// exec - string of OA language statements to execute directly
	exec string

	// lib - OA library file name to access
	lib string

	// profile - file name to use for profile data
	profile string

	// log - file name to use for logging
	log string

	// loq - loquacity of INFO messages in log
	loq int

	// name - name used for logging
	name string

	// limit - sets a memory size limit for the executing OA program
	limit int

	// inter - use interactive mode
	inter bool

	// syntax - debugging tool, displays syntax tree information
	syntax bool

	// icode - debugging tool, displays internal instruction code list
	icode bool

	// h - (help) show the help then exit
	h bool

	// xhelp show the extended help then exit
	xhelp bool

	// q - (quiet) hide program output
	q bool

	// verbose - display more information for many functions that have display output
	verbose bool

	// v - (version) display version information then exit
	v bool

	// debug - show debug output on log
	debug bool

	// run - not real flag option - indicates run mode
	// Can be in any style from- q [run] [-pgm fn] [fn.oa]	  For example:
	//   q run test.oa
	//   q test.oa
	//   q -pgm test.oa
	//   q run -pgm test.oa
	run bool
}

var (
	flgs    *flag.FlagSet // all command line flags
	status  int           = RCOK
	oaArgs  []string      // q args before -- arg
	scrArgs []string      // script args after -- arg

	// TODO change log writer
	/*
		ilg func(int, string, bool, string, string, string, string, bool,
			bool, int, int) *bytes.Buffer = xl.InitLogWtr

		// alias logging functions
		lgi func(int, string, ...interface{}) = xl.Lgi // INFO
		lgw func(string, ...interface{})      = xl.Lgw // WARN
		lge func(string, ...interface{})      = xl.Lge // ERROR
		lgf func(string, ...interface{})      = xl.Lgf // FATAL
		lgd func(string, ...interface{})      = xl.Lgd // DEBUG

	*/
)

func init() {
	// assert q command not embedded, running in its own process
	qs.QsEmbedded = false

	flgs = flag.NewFlagSet("", flag.ExitOnError)
	flgs.StringVar(&u.pgm, nmPgm, "", msgPgm)                 // file name containing OA program
	flgs.StringVar(&u.exec, nmExec, "", msgExec)              // string of OA language statements to execute directly
	flgs.StringVar(&u.lib, nmLib, "", msgLib)                 // OA library file name to access
	flgs.StringVar(&u.profile, nmProfile, "", msgProfile)     // file name to use for profile data
	flgs.StringVar(&u.log, nmLog, "", msgLog)                 // file name to use for logging
	flgs.StringVar(&u.name, nmName, NAME, msgName)            // name used as tag for logging and internal ref _NAME
	flgs.IntVar(&u.limit, nmLimit, MEMDFLT, msgLimit)         // sets a memory size limit for the executing OA program
	flgs.IntVar(&u.loq, nmLoquacity, LOQUACITY, msgLoquacity) // level of INFO messages to log 0=low .. 9=high
	flgs.BoolVar(&u.inter, nmInter, false, msgInter)          // use interactive mode
	flgs.BoolVar(&u.syntax, nmSyntax, false, msgSyntax)       // debugging tool, displays syntax tree information
	flgs.BoolVar(&u.icode, nmIcode, false, msgIcode)          // debugging tool, displays internal instruction code list

	// commenting the following two lines will enable the standard Golang flag help
	flgs.BoolVar(&u.verbose, nmVerbose, false, msgVerbose) // display more information for many functions that have display output
	flgs.BoolVar(&u.h, nmH, false, msgHelp)                // display help then exit
	flgs.BoolVar(&u.h, nmHelp, false, msgHelp)             //  "
	flgs.BoolVar(&u.xhelp, nmXHelp, false, msgXHelp)       // display extended help then exit
	flgs.BoolVar(&u.v, nmV, false, msgVersion)             // display version then exit
	flgs.BoolVar(&u.v, nmVersion, false, msgVersion)       //  "
	flgs.BoolVar(&u.q, nmQuiet, false, msgQuiet)           // hide program output
	flgs.BoolVar(&u.debug, nmDebug, false, msgDebug)       // show additional debugging information

	flgs.Usage = func() {
		fmt.Println(`Usage: ` + PGM + ` [options] [oa-program [oa-program-args]]
Use ` + PGM + ` -help for more information`)
	}

}

func Main() int {
	if strings.Contains(fmt.Sprintf("%v", os.Args[:]), "-debug") { // preparse for debug
		u.debug = true
		fmt.Printf("%s DEBUG mode, started with %d parameters:\n", PGM, len(os.Args))
		for i, v := range os.Args {
			fmt.Printf("  [%d] %s\n", i, v)
		}
	}

	scrArgFlg := false
	subCmd := ""
	if len(os.Args) > 1 {
		// could have a sub command: run help version xhelp
		subCmd = strings.ToLower(os.Args[1])

		if subCmd == "help" {
			u.h = true
		} else if subCmd == "xhelp" {
			u.xhelp = true
		} else if subCmd == "version" {
			u.v = true
		} else if subCmd == "run" {
			u.run = true
		} else if subCmd == "int" {
			u.inter = true
		} else {
			subCmd = ""
		}

		// seperate Args before -- and after --
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

		// parse the q args
		if subCmd == "" { // first arg not sub command
			u.run = true
			subCmd = "run" // default to run sub command

			// get args

			if !strings.Contains(oaArgs[1], EXTN) && !strings.HasPrefix(oaArgs[1], "-") {
				// probably a non.oa extension #!/q script
				u.pgm = oaArgs[1] // get oa pgm
				if scrArgFlg {
					flgs.Parse(oaArgs[2:]) // extract options
				} else {
					scrArgs = os.Args[1:]
				}

			} else if strings.HasPrefix(oaArgs[1], "-") {
				flgs.Parse(oaArgs[1:]) // extract possible options
			} else {
				u.pgm = oaArgs[1] // get oa pgm

				if len(oaArgs) >= 2 { // might have more options
					flgs.Parse(oaArgs[2:]) // extract options
				}
			}

		} else if subCmd == "int" { // got subcmd - check for -options
			if !strings.HasPrefix(oaArgs[2], "-") {
				u.pgm = os.Args[2] // assume arg 2 is the oa pgm
			}
			if len(oaArgs) >= 2 {
				flgs.Parse(oaArgs[2:]) // extract all additional -options
			}
		} else { // got subcmd - check for -options
			if subCmd == "run" {
				if !strings.HasPrefix(oaArgs[2], "-") {
					u.pgm = os.Args[2] // assume arg 2 is the oa pgm
				}
			}
			if len(oaArgs) >= 2 {
				flgs.Parse(oaArgs[2:]) // extract all additional -options
			}
		}
	} else {
		// no command, args, or options
		u.inter = true
	}

	if u.debug {
		fmt.Printf("subCmd..: %s\n", subCmd)
		fmt.Printf("u.run...: %v\n", u.run)
		fmt.Printf("u.inter.: %v\n", u.inter)
		fmt.Printf("u.pgm...: %s\n", u.pgm)
		fmt.Printf("u.v.....: %v\n", u.v)
		fmt.Printf("u.h.....: %v\n", u.h)
	}

	// The following lines implements over-riding the lame ass standard help printing
	if u.h {
		// print main command blurb
		fmt.Printf("Help for %s Ver %s (%s)\n", PGM, VER, VERDATE)
		if u.verbose {
			fmt.Printf("%s\n", qs.HelpOverViewText)
			qs.VerboseHelp = true
		}
		fmt.Printf("%s\n", qs.HelpText)
		fmt.Printf("%s Option Details:\n", PGM)
		flgs.VisitAll(qs.PrtHelpFlag) // print each flag option
		if u.verbose {
			qs.VerboseHelp = true
			fmt.Printf("%s\n", qs.HelpTextVerbose)
		}
		os.Exit(RCOK)
	}

	// extended help options -X...
	if u.xhelp {
		// print main command blurb
		fmt.Printf("Extended help for %s Ver %s (%s)\n", PGM, VER, VERDATE)
		if u.verbose {
			fmt.Printf("%s\n", qs.HelpOverViewText)
			qs.VerboseHelp = true
		}
		fmt.Printf("%s\n", qs.HelpText)
		fmt.Printf("%s Option Details:\n", PGM)
		flgs.VisitAll(qs.PrtHelpFlag)  // print each flag option
		flgs.VisitAll(qs.PrtXHelpFlag) // print each extenmded flag option
		if u.verbose {
			fmt.Printf("%s\n", qs.HelpTextVerbose)
		}
		os.Exit(RCOK)
	}

	if u.v {
		// Use shared version output:
		fmt.Printf("Program: %s version: %s \n", PGM, VER)
		os.Exit(RCOK)
	}

	// set up logging
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Debug().
		Str("pgm", PGM).Str("ver", VER).Str("verdate", VERDATE).
		Msg("Q log started")

	// is OA program profiling required
	if len(u.profile) != 0 { // do profiling on OA code?
		prof, err := os.Create(u.profile)
		if err != nil {
			log.Error().Err(err).Msg("profile open error")
			os.Exit(RCERROR)
		}
		pprof.StartCPUProfile(prof)
		defer pprof.StopCPUProfile()
	}

	L := qs.NewState()
	defer L.Close()
	// optionally set memory limit in MB
	if u.limit > 0 {
		log.Debug().Int("limit", u.limit).Msg("set mem limit")
		L.SetMx(u.limit)
	}

	// display version information if requested or interactive
	if u.v || u.inter {
		fmt.Printf("%s, Version - %s, Build Date - %s\n", PGM, VER, VERDATE)
	}

	if len(u.lib) > 0 {
		if err := L.DoFile(u.lib); err != nil {
			status = RCWARN
			log.Warn().Err(err).Msg("library open error")
		}
	}

	if u.run {
		// check program extension type
		if !strings.Contains(u.pgm, EXTN) {
			// program with no extension, probably using #!/q
			// use the entire command args set minus the q executable
			scrArgs = os.Args[1:]
		}

		// create environment args list
		nargs := len(scrArgs)
		argtb := L.NewOAList()
		for i := 1; i < nargs; i++ {
			if u.debug {
				fmt.Printf("arg[%d]..: %s\n", i, scrArgs[i])
			}
			L.RawSet(argtb, qs.LNumber(i), qs.LString(scrArgs[i]))
		}
		L.SetGlobal("_NAME", qs.LString(u.name))
		L.SetGlobal("arg", argtb) // pass args to the OA script

		// run script pgm through script debugging tools
		if u.syntax || u.icode {
			file, err := os.Open(u.pgm)
			if err != nil {
				log.Error().Err(err).Str("pgm", u.pgm).Msg("Q srcipt open error")
				return RCERROR
			}
			segment, err2 := qsp.Parse(file, u.pgm)
			if err2 != nil {
				log.Error().Err(err2).Str("pgm", u.pgm).Msg("Q srcipt parse error")
				return RCERROR
			}
			if u.syntax {
				fmt.Println(qsp.Dump(segment))
			}
			if u.icode {
				icode, err3 := qs.Compile(segment, u.pgm)
				if err3 != nil {
					log.Error().Err(err3).Str("pgm", u.pgm).Msg("Q srcipt compile error")
					return RCERROR
				}
				fmt.Println(icode.String())
			}
		}
		// execute oa script from file
		if err := L.DoFile(u.pgm); err != nil {
			log.Error().Err(err).Str("pgm", u.pgm).Msg("Q srcipt execution error")
			return RCERROR
		}
	}

	// execute string of oa script from -exec str option
	if len(u.exec) > 0 {
		if err := L.DoString(u.exec); err != nil {
			log.Error().Err(err).Msg("Q string execution error")
			return RCERROR
		}
	}

	// use interactive mode
	if u.inter {
		// use the entire command args set minus the q executable
		scrArgs = os.Args[1:]
		// create environment args list
		nargs := len(scrArgs)
		argtb := L.NewOAList()
		for i := 1; i < nargs; i++ {
			if u.debug {
				fmt.Printf("arg[%d]..: %s\n", i, scrArgs[i])
			}
			L.RawSet(argtb, qs.LNumber(i), qs.LString(scrArgs[i]))
		}
		L.SetGlobal("arg", argtb) // pass args to the OA script

		fmt.Printf("%s interactive mode, enter exit() to exit, or help() for help. \n", PGM)
		doREPL(L)
	}
	return status
}

// do read/eval/print/loop
func doREPL(L *qs.LState) {
	reader := bufio.NewReader(os.Stdin)
	for {
		if str, err := loadline(reader, L); err == nil {
			if err := L.DoString(str); err != nil {
				status = RCWARN
				log.Error().Err(err).Msg("Q load line do string error")
			}
		} else { // error on loadline
			status = RCWARN
			log.Error().Err(err).Msg("Q load line error")
			return
		}
	}
}

func incomplete(err error) bool {
	if lerr, ok := err.(*qs.ApiError); ok {
		if perr, ok := lerr.Cause.(*qsp.Error); ok {
			return perr.Pos.Line == qsp.EOF
		}
	}
	return false
}

func loadline(reader *bufio.Reader, L *qs.LState) (string, error) {
	fmt.Print("> ")
	if line, err := reader.ReadString('\n'); err == nil {
		if _, err := L.LoadString("return " + line); err == nil { // try add return <...> then compile
			return line, nil
		} else {
			return multiline(line, reader, L)
		}
	} else {
		return "", err
	}
}

func multiline(ml string, reader *bufio.Reader, L *qs.LState) (string, error) {
	for {
		if _, err := L.LoadString(ml); err == nil { // try compile
			return ml, nil
		} else if !incomplete(err) { // syntax error , but not EOF
			return ml, nil
		} else {
			fmt.Print(">> ")
			if line, err := reader.ReadString('\n'); err == nil {
				ml = ml + "\n" + line
			} else {
				return "", err
			}
		}
	}
}
