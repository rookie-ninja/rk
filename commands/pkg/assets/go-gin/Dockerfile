FROM alpine

ENV WD=/usr/service/

RUN mkdir -p $WD
WORKDIR $WD
COPY target/ $WD

CMD ["bin/demo"]
