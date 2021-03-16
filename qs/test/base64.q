/* 
  Script:   base64.q
  Language: q -- Q scripting control language.	
  Output:

*/ 
PGM = "filetest.q" ;   // PGM is a string variable
VER = "0.0.1" ;        // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

a="This is it" ;
put("Start text:",a) ;
b=encodebase64(a) ;
put("Text in base64:",b) ;
c=decodebase64(b) ;
put("End text:",c) ;

logi("Program:" || PGM || " ended.") ;
exit(0)