/*  
  Script:   consulall.q
  Language: q -- ops-agent Q scripting control language.	
  Purpose:  Consul key load, access, delete test
  Date:     31Aug2017
  Output:
		Q script Consul load, access, delete test.

		Base key:       ops/emitest
		Consul keys loaded.
		Number of Consul keys loaded:   4
		  Key:  blue    Val:    cat
		  Key:  red     Val:    dog
		  Key:  green   Val:    frog
		  Key:  yellow  Val:    bird
		Consul keys at: ops/emitest     deleted.
		Q script Consul load, access, delete test ended.
*/

base = "ops/emitest"  // select base thats NOT already used !!

put("Q script Consul load, access, delete test.\n")
put("Base key:",base)

/* make some new KV pairs in list q */
q = {red="dog",blue="cat",green="frog",yellow="bird"}
// load new KV pairs to base location in Consul 
err = consulsetkeys(base,q)
// check that the load it worked
if err == nil then
  put("Consul keys loaded.")
else
  put("Consul keys not loaded.")
  put(err)
  exit(1)
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

/* delete the the Consul keys */
err = consuldeletekeys(base)
if err == nil then
  put("Consul keys at:",base,"deleted.")
else
  put("Consul keys at:",base,"NOT deleted.")
  put(err)
  exit(1)  
end

put("Q script Consul load, access, delete test ended.")
exit(0)
/* end test */