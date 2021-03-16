/*
  Script:   stringtest.q
  Language: q -- Q scripting control language.	
  Output:
  

*/

PGM = "stringtest.q" ;   // PGM is a string variable - This is a line comment
VER = "0.0.1" ;          // the semi-colon ';' is optional
logi("Program:" || PGM || " version:" || VER) ;

// Test banner.
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

logi("Program:" || PGM || " ended.") ;