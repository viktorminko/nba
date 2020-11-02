### Problem

Build simulation of NBA games. All games are started and finished at the same time. During the game every team 
can randomly score 0, 2 or 3 points within every N seconds period. Once team scored the event is fired and statistics is updated 
accordingly. 


### Architecture

There are 3 separate services: 
  1. Simulator: this service generates games schedule and runs simulation. During the simulation the services sends events to the 
  queue. MQTT is used as a queue provider. 
  2. MQTT broker - accepts and processes events messages. 
  3. Statistics - subscribes to broker, handles games events, updates statistics and runs http server to view statistics.
  
There is no DB involved, teams data (names, players) is read to Simulator memory on service start, pre defined json file with 
team data is used. Messages are processed in a queue with QoS 2, which means that every message is guaranteed to be delivered to 
subscriber exactly once. In addition to that Simulator and Statistic services are independent, so can be run separately and in any order. 

Simulation is started immediately once service is started and service is stopped after simulation is finished. Statistic service 
never stops waiting for new events to arrive. 

There are 2 types of events sent by Simulator: 

  1. Game status update - indicates that game status was updated, e.g. game started or finished. Contains information about teams: name, 
  players etc, home/guest team etc. 
  2. Scoring event - indicates that one of teams scored 2 or 3 points in the game. Only contains minimum amount of data. 
  
 Messages are encoded/decoded in binary using golang encoding/gob package.
  
Next statistic is being calculated (and it can be expanded further in obvious way): 

1. Game score - how many points each team (home and guest) scored so far in the game. 
2. Game start/finish time
3. Last event timestamp - when last event in the game happened (team scored or status changed). Indicates time elapsed since the 
beginning of the game. 
4. Overall score for home/guest teams within all games. 

Frontend is represented by simple html template https://github.com/viktorminko/nba/blob/master/pkg/statistic/frontend/html/layout.html 
Page is updated automatically. 

### Usage

Build and run containers using Make 

`
make up
`

This will build executables from source and start 3 services mentioned above. Simulation will start once container is up and running. 
To see statistics navigate to http://0.0.0.0:8080/ in your browser. Page will be updated automatically. 
Simulation container stops after simulation is done. Statistics container never stops and waiting for new simulation events. 
Statistics will be updated when simulation is started again. 


Run existing containers with docker without building from source

`
docker-compose up 
`
This will run containers using existing images. Use this to save time and avoid building from source.


You can use Make to build/run services separately, run unit tests, linter or update dependencies. 
See https://github.com/viktorminko/nba/blob/master/Makefile

You can change app configuration using environment variables in docker-compose. 
See available configuration options here https://github.com/viktorminko/nba/blob/master/pkg/simulation/opts/env.go and here https://github.com/viktorminko/nba/blob/master/pkg/statistic/opts/env.go