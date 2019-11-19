FROM node
WORKDIR /source
ADD ./runTest.sh /source
RUN chmod +x ./runTest.sh
RUN mkdir /data
RUN npm install -g mocha
RUN npm install -g https://github.com/craftybones/mocha-json-reporter
ENTRYPOINT ["./runTest.sh"]