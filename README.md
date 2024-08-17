# Telegram Translator Bot

This bot automatically translates messages in Telegram chats. It uses the DeepL API for translations and can be configured to work with specific groups and users.

## Features

- Automatic message translation
- Configurable target language
- Whitelist for allowed groups and users
- Customizable API endpoint
- Option to ignore specific languages

## Usage

### Command Line

You can run the bot using Docker with the following command:

```bash
./bot -token=your_bot_token_here -target=DE -api=http://127.0.0.1:1188/translate -ignore=ZH,EN,DE -groups=-1001652593847 -users=890315416,123456789
```

#### Available Options

- `-token`: Telegram Bot Token (required)
- `-target`: Target language for translation (default: "ZH")
- `-api`: API URL for translation service (default: "http://127.0.0.1:1188/translate")
- `-ignore`: Comma-separated list of languages to ignore (default: "ZH")
- `-groups`: Comma-separated list of allowed group IDs
- `-users`: Comma-separated list of allowed user IDs

#### Examples

1. Basic usage (only setting bot token):
   ```bash
   ./bot -token=your_bot_token_here
   ```

2. Set target language and API URL:
   ```bash
   ./bot -token=your_bot_token_here -target=EN -api=http://127.0.0.1:1188/translate
   ```

3. Set ignored languages and allowed groups:
   ```bash
   ./bot -token=your_bot_token_here -ignore=ZH,EN -groups=-1001652593847,-1002345678901
   ```

4. Set allowed users and target language:
   ```bash
   ./bot -token=your_bot_token_here -users=890315416,123456789 -target=FR
   ```

5. Full example with all parameters:
   ```bash
   ./bot -token=your_bot_token_here -target=DE -api=http://127.0.0.1:1188/translate -ignore=ZH,EN,DE -groups=-1001652593847 -users=890315416,123456789
   ```

### Docker Compose

1. Create a `compose.yaml` file with the following content:

```yaml
services:
  deeplx-bot:
    image: missuo/deeplx-bot:latest
    container_name: telegram-translator-bot
    environment:
      - BOT_TOKEN=your_bot_token_here
      - TARGET_LANG=ZH
      - API_URL=http://127.0.0.1:1188/translate
      - IGNORE_LANGS=ZH
      - ALLOWED_GROUPS=-1001652593847
      - ALLOWED_USERS=890315416
    restart: always

  # Optional: Uncomment the following to use the official DeepLX service
  # deeplx:
  #   image: ghcr.io/owo-network/deeplx:latest
  #   restart: always
  #   ports:
  #     - "1188:1188"
```

2. Replace `your_bot_token_here` with your actual Telegram bot token.

3. Adjust other environment variables as needed.

4. Run the following command in the same directory as your `compose.yaml`:

```bash
docker compose up -d
```

This will start the bot and the DeepL API proxy service in the background.

To view logs:
```bash
docker compose logs -f
```

To stop the services:
```bash
docker compose stop
```

## Configuration

You can configure the bot using environment variables or command-line arguments. When using Docker Compose, set the environment variables in the `compose.yaml` file.

| Environment Variable | Command-Line Argument | Description                            | Default Value              |
|----------------------|-----------------------|----------------------------------------|----------------------------|
| BOT_TOKEN            | -token                | Telegram Bot Token (required)          | -                          |
| TARGET_LANG          | -target               | Target language for translation        | ZH                         |
| API_URL              | -api                  | API URL for translation service        | http://deeplx:1188/translate |
| IGNORE_LANGS         | -ignore               | Comma-separated list of languages to ignore | ZH                    |
| ALLOWED_GROUPS       | -groups               | Comma-separated list of allowed group IDs | -                      |
| ALLOWED_USERS        | -users                | Comma-separated list of allowed user IDs | -                       |

## Support

For issues, questions, or contributions, please open an issue in the GitHub repository.

## License

[MIT License](LICENSE)