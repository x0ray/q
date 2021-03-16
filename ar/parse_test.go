package ar

import (
	"testing"
)

// TestParsedValue - Test values returned from getlist function after parse
func TestParsedValue(t *testing.T) {
	var pass int
	t.Log("Test TestParsedValue() function started")

	/* perform tests */
	inp := `-in the --cap 'Library, or watch-devil "in the details"' -empty "" -hyp "molly-coddle" -the  video --fox 55 -stage 'Tutorial.'`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedValue() error: " + err.Error())
	} else {
		sl := args.GetList()
		if sl[1] == "the" {
			pass++
		} else {
			t.Error("Test TestParsedValue() 'the' not at index 1")
		}
		if sl[9] == "video" {
			pass++
		} else {
			t.Error("Test TestParsedValue() 'video' not at index 9")
		}
		if sl[11] == "55" {
			pass++
		} else {
			t.Error("Test TestParsedValue() '55' not at index 11")
		}
	}

	/* publish test success or failure */
	if pass != 3 {
		t.Error("Test TestParsedValue() failed.")
	} else {
		t.Log("Test TestParsedValue() OK.")
	}
}

// ExampleParsedValueDashes - Test when dashes are in the value the value must be quoted
func ExampleParsedValueDashes() {
	inp := `-in the --cap 'Library, or watch-devil "in the details"' -hyp "molly-coddle" -the  video --fox 55 -stage 'Tutorial.'`
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : -in
	// 1 : the
	// 2 : --cap
	// 3 : Library, or watch-devil "in the details"
	// 4 : -hyp
	// 5 : molly-coddle
	// 6 : -the
	// 7 : video
	// 8 : --fox
	// 9 : 55
	// 10 : -stage
	// 11 : Tutorial.
}

// TestParsedValueDashes - Test map values when dashes are in the value the value must be quoted
func TestParsedValueDashes(t *testing.T) {
	var pass int
	t.Log("Test TestParsedValue() function started")

	/* perform tests */
	inp := `-in the --cap 'Library, or watch-devil "in the details"' -empty "" -hyp "molly-coddle" -the  video --fox 55 -stage 'Tutorial.'`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedValue() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["in"] == "the" {
			pass++
		} else {
			t.Error("Test TestParsedValueDashes() 'the' not in opmap['in']")
		}
		if opmap["hyp"] == "molly-coddle" {
			pass++
		} else {
			t.Error("Test TestParsedValueDashes() 'molly-coddle' not in opmap['hyp']")
		}
		if opmap["stage"] == "Tutorial." {
			pass++
		} else {
			t.Error("Test TestParsedValueDashes() 'Tutorial.' not in opmap['stage']")
		}
		if opmap["fox"] == "55" {
			pass++
		} else {
			t.Error("Test TestParsedValueDashes() '55' not in opmap['fox']")
		}
	}

	/* publish test success or failure */
	if pass != 4 {
		t.Error("Test TestParsedValueDashes() failed.")
	} else {
		t.Log("Test TestParsedValueDashes() OK.")
	}
}

// ExampleParsedValueFiles - Test values that look like file names do not need to be quoted unless they have special or dash characters
func ExampleParsedValueFiles() {
	inp := `arg/stein.txt -freaky "http://pig.dog.org/#bird" -file "../yes.man" -float 42.012 -fs "ext4 xfs" -math "55.0 + 92 / 68"`
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : arg/stein.txt
	// 1 : -freaky
	// 2 : http://pig.dog.org/#bird
	// 3 : -file
	// 4 : ../yes.man
	// 5 : -float
	// 6 : 42.012
	// 7 : -fs
	// 8 : ext4 xfs
	// 9 : -math
	// 10 : 55.0 + 92 / 68
}

// TestParsedValueFiles - Test map values that look like file names do not need to be quoted unless they have special or dash characters
func TestParsedValueFiles(t *testing.T) {
	var pass int
	t.Log("Test TestParsedValueFiles() function started")

	/* perform tests */
	inp := `arg/stein.txt -freaky "http://pig.dog.org/#bird" -file "../yes.man" -float 42.012 -fs "ext4 xfs" -math "55.0 + 92 / 68"`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedValueFiles() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["freaky"] == "http://pig.dog.org/#bird" {
			pass++
		} else {
			t.Error("Test TestParsedValueFiles() 'http://pig.dog.org/#bird' not in opmap['freaky']")
		}
		if opmap["file"] == "../yes.man" {
			pass++
		} else {
			t.Error("Test TestParsedValueFiles() '../yes.man' not in opmap['file']")
		}
		if opmap["fs"] == "ext4 xfs" {
			pass++
		} else {
			t.Error("Test TestParsedValueFiles() 'ext4 xfs' not in opmap['fs']")
		}
	}

	/* publish test success or failure */
	if pass != 3 {
		t.Error("Test TestParsedValueFiles() failed.")
	} else {
		t.Log("Test TestParsedValueFiles() OK.")
	}
}

// ExampleParsedValuePathSep - Test the colon path seperator character is allowed in a value other specials need quotes
func ExampleParsedValuePathSep() {
	inp := `"bob@mail.com" -path /here:/there:/every/where -pound "#" -dot "." -tilde "~" -bang "!" -inq 'fr og' -apos "pilgrim's"`
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : bob@mail.com
	// 1 : -path
	// 2 : /here:/there:/every/where
	// 3 : -pound
	// 4 : #
	// 5 : -dot
	// 6 : .
	// 7 : -tilde
	// 8 : ~
	// 9 : -bang
	// 10 : !
	// 11 : -inq
	// 12 : fr og
	// 13 : -apos
	// 14 : pilgrim's
}

// TestParsedValuePathSep - Test map values with colon path seperator character is allowed in a value other specials need quotes
func TestParsedValuePathSep(t *testing.T) {
	var pass int
	t.Log("Test TestParsedValuePathSep() function started")

	/* perform tests */
	inp := `"bob@mail.com" -path /here:/there:/every/where -pound "#" -dot "." -tilde "~" -bang "!" -inq 'fr og' -apos "pilgrim's"`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedValuePathSep() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["#0"] == "bob@mail.com" {
			pass++
		} else {
			t.Error("Test TestParsedValuePathSep() 'bob@mail.com' not in opmap['#1']")
		}
		if opmap["path"] == "/here:/there:/every/where" {
			pass++
		} else {
			t.Error("Test TestParsedValuePathSep() '/here:/there:/every/where' not in opmap['path']")
		}
		if opmap["apos"] == "pilgrim's" {
			pass++
		} else {
			t.Error("Test TestParsedValuePathSep() 'pilgrim's' not in opmap['apos']")
		}
	}

	/* publish test success or failure */
	if pass != 3 {
		t.Error("Test TestParsedValuePathSep() failed.")
	} else {
		t.Log("Test TestParsedValuePathSep() OK.")
	}
}

// ExampleParsedValueSingleQuote - Test single quoted strings are ok and single quotes can be enclosed in double quotes
func ExampleParsedValueSingleQuote() {
	inp := `infile.txt outfile.txt -keep 'red green white gray' -fix 44 -upcase -quote "'" -quiet -- gparm -right 88 -up 48.3 -down 22.5`
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : infile.txt
	// 1 : outfile.txt
	// 2 : -keep
	// 3 : red green white gray
	// 4 : -fix
	// 5 : 44
	// 6 : -upcase
	// 7 : -quote
	// 8 : '
	// 9 : -quiet
	// 10 : --
	// 11 : gparm
	// 12 : -right
	// 13 : 88
	// 14 : -up
	// 15 : 48.3
	// 16 : -down
	// 17 : 22.5
}

// TestParsedValueSingleQuote - Test map values with single quoted strings are ok and single quotes can be enclosed in double quotes
func TestParsedValueSingleQuote(t *testing.T) {
	var pass int
	t.Log("Test TestParsedValueSingleQuote() function started")

	/* perform tests */
	inp := `infile.txt outfile.txt -keep 'red green white gray' -fix 44 -upcase -quote "'" -quiet -- gparm -right 88 -up 48.3 -down 22.5`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedValueSingleQuote() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["quote"] == "'" {
			pass++
		} else {
			t.Error("Test TestParsedValueSingleQuote() ''' not in opmap['quote']")
		}
		if opmap["#1"] == "outfile.txt" {
			pass++
		} else {
			t.Error("Test TestParsedValueSingleQuote() 'outfile.txt' not in opmap['#1']")
		}
		if opmap["quiet"] == "" {
			pass++
		} else {
			t.Error("Test TestParsedValueSingleQuote() '' not in opmap['quiet']")
		}
		if opmap["right"] == "88" {
			pass++
		} else {
			t.Error("Test TestParsedValueSingleQuote() '88' not in opmap['right']")
		}
	}

	/* publish test success or failure */
	if pass != 4 {
		t.Error("Test TestParsedValueSingleQuote() failed.")
	} else {
		t.Log("Test TestParsedValueSingleQuote() OK.")
	}
}

// ExampleParsedValueBracketing - Test argumants are split on blanks and bracketing characters, dash int is considdered a flag name
func ExampleParsedValueBracketing() {
	inp := `1234 -22.6 5678,334 $58.99 -555 58% -333.22 (77) TT[99] map(66) {curl} possable values here`
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : 1234
	// 1 : -22.6
	// 2 : 5678,334
	// 3 : $58.99
	// 4 : -555
	// 5 : 58%
	// 6 : -333.22
	// 7 : (77)
	// 8 : TT[99]
	// 9 : map(66)
	// 10 : {curl}
	// 11 : possable
	// 12 : values
	// 13 : here
}

// TestParsedValueBracketing - Test map values with argumants are split on blanks and bracketing characters, dash int is considdered a flag name
func TestParsedValueBracketing(t *testing.T) {
	var pass int
	t.Log("Test TestParsedValueBracketing() function started")

	/* perform tests */
	inp := `1234 -22.6 5678,334 $58.99 -555 58% -333.22 (77) TT[99] map(66) {curl} possable values here`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedValueBracketing() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["#0"] == "1234" {
			pass++
		} else {
			t.Error("Test TestParsedValueBracketing() '1234' not in opmap['#0']")
		}
	}

	/* publish test success or failure */
	if pass != 1 {
		t.Error("Test TestParsedValueBracketing() failed.")
	} else {
		t.Log("Test TestParsedValueBracketing() OK.")
	}
}

// ExampleParsedValuePlusBang - Test plus + and bang ! are considdered flag name indicators
func ExampleParsedValuePlusBang() {
	inp := `456 'abc' +246 "jkl" 1.56e-28 -loa1d+runner "0.0.34" !name -f high=low -hamer=false --revenge=true -size 79 -w 92.3 -c red`
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : 456
	// 1 : abc
	// 2 : +246
	// 3 : jkl
	// 4 : 1.56e-28
	// 5 : -loa1d+runner
	// 6 : 0.0.34
	// 7 : !name
	// 8 : -f
	// 9 : high=low
	// 10 : -hamer=false
	// 11 : --revenge=true
	// 12 : -size
	// 13 : 79
	// 14 : -w
	// 15 : 92.3
	// 16 : -c
	// 17 : red
}

// TestParsedValuePlusBang - Test map values with plus + and bang ! are considdered flag name indicators
func TestParsedValuePlusBang(t *testing.T) {
	var pass int
	t.Log("Test TestParsedValuePlusBang() function started")

	/* perform tests */
	inp := `456 'abc' +246 "jkl" 1.56e-28 -load+runner "0.0.34" !name -f high=low -hamer=false --revenge=true -size 79 -w 92.3 -c red`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedValuePlusBang() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["#0"] == "456" {
			pass++
		} else {
			t.Error("Test TestParsedValuePlusBang() '456' not in opmap['#0']")
		}
		if opmap["#4"] == "1.56e-28" {
			pass++
		} else {
			t.Error("Test TestParsedValuePlusBang() '1.56e-28' not in opmap['#4']")
		}
		if opmap["hamer"] == "false" {
			pass++
		} else {
			t.Error("Test TestParsedValuePlusBang() 'false' not in opmap['hamer']")
		}
		if opmap["revenge"] == "true" {
			pass++
		} else {
			t.Error("Test TestParsedValuePlusBang() 'true' not in opmap['revenge']")
		}
		if opmap["246"] == "jkl" {
			pass++
		} else {
			t.Error("Test TestParsedValuePlusBang() 'jkl' not in opmap['246']")
		}
	}

	/* publish test success or failure */
	if pass != 5 {
		t.Error("Test TestParsedValuePlusBang() failed.")
	} else {
		t.Log("Test TestParsedValuePlusBang() OK.")
	}
}

// ExampleParsedNameNumber - Test flag names can contain numbers
func ExampleParsedNameNumber() {
	inp := `-name456 88 -char99 "ninety nine" -99 'bottles'`
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : -name456
	// 1 : 88
	// 2 : -char99
	// 3 : ninety nine
	// 4 : -99
	// 5 : bottles
}

// TestParsedNameNumber - Test map values with flag names can contain numbers
func TestParsedNameNumber(t *testing.T) {
	var pass int
	t.Log("Test TestParsedNameNumber() function started")

	/* perform tests */
	inp := `-name456 88 -char99 "ninety nine" -99 'bottles'`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedNameNumber() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["name456"] == "88" {
			pass++
		} else {
			t.Error("Test TestParsedNameNumber() '1234' not in opmap['#1']")
		}
	}

	/* publish test success or failure */
	if pass != 1 {
		t.Error("Test TestParsedNameNumber() failed.")
	} else {
		t.Log("Test TestParsedNameNumber() OK.")
	}
}

// ExampleParsedBadName - Test bad flag names have some special characters
func ExampleParsedBadName() {
	inp := `-bad.name 49 -my-name "bill smith" -no$man 0 -what/k `
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : -bad.name
	// 1 : 49
	// 2 : -my-name
	// 3 : bill smith
	// 4 : -no$man
	// 5 : 0
	// 6 : -what/k
}

// TestParsedBadName - Test map values with bad flag names have some special characters
func TestParsedBadName(t *testing.T) {
	var pass int
	t.Log("Test TestParsedBadName() function started")

	/* perform tests */
	inp := `-bad.name 49 -my-name "bill smith" -no$man 0 -what/k `
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedBadName() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["bad.name"] == "49" {
			pass++
		} else {
			t.Error("Test TestParsedBadName() '.name49' not in opmap['bad']")
		}
		if opmap["my-name"] == "bill smith" {
			pass++
		} else {
			t.Error("Test TestParsedBadName() 'bill smith' not in opmap['my-name']")
		}
		if opmap["what/k"] == "" {
			pass++
		} else {
			t.Error("Test TestParsedBadName() '/k' not in opmap['what']")
		}
	}

	/* publish test success or failure */
	if pass != 3 {
		t.Error("Test TestParsedBadName() failed.")
	} else {
		t.Log("Test TestParsedBadName() OK.")
	}
}

// ExampleParsedArgStringSeq - Test string arg sequences
func ExampleParsedArgStringSeq() {
	inp := `"Mary had a little" "it\"s flees were" 'Wise men came from a fire' 'Harry\'s name was mud' "for the MCP"`
	args := new(Args)
	err := args.ParseArg(inp)
	if err == nil {
		args.PrtOpts()
	}
	// Output:
	// Parsed Option Flags:
	// 0 : Mary had a little
	// 1 : it\"s flees were
	// 2 : Wise men came from a fire
	// 3 : Harry\'s name was mud
	// 4 : for the MCP
}

// TestParsedArgStringSeq - Test map values with string arg sequences
func TestParsedArgStringSeq(t *testing.T) {
	var pass int
	t.Log("Test TestParsedArgStringSeq() function started")

	/* perform tests */
	inp := `"Mary had a little" "it\"s flees were" 'Wise men came from a fire' 'Harry\'s name was mud' "for the MCP"`
	args := new(Args)
	err := args.ParseArg(inp)
	if err != nil {
		t.Error("Test TestParsedArgStringSeq() error: " + err.Error())
	} else {
		opmap := args.GetMap()
		if opmap["#1"] == `it\"s flees were` {
			pass++
		} else {
			t.Error("Test TestParsedArgStringSeq() 'it\"s flees were' not in opmap['#1']")
		}
		if opmap["#4"] == "for the MCP" {
			pass++
		} else {
			t.Error("Test TestParsedArgStringSeq() 'for the MCP' not in opmap['#4']")
		}
	}

	/* publish test success or failure */
	if pass != 2 {
		t.Error("Test TestParsedArgStringSeq() failed.")
	} else {
		t.Log("Test TestParsedArgStringSeq() OK.")
	}
}
