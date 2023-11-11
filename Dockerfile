FROM alpine:3.18.4

#WORKDIR /gateway
ADD output/swiss /swiss
ADD output/wintun.dll /wintun.dll

RUN chmod +x /swiss

EXPOSE 8081
ENTRYPOINT [ "./swiss" ]

#CMD [ "cat", "./config.yaml" ]
