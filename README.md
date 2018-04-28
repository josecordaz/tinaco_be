
```bash
    docker build -t "tinaco_be" .
```

```bash
    docker run -d -p 8000:8000 --name tinaco_be tinaco_be
```

```bash
    docker rm -f tinaco_be
```

```bash
    docker tag ID josecordaz/tinaco_be:1.0
```

```bash
    docker push josecordaz/tinaco_be:1.0
```

```bash
    docker run -d -p 8000:8000 --name tinaco_be josecordaz/tinaco_be:1.0
```