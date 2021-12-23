#!/bin/sh

# run as a daemon
daemonize -c $PWD -e /tmp/mss.err -o /tmp/mss.out -p /tmp/mss.pid -l /tmp/mss.lock \
	$PWD/kiam-ms -c .config/config.yaml -watch
