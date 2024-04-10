package main

import (
  "fmt"
  "encoding/json"
  "os"
  "net"
  "strconv"
)

type Config struct {
  Port int 
  Fields map[string]string
}
func LoadConfigs() (*Config, error){
  data, err := os.ReadFile("config.json")
  if err != nil {
    return nil, err
  }

  var config Config
  err = json.Unmarshal(data, &config)
  if err != nil {
    return nil, err
  }

  return &config, nil
  
}


func main() {
  //Load configurations
  config, err := LoadConfigs()
  if err != nil {
    fmt.Printf("Error loading configs: %s\n", err)
    return
  }
  
  //Open Listener
  listener, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
  if err != nil {
    fmt.Printf("Error setting up lister: %s\n", err)
    os.Exit(1)
  }
  defer listener.Close()
  fmt.Printf("HL7 Validator started and listening on port %n\n", config.Port)


  //Validate messages based on configurations

} 
