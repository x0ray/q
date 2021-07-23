# Q scripting language

## Introduction

The Q command runs a Q script program. These scripts can also be 
run within the ubiquity command (see: ubiquity "commandType": 
"qsl"). The format of the Q scripting language is described in detail here.

The Q script language includes compatible features and syntax originally 
used by many other languages, including PL1, Go, Lua, C, C++, 
Fortran. If you recognize features from these languages it is likely they 
will work as you remember, however, there are differences. Check for 
differences, by testing small examples, and referring to the 'help -verbose' 
text for correct Q script syntax.

## Contents

* [Introduction](#introduction)
* [Usage summary](#usage-summary)
    * [Command and Option Details](#command-and-option-details)
        * [Command Format](#command-format)
    * [Option Details](#option-details)
        * [-debug](#debug) 	
        * [-exec](#exec) 		
        * [-h](#h) 	
        * [-help](#help) 	
        * [-inter](#inter)  
        * [-lib](#lib)  		
        * [-limit](#limit) 
        * [-log](#log) 	
        * [-name](#name) 		
        * [-pgm](#pgm)  	
        * [-profile](#profile) 	
        * [-quiet](#quiet) 	
        * [-v](#v) 	
        * [-verbose](#verbose)  	
        * [-version](#version)  
* [Q Language syntax](#language-syntax)  
    * [Comments](#comments)  
        * [Comment examples](#comment-examples)  
    * [Reserved words](#reserved-words)  
    * [Operators](#operators)  
    * [Variables](#variables)  
    * [Control Structures](#control-structures)  
    * [Built In Procedures and Functions](#q-Language-procedures-and-functions)  
        * [Standard](#standard-procs)  
           [assert](#assert) [bye](#bye) [collectgarbage](#collectgarbage) [error](#error) [getfenv](#getfenv) [getmetalist](#getmetalist) [help](#help) [load](#load)
           [loadfile](#loadfile) [loadstring](#loadstring) [log](#log) [logd](#logd) [loge](#loge) [logi](#logi) [logw](#logw) [next](#next)  [pcall](#pcall) [put](#put) [quit](#quit)
           [rawequal](#rawequal) [rawget](#rawget) [rawset](#rawset) [run](#run) [stop](#stop) [tonumber](#tonumber) [tostring](#tostring) [type](#type) [xpcall](#xpcall) 

        * [Input and Output](#input-and-output-procs)  
           [close](#close) [flush](#flush) [lines](#lines) [input](#input) [output](#output) [open](#open) [read](#read) [type](#type) [tmpfile](#tmpfile) [write](#write)

        * [Mathematics](#mathematics-procs-and-constants)  

            * [Constants](#math-constants)  
               [pi](#pi) [e](#e) [phi](#phi) [sqrt2](#sqrt2) [sqrte](#sqrte) [huge](#huge) [small](#small) 

            * [Math Procs](#math-procs)  
               [abs](#abs) [acos](#acos) [asin](#asin) [atan](#atan) [ceil](#ceil) [cos](#cos) [cosh](#cosh) [deg](#deg) [exp](#exp) [fact](#fact) [fib](#fib) [floor](#floor) [fmod](#fmod) [frexp](#frexp)
               [ldexp](#ldexp) [log](#log) [log10](#log10) [max](#max) [mean](#mean) [median](#median) [min](#min) [mod](#mod) [mode](#mode) [modf](#modf) [pow](#pow) [rad](#rad) [random](#random)
               [randomseed](#randomseed) [range](#range) [rms](#rms) [sin](#sin) [sinh](#sinh) [sqrt](#sqrt) [stddev](#stddev) [sum](#sum) [tan](#tan) [tanh](#tanh) [uuidgen](#uuidgen)
               [uuidgenfmt](#uuidgenfmt) [variance](#variance)


        * [Operating System](#operating-system-procs)  
           [argstr](#argstr) [arglist](#arglist) [argopts](#argopts) [chdir](#chdir) [clearenv](#clearenv) [clock](#clock) [difftime](#difftime) [execute](#execute) [exist](#exist)
           [exit](#exit) [date](#date) [getenv](#getenv) [geteuid](#geteuid) [gethome](#gethome) [getuser](#getuser) [getpid](#getpid) [getppid](#getppid) [getuid](#getuid) 
           [getwd](#getwd) [hostname](#hostname) [remove](#remove) [rename](#rename) [setenv](#setenv) [sleep](#sleep) [stat](#stat) [statfs](#statfs) [time](#time) [tmpname](#tmpname)
           [unsetenv](#unsetenv)

        * [String Handling](#string-procs)  
           [after](#after) [before](#before) [byte](#byte) [char](#char) [contains](#contains) [containsany](#containsany) [count](#count) [decodebase64](#decodebase64) 
           [dump](#dump) [encodebase64](#encodebase64) [find](#find) [format](#format) [gsub](#gsub) [hasprefix](#hasprefix) [bhassuffixye](#hassuffix) [index](#index) 
           [indexany](#indexany) [lastindex](#lastindex) [lastindexany](#lastindexany) [len](#len) [length](#length) [lower](#lower) [match](#match) [prskdtch](#prskdtch) 
           [prxchange](#prxchange) [rep](#rep) [replace](#replace) [reverse](#reverse) [scan](#scan) [scanall](#scanall) [sub](#sub) [substr](#substr) [trim](#trim) [trimleft](#trimleft) 
           [trimprefix](#trimprefix) [trimright](#trimright) [trimspace](#trimspace) [trimsuffix](#trimsuffix) [title](#title) [upper](#upper) 

        * [List Handling](#qList-procs)  
           [dumpl](#dumpl) [getn](#getn) [concat](#concat) [insert](#insert) [maxn](#maxn) [erase](#erase) 
           [marshal](#marshal) [marshalxml](#marshalxml) [unmarshal](#unmarshal) [sort](#sort) 

* [Script examples](#script-examples)  
    * [Example 1 Comments, Variables, Procs](#example-1-comments-variables-procs)  
    * [Example 2 Input from file](#example-2-input-from-file)  
    * [Example 3 Interactive](#example-3-interactive)  
* [Return Codes](#return-codes)  


## Usage Summary

``` bash
q [ command ]
	[-pgm <program> ]            File name containing Q program to run.
	[-exec <script-segment> ]    String of Q language statements.
	[-lib <Q-lib-file> ]         Q library file name to access.
	[-profile <profile-file> ]   File name to use for profile data.
	[-log <log-file> ]           File name to use for logging.
	[-name <name-string> ]       Name tag for logging, and _NAME script variable.
	[-limit <nnn> ]              Sets a memory size limit for Q program.
	[-inter]                     Use interactive mode.
	[-help | -h]                 Display this help then exit.
	[-quiet ]                    Do not show output.
	[-version | -v]              Show version information then exit.
	[-verbose]                   Show more information for functions with output.
	[-debug]                     Show debugging info.

```

## Command and Option Details

### Command Format

`q [ command | script.q ] [ -options ] [ -- script options ]`

If no arguments are specified the Q language interpreter runs interactively, 
reading statements from stdin. Exiting from interactive mode is done using 
bye(), quit(), exit() functions, or using ctrl-C.     

If a file containing an Q script is marked executable and is located in 
the path that is searched by the command shell, and if the first line of the 
script has the following format:
``` d
#!/your-path/bin/q 
// Script: hello
put("Hello World")
```
Then the script can be executed from the shell command line simply by typing 
the script files name. 

## Commands

* `help` - Displays brief help information on stdout. More detailed 
    information is displayed if the -verbose option is also specified.
* `version` - Displays the Q programs version information on stdout.
* `run` - Runs the program named by the -pgm option
* `int` - Run in interactive mode. Allows both Q and script options to be
    passed to the interactive script environment to facilitate option testing.  

## Option Details

#### debug 
Type: bool (default false)
	
Display debugging information during Q program execution.		
#### exec
Type: string (set ) (default )
	
String of Q language statements to execute directly.		
#### h
Type: bool (set true) (default false)
	
Display help information and exit.		
#### help
Type: bool (set true) (default false)
	
Display help information and exit.		
#### inter 
Type: bool (set false) (default false)

File name to use for profile data.		
#### lib 
Type: string (set ) (default )

Q library file name to access.		
#### limit
Type: int (set 0) (default 0)
	
Sets a memory limit for the executing Q program in MB. If 0 is 
used, no limit is set.		
#### log
Type: string (set ) (default )
	
File name to use for logging.		
#### name
Type: string (set Q) (default Q)
	
Name used to name script for logging, and internal reference.		
#### pgm 
Type: string (set ) (default )

Name of file containing Q program.		
#### profile
Type: string (set ) (default )

File name to use for profile data.		
#### quiet
Type: bool (set false) (default false)
	
Hide program output during Q program execution.		
#### v
Type: bool (set false) (default false)
	
Display version information and exit.		
#### verbose 
Type: bool (set true) (default false)
	
Displays more information for many functions that have display output.		
#### version 
Type: bool (set false) (default false)
	
Display version information and exit.
		

## Language syntax

### Comments

Three forms of comment exist, two line formats and one block format:

* The characters   // start a line comment and the end of line marks the comment end. 
* The characters   #! start an alternate line comment. 

    The end of line marks the of the comment. This form can be used on 
	some systems, when used to name the location of Q executable. 
	For example:  #!/usr/local/bin/q
	
* The characters   /* start a block comment and */   ends a block comment.

#### Comment examples
```` d
    // line comment
        a = 55  // end of line comment

    /* a block comment can be part of a line,
      even between statement tokens
      or span multiple lines */
    a = /* 123 */ 42
````

### Reserved words

* `and` - logical operator 
* `break` - statement identifier
* `dcl` - statement identifier
* `do` - statement identifier
* `else` - statement identifier
* `elseif` - statement identifier
* `end` - statement identifier
* `false` - boolean value
* `for` - statement identifier
* `func` - alternate procedure identifier
* `if` - statement identifier
* `in` - for expression list operator
* `nil` - nil value
* `not` - logical operator 
* `or` - logical operator 
* `proc` - procedure identifier
* `repeat` - statement identifier
* `return` - statement identifier
* `then` - statement identifier
* `true` - boolean value
* `until` - statement identifier
* `while`	- statement identifier


### Operators 

| Symbol | Function         | Type    |	Association | Precedence |
| -------|------------------|---------|-------------|------------|
| `+`    | addition         | binary  |	left        | 5          |
| `-`    | subtraction      | binary  | left        | 5          |
| `*`    | multiplication   | binary  |	left        | 6          |
| `/`    | division         | binary  |	left        | 6          |
| `%`    | modulus          | binary  |	left        | 6          |
| `and`  | logical and      | binary  | left        | 2          |
| `or`   | logical or       | binary  | left        | 1          |
| `<`    | less             | binary  | left        | 3          |
| `>`    | greater          | binary  | left        | 3          |
| `<=`   | less or equal    | binary  | left        | 3          |
| `>=`   | greater or equal | binary  | left        | 3          |
| `!=`   | not equal        | binary  | left        | 3          |
| `==`   | equal            | binary  | left        | 3          |
| `!!`   | concatenate      | binary  | right       | 4          |
| `#`    | length           | unary   | left        | 7          |
| `not`  | logical not      | unary   | right       | 7          |
| `-`    | unary negative   | unary   | right       | 7          |
| `^`    | exponent         | binary  | right       | 8          |

A higher precedence is interpreted ahead of a lower precedence.

### Variables

* Variable names can use the following characters: A..Z, a..z, 0..9, and _
* A variable can be created by initial assignment which will create a 
  globally accessible variable seen by all procs, or a variable can be 
  created using the dcl statement which will make the variable accessible 
  only by the proc that executed the dcl statement.
* Variable types can be:
  * string s = "Hello Zaphod"
  * number n = 42
  * list   a = {1,2,3}  or  b = {e=33,t="time",c=os.clock()}
  * proc   f = os.getenv("PATH")

Variables are created by assigning by assigning a value on the right hand 
side of the assignment operator '=' to a variable name on the left hand side
of the '='. For example: 
```
> ME=44
> Me="Text"
> me={23.69,"Twenty Three"}
> put(ME,Me,me)
44      Text    list: 0xc000548840
> put(type(ME),type(Me),type(me))
num     str     list
> put(#Me,#me)
4       2
```  
Notice the variable names are case sensitive, and the type of the variable
is set when the variable is created or re-created. The length operator '#'
is only valid for the string and list variable types. 

### Control Structures

The following logical control directives are available:

* `while` <_expression_> `do` <_block_> `end`
* `repeat` <_block_> `until` <_expression_>
* `if` <_expression_> `then` <_block_> `end`
* `if` <_expression_> `then` <_block_> `else` <_block_> `end`
* `if` <_expression_> `then` <_block_> `elseif` <_expression_> `then` <_block_> `else` <_block_> `end`
* `for name` = <_expression_> , <_expression_> [, <_expression_>] `do` <_block_> `end` 
* `for` <_variable_> [, <_variable_>] `in` <_expression_> [, <_expression_>]  `do` <_block_> `end`
* `break`
* `return` [_values_]

#### while

`while` <_expression_> `do` <_block_> `end`

For example:
```
> i = 5
> while i <= 12 do
>> put(i,fact(i))
>> i=i+1
>> end
5       120
6       720
7       5040
8       40320
9       362880
10      3628800
11      39916800
12      479001600
```

#### repeat

`repeat` <_block_> `until` <_expression_>

For example:
```
> n = 15
> repeat
>> f = fib(n)
>> put(n,f)
>> n = n-1
>> until n < 6
15      610
14      377
13      233
12      144
11      89
10      55
9       34
8       21
7       13
6       8
```

#### if_then

`if` <_expression_> `then` <_block_> `end`

For example:
```
> func roll()
>> v = random(1,6)
>> return v
>> end
>
> color =  "black"
> if roll() == 4 then
>> color = "yellow"
>> end
> put(color)
black
```

#### if_then_else

`if` <_expression_> `then` <_block_> `else` <_block_> `end`

For example:
```
> func roll()
>> v = random(1,6)
>> return v
>> end
>
>  if roll() < 3 then
>> color = "red"
>> else
>> color = "blue"
>> end
> put(color)
blue
```

#### if_then_elseif_else

`if` <_expression_> `then` <_block_> `elseif` <_expression_> `then` <_block_> `else` <_block_> `end

For example:
```
> func roll()
>> v = random(1,6)
>> return v
>> end
>
> r = roll()
> if r == 1 then color = "red"
>> elseif r == 2 then color = "green"
>> elseif r == 3 then color = "blue"
>> elseif r == 4 then color = "yellow"
>> elseif r == 5 then color = "orange"
>> else color = "black"
>> end
> put(r,color)
2       green
```

#### for

`for name` = <_start expression_> , <_end expression_> [, <_increment expression_>] `do` <_block_> `end`

For example:
```
> func roll()
>> v = random(1,6)
>> return v
>> end
>
> r = roll()
> for val=1,r do   // 1=start val and r=end val
>> put(val)
>> end
1
2
3
4
5
> // Golden ratio - NB last converging iterations are very expensive 
> for i=4,44,4 do   // 4=start i, 44=end i, and 4=increment
>> put(i,i+1,fib(i+1)/fib(i))
>> end
4       5       1.6666666666666667
8       9       1.619047619047619
12      13      1.6180555555555556
16      17      1.618034447821682
20      21      1.6180339985218033
24      25      1.618033988957902
28      29      1.6180339887543225
32      33      1.618033988749989
36      37      1.618033988749897
40      41      1.618033988749895
44      45      1.618033988749895
>
```

#### for_in

`for` <_variable_> [, <_variable_>] `in` <_expression_> [, <_expression_>]  `do` <_block_> `end`

For example:
```
> q = {how=3,many=5,roads=9,must=2,man=6,walk=1}
> for k,v in pairs(q) do
>> put(k,v)
>> end
roads   9
must    2
man     6
walk    1
how     3
many    5
```

#### break

`break`

For example:
```
> for i=4,44,4 do
>> c1 = clock()
>> put(i,i+1,fib(i+1)/fib(i))
>> c2 = clock()
>> tm = c2 - c1
>> put("Time:",tm,"(s)")
>> if tm > 1.0 then break end
>> end
4       5       1.6666666666666667
Time:   4.878699996879732e-05   (s)
8       9       1.619047619047619
Time:   1.6770999991422286e-05  (s)
12      13      1.6180555555555556
Time:   2.134700002898171e-05   (s)
16      17      1.618034447821682
Time:   7.099399999788147e-05   (s)
20      21      1.6180339985218033
Time:   0.00028661199985435815  (s)
24      25      1.618033988957902
Time:   0.0018965350000144099   (s)
28      29      1.6180339887543225
Time:   0.01369292399999722     (s)
32      33      1.618033988749989
Time:   0.09743848399989474     (s)
36      37      1.618033988749897
Time:   0.6728419630001099      (s)
40      41      1.618033988749895
Time:   4.602405675     (s)
>
```

### Q Language Procedures and Functions

Q script has many procedures (procs) and functions built in. There are standard
procs and auxiliary procs. The standard procs do not need a name space prefix
to be referenced. 

For example: using proc bye to exit the script is specified:
    `bye()`
However some IO procs (see below) need to be specified with a prefix 
like “i” shown here:
     `f = i.open("test.txt","r")`

#### Standard procs

##### assert
```
[z:bool =] assert(a:bool,b:str)
```

Asserts 'a' is true and if not issues message 'b' Returns 'a'.

##### bye
```
bye([a:num])
```

Stop the executing script, returns code 'a' to system.
```
$ q
q, Version - 2.0.0, Build Date - 16JUL2021
q interactive mode, enter exit() to exit, or help() for help.
> bye()
```

##### collectgarbage
```
collectgarbage( [a:str("stop"|"restart"|"collect"|"step"|"count"]) )
```

Runs specified mode of garbage collection to free used resources.

##### error
```
error(a:str)
```
Issues error message 'a'.
```
> for i=1,10 do
>> put("I=",i)
>> if i==5 then error("Bad number 5") end
>> end
I=      1
I=      2
I=      3
I=      4
I=      5
2019-06-17 11:14:34.922 ERROR [qmain.go:515] [Q] - Load segment error <string>:5: Bad number 5
stack traceback:
        [G]: in proc 'error'
        <string>:5: in main segment
        [G]: ?
```

##### getfenv 
```
z = getfenv(a:proc)
```

Returns the current environment into 'z' for the proc 'a'.

##### getmetalist
```
z = getmetalist(a:list)
```

Returns meta info into 'z' for list 'a'.

##### help
```
help()
```

Displays the help text.

##### load
```
z = load(a:proc [,b:str])
```
Loads proc 'z', using proc"nil", "bool", "num", "str", "proc", 
"data", "thread", "list", "chan" 'a' to load multiple segment 
strings containing a Q script.

##### loadfile
```
z = loadfile(a:str)
```

Loads proc 'z' from a file named 'a' which contains an Q script.
```
> f=i.open("func.dat","w")
> f:write("q=sin(6);return q")
> f:close()
> sin6 = loadfile("func.dat")
> put(sin6())
-0.27941549819892586
```

##### loadstring
``` 
z = loadstring(a:str)
```

Loads proc 'z' from a string 'a'.
```
> f = loadstring("c=cos(22);return c")
> put(f)
proc: 0xc00056b600
> put(f())
-0.9999608263946371
```

##### logd
```
logd(a:str)
```

Writes message string 'a' on the log as a debug message, only if the
-debug option is turned on.
```
> logd("How many roads")
```
In this example, no output was produced on stderr as -debug option was 
not on.

##### loge
```
loge(str)
```

Writes message string 'a' on the log as an error message.
```
> loge("This is it")
2019-06-11 10:05:18.118 ERROR [qslibbase.go:467] [Q] - This is it
```

##### logi
```
logi(str)
```

Writes message string 'a' on the log as an info message.
```
> logi("Excellent Frankie")
2019-06-11 10:15:25.542  INFO [qslibbase.go:447] [Q] - Excellent Frankie
```

##### logw
```
logw(str)
```

Writes message string 'a' on the log as a warning message
```
> logw("Sounds good, but doesn't actually mean anything!")
2019-06-11 10:17:34.998  WARN [qslibbase.go:457] [Q] - Sounds good, but doesn't actually mean anything!
```

##### next
```
i,v = next(a:list,b:*)
```

Get the next index 'i' and value 'v' from the list 'a'.
```
> qs = {"one","two","three","four","five","six"}
> for i,v in next, qs do
>> put(i,v)
>> end
1       one
2       two
3       three
4       four
5       five
6       six
>
```

##### pcall
```
ok,r = pcall(a:proc,b:*,c:*,...)
```

Call proc 'a' with args 'b', 'c', etc. Catch any runtime 
errors. If no error 'ok' is true.
```
> proc divide(a,b)
>> r = a / b
>> return r
>> end
>
> ok,r = pcall(divide,20,2)
> put(ok,r)
true    10
> ok,r = pcall(divide,20,0)
> put(ok,r)
true    +Inf
> ok,r = pcall(divide,20,"bad")
> put(ok,r)
false   <string>:3: cannot perform div operation between num and str
>
```

##### put
```
put(a:str,b:str,...)
```

Write strings 'a', 'b', etc to stdout.
```
> put("Forty two",42,a,log(a))
Forty two       42      4.477336814478207       1.4990284086108943
> put("One\nTwo\nThree\n")
One
Two
Three

> put("a\tb\nc\td")
a       b
c       d
```

##### quit
```
quit([a:num])
```

Stpqt the executing script, returns code 'a' to system.

##### rawequal
```
z:bool = rawequal(a:*, b:*)
```

Checks if 'a' equals 'b' and returns true in 'z' if they are equal.
```
> a=log(88)
> b=log10(66)
> q=rawequal(a,b)
> put(q)
false
> put(rawequal(1,1))
true
> c=1
>  put(rawequal(1,c))
true
> d=c
>  put(rawequal(d,c))
true
>  put(rawequal(1,1.0))
true
```

##### rawget
```
z = rawget(a:list, i:*)
```

Gets the real value of 'a[i]' into 'z'.
```
> list = {1,2,3,"smith",5,6,7,8,"dog"}
> val = rawget(list,4)
> put(val)
smith
```

##### rawset
```
rawset(a:list, i:*, v:*)
```

Sets the value 'v' into list 'a' at index 'i'.
```
> list = {1,2,3,4,5,6}
> rawset(list,3,"junk")
> dumpl(list)
(array): [
(num)1: 1
(num)2: 2
(str)3: "junk"
(num)4: 4
(num)5: 5
(num)6: 6
]
```

##### run
```
run(a:str)
```

Run the Q srcipt from file name 'a'.
```
> f = i.open("test.q","w")
> f:write('qs = {"one","two","three","four","five","six"}')
> f:write("\n")
> f:write("for i,v in pairs(qs) do \n")
> f:write("put(i,v) \n")
> f:write("end \n")
> f:close()
> run("test.q")
1       one
2       two
3       three
4       four
5       five
6       six
>
```

##### stop
```
stop([a:num])
```

Stpqt the executing script, returns code 'a' to system.
```
> stop(99)
dingo@mach:~ > echo $?
99
```

##### tonumber
```
n = tonumber(s:str)
```

Returns string 's' converted to a number.
```
> a="77"
> put(a)
77
> b=tonumber(a)
> put(b)
77
> put(type(a))
str
> put(type(b))
num
```

##### tostring
```
* `s = tostring(n:num)
```

Returns number 'n' converted to a string.
```
> c=cos(42)
> put(c)
-0.3999853149883513
> d=tostring(c)
> put(substr(d,3,5))
.3999
>
> proc fibstr(n)
>> s = ""
>> for i=0,n do
>> s = s || tostring(fib(i))
>> end
>> return s
>> end
>
> put(fibstr(40))
01123581321345589144233377610987159725844181676510946177112865746368750251213931964183178115142298320401346269217830935245785702887922746514930352241578173908816963245986102334155

```

##### type
```
t = type(a:*)
```

Returns the type of 'a', one of: "nil", "bool", "num", "str", "proc", 
"data", "thread", "list", "chan".
```
> put(type(22))
num
> put(type(false))
bool
> put(type("OK"))
str
> put(type(log10))
proc
> put(type({3,4.9,"Q",007}))
list
> put(type(undefinedvar))
nil

```

##### unpack
```
a[i], ..., a[j] = unpack(a:list[,i:num[,j:num]])
```

Returns a slice of the list 'a' from 'i' to 'j'
```
> list = {1,2,3,"smith",5,6,7,8,"dog","pig","cat",1000,2000}
> a,b,c,d,e = unpack(list,7,11)
> put(a,b,c,d,e)
7       8       dog     pig     cat
```

##### xpcall
```
ok = xpcall(a:proc, e:proc)
```

Call proc 'a' catching any runtime errors. If an error occurs 
'ok' is set to false and proc 'e' is called.
```
> proc err()
>> res=0
>> put("Error in add() proc - result 0")
>> end
>
> proc add()
>> res=num+5
>> end
>
> num=3
> ok = xpcall(add,err)
> put(ok,res)
true    8
>
> num="bad"
> ok = xpcall(add,err)
Error in add() proc - result 0
>
> num=7
> ok = xpcall(add,err)
> put(ok,res)
true    12
```

#### Input and Output procs

##### close
```
f.close()
```

Close file handle 'f'.
```
>  f=i.open("my.txt","r")
>  f:close()
```

##### flush
```
i.flush(f)
```

Save any buffered data written to file handle 'f' to disk.
```
> out = i.open("out.txt","w")
> for i=1,10 do
>> out:write(tostring(i))
>> out:flush()
>> end
> out:close()
> stop()
dingo@mach:~ > cat out.txt
12345678910
```

##### lines
```
s:str = i.lines(f)
```

Read next line 's' from file handle 'f'.
```
> f=i.open("my.txt","r")
> put(type(f))
data
> q=f:lines()
> put(type(q))
proc
> l=q()
> put(l)
This
>  l=q()
>  put(l)
nil
```

##### input
```
f = i.input([f|n:str])
```

Returns default input file or assigns new default from 
name 'n' or handle 'f'.
```
> h=i.input()
> put(type(h))
data
```

##### output
```
f = i.output([f|n:str])
```

Returns default output file or assigns new default from 
name 'n' or handle 'f'.
```
> f = i.output("out.txt")
> f:write("This is a test\n")
> f:close()
>
> z = i.open("out.txt")
> v = z:read()
> put(v)
This is a test
```

##### open
```
f = i.open(n:str[,m:str])
```

Opens file name 'n' using mode 'm', returns file handle 'f'. 
Mode can be "r", "w", "a", "r+", "w+", "a+", indicating read, 
write, append, and "+" indicates update mode. Modes may have 
"b" concatenated for binary.
```
> f=i.open("my.txt","r")
> put(type(f))
data
```

##### read
```
d = f:read()
```
Read data from file 'f' into 'd'.
```
>  f=i.open("my.txt","r")
> dat=f:read()
> put(type(dat))
str
> put(dat)
This
> f:close()
```

##### type
```
f:type()
```

##### tmpfile
```
f:tmpfile()
```

Returns handle 'f' to a temporary file, opened in update mode. 
The file is deleted on Q script termination.
```
> t=i.tmpfile()
> start = t:seek("cur")
> t:write("This\nis\nit")
> t:seek("set",start)
> dat = t:read()
> put(dat)
This
> dat = t:read()
>  put(dat)
is
> dat = t:read()
>  put(dat)
it
> dat = t:read()
> put(dat)
nil
```

##### write
```
f:write()
```

Write a string into file 'f'
```
> f=i.open("my.txt","w")
> f:write("This\n")
> f:close()
> exit()
dingo@mach:~ > cat my.txt
This
```


### Mathematics procs and constants
#### Math Constants

* `pi` The constant for pi.
* `e` The constant for e.
* `phi` The constant for Phi.
* `sqrt2` The constant for the square root of two.
* `sqrte` The constant for the square root of e.
* `huge` The constant for the largest number available.
* `small` The constant for the largest number available.

The following copies all the constant values into a key=value list and then
uses the dumpl() builtin proc to print the list of constants.
```
> const = {pi=pi,e=e,phi=phi,sqrt2=sqrt2,sqrte=sqrte,huge=huge,small=small}
> dumpl(const)
(strdict): {
(num)sqrte: 1.6487212707001282
(num)huge: 1.7976931348623157e+308
(num)small: 5e-324
(num)pi: 3.141592653589793
(num)e: 2.718281828459045
(num)phi: 1.618033988749895
(num)sqrt2: 1.4142135623730951
}
```

#### Math Procs

##### abs
```
z:num = abs(a:num)
```

Returns the absolute value of 'a' in 'z'.
```
> a = abs(-88.234)
> put(a)
88.234
> b = abs(44)
> put(b)
44
> c = abs(-12)
> put(c)
12
```

##### acos
```
z:num = acos(a:num)
```

Returns the arc cosine value of 'a' in 'z'.
```
>  for i=0,1,0.1 do put(acos(i)) end
1.5707963267948966
1.4706289056333368
1.3694384060045657
1.266103672779499
1.1592794807274085
1.0471975511965976
0.9272952180016123
0.7953988301841436
0.6435011087932845
0.4510268117962626
1.4901161193847656e-08
```

##### asin
```
z:num = asin(a:num)
```

Returns the arc sine value of 'a' in 'z'.

##### atan
```
z:num = atan(a:num)
```

Returns the arc tangent value of 'a' in 'z'.

##### ceil
```
z:num = ceil(a:num)
```

Returns the smallest integer larger than or equal to 'a' in 'z'.
```
> a = 5.856
> b = ceil(a)
> put(a,b)
5.856   6
```

##### cos
```
z:num = cos(a:num)
```

Returns the cosine value of 'a' in 'z'.
```
> c = cos(rad(45))    // cosine of 45 degrees
> put(c)
0.7071067811865476
```

##### cosh
```
z:num = cosh(a:num)
```

Returns the hyperbolic cosine value of 'a' in 'z'.

##### deg
```
z:num = deg(a:num)
```

Returns the angle of 'a' specified in radians as 'z' in degrees.
```
> d = deg(3.141)
> put(d)
179.96604345059157
```

##### exp
```
z:num = exp(a:num)
```

Returns the value of e raised to the power of 'a' in 'z'.
```
> e5 = exp(5)
> put(e5)
148.4131591025766
```

##### fact
```
z:num = fact(a:num)
```

Returns the factorial of 'a' in 'z'.
```
>  for i=0,9 do put(i,fact(i)) end
0       1
1       1
2       2
3       6
4       24
5       120
6       720
7       5040
8       40320
9       362880
```

##### fib
```
z:num = fib(a:num)
```

Returns the Fibonacci number for 'a' in 'z'.
```
> for i=0,9 do put(i,fib(i)) end
0       0
1       1
2       1
3       2
4       3
5       5
6       8
7       13
8       21
9       34
```

##### floor
```
z:num = floor(a:num)
```

Returns the largest integer smaller than or equal to 'a' in 'z'.
```
> f = floor(4.987)
> put(f)
4
```

##### fmod
```
z:num = fmod(a:num,b:num)
```

Returns the remainder of 'a' divided by 'b' in 'z'.
```
> m = fmod(19,7)
> put(m)
5
```

##### frexp
```
z:num,y:num = frexp(a:num)
```

Returns 'z' and 'y' such that 'a' = z2**e, where e is an integer and 
the absolute value of 'z' is in the range [0.5, 1] or 0 when 'a' is 0.
```
> e,f = frexp(4)
> put(e,f)
0.5     3
```

##### ldexp
```
z:num = ldexp(a:num,b:num)
```

Returns 'a'*2**'b' in 'z'.
```
> a = ldexp(3,4)
> put(a)
48
```

##### log
```
z:num = log(a:num)
```

Returns the natural log of 'a' in 'z'.
```
> a=log(88)
> put(a)
4.477336814478207
```

##### log10
```
z:num = log10(a:num)
```

Returns the log base 10 of 'a' in 'z'.
```
> b=log10(66)
> put(b)
1.8195439355418686
```

##### max
```
z:num = max(a:num, b:num [,...])
```

Returns the maximum of 'a', 'b' etc in 'z'.
```
> mx = max(4,23,5.6,29,34,2,7.823,9,15.7,33)
> put(mx)
34
```

##### mean
```
z:num = mean(a:num, b:num [,...])
```

Returns the mean of 'a', 'b' etc in 'z'.
```
> x = mean(12,15,21,18,9,14)
> put(x)
14.833333333333334
```

##### median
```
z:num = median(a:num, b:num [,...])
```

Returns the median of 'a', 'b' etc in 'z'.
```
> md = median(4,23,5.6,29,34,2,7.823,9,15.7,33)
> put(md)
19.35
> md = median(4,23,5.6,29,34,2,7.823,9,15.7,33,58)
> put(md)
23
```

##### min
```
z:num = min(a:num, b:num [,...])
```

Returns the minimum of 'a', 'b' etc in 'z'.
```
> mn = min(4,23,5.6,29,34,2,7.823,9,15.7,33)
> put(mn)
2
```

##### mod
```
z:num = mod(a:num,b:num)
```

Returns the remainder of 'a' divided by 'b' in 'z'.
```
> a = mod(11,3)
> put(a)
2
```

##### mode
```
z:num = mode(a:num, b:num [,...])
```

Returns the mode of 'a', 'b' etc in 'z'.
```
> mo = mode(4,23,5.6,29,34,2,7.823,9,15.7,33,58,23,5,8,23,31)
> put(mo)
23
```

##### modf
```
z:num,y:num = modf(a:num)
```

Returns the integer part of 'a' in 'z' and the fractional part in 'y'.
```
> a,b = modf(11.683)
> put(a,b)
11      0.6829999999999998
```

##### pow
```
z:num = pow(a:num,b:num)
```

Returns 'a'**'b' in 'z'.
```
> p = pow(5,3)
> put(p)
125
```

##### rad
```
z:num = rad(a:num)
```

Returns the angle 'a' specified in degrees in 'z' in radians.
```
> a = sin(rad(45))
> put(a)
0.7071067811865475
```

##### random
```
z:num = random([a:num[,b:num]])
```

Returns a random number in 'z', in range 0..1, 1..'a', or 'a'..'b' 
depending on wether 'a' and 'b' are supplied. 'a' must be greater than 'b'.
```
> randomseed(675487)
> put(random(1,6))
5
> put(random(1,6))
5
> put(random(1,6))
2
> put(random(1,6))
6
> put(random(1,6))
3
> put(random(1,6))
4
> put(random(1,6))
3
```

##### randomseed
```
randomseed(a:num)
```

Sets the internal random seed. 'a' must be a number.

##### range
```
z:num = range(a:num, b:num [,...])
```

Returns the range of 'a', 'b' etc in 'z'.
```
> ra = range(29,34,7.823,9,15.7,33)
> put(ra)
26.177
```

##### rms
```
z:num = rms(a:num, b:num [,...])
```

Returns the root mean square of 'a', 'b' etc in 'z'.
```
> rm = rms(4,23,5.6,29,34,2,7.823,9,15.7,33)
> put(rm)
20.13715304853196
```

##### sin
```
z:num = sin(a:num)
```

Returns the sine value of 'a' in 'z'.
```
> v = sin(pi/2)
> put(v)
1
> v = sin(pi/4)
> put(v)
0.7071067811865475
```

##### sinh
```
z:num = sinh(a:num)
```

Returns the hyperbolic sine value of 'a' in 'z'.

##### sqrt
```
z:num = sqrt(a:num)
```

Returns the square root of 'a' in 'z'.
```
> q = sqrt(43)
> put(q)
6.557438524302
> q = sqrt(4)
> put(q)
2
```

##### stddev
```
z:num = stddev(a:num, b:num [,...])
```

Returns the standard deviation of 'a', 'b' etc in 'z'.
```
> s = stddev(4,23,5.6,29,34,2,7.823,9,15.7,33)
> put(s)
12.446052547338498
```

##### sum
```
z:num = sum(a:num, b:num [,...])
```

Returns the sum of 'a', 'b' etc in 'z'.
```
> su = sum(4,23,5.6,29,34,2,7.823,9,15.7,33)
> put(su)
163.123
```

##### tan
```
z:num = tan(a:num)
```

Returns the tangent value of 'a' in 'z'.
```
> t = tan(rad(45))
> put(t)
1
> t = tan(rad(60))
> put(t)
1.7320508075688767
> t = tan(rad(89))
> put(t)
57.28996163075987
```

##### tanh
```
z:num = tanh(a:num)
```

Returns the hyperbolic tangent value of 'a' in 'z'.
```
> k = tanh(2)
> put(k)
0.9640275800758169
```

##### uuidgen
```
z:str = uuidgen()
```

Returns a generated UUID in 'z'.
```
> s = dump(uuidgen())
> put(s)
00000000  f3 e0 b5 f7 3c d1 4d 9d  95 d1 24 1d 04 41 84 7f  |....<.M...$..A..|
```

##### uuidgenfmt
```
z:str = uuidgenfmt()
```

Returns a generated and formatted UUID in 'z'.
```
> w = uuidgenfmt()
> put(w)
E68897FD-F852-4E74-94A4-33712D0C8801
```

##### variance
```
z:num = variance(a:num, b:num [,...])
```

Returns the variance of 'a', 'b' etc in 'z'.
```
> v = variance(4,23,5.6,29,34,2,7.823,9,15.7,33)
> put(v)
154.90422401111113
```

#### Operating System procs

##### argstr
```
z:str = argstr()
```

Returns the arguments passed to the script as a string in 'z'.
```
> q int -name test -- -this "is it" -how 99 -many roads -are 'there on the way'
q, Version - 2.0.0, Build Date - 16JUL2021
q interactive mode, enter exit() to exit, or help() for help.
> args = argstr()
> put(args)
-this "is it" -how 99 -many roads -are "there on the way"
> stop(0)
>

```

##### arglist
```
z:str = arglist(a:str)
```

Returns the arguments from string 'a' parsed as a list in 'z'.
```
> q int -name test -- -this "is it" -how 99 -many roads -are 'there on the way'
q, Version - 2.0.0, Build Date - 16JUL2021
q interactive mode, enter exit() to exit, or help() for help.
> args = argstr()
> put(args)
-this "is it" -how 99 -many roads -are "there on the way"
> alist = arglist(args)
> dumpl(alist)
(array): [
(str)1: "-this"
(str)2: "is it"
(str)3: "-how"
(str)4: "99"
(str)5: "-many"
(str)6: "roads"
(str)7: "-are"
(str)8: "there on the way"
]
> stop(0)
>

```

##### argopts
```
z:str = argopts(a:str)
```

Returns the arguments from string 'a' parsed as a list of named 
options and arguments in 'z'.
```
> q int -name test -- -this "is it" -how 99 -many roads -are 'there on the way'
q, Version - 2.0.0, Build Date - 16JUL2021
q interactive mode, enter exit() to exit, or help() for help.
> args = argstr()
> put(args)
-this "is it" -how 99 -many roads -are "there on the way"
> opts = argopts(args)
> dumpl(opts)
(array): [
(str)how: "99"
(str)many: "roads"
(str)this: "is it"
(str)are: "there on the way"
]
> stop(0)
>
```

##### chdir
```
z:bool = chdir(a:str)
```

Changes the current working directory to 'a'. Returns true in 
'z' if the directory change was a success.

##### clearenv
```
clearenv()
```

Deletes all environment variables.

##### clock
```
z:num = clock()
```

Returns the time in seconds since the script started executing.
```
$ q
q, Version - 2.0.0, Build Date - 16JUL2021
q interactive mode, enter exit() to exit, or help() for help.
> put(clock())
13.512708397
> put(clock())
16.848223078
> put(clock())
19.120163765
> stop()
$
```

##### difftime
```
difftime(a:num,b:num)
```

Returns the difference between time 'a' and time 'b' in 'z'.

##### execute
```
z:num = execute(a:str)
```

Executes command 'a' as an OS process. Returns 0 if no errors 
occur or 1 if there were errors.
```
> rc = execute("ls")
a_demo  b_demo  c_demo  d.conf  func.dat  omg.yaml       rat_etl.json
a.json  b.yaml  c.xml   d_demo  omg       omg.yaml.json  ubiquity.pid
> put(rc)
0
```

##### exist
```
z:bool = exist(f:str)
```

Checks the existence of file 'f', returns true if file 'f' 
exists, otherwise returns false.
```
> e = exist("func.dat")
> put(e)
true
```

##### exit
```
exit([a:num])
```

Stpqt executing the Q script, and exits to the OS, optionally 
returning the code in 'a' to the OS shell.

##### date
```
z = date([a:str[,b:num]])
```

Returns a date in a list, or as a number depending om the format 
string 'a'. If the format 'a' is not supplied the date and time 
are in local time, as the number of seconds from the epoch. If the 
format string 'a' begins with "!" the datetime is UTC. If the 
format string 'a' begins with "*t" the date time is returned in 
a list containing elements: "year","month","day","hour", "min",
"sec","weekday","yearday","isdst".
```
> d=date()
> put(d)
13 Jun 19 12:32 EDT
> d=date("*t")
> put(d)
list: 0xc0001b0720
> dumpl(d)
(array): [
(num)min: 39
(num)yearday: 0
(bool)isdst: false
(num)month: 6
(num)year: 2019
(num)day: 13
(num)sec: 8
(num)hour: 12
(num)weekday: 4
]
> s = date("Today is %A in %B")
> put(s)
Today is Thursday in June
```

Format patterns for dates

| Pattern | Formatting function         |
|---------|-----------------------------|
| %a      | short week day name         |
| %A      | week day name               |
| %b      | short month name            |
| %B      | month name                  |
| %c      | date and time               |
| %d      | day of month                |
| %H      | hour mod 24                 |
| %I      | hour mod 12                 |
| %M      | minute                      |
| %m      | month                       |
| %p      | am or pm                    |
| %S      | second                      |
| %w      | weekday                     |
| %x      | date                        |
| %X      | time                        |
| %Y      | year four digits            |
| %y      | year two digits             |
| %%      | character %                 |


##### getenv
```
z:str = getenv(a:str)
```

Returns the value of the environment variable named 'a' in 'z'.
```
> env = getenv("PATH")
> put(env)
/usr/local/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/usr/local/go/bin:/home/dingo/bin
> dirs=scanall(env,":")
> dumpl(dirs)
(array): [
(str)1: "/usr/local/bin"
(str)2: "/usr/bin"
(str)3: "/usr/local/sbin"
(str)4: "/usr/sbin"
(str)5: "/usr/local/go/bin"
(str)6: "/home/dingo/bin"
]
```

##### geteuid
```
z:num = geteuid()
```

Returns the effective user ID as a number in 'z'.
```
> euid = geteuid()
> put(euid)
2306
```

##### gethome
```
z:str = gethome()
```

Returns the users home path string in 'z'. This is system 
dependent. On some systems this may be the users name.

##### getuser
```
z:str = getuser()
```

Returns the user ID string in 'z'.
```
> user = getuser()
> put(user)
dingo
```

##### getpid
```
z:num = getpid()
```

Returns the current process ID as a number in 'z'.
```
> pid = getpid()
> put(pid)
26677
> ps = execute("ps -a")
  PID TTY          TIME CMD
 2032 pts/1    00:00:00 dbus-launch
 2642 pts/0    00:00:00 dbus-launch
 9774 pts/5    00:00:00 ps
18352 pts/1    00:00:00 su
18387 pts/1    00:00:00 bash
19222 pts/4    00:00:00 su
19562 pts/4    00:00:00 bash
26677 pts/5    00:00:00 q
```

##### getppid
```
z:num = getppid()
```

Returns the current processes parent process ID as a number in 'z'.
```
> ppid = getppid()
> put(ppid)
30246
> ps = execute("ps -ef | grep 30246")
dingo    23477 26677  0 13:01 pts/5    00:00:00 /bin/sh -c ps -ef | grep 30246
dingo    23479 23477  0 13:01 pts/5    00:00:00 grep 30246
dingo    26677 30246  0 12:50 pts/5    00:00:00 q
dingo    30246 30244  0 Jun07 pts/5    00:00:00 -bash
```

##### getuid
```
z:num = getuid()
```

Returns the user ID as a number in 'z'.
```
> uid = getuid()
> put(uid)
2306
```

##### getwd
```
z:str = getwd()
```

Returns the current working directory in 'z'.
```
> wd = getwd()
> put(wd)
/home/dingo
```

##### hostname
```
z:str = hostname()
```

Returns the host name in 'z'.
```
> host = hostname()
> put(host)
mach.test.github.com
```

##### remove
```
z:bool = remove(a:str)
```

Removes the file name 'a'. If successful true is returned in 'z'.
```
> proc ls()
>> execute("ls")
>> end
>
> ls()
other.file
> remove("other.file")
> ls()
>
```

##### rename
```
z:bool = rename(a:str,b:str)
```

Renames file name 'a' to file name 'b'. If successful true 
is returned in 'z'.
```
> rc = execute("touch some.file; ls -ial")
total 16
   653228 drwxr-xr-x   2 dingo dingo   80 Jun 13 13:05 .
134218195 drwxr-xr-x. 25 root  root  4096 Jul 19  2018 ..
  6417324 -rw-------   1 dingo dingo  876 Jun  7 16:31 .bash_history
  1993727 -rw-rw-r--   1 dingo dingo  201 Apr  3 15:13 .bash_profile
  1993741 -rw-rw-r--   1 dingo dingo  602 Apr  3 14:42 .bashrc
  1780913 -rw-rw-r--   1 dingo dingo    0 Jun 13 13:06 some.file
> ok = rename("some.file","other.file")
> put(ok)
true
> rc = execute("ls -ial")
total 16
   653228 drwxr-xr-x   2 dingo dingo   81 Jun 13 13:08 .
134218195 drwxr-xr-x. 25 root  root  4096 Jul 19  2018 ..
  6417324 -rw-------   1 dingo dingo  876 Jun  7 16:31 .bash_history
  1993727 -rw-rw-r--   1 dingo dingo  201 Apr  3 15:13 .bash_profile
  1993741 -rw-rw-r--   1 dingo dingo  602 Apr  3 14:42 .bashrc
  1780913 -rw-rw-r--   1 dingo dingo    0 Jun 13 13:06 other.file
```

##### setenv
```
z:bool = setenv(a:str,b:str)
```

Sets environment variable name 'a' to value 'b'. If successful 
true is returned in 'z'.
```
> tst = getenv("TEST")
> put(tst)
nil
> ok = setenv("TEST","My dog skip")
> put(ok)
true
> tst = getenv("TEST")
> put(tst)
My dog skip
> rc = execute("echo $TEST")
My dog skip
> exit()
dingo@mach:~ > echo $TEST

dingo@mach:~ >
```

##### sleep
```
sleep(a:num)
```

The script stpqt executing for 'a' seconds.
```
> put(clock()) ; sleep(5) ; put(clock())
40.754074187
45.754221194
```

##### stat
```
z:list = stat(f:str)
```

Returns a list in 'z' containing file information for file 'f'.
```
> wd=getwd()
> st=stat(wd)
> dumpl(st)
(array): [
(num)rdev: 0
(str)atim: "13 Jun 19 13:09 EDT"
(str)ctim: "13 Jun 19 13:08 EDT"
(num)ino: 653228
(num)size: 81
(num)nlink: 2
(num)gid: 2306
(str)name: "/home/dingo"
(num)mode: 16877
(num)blocks: 0
(str)mtim: "13 Jun 19 13:08 EDT"
(num)dev: 2050
(num)blksize: 4096
(str)type: "Directory"
(num)uid: 2306
]
```

##### statfs
```
z:list = statfs(p:str)
```

Returns a list in 'z' containing file system information for path 'p'.
```
> stroot = statfs("/")
> dumpl(stroot)
(array): [
(num)freebytes: 196484440064
(num)percentfreenodes: 99.69312421875
(num)percentusedbytes: 25.010518417195897
(num)bfree: 47969834
(num)frsize: 4096
(num)ffree: 127607199
(num)usedblocks: 15998916
(num)fsid2: 0
(num)percentusednodes: 0.3068757812499996
(num)bsize: 4096
(num)bavail: 47969834
(num)totalbytes: 262016000000
(num)files: 128000000
(str)typename: "XFS"
(num)namelen: 255
(num)usedbytes: 65531559936
(num)percentfreebytes: 74.9894815828041
(str)name: "/"
(num)fsid1: 2050
(num)type: 1481003842
(num)blocks: 63968750
]
```

##### time
```
z = time([a:list])
```

Returns the time in 'z'. If list 'a' is specified, 'z' returns 
a list, else 'z' is a Unix time value as a number.
```
> tm = time()
> put(tm)
1560446536
> t = time({year=1960,month=1,day=16,hour=6,min=30,sec=0})
> put(t)
-314281800
```

##### tmpname
```
z:str = tmpname()
```

Returns a temporary file name from the current working directory 
in 'z'.

##### unsetenv
```
z:bool = unsetenv(a:str)
```

Un-sets the environment variable name 'a' and returns true 
in 'z' if successful.

#### String procs

##### after
```
z:str = after(a:str,b:str)
```

Returns a sub-string of string 'a' in 'z' containing all characters
after substring 'b'. If 'b' is not located 'z' is the empty string.
```
> s="Kin folks said: Jed, git away from there!"
> w=after(s,"Jed")
> put(w)
, git away from there!
```

##### before
```
z:str = before(a:str,b:str)
```

Returns a sub-string of string 'a' in 'z' containing all characters 
before substring 'b'. If 'b' is not located 'z' is the empty string.
```
> s="Kin folks said: Jed, git away from there!"
> v=before(s,"Jed")
> put(v)
Kin folks said:
```

##### byte
```
z:num = byte(a:str[,b:int[,c:int]])
```

Returns a code point for 'a[b]'.
```
> a2k = "abcdefghijk"
> for i=1,#a2k do
>> put(substr(a2k,i,1),"=",byte(a2k,i))
>> end
a       =       97
b       =       98
c       =       99
d       =       100
e       =       101
f       =       102
g       =       103
h       =       104
i       =       105
j       =       106
k       =       107
```

##### char
```
z:str = char(a:num)
```

Returns character in 'z' for code point 'a'.
```
> for i=100,110 do
>> put(i,char(i))
>> end
100     d
101     e
102     f
103     g
104     h
105     i
106     j
107     k
108     l
109     m
110     n
```

##### contains
```
z:bool = contains(a:str,b:str)
```

Returns true in 'z' if string 'a' contains string 'b'.
```
> s="...was a man barely alive..."
> man = contains(s,"man")
> put(man)
true
```

##### containsany
```
z:bool = containsany(a:str,b:str)
```

Returns true in 'z' if string 'a' contains any characters 
from string 'b'.

##### count
```
z:num = count(a:str,b:str)
```

Returns the count in 'z' of strings 'b' in string 'a'.
```
> s="Around the rugged rock the ragged rascal ran"
> c=count(s,"the")
> put(c)
2
> r=count(s,"r")
> put(r)
6
```

##### decodebase64
```
z:str = decodebase64(a:str)
```

Returns a string into 'z' from a base64 encoded string in 'a'.
```
> e="VGhpcyBpcyBpdA=="
> d=decodebase64(e)
> put(d)
This is it
```

##### dump
```
z:str = dump(a:str)
```

Creates a formatted dump of 'a' in string 'z'.
```
> q="How many roads must a man walk down?"
> fd=dump(q)
> put(fd)
00000000  48 6f 77 20 6d 61 6e 79  20 72 6f 61 64 73 20 6d  |How many roads m|
00000010  75 73 74 20 61 20 6d 61  6e 20 77 61 6c 6b 20 64  |ust a man walk d|
00000020  6f 77 6e 3f                                       |own?|
```

##### encodebase64
```
z:str = encodebase64(a:str)
```

Creates a base64 string in 'z' from the string in 'a'.
```
> a="This is it" ;
> b=encodebase64(a) ;
> put(b)
VGhpcyBpcyBpdA==
```

##### find
```
s:num,e:num = find(a:str,p:str[,init[,plain]])
```

Returns the start 's' and end 'e' position of pattern 'p' in string 'a'.
```
> t="this is it folks"
> a,b = find(t,"it")     /* find 'it' */
> put(a,b)
9       10
> a,b = find(t,"i%a")    /* find 'i' followed by a letter */
> put(a,b)
3       4
```

Character pattern selection classes

| Class | Characters selected      |
|-------|--------------------------|
| .     | all characters           |
| %a    | letters [a-zA-Z]         |
| %b    | balanced string          |
| %c    | control characters       |
| %d    | digits [0-9]             |
| %l    | lower case letters [a-z] |
| %p    | punctuation characters   |
| %s    | space characters         |
| %u    | upper case letters [A-Z] |
| %w    | alphanumeric characters  |
| %x    | hexadecimal characters   |
| %z    | null or zero character   |

If the class selection code is in upper case, then the characters 
selected will be the complement of the characters selected in the 
table. For example: %C would select non-control characters instead 
of control characters.

Pattern selection modifiers

| Modifier | Function        |
|----------|-----------------|
| +        | one or more     |
| *        | zero or more    | 
| -        | zero or more    |
| ?        | zero or one     |
| [ ]      | any of set      |
| ( )      | capture         |
| ^        | not             |
| %        | modifier escape |

##### format
```
z:str = format(f:str,s:str)
```

Formats string 's' depending on format string 'f' returning 
formatted string in 'z'.

##### gsub
```
gsub()
```
    
##### hasprefix
```    
z:bool = hasprefix(a:str,b:str)
```

Returns true in 'z' if string 'a' has string 'b' as a prefix.
```
> q="Test string"
> a=hasprefix(q,"Te")
> put(a)
true
```

##### hassuffix
```
z:bool = hassuffix(a:str,b:str)
```

Returns true in 'z' if string 'a' has string 'b' as a prefix.
```
> q="Test string"
> b=hassuffix(q,"ing")
> put(b)
true
> c=hassuffix(q,"dog")
> put(c)
false
```

##### index
```
z:num = index(a:str,b:str)
```

Returns the first position in 'z' of string 'b' in string 'a'. 
If 'b' is not in 'a' then -1 is returned.
```
> q="Test string"
> p=index(q,"st")
> put(p)
2
> ix=index(q,"z")
> put(ix)
-1
```

##### indexany
```
z:num = indexany(a:str,b:str)
```

Returns the first position as 'z' of any of the characters in string 'b' 
that are present in string 'a'.
```
> st="The river bank was steep"
> pe=indexany(st,"aeiou")
> put(pe)
2
> pe=indexany(st,"rst")
> put(pe)
4
```

##### lastindex
```
z:num = lastindex(a:str,b:str)
```

Returns the last position in 'z' of string 'b' in string 'a'. 
If 'b' is not in 'a' then -1 is returned.

##### lastindexany
```
z:num = lastindexany(a:str,b:str)
```

Returns the last position as 'z' of any character in string 'b' 
that is in string 'a'.
```
> v="many other worlds exist"
> px=lastindex(v,"e")
> put(px)
18
```

##### len
```
z:int = len(a:str)
```

Returns the length of string 'a' in 'z'. This is an alternative to the 
'#' operator.
```
> s="a frog was lost"
> put(len(s))
15
> put(#s)
15
```

##### length
```
z:int = length(a:str)
```

Returns the length of string 'a' in 'z'. This is an alias for 
len(). See also the length operator # above.
```
> s="a frog was lost"
> put(length(s))
15
> put(len(s))
15
```

##### lower
```
z:str = lower(a:str)
```

Returns in 'z' the string 'a' converted to lower case.

##### match
```
z:str = match(a:str,p:str,i:num)
```

Returns capture from 's' using pattern 'p' into 'z' optionally 
starting at 'i'.

##### prskdtch
```
z:list = prskdtch(rx:str,a:str)
```

Returns a list of positions in 'z' of the regular expression 'rx' 
match results found while matching string 'a'.

##### prxchange
```
z:str = prxchange(rx:str,a:str,b:str)
```

Replaces regular expression 'rx' matches found in 'a' with value 
from 'b' returning the results in 'z'.

##### rep
```
z:str = rep(a:str,b:num)
```

Returns 'b' concatenated copies of string 'a' in string 'z'.
```
> s="Mouse "
> ss=rep(s,10)
> put(ss)
Mouse Mouse Mouse Mouse Mouse Mouse Mouse Mouse Mouse Mouse
```

##### replace
```
z:str = replace(a:str,b:str,c:str,d:num)
```

Replaces 'd' strings of 'b' found in 'a' with strings of 'c' 
and returns the result in 'z'. If 'd' is -1 all strings of 'b' 
are replaced.
```
> s="Mouse "
> ss=rep(s,10)
> put(ss)
Mouse Mouse Mouse Mouse Mouse Mouse Mouse Mouse Mouse Mouse
> ss=replace(ss,"use","te",-1)
> put(ss)
Mote Mote Mote Mote Mote Mote Mote Mote Mote Mote
```

##### reverse
```
z:str = reverse(a:str)
```

Returns in 'z' all characters of string 'a' in reverse order.
```
> s="pig dog"
> s=reverse(s)
> put(s)
god gip
```

##### scan
```
z:str = scan(a:str,b:str,c:int)
```

Returns a sub-string of string 'a' in 'z' which is the 'c'th 
substring from the left, delimited by one or more of the characters 
in string 'b'. If 'c' is a negative value the scan proceeds from 
right to left, and 'c' is counted from the right.
```
> s="Kin folks said: Jed, git away from there!"
> word5=scan(s," ?:,!.",5)
> put(word5)
git
```

##### scanall
```
z:list = scanall(a:str,b:str)
```

Returns a list of all sub-strings of string 'a' in 'z' which are 
delimited by one or more of the characters in string 'b'.
```
> s="Kin folks said: Jed, git away from there!"
> words=scanall(s," ,.?:!")
> dumpl(words)
(array): [
(str)1: "Kin"
(str)2: "folks"
(str)3: "said"
(str)4: "Jed"
(str)5: "git"
(str)6: "away"
(str)7: "from"
(str)8: "there"
]
> put(words[4])
Jed
```

##### sub
```
z:str = sub(a:str,b:int[,c:int])
```

Returns a sub-string of string 'a' in 'z' from position 'b' 
inclusive, **to position** 'c'.
```
> s="Kin folks said: Jed, git away from there!"
> v=sub(s,5)
> put(v)
folks said: Jed, git away from there!
> w=sub(s,5,8)
> put(w)
folk
```

##### substr
```
z:str = substr(a:str,b:int[,c:int])
```

Returns a sub-string of string 'a' in 'z' from position 'b' 
inclusive, **for length** 'c'.
```
> s="Kin folks said: Jed, git away from there!"
> w=substr(s,5,3)
> put(w)
fol
```

##### trim
```
z:str = trim(a:str,b:str)
```

Returns in 'z' the string 'a' with all characters from 'b' 
removed from the left and right side.

##### trimleft
```
z:str = trimleft(a:str,b:str)
```

Returns in 'z' the string 'a' with all characters from 'b' 
removed from the left side removed.

##### trimprefix
```
z:str = trimprefix(a:str,b:str)
```

Returns in 'z' the string 'a' with prefix string 'b' removed 
from the left side.

##### trimright
```
z:str = trimright(a:str,b:str)
```

Returns in 'z' the string 'a' with all characters from 'b' 
removed from the right side removed.

##### trimspace
```
z:str = trimspace(a:str)
```

Returns in 'z' the string 'a' with all white space characters 
on the left and right side removed.

##### trimsuffix
```
z:str = trimsuffix(a:str,b:str)
```

Returns in 'z' the string 'a' with suffix string 'b' removed 
from the right side.

##### title
```
z:str = title(a:str)
```

Returns in 'z' the string 'a' converted to title case.
```
> s="How many roads"
> put(title(s))
How Many Roads
```

##### upper
```
z:str = upper(a:str)
```

Returns in 'z' the string 'a' converted to upper case.
```
> s="How many roads"
> put(upper(s))
HOW MANY ROADS
```


#### QList procs

##### dumpl
```
dumpl(a:list)
```

Dumps the contents of list ‘a’ in a readable format to stdout. 
If list ‘a’ contains lists those lists are recursively dumped also.
```
> lst={one=1,str="dog",pi=3.14159,sub_ls={f="fox",c="bat"},array={1,2,3,4},carray={"one","two","three"}}
> dumpl(lst)
(strdict): {
(strdict)sub_ls: {
(str)f: "fox"
(str)c: "bat"
}
(strdict)array: {
(num)1: 1
(num)2: 2
(num)3: 3
(num)4: 4
}
(strdict)carray: {
(str)1: "one"
(str)2: "two"
(str)3: "three"
}
(num)one: 1
(str)str: "dog"
(num)pi: 3.14159
}
```

##### getn
```
z:num = getn(a:list)
```

Returns the number of elements in ‘z’ contained in list ‘a’.
```
> a={2,4,6,8}
> n=getn(a)
> put(n)
4
```

##### concat
```
z:str = concat(a:list [,b:str [,c:num [,d:num ] ] ] )
```

Returns in ‘z’ the concatenation of all the elements in list ‘a’. 
The string ‘b’ is an optional separator that is inserted between 
the list elements. The number ‘k’ is an optional starting position 
in the list and the number ‘d’ is an optional ending position in 
the list. 

##### insert
```
insert(a:list, [b:num,] c:*)
```

Inserts element ‘c’ into list ‘a’ at position ‘b’. If ‘b’ is 
omitted the insertion of ‘c’ is at the end of the list ‘a’.

##### marshal
```
z:str = marshal(a:list)
```

Marshal out the list 'a' into a string 'z' using JSON format.
```
> a = {st={a=22,b=3,c=45,d=6,e=92},b={dog=232,t="ccc"},four=4,five=5,real=3.141,str="This guy"}
> dumpl(a)
  (strdict): {
    (strdict)st: {
      (num)e: 92
      (num)a: 22
      (num)b: 3
      (num)c: 45
      (num)d: 6
    }
    (strdict)b: {
      (num)dog: 232
      (str)t: "ccc"
    }
    (num)four: 4
    (num)five: 5
    (num)real: 3.141
    (str)str: "This guy"
  }
> s = marshal(a)
> put(s)
{"st":{"a":22,"b":3,"c":45,"d":6,"e":92},"b":{"dog":232,"t":"ccc"},"four":4,"five":5,"real":3.141,"str":"This guy"}
> s = marshal(a,"  ")
> put(s)
{
  "st":{
    "e":92,
    "a":22,
    "b":3,
    "c":45,
    "d":6
  },
  "b":{
    "dog":232,
    "t":"ccc"
  },
  "four":4,
  "five":5,
  "real":3.141,
  "str":"This guy"
}
```

##### marshalxml
```
z:str = marshalxml(a:list)
```

Marshal out the list 'a' into a string 'z' using XML format.
```
> a = {st={a=22,b=3,c=45,d=6,e=92},b={dog=232,t="ccc"},four=4,five=5,real=3.141,str="This guy"}
> s = marshalxml(a)
> put(s)
<list><st><a>22</a><b>3</b><c>45</c><d>6</d><e>92</e></st><b><dog>232</dog><t>ccc</t></b><four>4</four><five>5</five><real>3.141</real><str>This guy</str></list>
```

##### unmarshal
```
a:list = unmarshal(z:str)
```

Unmarshal string 'z' in JSON format into a list 'a'.
```
> str = '{"st":{"a":22,"b":3,"c":45,"d":6,"e":92},"b":{"dog":232,"t":"ccc"},"four":4,"five":5,"real":3.141,"str":"This guy"}'
> lst = unmarshal(str)
> dumpl(lst)
(array): [
(array): [
(num)five: 5
(num)real: 3.141
(str)str: "This guy"
(strdict)st: {
(num)a: 22
(num)b: 3
(num)d: 6
(num)e: 92
(num)c: 45
}
(strdict)b: {
(str)t: "ccc"
(num)dog: 232
}
(num)four: 4
]
]
```

##### maxn
```
z:num = maxn(a:list)
```

Returns the largest position in list ‘a’ where an element can 
be inserted. If no items are in the list then the value returned 
is 0.
```
> a = {22,3,45,6,92,53,8,7,43,13,99}
> b = maxn(a)
> put(b)
11
```

##### erase
```
z = erase(a:list,b:num)
```

Removes the item at position ‘b’ from list ‘a’ and returns the 
removed item in ‘z’.
```
> qlist = {1,2,3,4,5,6,7}
> item = erase(qlist,4)
> put(item)
4
> dumpl(qlist)
(array): [
(num)1: 1
(num)2: 2
(num)3: 3
(num)4: 5
(num)5: 6
(num)6: 7
]
```

##### sort
```
sort(a:list [, b(c,d):bool ] )
```

Sorts the items in list ‘a’. The second parameter ‘b()’ is a call 
back comparison function which accepts two parameters. This is called
during sorting with two elements from list ‘a’. This function must
return a boolean value. When the returned value is true, this indicates 
the list element value ‘c’ is less than the list element value ‘d’. 
```
>  a = {22,3,45,6,92,53,8,7,43,13,99}
>  dumpl(a)
(array): [
(num)1: 22
(num)2: 3
(num)3: 45
(num)4: 6
(num)5: 92
(num)6: 53
(num)7: 8
(num)8: 7
(num)9: 43
(num)10: 13
(num)11: 99
]
> proc cproc(c,d)
>> if c < d then return true else return false end
>> end
> sort(a,cproc)
> dumpl(a)
(array): [
(num)1: 3
(num)2: 6
(num)3: 7
(num)4: 8
(num)5: 13
(num)6: 22
(num)7: 43
(num)8: 45
(num)9: 53
(num)10: 92
(num)11: 99
]
```

## Script examples

### Example 1 Comments Variables Procs
``` js
/*
  Script:   test.q
  Language: Q -- ubiquity scripting control language.
  Note: This is a block style comment.
*/

// This is a line comment

PGM = "test.q" ;   // PGM is a string variable
VER = "0.0.1" ;     // the sfnj-colon ';' is optional
// The built-in put() proc writes one or more comma
// separated values on the stdout file.
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
```

### Example 2 Input from file

`dingo@mach:~ > cat test.q`
``` js
put("File test Q program")       // test program
f = i.open("test.txt","r")
if i.type(f) == "file" then       /* if uses then and end */  
  for line in f:lines() do
    t = sub(line,26)
    if t != "" then
      put("::>",t)
    end
  end    /* of for */
end
f:close()
put("File test program ended.")
```

```
> q test.q
File test Q program
::>      202 Apr 19 12:15 c.json
::>      242 Apr 20 15:49 c.txt
::>        6 Sep  9  2016 Desktop/
::>       27 Jul 14  2016 gopath -> /dingo/gopath/
::>      175 Apr  4  2016 junk.txt
::>       11 Jul 18 17:05 pids.txt
::>      163 Aug  8  2016 test.js
::>      184 Aug 14 16:35 test.q
::>        0 Aug 14 16:53 test.txt
::>     4096 Sep  9  2016 tmp/
::>     2591 Jul 25  2016 vmstat.txt
File test program ended.
dingo@mach:~ >
```

### Example 3 Interactive

```
Microsoft Windows [Version 10.0.14393]
(c) 2016 Microsoft Corporation. All rights reserved.

C:\Users\user>q version
Program: q version: 2.0.0

C:\Users\user>q
q, Version - 2.0.0, Build Date - 16JUL2021
q interactive mode, enter exit() to exit, or help() for help.
> put(m.pi)
3.141592653589793
> dingo="I am here"
> put(dingo)
I am here
> p=m.pi
> q={1,2,3,4,5,6,7,8,9,10}
> put(q[5])
5
> for i=1,#amazing do
>> put(i,amazing[i])
>> end
1       1
2       2
3       3
4       4
5       dave
6       6
7       7
8       8
>
> all={
>> one="this",
>> two=55,
>> three=92.88}
>
> put(all.one)
this
> put(all.two)
55
> put(all.three)
92.88
> sum=all.two+all.three
> put(sum)
147.88
> all={
>> one="this",
>> two=55,
>> three=92.88}
>
> put(all.one)
this
> put(all.two)
55
> put(all.three)
92.88
> sum=all.two+all.three
> put(sum)
147.88
> all.dog="Lucy"
> put(all[4])
nil
> all[1]=45
> all[4]=96
> put(all[4])
96
> put(type(all))
list
> put(type(all[1]))
num
> for i=1,#all do
>> if type(all[i]) == "num" then
>> put(all[i])
>> end
>> end
45
96
>
> bye()
C:\Users\user>
```


## Return Codes

Return codes that may be issued by q.

| Return Code | Meaning                                  |
|-------------|------------------------------------------|
| 0           | OK (command completed successfully)      |
| 1           | Warnings (warning messages were logged)  |
| 2           | Errors (error messages were logged)      |
| 3           | Critical (command failed unexpectedly)   |


