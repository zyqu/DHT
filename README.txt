************************************************
* Distributed Wikipedia Command Line inferface *
************************************************

Make sure to configure the maximal number of allowed to be opened files in your OS before conduct the large-scale test for multi-user mode. 4000 might be a good choice

Command line:

single-user mode, where one node fetching webpages:
batchtest filenameofwikiindex 0 

multi-user mode, where multiple nodes fetching/sharing webpages:
batchtest filenameofwikiindex 1
startindex endindex

if the endindex is less than startindex, then there is a loop. end index -> 0 -> startindex

filenameofwikiindex is the jsonfile which has the indexes of wiki webpage to fetch "/wikiindex/filenameofwikiindex"

In the multi-user mode, the way to test is as follows:

First:
Device A ping Device B and Device C
Device B ping Device A and Device C
Device C ping Device A and Device B

Second:
Device A 
batchtest 30000wikiindex.json 1
0 2999

Device B
batchtest 30000wikiindex.json 1
1000 999

Device C
batchtest 30000wikiindex.json 1
2000 1999



************
* BUILDING *
************

Go's build tools depend on the value of the GOPATH environment variable. $GOPATH
should be the project root: the absolute path of the directory containing
{bin,pkg,src}.

Once you've set that, you should be able to build the skeleton and create an
executable at bin/main with:

    go install main

Running main as

    main localhost:7890 localhost:7890

will cause it to start up a server bound to localhost:7890 (the first argument)
and then connect as a client to itself (the second argument). All it does by
default is perform a PING RPC and exit.



**************************
* COMMAND-LINE INTERFACE *
**************************

whoami
    Print your node ID.

local_find_value key
    If your node has data for the given key, print it.
    If your node does not have data for the given key, you should print "ERR".

get_contact ID
    If your buckets contain a node with the given ID,
        printf("%v %v\n", theNode.addr, theNode.port)
    If your buckers do not contain any such node, print "ERR".

iterativeStore key value
    Perform the iterativeStore operation and then print the ID of the node that
    received the final STORE operation.

iterativeFindNode ID
    Print a list of ≤ k closest nodes and print their IDs. You should collect
    the IDs in a slice and print that.

iterativeFindValue key
    printf("%v %v\n", ID, value), where ID refers to the node that finally
    returned the value. If you do not find a value, print "ERR".

ping nodeID
ping host:port
    Perform a ping. 

store nodeID key value 
    Perform a store and print a blank line.

find_node nodeID key
    Perform a find_node and print its results as for iterativeFindNode.

find_value nodeID key
    Perform a find_value. If it returns nodes, print them as for find_node. If
    it returns a value, print the value as in iterativeFindValue.

