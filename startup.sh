
db_creds_file="/vault/secrets/treco_db_creds"
if [[ -f $db_creds_file ]]
then
    source $db_creds_file
else
    echo "$db_creds_file doesn't exists"
fi

./treco serve