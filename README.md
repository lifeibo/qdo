QDo
======================
A LevelDB queue implementation written in Golang.


Create conveyor
----------------------
    curl http://127.0.0.1:8080/api/conveyor \
       -d conveyor_id=foo \
       -d max_concurrent=2 \
       -d max_rate=100 \
       -d task_timeout=60 \
       -d task_max_tries=10

Update conveyor
----------------------
    curl http://127.0.0.1:8080/api/conveyor/foo \
       -d paused=true

Delete conveyor
----------------------
    curl -X DELETE http://127.0.0.1:8080/api/conveyor/foo

Create task
----------------------
    curl http://127.0.0.1:8080/api/conveyor/foo/task \
       -d target=http://127.0.0.1/mytask \
       -d "payload={'foo': 'bar'}"

Create scheduled task
----------------------
    curl http://127.0.0.1:8080/api/conveyor/foo/task \
       -d target=http://127.0.0.1/mytask \
       -d scheduled=1399999999 \
       -d "payload={'foo': 'bar'}"

Delete all tasks
----------------------
    curl -X DELETE http://127.0.0.1:8080/api/conveyor/foo/task


Get conveyor stats
----------------------
    curl http://127.0.0.1:8080/api/conveyor/foo/stats
