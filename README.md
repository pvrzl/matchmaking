### Run API
- Make sure you have docker installed on your system
- Make sure port 5432 and 6379 is available
- setup .env file
- Run API command: ```docker-compose up```
- Run db/redis only: ```docker-compose up db redis```

### Migrations
- to create new migration : ```make migrate-create```
- to execute new migration : ```make migrate-up```
- to rollback to last migration : ```make migrate-down```
- to force migration to certain version: ```make migrate-fix```
  
### Run Test
- All file should have test file, example : file_name.go should have file_name_test.go
- Make sure test SUCCESS and FAILED case with coverage 80%
- Run test with coverage locally

### Helper Tools
- ```make generator-repository``` : generate simple repository
- ```make generator-interface``` : generate interface from selected file
- ```make generator-usecase``` : generate simple usecase

### Unit Test
Refer to https://github.com/golang/mock to create Mock Object. Set destination of Mock Object to ```mock``` folder in this project.
