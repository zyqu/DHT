package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
    "net"
    "fmt"
    )


// Host identification.
type Contact struct {
    NodeID ID
    Host net.IP
    Port uint16
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
    go Update(k, ping.Sender)
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
}

type StoreResult struct {
    MsgID ID
    Err error
}

func (k *Kademlia) Store(req StoreRequest, res *StoreResult) error {
    // TODO: Implement.
    go Update(k, req.Sender)
    k.Localmap[req.Key]=req.Value
    fmt.Println("\n")
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
    Update(k, req.Sender)

    bitindex := k.NodeID.Xor(req.NodeID).PrefixLen()
    tempFoundNode:=new(FoundNode)
    FoundNodelst := make([]FoundNode, K)
    for i:=0;i<len(k.AddrTab[bitindex].ContactLst);i++{
        tempFoundNode.NodeID=k.AddrTab[bitindex].ContactLst[i].NodeID
        tempFoundNode.Port=k.AddrTab[bitindex].ContactLst[i].Port
        tempFoundNode.IPAddr=k.AddrTab[bitindex].ContactLst[i].Host.String()
        FoundNodelst[i]=*tempFoundNode
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
}

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
    // TODO: Implement.
    Update(k, req.Sender)
    found, val:=Local_Find_Value(k,req.Key)
    res.Value=val
    if found == false{
        bitindex := k.NodeID.Xor(req.Key).PrefixLen()
        tempFoundNode:=new(FoundNode)
        FoundNodelst := make([]FoundNode, K)
        for i:=0;i<len(k.AddrTab[bitindex].ContactLst);i++{
            tempFoundNode.NodeID=k.AddrTab[bitindex].ContactLst[i].NodeID
            tempFoundNode.Port=k.AddrTab[bitindex].ContactLst[i].Port
            tempFoundNode.IPAddr=k.AddrTab[bitindex].ContactLst[i].Host.String()
            FoundNodelst[i]=*tempFoundNode
        }
        res.Nodes=FoundNodelst
    }

    res.MsgID=CopyID(req.MsgID)
    res.Err=nil
    return nil
}

