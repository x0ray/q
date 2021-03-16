/* 
  Script:   deploy.q
  Language: q -- Q scripting control language.	
  Output:

*/
PGM = "deploy.q" ;      // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

// deploymentinfo test
d = deploymentinfo()
put(d)
for k,v in pairs(d) do
  put(k,v)
end

put("Program:" || PGM || " ended.") ;