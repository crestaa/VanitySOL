# VanitySOL
_Vanity wallet generator on the Solana blockchain_

This project was built using Go v1.22.2

## First setup
Use `go mod tidy` to setup the project while being inside the main folder.

## Usage
Run the project with `go run .` followed by the next flags (optional):
- `-s string`: a sequence of letters (upper/lower case) that defines what your address is going to start with. Leave empty if you don't want to filter this (default: empty string);
- `-e string`: a sequence of letters (upper/lower case) that defines what your address is going to end with. Leave empty if you don't want to filter this (default: empty string);
- `-c`: use if you want your generation to be case sensitive (default: false);
- `-t number`: the number of threads you want to use for the address generation (default: half of your logical cores).

_i.e.: `go run . -s van -e SOL -c -t 15` will start generating addresses starting with "van" and ending with "SOL" (case sensitive) with 15 threads._

The process of generating a vanity address can be time and hardware intense if you try to match more than 4 characters, growing exponentially. Case sensitivity adds another layer of complexity to the generation.
