# rootless-orchestration
rootless-orchestration does:
- manages docker daemon.
- constructs/destructs orchestration (using Docker Swarm mode).

## usage
we use two terminal (or background jobs).

export ip address of hosts:
```
$ # both terminal 1 and 2
$ export ROOTLESS_ORCHESTRATION_HOSTS="172.16.0.1,172.16.0.2"
```

start rootless-orchestration server:
```
$ # terminal 1
$ rootless-orchestration
```

start docker daemon using rootless mode and constructs swarm cluster:
```
$ # terminal 2
$ rootless-orchestration dockerd
$ rootless-orchestration swarm
```

to exit rootless-orchestration, send ctrl+c at terminal 1:
```
^C
$
``` 
