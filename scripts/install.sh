#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installMyMetadata() {
    curl -Lo- https://github.com/sunshineplan/metadata/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/metadata* /var/www/metadata
    cd /var/www/metadata
    go build -ldflags "-s -w"
    ./metadata install
}

configMyMetadata() {
    read -p 'Please enter metadata database server address: ' dbserver
    read -p 'Please enter metadata server port: ' dbport
    read -p 'Please enter metadata database name: ' database
    read -p 'Please enter metadata collection name: ' collection
    read -p 'Please enter metadata server username: ' username
    read -sp 'Please enter metadata server password: ' password
    read -p 'Please enter unix socket(default: /run/metadata.sock): ' unix
    [ -z $unix ] && unix=/run/metadata.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/metadata.log): ' log
    [ -z $log ] && log=/var/log/app/metadata.log
    mkdir -p $(dirname $log)
    sed "s/\$dbserver/$dbserver/" /var/www/metadata/config.ini.default > /var/www/metadata/config.ini
    sed -i "s/\$dbport/$dbport/" /var/www/metadata/config.ini
    sed -i "s/\$database/$database/" /var/www/metadata/config.ini
    sed -i "s/\$collection/$collection/" /var/www/metadata/config.ini
    sed -i "s/\$username/$username/" /var/www/metadata/config.ini
    sed -i "s/\$password/$password/" /var/www/metadata/config.ini
    sed -i "s,\$unix,$unix," /var/www/metadata/config.ini
    sed -i "s,\$log,$log," /var/www/metadata/config.ini
    sed -i "s/\$host/$host/" /var/www/metadata/config.ini
    sed -i "s/\$port/$port/" /var/www/metadata/config.ini
    service metadata start
}

writeLogrotateScrip() {
    if [ ! -f '/etc/logrotate.d/app' ]; then
	cat >/etc/logrotate.d/app <<-EOF
		/var/log/app/*.log {
		    copytruncate
		    rotate 12
		    compress
		    delaycompress
		    missingok
		    notifempty
		}
		EOF
    fi
}

createCronTask() {
    cp -s /var/www/metadata/scripts/metadata.cron /etc/cron.monthly/metadata
    chmod +x /var/www/metadata/scripts/metadata.cron
}

setupNGINX() {
    cp -s /var/www/metadata/scripts/metadata.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/metadata/scripts/metadata.conf
    sed -i "s,\$unix,$unix," /var/www/metadata/scripts/metadata.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installMyMetadata
    configMyMetadata
    writeLogrotateScrip
    createCronTask
    setupNGINX
}

main