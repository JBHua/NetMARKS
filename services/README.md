# Readme for Mock Services

## Env Variable
1. Protocol used (gRPC/HTTP)
2. Constant latency (simulated), denoted in ms. Default to 0ms. Set using env variable

## Request & Response
- Request contains two fields:
    1. Quantity. Default to 1
    2. Response size. denoted in Byte. Default to 1. Max message size is 64MB, see: https://stackoverflow.com/questions/34128872/google-protobuf-maximum-size
Note: The base response size is around 215 bytes. So the total response size will be `base_size + response_size`
