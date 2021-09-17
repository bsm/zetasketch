default: test

test:
	go test ./...

lint:
	golangci-lint run

bench:
	go test ./... -run=NONE -bench=. -benchmem

GOPROTO_PACKAGE=github.com/bsm/zetasketch/internal/zetasketch

# proto task fetches and compiles zetasketch protobuf.
#
# To install protoc-gen-go:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@v1
#
# Protos are explicitly fetched because ALL of them are required
# before any `protoc` calls (roughly - they require each other).
proto: internal/zetasketch/hllplus-unique.proto internal/zetasketch/aggregator.proto internal/zetasketch/unique-stats.proto
	@mkdir -p $(dir $@)
	protoc \
		-I=internal/protobuf \
		-I=internal/zetasketch \
		--go_out=internal/zetasketch \
		--go_opt=paths=source_relative \
		--go_opt=Mhllplus-unique.proto=$(GOPROTO_PACKAGE) \
		--go_opt=Maggregator.proto=$(GOPROTO_PACKAGE) \
		--go_opt=Munique-stats.proto=$(GOPROTO_PACKAGE) \
		$^

internal/zetasketch/%.proto:
	@mkdir -p $(dir $@)
	curl -so $@ https://raw.githubusercontent.com/google/zetasketch/master/proto/$*.proto

internal/protobuf/google/protobuf/descriptor.proto:
	@mkdir -p $(dir $@)
	curl -so $@ https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/descriptor.proto
