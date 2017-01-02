# pipelines
This project contains some Go concurrent pipeline examples based on this [article](https://blog.golang.org/pipelines).

A pipeline is a kind of concurrent program. It is a series of stages connected by channels, where each stage is a group of goroutines running the same function. In each stage, the goroutines:

* Receive values from upstream via inbound channels
* Perform some function on that data, usually producing new values
* Send values downstream via outbound channels

Each stage has any number of inbound and outbound channels, except the first and last stages, which have only outbound or inbound channels, respectively. The first stage is sometimes called the _source_ or _producer_; the last stage, the _sink_ or _consumer_.

## Getting Started

```sh
$ go get -u github.com/ihcsim/pipelines
$ go test -v -cover -race ./...
```

`square`: A pipeline that extracts the list of integers into a channel, and perform square operation on each integer.
`md5dir`: A pipeline that compute the MD5 checksum of all the files in a give directory.

* `md5sum/serial`: Computes the MD5 checksum of files in a folder in sequential lexical order.
* `md5sum/parallel`: Computes the MD5 checksum of files in a folder using a two-stage pipeline.
* `md5sum/bounded`: Compute the MD5 checksum of files in a folder using a three-stage pipeline where the number of goroutines are bounded.

## LICENSE
MIT. Refer to the [LICENSE](LICENSE) file.
