.PHONY: test, publish_version

test:
	go test ./...


publish_version: test
	git tag ${TAG_VERSION}
	git push --tags
	GOPROXY=proxy.golang.org go list -m github.com/sohail288/go_payload_parser@${TAG_VERSION}
