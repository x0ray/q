/*
  Script:   loop.q
  Language: q -- ops-agent Q scripting control language.
  Purpose:  Measure the time of counting up to ten million by one.
  
  Output:
						2018-02-15T11:12:33.898226-05:00 INFO  [qslibbase.go:441] [Q] - Program:loop.q version:0.0.1
						Start: 02/15/2018 11:12:33
						Elapsed time: 2.81
						Finish: 02/15/2018 11:12:36
						2018-02-15T11:12:36.795748-05:00 INFO  [qslibbase.go:441] [Q] - Program:loop.q ended.

*/
PGM = "loop.q" ;        // PGM is a string variable
VER = "0.0.1" ;          // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

put(date("Start: %m/%d/%Y %X"))
dcl x = clock()
dcl s = 0
for i=1,10000000 do  // perform tem million times
  s = s + 1
end
e = clock()
put(format("Elapsed time: %.2f",e-x))
put(date("Finish: %m/%d/%Y %X"))

// thats all folks
logi("Program:" || PGM || " ended.") ;
exit(0)   // exit the script with code zero
