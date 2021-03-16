/*  
  Script:   consuload.q
  Language: q -- ops-agent Q scripting control language.	
  Purpose:  Consul key load, access, delete test
  Date:     31Aug2017
  Output:
						Q script Consul load test.

						Base key:       ops/emitest
						Consul keys loaded.
						Number of Consul keys loaded:   9
						  Key:  black   Val:    pig
						  Key:  pink    Val:    mouse
						  Key:  yellow  Val:    bird
						  Key:  orange  Val:    fish
						  Key:  deep/blue/steel Val:    ardvark
						  Key:  green   Val:    frog
						  Key:  red     Val:    dog
						  Key:  blue    Val:    cat
						  Key:  deep/blue       Val:    whale
						Q script Consul load test ended.
*/
base = "ops/emitest"

put("Q script Consul load test.\n")
put("Base key:",base)

/* make some new KV pairs in list q */
q = {black="pig",orange="fish",red="dog",blue="cat",green="frog",yellow="bird"}
q["deep/blue"]="whale"
q["deep/blue/steel"]="ardvark"
// load new KV pairs to base location in Consul 
err = consulsetkeys(base,q)
// check that the load it worked
if err == nil then
  put("Consul keys loaded.")
else
  put("Consul keys not loaded.")
  put(err)
end

/* access new Consul keys */
p = consulgetkeys(base)
// check that consul KV pairs were loaded into list p
if p != nil then
  put("Number of Consul keys loaded:",keys(q))
  // print out all the keys
  for key, val in pairs(p) do
    put("  Key:",key,"Val:",val)
  end
else
  put("No Consul keys found for: ",base)
  err = consulcheckerror(base)
  put(err)
  exit(1)
end
put("Q script Consul load test ended.")