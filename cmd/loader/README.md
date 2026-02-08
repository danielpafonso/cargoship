# Loader

Compresses and cleans up files from the local file system. It supports two modes of operation: archiving files into compressed archives, and deleting files after they've been processed.

- **Compress Mode**: Create archives from source files matching specific patterns
- **Clean Mode**: Delete files from specified directories
- **Time-based Filtering**: Process only files within specified time windows
- **Manifest Tracking**: Track processed files to avoid reprocessing

## Configuration

| Field Name               | Type    | Description                                                                                                     |
| ------------------------ | ------- | --------------------------------------------------------------------------------------------------------------- |
| log2console              | boolean | Flag indicating if the logging should be duplicated to the console                                              |
| manifest                 | string  | File path to the manifest file that tracks processing state and timestamps                                      |
| logging                  | object  |                                                                                                                 |
| &ensp; script            | string  | Path to script execution log file, supports [dynamic timestamps](../../README.md#timestamp-formating)           |
| &ensp; files             | string  | Path to processed files log, supports [dynamic timestamps](../../README.md#timestamp-formating)                 |
| services                 | array   |                                                                                                                 |
| &ensp; name              | string  | Service identifier name                                                                                         |
| &ensp; enable            | boolean | Flag to enable service to run                                                                                   |
| &ensp; mode              | string  | Operation mode: `"compress"` to archive files, `"clean"` to delete files                                        |
| &ensp; sourceFolder      | string  | Path to the source directory containing files to process                                                        |
| &ensp; destinationFolder | string  | Path to destination directory for compressed archives (required for compress mode)                              |
| &ensp; filePrefix        | string  | Filter files by prefix (e.g., `"data_"` matches `data_2024.txt`). Leave empty to match all files                |
| &ensp; fileExtension     | string  | Filter files by extension (e.g., `".txt"`, `".log"`). Leave empty to match all extensions                       |
| &ensp; maxTime           | int     | Maximum time window in minutes for processing files. See [maxTime](../../README.md#maxtime) for details         |
| &ensp; windowLimit       | int     | Exclude files newer than this many minutes from now. See [windowLimit](../../README.md#windowlimit) for details |

### Mode Options

- **`compress`**: Archives matching files into a `.tar.gz` file in the destination folder
- **`clean`**: Deletes matching files from the source folder
