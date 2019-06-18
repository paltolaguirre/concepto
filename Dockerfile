FROM scratch

# Se setea variable de entorno
ARG RECENV=P
ENV RECENV $RECENV

# Se instala app
ADD main /
ADD configuracion.json /
# Se copia certificado de maquina host
#COPY ca-certificates.crt /etc/ssl/certs/

# Se habilita puerto
EXPOSE 80

# Se ejecuta solo cuando se corre "docker run"
ENTRYPOINT ["/main"]
