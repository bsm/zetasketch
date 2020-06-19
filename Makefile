default: vet test

vet:
	go vet ./...

test:
	go test ./...

# cd ~ && go get -u github.com/golang/protobuf/protoc-gen-go
proto: \
	internal/zetasketch/hllplus-unique.pb.go \
	internal/zetasketch/aggregator.pb.go \
	internal/zetasketch/unique-stats.pb.go

internal/zetasketch/%.pb.go: \
	tmp/zetasketch/proto/hllplus-unique.proto \
	tmp/zetasketch/proto/aggregator.proto \
	tmp/zetasketch/proto/unique-stats.proto \
	tmp/protobuf/src/google/protobuf/descriptor.proto
		@mkdir -p $(dir $@)
		protoc \
			-I=tmp/zetasketch/proto:tmp/protobuf/src \
			--go_out=internal/zetasketch \
			--go_opt=paths=source_relative \
			tmp/zetasketch/proto/$*.proto

tmp/zetasketch/proto/%.proto:
	@mkdir -p $(dir $@)
	curl -so $@ https://raw.githubusercontent.com/google/zetasketch/master/proto/$*.proto

tmp/protobuf/src/google/protobuf/descriptor.proto:
	@mkdir -p $(dir $@)
	curl -so $@ https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/descriptor.proto \

clean:
	rm -rf tmp
	rm -rf internal/zetasketch/*.pb.go
