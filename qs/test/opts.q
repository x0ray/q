/*
  Script:   opts.q
  Language: q -- Q scripting control language.
  Purpose:  Command and options scanning
  
  Input:    q run opts.oa -- -my -dog spot -and "erin brockovitch" 
  Output:

*/
PGM = "opts.q" ;        // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

// get command line parameters after --
put("Command line args as a string")
a = argstr()
put(a)

// put parameters into a list by position
put("Parameters in a list by position")
al = arglist(a)
put(al)
for key,val in pairs(al) do
  put(key,val)
end

// put parameters into a list by option
put("Parameters in a list by option")
kl = argopts(a)
for k,v in pairs(kl) do
  put(k,v)
end

// thats all folks
logi("Program:" || PGM || " ended.") ;
exit(0)