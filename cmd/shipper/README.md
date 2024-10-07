# Shipper
Scripts to extract and sends files to ftps and sftp servers

## Configuration
| Field Name               | Type          | Description                                                                                                                                          |
|--------------------------|---------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| log2console              | boolean       | Flag indicating if the logging should be duplicated to the console                                                                                   |
| timesFilePath            | string        | File Path to the Times File                                                                                                                          |
| logging                  | object        |                                                                                                                                                      |
| &ensp; script            | string        | Path to Script Logging,  which can have a dynamic timestamp, see [Dynamic Timestamp](../../README.md#dynamic-timestamp) for more information         |
| &ensp; files             | string        | Path to processed files Logging, which can have a dynamic timestamp, see [Dynamic Timestamp](../../README.md#dynamic-timestamp) for more information |
| ftps                     | array         |                                                                                                                                                      |
| &ensp; name              | string        | Server identifier name                                                                                                                               |
| &ensp; hostname          | string        | Server Hostname                                                                                                                                      |
| &ensp; port              | int           | Server port                                                                                                                                          |
| &ensp; username          | string        | Username to authenticate on the server                                                                                                               |
| &ensp; password          | string        | Password to authenticate on the server                                                                                                               |
| &ensp; protocol          | string        | Server trasferer protocol: ftp, sftp                                                                                                                 |
| services                 | array         |                                                                                                                                                      |
| &ensp; name              | string        | Service identifier name                                                                                                                              |
| &ensp; enable            | boolean       | Flag to enable the servic ro run                                                                                                                     |
| &ensp; ftpConfig         | string array  | List of servers to run the services againts                                                                                                          |
| &ensp; sourceFolder      | string        | Source Folder                                                                                                                                        |
| &ensp; destinationFolder | string        | Destination Folder                                                                                                                                   |
| &ensp; filePrefix        | string        | File prefix to filter source files                                                                                                                   |
| &ensp; fileExtension     | string        | File extention to filter source files                                                                                                                |
| &ensp; historyFolder     | string        | Folder to move/archive process files                                                                                                                 |
| &ensp; maxTime           | int           | Max time (minutes) windows of files to process, see [Time Windows](../../README.md#time-windows) for more information                                |
| &ensp; windowLimit       | int           | Limit (minutes) in relation to NOW where newer files won't be process, see [Time Windows](../../README.md#time-windows) for more information         |
