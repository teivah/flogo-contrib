# tibco-etcd
This activity provides your flogo application the ability to connect to an etcd cluster

## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/etcd
```

## Schema
Inputs and Outputs:

```json
{
   "inputs":[
      {
         "name":"key",
         "type":"string",
         "required":true
      },
      {
         "name":"value",
         "type":"string"
      },
      {
         "name":"method",
         "type":"string",
         "allowed":[
            "Create",
            "Get",
            "Update",
            "Delete"
         ],
         "value":"Create",
         "required":true
      },
      {
         "name":"servers",
         "type":"string",
         "required":true
      }
   ],
   "outputs":[
      {
         "name":"output",
         "type":"any"
      }
   ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| key | The key identifier |
| value   | The value to be stored as a String |
| method       | The method type (Create, Get, Update or Delete) |
| servers  | The etcd servers, delimited by a *;* character (e.g. *http://etcd1:2379;http://etcd2:2379*) |
Note: if method is set to Get, value is ignored
## Configuration Examples
Configure an upsert method:

```json
{
   "key":"foo",
   "data":"bar",
   "method":"Create",
   "server":"http://etcd:2379"
}
```