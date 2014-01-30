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
    //firstPeerStr := args[1]

    fmt.Printf("kademlia starting up!\n")
    //kademClient := kademlia.NewKademlia()
    kademServer := kademlia.NewKademlia()
    kademlia.StartServ(kademServer,listenStr)
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

        case "ping":
            if len(tokens) != 2{
                fmt.Println("ping takes 1 argument, IP:port or NodeID")
                break
            }
            if strings.Contains(tokens[1],":"){
                iptokens:=strings.Split(tokens[1], ":")
                kademlia.DoPing(kademServer, net.ParseIP(iptokens[0]), kademlia.Str2Port(iptokens[1]))
            }

        case "whoami":
            if len(tokens) != 1{
                fmt.Println("Whoami takes no argument")
                break
            }
            fmt.Println(kademServer.NodeID.AsString())

        case "local_find_value":
            if len(tokens) != 2{
                fmt.Println("local_find_value takes 1 argument, the key in DHT")
                break
            }
            keyID, _ := kademlia.FromString(tokens[1])
            kademlia.Local_Find_Value(kademServer, keyID)

        case "get_contact":
            if len(tokens) != 2{
                fmt.Println("get_contact takes 1 argument, the node ID")
                break
            }
            targetID, _ := kademlia.FromString(tokens[1])
            kademlia.Get_Contact(kademServer, targetID)

        case "iterativeStore":
            if len(tokens) != 3{
                fmt.Println("iterativeStore takes 2 arguments, the key and value you put in DHT")
                break
            }
            fmt.Println("iterativeStore")

        case "iterativeFindNode":
            if len(tokens) != 2{
                fmt.Println("iterativeFindNode takes 1 argument, the node ID")
                break
            }
            fmt.Println("iterativeFindNode")

        case "iterativeFindValue":
            if len(tokens) != 2{
                fmt.Println("iterativeFindValue takes 1 argument, the key in DHT")
                break
            }
            fmt.Println("iterativeFindValue")

        default:
            fmt.Println("Invalid command, the correct operation: whoami, local_find_value, get_contact, iterativeStore, iterativeFindNode, iterativeFindValue")
            
        }
        
    }
}

