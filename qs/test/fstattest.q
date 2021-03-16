/* 
  Script:   fstattest.q
  Language: q -- Q scripting control language.	
  Output:
  
*/
PGM = "fstattest.q" ;   // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;
path = "/home"
sl = statfs(path)
put("File system path:",path, sl)
put("  Name.............:",sl.name)
put("  Type.............:",sl.type)
put("  BSize............:",sl.bsize)
put("  Blocks...........:",sl.blocks)
put("  Bfree............:",sl.bfree)
put("  Bavail...........:",sl.bavail)
put("  Files............:",sl.files)
put("  Ffree............:",sl.ffree)
put("  Fsid.............:",sl.fsid1,sl.fsid2)
put("  Namelen..........:",sl.namelen)
put("  Frsize...........:",sl.frsize)
put("  Files............:",sl.files)
put("  Typename.........:",sl.typename)
put("  Usedblocks.......:",sl.usedblocks)
put("  Totalbytes.......:",sl.totalbytes)
put("  Freebytes........:",sl.freebytes)
put("  Usedbytes........:",sl.usedbytes)
put("  Percentfreebytes.:",sl.percentfreebytes)
put("  Percentfreenodes.:",sl.percentfreenodes)
put("  Percentusedbytes.:",sl.percentusedbytes)
put("  Percentusednodes.:",sl.percentusednodes)

logi("Program:" || PGM || " ended.") ;