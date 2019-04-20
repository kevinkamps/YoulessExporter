FROM busybox:1.30-glibc

COPY bin/application_linux_amd64 /bin/application

ENTRYPOINT ["/bin/application"]
