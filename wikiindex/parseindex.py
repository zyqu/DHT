#!/bin/python

import json

lines=open('enwiki-20130503-pages-articles-multistream-index.txt','r').readlines()

numofindex=90000
indexlst=[]

for l in lines:
	l=l.strip()
	l=l[l.find(':')+1:]
	l=l[l.find(':')+1:]
	l=l.replace(' ','_')
	if l not in indexlst:
		indexlst.append(l)

	if len(indexlst)>numofindex-1:
		break

open('%s'%numofindex+'wikiindex.json','w').write(json.dumps(indexlst))






