## Check connections limitation on RDS Databases.

Simulate simultaneous MySQL Connections & test its limitation.

- Write program that opens multiple connections on MySQL Database.
- Keep the connections open instead of closing them after the query.
- Confirm on Database side whether it shows multiple connections are open.
- Check how many such multiple connections are possible to be kept open.

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
db.OpenConnections :  40

ERROR :  Error 1040: Too many connections
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

Same as above for the DatabaseConnections metric, i.e. Threads_connected = 45.

When we see our programs output logs we see that:

```sh
# Output logs from our go program.
db.OpenConnections :  100
...
Total errorCount:  0 # This is a total count of how many goroutines i.e. our concurrent connection requests failed.
```

When we run the program with 1000 concurrent connection requests, the stats are:

Same as above for the DatabaseConnections metric, i.e. Threads_connected = 45.

When we see our programs output logs we see that:

```sh
# Output logs from our go program.
db.OpenConnections :  1000

ERROR :  Error 9501: Timed-out waiting to acquire database connection
...
Total errorCount:  818 # This is a total count of how many goroutines i.e. our concurrent connection requests failed.
```

If we look up this error code in AWS, it states the cause to be:
```
ERROR 9501 (HY000): Timed-out waiting to acquire database connection
-
The proxy timed-out waiting to acquire a database connection. Some possible reasons include the following:
- The proxy is unable to establish a database connection because the maximum connections have been reached
- The proxy is unable to establish a database connection because the database is unavailable.
```

Does this imply that for our `db.t2.small instance` which has a set value of `max_connections` to be 45, can have it's RDS Proxy handle a maximum of around `~180 (~4 times the max_connections value)` concurrent connection requests?

Funnily, when we run the program with 200 concurrent connection requests, the stats are:

```sh
# Output logs from our go program.
db.OpenConnections :  200

ERROR :  Error 9501: Timed-out waiting to acquire database connection
...
Total errorCount:  54 # This is a total count of how many goroutines i.e. our concurrent connection requests failed.
```

This means that this time `~146` connections were successful.

In either cases, we were able to achieve successful request results from RDS Proxy at much higher scale than directly connecting to the database instance.

Another important note that can be added is, currently in the above test scenarios, we keep the connection open for a long time (kept open for 30 seconds).
If we change it to shorter durations, say 3 seconds, we note the results to be:

For the case of 200 concurrent requests:
```sh
# Output logs from our go program.
db.OpenConnections :  200

Total errorCount:  0 # This is a total count of how many goroutines i.e. our concurrent connection requests failed.
```

For the case of 1000 concurrent requests:
```sh
# Output logs from our go program.
db.OpenConnections :  1000

Total errorCount:  0 # This is a total count of how many goroutines i.e. our concurrent connection requests failed.
```

We see that both the tests ran successfully! RDS Proxy was able to manage this huge number of concurrent requests even for a database with such a minimal spec (db.t2.small = 1 vCPU, 2 GB Mem). Had we run the same test directly against the database, we would have immediately run into the error: `Error 1040: Too many connections`.

Another point to note is, it would be a great idea to use RDS Proxy when serving serverless programs, since they tend to frequently open and close database connections. RDS Proxy would take away connection management from the database and do it by itself, hence wasting less memory resources of the database.