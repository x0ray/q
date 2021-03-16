/*
  Script:   scantest.q
  Language: q -- Q scripting control language.
  Purpose:  String scanning test
  Output:
  

*/
PGM = "scantest.q" ;    // PGM is a string variable
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

word7 = scan(a,". \n!",7)
put("word 7",word7)

wordminus4 = scan(a,". \n!",-4)
put("word minus 4",wordminus4)

all = scanall(a,". \n!")   // scan for all blank delimited words
put("Number of blank delimited words:",#all)   // length of words list
for i=1,#all do
  put(i,all[i])
end

all = scanall(a,"aeiou \n.!")   // scan for all blank delimited words
put("Number of vowel delimited word parts:",#all)   // length of words list
for i=1,#all do
  put(i,all[i])
end

logi("Program:" || PGM || " ended.") ;
