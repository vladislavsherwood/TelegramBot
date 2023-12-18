package telegram

import (
	"encoding/json"
	"github.com/vladislavsherwood/TelegramBot/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	// Хост API сервиса тг
	host string
	// Базовый путь, с которого начинаются все запросы
	basePath string
	//
	client http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "getUpdates"
)

// New - Функция, создающая клиент
func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// newBasePath - Базовый путь, формирующийся из токена
func newBasePath(token string) string {
	return "bot" + token
}

// Updates - Получение сообщений клиентом
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	// Сообщение серверу, с какого апдейта (или сообщения) начинать получение информации.
	q.Add("offset" /*Конвертирует в string "Integer (I) to (to) ASCII (a)"*/, strconv.Itoa(offset))
	// Максимальное количество апдейтов, получаемое за один запрос
	q.Add("limit" /*Конвертирует в string "Integer (I) to (to) ASCII (a)"*/, strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chatId", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	// Дефер ошибки в конце функции doRequest
	defer func() { err = e.WrapIfErr("cant do request: %w", err) }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		// path.Join умно вставляет слэш когда надо
		Path: path.Join(c.basePath, method),
	}
	// Создание объекта запроса
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Енкодинг объекта запроса
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
