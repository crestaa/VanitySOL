# VanitySOL
_Vanity wallet generator on the Solana blockchain_

This project was built using Go v1.22.2

## First setup
Use `go mod tidy` to setup the project while being inside the main folder.

## Usage
Run the project with `go run .` while being inside the project folder: this will open a CLI for further inputs.
The following inputs are needed:
- **start of the address**: a sequence of letters (upper/lower case) that defines what your address is going to start with. Leave empty if you don't want to filter this;
- **end of the address**: a sequence of letters (upper/lower case) that defines what your address is going to end with. Leave empty if you don't want to filter this;
- **case sensitivity**: a boolean value (true/false) that defines whether the previous inputs are intended to be case sensitive or not.

This project will run by using half of the CPUs available on your hardware, while being light on RAM usage.

The process of generating a vanity address can be time and hardware intense if you try to match more than 4 characters, growing exponentially. Case sensitivity adds another layer of complexity to the generation.

## Future updates
- using parameters instead of asking for inputs
- manual setup of the number of CPUs used
