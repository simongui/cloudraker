project = cloudraker
projectpath = ${PWD}
glidepath = ${PWD}/vendor/github.com/Masterminds/glide

target:
	@go build

test:
	@go test

integration: test
	@go test -tags=integration

install:
	@cp $(project) /usr/local/bin
	@chmod 755 /usr/local/bin/$(project)
	@mkdir -p /usr/local/share/$(project)
	@cp files/my.cnf /usr/local/share/$(project)
	@chmod -R 755 /usr/local/share/$(project)

builddockercontainer:
	cd files/docker;docker build -t cloudraker/mysql-server:5.7 . > /dev/null 2> /dev/null;

$(glidepath)/glide:
	@git clone https://github.com/Masterminds/glide.git $(glidepath)
	@cd $(glidepath);make build
	@cp $(glidepath)/glide .

deps: builddockercontainer $(glidepath)/glide
	@$(glidepath)/glide install
	@sudo gem install ghost
