

function down() {
  migrate -path=./migrations -database postgres://localhost:5432/persistsql\?sslmode=disable down "$1";
}

function up() {
  migrate -path=./migrations -database postgres://localhost:5432/persistsql\?sslmode=disable up "$1";
}

function force() {
  migrate -path=./migrations -database postgres://localhost:5432/persistsql\?sslmode=disable force "$1";
}

function downup() {
  down "$@"
  up "$@"
}

function updown() {
  up "$@"
  down "$@"
}

function forcedown() {
  force "$@"
  down "$@"
}

function forcedownup() {
  force "$@"
  downup "$@"
}
