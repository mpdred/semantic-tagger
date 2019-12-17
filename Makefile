build:
	go build
#	docker build --tag semtag .
#	docker create --name tmp_semtag semtag
#	docker cp tmp_semtag:/semtag ./
#	docker rm tmp_semtag

clean:
	rm -v semtag
