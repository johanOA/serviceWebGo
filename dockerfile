FROM debian:latest


RUN apt-get update -y && apt-get install -y  \
    curl \
    build-essential \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

    
#Instalacion de Go
RUN curl -OL https://go.dev/dl/go1.20.6.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.20.6.linux-arm64.tar.gz && \
    rm go1.20.6.linux-arm64.tar.gz

#Anadir Go al Path
ENV PATH="/usr/local/go/bin:${PATH}"

#Instalacion Node JS y npm
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm

#Crea el directorio de trabajo
WORKDIR /goWeb

#Crea la carpeta de imagenes
RUN mkdir imagenes

COPY . .
# Copiar package.json y otros archivos necesarios
COPY /goWeb/package.json .
COPY /goWeb/tailwind.config.js .

RUN npm install

RUN npm run build-css


# Inicializar el módulo Go
RUN go mod init app

RUN go mod tidy

# Compilar la aplicación Go
RUN go build -o app ./main.go

#Instala wget para las imagenes
RUN apt-get install wget 

RUN wget -P /imagenes https://i.ytimg.com/vi/KPLWWIOCOOQ/maxresdefault.jpg && \
    wget -P /imagenes https://static.hbo.com/game-of-thrones-1-1920x1080.jpg && \
    wget -P /imagenes https://c.files.bbci.co.uk/53B2/production/_106462412_fb-heroinas-got.jpg

RUN pwd    

EXPOSE 8080

# Comando para ejecutar cuando inicie el contenedor
CMD ["go", "run", "main.go", "./goWeb/imagenes"]