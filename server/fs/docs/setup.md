```shell
    go mod init MIA_P2_202202410_1VAC1S2025
    go run ./main.go
    go mod tidy # solve dependencies
```

### Connect to the EC2 Instance
```shell
    ssh -i mia_202202410.pem ubuntu@18.226.34.197
```