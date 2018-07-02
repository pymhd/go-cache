## go-cache
#### Key value  operations
 - **Add Item:**
`curl -X POST -H "Content-Type: application/json" -d '{"value": "VALUE", "ttl": 3600}''http://127.0.0.1:9000/map/keyname'`
 - **Del Item:**
`curl -X DELETE 'http://127.0.0.1:9000/map/KEYNAME'`
  - **Get Item:**
`curl  -X GET 'http://127.0.0.1:9000/map/KEYNAME'`


#### Key List of Values operations
 - **Add item to list**
 `curl -X PUT 'http://127.0.0.1:9000/list/KEYNAME' -H "Content-Type: application/json" -d '{"value": "VALUE"}'`
 - **Check if item is in list**
 `curl -X GET 'http://127.0.0.1:9000/list/KEYNAME' -H "Content-Type: application/json" -d '{"value": "VALUE"}'`
  Response: {"value":false} or {"value":true}
 - ** Delete from list**
 `curl -X DELETE 'http://127.0.0.1:9000/list/KEYNAME' -H "Content-Type: application/json" -d '{"value": "VALUE"}'`

### Health
`curl -X GET http://127.0.0.1:9000/healthz`
