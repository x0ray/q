/* 
  Script:   filetest.q
  Language: q -- ops-agent Q scripting control language.	
  Output:
david@ptnode17:~ > emi-oa filetest.q
2017-08-22T13:08:31.256748-04:00 INFO  [qslibbase.go:255] [Q] - Program:filetest.q version:0.0.1
::>      202 Apr 19 12:15 c.json
::>      242 Apr 20 15:49 c.txt
::>        6 Sep  9  2016 Desktop/
::>       27 Jul 14  2016 gopath -> /dept/fitforit/david/gopath/
::>      175 Apr  4  2016 junk.txt
::>       11 Jul 18 17:05 pids.txt
::>      163 Aug  8  2016 test.js
::>      184 Aug 14 16:35 test.oa
::>        0 Aug 14 16:53 test.txt
::>     4096 Sep  9  2016 tmp/
::>     2591 Jul 25  2016 vmstat.txt
Program:filetest.q ended.

*/ 
PGM = "filetest.q" ;   // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

f = i.open("test.txt","r")
if type(f) == "data" then
  for line in f:lines() do
    t = sub(line,26)
    if t != "" then
      put("::>",t)
    end
  end
else
  put("Type not file, is: ",type(f))	
end
    
f:close()

logi("Program:" || PGM || " ended.") ;