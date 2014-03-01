from bs4 import BeautifulSoup
import sys
filename=sys.argv[1]
file=open(filename, 'r')
html_doc=str()
for line in file:
    html_doc+=line
soup=BeautifulSoup(html_doc)
#print soup.prettify()
for link in soup.find_all('link'):
    if link['href'][0]!='#' and link['href'][:2]!='//' and link['href'][:4]!='/www' and link['href'][:4]!='http':
        href='http://en.wikipedia.org'
        href+=link['href']
        link['href']=href
for link in soup.find_all('a'):
    if link.get('href')!=None and link['href'][0]!='#' and link['href'][:2]!='//' and link['href'][:3]!='www' and link['href'][:4]!='http':
        href='http://en.wikipedia.org'
        href+=link['href']
        link['href']=href
html=soup.prettify('utf-8')
file.close()
with open(filename, 'wb') as file:
    file.write(html)