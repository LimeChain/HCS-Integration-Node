FILE=peer3.env
if [ ! -f "$FILE" ]; then
    echo "$FILE does not exist. Please create it based on the .env example to run peer 2"
	exit 1
fi

go run cmd/*.go $FILE