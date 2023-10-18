

put() {
    echo "Post To Docker"
    curl -X POST -d "$*" http://localhost:8081/put/string --output -

}

get() {
    echo "Get From Docker"
    curl http://localhost:8081/get?hash=$@ --output -

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
    *)
        exit 1
        ;;
esac

exit 0