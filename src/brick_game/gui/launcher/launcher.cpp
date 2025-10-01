#include <chrono>
#include <iostream>
#include <thread>

#include "../../client/client.h"

int main() {
  BrickClient* client = getClient();
  if (!client) return 1;

  std::vector<std::string> gameList;
  while (true) {
    try {
      gameList = client->getGamesList();
      break;
    } catch (const std::exception& e) {
      std::cerr << "Cannot connect to server, retrying in 5 seconds...\n";
      std::this_thread::sleep_for(std::chrono::seconds(5));
    }
  }

  std::cout << "Select a game from server:\n";
  for (size_t i = 0; i < gameList.size(); i++) {
    std::cout << i + 1 << ". " << gameList[i] << "\n";
  }
  std::cout << "\n";
  int choice = -1;

  if (!(std::cin >> choice) || choice < 1 || choice > 3) {
    std::cout << "Incorrect input\n";
  } else {
    std::cin.ignore(1000, '\n');
    std::cout << "\nSelect an interface:\n";
    std::cout << "1. Console\n2. Desktop\n";

    int interface = -1;
    if (!(std::cin >> interface) || interface < 1 || interface > 2) {
      std::cout << "Incorrect input\n";
    } else {
      client->chooseGame(choice);
      if (interface == 1)
        system("./build/console/brickgame");
      else if (interface == 2)
        system("./build/desktop/brickgame");
    }
  }

  return 0;
}
