FROM alpine

COPY dnsimple /usr/bin

ENTRYPOINT [ "dnsimple" ]

CMD [ "dnsimple" ]
