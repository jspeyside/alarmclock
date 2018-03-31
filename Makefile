build: clean
	docker run --rm -it -v `pwd`:/go/src/github.com/jspeyside/alarmclock speyside/golang 'set -x && cd /go/src/github.com/jspeyside/alarmclock && export VERSION=`cat VERSION.txt` && go get ./... && go build -o alarmclock -ldflags "-X github.com/jspeyside/alarmclock/domain.Version=$$VERSION"'
	docker build -t speyside/alarmclock:`cat VERSION.txt` .

release: build
	docker tag speyside/alarmclock:`cat VERSION.txt` speyside/alarmclock:latest
	docker push speyside/alarmclock:`cat VERSION.txt` speyside/alarmclock:latest

clean:
	rm -rf build
	rm -f alarmclock
