project = cloudraker
projectpath = ${PWD}
glidepath = ${PWD}/vendor/github.com/Masterminds/glide
keyspath = ${PWD}/files/docker

target:
	@go build

test:
	@go test

integration: test
	@go test -tags=integration

install:
	cp $(project) /usr/local/bin
	chmod 755 /usr/local/bin/$(project)
	mkdir -p /usr/local/share/$(project)
	cp files/my.cnf /usr/local/share/$(project)
	chmod -R 755 /usr/local/share/$(project)

$(keyspath)/id_rsa:
	ssh-keygen -t rsa -N "" -f files/docker/id_rsa

builddockercontainer:
	cd files/docker;docker build -t cloudraker/mysql-server:5.7 .

$(glidepath)/glide:
	git clone https://github.com/Masterminds/glide.git $(glidepath)
	cd $(glidepath);make build
	cp $(glidepath)/glide .

libs: $(glidepath)/glide
	$(glidepath)/glide install

deps: $(keyspath)/id_rsa builddockercontainer libs
	sudo gem install ghost
