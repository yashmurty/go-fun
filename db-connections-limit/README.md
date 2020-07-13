## Check connections limitation on RDS DBs.

Simulate simultaneous MySQL Connections & test its limitation.

- Write program that opens multiple connections on MySQL Database.
- Confirm on Database side whether it shows multiple connections are open.

### DB Access Program

**Start the program**:
```sh
# Load the env variables.
bash load_env.sh

# Run the go program.
go run main.go
```

**Note:** Go mod will automatically install missing dependencies.