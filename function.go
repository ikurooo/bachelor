package main

type Request struct {
	Method string
	Path   string
	Body   []byte
}

type Response struct {
	Status int
	Body   []byte
}

func Handle(request Request) Response {
	if request.Method == "GET" && request.Path == "/hello" {
		return Response{
			Status: 200,
			Body:   []byte("Hello from my Wasm function!\n"),
		}
	}

	return Response{
		Status: 404,
		Body:   []byte("Not found\n"),
	}
}
