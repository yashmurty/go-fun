set -o allexport
eval $(grep -v '^#' .env | sed 's/^/export /')
set +o allexport