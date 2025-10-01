QT       += core gui

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

CONFIG += c++17 static release

QMAKE_CXXFLAGS += -Wall -Wextra -Werror -std=c++17

TEMPLATE = app

LIBS_PATH = $$PWD/../../../build/libs

contains(DEFINES, SNAKE) {
    LIBS += $$LIBS_PATH/snake.a
    TARGET = snake_desktop
} else:contains(DEFINES, TETRIS) {
    LIBS += $$LIBS_PATH/tetris.a
    TARGET = tetris_desktop
} else:contains(DEFINES, RACE) {
    LIBS += $$LIBS_PATH/race.a
    TARGET = race_desktop
} else:contains(DEFINES, REST_CLIENT) {
    LIBS += $$LIBS_PATH/client.a -lcurl
    QMAKE_CXXFLAGS += -lcurl
    TARGET = brickgame
}

SOURCES += \
    gamefield.cpp \
    main.cpp \
    gamewindow.cpp

HEADERS += \
    gamefield.h \
    gamewindow.h \
    ../../brick_game/brickgame_api.h \

# Default rules for deployment.
qnx: target.path = /tmp/$${TARGET}/bin
else: unix:!android: target.path = /opt/$${TARGET}/bin
!isEmpty(target.path): INSTALLS += target

RESOURCES += resources/resourses.qrc
