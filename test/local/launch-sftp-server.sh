# run server
docker run -d -v $(pwd)/ftpdata:/home/ftpuser \
-p 22:22 \
--name sftpd --rm sftpd
