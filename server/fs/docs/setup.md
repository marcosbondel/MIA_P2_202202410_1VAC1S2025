```shell
    go mod init MIA_P2_202202410_1VAC1S2025
    go run ./main.go
    go mod tidy # solve dependencies
```
### Generate binary
```shell
    go build -o <myapp>
```

### Connect to the EC2 Instance
```shell
    ssh -i mia_202202410.pem ubuntu@3.140.85.187
```

### Systemd service
```bash
    sudo nano /etc/systemd/system/myapp.service

    # Enable and start service
    sudo systemctl daemon-reexec
    sudo systemctl daemon-reload
    sudo systemctl enable myapp.service
    sudo systemctl start myapp.service

    # Check logs and status
    sudo systemctl status myapp.service
    journalctl -u myapp.service -f
```