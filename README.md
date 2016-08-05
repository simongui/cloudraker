# Cloudraker
Cloudraker is a MySQL cluster management tool for macOS using Docker for Mac. Cloudraker makes it easy to work with a cluster of MySQL instances that can replicate between each other.

Cloudraker is something I created out of necessity while working on MySQL clusters in development. I wanted a local development environment that could be created quickly and closely match what was in production for shards, datacenters and replication topologies.

# Features
- Fast. You can add and remove 5 MySQL instances that replicate between each other in less than `30 seconds`.
- Creates local IP aliases on `lo0` on the specified subnet.
- Creates production like hostnames like `db1.us-east-1.com` that resolve to the local IP address aliases created on `lo0`.
- Adds `sshd` to each container for SSH access.
- Creates MySQL replication grants.

# TODO
- Replication
- Proxies local MySQL and SSH ports to the hostnames (`db1.us-east-1.com:3306`, `db2.us-east1.com:3306`) that routes to the proper container for production like testing.
- Failure injection.
- With SSH access to each MySQL instance already setup, the `MHA` failover tool could be supported.

# Usage
```
usage: cloudraker [<flags>] <command> [<args> ...]

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  ls [<flags>]
    List nodes in a specific or all clusters.

  add [<flags>]
    Add a node to the specified cluster and datacenter.

  addbatch [<flags>]
    Add a batch of nodes to the specified cluster and datacenter.

  remove [<flags>] [<host>]
    Remove a node from the cluster.

  removebatch [<flags>]
    Add a batch of nodes to the specified cluster and datacenter.
```

#### Adding and removing a single node
```
$ sudo ./cloudraker add --cluster=cluster --datacenter=us-east-1 --host=db6.us-east-1.com --ipaddress=10.0.0.6
 8 / 8 [====================================================================================] 100.00% 12s
MySQL running
	host: db6.us-east-1.com
	id: 6
	read_only: true
	
$ sudo ./cloudraker remove db6.us-east-1.com
 4 / 4 [====================================================================================] 100.00% 1s
```

#### Adding and removing batches of nodes
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
