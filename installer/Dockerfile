FROM fedora

RUN yum install -y ruby-devel gcc make rpm-build rubygems

RUN gem install --no-document fpm

WORKDIR /usr/src/app

ENTRYPOINT [ "fpm","-s", "dir", "-t", "rpm", "-n", "vauth", "--after-install", "./install.sh", "-v", "0.0.1"]