FROM golang:latest as builder
LABEL maintainer="Arman Shah <ashah360@uw.edu>"
ARG GH_TOKEN
RUN go env -w GOPRIVATE="github.com/ashah360/*"
RUN git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com".insteadOf "https://github.com"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../main .

FROM alpine:latest
ARG DOPPLER_TOKEN
RUN apk --no-cache add curl
RUN apk --no-cache add gnupg
RUN apk --no-cache add ca-certificates
# Install the Doppler CLI
RUN (curl -Ls https://cli.doppler.com/install.sh || wget -qO- https://cli.doppler.com/install.sh) | sh
WORKDIR /root/
COPY --from=builder /app/main .
ENV PORT 3000
ENV DOPPLER_TOKEN ${DOPPLER_TOKEN}
EXPOSE 3000
CMD ["doppler", "run", "--", "./main"]