#!/bin/sh

set -e

while getopts “d” opt; do
  case $opt in
    d) ECHO=echo ;;
  esac
done

echo "Enter your username (Github org):"
read ORG

echo "Enter the name of your project:"
read NAME

update() {
    sed -i -e "s/dnnrly/${ORG}/g" ${1}
    sed -i -e "s/layli/${NAME}/g" ${1}
}

for f in `find . -type f -not -iwholename '*.git/*' -not -iwholename '*tmp*' -not -iwholename '*libexec*' -not -iwholename '*bin*' -not -iwholename '*share*' -not -iwholename '*setup*'`
do
    ${ECHO} update $f
done

${ECHO} mv ./cmd/layli ./cmd/${NAME}
${ECHO} mv ./layli.go ./${NAME}.go
${ECHO} mv ./layli_test.go ./${NAME}_test.go

