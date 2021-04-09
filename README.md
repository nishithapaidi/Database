# Database
# Code to build database and build it into a webtool


## Building a single binary containing templates

This is a complete example to create a single binary with the
[gin-gonic/gin][gin] Web Server with HTML templates.

[gin]: https://github.com/gin-gonic/gin

## How to use

### Prepare Packages

```
go get github.com/gin-gonic/gin
go get github.com/jessevdk/go-assets-builder
```

### Generate assets.go

```
go-assets-builder html -o assets.go
```

### Quick Run

```
go run . > output.txt
```

### Build the server

```
go build -o assets-in-binary
```

### Run

```
./assets-in-binary
```

### Installing Service

You'll need to configure the file_config.go.example to file_config.go. The same is with bin/database.config.example.

Add data files to bin folder.

```
ln bin/morrislab.service /etc/systemd/system/
```

```
systemctl daemon-reload
```

```
systemctl start morrislab.service
```

## Adding New Entries
Data in the database is split between mRNA, LncRna, and human.
The scripts that start with populate... are used to populate the database.
* populateRnaDb.go takes on the input from the mouse xlsv excel sheet. The sample id is the column store in gene.sample
and the the value of the sample is stored in gene.trial. The purpose for separating the sample and trial tables is to
allow customizable columns. At the time, I thought this was a good idea, however, looking back it makes the solution
far more complicated than it should be.
* populateMouseInfo.go takes care of the mouse info xlsv excel sheet. It makes use of the gene.mouse_info table which is
relatively straightforward.
* populateHuman.go takes care of the human data set. It makes use of the gene.human_gene_data table and is relatively 
straightforward in implementation as with mouse info.
