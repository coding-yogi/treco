if [[ -f "/vault/secrets/dev-treco-db01-rw" ]]
then
    source "/vault/secrets/dev-treco-db01-rw"
fi

./treco serve