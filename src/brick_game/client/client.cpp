#include "client.h"

static size_t WriteCallback(void* contents, size_t size, size_t nmemb,
                            void* userp) {
  ((std::string*)userp)->append((char*)contents, size * nmemb);
  return size * nmemb;
}

BrickClient::BrickClient(std::string url) {
  if (url == "")
    m_serverUrl = defaultServerUrl;
  else
    m_serverUrl = url;
  initFields();
}

BrickClient::~BrickClient() { freeFields(); }

void BrickClient::initFields() {
  if (m_gameState.field == nullptr) {
    m_gameState.field = new int*[HEIGHT]();
    for (int i = 0; i < HEIGHT; i++) {
      m_gameState.field[i] = new int[WIDTH]();
    }
  }

  if (m_gameState.next == nullptr) {
    m_gameState.next = new int*[4]();
    for (int i = 0; i < 4; i++) {
      m_gameState.next[i] = new int[4]();
    }
    m_gameState.pause = START_;
  }
}

void BrickClient::freeFields() {
  if (m_gameState.field != nullptr) {
    for (int i = 0; i < HEIGHT; i++) {
      delete[] m_gameState.field[i];
      m_gameState.field[i] = nullptr;
    }
    delete[] m_gameState.field;
    m_gameState.field = nullptr;
  }

  if (m_gameState.next != nullptr) {
    for (int i = 0; i < 4; i++) {
      delete[] m_gameState.next[i];
      m_gameState.next[i] = nullptr;
    }
    delete[] m_gameState.next;
    m_gameState.next = nullptr;
  }
}

std::vector<std::string> BrickClient::getGamesList() {
  std::string url = m_serverUrl + "/games";
  std::string response = httpGet(url);
  auto j = json::parse(response);

  std::vector<std::string> games;
  for (auto& g : j["games"]) {
    games.push_back({g["name"]});
  }
  return games;
}

void BrickClient::chooseGame(int gameID) {
  std::string url = m_serverUrl + "/games/" + std::to_string(gameID);
  httpPost(url, "{}");
}

void BrickClient::sendAction(const UserAction_t& action) {
  json j = {{"action_id", static_cast<int>(action)}, {"hold", false}};

  try {
    httpPost(m_serverUrl + "/actions", j.dump());
  } catch (const std::exception& e) {
  }
}

GameInfo_t BrickClient::getState() {
  std::string url = m_serverUrl + "/state";
  try {
    std::string response = httpGet(url);
    auto j = json::parse(response);

    m_gameState.score = j["score"];
    m_gameState.high_score = j["high_score"];
    m_gameState.level = j["level"];
    m_gameState.speed = j["speed"];
    m_gameState.pause = j["pause"];

    auto field = j["field"].get<std::vector<std::vector<int>>>();
    for (int i = 0; i < HEIGHT; i++) {
      for (int j = 0; j < WIDTH; j++) {
        m_gameState.field[i][j] = field[i][j];
      }
    }

    auto next = j["next"].get<std::vector<std::vector<int>>>();
    for (int i = 0; i < NEXT_SIZE; i++) {
      for (int j = 0; j < NEXT_SIZE; j++) {
        m_gameState.next[i][j] = next[i][j];
      }
    }
  } catch (const std::exception& e) {
  }

  return m_gameState;
}

std::string BrickClient::httpGet(const std::string& url) {
  CURL* curl = curl_easy_init();
  if (!curl) throw std::runtime_error("CURL init failed");

  std::string readBuffer;
  curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, &readBuffer);

  CURLcode res = curl_easy_perform(curl);
  if (res != CURLE_OK) {
    curl_easy_cleanup(curl);
    throw std::runtime_error(std::string("httpGet failed: ") +
                             curl_easy_strerror(res));
  }

  curl_easy_cleanup(curl);

  return readBuffer;
}

void BrickClient::httpPost(const std::string& url, const std::string& payload) {
  CURL* curl = curl_easy_init();
  if (!curl) throw std::runtime_error("CURL init failed");

  struct curl_slist* headers = nullptr;
  headers = curl_slist_append(headers, "Content-Type: application/json");

  curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
  curl_easy_setopt(curl, CURLOPT_POSTFIELDS, payload.c_str());
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
  std::string readBuffer;
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, &readBuffer);

  CURLcode res = curl_easy_perform(curl);
  curl_slist_free_all(headers);

  if (res != CURLE_OK) {
    curl_easy_cleanup(curl);
    throw std::runtime_error(std::string("httpPost failed: ") +
                             curl_easy_strerror(res));
  }

  curl_easy_cleanup(curl);
}

GameInfo_t updateCurrentState() { return getClient()->getState(); }

void userInput(UserAction_t action, bool hold) {
  getClient()->sendAction(action);
  (void)hold;
}

BrickClient* getClient() {
  static BrickClient client{};
  return &client;
}

void client() { getClient(); }