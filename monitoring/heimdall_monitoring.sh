#!/bin/bash

#############################################################################
#									    #
#      This script does 						    #
#            1: Enable prometheus for Heimdall.                             #
#            2: setting up Prometheus and Grafana locally.                  #
#            3: Create Grafana Dashboard for Heimdall.		            #
#									    #
#############################################################################

## path of Grafana datasorce and dashboard json's
json_path=$(pwd)
echo $json_path

## Enabling Prometheus for Heimdall
sudo sed -i 's/prometheus = false/prometheus = true/' /home/ubuntu/.heimdalld/config/config.toml
## Restaring Heimdalld
sudo service heimdalld restart

## Installing Prometheus to fetch the Heimdall metrics 

sudo apt-get update -y
sudo useradd --no-create-home --shell /bin/false prometheus
sudo mkdir /etc/prometheus
sudo mkdir /var/lib/prometheus
sudo chown prometheus:prometheus /etc/prometheus
sudo chown prometheus:prometheus /var/lib/prometheus
cd ~
curl -LO https://github.com/prometheus/prometheus/releases/download/v2.16.0/prometheus-2.16.0.linux-amd64.tar.gz
tar xvf prometheus-2.16.0.linux-amd64.tar.gz
sudo cp prometheus-2.16.0.linux-amd64/prometheus /usr/local/bin/
sudo cp prometheus-2.16.0.linux-amd64/promtool /usr/local/bin/
sudo chown prometheus:prometheus /usr/local/bin/prometheus
sudo chown prometheus:prometheus /usr/local/bin/promtool
sudo cp -r prometheus-2.16.0.linux-amd64/consoles /etc/prometheus
sudo cp -r prometheus-2.16.0.linux-amd64/console_libraries /etc/prometheus
sudo chown -R prometheus:prometheus /etc/prometheus/consoles
sudo chown -R prometheus:prometheus /etc/prometheus/console_libraries
rm -rf prometheus-2.16.0.linux-amd64.tar.gz prometheus-2.16.0.linux-amd64

## Adding localhost to prometheus ##
cat <<EOF >/etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:26660']
EOF

sudo chown prometheus:prometheus /etc/prometheus/prometheus.yml

## Configuring Prometheus as systemd service
cat <<EOF >/etc/systemd/system/prometheus.service
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl start prometheus
sudo systemctl enable prometheus


### Configuring Grafana locally
echo "deb https://packages.grafana.com/oss/deb stable main" | sudo tee -a /etc/apt/sources.list.d/grafana.list
sudo apt-get install -y adduser libfontconfig1
wget https://dl.grafana.com/oss/release/grafana_6.6.2_amd64.deb
sudo dpkg -i grafana_6.6.2_amd64.deb
sudo systemctl daemon-reload
sudo systemctl start grafana-server
sudo systemctl enable grafana-server.service

#Adding prometheus DB as Datasourse for grafana
curl -k -X POST  http://admin:admin@localhost:3000/api/datasources -H 'Content-Type: application/json' -H 'Accept: application/json' -d "$(cat $json_path/datasource.json)"


#creating grafana dashboard for heimdall
curl -k -X POST http://admin:admin@localhost:3000/api/dashboards/db -H 'Content-Type: application/json' -H 'Accept: application/json' -d "$(cat $json_path/dashboard.json)"
