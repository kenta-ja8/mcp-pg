# mcp–pg

### build
```
go build -o dist/ .
```

### show list
```
go build -o dist/ .
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./dist/mcp-pg | jq
echo '{"jsonrpc":"2.0","id":2,"method":"resources/list","params":{}}' | ./dist/mcp-pg | jq
echo '{"jsonrpc":"2.0","id":3,"method":"prompts/list","params":{}}' | ./dist/mcp-pg | jq
```

### call tools
```
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"add","x":5,"y":3}}}' | \
./dist/mcp-pg |\
jq
```
