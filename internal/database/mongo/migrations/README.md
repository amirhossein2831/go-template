to create migration file do this
first cd to migration dir
the run command 
```bash
     migrate create -ext json -dir . -seq migration_name
```
also in application.yml you can have three type of migration type

up
down
step

if you using step you should set the number of step you want go forward or backward positive for forward for example 4 and negative for back ward for example -2 to back rollback two last migration