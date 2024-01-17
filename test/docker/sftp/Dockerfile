FROM alpine:latest
## install openssh, set password login and create user
RUN apk add --update --no-cache openssh; \
	echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config; \
	adduser -s /bin/sh -D ftpuser; \
	echo -n 'ftpuser:passwd' | chpasswd; \
	mkdir -p /home/ftpuser; \
	echo "hello" >> /home/ftpuser/welcome.txt

COPY entrypoint.sh /
EXPOSE 22
ENTRYPOINT ["/entrypoint.sh"]
