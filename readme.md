# htmxxx

> htmx playground

## usage

```sh
npm ci
NODE_ENV=dev npm start
while inotifywait -q -e modify index.html > /dev/null; do curl http://localhost:3000/refresh -f || break; sleep 1; done
```

## author

balazs4
