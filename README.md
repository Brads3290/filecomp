# FileComp

FileComp is a recursive file hash comparer currently supporting the following hash algorithms:
 * MD5 (default)
 * SHA1
 * SHA256
 * SHA512
 
The following output formats are supported:
 * Plain text, to stdout (default)
 * XML
 
## Installation
Filecomp simply needs to be built using the Go SDK:
```$xslt
# Windows:
go build -o bin/filecomp.exe filecomp/main

# Unix
go build -o bin/filecomp filecomp/main
```
The Windows/Unix build commands are identical except for the output file name.

 
## Usage
The simplest usage is as follows:
```
filecomp -s1 /path/to/src/one -s2 /path/to/src/two
```

If the paths are directories, they will be recursively walked for files with the same path (relative to each respective source), and each file's MD5 has will be compared.

If the paths point to files, then the two files will be compared.

### Output Formats
The default output is to stdout, and looks similar to the following:
```
===== SAME FILES =====
[1] /same_one.xml
[2] /same_two.xml

===== DIFFERENT FILES =====
[1] /not_same_one.xml
```

To output to XML, use the `-o <filename>` flag. E.g. `filecomp -s1 /path/to/src/one -s2 /path/to/src/two -o results.xml`

**results.xml:**
```xml
<fileList>
  <same>
    <file relPath="/same_one.xml" hashType="md5" hash1="9da72578e88cf6a15e011d3868b39f2f" hash2="9da72578e88cf6a15e011d3868b39f2f"></file>
    <file relPath="/same_two.xml" hashType="md5" hash1="cefdcfd93de05e6a1e3aae6f376c86e4" hash2="cefdcfd93de05e6a1e3aae6f376c86e4"></file>
  </same>
  <different>
    <file relPath="/not_same_one.xml" hashType="md5" hash1="001f62e1aff6376f43f85db3b33fc684" hash2="78b9b7f9e4171bc3d16ca9831a5042ac"></file>
  </different>
</fileList>
```

### Hash Type
To change the hash type, use the `-h <hash_type>` flag. The following values can be used for `<hash_type>`:
* `md5`
* `sha1`
* `sha256`
* `sha512`

Adding support for new hash types is also possible. Simply add a case to the switch statement in `GetHashFunc` at the bottom of `main/main.go`, and return a `func() hash.Hash`, which should return your desired hash algorithm.

### Multithreading
Filecomp will run at least two threads when calculating hashes: at least one worker thread, and the main thread which will allocate files to the worker(s). To increase the number of workers, use the `-t <thread_count>` flag. If you pass in a value less than 1, **it will be ignored.**