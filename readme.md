# nyukan

入室退室管理システム

- nfcで入室、退室ができる
- 施設の入り口に置いてある
- raspberry pi

# setup

```bash
git clone <this repo link> /home/finyl/nyukan
cd /home/finyl/nyukan
go build -o main main.go
// TODO: enable nyukan systemd service (see ./nyukan.service)
```

# ssh

ssh finyl@192.168.2.105
