// Copyright 2017 Google Inc. All Rights Reserved.
//
// Distributed under MIT license.
// See file LICENSE for detail or copy at https://opensource.org/licenses/MIT

package cbrotli

// Inform golang build system that it should link brotli libraries.

// #cgo CFLAGS: -I /Users/acorn/workspace/programming/clang/brotli/out/installed/include
// #cgo LDFLAGS:-L/Users/acorn/workspace/programming/clang/brotli/out/installed/lib -lbrotlicommon
// #cgo LDFLAGS: -L/Users/acorn/workspace/programming/clang/brotli/out/installed/lib -lbrotlienc
// #cgo LDFLAGS: -L/Users/acorn/workspace/programming/clang/brotli/out/installed/lib -lbrotlidec
import "C"
