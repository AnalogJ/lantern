FROM python:2-alpine
WORKDIR /srv/lantern

# see https://groups.google.com/forum/#!searchin/google-chrome-developer-tools/compile|sort:date/google-chrome-developer-tools/Lh5dgrfJ8c0/ihfnWUUrBQAJ

RUN apk add curl \
    && curl -O -L https://github.com/ChromeDevTools/devtools-frontend/archive/6bd6d4996c0e1dd424c85540e298338c9aa913ad.tar.gz \
    && tar -xvf 6bd6d4996c0e1dd424c85540e298338c9aa913ad.tar.gz --strip 1 \
    && rm -rf 6bd6d4996c0e1dd424c85540e298338c9aa913ad.tar.gz \
    && curl -o front_end/InspectorBackendCommands.js -L https://chrome-devtools-frontend.appspot.com/serve_file/@38db055e5fc20b2eddca2c829c324fb49de07cbf/InspectorBackendCommands.js \
    && curl -o front_end/SupportedCSSProperties.js -L https://chrome-devtools-frontend.appspot.com/serve_file/@38db055e5fc20b2eddca2c829c324fb49de07cbf/SupportedCSSProperties.js

COPY index.html index.html

RUN curl -O -L https://github.com/m3ng9i/ran/releases/download/v0.1.3/ran_linux_amd64.zip \
    && unzip ran_linux_amd64.zip -d /usr/local/bin \
    && rm -rf ran_linux_amd64.zip \
    && ran_linux_amd64 --version \
    && ls -alt /srv/lantern

CMD ["ran_linux_amd64", "-r", "/srv/lantern/", "-p", "8080", "-index", "index.html"]

