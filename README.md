# Contacts Book

A command-line contact management application written in Go.

## Features

- Add new contacts (name, email, mobile)
- List all contacts
- Search contacts by name (partial matching)
- Delete contacts
- Edit existing contacts

## Input Validation

- Email: Must be a valid email format
- Mobile: Must be 10 digits starting with 0 (automatically converts to +971 format)
- Names: Automatically formatted to title case

## Usage

Run the application:

```
go run contacts.go
```

Or build and run:

```
go build contacts.go
./contacts
```

## Menu Options

```
(1). Add Contact
(2). List Contacts
(3). Search
(4). Delete Contact
(5). Edit Contact
(6). Exit
```

## Data Storage

Contacts are stored in `contacts.txt` in CSV format:

```
Name,Email,Mobile
```

## Dependencies

- golang.org/x/text/cases
- golang.org/x/text/language

Install dependencies:

```
go mod init contacts-book
go mod tidy
```

## Requirements

- Go 1.18 or later
