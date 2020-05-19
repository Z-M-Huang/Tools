# Tools
Tools for fun

![Docker Release CI](![Docker Release CI](https://github.com/Z-M-Huang/Tools/workflows/Docker%20Release%20CI/badge.svg))

# For future developers who found this interesting.
Please add environment variables to have this repo run properly.
```
"env": {
  "DEBUG": "1",
  "HOST": "localhsot",
  "JWT_KEY": "aLNDHKx6NftYrKsPs4GAdQc3ugCjrLbh",
  "GOOGLE_CLIENT_ID": "not empty sample string",
  "GOOGLE_CLIENT_SECRET": "not empty sample string",
  "REDIS_ADDR": "ip:port",
  "CONNECTION_STRING": "./db.db",
  "DB_DRIVER": "sqlite3"
}
```
Some variables cannot be empty but set a random string would help you start faster. Also, JWT_KEY can be random string for you to start. And I hope you change it.