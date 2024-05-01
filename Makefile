.PHONY: coverage

coverage:
	go test -short -v -covermode=count -coverprofile=coverage.out  | tee test_output.txt
	go tool cover -func=coverage.out | awk '/total/ {print "| **" $$1 "** | **" $$3 "** |"}' | tee coverage.txt
	cat test_output.txt | grep 'ok.*coverage' | awk '{print "| " $$2 " | " $$5 " |"}' | tee -a coverage.txt
	echo "## Test Coverage Report" > coverage_with_header.txt
	echo "| Package           | Coverage |" >> coverage_with_header.txt
	echo "|-------------------|----------|" >> coverage_with_header.txt
	cat coverage.txt >> coverage_with_header.txt
	mv coverage_with_header.txt coverage.txt
	rm test_output.txt
