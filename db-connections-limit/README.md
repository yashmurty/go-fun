## Check connections limitation on RDS Databases.

Simulate simultaneous MySQL Connections & test its limitation.

- Write program that opens multiple connections on MySQL Database.
- Confirm on Database side whether it shows multiple connections are open.

### Connection Pool in Go

To understand how it works under the hood, refer to these links:
- http://go-database-sql.org/accessing.html
- http://go-database-sql.org/connection-pool.html

### Run the program

**Start the program**:
```sh
# Load the env variables.
bash load_env.sh

# Run the go program.
go run main.go
```

**Note:** Go mod will automatically install missing dependencies.

## Connections Limit - Test Results

### AWS RDS Aurora (db.t2.small instance)

This instance type has default value for `max_connections` parameter set as 45.
https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Managing.Performance.html

Before performing the limits test, the status is:

```sh
# This is equivalent to the DatabaseConnections metric.
mysql> show status where `variable_name` = 'Threads_connected';
+-------------------+-------+
| Variable_name     | Value |
+-------------------+-------+
| Threads_connected | 5     |
+-------------------+-------+
1 row in set (0.20 sec)
```

Out of these 5 connections, 4 belong to system process by AWS RDS (this number can vary at times). 1 belongs to our MySQL Client that we are using to monitor these stats ourselves.

When we run the program with 50 concurrent connection requests, the stats are:

```sh
# This is equivalent to the DatabaseConnections metric.
mysql> show status where `variable_name` = 'Threads_connected';
+-------------------+-------+
| Variable_name     | Value |
+-------------------+-------+
| Threads_connected | 45    |
+-------------------+-------+
1 row in set (0.21 sec)
```

As expected, it matches the `max_connections` parameter value of 45.

Out of these 45, 5 are the same as stated above. So the remaining 40 belong to our program which tried to establish 50 concurrent connections.
This shows that 10 must have failed to connect.

We can confirm this by checking our console output, where we see logs like:

```sh
# Output logs from our go program.
ERROR :  Error 1040: Too many connections
db.OpenConnections :  40
...
Total errorCount:  10 # This is a total count of how many goroutines i.e. our concurrent connection requests failed.
```

> One easy way to avoid running into such limitations in `go` would be to set the `db.SetMaxOpenConns(N)` value, where N could be 40 to stay under the max limit of 45. This will make sure that other concurrent requests don't open any connection when already 40 are open, and they wait for connections to become available again.

### AWS RDS Proxy (for db.t2.small instance)

> **NOTE:** RDS Proxy is not accessible publicly. We RDS Proxy endpoint resolves only within the VPC, unlike the RDS instance itself, which can be set to be accessible publicly.

Before performing the limits test, the status is:

```sh
# This is equivalent to the DatabaseConnections metric.
mysql> show status where `variable_name` = 'Threads_connected';
+-------------------+-------+
| Variable_name     | Value |
+-------------------+-------+
| Threads_connected | 7     |
+-------------------+-------+
1 row in set (0.36 sec)
```

Out of these 5 connections, 4 belong to system process by AWS RDS (this number can vary at times). 1 belongs to our MySQL Client that we are using to monitor these stats ourselves. 2 new connections belong to `rdsproxyadmin`.

When we run the program with 50 concurrent connection requests, the stats are:

```sh
# This is equivalent to the DatabaseConnections metric.
mysql> show status where `variable_name` = 'Threads_connected';
+-------------------+-------+
| Variable_name     | Value |
+-------------------+-------+
| Threads_connected | 45    |
+-------------------+-------+
1 row in set (0.19 sec)
```

As expected, it matches the `max_connections` parameter value of 45.
This shows that 10 must have failed to connect, but...

But unlike before, now we are using a `RDS Proxy` to connect to the database.
When we see our programs output logs we see that:

```sh
# Output logs from our go program.
db.OpenConnections :  50
...
Total errorCount:  0 # This is a total count of how many goroutines i.e. our concurrent connection requests failed.
```

As we can see, the `RDS Proxy` flawlessly handled the connection pooling and we did not run into the `too many open connections` issue like earlier. 
Please note that we were able to request a higher number of concurrent connection requests than the permissible `max_connections` value of 45. This is because under the hood `RDS Proxy` will automatically queue any requests higher than the max limit and wait for connections to be free again before resolving it.

When we run the program with 100 concurrent connection requests, the stats are:
