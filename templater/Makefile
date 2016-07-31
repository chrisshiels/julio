# 'Makefile'.
# Chris Shiels.


VERSION = latest


templater:	templater.go
	go build -ldflags "-X main.version $(VERSION)" templater.go


clean:
	rm -fv templater
