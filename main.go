package main

import (
  "fmt"
  "encoding/json"
  "os"
  "net"
  "log"
  "strconv"
)

type Config struct {
  Port int 
  SendAcks bool
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

type Message struct {
  sender net.Conn
  body []byte
}

type Server struct {
  listenPort string
  sendAcks bool
  listener net.Listener
  quitChan chan struct{}
  msgChan chan Message 
}

func NewServer(port string, ack bool) *Server {
  return &Server{
    listenPort: port,
    sendAcks: ack,
    quitChan: make(chan struct{}),
    msgChan: make(chan Message, 10),
  }
}

func (s *Server) Start () error {
  listener, err := net.Listen("tcp", s.listenPort)
  if err != nil {
    return err
  }
  defer listener.Close()
  s.listener = listener

  go s.acceptLoop()

  <-s.quitChan
  close(s.msgChan)

  return nil
}

func (s *Server) acceptLoop() {
  for {
    conn, err := s.listener.Accept()
    if err != nil {
      fmt.Printf("Error accepting the connection %s \n", err)
      continue
    }
    fmt.Printf("New Connection made. %s\n", conn.RemoteAddr())
    go s.readLoop(conn)
  }
}

func (s *Server) readLoop(conn net.Conn) {
  defer conn.Close()
  buf := make([]byte, 2048)
  for {
    n, err := conn.Read(buf)
    if err != nil {
      fmt.Printf("Error reading from connection to buffer %s", err)
      continue
    }
    s.msgChan <- Message {
      sender: conn,
      body:buf[:n],
    }

    //Ack feature here? 
    if(s.sendAcks){
      conn.Write([]byte("ack"))
    }

  }
}


func main() {
  //Load configurations
  config, err := LoadConfigs()
  if err != nil {
    fmt.Printf("Error loading configs: %s\n", err)
    return
  }
  
  //Open Listener
  server := NewServer(":" + strconv.Itoa(config.Port),config.SendAcks)
  log.Fatal(server.Start())

  fmt.Printf("HL7 Validator started and listening on port %s\n", strconv.Itoa(config.Port))

  go func(){
    for msg := range server.msgChan{
      fmt.Printf("Received message %s \n \n ", string(msg.body))
    }
  }()


  //Validate messages based on configurations

} 
