# CargoShip

Applications that Extract, Process, and Send files to and from FTP and SFTP servers

| Components                         |                                                     |
| ---------------------------------- | --------------------------------------------------- |
| [Shipper](cmd/shipper/README.md)   | Download and upload files to/from (S)FTP servers    |
| [Loader](cmd/loader/README.md)     | Compress and remove files from local storage        |
| [Packager](cmd/packager/README.md) | Apply file processors and generate new parsed files |

## Timestamp Formatting

In configurations the files can be configured with a dynamic timestamp that is replaced on file creation.

Since this project uses golang the timestamp formatting is the same as golang's [time package](https://pkg.go.dev/time#pkg-constants).

For the more used formats see the table below:

| Time Part | Time format | Value |
| --------- | ----------- | ----- |
| year      | yyyy - 2020 | 2006  |
| month     | mm   - 12   | 01    |
| day       | dd   - 23   | 02    |
| hours     | HH   - 14   | 15    |
| minutes   | MM   - 59   | 04    |
| seconds   | SS   - 45   | 05    |

## Building the Project

To build all applications:

```bash
make build
```

To build individual applications:

```bash
make shipper   # Build only the shipper application
make loader    # Build only the loader application
make packager  # Build only the packager application
```

Built binaries will be placed in the `build/` directory along with their configuration files.

## Common Configuration Fields

### maxTime

Maximum time window (in minutes) for processing files. This creates a time limit by using the first valid file to download and adding minutes equal to the maxTime value. Files outside this window will not be processed.

### windowLimit

Time limit (in minutes) relative to the current date/time. Files newer than this window will not be processed. This is calculated by subtracting minutes equal to windowLimit value from current date.
