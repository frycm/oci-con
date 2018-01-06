# Containers from scratch

Code inspired by _Containers from scratch_ youtube presentations and [What even is container](https://jvns.ca/blog/2016/10/10/what-even-is-a-container/).

## Commands

Target run commands

    docker    run <image> <cmd> <args>
    ocscratch run         <cmd> <args>

Run shell

    go run ocscratch.go run /bin/sh

## Preparation

Download and unpack Alpine Linux Mini Root File System from https://alpinelinux.org/downloads/
and set it's location in **rootFsLocation** constant

## History

2000 - FreeBDS: Jails (chroot files isolation)
2001 - Linux: vServer (port to Linux)
2004 - Solaris: Zones (Snapshots)
2006 - Google: Linux Process containers (cgroups) -> Borg -> Kubernetes
2008 - RedHat: Namespaces, limitng root in containers
2008 - IBM: LXC (user tools for containers and cgroups)
2013 - Docker Inc: Docker (simple user tools and images)