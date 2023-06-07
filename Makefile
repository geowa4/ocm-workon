BIN_DIR:=./bin
BIN_NAME:=ocm-workon

"${BIN_DIR}/${BIN_NAME}": build

.PHONY: build
build: clean
	mkdir -p ${BIN_DIR}
	go build -mod=mod -o "${BIN_DIR}/${BIN_NAME}" main.go

.PHONY: clean
clean:
	rm -rf "${BIN_DIR}"

.PHONY: test
test:
	go test ./pkg/... ./cmd/... -count=1 -mod=mod

.PHONE: usage
usage: "${BIN_DIR}/${BIN_NAME}"
	"${BIN_DIR}/${BIN_NAME}" help
