/*
  Script:   mailtest.q
  Language: q -- Q scripting control language.	
  Output:
*/

PGM = "mailtest.q" ;    // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;
// email test
message = "This is a test message using the Q language !"

setmailhost("nzkiwi1g@gmail.com","********","gmail.com")

sendmail("nzkiwi1g@gmail.com","nzkiwi1g@gmail.com",message)

logi("Program:" || PGM || " ended.") ;
