# Search And Sort Movies or Series

## Service example
```
[Unit]
Description = Search and Sort Movies Service
After = network.target

[Service]
WorkingDirectory=/media/hdd/app
ExecStart =/bin/bash -c "/media/hdd/app/search-and-sort-movies-linux-amd64 -scan"
User=root
Group=root
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```