default: vet test

vet:
	go vet ./...

test:
	go test ./...

# cd ~ && go get -u github.com/golang/protobuf/protoc-gen-go
proto: \
	internal/zetasketch/hllplus-unique.proto \
	internal/zetasketch/aggregator.proto \
	internal/zetasketch/unique-stats.proto \
	tmp/protobuf/src/google/protobuf/descriptor.proto \
	\
	internal/zetasketch/hllplus-unique.pb.go \
	internal/zetasketch/aggregator.pb.go \
	internal/zetasketch/unique-stats.pb.go

internal/zetasketch/%.pb.go: internal/zetasketch/%.proto
		@mkdir -p $(dir $@)
		protoc \
			-I=internal/zetasketch:tmp/protobuf/src \
			--go_out=internal/zetasketch \
			--go_opt=paths=source_relative \
			$<

internal/zetasketch/%.proto:
	@mkdir -p $(dir $@)
	curl -so $@ https://raw.githubusercontent.com/google/zetasketch/master/proto/$*.proto

tmp/protobuf/src/google/protobuf/descriptor.proto:
	@mkdir -p $(dir $@)
	curl -so $@ https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/descriptor.proto \

clean:
	rm -rf tmp
	rm -rf internal/zetasketch/*.proto
	rm -rf internal/zetasketch/*.pb.go
