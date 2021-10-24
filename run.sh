#!/bin/bash

if [ -z "$SEEDS" ]; then
  echo "SEEDS not defined"
  exit 1
fi

if [ -z "$CHAIN_ID" ]; then
  echo "CHAIN_ID not defined"
  exit 1
fi

tenderseed -seeds="$SEEDS" -chain-id "$CHAIN_ID" start
