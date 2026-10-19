

## Objetivo: 
criar um sistema de migrações entre versões do banco de dados


## Novo banco: 
criar um banco especifico na mecanica de database chamado versions , que ira armazenar a versão de cada banco de dados, comecando na 1.

## comands-inc: 
increments the version of the database, then create a migration function with the new version.
```bash 
agnos inc-database <db-nam> 
```

## update:
runs the update function for the database, that updates from the old version to the newest one.
## comands-update:
```bash 
agnos update-database <db-nam> 
```

## update-all:

updates all the databases to the newest version.
## comands-update-all:
```bash 
agnos update-database-all 
```


## on start 
it must have a middleware that checks if the version of the database is updated, if not, it terminate the program.