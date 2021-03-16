// Package ar parses a command parameter string
package ar

/*
Parses a command parameter string into a slice of strings and also
a map of flags and argumants. This package can split a command
parameter string into a slice of strings which is suitable for
passing to exec.Command()
	see:  https://golang.org/pkg/os/exec/#Command

Notes:     Flag option names:
             Flag option names have the following form:
               -flagname   = any option type, usually used for single letter
                 option names, often used for interactive commands.
               --flagname  = any option type, usually used for long option names
                 which are usually used for batch or commands suited for scripts.
               +flagname   = boolean option type with value of true, this is a
                 acknowledged short hand for -flagname=true
               !flagname   = boolean option type with value of false, this is a
                 acknowledged short hand for -flagname=false
             The name can be any alphanumeric combination of characters. Upper
             and lower case characters are considdered different.

           Arguments:
             An argument is any value in the input string that is not preceeded
			 by a flag option name.

           String option value pair:
             A string value is any quoted value or any non-special and non-whitespace
             character sequence that is not a number or true or false

           Numeric options:
             A numeric value is any character sequence that is numbers or float
             numbers. Leading negative or positive symbols are not allowed, and
             if needed must be passed in quotes, ie "-58"

           Boolean options:
             A boolean value is specified in the following ways:
               -flagname -otherflagname
               --flagname -otherflagname
               +flagname
               !flagname
               -flagname=false     The must not be any blanks next to the "=" character.
               -flagname=true
               --flagname=false
               --flagname=true

           Following are 10 examples of input parameter strings and their
           parsed outputs available from GetList and GetMap

		Input string: -in the --cap 'Library, or watch-devil "in the details"' -hyp "molly-coddle" -the  video --fox 55 -stage 'Tutorial.'
			Parsed Option Flags:
			0 : -in
			1 : the
			2 : --cap
			3 : Library, or watch-devil "in the details"
			4 : -hyp
			5 : molly-coddle
			[stage]:Tutorial.
			[in]:the
			[cap]:Library, or watch-devil "in the details"
			[hyp]:molly-coddle
			[the]:video
			[fox]:55

		Input string: arg/stein.txt -freaky "http://pig.dog.org/#bird" -file "../yes.man" -float 42.012 -fs "ext4 xfs" -math "55.0 + 92 / 68"
			Parsed Option Flags:
			0 : arg/stein.txt
			1 : -freaky
			2 : http://pig.dog.org/#bird
			3 : -file
			4 : ../yes.man
			5 : -float
			6 : 42.012
			7 : -fs
			8 : ext4 xfs
			9 : -math
			10 : 55.0 + 92 / 68

			Parsed Option Map:
			[fs]:ext4 xfs
			[math]:55.0 + 92 / 68
			[#0]:arg/stein.txt
			[freaky]:http://pig.dog.org/#bird
			[file]:../yes.man
			[float]:42.012

		Input string: "bob@mail.com" -path /here:/there:/every/where -pound "#" -empty "" -tilde "~" -bang "!" -inq 'fr og' -apos "pilgrim's"
			Parsed Option Flags:
			0 : bob@mail.com
			1 : -path
			2 : /here:/there:/every/where
			3 : -pound
			4 : #
			5 : -empty
			6 :
			7 : -tilde
			8 : ~
			9 : -bang
			10 : !
			11 : -inq
			12 : fr og
			13 : -apos
			14 : pilgrim's

			Parsed Option Map:
			[apos]:pilgrim's
			[#0]:bob@mail.com5
			[path]:/here:/there:/every/where
			[pound]:#
			[empty]:
			[tilde]:~
			[bang]:!
			[inq]:fr og

		Input string: infile.txt outfile.txt -keep 'red green white gray' -fix 44 -upcase -quote "'" -quiet -- gparm -right 88 -up 48.3 -down 22.5
			Parsed Option Flags:
			0 : infile.txt
			1 : outfile.txt
			2 : -keep
			3 : red green white gray
			4 : -fix
			5 : 44
			6 : -upcase
			7 : -quote
			8 : '
			9 : -quiet
			10 : --
			11 : gparm
			12 : -right
			13 : 88
			14 : -up
			15 : 48.3
			16 : -down
			17 : 22.5

			Parsed Option Map:
			[keep]:red green white gray
			[fix]:44
			[upcase]:
			[quote]:'
			[quiet]:
			[#11]:gparm
			[right]:88
			[#1]:outfile.txt
			[down]:22.5
			[up]:48.3
			[#0]:infile.txt

		Input string: 1234 -22.6 5678,334 $58.99 -555 58% -333.22 (77) TT[99] map(66) {curl} possable values here
			Parsed Option Flags:
			0 : 1234
			1 : -22.6
			2 : 5678,334
			3 : $58.99
			4 : -555
			5 : 58%
			6 : -333.22
			7 : (77)
			8 : TT[99]
			9 : map(66)
			10 : {curl}
			11 : possable
			12 : values
			13 : here

			Parsed Option Map:
			[333.22]:(77)
			[#10]:{curl}
			[#12]:values
			[#13]:here
			[22.6]:5678,334
			[#3]:$58.99
			[555]:58%
			[#8]:TT[99]
			[#9]:map(66)
			[#11]:possable
			[#0]:1234

		Input string: 456 'abc' +246 "jkl" 1.56e-28 -load+runner "0.0.34" !name -f high=low -hamer=false --revenge=true -size 79 -w 92.3 -c red
			Parsed Option Flags:
			0 : 456
			1 : abc
			2 : +246
			3 : jkl
			4 : 1.56e-28
			5 : -load+runner
			6 : 0.0.34
			7 : !name
			8 : -f
			9 : high=low
			10 : -hamer=false
			11 : --revenge=true
			12 : -size
			13 : 79
			14 : -w
			15 : 92.3
			16 : -c
			17 : red

			Parsed Option Map:
			[w]:92.3
			[c]:red
			[#4]:1.56e-28
			[f]:high=low
			[size]:79
			[loadrunner]:0.0.34
			[name]:
			[hamer]:false
			[revenge]:true
			[#0]:456
			[#1]:abc
			[246]:jkl

		Input string: -name456 88 -char99 "ninety nine" -99 'bottles'
			Parsed Option Flags:
			0 : -name456
			1 : 88
			2 : -char99
			3 : ninety nine
			4 : -99
			5 : bottles

			Parsed Option Map:
			[name456]:88
			[char99]:ninety nine
			[99]:bottles

		Input string: -bad.name 49 -my-name "bill smith" -no$man 0 -what/k
			Parsed Option Flags:
			0 : -bad.name
			1 : 49
			2 : -my-name
			3 : bill smith
			4 : -no$man
			5 : 0
			6 : -what/k

			Parsed Option Map:
			[myname]:bill smith
			[no$man]:0
			[what/k]:
			[bad.name]:49

		Input string: "Mary had a little" "it\"s flees were" 'Wise men came from a fire' 'Harry\'s name was mud' "for the MCP"
			Parsed Option Flags:
			0 : Mary had a little
			1 : it\"s flees were
			2 : Wise men came from a fire
			3 : Harry\'s name was mud
			4 : for the MCP

			Parsed Option Map:
			[#3]:Harry\'s name was mud
			[#4]:for the MCP
			[#0]:Mary had a little
			[#1]:it\"s flees were
			[#2]:Wise men came from a fire

			Input string: check process -format pretty -debug=true
			Parsed Option Flags:
			0 : check
			1 : process
			2 : -format
			3 : pretty
			4 : -debug=true

			Parsed Option Map:
			[debug]:true
			[#0]:check
			[#1]:process
			[format]:pretty

*/

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// global variables
type Args struct {
	om     map[string]string // map of flags/#args and values
	out    []string          // slice of all arg,flag,values
	argcnt int               // number of args
	inp    string            // input arg string
}

// emitFlag - parses a command parameter string into a slice
func (a *Args) emitFlag(flag string, arg string) error {
	var err error
	if flag == "" {
		a.om["#"+strconv.Itoa(a.argcnt)] = arg // add flag to flag map
	} else if flag == "--" {
		a.om["#"+strconv.Itoa(a.argcnt)] = arg // add flag to flag map
	} else {
		if len(flag) < 2 {
			err = errors.New("One character option: " + flag + " not valid.")
			return err
		} else {
			f := flag[1:] // remove - + !
			if len(f) > 0 {
				if f[0:1] == "-" {
					f = f[1:] // remove remaining -
				}
				if len(f) > 0 {
					if strings.Contains(f, "=") {
						fs := strings.Split(f, "=")
						a.om[fs[0]] = fs[1] // add flag to flag map
					} else {
						a.om[f] = arg // add flag to flag map
					}
				} else {
					err = errors.New("Missing flag name.")
					return err
				}
			}
		}
	}
	return err
}

// emitArg - parses a command parameter string into a slice
func (a *Args) emitArg(flag string, arg string) error {
	var err error
	a.out = append(a.out, arg)
	if !strings.Contains(flag, "=") {
		err = a.emitFlag(flag, arg)
	}
	a.argcnt++
	return err
}

// parmParse - parses a command parameter string into a slice
//   This function splits the command parameter string into a slice of
//   strings which is suitable for passing to exec.Command()
//     see:  https://golang.org/pkg/os/exec/#Command
//   The Go compiler scanner is used to make this a simpler matter, and
//   allows for nested strings, strings with blanks. The \ escape character
//   is inserted in strings when necessary.
func (a *Args) ParseArg(in string) error {
	var err error
	// init global vbalues
	a.inp = in
	a.out = make([]string, 0)
	a.argcnt = 0
	a.om = make(map[string]string)

	arg := ""  // the current arg being collected
	flag := "" // current flag, reset after an arg
	// only one or none of these 3 following can be true
	inarg := false  // an arg is an un quoted value for a flag or a positional arg
	inqstr := false // in a quoted string: "...'d'..." or '..."k"...' or "..\"..." or '..\'..'
	inflag := false // in a flag (starts with - -- + !

	eqchar := ""      // end of quoted string character `"` or `'` or ``
	c := ""           // current character being examined
	lastc := ""       // last character examined
	lin := len(a.inp) // length of input
	p := 0            // position in input string
	for p < lin {
		c = a.inp[p : p+1]

		if c == " " { // arg end ?
			if inarg { // end of an arg
				err = a.emitArg(flag, arg)
				if err != nil {
					return err
				}
				arg = ""
				flag = ""
				inarg = false
				inflag = false
			} else if inqstr {
				arg = arg + c
			} else if inflag { // end of flagname, including -bool=true style
				a.out = append(a.out, flag)
				a.argcnt++
				inflag = false
			}

		} else if (c == "-") || (c == "+") || (c == "!") { // flag start ?
			if inarg {
				arg = arg + c
			} else if inqstr {
				arg = arg + c
			} else if inflag {
				flag = flag + c
			} else { // new flag
				if flag != "" { // got new flag with no arg before emitting last flag?
					err = a.emitFlag(flag, "") // emit last flag now
				}
				inflag = true
				flag = c
			}

		} else if (c == `"`) || (c == `'`) {
			if inqstr {
				if c == eqchar && lastc != `\` { // end qstr
					// emit the qstr as an arg
					err = a.emitArg(flag, arg)
					if err != nil {
						return err
					}
					arg = ""
					flag = ""
					inarg = false
					inflag = false
					// indicate no longer in qstr
					inqstr = false
					eqchar = ""
				} else {
					arg = arg + c
				}
			} else { // new qstr
				eqchar = c
				inqstr = true
			}

		} else { // accumulate arg / flag chars
			if inflag {
				flag = flag + c
			} else if inqstr {
				arg = arg + c
			} else {
				inarg = true
				arg = arg + c
			}
		}

		lastc = c
		p++
	} // of for

	// catch last arg or flag
	if inarg {
		err = a.emitArg(flag, arg)
	} else if inflag {
		a.out = append(a.out, flag)
		err = a.emitFlag(flag, arg)
	}

	return err
}

// GetMap - returns a map of flags and their values, allowing for direct named
//   flag value access. If no preceeding --flagname -flagname +flagname or
//   !flagname is specified a value is assumed to be just an argument and its
//   map name becomes '#nnn' where nnn is the positional number of the argument
//   relative to the other argumants, not including flags and their values.
func (a *Args) GetMap() map[string]string {
	return a.om
}

// GetList - returns a slice of args suitable for calling os/exec.Command(name, args)
func (a *Args) GetList() []string {
	return a.out
}

// Print - prints the input string
func (a *Args) Print() {
	if len(a.out) > 0 {
		fmt.Printf("Input string: %s\n", a.inp)
	}
}

// PrtOpts - prints the slice of options and arguments
func (a *Args) PrtOpts() {
	if len(a.out) > 0 {
		fmt.Printf("Parsed Option Flags:\n")
		for i, v := range a.out {
			fmt.Println(i, ":", v)
		}
		fmt.Printf("\n")
	}
}

// PrtOpMap - prints the map of flags and their values including #nnn arguments
func (a *Args) PrtOpMap() {
	if len(a.om) > 0 {
		fmt.Printf("Parsed Option Map:\n")
		for k, v := range a.om {
			fmt.Printf("[%s]:%s\n", k, v)
		}
		fmt.Printf("\n")
	}
}
