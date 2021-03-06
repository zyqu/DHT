package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
    "net"
    "fmt"
	"io/ioutil"
	"os"
	//"errors"
    )


// Host identification.
type Contact struct {
    NodeID ID
    Host net.IP
    Port uint16
	queried bool
}


// PING
type Ping struct {
    Sender Contact
    MsgID ID
}

type Pong struct {
    MsgID ID
    Sender Contact
}


func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
    // This one's a freebie.
    //go Update(k, ping.Sender)
	k.ch<-ping.Sender
	go Update2(k)
	
    fmt.Println("Sever NodeID: ", k.NodeID.AsString())
    pong.MsgID = CopyID(ping.MsgID)
    pong.Sender.NodeID=k.NodeID
    pong.Sender.Host=k.Host
    pong.Sender.Port=k.Port

    return nil
}


// STORE
type StoreRequest struct {
    Sender Contact
    MsgID ID
    Key ID
    Value []byte
	Body string
}

type StoreResult struct {
    MsgID ID
    Err error
	
}

func (k *Kademlia) Store(req StoreRequest, res *StoreResult) error {
    // TODO: Implement.
    //go Update(k, req.Sender)
	k.ch<-req.Sender
	go Update2(k)
	
	k.Localmap[req.Key]=req.Value
    //fmt.Println("\n")
    res.MsgID=CopyID(req.MsgID)
    res.Err=nil
    return nil
}

func (k *Kademlia) Store2(req StoreRequest, res *StoreResult) error {
    // TODO: Implement.
    //go Update(k, req.Sender)
	k.ch<-req.Sender
	go Update2(k)
	
	k.Localmap[req.Key]=req.Value
	webpageDSroot:="./webpageDS/"
	filename:=webpageDSroot+string(req.Value)+".html"
    f, openerr := os.Create(filename)
    check(openerr)

    strbody:=req.Body[:]

    /*
    stringconvert , _ := HTMLParser(strbody)
    byteconvert:=[]byte(stringconvert)

    _,writeerr := f.Write(byteconvert)
    check(writeerr)
	*/
	_, writeerr:=f.Write([]byte(strbody))
    check(writeerr)
    f.Sync()
    f.Close()
	perr:=HTMLParser(filename)
  	if perr!=nil{
		fmt.Println(perr)
	}

    f.Sync()
    f.Close()
	
    //fmt.Println("\n")
    res.MsgID=CopyID(req.MsgID)
    res.Err=nil
    return nil
}


// FIND_NODE
type FindNodeRequest struct {
    Sender Contact
    MsgID ID
    NodeID ID
}

type FoundNode struct {
    IPAddr string
    Port uint16
    NodeID ID
}

type FindNodeResult struct {
    MsgID ID
    Nodes []FoundNode
    Err error
}


func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
    // TODO: Implement.
    //go Update(k, req.Sender)
	k.ch<-req.Sender
	go Update2(k)

    bitindex := k.NodeID.Xor(req.NodeID).PrefixLen()-1
	if bitindex<0{
		bitindex=0
	}
    tempFoundNode:=new(FoundNode)
    FoundNodelst := make([]FoundNode, K)
    for i:=0;i<len(k.AddrTab[bitindex].ContactLst);i++{
		if(k.AddrTab[bitindex].ContactLst[i].Host!=nil&&k.AddrTab[bitindex].ContactLst[i].NodeID!=req.Sender.NodeID){
			tempFoundNode.NodeID=k.AddrTab[bitindex].ContactLst[i].NodeID
			tempFoundNode.Port=k.AddrTab[bitindex].ContactLst[i].Port
			tempFoundNode.IPAddr=k.AddrTab[bitindex].ContactLst[i].Host.String()
			FoundNodelst[i]=*tempFoundNode
		}
    }
    res.Nodes=FoundNodelst


    res.MsgID=CopyID(req.MsgID)
    res.Err=nil
    return nil
}


// FIND_VALUE
type FindValueRequest struct {
    Sender Contact
    MsgID ID
    Key ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
    MsgID ID
    Value []byte
    Nodes []FoundNode
    Err error
	Body string
}

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
    // TODO: Implement.
    //go Update(k, req.Sender)
	k.ch<-req.Sender
	go Update2(k)
	
	found, val:=Local_Find_Value(k,req.Key)
    res.Value=val
    if found == false{
        bitindex := k.NodeID.Xor(req.Key).PrefixLen()-1
		if bitindex<0{
			bitindex=0
		}
		tempFoundNode:=new(FoundNode)
        FoundNodelst := make([]FoundNode, K)
        for i:=0;i<len(k.AddrTab[bitindex].ContactLst);i++{
			if k.AddrTab[bitindex].ContactLst[i].Host!=nil&&k.AddrTab[bitindex].ContactLst[i].NodeID!=req.Sender.NodeID{
				tempFoundNode.NodeID=k.AddrTab[bitindex].ContactLst[i].NodeID
				tempFoundNode.Port=k.AddrTab[bitindex].ContactLst[i].Port
				tempFoundNode.IPAddr=k.AddrTab[bitindex].ContactLst[i].Host.String()
				FoundNodelst[i]=*tempFoundNode
			}
        }
        res.Nodes=FoundNodelst
    }
	
    res.MsgID=CopyID(req.MsgID)
    res.Err=nil
    return nil
}

/////////////
func (k *Kademlia) FindValue2(req FindValueRequest, res *FindValueResult) error {
    // TODO: Implement.
    //go Update(k, req.Sender)
	k.ch<-req.Sender
	go Update2(k)
	//fmt.Println("Finding....")
	
	found, val:=Local_Find_Value(k,req.Key)
	if found==true{
	    res.Value=val
		webpageDSroot:="./webpageDS/"
		filename:=webpageDSroot+string(val)+".html"
		body, err:=ioutil.ReadFile(filename)
		if err!=nil{
			return err
		}
		res.Body=string(body)
		//fmt.Println("Find value: ", string(val))
	}else if found == false{
        bitindex := k.NodeID.Xor(req.Key).PrefixLen()-1
		if bitindex<0{
			bitindex=0
		}
		tempFoundNode:=new(FoundNode)
        FoundNodelst := make([]FoundNode, K)
        for i:=0;i<len(k.AddrTab[bitindex].ContactLst);i++{
			if k.AddrTab[bitindex].ContactLst[i].Host!=nil&&k.AddrTab[bitindex].ContactLst[i].NodeID!=req.Sender.NodeID{
				tempFoundNode.NodeID=k.AddrTab[bitindex].ContactLst[i].NodeID
				tempFoundNode.Port=k.AddrTab[bitindex].ContactLst[i].Port
				tempFoundNode.IPAddr=k.AddrTab[bitindex].ContactLst[i].Host.String()
				FoundNodelst[i]=*tempFoundNode
			}
        }
        res.Nodes=FoundNodelst
    }
	
    res.MsgID=CopyID(req.MsgID)
    res.Err=nil
    return nil
}


