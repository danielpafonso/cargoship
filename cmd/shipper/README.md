# Shipper

Transfers files to and from FTP and SFTP servers. It supports both importing (downloading) and exporting (uploading) files with flexible filtering and time-based selection

- **Multi-Protocol Support**: Connect to FTP and SFTP servers
- **Bidirectional Transfer**: Download (import) or upload (export) files
- **Pattern Matching**: Filter files by prefix and extension
- **Time-based Selection**: Transfer files within specific time windows
- **History Tracking**: Archive transferred files to prevent reprocessing
- **Multi-Server Support**: Configure multiple servers and services

## Configuration

| Field Name               | Type         | Description                                                                                                     |
| ------------------------ | ------------ | --------------------------------------------------------------------------------------------------------------- |
| log2console              | boolean      | Flag indicating if the logging should be duplicated to the console                                              |
| timesFilePath            | string       | File path to the manifest file that tracks processing state and timestamps                                      |
| logging                  | object       |                                                                                                                 |
| &ensp; script            | string       | Path to script execution log file, supports [dynamic timestamps](../../README.md#timestamp-formatting)          |
| &ensp; files             | string       | Path to processed files log, supports [dynamic timestamps](../../README.md#timestamp-formatting)                |
| ftps                     | array        |                                                                                                                 |
| &ensp; name              | string       | Server identifier name                                                                                          |
| &ensp; hostname          | string       | Server hostname or IP address                                                                                   |
| &ensp; port              | int          | Server port number                                                                                              |
| &ensp; username          | string       | Username for server authentication                                                                              |
| &ensp; password          | string       | Password for server authentication                                                                              |
| &ensp; protocol          | string       | Transfer protocol: `"ftp"` or `"sftp"`                                                                          |
| services                 | array        |                                                                                                                 |
| &ensp; name              | string       | Service identifier name                                                                                         |
| &ensp; enable            | boolean      | Flag to enable the service to run                                                                               |
| &ensp; ftpConfig         | string array | List of server names (from `ftpServers`) to execute this service against.                                       |
| &ensp; sourceFolder      | string       | For `import`: remote server path. For `export`: local filesystem path                                           |
| &ensp; destinationFolder | string       | For `import`: local filesystem path. For `export`: remote server path                                           |
| &ensp; filePrefix        | string       | Filter files by prefix (e.g., `"data_"` matches `data_2024.txt`). Leave empty to match all files                |
| &ensp; fileExtension     | string       | Filter files by extension (e.g., `".txt"`, `".csv"`). Leave empty to match all extensions                       |
| &ensp; historyFolder     | string       | Path to archive directory where processed files are moved after successful transfer                             |
| &ensp; maxTime           | int          | Maximum time window in minutes for transferring files. See [maxTime](../../README.md#maxtime) for details       |
| &ensp; windowLimit       | int          | Exclude files newer than this many minutes from now. See [windowLimit](../../README.md#windowlimit) for details |
