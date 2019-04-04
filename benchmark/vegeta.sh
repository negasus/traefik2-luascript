#!/bin/bash

echo "GET http://localhost/" | vegeta attack -rate 2000 -duration=60s | tee results.bin | vegeta report