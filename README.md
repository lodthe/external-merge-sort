# External merge sort

An implementation of external merge sort algorithm (with 2-way merge) written in Go. This particular implementation sorts strings. The implementation is located [here](./pkg/algo/merge_sort.go).

As it's implementation of an algorithm in external memory, we are interested in disk usage. If the size of main memory is M, and we read/write B bytes at a time, 

The tool is not for production usage. If you want to sort huge files with high performance, probably, you want to find an implementation written in C or C++ (or any other programming language without a garbage collector).

## Usage

The next example shows how to generate a huge file, sort it and validate the result:
```bash
make all

# Generate 42000000 lines. 
# Allowed line lengths: 126-128. (~5GB file will be created)
# Each line will have the same prefix of length 110.
./bin/generator -min-length 126 -max-length 128 -equal-prefix-length 110 -alphabet hex -output input.txt -count 42000000

# Sort input.txt and save the result in output.txt.
# Main memory limit: 500MB. Block size: 1MB.
# (in fact, more memory may be allocated due to the presence of a garbage collector)
./bin/sort -memory 500000000 -blocksize 1000000 -input input.txt -output output.txt

# Validate the result.
./bin/validator -input input.txt -output output.txt
```

Use `--help` for more detail.

## Tools

There are three useful tools in this repository.

Run make to build them:
```bash
make all
```

After that, three executable files will be created in the `bin` directory. Use `--help` to list supported arguments.

### Generator

Generator produces files with random tokens (strings). You can specify count of tokens, their length and a set of allowed characters.

```text
Usage of ./bin/generator:
  -alphabet string
        You can specify one of the supported sets of characters used for generation:
        binary - 01;
        lower - abc..xyz;
        upper - ABC...XYZ;
        numbers - 012345689;
        alnum - abc...xyz0123456789;
        hex - 0123456789ABCDEF;
        non-space - ASCII code in range (32, 128).
        
        If the specified value doesn't match with any of the predefined alphabet names, this value will be used as the set of characters. (default "lower")
  -count int
        How many tokens must be generated?. (default 1000)
  -delimiter string
        A character used to separate tokens. (default "\n")
  -equal-prefix-length int
        All tokens will begin with the same prefix, and you can specify the length of this prefix.
  -max-length int
        Maximum allowed token length. (default 128)
  -min-length int
        Minimum allowed token length. (default 128)
  -output string
        Where the generated data will be saved. (default "input.txt")
```

### Sort

Sort is an implementation of external merge sort algorithm. 

Choose memory limit and block size according to your setup and limitations.

```text
Usage of ./bin/sort:
  -blocksize int
        Size of one block (in bytes). (default 1048576)
  -delimiter string
        A character used to separate tokens. (default "\n")
  -input string
        Input file path. (default "input.txt")
  -memory int
        The algorithm will use at most O(memory) main memory. (default 536870912)
  -order string
        Sort order. Supported values: ASC, DESC. (default "ASC")
  -output string
        Output file path. (default "output.txt")
  -tempdir string
        Where temporary files can be created. If you use /tmp, make sure there is enough space for two copies of the input file. (default ".")
```

### Validator

Validator helps to check whether the implementation is correct. It takes input and output files and makes the following checks:
- the number of tokens matches;
- hashes of the set of tokens are the same;
- order of tokens in the output file is correct.

```text
Usage of ./bin/validator:
  -delimiter string
        A character used to separate tokens. (default "\n")
  -input string
        Input file path. (default "input.txt")
  -order string
        Sort order. Supported values: ASC, DESC. (default "ASC")
  -output string
        Output file path. (default "output.txt")
```