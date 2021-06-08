## Database setup

Create the user and the database, then apply the schema file:
```
kubectl exec -it postgres-0 -- createuser -P todos
kubectl exec -it postgres-0 -- createdb todos -O todos
cat db/schema.sql | kubectl exec -i postgres-0 -- psql -Utodos -dtodos
```
