/*
  Script:   mathtest.q
  Language: q -- Q scripting control language.	
  Output:
    
*/

PGM = "mathtest.q" ;    // PGM is a string variable
VER = "0.0.1" ;         // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;
// simple math
put("simple math:")
a = 44 ;
b = 6 ;
put("a:",a,"b:",b)
c = a - b ;
d = a / b ;
e = a * b ;
f = a + b ;
put("a-b=c:",c,"a/b=d:",d,"a*b=e:",e,"a+b=f:",f) ;

// builtin constants
put("math constants:")
put("pi", pi)
put("e", e)
put("phi", phi)
put("sqrt2", sqrt2)
put("sqrte", sqrte)
put("huge", huge)
put("small", small)

// various built in math library procs
put("math builtin procs:")
put("abs(-77):",abs(-77)) 
put("acos(0.3):",acos(0.3)) 
put("asin(0.3):",asin(0.3)) 
put("atan(3.0):",atan(3.0)) 
put("atan2(3.0,6):",atan2(3.0,6)) 
put("ceil(3.5283):",ceil(3.5283)) 
put("cos(0.3):",cos(0.3)) 
put("cosh(0.3):",cosh(0.3)) 
put("deg(pi/2):",deg(pi/2)) 
put("exp(6.05):",exp(6.05)) 
put("fact(6):",fact(6)) 
put("fib(11):",fib(11)) 
put("floor(7.6892):",floor(7.6892)) 
put("fmod(78.345,3):",fmod(78.345,3)) 
put("frexp(78.345):",frexp(78.345)) 
put("ldexp(1.345,4):",ldexp(1.345,4)) 
put("log(300):",log(300)) 
put("log10(300):",log10(300)) 
put("max(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",max(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("mean(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",mean(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("median(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",median(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("min(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",min(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("mod(78.345,3):",mod(78.345,3)) 
put("mode(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",mode(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("modf(78.345,3):",modf(78.345)) 
put("pow(1.345,4):",pow(1.345,4)) 
put("rad(90):",rad(90)) 
put("randomseed(12345):",randomseed(12345)) 
put("random(),random(),random(),random():",random(),random(),random(),random()) 
put("range(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",range(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("rms(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",rms(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("sin(0.3):",sin(0.3)) 
put("sinh(0.3):",sinh(0.3)) 
put("sqrt(144):",sqrt(144)) 
put("stddev(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",stddev(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("sum(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",sum(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 
put("tan(0.3):",tan(0.3)) 
put("tanh(0.3):",tanh(0.3)) 
put("variance(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23):",variance(30,40,9,4,56,76,33,33,24,3,543,66.08,3,1.23)) 

logi("Program:" || PGM || " ended.") ;
