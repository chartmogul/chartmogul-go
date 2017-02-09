// For testing requests with node.js
// node dump-server.js

var http = require('http');

var server = http.createServer().listen(8000);

server.on('request', function(request, response) {
    console.log(request.method);
    console.log(request.url);
    console.log(request.headers);
    var body = [];
    request.on('data', function(chunk) {
      body.push(chunk);
    }).on('end', function() {
      body = Buffer.concat(body).toString();
      // at this point, `body` has the entire request body stored in it as a string
      console.log("Body:");
      console.log(body);
      response.end();
    });
});
