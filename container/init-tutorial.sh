#
# create the tutorial
#
#bootstrap needs gb and gopherjs
cd /seven5/bootstrap
export GOPATH=/seven5/bootstrap
go get github.com/constabulary/gb/...
go get github.com/gopherjs/gopherjs

# no we are ready to get the source code
unset GOPATH
mkdir -p /seven5/tutorial/src/github.com/seven5
cd /seven5/tutorial/src/github.com/seven5
git clone -b tutorial https://github.com/seven5/tutorial.git
git clone https://github.com/seven5/seven5.git
git clone https://github.com/seven5/gb-seven5.git

#configuration
cp tutorial/manifest /seven5/tutorial/vendor/manifest
mv /enable-tutorial.sh /seven5/enable-tutorial

#inflate deps
cd /seven5/tutorial
gb vendor restore

#build source code
gb build github.com/seven5/tutorial/...

#create and init db
#because this is ONlY used at container build time, we let postgres run on 5432
/etc/init.d/postgresql start
sleep 5
su postgres bash -c "psql -c \"CREATE USER root WITH PASSWORD '';\""
su postgres -c "createdb -O root fresno"
export DATABASE_URL="postgres://root@localhost:5432/fresno"
migrate --up
/etc/init.d/postgresql stop
