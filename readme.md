# htmxxx

> htmx playground

## usage

```sh
npm ci
```

```
npm start
```

## dev

### auto-refresh

```sh
while inotifywait -q index.html > /dev/null; do curl http://localhost:3000/refresh -f || break; sleep 1; done
```

## author

balazs4
