certs=$1
db=$2

/cockroach/cockroach.sh sql --certs-dir=${certs} --host=cockroachdb-public:26257 -e "DROP DATABASE IF EXISTS ${db} CASCADE;"
