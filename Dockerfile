
FROM faddat/archlinux

ENV PATH $PATH:/root/go/bin
ENV GOPATH /root/go/

RUN pacman -Syyu --noconfirm go

COPY . . 
RUN go mod tidy
RUN go install .
RUN mv ~/go/bin/tenderseed /usr/bin/

COPY run.sh /usr/bin/
RUN chmod +x /usr/bin/run.sh
ENTRYPOINT ["run.sh"]
