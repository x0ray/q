// Package qs - q scripting language
package qs

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Following is the help over view - supplied first with -verbose option
const HelpOverViewText = "" +
	`The ` + PGM + ` program runs a Q script program. The form of the Q scripting 
language is described in detail if the -verbose option is used along with -help
option. 

Notes:
   1. The Q script language includes compatible features and syntax 
	  originally used by many other languages, including PL1, Go, 
	  Lua, C, C++, Fortran. If you recognize features from these
      languages it is likely they will work as you remember, however, there
      are differences. Check for differences, by testing small examples, and 
      referring to the 'help -verbose' text for correct Q script syntax.  
   2. Due to the experimental nature of the ` + PGM + ` ` + VER + ` scripting 
      program, its specification subject to change, and backwards compatibility
      may not be provided.  
`

// Following is the main command help blurb - always supplied
const HelpText = `Usage: ` + PGM + `
   [ command ]
   [ -pgm ]     <program>        File name containing Q program to run.
   [ -exec ]    <script-segment> String of Q language statements.
   [ -lib ]     <q-lib-file>     Q library file name to access.
   [ -profile ] <profile-file>   File name to use for profile data.
   [ -log ]     <log-file>       File name to use for logging.
   [ -name ]    <name-string>    Name tag used in logging, and _NAME script variable.
   [ -limit ]   <nnn>            Sets a memory size limit for Q program.
   [ -inter ]                    Use interactive mode.
   [ -help | -h ]                Display this help then exit.
   [ -quiet ]                    Do not show output.
   [ -version | -v ]             Show version information then exit.
   [ -verbose ]                  Show more information for functions with output.
   [ -debug ]                    Show debugging info	

Commands:
   help:     Display help information and quit.
   int:      Run ` + PGM + ` in interactive mode.
   run:      Run the ` + PGM + ` script named on the -pgm option.
   version:  Display program version information and quit.
	
Use '` + PGM + ` help -verbose' for help additional ` + PGM + ` help topics:	
   * Language syntax
   * Language built in procs
   * Script examples
   * Environment variables
   * Return Codes	
`

// TODO expand help details
// Following is the main command help blurb
const HelpTextVerbose = languageSyntax + builtInProcs + scriptExamples + environmentVariables + returnCodes

const languageSyntax = `
Language syntax:
` + comments + reservedWords + operators + variables + controlStructures

const reservedWords = `
  Reserved words: 
    and break dcl do else elseif end false for if in nil 
	not or proc repeat return then true until while
	
 `

const comments = `
  Comments: 
    Three forms of comment exist, line and block:
	- The characters // start a line comment and the end of line marks the comment end.
	- The characters #! start an alternate line comment. The end of line marks 
      the of the comment. This form can be used on some systems, when used to name 
      the location of ` + PGM + `. For example #!/usr/local/bin/` + PGM + ` 
	- The characters /* start a block comment and */ ends a block comment.
	
  Comment examples:
    // line comment
	a = 55  // end of line comment
	
    /* a block comment can be part of a line, 
      even between statement tokens
      or span multiple lines */
    a = /* 123 */ 42
		
 `

const operators = `
  Operators: 
    Symbol  Function          Type     Association  Precedence
    ----------------------------------------------------------
    +       addition          binary   left         5
    -       subtraction       binary   left         5
    *       multiplication    binary   left         6
    /       division          binary   left         6
    %       modulus           binary   left         6
    and     logical and       binary   left         2
    or      logical or        binary   left         1
    <       less              binary   left         3 
    >       greater           binary   left         3       
    <=      less or equal     binary   left         3
    >=      greater or equal  binary   left         3 
    !=      not equal         binary   left         3  
    ==      equal             binary   left         3
    ||      concatenate       binary   right        4
    #       length            unary    left         7
    not     logical not       unary    right        7   
    -       unary negative    unary    right        7 
    ^       exponent          binary   right        8

    A higher precedence is interpreted ahead of a lower precedence.
	
 `

const variables = `
  Variables: 
    - Variable names can use the following characters: A..Z, a..z, 0..9, and _
    - A variable can be created by initial assignment which will create a globally
      accessible variable seen by all procs, or a variable can be created using
      the dcl statement which will make the variable accessible only by the proc 
      that executed the dcl statement.
	- Variable types can be:
        string      s = "Hello Zaphod"
        number      n = 42
        struct      a = {1,2,3}  or  b = {e=33,t="time",c=os.clock()}
        proc        f = os.getenv("PATH")

  Builtin global variables
	_VERSION  - vesion of ` + PGM + `
	_NAME     - value presented in -name option or ` + NAME + ` 
	os        - name of operating system, linux, windows
	arch      - architecture name of executing hardware. One of 386, amd64, arm, ...
	pi        - value of mathematical constant pi (limited by system architecture)
	e         - value of mathematical constant e (limited by system architecture)
	phi       - value of mathematical constant phi (limited by system architecture)
	sqrt2     - value of the square root of 2 (limited by system architecture)
	sqrte     - value of the square root of e (limited by system architecture)
	huge      - the largest floating point value 
	small     - the smallest floating point value
	arg       - a list containing the scripts arguments 
 `

const controlStructures = `
  Control Structures: 
    The following logical control directives are available:
	 
    - while <expression> do <block> end

    - repeat <block> until <expression> 

    - if <expression> then <block> end
    - if <expression> then <block> else <block> end
    - if <expression> then <block> elseif <expression> then <block> else <block> end

    - for name = <start-expression> , <end-expression> [, <inc-expression>] do <block> end 

    - break
    - return [values]
	
 `

const builtInProcs = `
Language built in procs
  Q script has many procs built in. There are standard procs and auxillary procs. 
  The standard procs do not need a name space prefix in order to be referenced. 
  For example using proc bye to exit the script is specified: 
    bye() 
  However the auxillary IO proc 'open' is specified with a prefix 'i.':
	f = i.open("myfile.txt","r")

  Standard procs 
	[z:bool =] assert(a:bool,b:str)
		Asserts 'a' is true and if not issues message 'b' Returns 'a'.
		
	bye([a:num]) 
		Stops the executing script, returns code 'a' to system.
		
	collectgarbage( [a:str("stop"|"restart"|"collect"|"step"|"count"]) )
		Runs specified mode of garbage collection to free used resources.
	
	error(a:str)               
		Issues error message 'a'.
	
	z = getfenv(a:proc) 
		Returns the current environment into 'z' for the proc 'a'.
		
	z = getmetalist(a:list)
		Returns meta info into 'z' for list 'a'.  
	
	help() 
		Displays this help text.
		
	z = load(a:proc [,b:str])
		Loads proc 'z', using proc"nil", "bool", "num", "str", "proc", "data", 
		"thread", "list", "chan" 'a' to load multiple segment strings containing 
		an Q script.
		
	z = loadfile(a:str)
		Loads proc 'z' from a file named 'a' which contains an Q script.
	
	z = loadstring(a:str)
		Loads proc 'z' from a string 'a'.
	
	log(a:str)
		Writes message string 'a' on the log as an info message.
	
	logd(a:str)
		Writes message string 'a' on the log as a debug message.
	
	loge(str)
		Writes message string 'a' on the log as an error message.
		
	logi(str)
		Writes message string 'a' on the log as an info message.
		
	logw(str)
		Writes message string 'a' on the log as a warning message
		
	i,v = next(a:list,b:*)
		Get the next index 'i' and value 'v' from the list 'a'.
	
	ok,r = pcall(a:proc,b:*,c:*,...)
		Call proc 'a' with args 'b', 'c', etc. Catch any runtime errors. If no 
		error 'ok' is true. 
	
	put(a:str,b:str,...)
		Write strings 'a', 'b', etc to stdout.
	
	quit([a:num]) 
		Stops the executing script, returns code 'a' to system.
		
	z:bool = rawequal(a:*, b:*)
		Checks if 'a' equals 'b' and returns true in 'z' if they are equal.
	
	z = rawget(a:list, i:*)
		Gets the real value of 'a[i]' into 'z'.
		
	rawset(a:list, i:*, v:*)
		Sets the value 'v' into list 'a' at index 'i'.
	
	run(a:str)
		Run the Q srcipt from file name 'a'.
` + // TODO Fix help
	/*
		select()
		_printregs()
		setfenv()
		setmetalist()
	*/`
	
	stop([a:num]) 
		Stops the executing script, returns code 'a' to system.
		
	n = tonumber(s:str)
		Returns string 's' converted to a number.
	
	s = tostring(n:num)
		Returns number 'n' converted to a string.
	
	t = type(a:*)
		Returns the type of 'a', one of: "nil", "bool", "num", "str", "proc", 
		"data", "thread", "list", "chan".	
	
	a[i], ..., a[j] = unpack(a:list[,i:num[,j:num]])
		Returns a slice if the list 'a' from 'i' to 'j'
	
	ok, r = xpcall(a:proc, e:proc)
		Call proc 'a' catching any runtime errors. If an error occurs 'ok' is 
		set to false and proc 'e' is called.  
` + // TODO Fix help
	/*
				module()

				require()
		  Dbg Lib     (library name not required to call)
			getfenv()
			getinfo()
			getlocal()
			getmetalist()
			getupvalue()
			setfenv()
			setlocal()
			setmetalist()
			setupvalue()
			traceback()
	*/`

  Input and Output procs:     
	i.close(f)
		Close file handle 'f'.
	
	i.flush(f)
		Save any buffered data written to file handle 'f' to disk.
	
	s:str = i.lines(f)
		Read next line 's' from file handle 'h'. 
	
	f = i.input([f|n:str])
		Returns default input file, or assigns new default fron name 'n' or 
		handle 'f'.
	
	f = i.output([f|n:str])
		Returns default output file, or assigns new default fron name 'n' or 
		handle 'f'.
	
	f = i.open(n:str[,m:str])
		Opens file name 'n' using mode 'm', returns file handle 'f'. Mode can be
		"r", "w", "a", "r+", "w+", "a+", indicating read, write, append, and "+"
		indicates update mode. Modes may have "b" concatenated for binary. 
		
` + // TODO Not included - os dependant
	/*
		io.popen()
	*/`	
	f:read()
	
	f:type()
	
	f:tmpfile()
		Returns handle 'f' to a temporary file, opened in update mode. The file is
		deleted on Q script termination. 
	
	f:write()
	
	f:close()
	
  Mathematics procs:    	
	z:num = abs(a:num)
		Returns the absolute value of 'a' in 'z'.
	
	z:num = acos(a:num)
		Returns the arc cosine value of 'a' in 'z'.
		
	z:num = asin(a:num)
		Returns the arc sine value of 'a' in 'z'.
		
	z:num = atan(a:num)
		Returns the arc tangent value of 'a' in 'z'.
` + // TODO Not included - os dependant
	/*
		z:num = atan2(a:num,b:num)
	*/`	
	z:num = ceil(a:num)
		Returns the smallest integer larger than or equal to 'a' in 'z'.
		
	z:num = cos(a:num)
		Returns the cosine value of 'a' in 'z'.
		
	z:num = cosh(a:num)
		Returns the hyperbolic cosine value of 'a' in 'z'.
		
	z:num = deg(a:num)
		Returns the angle of 'a' specified in radians as 'z' in degrees.
		
	z:num = exp(a:num)
		Returns the value of e raised to the power of 'a' in 'z'.
		
	z:num = fact(a:num)
		Returns the factorial of 'a' in 'z'.
		
	z:num = fib(a:num)
		Returns the Fibonacci number for 'a' in 'z'.
		
	z:num = floor(a:num)
		Returns the largest integer smaller than or equal to 'a' in 'z'.
		
	z:num = fmod(a:num,b:num)
		Returns the remainder of 'a' divided by 'b' in 'z'.
	
	z:num,y:num = frexp(a:num)
		Returns 'z' and 'y' such that 'a' = z2**e, where e is an integer and the
		absolute value of 'z' is in the range [0.5, 1] or 0 when 'a' is 0.
	
	z:num = ldexp(a:num,b:num)
		Returns 'a'*2**'b' in 'z'.	
	
	z:num = log(a:num)
		Returns the natural log of 'a' in 'z'.
	
	z:num = log10(a:num)
		Returns the log base 10 of 'a' in 'z'.
		
	z:num = max(a:num, b:num [,...])
		Returns the maximum of 'a', 'b' etc in 'z'.
		
	z:num = mean(a:num, b:num [,...])
		Returns the mean of 'a', 'b' etc in 'z'.
		
	z:num = median(a:num, b:num [,...])
		Returns the median of 'a', 'b' etc in 'z'.
		
	z:num = min(a:num, b:num [,...])
		Returns the minimum of 'a', 'b' etc in 'z'.
		
	z:num,y:num = mod(a:num,b:num)
		Returns the remainder of 'a' divided by 'b' in 'z'.	
	
	z:num = mode(a:num, b:num [,...])
		Returns the mode of 'a', 'b' etc in 'z'.
		
	z:num,y:num = modf(a:num)
		Returns the integer part of 'a' in 'z' and the fractional part in 'y'.
		
	z:num = pow(a:num,b:num)
		Returns 'a'**'b' in 'z'.	
	
	z:num = rad(a:num)
		Returns the angle 'a' specified in degrees in 'z' in radians.	
	
	z:num = random([a:num[,b:num]])
		Returns a random number in 'z', in range 0..1, 1..'a', or 'a'..'b' 
		depending on wether 'a' and 'b' are supplied. 'a' must be greater than 'b'.	
	
	randomseed(a:num)
		Sets the internal random seed. 'a' must be a number.
	
	z:num = range(a:num, b:num [,...])
		Returns the range of 'a', 'b' etc in 'z'.
		
	z:num = rms(a:num, b:num [,...])
		Returns the root mean square of 'a', 'b' etc in 'z'.
		
	z:num = sin(a:num)
		Returns the sine value of 'a' in 'z'.
		
	z:num = sinh(a:num)
		Returns the hyperbolic sine value of 'a' in 'z'.
		
	z:num = sqrt(a:num)
		Returns the square root of 'a' in 'z'.
		
	z:num = stddev(a:num, b:num [,...])
		Returns the standard deviation of 'a', 'b' etc in 'z'.
		
	z:num = sum(a:num, b:num [,...])
		Returns the sum of 'a', 'b' etc in 'z'.
		
	z:num = tan(a:num)
		Returns the tangent value of 'a' in 'z'.
	
	z:num = tanh(a:num)
		Returns the hyperbolic tangent value of 'a' in 'z'.
  
	z:str = uuidgen()
		Returns a generated UUID in 'z'.
		
	z:str = uuidgenfmt()
		Returns a generated and formatted UUID in 'z'.
		
	z:num = variance(a:num, b:num [,...])
		Returns the variance of 'a', 'b' etc in 'z'.
		
  Operating System procs:  
	z:str = argstr()
		Returns the arguments passed to the script as a string in 'z'.
	     
	z:str = arglist(a:str)
		Returns the arguments from string 'a' parsed as a list in 'z'.
	     
	z:str = argopts(a:str)
		Returns the arguments from string 'a' parsed as a list of named options and arguments in 'z'.
	     
	z:bool = chdir(a:str) 
		Changes the current working directory to 'a'. Returns true in 'z' if the 
		directory change was a success.	
	
	o.clearenv()
		Deletes all environment variables.
	
	z:num = clock()
		Returns the time in seconds since the script started executing
	
	difftime(a:num,b:num)
		Returns the difference between time 'a' and time 'b' in 'z'.
	
	z:num = execute(a:str)
		Executes command 'a' as an OS process. Returns 0 if no errors occur or 1 if 
		there were errors.		
	
	z:bool = exist(f:str)
		Checks the existance of file 'f', returns true if file 'f' exists, otherwise 
		returns false.
	
	exit([a:num])
		Stops executing the Q script, and exits to the OS, optionally returning 
		the code in 'a' to the OS shell. 
	
	z = date([a:str[,b:num]])
		Returns a date in a list, or as a number depending om the format string 'a'.
		If the format 'a' is not supplied the date and time are in local time as a 
		number of seconds from the epoch. If the format string 'a' begins with "!"
		the datetime is UTC. If the format string 'a' begins with "*t" the date time
		is returned in a list containing elements: "year","month","day","hour",
		"min","sec","weekday","yearday","isdst".	
	
	z:str = getenv(a:str)
		Returns the value of the environment variable named 'a' in 'z'. 
	
	z:num = geteuid()
		Returns the effective user ID as a number in 'z'.
		
	z:str = gethome()
		Returns the users home path string in 'z'. This is system dependant. On 
		some systems this may be the users name.
		
	z:str = getuser()
		Returns the user ID string in 'z'.
			
	z:num = getpid()
		Returns the current process ID as a number in 'z'.
			
	z:num = getppid()
		Returns the current processes parent process ID as a number in 'z'.
	
	z:num = getuid()
		Returns the user ID as a number in 'z'.	
	
	z:str = getwd()
		Returns the current working directory in 'z'.
	
	z:str = hostname()
		Returns the host name in 'z'.	
	
	z:bool = remove(a:str)
		Removes the file name 'a'. If successful true is returned in 'z'.		
	
	z:bool = rename(a:str,b:str)
		Renames file name 'a' to file name 'b'. If successful true is returned in 'z'.		
		
	z:bool = setenv(a:str,b:str)
		Sets environment variable name 'a' to value 'b'. If successful true is returned in 'z'.			
	
	sleep(a:num)
		The script stops executing for 'a' seconds.			
	
	z:list = stat(f:str)
		Returns a list in 'z' containing file information for file 'f'.			
	
	z:list = statfs(p:str)
		Returns a list in 'z' containing file system information for path 'p'.			
	
` + // TODO Not implemented
	/*
		o.setlocale()
	*/`		
	z = time([a:list])
		Returns the time nin 'z'. If list 'a' is specified 'z' is returnes a list, else
		'z' is a Unix time value as a number.
	
	z:str = tmpname()
		Returns a temporary file name from the current working directory in 'z'.
	
	z:bool = unsetenv(a:str)
		Un-sets the environment variable name 'a' and returns true in 'z' if successful. 
	

  String procs:
	z:str = after(a:str,b:str)
		Returns a sub-string of string 'a' in 'z' containing all characters after 
		substring 'b'. If 'b' is not located 'z' is the empty string. 
	
	z:str = before(a:str,b:str)
		Returns a sub-string of string 'a' in 'z' containing all characters before 
		substring 'b'. If 'b' is not located 'z' is the empty string. 
	
	z:num = byte(a:str[,b:int[,c:int]])
		Returns a code point for 'a[b]'.
	
	z:str = char(a:num)
		Returns character in 'z' for code point 'a'.
	
	z:bool = contains(a:str,b:str)
		Returns true in 'z' if string 'a' contains string 'b'.
	
	z:bool = containsany(a:str,b:str)
		Returns true in 'z' if string 'a' contains any characters from string 'b'.
	
	z:num = count(a:str,b:str)
		Returns the count in 'z' of strings 'b' in string 'a'.
		
	z:str = decodebase64(a:str)
		Decodes a base 64 encoded string in 'a' to a text string in 'z'.	
	
	z:str = dump(a:str)
		Creates a formatted dump of 'a' in string 'z'.
		
	z:str = encodebase64(a:str)
		Creates a base 64 encoded form of 'a' in string 'z'.		
	
	s:num,e:num = find(a:str,p:str[,init[,plain]])
		Returns the start 's' and end 'e' position of pattern 'p' in string 'a'.
	
	z:str = format(f:str,s:str)
		Formats string 's' depending on format string 'f' returning formatted string
		in 'z'.
	
	gsub()
	
	z:bool = hasprefix(a:str,b:str)
		Returns true in 'z' if string 'a' has string 'b' as a prefix.	
	
	z:bool = hassuffix(a:str,b:str)
		Returns true in 'z' if string 'a' has string 'b' as a prefix.	
	
	z:num = index(a:str,b:str)
		Returns the first position in 'z' of string 'b' in string 'a'. If 'b' is 
		not in 'a' then -1 is returned.
	
	z:num = indexany(a:str,b:str)
		Returns the first position as 'z' of any character in string 'b' that is in 
		string 'a'.	
	
	z:num = lastindex(a:str,b:str)
		Returns the last position in 'z' of string 'b' in string 'a'. If 'b' is 
		not in 'a' then -1 is returned.		
	
	z:num = lastindexany(a:str,b:str)
		Returns the last position as 'z' of any character in string 'b' that is in 
		string 'a'.	
	
	z:int = len(a:str)
		Returns the length of string 'a' in 'z'.
	
	z:int = length(a:str)
		Returns the length of string 'a' in 'z'. This is an alias for len(). See 
		also the length operator # above.
	
	z:str = lower(a:str)
		Returns in 'z' the string 'a' converted to lower case.
	
	z:str = match(a:str,p:str,i:num)
		Returns capture from 's' using pattern 'p' into 'z' optionally starting at 'i'.
	
	z:list = prxmatch(rx:str,a:str)
		Returns a list of positions in 'z' of the regular expression 'rx' match 
		results found while matching string 'a'.
	
	z:str = prxchange(rx:str,a:str,b:str)
		Replaces regular expression 'rx' matches found in 'a' with value from 'b' 
		returning the results in 'z'.
	
	z:str = rep(a:str,b:num)
		Returns 'b' concatenated copies of string 'a' in string 'z'.
	
	z:str = replace(a:str,b:str,c:str,d:num)
		Replaces 'd' strings of 'b' found in 'a' with strings of 'c' and returns
		the resuly in 'z'. If 'd' is -1 all strings of 'b' are replaced. 
	
	z:str = reverse(a:str)
		Returnes in 'z' all characters of string 'a' in reverse order.
	
	z:str = scan(a:str,b:str,c:int)
		Returns a sub-string of string 'a' in 'z' which is the 'c'th substring 
		from the left, delimited by one or more of the characers in string 'b'. 
		If 'c' is a negative value the scan proceeds from right to left, and 
		'c' is counted from the right. 
	
	z:list = scanall(a:str,b:str)
		Returns a list of all sub-strings of string 'a' in 'z' which are delimited 
		by one or more of the characers in string 'b'. 
	
	z:str = sub(a:str,b:int[,c:int])
		Returns a sub-string of string 'a' in 'z' from position 'b' inclusive, 
		to position 'c'.
	
	z:str = substr(a:str,b:int[,c:int])
		Returns a sub-string of string 'a' in 'z' from position 'b' inclusive,
		for length 'c'.
	
	z:str = trim(a:str,b:str)
		Returns in 'z' the string 'a' with all characters from 'b' removed from
		the left and right side .
		
	z:str = trimleft(a:str,b:str)
		Returns in 'z' the string 'a' with all characters from 'b' removed from 
		the left side removed.
		
	z:str = trimprefix(a:str,b:str)
		Returns in 'z' the string 'a' with prefix string 'b' removed from the 
		left side.
	
	z:str = trimright(a:str,b:str)
		Returns in 'z' the string 'a' with all characters from 'b' removed from
		the right side removed.
		
	z:str = trimspace(a:str)
		Returns in 'z' the string 'a' with all white space characters on the 
		left and right side removed.
	
	z:str = trimsuffix(a:str,b:str)
		Returns in 'z' the string 'a' with suffix string 'b' removed from the 
		right side.
	
	z:str = title(a:str)
		Returns in 'z' the string 'a' converted to title case.
	
	z:str = upper(a:str)
		Returns in 'z' the string 'a' converted to upper case.
  

  EMI procs:
	z:bool = consulavailable()
		Check if Consul is available
	
	z:str = consulcheckerror()
		Get the Consul error message
	
	z:bool = consulcheckkey(k:str)
		Check if Consul key exists
	
	z:str = consuldeletekeys(kb:str)
		Delete Consul keys
	
	z:str = consulgetkey(k:str)
		Get Consul key value
	
	z:list = consulgetkeys(kb:str)
		Get list of Consul key values
	
	z:str = consulputkey(k:str,v:str)
		Set a Consul key value
	
	z:str = consulsetkeys(kb:str,ks:list)
		Set list of Consul key values
		
	z:list = deploymentinfo()
		Create a list containing deployment values		
		
	z:str = setmailhost(user:str,pwd:str,c:host)
		Set mailer credentials
	
	z:str = sendmail(from:str,to:str,c:msg)
		Send an email message	
				
	
  QList procs:
	l.getn()
	l.concat()
	l.insert()
	l.maxn()
	l.remove()
	l.sort()

`

const scriptExamples = `
Script examples:

Example 1:
	/*
	  Script:   test.q
	  Language: q -- scripting control language.	
	  Note: This is a block style comment.
	*/
	
	// This is a line comment  
	
	PGM = "test.q" ;   // PGM is a string variable
	VER = "0.0.1" ;     // the semi-colon ';' is optional
	// The builtin put() proc writes one or more comma 
	// seperated values on the stdout file.
	// The || operator concatenates strings.
	log("Program:" || PGM || " version:" || VER) ;
	
	// get clock time using os library proc
	st = clock()   // time stamp in seconds for later use
	
	// Create a proc (sub-routine / function) 
	// Procedures must have an associated end statement
	// Types are not required for variables
	proc blabla(a,b)
	  dcl c        // define local variable called c
	  c = a + b    // add parameters a and b giving c
	  return c     // return c to the caller
	end
	
	// use the new proc blabla
	pten = blabla(567,9)
	put("pten: ",pten) 
	put("blabla: ",blabla(99,86)) 
	put("The","answer",blabla(38,4))
	
	// make a local alias called 'pid' of os library proc
	dcl pid = getpid() ;  // declare pid a local variable
	put("current pid:" || pid) ;
	
	// various built in string library procs
	ss = "this is a story about a man named Jed whose kin folks were bankers!"
	put(ss)
	put("String has prefix 'thi'? ",hasprefix(ss,"thi")) 
	put("String has prefix 'quack'? ",hasprefix(ss,"quack")) 
	put("String has suffix 'ers!'? ",hassuffix(ss,"ers!")) 
	put("String has suffix 'dog'? ",hassuffix(ss,"dog")) 
	put("String contains 'named'? ",contains(ss,"named")) 
	put("String contains 'Borris'? ",contains(ss,"Borris")) 
	put("Count of 'is' is:",count(ss,"is"))
	put("Count of 'a' is:",count(ss,"a"))
	ss = title(ss)
	put(ss)  
	rr = reverse(ss)     // reverse
	put(rr)
	ss = rep(rr,3)       // repeat  
	rr = replace(ss,"deJ","Harry Potter",2)
	du = dump(rr)        // dump
	put(du)
	put("String rr:",rr)
	put("String containsany 'x y z'? ",containsany(rr,"x y z")) 
	put("String containsany 'q'? ",containsany(rr,"q")) 
	put("Index of 'Harry'? ",index(rr,"Harry")) 
	put("Index of 'Barry'? ",index(rr,"Barry")) 
	put("Index any for 'aeiou' ",indexany(rr,"aeiou")) 
	put("Index any for 'xyz' ",indexany(rr,"xyz")) 
	put("Last index any for 'aeiou' ",lastindexany(rr,"aeiou")) 
	put("Last index any for 'xyz' ",lastindexany(rr,"xyz")) 
	
	put("String rr:",rr)
	tr = trim(rr,"aeiou")
	put("Trim:",tr)
	tr = trimleft(rr," !sreknaB")
	put("Trimleft:",tr)
	tr = trimprefix(rr," !sreknaB")
	put("Trimprefix:",tr)
	tr = trimright(rr," IsihT")
	put("Trimpright:",tr)
	
	rr = "    this is a test    "
	put("String rr:",rr)
	tr = trimspace(rr)
	put("Trimspace:",tr)
	
	rr = "this is a test"
	put("String rr:",rr)
	tr = trimsuffix(rr,"test")
	put("Trimpsuffix:",tr)
	
	// various built in math library procs
	a = 44 ;
	b = 6 ;
	c = a - b ;
	put("c:",c) ;
	
	// a variable can be a proc
	add2 = blabla
	f = add2(a,b)
	put("a add2 b:",f)
	
	// calculate and print elapsed time in seconds as float
	et = clock()
	put("Elapsed time:",et-st,"(s)")
	
	put("Program:" || PGM || " ended.") ;
	
Example 2:
	put("File test Q program")
	f = i.open("test.txt","r")
	if type(f) == "data" then
	  for line in f:lines() do
	    t = sub(line,26)
	    if t != "" then
	      put("::>",t)
	    end
	  end
	end
	f:close()
	put("File test program ended.")


`

const environmentVariables = `
Environment variables used by ` + PGM + `: none

`

const returnCodes = `
Return Codes:
  ok       = 0    
  warning  = 1   
  critical = 2    
  fatal    = 3  
`

// Following is the extended command help blurb - always supplied
const XHelpText = `Extended usage: ` + PGM + ` 
   [-Xhelp]                    Display extended help options.         
   [-Xsyntax]                  Debug tool, displays syntax tree information.
   [-Xicode]                   Debug tool, shows internal instruction code list.
				 			 	
`

// Following is the extended command help blurb
const XHelpTextVerbose = `
` + PGM + ` Extended Details

`

var VerboseHelp bool

// prtHelpFlag is used to format the programs parameters for flag.VisitAll()
// when it processes the -help or -h parameter
func PrtHelpFlag(f *flag.Flag) {
	if !strings.HasPrefix(f.Name, "X") {
		// get the flag type to display
		typstr := fmt.Sprintf("%v\n", reflect.TypeOf(f.Value))
		// remove the excess junk from the type
		typstr = strings.Split(typstr, ".")[1] // ie:  *flag.boolValue --> boolValue
		typstr = typstr[:len(typstr)-6]        // ie:  boolValue --> bool
		defVal, _ := strconv.Atoi(f.DefValue)
		if defVal == NOTEST {
			if VerboseHelp {
				fmt.Printf("  -%s %s (set %v) No default value \n\t%s\n", f.Name, typstr, f.Value, f.Usage)
			} else {
				fmt.Printf("  -%s %s No default value \n\t%s\n", f.Name, typstr, f.Usage)
			}
		} else {
			if VerboseHelp {
				fmt.Printf("  -%s %s (set %v) (default %v) \n\t%s\n", f.Name, typstr, f.Value, f.DefValue, f.Usage)
			} else {
				fmt.Printf("  -%s %s (default %v) \n\t%s\n", f.Name, typstr, f.DefValue, f.Usage)
			}
		}
	}
}

// prtXHelpFlag is used to format the programs parameters for flag.VisitAll()
// when it processes the -help or -h parameter
func PrtXHelpFlag(f *flag.Flag) {
	if strings.HasPrefix(f.Name, "X") {
		// get the flag type to display
		typstr := fmt.Sprintf("%v\n", reflect.TypeOf(f.Value))
		// remove the excess junk from the type
		typstr = strings.Split(typstr, ".")[1] // ie:  *flag.boolValue --> boolValue
		typstr = typstr[:len(typstr)-6]        // ie:  boolValue --> bool
		defVal, _ := strconv.Atoi(f.DefValue)
		if defVal == NOTEST {
			if VerboseHelp {
				fmt.Printf("  -%s %s (set %v) No default value \n\t%s\n", f.Name, typstr, f.Value, f.Usage)
			} else {
				fmt.Printf("  -%s %s No default value \n\t%s\n", f.Name, typstr, f.Usage)
			}
		} else {
			if VerboseHelp {
				fmt.Printf("  -%s %s (set %v) (default %v) \n\t%s\n", f.Name, typstr, f.Value, f.DefValue, f.Usage)
			} else {
				fmt.Printf("  -%s %s (default %v) \n\t%s\n", f.Name, typstr, f.DefValue, f.Usage)
			}
		}
	}
}
