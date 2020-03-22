# covid-socialgraph

## Development

### Docker Neo4j
```
docker run --rm \
    --publish=7474:7474 --publish=7687:7687 \
    --volume=$HOME/neo4j/data:/data \
    neo4j
```

### Docker Socialgraph
```
docker build -t socialgraph .
docker run --rm \
    --env COVIDHUB_DB_USER=neo4j \
    --env COVIDHUB_DB_PASSWORD=neo4j \
    socialgraph
```
