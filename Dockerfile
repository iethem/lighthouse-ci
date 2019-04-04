FROM golang
# Install git.
# Git is required for fetching the dependencies.
RUN apt-get update --fix-missing && apt-get -y upgrade

RUN apt-get install --yes curl
RUN curl --silent --location https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get install --yes nodejs
RUN apt-get install --yes build-essential

RUN npm i lighthouse -g

# Install latest chrome dev package.
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list' \
    && apt-get update && apt-get install -y procps vim\
    && apt-get install -y google-chrome-unstable --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /src/*.deb

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/lighthouse-ci/

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["lighthouse-ci"]
