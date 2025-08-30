# Packager

Script that process source files applying some command and generate a new concatenation file

## Configuration

| Field Name                   | Type    | Description                                                                                                                                          |
| ---------------------------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| log2console                  | boolean | Flag indicating if the logging should be duplicated to the console                                                                                   |
| manifest                     | string  | File Path to the Times File                                                                                                                          |
| logging                      | object  |                                                                                                                                                      |
| &ensp; script                | string  | Path to Script Logging, which can have a dynamic timestamp, see [Dynamic Timestamp](../../README.md#dynamic-timestamp) for more information          |
| &ensp; files                 | string  | Path to processed files Logging, which can have a dynamic timestamp, see [Dynamic Timestamp](../../README.md#dynamic-timestamp) for more information |
| services                     | array   |                                                                                                                                                      |
| &ensp; name                  | string  | Service identifier name                                                                                                                              |
| &ensp; enable                | boolean | Flag to enable service to run                                                                                                                        |
| &ensp; mode                  | string  | Mode of processing files                                                                                                                             |
| &ensp; cmd                   | string  | Command to run, on each file                                                                                                                         |
| &ensp; sourceFolder          | string  | Source Folder                                                                                                                                        |
| &ensp; filePrefix            | string  | Source File prefix, used for filtering files list                                                                                                    |
| &ensp; fileExtension         | string  | Source File extension, used for filtering files list                                                                                                 |
| &ensp; destinationFolder     | string  | Destination Folder                                                                                                                                   |
| &ensp; destinationFile       | string  | Destination processed file name format                                                                                                               |
| &ensp; destinationDateFormat | string  | Destination processed file date format                                                                                                               |
| &ensp; historyFolder         | string  | History Folder where original files are stored after processing                                                                                      |
| &ensp; newline               | string  | Add newline after processing files                                                                                                                   |
| &ensp; maxTime               | int     | Max time (minutes) windows of files to process, see [Time Windows](../../README.md#time-windows) for more information                                |
| &ensp; windowLimit           | int     | Limit (minutes) in relation to NOW where newer files won't be process, see [Time Windows](../../README.md#time-windows) for more information         |

## Destination file placeholders

| Placeholder | Description                                                             |
| ----------- | ----------------------------------------------------------------------- |
| {date}      | Timestamp in UTC with format described in field `destinationDateFormat` |
| {files}     | Number of file process in the run                                       |
