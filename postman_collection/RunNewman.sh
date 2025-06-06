set -e;
onExit() {
    exitCode=$?
    echo  "Exit Code ==> "$exitCode
    if [ $exitCode != "0" ]; then
        echo "Tests failed"
        # build failed, don't deploy
        exit 1
    else
        echo "Tests passed"
        # deploy build
    fi
}

trap onExit EXIT
# --bail
sleep 3
newman run GoLangProgram.postman_collection.json \
            --delay-request=300  \
            --reporters=cli,junit,json \
            --iteration-count=3 \
            --environment=GoLangProgram.postman_environment.json