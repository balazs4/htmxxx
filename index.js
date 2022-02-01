let clients = [];

require('http')
  .createServer((req, res) => {
    if (req.url === '/hello') {
      res.end(JSON.stringify({ date: new Date().toJSON() }));
      return;
    }

    if (req.url === '/event') {
      res.writeHead(200, {
        'content-type': 'text/event-stream',
        connection: 'keep-alive',
        'cache-control': 'no-cache',
      });
      const id = Date.now();
      clients.push({ id, res });

      req.on('close', () => {
        console.log(`${id} connection closed`);
        clients = clients.filter((client) => client.id !== id);
      });
      return; // ????
    }

    if (req.url === '/refresh') {
      clients.forEach((client) => {
        client.res.write(`data: refresh ${Date.now()}\n\n`);
      });
      res.end();
      return;
    }

    const file = require('path').join('.', req.url);

    require('fs').createReadStream(file).pipe(res);
  })
  .listen(process.env.PORT, () =>
    console.log(`http://localhost:${process.env.PORT}`)
  );
