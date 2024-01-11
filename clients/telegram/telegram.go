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

// Client представляет клиент Telegram API.
type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

// New создает новый экземпляр клиента Telegram API.
func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// newBasePath генерирует базовый путь из токена.
func newBasePath(token string) string {
	return "bot" + token
}

// Updates получает сообщения для клиента.
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	// Указывает серверу ID апдейта или сообщения, с которого начать получение информации.
	q.Add("offset", strconv.Itoa(offset))
	// Максимальное количество апдейтов, получаемое за один запрос.
	q.Add("limit", strconv.Itoa(limit))

	// doRequest выполняет HTTP-запрос для получения информации с использованием метода getUpdatesMethod и параметров q.
	// В случае ошибки в процессе выполнения запроса возвращается nil и ошибка.
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	// json.Unmarshal распаковывает данные, полученные от сервера, и сохраняет их в переменной res.
	// В случае ошибки в процессе распаковки возвращается nil и ошибка.
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

// SendMessage отправляет сообщение клиенту, используя метод API sendMessageMethod.
// chatId - идентификатор чата, куда будет отправлено сообщение.
// text - текст сообщения, которое будет отправлено.
// Создается объект url.Values для передачи параметров запроса, включая chatId и text.
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	// Вызывается функция doRequest, которая отправляет запрос на сервер Telegram с методом sendMessageMethod и параметрами q.
	_, err := c.doRequest(sendMessageMethod, q)
	// Если произошла ошибка при выполнении запроса, функция возвращает ошибку, обернутую в дополнительное описание.
	if err != nil {
		return e.Wrap("не удалось отправить сообщение", err)
	}
	return nil
}

// doRequest выполняет HTTP-запрос к Telegram API.
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	// Обработка ошибок в конце функции doRequest.
	defer func() { err = e.WrapIfErr("не удалось выполнить запрос: %w", err) }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		// path.Join вставляет слеш, если необходимо.
		Path: path.Join(c.basePath, method),
	}
	// Создание объекта запроса.
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Кодирование объекта запроса.
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
