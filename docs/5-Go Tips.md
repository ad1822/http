## Some Error Handing Thing

### Rule of thumb

- `log.Fatalf` → startup failures, unrecoverable errors.

- `log.Printf` → runtime errors, warnings, diagnostics.

- `fmt.Printf` → normal program output (what the program intends to show).
