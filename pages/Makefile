#gopherjs does the same dependency checking as go tool, so there is no harm
#in just building things with it everytime (it's also fast)
.PHONY: ../static/en/web/signup.js
.PHONY: ../static/en/web/login.js
.PHONY: ../static/en/web/post/index.js
.PHONY: ../static/en/web/post/new.js
.PHONY: ../static/en/web/post/edit.js
.PHONY: ../static/en/web/post/view.js

GOPHERJSOPTS:=

##
## WARNING: This makefile does not have dependencies on libraries used by the
## WARNING: go packages in pages such as seven5.  If you change those libraries
## WARNING: you have to do "make clean; make" to rebuild and include the new
## WARNING: library versions.  This is rare.
##

PAGES:=../static/en/web/signup.html ../static/en/web/login.html ../static/en/web/post/index.html \
	../static/en/web/post/new.html ../static/en/web/post/edit.html ../static/en/web/post/view.html

JS:=../static/en/web/signup.js ../static/en/web/login.js ../static/en/web/post/index.js \
	../static/en/web/post/new.js ../static/en/web/post/edit.js ../static/en/web/post/view.js

JSMAPS:=../static/en/web/signup.js.map ../static/en/web/login.js.map ../static/en/web/post/index.js.map \
	../static/en/web/post/new.js.map ../static/en/web/post/edit.js.map ../static/en/web/post/view.js.map

all: $(PAGES) $(JS) $(JSMAPS)

clean:
	rm -rf $(PAGES) $(JS) $(JSMAPS)

html: $(PAGES)

js: $(JS)


##
## CLIENT SIDE
##
../static/en/web/signup.js:
	gopherjs build $(GOPHERJSOPTS) -o  ../static/en/web/signup.js ../client/signup.go

../static/en/web/login.js:
	gopherjs build $(GOPHERJSOPTS) -o  ../static/en/web/login.js ../client/login.go

../static/en/web/post/index.js:
	gopherjs build $(GOPHERJSOPTS) -o  ../static/en/web/post/index.js ../client/post/index.go

../static/en/web/post/new.js:
	gopherjs build $(GOPHERJSOPTS) -o  ../static/en/web/post/new.js ../client/post/new.go

../static/en/web/post/edit.js:
	gopherjs build $(GOPHERJSOPTS) -o  ../static/en/web/post/edit.js ../client/post/edit.go

../static/en/web/post/view.js:
	gopherjs build $(GOPHERJSOPTS) -o  ../static/en/web/post/view.js ../client/post/view.go


##
## HTML PAGES
##
../static/en/web/signup.html: template/signup.html template/signup.json template/support/*.tmpl
	pagegen --support=support --dir=template --start=signup.html --json=signup.json > ../static/en/web/signup.html

../static/en/web/login.html: template/login.html template/login.json template/support/*.tmpl
	pagegen --support=support --dir=template --start=login.html --json=login.json > ../static/en/web/login.html

../static/en/web/post/index.html: template/post/index.html template/post/index.json template/support/*.tmpl
	pagegen --support=support --dir=template --start=post/index.html --json=post/index.json > ../static/en/web/post/index.html

../static/en/web/post/new.html: template/post/new.html template/post/new.json template/support/*.tmpl
	pagegen --support=support --dir=template --start=post/new.html --json=post/new.json > ../static/en/web/post/new.html

../static/en/web/post/edit.html: template/post/edit.html template/post/edit.json template/support/*.tmpl
	pagegen --support=support --dir=template --start=post/edit.html --json=post/edit.json > ../static/en/web/post/edit.html

../static/en/web/post/view.html: template/post/view.html template/post/view.json template/support/*.tmpl
	pagegen --support=support --dir=template --start=post/view.html --json=post/view.json > ../static/en/web/post/view.html
