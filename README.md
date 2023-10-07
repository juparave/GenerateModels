# GenModels - Gorm Model Generator

GenModels is a command-line utility for generating Gorm models from a MySQL database schema. It simplifies the process of creating models for your Go applications based on your database structure.

## Installation

You can install GenModels using `go get`:

```bash
go get github.com/juparave/genmodels
```

## Usage

To generate Gorm models from your MySQL database, run the following command:

```bash
genmodels [flags]
```

### Flags

* -d, --database: Database name.
* -u, --user: Database user.
* -p, --password: Database password.
* -H, --host: Database host (default: localhost).
* -P, --port: Database port (default: 3306).

### Example
Generate Gorm models for a MySQL database named mydb with the user user and password password:

```bash
genmodels -d mydb -u user -p password
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
Cobra - A Commander for modern Go CLI interactions.
Gorm - The fantastic ORM library for Golang.
