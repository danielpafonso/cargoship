# Packager

Processes files by applying commands to each file and concatenating the results into a single output file.

- **Command Execution**: Run custom commands on each matching file
- **File Concatenation**: Combine processed results into a single output file
- **Pattern Matching**: Filter files by prefix and extension
- **Time-based Processing**: Process files within specific time windows
- **History Tracking**: Archive original files after processing
- **Dynamic Naming**: Generate output filenames with timestamps and file counts

## Configuration

| Field Name                   | Type    | Description                                                                                                            |
| ---------------------------- | ------- | ---------------------------------------------------------------------------------------------------------------------- |
| log2console                  | boolean | Flag indicating if the logging should be duplicated to the console                                                     |
| manifest                     | string  | File path to the manifest file that tracks processing state and timestamps                                             |
| logging                      | object  |                                                                                                                        |
| &ensp; script                | string  | Path to script execution log file, supports [dynamic timestamps](../../README.md#timestamp-formating)                  |
| &ensp; files                 | string  | Path to processed files log, supports [dynamic timestamps](../../README.md#dtimestamp-formating)                       |
| services                     | array   |                                                                                                                        |
| &ensp; name                  | string  | Service identifier name                                                                                                |
| &ensp; enable                | boolean | Flag to enable service to run                                                                                          |
| &ensp; mode                  | string  | Processing mode: `"command"` to run commands on files, `"concat"` to concatenate files without commands                |
| &ensp; cmd                   | string  | Shell command to execute on each file. Use `{file}` as placeholder for the filename. Only applicable in `command` mode |
| &ensp; sourceFolder          | string  | Path to the source directory containing files to process                                                               |
| &ensp; filePrefix            | string  | Filter files by prefix (e.g., `"data_"` matches `data_2024.txt`). Leave empty to match all files                       |
| &ensp; fileExtension         | string  | Filter files by extension (e.g., `".txt"`, `".json"`). Leave empty to match all extensions                             |
| &ensp; destinationFolder     | string  | Path to destination directory for the generated output file                                                            |
| &ensp; destinationFile       | string  | Output filename template                                                                                               |
| &ensp; destinationDateFormat | string  | Timestamp format for `{date}` placeholder                                                                              |
| &ensp; historyFolder         | string  | Path to archive directory where original source files are moved after successful processing                            |
| &ensp; newline               | string  | Line separator to add between processed files                                                                          |
| &ensp; maxTime               | int     | Maximum time window in minutes for processing files. See [maxTime](../../README.md#maxtime) for details                |
| &ensp; windowLimit           | int     | Exclude files newer than this many minutes from now. See [windowLimit](../../README.md#windowlimit) for details        |

### Mode Options

- **`command`**: Executes the specified `cmd` on each file and concatenates the command output
- **`concat`**: Concatenates file contents directly without running commands

## Destination File Placeholders

Use these placeholders in the `destinationFile` field:

| Placeholder | Description                                                                  |
| ----------- | ---------------------------------------------------------------------------- |
| `{date}`    | Replaced with timestamp in UTC using the format from `destinationDateFormat` |
| `{files}`   | Replaced with the number of files processed in the current run               |
