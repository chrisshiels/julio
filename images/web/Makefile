# 'Makefile'.
# Chris Shiels.


VERSION = latest


run:	build
	sudo docker run --rm -i -t -p 7000:7000 julio/web:$(VERSION)


build:	Dockerfile src/github.com/chris/web/web.go
	GOPATH=`pwd` go install -ldflags "-X main.version $(VERSION)" github.com/chris/web
	sudo docker build -t julio/web:$(VERSION) .
	touch $@


clean:
	rm -fv build bin/web
