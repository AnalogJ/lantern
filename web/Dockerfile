FROM python:alpine
WORKDIR /srv/lantern

RUN apk add curl
#    && curl -O -L https://github.com/ChromeDevTools/devtools-frontend/archive/master.tar.gz \
#    && tar -xvf master.tar.gz --strip 1 \
#    && rm -rf master.tar.gz \
#    && python scripts/build/code_generator_frontend.py /srv/lantern/inspector/protocol.json --output_js_dir /srv/lantern/front_end \

COPY index.html index.html

RUN curl -O -L https://github.com/m3ng9i/ran/releases/download/v0.1.3/ran_linux_amd64.zip \
    && unzip ran_linux_amd64.zip -d /usr/local/bin \
    && rm -rf ran_linux_amd64.zip \
    && ran_linux_amd64 --version \
    && ls -alt /srv/lantern

CMD ["ran_linux_amd64", "-r", "/srv/lantern/", "-p", "8080", "-index", "index.html"]

