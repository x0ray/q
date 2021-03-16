/*
  Script:   test.q
  Language: q -- Q scripting control language.	
  Note: This is a block style comment.
  
  Output:
  
  
*/

// This is a line comment  

PGM = "test.q" ;    // PGM is a string variable
VER = "0.0.1" ;     // the semi-colon ';' is optional
// The builtin put() proc writes one or more comma 
// seperated values on the stdout file.
// The || operator concatenates strings.
logi("Program:" || PGM || " version:" || VER) ;

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

// a variable can be a proc
a = 59
b = 68.4521
add2 = blabla
f = add2(a,b)
put("a add2 b:",f)

// calculate and print elapsed time in seconds as float
et = clock()
put("Elapsed time:",et-st,"(s)")

put("Program:" || PGM || " ended.") ;
