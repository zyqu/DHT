package kademlia

import (
  "fmt"
  "net"
  "time"
  "log"
  "strconv"
  "net/rpc"
  "strings"
  "errors"
  "os"
  "os/exec"
  "net/http"
  "io/ioutil"
  )

const K=20
const BitNum=160
const Alpha=3
const Timeout =300
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.


// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
	NodeID ID
	Host net.IP
	Port uint16
	Localmap map[ID][]byte
	AddrTab [BitNum]Bucket
	ch chan Contact

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

func  Update2(k *Kademlia) error{
	contact:=<-k.ch
    selfid:=k.NodeID
    requestid:=contact.NodeID

    bitindex := selfid.Xor(requestid).PrefixLen()-1

    if bitindex < 0{
      bitindex=0
    }
	
	if k.NodeID==contact.NodeID{
		return errors.New("Updating itself.")
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
        //fmt.Println("Add Contact")
        //fmt.Println(contact.NodeID.AsString())
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

  bitindex :=  kadem.NodeID.Xor(id).PrefixLen()-1
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


//////////
//for test
func ShowC(kadem *Kademlia){
	for bitindex:=0; bitindex<BitNum; bitindex++{

		contactlist:=kadem.AddrTab[bitindex].ContactLst
		for i:=0; i<K; i++{
			if(contactlist[i].Host!=nil){
				fmt.Println(contactlist[i].NodeID.AsString(), "  bitindex=", bitindex)
			}
		}

	}
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

func Find_Contact(kadem *Kademlia, id ID) (net.IP, uint16){
	bitindex:=kadem.NodeID.Xor(id).PrefixLen()-1
	if bitindex<0{
		bitindex=0
	}
	contactlst:=kadem.AddrTab[bitindex].ContactLst
	for i:=0; i<K; i++{
		if contactlst[i].NodeID.Equals(id)==true&&contactlst[i].Host!=nil{
			return contactlst[i].Host, contactlst[i].Port
		}
	}
	return nil, 0
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
  //os.Mkdir("./tmp/"+kadem.NodeID.AsString(), 0700)
  
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
        //go Update(kadem, pong.Sender)
		kadem.ch<-pong.Sender
		go Update2(kadem)

		//fmt.Println("Ping Success!")
        return true;
    }
    //fmt.Println("Ping Failure!")
    return false;
}
func DoStore(kadem *Kademlia, remoteContact *Contact, storeKey ID, storeValue []byte) bool{
    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)

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
        //go Update(kadem, *remoteContact)
		kadem.ch<-*remoteContact
		go Update2(kadem)
		
        return true;
    }
    fmt.Println(remoteContact.NodeID.AsString(), " Store Failure!")
    return false;

}
func DoStore2(kadem *Kademlia, remoteContact *Contact, storeKey ID, storeValue []byte, storeBody string) bool{
    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)

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
	storeReq.Body=storeBody

    var storeRes StoreResult

    err = client.Call("Kademlia.Store2", storeReq, &storeRes)
    if err != nil {
        //log.Fatal("Call: ", err)
        return false
    }
    if storeReq.MsgID.Equals(storeRes.MsgID){
        //go Update(kadem, *remoteContact)
		kadem.ch<-*remoteContact
		go Update2(kadem)
		
        return true;
    }
    fmt.Println(remoteContact.NodeID.AsString(), " Store Failure!")
    return false;

}

func DoFindNode(kadem *Kademlia, remoteContact *Contact, searchKey ID) bool{

    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    //fmt.Println("Start FoundNode:")
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
        //go Update(kadem, *remoteContact)
		kadem.ch<-*remoteContact
		go Update2(kadem)

		
		//fmt.Println("Found Node Success!")
		//fmt.Println("-------------------------------")
        //fmt.Println(findNodeRes.Nodes)
		for i:=0; i<len(findNodeRes.Nodes); i++{
			if findNodeRes.Nodes[i].IPAddr!=""{
				newcontact:=new(Contact)
				newcontact.NodeID=findNodeRes.Nodes[i].NodeID
				newcontact.Host=net.ParseIP(findNodeRes.Nodes[i].IPAddr)
				newcontact.Port=findNodeRes.Nodes[i].Port
				newcontact.queried=false
				//go Update(kadem, *newcontact)
				kadem.ch<-*newcontact
				go Update2(kadem)
				fmt.Println(findNodeRes.Nodes[i].NodeID.AsString())
			}
		}
		//fmt.Println("-------------------------------")
        return true;
    }
    fmt.Println("FoundNode Failure!")
    return false;

}
func DoFindNode2(kadem *Kademlia, remoteContact *Contact, searchKey ID, succ chan bool) {

    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    //fmt.Println("Start FoundNode:")
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
        //go Update(kadem, *remoteContact)
		kadem.ch<-*remoteContact
		go Update2(kadem)
		
		//fmt.Println("Found Node successfully from: ", remoteContact.NodeID.AsString())
		//fmt.Println("-------------------------------")
        //fmt.Println(findNodeRes.Nodes)
		for i:=0; i<len(findNodeRes.Nodes); i++{
			if findNodeRes.Nodes[i].IPAddr!=""{
				newcontact:=new(Contact)
				newcontact.NodeID=findNodeRes.Nodes[i].NodeID
				newcontact.Host=net.ParseIP(findNodeRes.Nodes[i].IPAddr)
				newcontact.Port=findNodeRes.Nodes[i].Port
				newcontact.queried=false
				//go Update(kadem, *newcontact)
				kadem.ch<-*newcontact
				go Update2(kadem)
				//fmt.Println(findNodeRes.Nodes[i].NodeID.AsString())
			}
		}
		//fmt.Println("-------------------------------")
        //return true;
		succ<-true
    }else{
		//fmt.Println("FoundNode Failure!")
		//return false;
		succ<-false
	}
}
func DoFindValue(kadem *Kademlia, remoteContact *Contact, searchKey ID) bool{
    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    //fmt.Println("Start FoundValue:")
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
        //go Update(kadem, *remoteContact)
		kadem.ch<-*remoteContact
		go Update2(kadem)
		
		//fmt.Println("Find Value Success!")
        if findValueRes.Value != nil{
          //fmt.Println("Value Found:")
          //fmt.Println(string(findValueRes.Value))
        }else{
          fmt.Println("Node Found:")
          //fmt.Println(findValueRes.Nodes)
		  for i:=0; i<len(findValueRes.Nodes); i++{
			  if findValueRes.Nodes[i].IPAddr!=""{
				  newcontact:=new(Contact)
				  newcontact.NodeID=findValueRes.Nodes[i].NodeID
				  newcontact.Host=net.ParseIP(findValueRes.Nodes[i].IPAddr)
				  newcontact.Port=findValueRes.Nodes[i].Port
				  newcontact.queried=false
				  //go Update(kadem, *newcontact)
				  kadem.ch<-*newcontact
				  go Update2(kadem)
				  //fmt.Println(findValueRes.Nodes[i].NodeID.AsString())
			  }
		  }
        }
        return true;
    }
    //fmt.Println("FindValue Failure!")
    return false;




}
type Ret struct{
	value []byte
	nodeFound bool
	from ID
	body string
}
func DoFindValue2(kadem *Kademlia, remoteContact *Contact, searchKey ID, ret chan Ret) bool{
    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    //fmt.Println("Start FoundValue:")
    client, err := rpc.DialHTTP("tcp", remoteHost.String()+":"+remotePortstr)
    if err != nil {
        //log.Fatal("DialHTTP: ", err)
		
		var retu Ret
		retu.value=nil
		retu.nodeFound=false
		ret<-retu
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
		var retu Ret
		retu.value=nil
		retu.nodeFound=false
		ret<-retu
        return false
    }

    if findValueReq.MsgID.Equals(findValueRes.MsgID){
        //go Update(kadem, *remoteContact)
		kadem.ch<-*remoteContact
		go Update2(kadem)
		
		//fmt.Println("Find Value Success!")
        if findValueRes.Value != nil{
          //fmt.Println("Value Found:", string(findValueRes.Value), " from ", remoteContact.NodeID.AsString())
		  var retu Ret
		  retu.value=findValueRes.Value
		  retu.nodeFound=false
		  retu.from=remoteContact.NodeID
		  ret<-retu
        }else{
          //fmt.Println("Node Found:")
          //fmt.Println(findValueRes.Nodes)
		  for i:=0; i<len(findValueRes.Nodes); i++{
			  if findValueRes.Nodes[i].IPAddr!=""{
				  newcontact:=new(Contact)
				  newcontact.NodeID=findValueRes.Nodes[i].NodeID
				  newcontact.Host=net.ParseIP(findValueRes.Nodes[i].IPAddr)
				  newcontact.Port=findValueRes.Nodes[i].Port
				  newcontact.queried=false
				  //go Update(kadem, *newcontact)
				  kadem.ch<-*newcontact
				  go Update2(kadem)
				  //fmt.Println(findValueRes.Nodes[i].NodeID.AsString())
			  }
		  }
		  var retu Ret
		  retu.value=nil
		  retu.nodeFound=true
		  ret<-retu
        }
        return true;
    }
    //fmt.Println("FindValue Failure!")
	var retu Ret
	retu.value=nil
	retu.nodeFound=false
	ret<-retu
    return false;




}
func IterativeStore(kadem *Kademlia, storeKey ID, storeValue []byte)  error{
	nodes, err:=IterativeFindNode(kadem, storeKey)
	if err!=nil{
		fmt.Println("IterativeStore Error: cannot perform IterativeFindNode")
		return errors.New("cannot perform IterativeFindNode")
	}
	fmt.Println("IterativeStore node list: ")
	fmt.Println("---------------------------------------")

	for i:=0; i<len(nodes); i++{
		r:=DoStore(kadem, &nodes[i], storeKey, storeValue)
		if !r{
			fmt.Println("---------------------------------------")
			return errors.New("DoStore error")
		}
		fmt.Println(nodes[i].NodeID.AsString())
	}
	fmt.Println("---------------------------------------")
	return nil
}
func IterativeFindNode(kadem *Kademlia, searchKey ID) ([]Contact, error) {
	finished:=false
	shortlist:=getClosestContacts(kadem, searchKey)
	if len(shortlist)==0{
		fmt.Println("No contact in list")
		return nil, errors.New("No contact")
	}
	closestNode:=shortlist[0]
	//to:=make(chan int, 1)
	//go setTimer(to)
	nodes:=make([]Contact, 0)
	for !finished{

		if(len(shortlist)==0){
			finished=true
			break
		}
		succ:=make(chan bool)
		foundKey:=false
		////////
		to:=make(chan int, 1)
		
		for i:=0; i<len(shortlist); i++{
			
			found, searchCon:=Search_Contact(kadem, shortlist[i].NodeID)
			if found{
				go DoFindNode2(kadem, &searchCon, searchKey, succ)
				go setTimer(to)
			}else{
				fmt.Println("Cannot perform FindNode")
				iterativeHelper(kadem)
				return nil, errors.New("Cannot perform FindNode")
			}
			
		}
		if foundKey{
			/*
			fmt.Println("Found: ", searchKey.AsString())
			fmt.Println("List of iterative find_node: ")
			for i:=0; i<len(nodes); i++{
				fmt.Println(nodes[i].NodeID.AsString())
			}
			iterativeHelper(kadem)
			return true, nil
			*/
			finished=true
			break
		}
		for i:=0; i<len(shortlist); i++{
			select{
			case <-succ:
				//if node found and successfully pinged
				if shortlist[i].NodeID==searchKey{
					nodes=append(nodes, *shortlist[i])
					fmt.Println("Searchkey is found.")
					finished=true
					break
				}
				setQueried(kadem, shortlist[i].NodeID)
				nodes=append(nodes, *shortlist[i])
				if len(nodes)>=K{
					finished=true
					break
				}
			case <-to:
				fmt.Println("timeout")
				finished=true
				break
			}
			
		}
		shortlist=getClosestContacts(kadem, searchKey)
		if len(shortlist)<=0{
			finished=true
			break
		}
		d1:=kadem.NodeID.Xor(closestNode.NodeID).PrefixLen()-1
		if d1<0{
			d1=0
		}
		d2:=kadem.NodeID.Xor(shortlist[0].NodeID).PrefixLen()-1
		if d2<0{
			d2=0
		}
		if d1<d2{
			//no node return closer than closest node already seen
			finished=true
			break
		}else{
			closestNode=shortlist[0]
		}
	}
	
	fmt.Println("List of iterativeFindNode: ")
	for i:=0; i<len(nodes); i++{
		fmt.Println(nodes[i].NodeID.AsString())
	}
	
	iterativeHelper(kadem)
	return nodes, nil
}
func getClosestContacts(kadem *Kademlia, key ID) []*Contact{
	//Alpha non-contacted closest contacts
	ret:=make([]*Contact, 0)
	bitindex:=kadem.NodeID.Xor(key).PrefixLen()-1
	if bitindex<0{
		bitindex=0
	}
    if bitindex < 0{
        bitindex=0
    }
	length:=0
	//find Alpha closest nodes

	for i:=0; i<K; i++{
		if kadem.AddrTab[bitindex].ContactLst[i].queried==false&&kadem.AddrTab[bitindex].ContactLst[i].Host!=nil{
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
func iterativeHelper(kadem *Kademlia){
	for bitindex:=0; bitindex<BitNum; bitindex++{
		for i:=0; i<K; i++{
			kadem.AddrTab[bitindex].ContactLst[i].queried=false
		}
	}
}
func setQueried(kadem *Kademlia, id ID){
	bitindex:=kadem.NodeID.Xor(id).PrefixLen()-1
	if bitindex<0{
		bitindex=0
	}
	for i:=0; i<len(kadem.AddrTab[bitindex].ContactLst); i++{
		if id.Equals(kadem.AddrTab[bitindex].ContactLst[i].NodeID){
			kadem.AddrTab[bitindex].ContactLst[i].queried=true
			break
		}
	}
}
func setTimer(t chan int){
	time.Sleep(time.Millisecond*Timeout)
	t<-1
}
func IterativeFindValue(kadem *Kademlia, searchKey ID) (bool, error) {
	finished:=false
	shortlist:=getClosestContacts(kadem, searchKey)
	
	if len(shortlist)==0{
		fmt.Println("No contact in list")
		return false, errors.New("No contact")
	}
	closestNode:=shortlist[0]
	to:=make(chan int, 1)
	go setTimer(to)
	nodes:=make([]Contact, 0)
	for !finished{

		if len(shortlist)==0{
			finished=true
			break
		}
		closestNode=shortlist[0]
		ret:=make(chan Ret)
		for i:=0; i<len(shortlist); i++{
			found, searchCon:=Search_Contact(kadem, shortlist[i].NodeID)
			if found{
				//fmt.Println("DoFindValue: ", searchCon.NodeID.AsString())
				go DoFindValue2(kadem, &searchCon, searchKey, ret)

			}else{
				fmt.Println("Cannot perform FindValue")
				iterativeHelper(kadem)
				return false, errors.New("Cannot perform FindValue")
			}
			
		}
		for i:=0; i<len(shortlist); i++{
			select{
			case r:=<-ret:
				if r.value!=nil{
					//vallue found
					fmt.Println("iterative find value: ", string(r.value), " from ", r.from.AsString())
					iterativeHelper(kadem)
					//store to closestNode
					DoStore(kadem, closestNode, searchKey, r.value)
					
					return true, nil
				}else{
					//shortlist[i].queried=true
					/////
					id:=shortlist[i].NodeID
					bitindex:=kadem.NodeID.Xor(id).PrefixLen()-1
					if bitindex<0{
						bitindex=0
					}
					for i:=0; i<len(kadem.AddrTab[bitindex].ContactLst); i++{
						if id.Equals(kadem.AddrTab[bitindex].ContactLst[i].NodeID){
							//fmt.Println("queried", shortlist[i].NodeID)
							kadem.AddrTab[bitindex].ContactLst[i].queried=true
							break
						}
					}
					nodes=append(nodes, *shortlist[i])
					if len(nodes)>=K{
						finished=true
						break
					}
				}
			case <-to:
				fmt.Println("timeout")
				finished=true
				break
			}
		}
		shortlist=getClosestContacts(kadem, searchKey)
		if len(shortlist)<=0{
			finished=true
			break
		}
		d1:=kadem.NodeID.Xor(closestNode.NodeID).PrefixLen()-1
		if d1<0{
			d1=0
		}
		d2:=kadem.NodeID.Xor(shortlist[0].NodeID).PrefixLen()-1
		if d2<0{
			d2=0
		}
		if d1<d2{
			//no node return closer than closest node already seen
			finished=true
			break
		}else{
			closestNode=shortlist[0]
		}
	}
	
	fmt.Println("iterativeFindValue ERR")
	iterativeHelper(kadem)
	return true, nil
}

/////////////
//DoFindValue3
func DoFindValue3(kadem *Kademlia, remoteContact *Contact, searchKey ID, ret chan Ret) bool{
    remoteHost:=remoteContact.Host
    remotePortstr:=Port2Str(remoteContact.Port)
    //fmt.Println("Start FoundValue:")
    client, err := rpc.DialHTTP("tcp", remoteHost.String()+":"+remotePortstr)
    if err != nil {
        //log.Fatal("DialHTTP: ", err)
		
		var retu Ret
		retu.value=nil
		retu.nodeFound=false
		ret<-retu
        return false
    }

    findValueReq := new(FindValueRequest)

    findValueReq.Sender.Host=kadem.Host
    findValueReq.Sender.Port=kadem.Port
    findValueReq.Sender.NodeID=kadem.NodeID

    findValueReq.MsgID=NewRandomID()

    findValueReq.Key=searchKey

    var findValueRes FindValueResult
    err = client.Call("Kademlia.FindValue2", findValueReq, &findValueRes)

    if err != nil {
        //log.Fatal("Call: ", err)
		var retu Ret
		retu.value=nil
		retu.nodeFound=false
		ret<-retu
        return false
    }

    if findValueReq.MsgID.Equals(findValueRes.MsgID){
        //go Update(kadem, *remoteContact)
		kadem.ch<-*remoteContact
		go Update2(kadem)
		
		//fmt.Println("Find Value Success!")
        if findValueRes.Value != nil{
          //fmt.Println("Value Found:", string(findValueRes.Value), " from ", remoteContact.NodeID.AsString())
		  var retu Ret
		  retu.value=findValueRes.Value
		  retu.body=findValueRes.Body
		  retu.nodeFound=false
		  retu.from=remoteContact.NodeID
		  ret<-retu
        }else{
          //fmt.Println("Node Found:")
          //fmt.Println(findValueRes.Nodes)
		  for i:=0; i<len(findValueRes.Nodes); i++{
			  if findValueRes.Nodes[i].IPAddr!=""{
				  newcontact:=new(Contact)
				  newcontact.NodeID=findValueRes.Nodes[i].NodeID
				  newcontact.Host=net.ParseIP(findValueRes.Nodes[i].IPAddr)
				  newcontact.Port=findValueRes.Nodes[i].Port
				  newcontact.queried=false
				  //go Update(kadem, *newcontact)
				  kadem.ch<-*newcontact
				  go Update2(kadem)
				  //fmt.Println(findValueRes.Nodes[i].NodeID.AsString())
			  }
		  }
		  var retu Ret
		  retu.value=nil
		  retu.body=""
		  retu.nodeFound=true
		  ret<-retu
        }
        return true;
    }
    //fmt.Println("FindValue Failure!")
	var retu Ret
	retu.value=nil
	retu.nodeFound=false
	ret<-retu
    return false;




}

func IterativeFindValue2(kadem *Kademlia, searchKey ID) (string, error) {
	finished:=false
	shortlist:=getClosestContacts(kadem, searchKey)
	
	if len(shortlist)==0{
		//fmt.Println("No contact in list")
		return "", errors.New("No contact")
	}
	closestNode:=shortlist[0]
	to:=make(chan int, 1)
	go setTimer(to)
	nodes:=make([]Contact, 0)
	for !finished{

		if len(shortlist)==0{
			finished=true
			break
		}
		closestNode=shortlist[0]
		ret:=make(chan Ret)
		for i:=0; i<len(shortlist); i++{
			found, searchCon:=Search_Contact(kadem, shortlist[i].NodeID)
			if found{
				//fmt.Println("DoFindValue: ", searchCon.NodeID.AsString())
				go DoFindValue3(kadem, &searchCon, searchKey, ret)

			}else{
				//fmt.Println("Cannot perform FindValue")
				iterativeHelper(kadem)
				return "", errors.New("Cannot perform FindValue")
			}
			
		}
		for i:=0; i<len(shortlist); i++{

			select{
			case r:=<-ret:
				
				if r.value!=nil{
					//vallue found
					//fmt.Println("iterative find value: ", string(r.value), " from ", r.from.AsString())
					iterativeHelper(kadem)
					//store to closestNode

					DoStore2(kadem, closestNode, searchKey, r.value, r.body)
					
					return string(r.body), nil
				}else{
					//shortlist[i].queried=true
					/////
					id:=shortlist[i].NodeID
					bitindex:=kadem.NodeID.Xor(id).PrefixLen()-1
					if bitindex<0{
						bitindex=0
					}
					for i:=0; i<len(kadem.AddrTab[bitindex].ContactLst); i++{
						if id.Equals(kadem.AddrTab[bitindex].ContactLst[i].NodeID){
							//fmt.Println("queried", shortlist[i].NodeID)
							kadem.AddrTab[bitindex].ContactLst[i].queried=true
							break
						}
					}
					nodes=append(nodes, *shortlist[i])
					if len(nodes)>=K{
						finished=true
						break
					}
				}
			case <-to:
				//fmt.Println("timeout")
				finished=true
				break
			}
		}

		shortlist=getClosestContacts(kadem, searchKey)
		if len(shortlist)<=0{
			finished=true
			break
		}
		d1:=kadem.NodeID.Xor(closestNode.NodeID).PrefixLen()-1
		if d1<0{
			d1=0
		}
		d2:=kadem.NodeID.Xor(shortlist[0].NodeID).PrefixLen()-1
		if d2<0{
			d2=0
		}
		if d1<d2{
			//no node return closer than closest node already seen
			finished=true
			break
		}else{
			closestNode=shortlist[0]
		}
	}
	
	//fmt.Println("iterativeFindValue ERR")
	iterativeHelper(kadem)
	return "", errors.New("Not found")
}




func NewKademlia() *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
    retNode := new(Kademlia)
    retNode.NodeID=NewRandomID()
    //retNode.NodeID,_ = FromString("c3744506eaee5ffe77b580a5676c59d5776587ca")
    retNode.Localmap = make(map[ID][]byte)
	  retNode.ch=make(chan Contact, 1)
    //retNode.AddrTab=make(Bucket, BitNum)
    //retNode.Host=
    //retNode.Port=
    //retNode.Localmap["Zhengyang"]=1
    //retNode.AddrTab[152].ContactLst[0].NodeID, _=FromString("c3744506eaee5ffe77b580a5676c59d5776587cb")
    //fmt.Println(retNode.AddrTab[152].ContactLst[0].Host==nil)
    //fmt.Println(retNode.AddrTab[152].ContactLst)
	
	
    return retNode
}

func check(e error){
  if e != nil{
    panic(e)
  }
}

func FetchUrl(kadem *Kademlia, url string)(int){
  webpageDSroot:="./webpageDS/"
  strkey:=strings.Replace(url,"http://en.wikipedia.org/wiki/","",1)
  filename:=webpageDSroot+strkey+".html"
  fmt.Println(filename)
  resp, err := http.Get(url)
  if err != nil{
    fmt.Printf("%s",err)
    return -1
    //url not accessabile
  }else{
    defer resp.Body.Close()
    body, bodyerr := ioutil.ReadAll(resp.Body)
    if bodyerr != nil{
      fmt.Printf("%s",bodyerr)
      return -1
      //read response body error
    }else{
      f, openerr := os.Create(filename)
      check(openerr)

      strbody:=string(body[:])

      
      stringconvert , _ := HTMLParser(strbody)
      byteconvert:=[]byte(stringconvert)

      //write to localmap
      urlID,_ := FromString(url)
      kadem.Localmap[urlID]=byteconvert

      _,writeerr := f.Write(byteconvert)
      check(writeerr)

      f.Sync()
      f.Close()

      return 0
    }
  }
  return 0
}

func HTMLParser(body string) (string, error){
	/*
	if _, err:=os.Stat(file); os.IsNotExist(err){
		fmt.Println("No such file")
		return false, errors.New("No such file")
	}
	cmd:=exec.Command("python", "src/href.py", file)
	cmd.Run()
	return true, nil
	*/
	cmd:=exec.Command("python", "src/href.py", string(body))
	out, err:=cmd.Output()
	if err!=nil{
		fmt.Println("ERROR with executing python script ", err)
		return "", err
	}
	//fmt.Println(string(out))
	return string(out), nil
	
}

func HandleClient(kadem *Kademlia, url string) string{
    webpageDSroot:="./webpageDS/"
    strkey:=strings.Replace(url,"http://en.wikipedia.org/wiki/","",1)
    filename:=webpageDSroot+strkey+".html"
	ret1, err1:=ioutil.ReadFile(filename)
	if err1==nil{
		return string(ret1)
	}
	key:=Hashcode(strkey)
	ret2, err2:=IterativeFindValue2(kadem, key)
	if err2==nil{
        
		return string(ret2)
	}
	FetchUrl(kadem, url)
	ret3, _:=ioutil.ReadFile(filename)
	
	kadem.Localmap[key]=[]byte(strkey)
	return string(ret3)
}

