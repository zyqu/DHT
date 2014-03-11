package main

import (
    "os"
    "bufio"
    "flag"
    "fmt"
    "log"
    "math/rand"
    "net"
    //"net/http"
    //"net/rpc"
    "time"
    "strings"
	//"io/ioutil"
    "strconv"
    "io/ioutil"
	"encoding/json"
)

import (
    "kademlia"
)


func main() {
    // By default, Go seeds its RNG with 1. This would cause every program to
    // generate the same sequence of IDs.
    rand.Seed(time.Now().UnixNano())
	

    // Get the bind and connect connection strings from command-line arguments.
    flag.Parse()
    args := flag.Args()
    if len(args) != 2 {
        log.Fatal("Must be invoked with exactly two arguments!\n")
    }
    listenStr := args[0]
    firstPeerStr := args[1]

    fmt.Printf("kademlia starting up!\n")
    kademClient := kademlia.NewKademlia()

    iptokens:=strings.Split(firstPeerStr, ":")
	ipaddr, err:=net.ResolveIPAddr("ip", iptokens[0])
	if err==nil{
    	kademClient.Host=net.ParseIP(ipaddr.String())
    	kademClient.Port=kademlia.Str2Port(iptokens[1])
	}else{
		fmt.Println("Cannot resolve client ip: ", iptokens[0])
		return
	}
	server:=strings.Split(listenStr, ":")
	serverIP, err:=net.ResolveIPAddr("ip", server[0])
	if err==nil{
		kademlia.StartServ(kademClient, serverIP.String()+":"+server[1])
	}else{
		fmt.Println("Cannot resolve server ip")
		return
	}
	//fmt.Println(kademlia.Hashcode("Kademlia").AsString())

    //kademServer := kademlia.NewKademlia()
    //kademlia.StartServ(kademServer,listenStr)
    //kademlia.StartServ(kademClient,firstPeerStr)
/*
    rpc.Register(kadem)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", listenStr)
    if err != nil {
        log.Fatal("Listen: ", err)
    }

    // Serve forever.
    go http.Serve(l, nil)
*/
    // Confirm our server is up with a PING request and then exit.
    // Your code should loop forever, reading instructions from stdin and
    // printing their results to stdout. See README.txt for more details.
    /*
    client, err := rpc.DialHTTP("tcp", firstPeerStr)
    if err != nil {
        log.Fatal("DialHTTP: ", err)
    }
    ping := new(kademlia.Ping)
    //fmt.Println(ping)
    ping.Sender.Host=net.IPv4(224, 0, 0, 1) 
    ping.MsgID = kademlia.NewRandomID()
    var pong kademlia.Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
        log.Fatal("Call: ", err)
    }

    

    log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
    log.Printf("pong msgID: %s\n", pong.MsgID.AsString())*/

    reader := bufio.NewReader(os.Stdin)
    //loop
    for{
        line, err := reader.ReadString ('\n')
        if err != nil{
            break
        }
        command := strings.Replace(line,"\n","",-1)
        tokens := strings.Fields(command)


        switch tokens[0] {

        case "batchtest":
            if len(tokens)!=3{
                fmt.Println("fetchurl takes 2 mode, filename mode")
                break
            }
            mode, err := strconv.Atoi(tokens[2])
            filename :=tokens[1]
            if err!=nil{
                break
            }

            

            if mode==1{
                fmt.Println("Multi-user Mode")
				//save mode1's log into log1.txt
	            logf, _ := os.Create("log1.txt")
	            log.SetOutput(logf)
	            defer logf.Close()
				J, err:=ioutil.ReadFile("./wikiindex/"+filename)
				if err!=nil{
					panic(err)
				}
				fmt.Println("Indicate start and end positions: ")
				line, err := reader.ReadString ('\n')
				if err!=nil{
					fmt.Println(err)
					break
				}
				command:=strings.Replace(line, "\n", "", -1)
				tokens:=strings.Fields(command)
				s, _:=strconv.Atoi(tokens[0])
				e, _:=strconv.Atoi(tokens[1])
				var items []string
				err=json.Unmarshal(J, &items)
				if err!=nil{
					fmt.Println("Error:", err)
				}
				for idx:=range items[s:e]{
					//fmt.Println(items[idx])
                    start := time.Now()
                    _, success:= kademlia.HandleClient(kademClient, "http://en.wikipedia.org/wiki/"+items[idx], mode)
                    elapsed := time.Since(start)
                    if success==true{
                    	log.Println(items[idx], elapsed)
                    	fmt.Println(items[idx], elapsed)
                	}else{
                    	log.Println(items[idx], "-1")
                    	fmt.Println(items[idx], "-1")
                	}
                    logf.Sync()
				}

            }else{
                fmt.Println("Single User Mode")
				//save mode0's log into log0.txt
	            logf, _ := os.Create("log0.txt")
	            log.SetOutput(logf)
	            defer logf.Close()
				J, err:=ioutil.ReadFile("./wikiindex/"+filename)
				if err!=nil{
					panic(err)
				}
				var items []string
				err=json.Unmarshal(J, &items)
				if err!=nil{
					fmt.Println("Error:", err)
				}
				for idx:=range items{
					//fmt.Println(items[idx])
                    start := time.Now()
                    _, success:= kademlia.HandleClient(kademClient, "http://en.wikipedia.org/wiki/"+items[idx], mode)
                    elapsed := time.Since(start)
                    if success==true{
                    	log.Println(items[idx], elapsed)
                    	fmt.Println(items[idx], elapsed)
                	}else{
                    	log.Println(items[idx], "-1")
                    	fmt.Println(items[idx], "-1")
                	}
                    logf.Sync()
				}
				/*
                f, err := os.OpenFile("./wikiindex/"+filename,os.O_RDONLY,0)
                if err != nil{
                    fmt.Printf("%v\n",err)
                    os.Exit(1)
                }
                defer f.Close()

                
                br := bufio.NewReader(f)
                for{
                    line, err := br.ReadString('\n')
                    line = strings.Replace(line,"\n","",-1)

                    if err == io.EOF {
                        break
                    }else{
                        //fmt.Printf("http://en.wikipedia.org/wiki/"+line)
                        start := time.Now()
                        _, success:= kademlia.HandleClient(kademClient, "http://en.wikipedia.org/wiki/"+line, mode)
                        elapsed := time.Since(start)
                        if success==true{
                        log.Println(line, elapsed)
                        fmt.Println(line, elapsed)
                    }else{
                        log.Println(line, "-1")
                        fmt.Println(line, "-1")
                    }
                        logf.Sync()
                    }
					
                }
            	*/
            

                //_, success:= kademlia.HandleClient(kademClient, tokens[1], mode)
				fmt.Println("-----------------------------------------------------")
				fmt.Println("Finished batchtest!")
            }



        case "get":
            if len(tokens)!=3{
                fmt.Println("fetchurl takes 2 arguments, url mode")
                break
            }
            mode, err := strconv.Atoi(tokens[2])
            if err!=nil{
                break
            }
            kademlia.HandleClient(kademClient, tokens[1], mode)

		case "showc":
			kademlia.ShowC(kademClient)

        case "find_value":
            if len(tokens) != 3{
                fmt.Println("find_node takes 2 arguments, nodeID key")
                break
            }
            nodeid, _ := kademlia.FromString(tokens[1])
            keyid, _ :=kademlia.FromString(tokens[2])
            found, targetCont := kademlia.Search_Contact(kademClient,nodeid)
            if found==true{
                kademlia.DoFindValue(kademClient, &targetCont, keyid)
            } else{
                fmt.Println("Contact not in bucket!")
            }
            

        case "find_node":
            if len(tokens) != 3{
                fmt.Println("find_node takes 2 arguments, nodeID key")
                break
            }
            nodeid, _ := kademlia.FromString(tokens[1])
            keyid, _ :=kademlia.FromString(tokens[2])
            found, targetCont := kademlia.Search_Contact(kademClient,nodeid)
            if found==true{
                kademlia.DoFindNode(kademClient, &targetCont, keyid)
            }else{
                fmt.Println("Contact not in bucket!")
            }



        case "store":
            if len(tokens) != 4{
                fmt.Println("store takes 3 arguments, nodeID key value")
                break
            }
            nodeid, _ := kademlia.FromString(tokens[1])
            keyid, _ :=kademlia.FromString(tokens[2])
            found, targetCont := kademlia.Search_Contact(kademClient,nodeid)
            if found==true{
                kademlia.DoStore(kademClient, &targetCont, keyid, []byte(tokens[3]))
                fmt.Println("")
            }else{
                fmt.Println("Contact not in bucket!")
            }


        case "ping":
            if len(tokens) != 2{
                fmt.Println("ping takes 1 argument, IP:port or NodeID")
                break
            }
            if strings.Contains(tokens[1],":"){
                iptokens:=strings.Split(tokens[1], ":")
				ipaddr, err:=net.ResolveIPAddr("ip", iptokens[0])
				if err==nil{
					kademlia.DoPing(kademClient, net.ParseIP(ipaddr.String()), kademlia.Str2Port(iptokens[1]))
				}else{
					fmt.Println("Cannot resolve IP")
				}
            }else{
				targetID, _:=kademlia.FromString(tokens[1])
				ipaddr, port:=kademlia.Find_Contact(kademClient, targetID)
				if ipaddr!=nil{
					kademlia.DoPing(kademClient, ipaddr, port)
				}else{
					fmt.Println("Cannot find contact")
				}
			}

        case "whoami":
            if len(tokens) != 1{
                fmt.Println("Whoami takes no argument")
                break
            }
            fmt.Println(kademClient.NodeID.AsString())

        case "local_find_value":
            if len(tokens) != 2{
                fmt.Println("local_find_value takes 1 argument, the key in DHT")
                break
            }
            keyID, _ := kademlia.FromString(tokens[1])
            kademlia.Local_Find_Value(kademClient, keyID)

        case "get_contact":
            if len(tokens) != 2{
                fmt.Println("get_contact takes 1 argument, the node ID")
                break
            }
            targetID, _ := kademlia.FromString(tokens[1])
            kademlia.Get_Contact(kademClient, targetID)

        case "iterativeStore":
            if len(tokens) != 3{
                fmt.Println("iterativeStore takes 2 arguments, the key and value you put in DHT")
                break
            }
			nodeid, _ :=kademlia.FromString(tokens[1])
			kademlia.IterativeStore(kademClient, nodeid, []byte(tokens[2]))

        case "iterativeFindNode":
            if len(tokens) != 2{
                fmt.Println("iterativeFindNode takes 1 argument, the node ID")
                break
            }
            nodeid, _ := kademlia.FromString(tokens[1])
            kademlia.IterativeFindNode(kademClient, nodeid)

        case "iterativeFindValue":
            if len(tokens) != 2{
                fmt.Println("iterativeFindValue takes 1 argument, the key in DHT")
                break
            }
			nodeid, _ :=kademlia.FromString(tokens[1])
			kademlia.IterativeFindValue(kademClient, nodeid)

        default:
            fmt.Println("Invalid command, the correct operation: get, whoami, local_find_value, get_contact, iterativeStore, iterativeFindNode, iterativeFindValue")
            
        }
        
    }
}

