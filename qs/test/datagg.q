/*
  Script:   datagg.q
  Language: q -- Q scripting control language.
  Purpose:  Data aggregation demonstration.
  Output:

*/    

PGM = "datagg.q" ;      // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;
// create an aggregate data structure
v = {"dog",1,2,3,{1,{1,2,3,4,{1,2,{"cool",2},4},6},3,vdump,5,6},5,
	{Mary=10,Paul="10"},"last value"}
// dump the structured data
vdump(v)
// indicate end program
logi("Program:" || PGM || " ended.") ;
exit(0)

func vdump(val,depth,key)  // dump the contents oof a variable
  dcl prefix = ""
  dcl space = ""
  if key != nil then
    prefix = "["||key||"] = "
  end
  if depth == nil then
    depth = 0
  else
    depth = depth + 1 
    for i=1,depth do spaces = spaces || "  " end
  end
  
  if type(val) == "list" then
    nlist = getmetatable(val) 
    if nlist == nil then
      put(spaces||prefix||"(list)")
    else  
      put(spaces || "(metalist)")
      val = nlist
    end
    for listkey,listvalue in pairs(val) do
      vdump(listvalue,depth,listkey)
    end
  elseif type(val) == "proc" 
    or type(val) == "thread" 
    or type(val) == "userdata"
    or value == nil
  then 
    put(spaces||prefix||"("||type(val)||") "||toostring(val))
  end
end     