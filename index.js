let clients = [];

require('http')
  .createServer((req, res) => {
    if (process.env.NODE_ENV === 'dev') {
      if (req.url === '/event') {
        res.writeHead(200, {
          'content-type': 'text/event-stream',
          connection: 'keep-alive',
          'cache-control': 'no-cache',
        });
        const id = Date.now();
        clients.push({ id, res });

        req.on('close', () => {
          clients = clients.filter((client) => client.id !== id);
        });
        return;
      }

      if (req.url === '/refresh') {
        clients.forEach((client) => {
          client.res.write(`data: refresh ${Date.now()}\n\n`);
        });
        res.end(`${clients.length} client(s) refreshed\n`);
        return;
      }
    }

    if (req.url === '/hello') {
      res.end(JSON.stringify({ date: new Date().toJSON() }));
      return;
    }

    const filename = req.url === '/' ? '/index.html' : req.url;
    const file = require('path').join('.', filename);

    require('stream').pipeline(
      require('fs').createReadStream(file),
      async function* (source) {
        for await (const chunk of source) {
          yield chunk;
        }
        if (filename === '/index.html' && process.env.NODE_ENV === 'dev') {
          yield `<script>new EventSource('/event').onmessage = () => location.reload();</script>`;
        }
      },
      res,
      (err) => {
        if (err) {
          console.log(
            `[ERR] ${req.url} cannot be stream to the client. ${err}`
          );
          return;
        }
        console.log(`[OK] ${req.url}`);
      }
    );
  })
  .listen(process.env.PORT, () =>
    console.log(`http://localhost:${process.env.PORT}`)
  );
