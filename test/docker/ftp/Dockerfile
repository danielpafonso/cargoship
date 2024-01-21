FROM alpine:latest
## install vsftp, set password login and create user
RUN apk add --update --no-cache vsftpd; \
	adduser -s /bin/sh -D ftpuser; \
	echo -n 'ftpuser:passwd' | chpasswd; \
	mkdir -p /home/ftpuser; \
	echo "hello" >> /home/ftpuser/welcome.txt

COPY vsftpd.conf /etc/vsftpd/vsftpd.conf

EXPOSE 21
#EXPOSE 21100-21110
ENTRYPOINT ["/usr/sbin/vsftpd", "/etc/vsftpd/vsftpd.conf"]
