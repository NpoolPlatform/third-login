FROM jrei/systemd-centos:7

COPY ThirdLoginGateway.yaml /etc/ThirdLoginGateway/ThirdLoginGateway.yaml
COPY third-login-gateway /usr/local/bin/third-login-gateway

CMD ["third-login-gateway", "run"]
