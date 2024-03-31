# rdetective

A rolling hash algorithm to generate a signature and delta between two files.

## Dependencies
- [golang 1.17](https://golang.org/dl/)

## Build
To build the project simply run `./build.sh`.

## Run
To start rdetective run the following command:

```bash
./bin/rdetective diff --original path/to/original_file --updated path/to/updated_file [flags]
```

Use `--help` for all available flags:

```bash
./bin/rdetective diff --help
```

## Caveats
- rdetective is only using a `weak` algorithm ([adler32](https://en.wikipedia.org/wiki/Adler-32)) for the sake of exercise and simplicity. In the real world a combination of a week and a strong algorithm (like `sha1`) would be prefered to avoid collisions.
The weak hash function would typically be used for efficiency purposes, as it's faster to compute but may have a higher chance of collisions. The strong hash function, on the other hand, is slower but provides a more reliable and unique hash value. Adding a `strong` algorithm to rdetective's data structures should be fairly trivial.
- A patch function is not implemented
- rdetective prints out the differences found relative to the signature. A more human readable way would be to display the differences using the data of the original file and not the chunks.
