# httpin-example

This is an example project for demonstrating how to use [ggicci/httpin](https://github.com/ggicci/httpin) package to parse HTTP requests in Go.

- [httpin Documentation](https://ggicci.github.io/httpin/)

## Run this demo locally

```bash
# Method 1
go install github.com/ggicci/httpin-example # install `httpin-example` binary locally
httpin-example

# Method 2
git clone https://github.com/ggicci/httpin-example.git
cd httpin-example
go build
./httpin-example
```

## Test the demo using `curl` command

```bash
curl -XPOST -s -H'Content-Type: application/json' --data '{"login":"ggicci","is_member":false,"age":18}' "http://localhost:8080/users" | python3 -m json.tool
# {
#     "input": {
#         "login": "ggicci",
#         "created_at": "2023-09-30T12:33:30.110696-04:00",
#         "is_member": false,
#         "age": 18
#     },
#     "users": [
#         {
#             "login": "ggicci",
#             "created_at": "2023-09-30T12:33:30.110696-04:00",
#             "is_member": false,
#             "age": 18
#         }
#     ]
# }
curl -s "http://localhost:8080/users?is_member=1"
curl -s "http://localhost:8080/users?is_member=1&sort_by[]=name&sort_desc[]=true"
curl -s "http://localhost:8080/users?is_member=1&sort_by[]=name&sort_by[]=age&sort_desc[]=true&sort_desc[]=false"
```
