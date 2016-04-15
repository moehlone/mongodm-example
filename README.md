# Example API for mongodm

This is an example REST API for the [mongodm package](https://github.com/zebresel-com/mongodm/) written in golang with use of the [beego web framework](https://github.com/astaxie/beego). To simplify response/error handling the [beego-response package](https://github.com/zebresel-com/beego-response) is also included as an example.

There are two controllers (User and Message) where you can create and get something.


### You will find an implementation for

- Connection setup
- Models
- Model registration
- Model creation (persisting/saving)
- Population
- Relations (1:1, 1:N)
- Querys (Find, FindOne, ..)
- Custom validation
- Default validation
- Localisation example
- Response-/Error-Handling
- Pagination
- Using namespaces, routers and filters with beego

### You won`t find an implementation for

- Model updating
- Model deleting
- Authentication

Maybe I will add the missing points later..

## Requirements

- golang setup
- a mongo database (by now password authentication is not supported)
- some REST client for testing (optional)
- mongodm package `go get github.com/zebresel-com/mongodm`
- beego-response package `go get github.com/zebresel-com/beego-response`


## Quickstart

Move to your project directory and run `go build; ./mongodm-sample` (Unix) and beego should start the webserver on port `8080` for you. You can change beego settings in `conf/app.conf`, if necessary. To find all implemented and available REST endpoints open up `apidoc/index.html` in your browser. Thats all, you can test now!

**Feel free to contribute!**

Ask if you think that there is something missing.
