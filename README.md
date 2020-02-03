# VRChat export & import Friends via Discord Bot
Discord bot to export friends and import them to a new account

## Build
```bash
go build .
```

## Usage
Adjust the config.json
```json
{
    "discordbot": {
        "botname": "",
        "token": "",
        "email": "",
        "password": ""
    }
}
```
Commands
```cmd
!export "username" "password"

To import a generated friends list file, attach it and write as a comment
!import "username" "password"

```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
