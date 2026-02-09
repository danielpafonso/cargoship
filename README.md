# CargoShip

A suite of applications that Extract, Process, and Send files to and from FTP and SFTP servers

CargoShip consists of three components that can work independently or together to create robust file processing pipelines:

| Components                         |                                                     |
| ---------------------------------- | --------------------------------------------------- |
| [Shipper](cmd/shipper/README.md)   | Download and upload files to/from (S)FTP servers    |
| [Loader](cmd/loader/README.md)     | Compress and remove files from local storage        |
| [Packager](cmd/packager/README.md) | Apply file processors and generate new parsed files |

## Timestamp formatting

Configurations files support dynamic timestamps that are replaced at runtime allowing creation of time-stamped log files or output files.

CargoShip uses Go's [time package](https://pkg.go.dev/time#pkg-constants) formatting convention.

For the more used formats see the table below:

| Time Part | Time format | Value |
| --------- | ----------- | ----- |
| Year      | yyyy - 2020 | 2006  |
| Month     | mm - 12     | 01    |
| Day       | dd - 23     | 02    |
| Hours     | HH - 14     | 15    |
| Minutes   | MM - 59     | 04    |
| Seconds   | SS - 45     | 05    |

# Common Configuration Fields

### maxTime

Defines the maximum time window for processing files. When processing a batch of files, CargoShip identifies the timestamp of the first valid file and creates a time window extending `maxTime` minutes forward from that point. Only files within this window will be processed.

### windowLimit

Defines how far back from the current time to look for files. This creates a time boundary by subtracting `windowLimit` minutes from the current date/time. Files newer than this boundary will be excluded from processing.

### Example

```yaml
maxTime: 120 # 2-hour processing window
windowLimit: 15 # Skip files from last 15 minutes
```
