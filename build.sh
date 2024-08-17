#!/bin/bash


mkdir -p output/bin

# 敏感词检测
go build -v -o ./output/bin/test_web_server ./command/test_program
