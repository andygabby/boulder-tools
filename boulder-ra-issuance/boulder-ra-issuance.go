//This is a nagios/sensu style check against the boulder-ra lastIssuance stat
//It determines whether the lastIssuance is over the defined time threshold. 
package main

import (
  "os"
  "flag"
  "time"
  "fmt"
  "net/http"
  "encoding/json"
)

// create a struct to get the json value of lastIssuance.
type Stats struct {
  LastIssuance int64 `json:"lastIssuance"`
}

func main () {

  //Flags for url to pull metrics and threshold of seconds to raise alert
  url := flag.String("i", "http://0.0.0.0:4444/metrics/vars", "debug stats url to parse")
  threshold := flag.Int64("t", 120, "alert threshold in seconds")
  flag.Parse()

  //Get json from api
  statsData := new(Stats)
  getJson(*url, statsData)

  //Get now time to compare against last issuance in unix epoch
  now := time.Now()
  secs := now.Unix()

  //Difference between now and last issuance in seconds
  timeDiff := secs - statsData.LastIssuance

  //Alert if 0 or  over threshold, else we are ok
  //Upon restart of boulder-ra the stat is 0 until first issuance.
  //It is also 0 if the stat cannot be retrieved (bad url or firewall rule blocking)
  if statsData.LastIssuance == 0 {
    fmt.Printf("Critical: lastIssuance value is 0 or cannot retrieve stat\n")
    os.Exit(2)
  } else if timeDiff > *threshold  {
    fmt.Printf("Critical: %d seconds since last issuance\n", timeDiff)
    os.Exit(2)
  } else {
    fmt.Printf("OK: %d seconds since last issuance\n", timeDiff)
    os.Exit(0)
  }
}

//Return the json string from the metrics url directly to the Stats struct
func getJson(url string, target interface{}) error {
  r, err := http.Get(url)
  if err != nil {
    return err
  }
  defer r.Body.Close()

  return json.NewDecoder(r.Body).Decode(target)
}
