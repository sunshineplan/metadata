#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
}

installMetadata() {
    mkdir -p /var/www/metadata
    curl -Lo- https://github.com/sunshineplan/metadata/releases/latest/download/release.tar.gz | tar zxC /var/www/metadata
    cd /var/www/metadata
    chmod +x metadata
}

configMetadata() {
    read -p 'Please enter metadata data source: ' source
    read -p 'Please enter metadata database (default: metadata): ' database
    [ -z $database ] && database=metadata
    read -p 'Please enter metadata collection (default: metadata): ' collection
    [ -z $collection ] && collection=metadata
    read -p 'Please enter app id: ' id
    read -p 'Please enter api key: ' key
    read -p 'Please enter unix socket (default: /run/metadata.sock): ' unix
    [ -z $unix ] && unix=/run/metadata.sock
    read -p 'Please enter host (default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port (default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path (default: /var/log/app/metadata.log): ' log
    [ -z $log ] && log=/var/log/app/metadata.log
    read -p 'Please enter update URL: ' update
    read -p 'Please enter exclude files: ' exclude
    mkdir -p $(dirname $log)
    sed "s/\$source/$source/" /var/www/metadata/config.ini.default > /var/www/metadata/config.ini
    sed -i "s/\$database/$database/" /var/www/metadata/config.ini
    sed -i "s/\$collection/$collection/" /var/www/metadata/config.ini
    sed -i "s/\$id/$id/" /var/www/metadata/config.ini
    sed -i "s/\$key/$key/" /var/www/metadata/config.ini
    sed -i "s,\$unix,$unix," /var/www/metadata/config.ini
    sed -i "s,\$log,$log," /var/www/metadata/config.ini
    sed -i "s/\$host/$host/" /var/www/metadata/config.ini
    sed -i "s/\$port/$port/" /var/www/metadata/config.ini
    sed -i "s,\$update,$update," /var/www/metadata/config.ini
    sed -i "s|\$exclude|$exclude|" /var/www/metadata/config.ini
    ./metadata install || exit 1
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

setupNGINX() {
    cp -s /var/www/metadata/scripts/metadata.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/metadata/scripts/metadata.conf
    sed -i "s,\$unix,$unix," /var/www/metadata/scripts/metadata.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installMetadata
    configMetadata
    writeLogrotateScrip
    setupNGINX
}

main