package kademlia

import (
  "fmt"
  "net"
  //"time"
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
    counter:=0
    for ; counter<cap(lst);counter++{
      if lst[counter].Host==nil{
        return counter
      }
    }
    return cap(lst)
}


func  Update(k *Kademlia, contact Contact) error{
    selfid:=k.NodeID
    requestid:=contact.NodeID
    bitindex := selfid.Xor(requestid).PrefixLen()

    full := (cap(k.AddrTab[bitindex].ContactLst)==GetIndexLst(k.AddrTab[bitindex].ContactLst))
    hasContact, currentindex := Get_Contact2(k, requestid)

    fmt.Println(full)
    fmt.Println(hasContact)

    if hasContact==true {

        tempcontact:=k.AddrTab[bitindex].ContactLst[currentindex]
        for j:=currentindex; j<cap(k.AddrTab[bitindex].ContactLst)-1; j++{
          k.AddrTab[bitindex].ContactLst[j]=k.AddrTab[bitindex].ContactLst[j+1]
        }
        if full==true{
          k.AddrTab[bitindex].ContactLst[K-1]=tempcontact
        } else {
          k.AddrTab[bitindex].ContactLst[GetIndexLst(k.AddrTab[bitindex].ContactLst)]=tempcontact
        }

    } else {
      if full==false{
        fmt.Println("Add Contact")
        fmt.Println(contact)
        k.AddrTab[bitindex].ContactLst[GetIndexLst(k.AddrTab[bitindex].ContactLst)]=contact
      } 
    }

    full = (cap(k.AddrTab[bitindex].ContactLst)==GetIndexLst(k.AddrTab[bitindex].ContactLst))
    hasContact, currentindex = Get_Contact2(k, requestid)

    fmt.Println(full)
    fmt.Println(hasContact)

    return nil
}

func Get_Contact2(kadem *Kademlia, id ID) (bool, int){

  bitindex :=  kadem.NodeID.Xor(id).PrefixLen()
  contactlst:=kadem.AddrTab[bitindex].ContactLst

  for i:=0; i<K; i++ {
      if contactlst[i].NodeID.Equals(id)==true && contactlst[i].Host != nil{
        //fmt.Println("%v %v\n", contactlst[i].Host, contactlst[i].Port)
        return true, i
      }
  }

  //fmt.Println("ERR")
  return false, -1
}



func Get_Contact(kadem *Kademlia, id ID) (bool, int){

  bitindex :=  kadem.NodeID.Xor(id).PrefixLen()
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

func Local_Find_Value(kadem *Kademlia, key ID) int{
  val, ok := kadem.Localmap[key]
  if ok ==false{
   	fmt.Println("ERR")
   	return -1
   }
	if ok == true{
		fmt.Println(val)
	}
	return 0
}



func NewKademlia() *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
    retNode := new(Kademlia)
    retNode.NodeID=NewRandomID()
    //retNode.NodeID,_ = FromString("c3744506eaee5ffe77b580a5676c59d5776587ca")
    retNode.Localmap = make(map[ID][]byte)
    //retNode.Localmap["Zhengyang"]=1
    //retNode.AddrTab[152].ContactLst[0].NodeID, _=FromString("c3744506eaee5ffe77b580a5676c59d5776587cb")
    //fmt.Println(retNode.AddrTab[152].ContactLst[0].Host==nil)
    //fmt.Println(retNode.AddrTab[152].ContactLst)
    return retNode
}


