# 'Makefile'.
# Chris Shiels.


SUBJ = '/C=GB/ST=London/L=London/O=Julio/OU=Julio/CN=vm1'


vm1.crt:	vm1.key
	openssl req \
	-x509 \
	-new  \
	-key vm1.key \
	-nodes \
	-sha256 \
	-subj $(SUBJ) \
	-days 365 \
	-out ./vm1.crt


vm1.key:
	openssl genrsa -out vm1.key 4096


clean:
	rm -f vm1.crt vm1.key
