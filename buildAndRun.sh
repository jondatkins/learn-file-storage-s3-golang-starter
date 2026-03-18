#!/usr/bin/env bash
# go build -o out && ./out
chromium  http://localhost:8091/app/
go build -gcflags="all=-N -l" -o tubely && ./tubely  
