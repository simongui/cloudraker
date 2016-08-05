# Cloudraker
Cloudraker is a MySQL database docker cluster management tool for macOS using Docker for Mac. Cloudraker makes it easy to work with a cluster of MySQL instances that can replicate between each other.

Cloudraker is something I created out of necessity while working on MySQL clusters in production. I wanted a local development environment that could closely match what was in production for shards and datacenters and replication topologies but easy to create and tear down quickly.

# Features
- Fast. You can add and remove 5 MySQL instances that replicate between each other in less than `30 seconds`.
- Creates production like hostnames like `db1.us-east-1.com` that resolve to internal IP addresses.

# Future
- Proxies local MySQL and SSH ports to the hostnames (`db1.us-east-1.com:3306`, `db2.us-east1.com:3306`) that routes to the proper container for production like testing.
- Failure injection.

<img src="cloudraker.gif"/>
