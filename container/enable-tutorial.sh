## setup the TUTROOT directory, you may want to make this an absolute path
## because this will blow up if you source this directory from somewhere
## other than "."
export TUTROOT=/seven5/tutorial

## setup the gopath correctly. the second element is needed to make gopherjs
## work correctly.  Note that this is primarily used to locate the source
## code when debugging.
export GOPATH=$TUTROOT:$TUTROOT/vendor

## this should point to your go install if it is not installed in a
## a system-wide location (ala via brew)
export GOROOT=/usr/lib/go-1.6/

## the path needs to include go1.6+
##
## you should put the exact path you want to use in here, don't use
## PATH=$PATH:...  because that will bite you later.

##if you set GOROOT
##export PATH=$GOROOT/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:$TUTROOT/bin
#otherwise
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/go-1.6/bin:/seven5/bootstrap/bin:$TUTROOT/bin

##
## ENVIRONMENT FOR LOCAL DEVELOPMENT
##
export FRESNO_TEST=true  ##allows you to differentiate running on heroku or not
export HEROKU_NAME=damp-sierra-7161 ##put your heroku app name
export PORT=5000

### It is actually a good idea to run on a different port (not 5432) to
### prevent your application from interfering with other things that expect
### postgres to be on the standard port.
export PGPORT=5433
##assumes USER is set correctly, and probably will break if it has spaces
export DATABASE_URL="postgres://root@localhost:$PGPORT/fresno"

##assumes USER is set correctly, and probably will break if it has spaces
export PGUSER="$USER"
export PGSSLMODE=disable

##for use in testing only, this a weak key!
#look into
# $TUTROOT/bin/key2hex foo
# and then follow the instructions
export SERVER_SESSION_KEY='777936327336633637393632726b4476'

##where is the static content we should be serving up? note that this is relative
##to the location where you run fresno, so you may want to make this an absolute
##path
export STATIC_DIR=src/github.com/seven5/tutorial/static
