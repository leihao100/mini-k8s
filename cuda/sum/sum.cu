#include <stdio.h>
#include <cuda_runtime.h>

#define CHECK_ERROR(call)\
{\
  const cudaError_t error = call;\
  if (error != cudaSuccess)\
  {\
      printf("ERROR: %s:%d, ", __FILE__, __LINE__);\
      printf("code:%d, reason:%s\n", error, cudaGetErrorString(error));\
      exit(1);\
  }\
}

const int M = 15;
const int N = 10;

__global__ void matrix_add(int *A, int *B, int *C, int M, int N) {
    int i = blockIdx.x * blockDim.x + threadIdx.x;
    int j = blockIdx.y * blockDim.y + threadIdx.y;

    if (i < M && j < N) {
        int index = i * N + j;
        C[index] = A[index] + B[index];
    }
}

int main() {
    int size = M * N * sizeof(int);

    int *h_A = (int *)malloc(size);
    int *h_B = (int *)malloc(size);
    int *h_C = (int *)malloc(size);

    for (int i = 0; i < M * N; i++) {
        h_A[i] = i;
        h_B[i] = i;
    }

    int *d_A, *d_B, *d_C;
    CHECK_ERROR(cudaMalloc((void **)&d_A, size));
    CHECK_ERROR(cudaMalloc((void **)&d_B, size));
    CHECK_ERROR(cudaMalloc((void **)&d_C, size));

    CHECK_ERROR(cudaMemcpy(d_A, h_A, size, cudaMemcpyHostToDevice));
    CHECK_ERROR(cudaMemcpy(d_B, h_B, size, cudaMemcpyHostToDevice));

    dim3 threadsPerBlock(16, 16);
    dim3 blocksPerGrid((M + threadsPerBlock.x - 1) / threadsPerBlock.x, (N + threadsPerBlock.y - 1) / threadsPerBlock.y);

    matrix_add<<<blocksPerGrid, threadsPerBlock>>>(d_A, d_B, d_C, M, N);

    CHECK_ERROR(cudaMemcpy(h_C, d_C, size, cudaMemcpyDeviceToHost));

    printf("Matrix A + Matrix B = Matrix C:\n");
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < N; j++) {
            printf("%d ", h_C[i * N + j]);
        }
        printf("\n");
    }

    free(h_A);
    free(h_B);
    free(h_C);
    cudaFree(d_A);
    cudaFree(d_B);
    cudaFree(d_C);

    return 0;
}
