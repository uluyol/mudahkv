FROM scratch
MAINTAINER Muhammed Uluyol <uluyol0@gmail.com>

ADD mudahd /
ADD mudahc /

ENTRYPOINT ["/mudahd"]
