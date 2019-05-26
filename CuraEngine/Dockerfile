FROM trzeci/emscripten

RUN git clone https://github.com/kwonoj/protobuf-wasm && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protobuf-all-3.6.1.zip && \
    unzip protobuf-all-3.6.1.zip && \
    ls && \
    cd protobuf-3.6.1 && \
    cp ../protobuf-wasm/*.patch ./ && \
    git apply *.patch && \
    echo PATCHED PROTOBUF

RUN apt-get update && apt-get -y install autoconf automake libtool curl make g++ unzip

RUN cd protobuf-3.6.1 && \
    sh autogen.sh && \
    emconfigure ./configure CXXFLAGS="-O2" && \
    emmake make 

# RUN cd protobuf-3.6.1 && \
#     cp -r ./src/.libs /src/

RUN mkdir -p hello && \
    echo "int main() { return 0; }" > hello/hello.c && \
    emcc hello/hello.c -s WASM=1 -s SINGLE_FILE=1 -o hello/hello.js && \
    rm -rf hello

RUN apt-get -y install python3 python3-dev python3-sip-dev


# real men don't need apt-get
RUN wget http://ftp.br.debian.org/debian/pool/main/p/protobuf/protobuf-compiler_3.6.1.3-2_amd64.deb && \
    wget http://ftp.br.debian.org/debian/pool/main/p/protobuf/libprotoc17_3.6.1.3-2_amd64.deb && \
    wget http://ftp.br.debian.org/debian/pool/main/p/protobuf/libprotobuf17_3.6.1.3-2_amd64.deb && \
    dpkg -i libprotobuf17_3.6.1.3-2_amd64.deb && \
    dpkg -i libprotoc17_3.6.1.3-2_amd64.deb && \
    dpkg -i protobuf-compiler_3.6.1.3-2_amd64.deb
RUN git clone https://github.com/Ultimaker/libArcus.git

RUN cd libArcus && \
    mkdir build && cd build && \
    emconfigure cmake .. -DProtobuf_LIBRARIES=/src/protobuf-3.6.1/src/.libs -DProtobuf_INCLUDE_DIR=/src/protobuf-3.6.1/src -DProtobuf_VERSION=3.6.1 -DBUILD_PYTHON=off -DBUILD_EXAMPLES=off && \
    emmake make

RUN git clone https://github.com/Ultimaker/CuraEngine.git
RUN cd CuraEngine && \
    mkdir build && cd build && \
    emconfigure cmake .. -DProtobuf_LIBRARIES=/src/protobuf-3.6.1/src/.libs -DProtobuf_INCLUDE_DIR=/src/protobuf-3.6.1/src -DProtobuf_VERSION=3.6.1 -DArcus_DIR=/src/libArcus/build -DENABLE_ARCUS=off && \
    emmake make



CMD ["/bin/bash"]