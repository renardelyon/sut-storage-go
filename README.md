# sut-storage-go
Microservice for managing user uploaded files

## How to run in local?

1. Create file env and name it `dev.env`. its content can be seen in code block below. 
```
PORT=:50053
DB_URL=
SERVICE_ACCOUNT_PATH=
TOKEN_PATH=
FOLDER_ID=
CLIENT_ID=
CLIENT_SECRET=
GDRIVE_API_REFRESH_TOKEN=
```

2. Execute command below
```
make init # initialize go.mod
make tidy # Tidy up go module
```

3. Adding go bin into path env variables
```
export PATH=$PATH:$(go env GOPATH)/bin
```

4. Adding folder with `pb` as name into ther project root directory

5. Generate protobuf by executing command below
```
make proto-gen
```

6. Run the application
```
make run
```
