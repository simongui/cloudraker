#!/bin/bash
mysql -u root -ppassword -e "DROP DATABASE sysbench;"
mysql -u root -ppassword -e "CREATE DATABASE sysbench;"
sysbench --test=oltp --db-driver=mysql --oltp-table-size=$1 --mysql-db=sysbench --mysql-user=root --mysql-password=password prepare
sysbench --test=oltp --db-driver=mysql --oltp-table-size=$1 --oltp-test-mode=complex --oltp-read-only=off --num-threads=6 --max-time=60 --max-requests=0 --mysql-db=sysbench --mysql-user=root --mysql-password=password run
