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
team data is used. Messages are processed in queue with QoS 2, which means that every message is guaranteed to be delivered to 
subscriber exactly once. In addition to that Simulator and Statistic services are independent, so can be run separately and in any order. 

Simulation is started immediatelly once service is started and service is stopped after simulation is finished. Statistic service 
never stops waiting for new events to arrive. 

There are 2 types of events sent by Simulator: 

  1. Game status update - indicates that game status was updated, e.g. game started or finished. Contains information about teams: name, 
  players etc, home/guest tean etc. 
  2. Scoring event - indicates that one of teams scored 2 or 3 points in the game. 
  
Next statistics is calculated (and it can be expanded further in an obvious way): 

1. Game score - how many points each team (home and guest) scored so far in the game. 
2. Game start/finish time
3. Last event timestamp - when last event in the game happened (team scored or status changed). Indicates time elapsed since the 
beginning of the game. 
4. Overall score for home/guest teams within all games. 