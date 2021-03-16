/* 
  Script:   json.q
  Language: q -- ops-agent Q scripting control language.	
  Output:

*/
PGM = "json.q" ;        // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

// json marshal test
a={1,2,3,"this",{"roger",5,6,{1,2},"bill",66},"hope",98}
marshal(a)

ba={one=1,two=2,three=3,This="this",struct1={Roger="roger",five=5,six=6,struct2={one2=1,two2=2},Billy="bill",Route=66},No="hope",End=98}
marshal(ba)
	
c={outer="This",name={1,2,3,4,5,6},name2={1.2,1.4,1.5,1.7,1.9}}
marshal(c)

d={ok=false,outer="This",array={1,2,3,4,5,6},{one=1,two=2,three=3,This="this",struct1={Roger="roger",five=5,six=6,struct2={one2=1,two2=2},Billy="bill",Route=66},No="hope",End=98}}
marshal(d)


put("Program:" || PGM || " ended.") ;