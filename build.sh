#!/bin/env bash
docker build . --output "${PWD}" --target copytohost
