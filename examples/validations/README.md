# Using Errors in Validations

Implement `Validator` interface for each of your structs, as was done in [users.go](./users/users.go#14)

Then in a handler to add a user, the validate will check the request and construct error messages that are easy to understand.

Run the example:
```
cd examples/validation
go run main.go -p 8090
```

Then post an invalid request:
```
curl -D /dev/stderr -XPOST 'http://localhost:8090/add' -d '{}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:03:02 GMT
Content-Length: 37

invalid request because missing name
```

And see in the server log on stderr the following:
```
HTTP POST /: main.go(39):invalid request because users.go(16):missing name
```

The error is eash to interpret by developes, with references to the code in `main.go(39)` and `users.go(16)`.

The error given to the user is also very clear:
```
invalid request because missing name
```
If a name is added, the next error indicate what else is required:
```
% curl -D /dev/stderr -XPOST 'http://localhost:8090/add' -d '{"name":1}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:28:41 GMT
Content-Length: 136

cannot parse JSON body into users.AddUserRequest: json: cannot unmarshal number into Go struct field AddUserRequest.name of type string

% curl -D /dev/stderr -XPOST 'http://localhost:8090/add' -d '{"name":"Jan"}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:05:48 GMT
Content-Length: 46

invalid request because missing date-of-birth

% curl -D /dev/stderr -XPOST 'http://localhost:8090/add' -d '{"name":"Jan", "date-of-birth":"18 Nov 1973"}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:06:09 GMT
Content-Length: 80

invalid request because date-of-birth:"18 Nov 1973" not formatted as CCYY-MM-DD

% curl -D /dev/stderr -XPOST 'http://localhost:8090/add' -d '{"name":"Jan", "date-of-birth":"1973-11-18"}'
HTTP/1.1 200 OK
Date: Sat, 01 Mar 2025 15:06:16 GMT
Content-Length: 0

```

Note that the error messages follow the go convention to start with a lowercase letter, which makes it natural text when concatenated into a nested error.

Also note the error messages refer to the json tags of the fields, and not the field names, e.g. `missing date-of-birth` rather than `missing Dob`, because the user will need to know what to add to the JSON document in the request. As for the developer, they can follow the code reference to see what the field is called in the struct.

The error messages in `Validate()` also does not say `invalid request` or `invalid user`. Instead they only refer to the field that was found not to comply. The caller (http handler in this case), adds the next context to say `invalid request`, and when that is joined, the error reads well.

If you have nested structs, e.g. an Address inside the user request, then add a `Validate()` method to that type to check its own fields, and then call it inside the parent struct's `Validate()` method as was done for `UpdateUserRequest`.

```
% curl -D /dev/stderr -XPOST 'http://localhost:8090/upd'
HTTP/1.1 405 Method Not Allowed
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:17:55 GMT
Content-Length: 18

this is not a put

% curl -D /dev/stderr -XPUT 'http://localhost:8090/upd'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:18:01 GMT
Content-Length: 21

cannot parse JSON body into users.UpdateUserRequest

% curl -D /dev/stderr -XPUT 'http://localhost:8090/upd' -d '{}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:18:05 GMT
Content-Length: 63

invalid request because missing both date-of-birth and address

% curl -D /dev/stderr -XPUT 'http://localhost:8090/upd' -d '{"address":44}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:18:14 GMT
Content-Length: 21

cannot parse JSON body into users.UpdateUserRequest: json: cannot unmarshal number into Go struct field Address.address.Street of type string

% curl -D /dev/stderr -XPUT 'http://localhost:8090/upd' -d '{"address":{}}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:18:25 GMT
Content-Length: 63

invalid request because invalid address because missing street

% curl -D /dev/stderr -XPUT 'http://localhost:8090/upd' -d '{"address":{"street":"44 Wide Street"}}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 01 Mar 2025 15:18:40 GMT
Content-Length: 63

% invalid request because invalid address because missing street

% curl -D /dev/stderr -XPUT 'http://localhost:8090/upd' -d '{"address":{"street":"44 Wide Street","country":"South Africa"}}'
HTTP/1.1 200 OK
Date: Sat, 01 Mar 2025 15:18:52 GMT
Content-Length: 0

```
