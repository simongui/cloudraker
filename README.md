# Cloudraker
Cloudraker is a MySQL cluster management tool for macOS using Docker for Mac. Cloudraker makes it easy to work with a cluster of MySQL instances that can replicate between each other.

Cloudraker is something I created out of necessity while working on MySQL clusters in production. I wanted a local development environment that could closely match what was in production for shards, datacenters and replication topologies that was easy to create and tear down quickly.

# Features
- Fast. You can add and remove 5 MySQL instances that replicate between each other in less than `30 seconds`.
- Creates production like hostnames like `db1.us-east-1.com` that resolve to internal IP addresses.
- Adds `sshd` to each container for SSH access.
- Creates MySQL replication grants.

# TODO
- Replication
- Proxies local MySQL and SSH ports to the hostnames (`db1.us-east-1.com:3306`, `db2.us-east1.com:3306`) that routes to the proper container for production like testing.
- Failure injection.

# Usage
```
$ sudo ./cloudraker addbatch \
                    --nodes=5 \
                    --cluster=cluster \
                    --datacenter=us-east-1 \
                    --hostformat=db%d.us-east-1.com \
                    --subnet=10.0.0.
                    
 41 / 41 [====================================================================================] 100.00% 22s

$ cloudraker ls

  CLUSTER |       HOST        |    DC     |    IP
+---------+-------------------+-----------+----------+
  cluster | db3.us-east-1.com | us-east-1 | 10.0.0.3
  cluster | db1.us-east-1.com | us-east-1 | 10.0.0.1
  cluster | db2.us-east-1.com | us-east-1 | 10.0.0.2
  cluster | db5.us-east-1.com | us-east-1 | 10.0.0.5
  cluster | db4.us-east-1.com | us-east-1 | 10.0.0.4
+---------+-------------------+-----------+----------+
                                  NODES   |    5
                              +-----------+----------+
                              
$ sudo ./cloudraker removebatch --nodes=5 --hostformat=db%d.us-east-1.com
21 / 21 [====================================================================================]  100.00% 5s

$ cloudraker ls

  CLUSTER | HOST |  DC   | IP
+---------+------+-------+----+
+---------+------+-------+----+
                   NODES | 0
                 +-------+----+
```

# Screencast
<img src="cloudraker.gif"/>
