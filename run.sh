#!/bin/bash
export PATH=$PATH:/usr/local/go/bin
cd /home/fakebear/go/src/CampuSafe/
go run /home/fakebear/go/src/CampuSafe/main.go <<< "postgresql://kaiser:asdfqwrlkjhqwer@free-tier.gcp-us-central1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&sslrootcert=$HOME/.postgresql/root.crt&options=--cluster%3Dcampusafe-3594"
