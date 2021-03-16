/*
  Script:   struct.q
  Language: q -- Q scripting control language.
  Purpose:  Language aggregate variable structure demonstration 
  
  Output:

  
*/
PGM = "struct.q" ;        // PGM is a string variable
VER = "0.0.1" ;           // version
// Test banner.
logi("Program:" || PGM || " version:" || VER) ;

// demonstrate multiline string - where the string contains various quotes
put("Multiline string")

t=`{
  "version": 1,
  "timeStamp": "2017-11-10T10:26:52.611993-05:00",
  "level": "info",
  "category": "system",
  "type": "threshold",
  "source": "ops",
  "message": "No thresholds tested for resource 'memory''memory' metric 'percentUsed', current value is 5.580623922100162 percent",
  "threshold": {
    "metricName": "percentUsed",
    "metricBaseType": "float",
    "metricType": "gauge",
    "metricUnit": "percent",
    "testOperator": "le"
  },
  "status": "pass",
  "results": [
    {
      "resourceType": "memory",
      "resourceId": "ptnode17",
      "checkStatus": "pass",
      "metricValue": "5.580623922100162"
    }
  ],
  "properties": {
    "checkName": "check-memory",
    "hostname": "ptnode17",
    "maxLevel": "pass"
  }
}
`
put(t)

put()
// create a complex Q list called q based on JSON emited by check command
put("Complex Q list")
q={version=1,
   timeStamp="2017-11-10T10:26:52.611993-05:00",  /* comments are OK */
   level="info",                                  // and line comments as well 
   category="system",
   type="threshold",
   source="ops",
   message="No thresholds tested for resource 'memory''memory' metric 'percentUsed', current value is 5.580623922100162 percent",
   threshold={
     metricName="percentUsed",
     metricBaseType="float",
     metricType="gauge",
     metricUnit="percent",
     testOperator="le"
   },
   status="pass",
   results={
    {
      resourceType="memory",
      resourceId="ptnode17",
      checkStatus="pass",
      metricValue="5.580623922100162"
    }
  },
  properties={
    checkName="check-memory",
    hostname="ptnode17",
    maxLevel="pass"
  }
}
// display the OA list q for verification
dumpl(q)

put()
// marshal out the list and print
put("Marshal the complex Q list and print")
z = marshal(q)
put(z)

put()
// marshal out the list with indenting format and print
put("Marshal the complex Q list with indenting format and print")
z = marshal(q,"  ")
put(z)

put()
// marshal out the list with non-white space indenting format and print
put("Marshal the complex Q list with non-white space indenting format and print")
z = marshal(q,"|  ")
put(z)

// thats all folks
logi("Program:" || PGM || " ended.") ;
exit(0)