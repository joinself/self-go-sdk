#!/bin/sh

content=`cat ./_support/sources.json`
gocontent="package fact\n"
gocontent+="\n"
gocontent+='var sourceDefinition = []byte(`'
gocontent+=$content
gocontent+='`)\n'	
echo "$gocontent"
