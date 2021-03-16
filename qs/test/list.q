/*
  Script:   list.q
  Language: q -- ops-agent Q scripting control language.
  Purpose:  Demonstration of Q scripting list aggregate data structure
  Output:
*/
PGM = "list.q" ;         // PGM is a string variable
VER = "0.0.1" ;          // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

put("An Q list as an array")
a={1,2,3}
put(a)
dumpl(a)

// demo proc
proc adder(a,b)
  return a+b
end  
// reference to proc
p=adder
// test the adder proc via the reference
put(p(4,7))
// create an Q list with the proc reference included
d={ok=false,outer="Outer",full=nil,array={1,2,3,4,5,6},q=p,{one=1,two=2,
  three=3,This="this",struct1={Roger="roger",five=5,six=6,
  struct2={one2=1,two2=2},Billy="bill",Route=66},
  No="hope",End=98}}
// display formatted dump of the list d to stdout
dumpl(d) 
// marshal the list d to a JSON string in q
q=marshal(d)
put(q)  // display the JSON

// add JSON to a string variable j
j='{"resourceType":"taskmanager","managerName":"ray","xmPgmVer":"0.0.86-PeachPuff","version":"1","nothing":null,"initImportCount":0,"tasks":[{"version":1,"taskName":"OpsAgentActivity","happy":true,"runType":"periodic_aligned","errorAction":"cancel","publisherType":"none"}],"osType":"linux"}'
put(j)   // display the list string
// read the JSON string j into an OA list called u
u=unmarshal(j)
// dump out the OA list u
dumpl(u)

// thats it folks
logi("Program:" || PGM || " ended.") ;
