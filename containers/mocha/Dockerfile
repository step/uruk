FROM usemtech/nodejs-mocha
WORKDIR /source
ADD ./runTest.sh /source
RUN chmod +x ./runTest.sh
ENTRYPOINT ["./runTest.sh"]