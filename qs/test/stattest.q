/* 
  Script:   stattest.q
  Language: q -- Q scripting control language.	
  Output:
  
*/
PGM = "stattest.q" ;    // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

fn = "test.q"
sl = stat(fn)
if sl != nil then 
	put("File name:",fn, sl)
	put("  Name........:",sl.name)
	put("  Mode........:",sl.mode)
	put("  Type........:",sl.type)
	put("  Size........:",sl.size)
	put("  Blocks......:",sl.blocks)
	put("  Blocksize...:",sl.blksize)
	put("  Access time.:",sl.atim)
	put("  Modify time.:",sl.mtim)
	put("  Change time.:",sl.ctim)
else
	loge("File name ",fn," not found.") 
end 

put("Program:" || PGM || " ended.") ;