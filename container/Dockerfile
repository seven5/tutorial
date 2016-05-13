FROM ubuntu:16.04
MAINTAINER iansmith@acm.org
RUN apt-get update && apt-get install -y openssh-server git postgresql-9.5 golang-1.6 supervisor

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY postgres.sh /start-postgres.sh

#no, as a matter of fact, I don't know why this is necessary
RUN mkdir -p /var/run/postgresql/9.5-main.pg_stat_tmp/
RUN chown postgres /var/run/postgresql/9.5-main.pg_stat_tmp/

#we are not starting sshd in the usual way, so we have to do this by hand
RUN mkdir /var/run/sshd
RUN chmod 0755 /var/run/sshd

##for ssh
RUN echo "root:ubuntu" | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

#tutorial setup
RUN mkdir -p /seven5/tutorial/src
RUN mkdir -p /seven5/tutorial/vendor/src
RUN mkdir -p /seven5/bootstrap/src

COPY init-tutorial.sh /init-tutorial.sh
COPY enable-tutorial.sh /enable-tutorial.sh

## these are here largely for convenience if you grab a shell on the container
## you can get this same effect with 'source /seven5/tutorial/enable-tutorial'
ENV GOROOT=/usr/lib/go-1.6/
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/go-1.6/bin:/seven5/bootstrap/bin:/seven5/tutorial/bin
ENV FRESNO_TEST=true
ENV PORT=80
ENV PGPORT=5433
ENV DATABASE_URL="postgres://root@localhost:5433/fresno"
ENV PGUSER=root
ENV PGSSLMODE=disable
ENV SERVER_SESSION_KEY='777936327336633637393632726b4476'
ENV STATIC_DIR=/seven5/tutorial/src/github.com/seven5/tutorial/static

RUN sh /init-tutorial.sh

EXPOSE 22 80
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
