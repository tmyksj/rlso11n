# rlso11n: rootless-orchestration
rlso11n does:
- manages docker daemons.
- constructs/destructs an cluster (using docker swarm mode).

note: any releases under version 1.0.0 may contain breaking changes.

## usage
rlso11n uses two terminals (or background jobs).

exports an ip address list of hosts:
```
$ # terminal 1
$ export RLSO11N_HOST_LIST="172.16.0.1,172.16.0.2"
```

starts rlso11n servers:
```
$ # terminal 1
$ rlso11n bg/start
```

starts docker daemons using rootless mode:
```
$ # terminal 2
$ rlso11n up/dockerd
```

constructs a swarm cluster:
```
$ # terminal 2
$ rlso11n up/swarm
```

to exit rlso11n, send ctrl+c at terminal 1:
```
^C
$
```

## installation
```
$ make
$ make install
$ make install_dependencies
```

## settings
rlso11n uses environment variables for settings:
```
RLSO11N_DIR              use this directory to run rootless containers.
RLSO11N_HOST_LIST        use these hosts to run rlso11n.
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
