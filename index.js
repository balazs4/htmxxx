require("http")
  .createServer((req, res) => {
    if (req.url === "/hello") {
      res.end(JSON.stringify({ date: new Date().toJSON() }));
      return;
    }

    const file = require("path").join(".", req.url);

    require("fs")
      .createReadStream(file)
      .pipe(res);
  })
  .listen(process.env.PORT, () =>
    console.log(`http://localhost:${process.env.PORT}`)
  );
