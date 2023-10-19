

put() {
    echo "Post To Docker"
    port="$1"
    shift
    curl -X POST -d "$*" http://localhost:$port/put/string --output -

}

get() {
    echo "Get From Docker"
    port="$1"
    shift
    curl http://localhost:$port/get?hash=$@ --output -

}

whoami() {
    echo "Get From Docker"
    curl http://localhost:$1/whoami?=$2 --output -

}

case "$1" in
    "put")
        shift
        put "$@"
        ;;
    "get")
        shift
        get "$@"
        ;;
    "whoami")
        shift
        whoami "$@"
        ;;
    *)
        exit 1
        ;;
esac

exit 0