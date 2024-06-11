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
const int P = 20;

__global__ void matrix_multiply(int *A, int *B, int *C, int M, int N, int P) {
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;

    if (row < M && col < P) {
        int value = 0;
        for (int k = 0; k < N; k++) {
            value += A[row * N + k] * B[k * P + col];
        }
        C[row * P + col] = value;
    }
}

int main() {
    int size_A = M * N * sizeof(int);
    int size_B = N * P * sizeof(int);
    int size_C = M * P * sizeof(int);

    int *h_A = (int *)malloc(size_A);
    int *h_B = (int *)malloc(size_B);
    int *h_C = (int *)malloc(size_C);

    for (int i = 0; i < M * N; i++) {
        h_A[i] = i % 10;
    }
    for (int i = 0; i < N * P; i++) {
        h_B[i] = i % 10;
    }

    int *d_A, *d_B, *d_C;
    CHECK_ERROR(cudaMalloc((void **)&d_A, size_A));
    CHECK_ERROR(cudaMalloc((void **)&d_B, size_B));
    CHECK_ERROR(cudaMalloc((void **)&d_C, size_C));

    CHECK_ERROR(cudaMemcpy(d_A, h_A, size_A, cudaMemcpyHostToDevice));
    CHECK_ERROR(cudaMemcpy(d_B, h_B, size_B, cudaMemcpyHostToDevice));

    dim3 threadsPerBlock(16, 16);
    dim3 blocksPerGrid((P + threadsPerBlock.x - 1) / threadsPerBlock.x, (M + threadsPerBlock.y - 1) / threadsPerBlock.y);

    matrix_multiply<<<blocksPerGrid, threadsPerBlock>>>(d_A, d_B, d_C, M, N, P);

    CHECK_ERROR(cudaMemcpy(h_C, d_C, size_C, cudaMemcpyDeviceToHost));

    printf("Matrix A * Matrix B = Matrix C:\n");
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < P; j++) {
            printf("%d ", h_C[i * P + j]);
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
