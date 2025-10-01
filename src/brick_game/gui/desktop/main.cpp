#include <stdlib.h>
#include <time.h>

#include "gamewindow.h"

int main(int argc, char *argv[]) {
  QApplication a(argc, argv);
  GameWindow window;

#ifdef REST_CLIENT
  client();
#endif

  window.show();
  return a.exec();
}
