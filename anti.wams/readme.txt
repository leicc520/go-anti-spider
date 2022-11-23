使用golang编译的情况
GOARCH=wasm GOOS=js go build -ldflags "-w -s -X main.GPKey=xxx.xxx.xxx.xx" -o _s.wasm .
使用tinygo编译，
tinygo.exe build -no-debug -ldflags "-X 'main.GPKey=xxx.xxx.xxx.xx'" -target wasm -o _s.wasm .

window环境报错的话：error: could not find wasm-opt, set the WASMOPT environment variable to override
下载这个https://github.com/WebAssembly/binaryen/releases 然后copy wasm-opt.exe 到tinygo/bin目录即可

