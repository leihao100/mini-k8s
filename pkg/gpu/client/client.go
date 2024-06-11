package gpuclient

import (
	"MiniK8S/config"
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
)

const HPCUsername = ""
const HPCPassword = ""

type Client interface {
	Run(ctx context.Context)
	SubmitCudaJob(jobUID uuid.UUID, cudaFilePath string, slurmFileContent string, objectFileName string) (jobID string, err error)
	CheckJobFinish(jobID string) (bool, error)
	GetJobState(jobID string) (string, error)
	DownloadResult(jobUID uuid.UUID, localFilePath string, resultFileName string) (bool, error)
	CreateAndWriteFile(filePath string, content string) error
}

type JobClient struct {
	gophClient *goph.Client
	sftpClient *sftp.Client
}

func New() Client {
	return &JobClient{
		gophClient: nil,
		sftpClient: nil,
	}
}

func (c *JobClient) Run(ctx context.Context) {
	syncChan := make(chan bool)
	go c.run(ctx, syncChan)
	<-syncChan
}

func (c *JobClient) run(ctx context.Context, syncChan chan bool) {

	if c.gophClient == nil {
		var err error
		c.gophClient, err = goph.New(HPCUsername, config.PiHost, goph.Password(HPCPassword))
		if err != nil {
			log.Fatal(fmt.Sprintf("[jobClient] New ssh client failed: %v\n", err))
		}
		c.newSftp()
	}

	defer c.gophClient.Close()

	log.Printf("[jobClient] ssh client connect success\n")

	syncChan <- true

	<-ctx.Done()
}

func (c *JobClient) newSftp() {
	var err error
	if c.sftpClient == nil {
		c.sftpClient, err = c.gophClient.NewSftp()
		if err != nil {
			log.Fatal(fmt.Sprintf("[jobClient] New Sftp client failed: %v\n", err))
		}
	}
}

func (c *JobClient) SubmitCudaJob(jobUID uuid.UUID, cudaFilePath string, slurmFileContent string, objectFileName string) (jobID string, err error) {

	dirName := config.HPCJobDirPrefix + jobUID.String()
	fullDirName := config.HPCHomeDir + dirName
	cuFileName := filepath.Base(cudaFilePath)
	slurmFileName := strings.TrimSuffix(cuFileName, config.CuFileSuffix) + config.SlurmFileSuffix
	cudaFileDstPath := filepath.ToSlash(filepath.Join(fullDirName, cuFileName))
	slurmFileDstPath := filepath.ToSlash(filepath.Join(fullDirName, slurmFileName))
	objectFileDstPath := filepath.ToSlash(filepath.Join(fullDirName, objectFileName))

	res, err := c.executeCommand("mkdir " + dirName)
	if err != nil {
		log.Printf("[jobClient] SubmitCudaJob executeCommand err: %v\n", err)
		return "-1", err
	}

	err = c.upload(cudaFilePath, cudaFileDstPath)
	if err != nil {
		log.Printf("[jobClient] SubmitCudaJob upload file %v to %v err: %v\n", cudaFilePath, cudaFileDstPath, err)
		return "-1", err
	}

	err = c.CreateAndWriteFile(slurmFileDstPath, slurmFileContent)
	if err != nil {
		log.Printf("[jobClient] SubmitCudaJob writing slurm script to %v err: %v\n", slurmFileDstPath, err)
		return "-1", err
	}

	cmd := fmt.Sprintf("module load gcc/8.3.0 cuda/10.1.243-gcc-8.3.0 && nvcc %s -o %s -lcublas && cd %s && sbatch %s", cudaFileDstPath, objectFileDstPath, dirName, slurmFileDstPath)
	res, err = c.executeCommand(cmd)
	if err != nil {
		log.Printf("[jobClient] SubmitCudaJob executeCommandsAndGetLastOutput err: %v\n", err)
		return "-1", err
	}

	n, err := fmt.Sscanf(string(res), "Submitted batch job %s", &jobID)
	if err != nil || n != 1 {
		return "-1", err
	}
	return jobID, nil
}

func (c *JobClient) CheckJobFinish(jobID string) (bool, error) {
	cmd := fmt.Sprintf("sacct -j %v | tail -n +3 | awk '{print $1, $2, $3, $4, $5, $6, $7}'", jobID)
	res, err := c.executeCommand(cmd)
	if err != nil {
		log.Printf("[jobClient] CheckJobFinish executeCommand err: %v\n", err)
		return false, err
	}
	resp := string(res)
	log.Printf("[jobClient] resp: %v\n", resp)
	rows := strings.Split(resp, "\n")
	if len(rows) > 0 {
		row := rows[0]
		cols := strings.Split(row, " ")
		if len(cols) == 7 {
			if cols[5] == "COMPLETED" {
				return true, nil
			} else {
				return false, nil
			}
		}
	}
	return false, errors.New(fmt.Sprintf("Job %v not found\n", jobID))
}

func (c *JobClient) GetJobState(jobID string) (string, error) {
	cmd := fmt.Sprintf("sacct -j %s | tail -n +3 | awk '{print $1, $2, $3, $4, $5, $6, $7}'", jobID)
	res, err := c.executeCommand(cmd)
	if err != nil {
		log.Printf("[jobClient] CheckJobFinish executeCommand err: %v\n", err)
		return "", err
	}

	resp := string(res)
	log.Printf("[jobClient] resp: %v\n", resp)
	rows := strings.Split(resp, "\n")
	if len(rows) > 0 {
		row := rows[0]
		cols := strings.Split(row, " ")
		if len(cols) == 7 {
			return cols[5], nil
		}
	}
	return "MISSING", errors.New(fmt.Sprintf("Job %s not found\n", jobID))
}

func (c *JobClient) DownloadResult(jobUID uuid.UUID, localFilePath string, resultFileName string) (bool, error) {

	dirName := config.HPCJobDirPrefix + jobUID.String()
	fullDirName := config.HPCHomeDir + dirName
	resultFileOutputName := resultFileName + config.OutputFileSuffix
	resultFileErrorName := resultFileName + config.ErrorFileSuffix
	resultFileOutputDstPath := filepath.ToSlash(filepath.Join(fullDirName, resultFileOutputName))
	resultFileErrorDstPath := filepath.ToSlash(filepath.Join(fullDirName, resultFileErrorName))

	err := c.download(filepath.Join(localFilePath, resultFileOutputName), resultFileOutputDstPath)
	if err != nil {
		return false, err
	}
	err = c.download(filepath.Join(localFilePath, resultFileErrorName), resultFileErrorDstPath)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *JobClient) upload(localFilePath string, remoteFilePath string) error {
	return c.gophClient.Upload(localFilePath, remoteFilePath)
}

func (c *JobClient) download(localFilePath string, remoteFilePath string) error {
	return c.gophClient.Download(remoteFilePath, localFilePath)
}

func (c *JobClient) executeCommand(cmd string) ([]byte, error) {
	return c.gophClient.Run(cmd)
}

func (c *JobClient) createFile(fullPath string) (*sftp.File, error) {
	return c.sftpClient.Create(fullPath)
}

func (c *JobClient) writeFile(file *sftp.File, content string) error {
	_, err := file.Write([]byte(content))
	return err
}

func (c *JobClient) closeFile(file *sftp.File) error {
	return file.Close()
}

func (c *JobClient) CreateAndWriteFile(filePath string, content string) error {
	file, err := c.createFile(filePath)
	if err != nil {
		return err
	}

	err = c.writeFile(file, content)
	if err != nil {
		return err
	}

	err = c.closeFile(file)
	if err != nil {
		return err
	}

	return nil
}

func (c *JobClient) executeCommandsAndGetLastOutput(cmds []string) (out []byte, err error) {
	for _, cmd := range cmds {
		out, err = c.executeCommand(cmd)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
