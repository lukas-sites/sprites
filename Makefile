
main.wasm: visualizer.go
	GOOS=js GOARCH=wasm go build -o $@

serve: main.wasm
	goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))'