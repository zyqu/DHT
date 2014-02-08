package kademlia

import (
  "fmt"
  "net"
  "time"
  "log"
  "strconv"
  "net/rpc"
  "strings"
  "net/http"
  "errors"
  )

const K=20
const BitNum=160
const Alpha=3
const Timeout=6
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

    bitindex := selfid.Xor(requestid).PrefixLen()

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
        fmt.Println(contact.NodeID.AsString())
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

func Search_Contact(kadem *Kademlia, id ID) (bool, Contact){

  bitindex :=  kadem.NodeID.Xor(id).PrefixLen()
  if bitindex < 0{
      bitindex=0
  }

  contactlst:=kadem.AddrTab[bitindex].ContactLst

  for i:=0; i<K; i++ {
      if contactlst[i].NodeID.Equals(id)==true && contactlst[i].Host != nil{
        return true, contactlst[i]
      }
  }

  return false, *new(Contact)
}


func Get_Contact2(kadem *Kademlia, id ID) (bool, int){

  bitindex :=  kadem.NodeID.Xor(id).PrefixLen()
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


//////////
func ShowC(kadem *Kademlia){
	for bitindex:=0; bitindex<2; bitindex++{

		contactlist:=kadem.AddrTab[bitindex].ContactLst
		for i:=0; i<K; i++{
			if(contactlist[i].Host!=nil){
				fmt.Println(contactlist[i].NodeID.AsString(), "  bitindex=", bitindex)
			}
		}

	}
}



func Get_Contact(kadem *Kademlia, id ID) (bool, int){

  bitindex :=  kadem.NodeID.Xor(id).PrefixLen()
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
		fmt.Println(string(val))
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
    fmt.Println("CLient NodeID: ",kadem.NodeID.AsString())
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
		fmt.Println("-------------------------------")
        //fmt.Println(findNodeRes.Nodes)
		for i:=0; i<len(findNodeRes.Nodes); i++{
			if findNodeRes.Nodes[i].IPAddr!=""{
				fmt.Println(findNodeRes.Nodes[i].NodeID.AsString())
			}
		}
		fmt.Println("-------------------------------")
        return true;
    }
    fmt.Println("FoundNode Failure!")
    return false;

}

func DoFindNode2(kadem *Kademlia, remoteContact *Contact, searchKey ID, succ chan bool) {

    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    fmt.Println("Start FoundNode:")
    client, err := rpc.DialHTTP("tcp", remoteHost.String()+":"+remotePortstr)
    if err != nil {
        //log.Fatal("DialHTTP: ", err)
		succ<-false
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
        //return false
		succ<-false
    }

    if findNodeReq.MsgID.Equals(findNodeRes.MsgID){
        go Update(kadem, *remoteContact)
		
		fmt.Println("Found Node successfully from: ", remoteContact.NodeID.AsString())
		fmt.Println("-------------------------------")
        //fmt.Println(findNodeRes.Nodes)
		for i:=0; i<len(findNodeRes.Nodes); i++{
			if findNodeRes.Nodes[i].IPAddr!=""{
				newcontact:=new(Contact)
				newcontact.NodeID=findNodeRes.Nodes[i].NodeID
				newcontact.Host=net.ParseIP(findNodeRes.Nodes[i].IPAddr)
				newcontact.Port=findNodeRes.Nodes[i].Port
				newcontact.queried=false
				go Update(kadem, *newcontact)
				fmt.Println(findNodeRes.Nodes[i].NodeID.AsString())
			}
		}
		fmt.Println("-------------------------------")
        //return true;
		succ<-true
    }else{
		fmt.Println("FoundNode Failure!")
		//return false;
		succ<-false
	}
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

		
		fmt.Println("Find Value Success!")
        if findValueRes.Value != nil{
          fmt.Println("Value Found:")
          fmt.Println(string(findValueRes.Value))
        }else{
          fmt.Println("Node Found:")
          //fmt.Println(findValueRes.Nodes)
		  for i:=0; i<len(findValueRes.Nodes); i++{
			  if findValueRes.Nodes[i].IPAddr!=""{
				  fmt.Println(findValueRes.Nodes[i].NodeID.AsString())
			  }
		  }
        }
        return true;
    }
    fmt.Println("FindValue Failure!")
    return false;




}





func IterativeStore(kadem *Kademlia, storeKey ID, storeValue []byte) {
	//bitindex:=kadem.NodeID.Xor(storeKey).PrefixLen()
	
	
	
}

func IterativeFindNode(kadem *Kademlia, searchKey ID) (bool, error) {
	finished:=false
	shortlist:=getClosestContacts(kadem, searchKey)
	if len(shortlist)==0{
		fmt.Println("No contact in list")
		return false, errors.New("No contact")
	}
	closestNode:=shortlist[0]
	nodes:=make([]Contact, 0)
	for !finished{
		fmt.Println("shortlist length=", len(shortlist))
		if(len(shortlist)==0){
			finished=true
			break
		}
		succ:=make(chan bool)
		to:=make(chan int, Alpha)
		for i:=0; i<len(shortlist); i++{
			found, searchCon:=Search_Contact(kadem, shortlist[i].NodeID)
			if found{
				go DoFindNode2(kadem, &searchCon, searchKey, succ)
			}else{
				fmt.Println("Cannot perform FindNode")
				return false, errors.New("Cannot perform FindNode")
			}
			
			time.Sleep(time.Second*Timeout)
			to<-1
		}
		for i:=0; i<len(shortlist); i++{
			select{
			case <-succ:
				shortlist[i].queried=true
				nodes=append(nodes, *shortlist[i])
				if len(nodes)>=K{
					finished=true
					break
				}
			case <-to:
				copy(shortlist[i:], shortlist[i+1:])
				shortlist[len(shortlist)-1]=nil
				shortlist=shortlist[:len(shortlist)-1]
				continue
			}
			
		}
		shortlist=getClosestContacts(kadem, searchKey)
		if len(shortlist)<=0{
			return false, errors.New("Cannot find it.")
		}
		d1:=kadem.NodeID.Xor(closestNode.NodeID).PrefixLen()
		d2:=kadem.NodeID.Xor(shortlist[0].NodeID).PrefixLen()
		if d1<=d2{
			//no node return closer than closest node already seen
			finished=true
			break
		}else{
			closestNode=shortlist[0]
		}
	}
	fmt.Println("List of iterative find_node: ")
	for i:=0; i<len(nodes); i++{
		fmt.Println(nodes[i].NodeID.AsString())
	}
	return true, nil
}
func getClosestContacts(kadem *Kademlia, key ID) []*Contact{
	//Alpha non-contacted closest contacts
	ret:=make([]*Contact, 0)
	bitindex:=kadem.NodeID.Xor(key).PrefixLen()
    if bitindex < 0{
        bitindex=0
    }
	length:=0
	//find Alpha closest nodes
	for i:=0; i<K; i++{
		if kadem.AddrTab[bitindex].ContactLst[i].Host!=nil{
			ret=append(ret, &kadem.AddrTab[bitindex].ContactLst[i])
			length++
			if(length>=Alpha){
				return ret
			}
		}
	}
	for i:=1; (i<bitindex||i+bitindex<BitNum)&&length<Alpha; i++{
		if i<=bitindex{
			for j:=0; j<K&&length<K; j++{
				if kadem.AddrTab[bitindex-i].ContactLst[j].queried==false&&kadem.AddrTab[bitindex-i].ContactLst[j].Host!=nil{
					ret=append(ret, &kadem.AddrTab[bitindex-i].ContactLst[j])
					length++
					if(length>=Alpha){
						return ret
					}
				}
			}
		}
		if i+bitindex<BitNum{
			for j:=0; j<K&&length<K; j++{
				if kadem.AddrTab[i+bitindex].ContactLst[j].queried==false&&kadem.AddrTab[i+bitindex].ContactLst[j].Host!=nil{
					ret=append(ret, &kadem.AddrTab[i+bitindex].ContactLst[j])
					length++
					if(length>=Alpha){
						return ret
					}
				}
			}
			
		}
	}
	return ret
}

func IterativeFindValue(kadem *Kademlia, searchKey ID) {
	
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


