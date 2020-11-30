# rlso11n: rootless-orchestration
rlso11n does:
- manages docker daemons.
- constructs/destructs a cluster (using docker swarm mode).

note: any releases under version 1.0.0 may contain breaking changes.

## usage
rlso11n needs two terminals (or background jobs).

starts rlso11n servers:
```
$ # terminal 1
$ rlso11n serve --hosts node-0,node-1
```

starts docker daemons using rootless mode:
```
$ # terminal 2
$ rlso11n up dockerd
```

constructs a swarm cluster:
```
$ # terminal 2
$ rlso11n up swarm
```

executes commands:
```
$ # terminal 2
$
$ # execute `hostname` at all nodes
$ rlso11n exec --nodes all hostname
$
$ # execute `hostname` at 0, 2-7 nodes
$ rlso11n exec --nodes 0,2-7 hostname
$
$ # execute `docker ps -a` at all nodes
$ # --docker option appends -H option automatically
$ rlso11n exec --nodes all --docker ps -a
$
$ # \$i is replaced to node index
$ rlso11n exec --nodes all echo \$i
```

to exit rlso11n, send ctrl+c at terminal 1:
```
^C
$
```

## settings
rlso11n uses command options passed to serve command for settings.
available options are:
```
--dir /tmp/rlso11n              use this directory to run rootless containers
--hosts node-0,node-1           use these hosts to run rlso11n
```

## installation
```
$ make
$ make install
$ make install_dependencies
```

## administrator's settings
### for centos
installs glib2-devel, libcap-devel, libseccomp-devel for slirp4netns:
```
# yum install -y glib2-devel libcap-devel libseccomp-devel
```

sets subuid, subgid for user(s):
```
# echo username:100000:65536 >> /etc/subuid
# echo username:100000:65536 >> /etc/subgid
```

sets max_user_namespaces:
```
# echo 28633 > /proc/sys/user/max_user_namespaces
```
