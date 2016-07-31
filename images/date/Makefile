# 'Makefile'.
# Chris Shiels.


VERSION = latest


run:	build
	sudo docker run --rm -i -t -p 7001:7001 julio/date:$(VERSION)


build:	Dockerfile src/github.com/chris/date/date.go
	GOPATH=`pwd` go install -ldflags "-X main.version $(VERSION)" github.com/chris/date
	sudo docker build -t julio/date:$(VERSION) .
	touch $@


clean:
	rm -fv build bin/date
