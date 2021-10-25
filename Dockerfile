FROM faddat/archlinux

ENV PATH $PATH:/root/go/bin
ENV GOPATH /root/go/

RUN pacman -Syyu --noconfirm go

COPY tinyseed .
RUN mv tinyseed /usr/bin

COPY run.sh /usr/bin/
RUN chmod +x /usr/bin/run.sh
ENTRYPOINT ["run.sh"]
