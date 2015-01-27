#gopherjs does the same dependency checking as go tool, so there is no harm
#in just building things with it everytime (its also fast)
.PHONY: ../static/en/web/signup.js

##
## WARNING: This makefile does not have dependencies on libraries used by the 
## WARNING: go packages in pages such as seven5.  If you change those libraries
## WARNING: you have to do "make clean; make" to rebuild and include the new
## WARNING: library versions.
##

all:  ../static/en/web/index.html \
	../static/en/web/signup.html \
	../static/en/web/signup.js 

clean:
	rm -rf \
	 ../static/en/web/index.html \
	../static/en/web/signup.html \
	../static/fixed/en/web/signup.js*

##
## CLIENT SIDE
## 
../static/en/web/signup.js:
	gopherjs build -o  ../static/en/web/signup.js ../client/signup.go

##
## HTML PAGES
## 
../static/en/web/index.html: template/index.html template/index.json template/support/*.tmpl
	pagegen --support=support --dir=template --start=index.html --json=index.json > ../static/en/web/index.html

../static/en/web/signup.html: template/signup.html template/signup.json template/support/*.tmpl
	pagegen --support=support --dir=template --start=signup.html --json=signup.json > ../static/en/web/signup.html
