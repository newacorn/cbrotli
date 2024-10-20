rm -rf out&&mkdir out && cd out && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=./installed .. && cmake --build . --config Release --target install

// #cgo CFLAGS: -I /Users/acorn/workspace/programming/clang/brotli/out/installed/include
// #cgo LDFLAGS:-L/Users/acorn/workspace/programming/clang/brotli/out/installed/lib -lbrotlicommon
// #cgo LDFLAGS: -L/Users/acorn/workspace/programming/clang/brotli/out/installed/lib -lbrotlienc
// #cgo LDFLAGS: -L/Users/acorn/workspace/programming/clang/brotli/out/installed/lib -lbrotlidec
import "C"
