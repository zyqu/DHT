from bs4 import BeautifulSoup
import sys

html_doc=sys.argv[1]
soup=BeautifulSoup(html_doc)

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
#html=soup.prettify('utf-8')
print soup