"""CUDA function library."""

from numba import jit, float32, cuda

@cuda.jit
def SampleIncrementCUDA(an_array):
    """A test CUDA function which increments input array by one."""
    # Thread id in a 1D block
    tx = cuda.threadIdx.x
    # Block id in a 1D grid
    ty = cuda.blockIdx.y
    # Block width, i.e. number of threads per block
    # bpg = cuda.gridDim.x
    # Block width, i.e. number of threads per block
    bw = cuda.blockDim.x
    # Compute flattened index inside the array
    pos = tx + ty * bw
    if pos < an_array.size:  # Check array boundaries
        an_array[pos] += 1