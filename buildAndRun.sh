#!/usr/bin/env bash
# go build -o out && ./out
go build -gcflags="all=-N -l" -o tubely && ./tubely
