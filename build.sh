go get github.com/jteeuwen/go-bindata/...

debug=""
if [ $1 = "debug" ]; then
   debug="-debug"
   echo debug
else 
    echo non debug
fi

go-bindata $debug -prefix "assets/" assets/...
