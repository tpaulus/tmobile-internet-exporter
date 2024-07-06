# TMobile Home Internet Gateway Exporter
A simple Prometheus exporter for T-Mobile Home Internet Gateways.
Tested on the following Models:
* Sagemcom FAST5688W

All of the metrics and device details are queried from:
`http://<gateway_ip>/TMI/v1/gateway?get=all`
