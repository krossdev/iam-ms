github.com=hub.fastgit.org
osarch=darwin
server=nats-server-v2.6.4-$(osarch)-amd64/nats-server

nats:
	nohup $(server) -c ./config.yaml 1>&2 2>>/tmp/nats-server.log &
stop:
	pkill -f "[.]*nats-server -c \.\/config[.]*"

down:
	curl -L -O https://$(github.com)/nats-io/nats-server/releases/download/v2.6.4/nats-server-v2.6.4-$(osarch)-amd64.tar.gz

