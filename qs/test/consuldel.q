/*  
  Script:   consuldel.q
  Language: q -- Q scripting control language.	
  Purpose:  Consul key load, access, delete test
  Output:
		Q script Consul load, access, delete test.

		Base key:       ops/emitest
		Consul keys at: ops/emitest     deleted.
		Q script Consul load, access, delete test ended.
*/

base = "ops/emitest"  // select base thats NOT already used !!

put("Q script Consul delete test.\n")
put("Base key:",base)

/* delete the the Consul keys */
err = consuldeletekeys(base)
if err == nil then
  put("Consul keys at:",base,"deleted.")
else
  put("Consul keys at:",base,"NOT deleted.")
  put(err)
  exit(1)  
end

put("Q script Consul delete test ended.")
exit(0)
/* end test */