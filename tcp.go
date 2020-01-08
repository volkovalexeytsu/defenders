package main

import (
	"net"
	"fmt"
	"bufio"
	"strings" // only needed below for sample processing
	"strconv"
	"log"
	"os"
	"math/rand"
	"time"
	"flag"
)

//Получение стартового ключа
func get_session_key() string {
  var result string 
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  for i:=1; i < 11; i++ {
      var ran int = r.Intn(9) + 1
      var stran string = strconv.Itoa(ran)
      result += stran
  }
  return result
}

//Получение Хеша
func get_hash_str() string {
 var li string
 r := rand.New(rand.NewSource(time.Now().UnixNano()))
 for i:=0; i < 5; i++ {
   var ran int = r.Intn(6) + 1
   var stran string = strconv.Itoa(ran)
   li += stran
 }
 return li
}


//Получение следующего ключа
func next_session_key(hash, key string) string {
  if hash == "" {
    fmt.Println("Hash code is empty")
    }
  var result int
  for idx:=0; idx < len(hash); idx++ {
    var s string = string(hash[idx])
    i, err := strconv.Atoi(s)
	  if err != nil {
		  log.Fatal(err)
	  }
    result += calc_hash(key,i)
  }
  var ult string = strconv.Itoa(result)
  var res string = "0000000000" + ult
  var s1,s2 string
  for i:=(len(res)-1); i>=(len(res)-10); i-- {
	  s1 += string(res[i])
	}
	for i:=9; i>=0; i-- {
		s2 += string(s1[i])
	}
  return s2
}

//Вычисление следующего ключа на основе Хеша
func calc_hash(key string, val int) int{
  if val == 1 {
    var res1 string = string(key[0:5])
    k1, err :=strconv.Atoi(res1)
    if err != nil {
      log.Fatal(err)
    }
    return k1%97
  } else if val == 2 {
    var res2 string
    for i := 9; i >= 0; i-- {
		  var s2 string = string(key[i])
		  res2 += s2
	  }
    k2, err :=strconv.Atoi(res2)
    if err != nil {
      log.Fatal(err)
    }
    return k2
  } else if val == 3{
    var k3 string = key[5:10] + key[:5]
    k30, err :=strconv.Atoi(k3)
    if err != nil {
      log.Fatal(err)
    }
    return k30
  } else if val == 4{
    var k4 int
	  for i:=1; i<9; i++{
		  var s4 string = string(key[i])
      res4, err := strconv.Atoi(s4)
	    if err != nil {
		    log.Fatal(err)
	 	  }
		  k4 += (res4+41)
	  }
    return k4
  } else if val == 5{
    var k5 int
    var res5 byte
    var s5 byte
    for i:=0; i<10; i++ {
	  	s5 = key[i]
	  	res5 = s5 ^ 43
	  	k5 += int(res5)
  	}
    return k5
  } else {
	  k6, err := strconv.Atoi(key[0:9])
	  if err != nil {
		  log.Fatal(err)
	  }
    var k10 int64 = int64(k6)
	  k7, err := strconv.Atoi(key[9:10])
	  if err != nil {
	  	log.Fatal(err)
  	}
	  var k20 int64 = int64(k7)
  	return int(k10*10 + k20)
  }
}


func serv(conn *net.Conn) {

	var hash_from_client string
	var key_from_client string
	var next_key string
	
// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(*conn).ReadString('\n')
		hash_from_client = string(message[0:5])
		key_from_client = string(message[5:15])
		next_key = next_session_key(hash_from_client, key_from_client)
		// output message received
		fmt.Print("Key received: ", string(message[5:15]), " Next key: ", next_key, " Message Received: ", string(message[15:]))
		// sample process for string received
		newmessage := strings.ToUpper(message[15:])
		// send new string back to client
		(*conn).Write([]byte(next_key + newmessage + "\n"))
	}
}




func main() {

	port := flag.String("port", ":8081", "port")
    IP := flag.String("ip:port", "", "ip/port")
    n := flag.Int("n", 100, "num of connections")
	flag.Parse()
	if *IP == "" {
		fmt.Println("Launching server...")
		// listen on all interfaces
		ln, _ := net.Listen("tcp", *port)
		var kolvo int = 1
		// accept connection on port
		for {
			conn, _ := ln.Accept()
			if kolvo <= *n {
				kolvo++
				go serv(&conn)
			} else{conn.Close()} 	
		}
	} else {
		// connect to this socket
		conn, _ := net.Dial("tcp", *IP)
		//Формируем хеш
		var hash_to_server string = get_hash_str()
		var key_to_server string
		var key_to_check string
		var key_from_server string
		for { 
		
			// read in input from stdin
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Text to send: ")
			text, _ := reader.ReadString('\n')
			//Формирование ключа
			key_to_server = get_session_key()
			key_to_check = next_session_key(hash_to_server, key_to_server)
			// send to socket
			fmt.Fprintf(conn, hash_to_server + key_to_server + text + "\n")
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')
			key_from_server = string(message[0:10])
			if key_from_server != key_to_check {break}
			fmt.Print("Key from server: ", key_from_server, " Message from server: "+message[10:])
			
		}
		
	} 
}