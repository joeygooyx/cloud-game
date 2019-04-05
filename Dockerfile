From golang:1.12

RUN mkdir -p /go/src/github.com/giongto35/game-online
COPY . /go/src/github.com/giongto35/game-online/
WORKDIR /go/src/github.com/giongto35/game-online

# Install server dependencies
RUN apt-get update
#RUN apt-get install portaudio19-dev -y
RUN apt-get install libvpx-dev -y
RUN go get github.com/pions/webrtc
#RUN go get github.com/gordonklaus/portaudio
RUN go get github.com/gorilla/mux
RUN go install github.com/giongto35/game-online

EXPOSE 8000