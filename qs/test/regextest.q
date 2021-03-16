/*
  Script:   regextest.q
  Language: q -- ops-agent Q scripting control language.
  Output:
2017-09-28T13:51:32.539038-04:00 INFO  [qslibbase.go:432] [Q] - Program:regextest.q version:0.0.1
Test string:    This is the story of a
man named Jed whose kin folks
were distinct wierdos. Most
folks came to this conclusion
early!
Test string length:     117
Reg.Exp:        (is)|(os)
prxmatch length:        7
prxmatch list:
1       2
2       5
3       39
4       59
5       72
6       77
7       97
2017-09-28T13:51:32.543974-04:00 INFO  [qslibbase.go:432] [Q] - Program:regextest.q ended.
*/
PGM = "regextest.q" ;   // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;


a=`This is the story of a
man named Jed whose kin folks
were distinct wierdos. Most
folks came to this conclusion
early!`

put("Test string:",a)   // length of words list
put("Test string length:",#a)   // length of test string

rxp="(is)|(os)"
put("Reg.Exp:",rxp) 

ls = prxmatch(rxp,a)
put("prxmatch length:",#ls) 
put("prxmatch list:") 
for i=1,#ls do
  put(i,ls[i])
end

logi("Program:" || PGM || " ended.") ;
