/*
  Script:   consul.q
  Language: q -- ops-agent Q scripting control language.
  Purpose:  Consul key access test
  Date:     31Aug2017
  Output:
*/
base = "/ops"

put("Q script Consul access test.\n")
put("Base key:",base)
q = consulgetkeys(base)

if q != nil then
  put("  Keys:",keys(q))

  // print out all the keys
  for key, val in pairs(q) do
    put(" Key:",key," Val:",val)
  end
else
  put("No Consul keys found for: ",base)
  err = consulerror(base)
  put(err)
end

put("Consul test ended.")
