#include <curl/curl.h>

#include <iostream>
#include <stdexcept>
#include <string>

#include "../brickgame_api.h"
#include "lib/json.hpp"

using json = nlohmann::json;

const std::string defaultServerUrl = "http://localhost:8080/api";

class BrickClient {
 public:
  BrickClient(std::string url = "");
  ~BrickClient();

  std::vector<std::string> getGamesList();
  void chooseGame(int gameID);
  GameInfo_t getState();
  void sendAction(const UserAction_t& action);

 private:
  std::string httpGet(const std::string& url);
  void httpPost(const std::string& url, const std::string& payload);
  void initFields();
  void freeFields();

  std::string m_serverUrl;
  GameInfo_t m_gameState;
};

BrickClient* getClient();
