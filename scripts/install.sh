#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installMyMetadata() {
    curl -Lo- https://github.com/sunshineplan/mymetadata-go/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/mymetadata-go* /var/www/mymetadata-go
    cd /var/www/mymetadata-go
    go build
}

configMyMetadata() {
    read -p 'Please enter metadata database server address: ' dbserver
    read -p 'Please enter metadata server port: ' dbport
    read -p 'Please enter metadata database name: ' database
    read -p 'Please enter metadata collection name: ' collection
    read -p 'Please enter metadata server username: ' username
    read -p 'Please enter metadata server password: ' password
    read -p 'Please enter unix socket(default: /run/mymetadata-go.sock): ' unix
    [ -z $unix ] && unix=/run/mymetadata-go.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/mymetadata-go.log): ' log
    [ -z $log ] && log=/var/log/app/mymetadata-go.log
    mkdir -p $(dirname $log)
    sed "s/\$dbserver/$dbserver/" /var/www/mymetadata-go/config.ini.default > /var/www/mymetadata-go/config.ini
    sed -i "s/\$dbport/$dbport/" /var/www/mymetadata-go/config.ini
    sed -i "s/\$database/$database/" /var/www/mymetadata-go/config.ini
    sed -i "s/\$collection/$collection/" /var/www/mymetadata-go/config.ini
    sed -i "s/\$username/$username/" /var/www/mymetadata-go/config.ini
    sed -i "s/\$password/$password/" /var/www/mymetadata-go/config.ini
    sed -i "s,\$unix,$unix," /var/www/mymetadata-go/config.ini
    sed -i "s,\$log,$log," /var/www/mymetadata-go/config.ini
    sed -i "s/\$host/$host/" /var/www/mymetadata-go/config.ini
    sed -i "s/\$port/$port/" /var/www/mymetadata-go/config.ini
}

setupsystemd() {
    cp -s /var/www/mymetadata-go/scripts/mymetadata-go.service /etc/systemd/system
    systemctl enable mymetadata-go
    service mymetadata-go start
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
    cp -s /var/www/mymetadata-go/scripts/mymetadata-go.cron /etc/cron.monthly/mymetadata-go
    chmod +x /var/www/mymetadata-go/scripts/mymetadata-go.cron
}

setupNGINX() {
    cp -s /var/www/mymetadata-go/scripts/mymetadata-go.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/mymetadata-go/scripts/mymetadata-go.conf
    sed -i "s,\$unix,$unix," /var/www/mymetadata-go/scripts/mymetadata-go.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installMyMetadata
    configMyMetadata
    setupsystemd
    writeLogrotateScrip
    createCronTask
    setupNGINX
}

main