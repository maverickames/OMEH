# omeh - One more error handler
Handle and return unique errors that can be searched quickly to find the exact error. 

Allows you to return unique errors to the UI that can be quickly found in a sea of logs with the exact location and details where
the error occured. 

<img width="390" alt="image" src="https://user-images.githubusercontent.com/2593364/121457646-dd2fbe80-c95c-11eb-8aeb-fc5af2a1f8b7.png">

```
[Example App] ErrorHandler: 2021/06/09 19:47:07 
  -- Function: main.(*App).brokenRoute
  -- SourceFile: /Users/matthewames/Development/go/omeh/example/routes.go
  -- LineNumber: 36
  -- ErrorDetails: did something and it broke
  -- RequestDetail: /broken
  -- ErrorCode: 5000732679545666175
```

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/maverickames/omeh) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/maverickames/omeh/blob/master/LICENSE)
