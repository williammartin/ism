all:
	go build -o sm cmd/sm/main.go

clean:
	rm -f sm

test:
	ginkgo -r integration
