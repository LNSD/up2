<h1>UP²</h1> 

> **U**pload **P**re-signed **U**RL **P**rovider → UPUP → UP²

Microservice written in GO that generates pre-signed file upload (and download) URLs through a Restful API.

## Tech stack:

* Go v1.17
* [Labstack Echo](https://echo.labstack.com/)
* OpenAPI Codegen (for boilerplate generation)
* [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/)
* [MinIO Go client](https://docs.min.io/docs/golang-client-api-reference.html)
* Testify
* [Testcontainers-Go](https://golang.testcontainers.org/) (for integration tests)

## License

```
MIT License

Copyright (c) 2022 Lorenzo Delgado Antequera

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```