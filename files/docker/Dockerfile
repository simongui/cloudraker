FROM mysql:5.7
#FROM mysql/mysql-server:5.7

RUN apt-get update && apt-get install -y wget nano git tar make gcc automake autoconf libtool openssh-server sysbench && rm -rf /var/lib/apt/lists/*

COPY my.cnf /etc/mysql/my.cnf
COPY id_rsa /etc/ssh/ssh_host_rsa_key
COPY id_rsa.pub /etc/ssh/ssh_host_rsa_key.pub
COPY ssh_config /etc/ssh/ssh_config
COPY sshd_config /etc/ssh/sshd_config
RUN chmod -R 400 /etc/ssh
COPY bench.sh /root
RUN chmod +x /root/bench.sh

# root setup
RUN mkdir /root/.ssh
COPY id_rsa.pub /root/.ssh/authorized_keys
COPY id_rsa /root/.ssh/id_rsa
COPY id_rsa.pub /root/.ssh/id_rsa.pub
RUN chown -R root.root /root/.ssh
RUN chmod 400 /root/.ssh/*
