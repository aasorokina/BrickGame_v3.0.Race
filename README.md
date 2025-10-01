# BrickGame_v3.0.Race

Implementation of the Race game in the Go programming language - game library, REST API server and a web interface

## Project Description

Games and interfaces from the BrickGame v1.0 and BrickGame v2.0 projects have been transferred to this project.

The implementation of the Race game in the Go programming language consists of several components: a library responsible for the game logic (which implemented using a finite-state machine), a web interface, server providing REST API for accessing games via the network and client library that encapsulates the interaction with the server for console and desktop interfaces.

Race game library also can be directly connected to the console interface from BrickGame v1.0 and the desktop interface from BrickGame v2.0, both interfaces fully support new game. Games from the BrickGame v1.0 and BrickGame v2.0 projects are accessible via the REST API server and web interface.

![web](images/brick_game_web.gif)  

## Features

Developed games:

- Tetris (C library);
- Snake (C++ library);
- Race (Golang library/package).

Developed interfaces:

- Console (C, ncurses library);
- Dektop (C++, Qt);
- Web (HTML, CSS, JavaScript).

REST API server - Golang, Gin framework;  
Client library - C++, libcurl, nlohmann/json;

## Race game features

Core mechanics:

- Player’s car can switch lanes with ← → keys.
- Collision with another car ends the game.
- Rival cars spawn at the top and move downward.
- Holding ↑ key speeds up the movement of rival cars.

Scoring system:

- Points are awarded for every avoided rival car.
- High score is stored in a file and persists across sessions.

Level system:

- Every 5 points increases the level by +1.
- Maximum level: 10.
- Higher levels increase rival car speed.

Controls:

- Start game - Enter;
- Pause - p;
- End game - q;
- Move into left - Left arrow;
- Moving into right - Right arrow;
- Speed up - Up arrow.

## REST API

Server runs at `localhost:8080`.

Main endpoints:

GET /api/games – list of available games.  
POST /api/games/:id – start a game.  
POST /api/actions – perform a player action.  
GET /api/state – get current game state.  

## Building project

Games libraries code located in the `src/brick_game/<game_name>` folders.  
Server code located in the `src/brick_game/server` folder.  
Client library code located in the `src/brick_game/client` folder.  
Interfaces code located in the `src/brick_game/gui` folder.  
Installation directory - `build`.

Main Makefile targets:

`install` - builds all games with console and desktop interfaces, server, console and desktop apps with client library and places in an installation directory;  
`uninstall` - removes the installation directory;  
`run` - launches the menu with all available game options;  
`test` - runs tests;  
`<game_name>.a` - builds static library of the selected game.

## Dependencies

- Golang >= 1.22.0;
- gin;
- Qt >= 5.15;
- ncurses (libncurses5-dev);
- lubcurl (libcurl4-openssl-dev).
