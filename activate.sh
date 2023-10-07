COMPOSE_FILE="docker/docker-compose.yml"
file="main"

docker_build() {
    echo "Building Docker services..."
    #go build -o main main.go
    cp $file docker/$file
    [ -f $file ] && rm $file
    docker compose -f $COMPOSE_FILE build
}

docker_run() {
    echo "Running Docker services..."
    docker compose -f $COMPOSE_FILE --compatibility up -d 
}

case "$1" in
    "build")
        docker_build
        ;;
    "run")
        docker_run
        ;;
    *)
        exit 1
        ;;
esac

exit 0