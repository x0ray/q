all : qsp.go

qsp.go : qsp.go.y
	go tool yacc -o $@ qsp.go.y; [ -f y.output ] && ( rm -f y.output )
