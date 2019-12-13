FROM node
WORKDIR /source
ADD ./runLint.sh /source
RUN chmod +x ./runLint.sh
RUN mkdir /data
RUN npm install -g eslint
ENTRYPOINT ["./runLint.sh"]