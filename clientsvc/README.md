# Client Service

This engine is to process data that pushed to kafka by `File Processing Engine`

Here is how to get the result
Access to `http://localhost:9999/get-summary/{{stock_code}}` and it will return json result
`stock_code` stand for BBCA, BBRI etc
```
http://localhost:9999/get-summary/BBCA
```
