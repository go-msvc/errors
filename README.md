# Errors, Codes and Retries

Offers:
* easy-to-read nested errors,
* source references,
* numerical error codes, and
* retryable errors.

## Return an Error

Respond with an error message using Error() or Errorf():
```
import "github.com/go-msvc/errors"
...
return errors.Error("failed to do something")
return errors.Errorf("failed to connect to %s", addr)
```

Make a new retryable error:
```
return errors.Retryf(time.Second*5, "cannot connect to %s", addr)
```

Make a new error with a code:
```
return errors.Codef(404, "id %s not found", id)
```

Wrap an existing error with one of the following:
```
return errors.Wrap(err, "my error")
return errors.Wrapf(err, "my error %s", errorType)
return errors.Retry(err, time.Second*5)
return errors.Code(err, 404)
```

## Handle Errors

Use the following to format your error messages:

`%s` error message without source reference
`%v` error message with source reference as base filename
`%V` error message with source reference as fullpath filename

Add a `+` (like `%+s`) to recurse into wrapped errors with `" because "` linking the wrapped errors in one line.
Add a `-` (like `%-s`) to recurse into wrapped errors with `"\n"` linking the wrapped errors over multiple lines.

Examples:
* `%+s` is easy to read error like "cannot create user _because_ invalid request _because_ missing surname"
* `%+v` is good for logging errors, on one line with basefile reference of each error in the stack
* `%-V` is full error stack spanning multiple lines with full path names on source references.


# Examples
## Errors from Other Packages
Wrap errors from other packages to make it easy to see what failed and where in your code the failure occured.

In this example, the file does not exist:
```
f,err := os.Open(filename)
if err != nil {
    return errors.Wrap(err, "cannot open file")
}
```

## Doing Validations

See [Validations Example](./examples/validations/README.md)

## Retryable Errors

To see if an error is retryable, use `errors.IsRetryable(err)`, or check that and get the time when it can be retried with `when,ok := errors.RetryableAt(err)`.

## Error Codes

To see if an error has a code, use `errors.HasCode(err)`, or check and get the code with `code, ok := errors.GetCode(err)`.

## Named Errors

Today it is more common to use named errors instead of numerical codes. Code is mostly used with things like HTTP. For named errors use the standard `errors.New(<name>)` or `errors.Error(<name>)`. It is the go way of doing it. That can be wrapped many times and then check if that is the error using `errors.Is()`.
