# Tableau Server CLI Tool

CLI tool for Tableau Server utilizing Tableau REST API.

Feature requests are welcome! :-)

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/davidlukac)


## Commands

| Command     | Description                                                                             |
|-------------|-----------------------------------------------------------------------------------------|
| create user | Create new user from given username.                                                    |
| delete user | Delete user by username, supports moving existing assets to another user.               |
| get user    | Get user info for given username, OR list all users. All users can be exported in YAML. |
| login       | Authenticate and provide token for further communication.                               |
| update user | Update existing user role by username, or read user(s) and role(s) from a YAML file.    |


## Configuration

Configuration is loaded from `.local.env` file:

```
TABLEAU_URL="https://tableau.my-domain.com/api/3.11"
TABLEAU_USERNAME="tableau-api-user"
TABLEAU_PASSWORD="super-secret-password"
TABLEAU_EXISTING_ASSETS_USER_NAME="john.smith"
LOG_LEVEL=INFO
```
