#!/bin/bash
echo "=============== Build Boong Start ==============="
go build -o out/boong ./cmd/boong
if [ -f out/boong ]; then
    echo "Success to build boong : out/boong"
    echo "Example: cd test/HelloWorld && ../../out/boong -compiler gcc && ninja -C out"
else
    echo "Failed to build boong..."
fi
echo "=============== Build Boong End ================="