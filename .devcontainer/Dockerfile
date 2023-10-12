FROM --platform=linux/amd64 mcr.microsoft.com/vscode/devcontainers/base:0-bullseye

ENV DENO_INSTALL=/deno
RUN mkdir -p /deno \
    && curl -fsSL https://deno.land/x/install/install.sh | sh \
    && chown -R vscode /deno

ENV PATH=${DENO_INSTALL}/bin:${PATH} \
    DENO_DIR=${DENO_INSTALL}/.cache/deno

ENV DEBIAN_FRONTEND=noninteractive
RUN apt update  &&  apt upgrade -y

RUN apt install -y pkg-config

ENV NODE_VERSION=v20.7.0

RUN wget "https://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}-linux-x64.tar.xz"

RUN tar Jxvf "node-${NODE_VERSION}-linux-x64.tar.xz" -C /usr/local/ --strip-components=1 && rm "node-${NODE_VERSION}-linux-x64.tar.xz"

ENV TYPESCRIPT_VERSION=5.2.2
RUN npm install -g typescript@${TYPESCRIPT_VERSION} 
RUN apt install -y neovim
