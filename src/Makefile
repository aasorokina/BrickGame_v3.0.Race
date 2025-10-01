SHELL = /bin/bash

CC = gcc
CFLAGS = -std=c11 -Wall -Werror -Wextra -g
CXXFLAGS = -Wall -Werror -Wextra -g -std=c++17
LFLAGS = -lcheck -lsubunit -lrt -lpthread -lm
GFLAGS = -fprofile-arcs -ftest-coverage

GUIDIR = brick_game/gui
BUILDDIR = build
LIBDIR = $(BUILDDIR)/libs
SERVERDIR = brick_game/server

TETRIS_EXE = tetris
SNAKE_EXE = snake
RACE_EXE = race
BRICKGAME_EXE = brickgame

TETRIS_TEST = tetris_test
SNAKE_TEST = snake_test
RACE_TEST = race_test

LIB_TETRIS = tetris.a
LIB_SNAKE = snake.a
LIB_RACE = race.a
LIB_CLIENT = client.a

LIB_TETRIS_SRC = $(wildcard brick_game/tetris/*.c)
LIB_SNAKE_SRC = $(wildcard brick_game/snake/*.cpp)
LIB_CLIENT_SRC = $(wildcard brick_game/client/*.cpp)

TEST_TETRIS_SRC = $(wildcard test/tetris/*.c)
TEST_SNAKE_SRC = $(wildcard test/snake/*.cpp)

LIB_TETRIS_O = $(LIB_TETRIS_SRC:.c=.o)
LIB_SNAKE_O = $(LIB_SNAKE_SRC:.cpp=.o)
LIB_CLIENT_O = $(LIB_CLIENT_SRC:.cpp=.o)

TEST_TETRIS_O = $(TEST_TETRIS_SRC:.c=.o)
TEST_SNAKE_O = $(TEST_SNAKE_SRC:.cpp=.o)

GCOV_NAME = gcov_tests.info
DIST_FOLDER = BrickGame-3.0
DIST_NAME = BrickGame-3.0

GAME = 0

all: clean install
.PHONY: all clean tetris.a snake.a race.a install uninstall dvi dist test gcov_report style_check style_fix valgrind_test

install:
	@echo ""
	@echo "Installing ..."
	@make -s tetris_console
	@make -s snake_console
	@make -s race_console
	@make -s client_console

	@make -s tetris_desktop
	@make -s snake_desktop
	@make -s race_desktop
	@make -s client_desktop

	@make -s launcher
	@make -s server
	@echo ""
	@echo "Installed. Use <make run> for running apps"

uninstall:
	@rm -rf $(BUILDDIR)

clean:
	@rm -rf *.o *.a *.gcno *.gcda *.info report *.txt $(TETRIS_TEST) $(SNAKE_TEST) $(RACE_TEST) latex html *.tgz 

########################################## BUILD LIBS ##########################################

# tetris.a
$(LIB_TETRIS): $(LIB_TETRIS_O)
	@mkdir -p $(LIBDIR)
	@ar rc $(LIB_TETRIS) $(LIB_TETRIS_O)
	@ranlib $(LIB_TETRIS)
	@rm -rf brick_game/tetris/*.o
	@mv $(LIB_TETRIS) $(LIBDIR)/$(LIB_TETRIS)
	@$(call wrap_lib,tetris)

# snake.a
$(LIB_SNAKE): $(LIB_SNAKE_O)
	@mkdir -p $(LIBDIR)
	@ar rc $(LIB_SNAKE) $(LIB_SNAKE_O)
	@ranlib $(LIB_SNAKE)
	@rm -rf brick_game/snake/*.o
	@mv $(LIB_SNAKE) $(LIBDIR)/$(LIB_SNAKE)
	@$(call wrap_lib,snake)

# race.a
$(LIB_RACE):
	@mkdir -p $(LIBDIR)
	@cd brick_game/race && go build -buildmode=c-archive -o $(LIB_RACE) race.go
	@mv brick_game/race/race.h $(LIBDIR)/race.h
	@mv brick_game/race/$(LIB_RACE) $(LIBDIR)/$(LIB_RACE)

# client.a
$(LIB_CLIENT): $(LIB_CLIENT_O)
	@mkdir -p $(LIBDIR)
	@ar rc $(LIB_CLIENT) $(LIB_CLIENT_O)
	@ranlib $(LIB_CLIENT)
	@rm -rf brick_game/client/*.o
	@mv $(LIB_CLIENT) $(LIBDIR)/$(LIB_CLIENT)

########################################## BUILD CONSOLE APPS ##########################################

# running without server
tetris_console: $(LIB_TETRIS)
	@mkdir -p $(BUILDDIR)/console
	@$(CC) $(CFLAGS) -DTETRIS $(GUIDIR)/cli/*.c -o $(TETRIS_EXE) -L./$(LIBDIR) -l:$(LIB_TETRIS) -lncurses
	@mv $(TETRIS_EXE) $(BUILDDIR)/console/$(TETRIS_EXE)
	@rm -rf *.o

snake_console: $(LIB_SNAKE)
	@mkdir -p $(BUILDDIR)/console
	@g++ $(CXXFLAGS) -DSNAKE $(GUIDIR)/cli/*.c -o $(SNAKE_EXE) -L./$(LIBDIR) -l:$(LIB_SNAKE) -lncurses
	@mv $(SNAKE_EXE) $(BUILDDIR)/console/$(SNAKE_EXE)
	@rm -rf *.o

race_console: $(LIB_RACE)
	@mkdir -p $(BUILDDIR)/console
	@$(CC) $(CFLAGS) -DRACE $(GUIDIR)/cli/*.c -o $(RACE_EXE) -L./$(LIBDIR) -l:$(LIB_RACE) -lncurses
	@mv $(RACE_EXE) $(BUILDDIR)/console/$(RACE_EXE)
	@rm -rf *.o

# running with server
client_console: $(LIB_CLIENT)
	@mkdir -p $(BUILDDIR)/console
	@g++ $(CXXFLAGS) -DREST_CLIENT $(GUIDIR)/cli/*.c -o $(BRICKGAME_EXE) -L./$(LIBDIR) -l:client.a -lncurses -lcurl
	@mv $(BRICKGAME_EXE) $(BUILDDIR)/console/$(BRICKGAME_EXE)

########################################## BUILD DESKTOP APPS ##########################################

# running without server
tetris_desktop: $(LIB_TETRIS)
	@mkdir -p $(BUILDDIR)/desktop
	@cd $(BUILDDIR)/desktop && qmake ../../$(GUIDIR)/desktop/desktop.pro DEFINES+=TETRIS && make clean && make

snake_desktop: $(LIB_SNAKE)
	@mkdir -p $(BUILDDIR)/desktop
	@cd $(BUILDDIR)/desktop && qmake ../../$(GUIDIR)/desktop/desktop.pro DEFINES+=SNAKE && make clean && make

race_desktop: $(LIB_RACE)
	@mkdir -p $(BUILDDIR)/desktop
	@cd $(BUILDDIR)/desktop && qmake ../../$(GUIDIR)/desktop/desktop.pro DEFINES+=RACE && make clean && make

# running with server
client_desktop: $(LIB_CLIENT)
	@mkdir -p $(BUILDDIR)/desktop
	@cd $(BUILDDIR)/desktop && qmake ../../$(GUIDIR)/desktop/desktop.pro DEFINES+=REST_CLIENT && make clean && make

####################################### BUILD LAUNCHER FOR APPS WITH SERVER #######################################

launcher: libcurl $(LIB_CLIENT) $(LIB_TETRIS) $(LIB_SNAKE) $(LIB_RACE) 
	@g++ $(CXXFLAGS) -DREST_CLIENT $(GUIDIR)/launcher/*.cpp -o launcher -L./$(LIBDIR) -l:$(LIB_CLIENT) -lncurses -lcurl
	@mv launcher $(BUILDDIR)/launcher

##################################################### SERVER ######################################################

server: $(LIB_TETRIS) $(LIB_SNAKE) $(LIB_RACE)
	@cd $(SERVERDIR) && go build -o server ./cmd/main.go && mv server ../../$(BUILDDIR)/server

server_run:
	@if [ -f server.pid ]; then \
		echo "Server already running!"; \
	else \
		./$(BUILDDIR)/server >/dev/null 2>&1 & echo $$! > server.pid; \
		echo "Server started!"; \
	fi
			
server_stop:
	@kill $$(cat server.pid) 2>/dev/null || true
	@echo "Server stopped!"
	@rm -f server.pid

web_gui:
	@echo "http://localhost:8080/"

###################################################### TESTS ######################################################

test: $(TEST_TETRIS_O) $(LIB_TETRIS) $(TEST_SNAKE_O) $(LIB_SNAKE)
	@mkdir -p $(BUILDDIR)/test
	@$(CC) $(CFLAGS) $(TEST_TETRIS_O) -o $(TETRIS_TEST) -L./$(LIBDIR) -l:$(LIB_TETRIS) $(LFLAGS)
	@mv $(TETRIS_TEST) $(BUILDDIR)/test/$(TETRIS_TEST)
	@./$(BUILDDIR)/test/$(TETRIS_TEST)

	@g++ $(CXXFLAGS) -DTEST $(TEST_SNAKE_O) -o $(SNAKE_TEST) -L./$(LIBDIR) -l:$(LIB_SNAKE) -lgtest -pthread -lstdc++ -lsupc++ -lgcc
	@mv $(SNAKE_TEST) $(BUILDDIR)/test/$(SNAKE_TEST)
	@./$(BUILDDIR)/test/$(SNAKE_TEST)

	@cd brick_game/race && go test -c -o $(RACE_TEST) -cover && mv $(RACE_TEST) ../../$(BUILDDIR)/test/$(RACE_TEST)
	@./$(BUILDDIR)/test/$(RACE_TEST) -test.v
	@rm -rf *.o test/snake/test_snake.o test/tetris/test_tetris.o

valgrind_test: test
	@valgrind --tool=memcheck --leak-check=yes --show-leak-kinds=all --track-origins=yes -s ./$(BUILDDIR)/test/$(TETRIS_TEST)
	@valgrind --tool=memcheck --leak-check=yes --show-leak-kinds=all --track-origins=yes -s ./$(BUILDDIR)/test/$(SNAKE_TEST)

###################################################### COVERAGE ######################################################

gcov_report_tetris: clean
	@$(CC) $(CFLAGS) -c $(LIB_TETRIS_SRC) --coverage
	@$(CC) $(CFLAGS) -c $(TEST_TETRIS_SRC)
	@$(CC) *.o -o $(TETRIS_TEST) $(LFLAGS) $(GFLAGS)
	@./$(TETRIS_TEST)
	lcov -t "gcov_tests" -o $(GCOV_NAME) -c -d .
	genhtml -o report $(GCOV_NAME)
	@rm -rf *.gcno *.gcda *.gcov $(GCOV_NAME) *.o

gcov_report_snake: clean
	@g++ $(CXXFLAGS) -c $(LIB_SNAKE_SRC) --coverage
	@g++ $(CXXFLAGS) -c $(TEST_SNAKE_SRC)
	@g++ *.o -o $(SNAKE_TEST) $(LFLAGS) $(GFLAGS) -lgtest -pthread -lstdc++ -lsupc++ -lgcc
	./$(SNAKE_TEST)
	@lcov -t "gcov_tests" -o $(GCOV_NAME) -c -d . \
		--exclude '/usr/include/c++/11/*' \
		--exclude '/usr/include/c++/11/bits/*' \
		--exclude '/usr/include/c++/11/ext/*' 
	genhtml -o report $(GCOV_NAME)
	@rm -rf *.gcno *.gcda *.gcov $(GCOV_NAME) *.o

###################################################### STYLE ######################################################

style_check:
	@cp ../materials/linters/.clang-format .clang-format
	@find . -type f \( -name "*.c" -o -name "*.h" -o -name "*.cpp" \) -not -path "./build/**/*" -exec clang-format -n {} +
	@rm -f .clang-format

style_fix:
	@cp ../materials/linters/.clang-format .clang-format
	@find . -type f \( -name "*.c" -o -name "*.h" -o -name "*.cpp" \) -not -path "./build/**/*" -exec clang-format -i {} +
	@rm -f .clang-format

fmt:
	@bash -c "cd brick_game/race && go fmt ./..."
	@bash -c "cd brick_game/server && go fmt ./..."

lint:
	@bash -c "cd brick_game/race && golint ./..."
	@bash -c "cd brick_game/server && golint ./..."

##################################################################################################################

dvi:
	@doxygen Doxyfile
	@cd latex && xelatex refman.tex && cd ../

xelatex_install:
	sudo apt-get install texlive-xetex
	sudo apt-get install texlive-lang-cyrillic

dist: #uninstall install
	@mkdir -p $(DIST_FOLDER)/console
	@mkdir -p $(DIST_FOLDER)/desktop

	@cp -r $(BUILDDIR)/console $(DIST_FOLDER)/
	@cp $(BUILDDIR)/desktop/tetris_desktop $(DIST_FOLDER)/desktop/tetris_desktop
	@cp $(BUILDDIR)/desktop/snake_desktop $(DIST_FOLDER)/desktop/snake_desktop
	@cp $(BUILDDIR)/desktop/race_desktop $(DIST_FOLDER)/desktop/race_desktop
	@cp $(BUILDDIR)/desktop/brickgame $(DIST_FOLDER)/desktop/brickgame

	@cp $(BUILDDIR)/launcher $(DIST_FOLDER)/launcher
	@cp $(BUILDDIR)/server $(DIST_FOLDER)/server

	tar cvzf $(DIST_NAME).tgz $(DIST_FOLDER)/
	@rm -rf $(DIST_FOLDER)/

##################################################################################################################

dependencies:
	sudo apt-get update
	@make -s libcurl
	@make -s libncurses

libcurl:
	sudo apt-get install -y libcurl4-openssl-dev

libncurses:
	sudo apt-get install -y libncurses5-dev

run: menu
	@$(eval GAME := $(shell read -p "" input && echo $$input))
	@if [ "$(GAME)" = 1 ]; then \
		./$(BUILDDIR)/console/tetris; \
		exit 0; \
	elif [ "$(GAME)" = 2 ]; then \
		./$(BUILDDIR)/console/snake; \
		exit 0; \
	elif [ "$(GAME)" = 3 ]; then \
		./$(BUILDDIR)/console/race; \
		exit 0; \
	elif [ "$(GAME)" = 4 ]; then \
		./$(BUILDDIR)/desktop/tetris_desktop; \
		exit 0; \
	elif [ "$(GAME)" = 5 ]; then \
		./$(BUILDDIR)/desktop/snake_desktop; \
		exit 0; \
	elif [ "$(GAME)" = 6 ]; then \
		./$(BUILDDIR)/desktop/race_desktop; \
		exit 0; \
	elif [ "$(GAME)" = 7 ]; then \
		echo ""; \
		make -s server_run; \
		echo ""; \
		echo "For stop server use <make server_stop>"; \
		echo ""; \
		./$(BUILDDIR)/launcher; \
		exit 0; \
	elif [ "$(GAME)" = 8 ]; then \
		echo ""; \
		make -s server_run; \
		echo ""; \
		echo "For stop server use <make server_stop>"; \
		echo ""; \
		make -s web_gui; \
		exit 0; \
	else \
		echo "Incorrect number"; \
	fi

menu:
	@echo "<<BrickGame>>"
	@echo ""
	@echo "1. Tetris (console)"
	@echo "2. Snake (console)"
	@echo "3. Race (console)"
	@echo ""
	@echo "4. Tetris (desktop)"
	@echo "5. Snake (desktop)"
	@echo "6. Race (desktop)"
	@echo ""
	@echo "7. Console & desktop games on server"
	@echo "8. Run server with web gui"
	@echo ""
	@echo -n "Enter a number: "

define wrap_lib
cp $(LIBDIR)/$1.a $(LIBDIR)/$1_wrap.a
objcopy --redefine-sym userInput=$1_userInput $(LIBDIR)/$1_wrap.a
objcopy --redefine-sym updateCurrentState=$1_updateCurrentState $(LIBDIR)/$1_wrap.a
endef