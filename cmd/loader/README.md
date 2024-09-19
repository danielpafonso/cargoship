# Loader

Script that compress and clean/delete files from local file system

## Configuration

| Field Name               | Type    | Description                                                                                                                                          |
| ------------------------ | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| log2console              | boolean | Flag indicating if the logging should be duplicated to the console                                                                                   |
| manifest                 | string  | File Path to the Times File                                                                                                                          |
| logging                  | object  |                                                                                                                                                      |
| &ensp; script            | string  | Path to Script Logging, which can have a dynamic timestamp, see [Dynamic Timestamp](../../README.md#dynamic-timestamp) for more information          |
| &ensp; files             | string  | Path to processed files Logging, which can have a dynamic timestamp, see [Dynamic Timestamp](../../README.md#dynamic-timestamp) for more information |
| services                 | array   |                                                                                                                                                      |
| &ensp; name              | string  | Service identifier name                                                                                                                              |
| &ensp; mode              | string  | List of services to execute againts                                                                                                                  |
| &ensp; sourceFolder      | string  | Source Folder                                                                                                                                        |
| &ensp; destinationFolder | string  | Destination Folder, only aplicable on compress mode                                                                                                  |
| &ensp; filePrefix        | string  | File prefix to filter source files                                                                                                                   |
| &ensp; fileExtension     | string  | File extention to filter source files                                                                                                                |
| &ensp; maxTime           | int     | Max time (minutes) windows of files to process, see [Time Windows](../../README.md#time-windows) for more information                                |
| &ensp; windowLimit       | int     | Limit (minutes) in relation to NOW where newer files won't be process, see [Time Windows](../../README.md#time-windows) for more information         |
