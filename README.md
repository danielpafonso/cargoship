# CargoShip

Applications that Extract, Process, and Send files to and from FTP and SFTP servers

| Components ||
| --- | --- |
| Shipper | Download and upload files to/from (S)FTP servers |
| Loader | Compress and remove files from local storage |
| ? | Apply file processors and generate new parsed files |

## Timestamp formating

In configurations the files can be configurated with a dynamic timestamp that is replaced on file creation.

Since this project uses golang the timestamp formating is the same as golang's [time package](https://pkg.go.dev/time#pkg-constants).

For the more used formats see the table below:

| Time Part | Value |
|-----------|-------|
| year      | 2006  |
| month     | 01    |
| day       | 02    |
| hours     | 15    |
| minutes   | 04    |
| seconds   | 05    |


# Common Configuration Fields

### maxTime

> _add more info_

Time limit calculating by using the first valid file to download and add minutes equal to maxTime value

### windowLimit

> _add more info_

Time limit calculated by substratcing minutes equal to windowLimit value to current date
