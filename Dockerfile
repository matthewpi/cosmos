#
# Copyright (c) 2021 Matthew Penner
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
#

# Stage 1 (Build)
FROM        --platform=$BUILDPLATFORM golang:1.16-alpine3.13

RUN         apk add --update --no-cache ca-certificates=20191127-r5 git=2.30.2-r0 tzdata=2021a-r0

WORKDIR     /app/
COPY        go.mod go.sum /app/
RUN         go mod download
COPY        . /app/

RUN         CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -v -o cosmos cmd/cosmos/main.go

# Stage 2 (Final)
FROM        alpine:3.13

LABEL       author="Matthew Penner" maintainer="me@matthewp.io"

RUN         apk add --update --no-cache ca-certificates=20191127-r5 tzdata=2021a-r0
COPY        --from=builder /app/cosmos /usr/bin/

ENTRYPOINT  [ "cosmos" ]
