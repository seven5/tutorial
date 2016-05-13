#!/bin/sh

# This script is run by Supervisor to start PostgreSQL 9.5 in foreground mode

if [ -d /var/run/postgresql ]; then
    chmod 2775 /var/run/postgresql
else
    install -d -m 2775 -o postgres -g postgres /var/run/postgresql
fi

exec su postgres -c "/usr/lib/postgresql/9.5/bin/postgres -p 5433 -D /var/lib/postgresql/9.5/main -c config_file=/etc/postgresql/9.5/main/postgresql.conf"
