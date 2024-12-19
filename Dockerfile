FROM conanio/gcc11-ubuntu16.04:2.9.3

RUN mkdir -p /tmp/conan2 && chmod -R 777 /tmp/conan2

COPY conanfile.txt /tmp/conan2/conanfile.txt

WORKDIR /tmp/conan2

ENTRYPOINT ["tail", "-f", "/dev/null"]