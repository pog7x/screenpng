FROM pog7x/gobasebrowser:latest

COPY . /screenpng
WORKDIR /screenpng

VOLUME [ "/screenpng" ]

ENV DISPLAY:=99

EXPOSE 8099

CMD ["go", "run", "main.go", "serve"]
