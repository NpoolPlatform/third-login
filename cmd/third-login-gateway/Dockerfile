FROM jrei/systemd-centos:7

WORKDIR /app

COPY third-login-gateway /app
COPY ThirdLoginGateway.viper /app

ENTRYPOINT ["./third-login-gateway"]

EXPOSE 50280 50281 50282

CMD ["run"]