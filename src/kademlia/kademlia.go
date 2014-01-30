package kademlia

import (
  "fmt"
  "net"
  //"time"
  "log"
  "strconv"
  "net/rpc"
  "strings"
  "net/http"
  )

const K=20
const BitNum=160

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.


// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
  NodeID ID
  Host net.IP
  Port uint16
  Localmap map[ID][]byte
  AddrTab [BitNum]Bucket
}


type Bucket struct{
  ContactLst [K]Contact
}


func GetIndexLst(lst [K]Contact) int{
    
    for counter:=0; counter<K;counter++{
      if lst[counter].Host == nil{
        return counter
      }
    }
    return K
}


func  Update(k *Kademlia, contact Contact) error{
    selfid:=k.NodeID
    requestid:=contact.NodeID

    bitindex := selfid.Xor(requestid).PrefixLen()-1
    if bitindex < 0{
      bitindex=0
    }

    full := (K==GetIndexLst(k.AddrTab[bitindex].ContactLst))
    hasContact, currentindex := Get_Contact2(k, requestid)



    if hasContact==true {

        tempcontact:=k.AddrTab[bitindex].ContactLst[currentindex]
        for j:=currentindex; j<cap(k.AddrTab[bitindex].ContactLst)-1; j++{
          k.AddrTab[bitindex].ContactLst[j]=k.AddrTab[bitindex].ContactLst[j+1]
        }
        if full==true{
          k.AddrTab[bitindex].ContactLst[K-1]=tempcontact
          return nil
        } else {
          k.AddrTab[bitindex].ContactLst[GetIndexLst(k.AddrTab[bitindex].ContactLst)]=tempcontact
          return nil
        }

    } else {
      if full==false{
        fmt.Println("Add Contact")
        fmt.Println(contact)
        k.AddrTab[bitindex].ContactLst[GetIndexLst(k.AddrTab[bitindex].ContactLst)]=contact
        return nil
      } else{
        topContact:=k.AddrTab[bitindex].ContactLst[0]
        pingSucc := DoPing(k, topContact.Host, topContact.Port)
        if pingSucc==true{
          return nil
        } else{
          for j:=0; j<K-1; j++{
            k.AddrTab[bitindex].ContactLst[j]=k.AddrTab[bitindex].ContactLst[j+1]
          }
          k.AddrTab[bitindex].ContactLst[K-1]=contact
        }
      }
    }

    full = (cap(k.AddrTab[bitindex].ContactLst)==GetIndexLst(k.AddrTab[bitindex].ContactLst))
    hasContact, currentindex = Get_Contact2(k, requestid)



    return nil
}

func Get_Contact2(kadem *Kademlia, id ID) (bool, int){

  bitindex :=  kadem.NodeID.Xor(id).PrefixLen()-1
  if bitindex < 0{
      bitindex=0
  }

  contactlst:=kadem.AddrTab[bitindex].ContactLst

  for i:=0; i<K; i++ {
      if contactlst[i].NodeID.Equals(id)==true && contactlst[i].Host != nil{
        return true, i
      }
  }

  return false, -1
}



func Get_Contact(kadem *Kademlia, id ID) (bool, int){

  bitindex :=  kadem.NodeID.Xor(id).PrefixLen()-1
  if bitindex < 0{
      bitindex=0
  }

  contactlst:=kadem.AddrTab[bitindex].ContactLst

  for i:=0; i<K; i++ {
      if contactlst[i].NodeID.Equals(id)==true && contactlst[i].Host != nil{
        fmt.Println("%v %v\n", contactlst[i].Host, contactlst[i].Port)
        return true, i
      }
  }

  fmt.Println("ERR")
  return false, -1
}

func Local_Find_Value(kadem *Kademlia, key ID) (bool, []byte){
  val, ok := kadem.Localmap[key]
  if ok ==false{
   	fmt.Println("ERR")
   } else {
		fmt.Println(val)
	}
  return ok, val
}

func Port2Str(port uint16) string{
  return strconv.Itoa(int(port))
}

func Str2Port(port string) uint16{
  i,error := strconv.Atoi(port)
  if error != nil{
    log.Fatal("Invalid port", error)
    return uint16(9999)
  }else{
    return  uint16(i)
  }
}

func StartServ(kadem *Kademlia, ipport string) bool{
  rpc.Register(kadem)
  rpc.HandleHTTP()
  l, err := net.Listen("tcp", ipport)
  if err != nil {
      log.Fatal("Listen: ", err)
      return false
  }
  iptokens:=strings.Split(ipport, ":")
  kadem.Host=net.ParseIP(iptokens[0])
  kadem.Port=Str2Port(iptokens[1])
  // Serve forever.
  go http.Serve(l, nil)
  return true
}

func DoPing(kadem *Kademlia, remoteHost net.IP, port uint16) bool{
    portstr:=Port2Str(port)
    fmt.Println("Start Ping:")
    fmt.Println(remoteHost.String()+":"+portstr)
    client, err := rpc.DialHTTP("tcp", remoteHost.String()+":"+portstr)
    if err != nil {
        //log.Fatal("DialHTTP: ", err)
        return false
    }
    ping := new(Ping)

    ping.Sender.NodeID=kadem.NodeID
    ping.Sender.Host=kadem.Host
    ping.Sender.Port=kadem.Port
    ping.MsgID = NewRandomID()

    var pong Pong

    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
        //log.Fatal("Call: ", err)
        return false
    }

    if pong.MsgID.Equals(ping.MsgID){
        go Update(kadem, pong.Sender)
        fmt.Println("Ping Success!")
        return true;
    }
    fmt.Println("Ping Failure!")
    return false;
}

func DoStore(kadem *Kademlia, remoteContact *Contact, storeKey ID, storeValue []byte) bool{
    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    fmt.Println("Start Store:")
    client, err := rpc.DialHTTP("tcp", remoteHost.String()+":"+remotePortstr)
    if err != nil {
        //log.Fatal("DialHTTP: ", err)
        return false
    }

    storeReq := new(StoreRequest)
    storeReq.Sender.NodeID=kadem.NodeID
    storeReq.Sender.Host=kadem.Host
    storeReq.Sender.Port=kadem.Port

    storeReq.MsgID=NewRandomID()

    storeReq.Key=storeKey
    storeReq.Value=storeValue

    var storeRes StoreResult

    err = client.Call("Kademlia.Store", storeReq, &storeRes)
    if err != nil {
        //log.Fatal("Call: ", err)
        return false
    }
    if storeReq.MsgID.Equals(storeRes.MsgID){
        go Update(kadem, *remoteContact)
        fmt.Println("Store Success!")
        return true;
    }
    fmt.Println("Store Failure!")
    return false;

}

func DoFindNode(kadem *Kademlia, remoteContact *Contact, searchKey ID) bool{

    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    fmt.Println("Start FoundNode:")
    client, err := rpc.DialHTTP("tcp", remoteHost.String()+":"+remotePortstr)
    if err != nil {
        //log.Fatal("DialHTTP: ", err)
        return false
    }

    findNodeReq := new(FindNodeRequest)
    findNodeReq.Sender.NodeID=kadem.NodeID
    findNodeReq.Sender.Host=kadem.Host
    findNodeReq.Sender.Port=kadem.Port

    findNodeReq.MsgID=NewRandomID()

    findNodeReq.NodeID=searchKey

    var findNodeRes FindNodeResult
    err = client.Call("Kademlia.FindNode", findNodeReq, &findNodeRes)

    if err != nil {
        //log.Fatal("Call: ", err)
        return false
    }

    if findNodeReq.MsgID.Equals(findNodeRes.MsgID){
        go Update(kadem, *remoteContact)
        fmt.Println("Found Node Success!")
        fmt.Println(findNodeRes.Nodes)
        return true;
    }
    fmt.Println("FondNode Failure!")
    return false;

}


func DoFindValue(kadem *Kademlia, remoteContact *Contact, searchKey ID) bool{
    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    fmt.Println("Start FoundValue:")
    client, err := rpc.DialHTTP("tcp", remoteHost.String()+":"+remotePortstr)
    if err != nil {
        //log.Fatal("DialHTTP: ", err)
        return false
    }

    findValueReq := new(FindValueRequest)

    findValueReq.Sender.Host=kadem.Host
    findValueReq.Sender.Port=kadem.Port
    findValueReq.Sender.NodeID=kadem.NodeID

    findValueReq.MsgID=NewRandomID()

    findValueReq.Key=searchKey

    var findValueRes FindValueResult

    err = client.Call("Kademlia.FindValue", findValueReq, &findValueRes)

    if err != nil {
        //log.Fatal("Call: ", err)
        return false
    }

    if findValueReq.MsgID.Equals(findValueRes.MsgID){
        go Update(kadem, *remoteContact)
        fmt.Println("Found Value Success!")
        if findValueRes.Value != nil{
          fmt.Println("Value Found:")
          fmt.Println(findValueRes.Value)
        }else{
          fmt.Println("Node Found:")
          fmt.Println(findValueRes.Nodes)
        }
        return true;
    }
    fmt.Println("FondValue Failure!")
    return false;




}





func NewKademlia() *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
    retNode := new(Kademlia)
    retNode.NodeID=NewRandomID()
    //retNode.NodeID,_ = FromString("c3744506eaee5ffe77b580a5676c59d5776587ca")
    retNode.Localmap = make(map[ID][]byte)
    //retNode.AddrTab=make(Bucket, BitNum)
    //retNode.Host=
    //retNode.Port=
    //retNode.Localmap["Zhengyang"]=1
    //retNode.AddrTab[152].ContactLst[0].NodeID, _=FromString("c3744506eaee5ffe77b580a5676c59d5776587cb")
    //fmt.Println(retNode.AddrTab[152].ContactLst[0].Host==nil)
    //fmt.Println(retNode.AddrTab[152].ContactLst)
    return retNode
}


