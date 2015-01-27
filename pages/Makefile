##
## WARNING: This makefile does not have dependencies on libraries used by the 
## WARNING: go packages in pages such as seven5.  If you change those libraries
## WARNING: you have to do "make clean; make" to rebuild and include the new
## WARNING: library versions.
##

all:  ../static/en/web/index.html



clean:
	rm -rf \
	 ../static/en/web/index.html 

##
## HTML PAGES
## 
## These are done this way largely because it's hard to uniquely identify a page
## from only the basename.  This does NOT encode dependencies among the templates so
## you have to do make clean;make if you change the dependent files.
##

../static/en/web/index.html: template/index.html template/index.json
	pagegen --support=support --dir=template --start=index.html --json=index.json > ../static/en/web/index.html
