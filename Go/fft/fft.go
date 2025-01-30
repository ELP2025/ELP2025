package fft

/*
#cgo LDFLAGS: -lfftw3 -lfftw3_threads
#include <fftw3.h>
#include <stdlib.h>
*/
import "C"
import (
	"sync"
	"unsafe"
)

var fftMutex sync.Mutex

// fftRealToComplex computes the FFT of a real-valued input []float64
func FFT(input []float64) []complex128 {
	fftMutex.Lock()  // Lock the mutex to ensure thread safety during FFT execution
	defer fftMutex.Unlock()

	n := len(input)

	// Allocate memory for input and output arrays
	in := (*C.double)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.double(0)))))
	out := (*C.fftw_complex)(C.malloc(C.size_t((n/2+1)) * C.size_t(unsafe.Sizeof(C.fftw_complex{}))))
	defer C.free(unsafe.Pointer(in))
	defer C.free(unsafe.Pointer(out))

	// Convert Go float64 slice to FFTW input array
	inSlice := (*[1 << 30]C.double)(unsafe.Pointer(in))[:n:n]
	for i := 0; i < n; i++ {
		inSlice[i] = C.double(input[i])
	}

	// Create FFTW plan for real-to-complex FFT
  C.fftw_plan_with_nthreads(C.int(12))
	plan := C.fftw_plan_dft_r2c_1d(C.int(n), in, out, C.FFTW_ESTIMATE)
	defer C.fftw_destroy_plan(plan)

	// Execute FFT
	C.fftw_execute(plan)

	// Convert FFTW output to Go []complex128
	outputSize := n/2 + 1
	output := make([]complex128, outputSize)
	outSlice := (*[1 << 30]C.fftw_complex)(unsafe.Pointer(out))[:outputSize:outputSize]
	for i := 0; i < outputSize; i++ {
		output[i] = complex(float64(outSlice[i][0]), float64(outSlice[i][1]))
	}

	return output
}

func IFFT(input []complex128, originalSize int) []float64 {
	fftMutex.Lock()  // Lock the mutex to ensure thread safety during IFFT execution
	defer fftMutex.Unlock()

	n := len(input)

	// Allocate memory for FFTW input (complex) and output (real)
	in := (*C.fftw_complex)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.fftw_complex{}))))
	out := (*C.double)(C.malloc(C.size_t(originalSize) * C.size_t(unsafe.Sizeof(C.double(0)))))
	defer C.free(unsafe.Pointer(in))
	defer C.free(unsafe.Pointer(out))

	// Convert Go complex128 slice to FFTW's input format
	inSlice := (*[1 << 30]C.fftw_complex)(unsafe.Pointer(in))[:n:n]
	for i := 0; i < n; i++ {
		inSlice[i][0] = C.double(real(input[i])) // Real part
		inSlice[i][1] = C.double(imag(input[i])) // Imaginary part
	}

	// Create FFTW plan for complex-to-real IFFT
	plan := C.fftw_plan_dft_c2r_1d(C.int(originalSize), in, out, C.FFTW_ESTIMATE)
	defer C.fftw_destroy_plan(plan)

	// Execute IFFT
	C.fftw_execute(plan)

	// Convert FFTW output to Go []float64 and normalize (IFFT output is unnormalized)
	output := make([]float64, originalSize)
	outSlice := (*[1 << 30]C.double)(unsafe.Pointer(out))[:originalSize:originalSize]
	for i := 0; i < originalSize; i++ {
		output[i] = float64(outSlice[i]) / float64(originalSize) // Normalize IFFT output
	}

	return output
}

