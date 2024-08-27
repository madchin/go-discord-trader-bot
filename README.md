### Trader Bot

* Bot for trading items

* Simple interface, send commands with /command
```
    Actual supported commands:
        /offer sell
            add
            remove
            update
            list
        /offer buy
            add
            remove
            update
            list
```
* Need help? ask bot with /help command **NOT AVAILABLE**

### Overview

![overview](./arch_overview2.png)


### How to run project

#### Prerequisities:

* Installed docker

* Installed ```make``` utility tool (not required, but preferred)

* .env file in project root dir with variables:
```md
GUILD_ID=
APPLICATION_ID=
BOT_TOKEN=
PUBLIC_KEY=

* Note:
    * In **DEV** environment, application commands are registered in **guild**, 
      therefore the **GUILD_ID** environment is **required**.

    * In **PROD** environment, application commands are registered **globally**,
      therefore the **GUILD_ID** environment is **NOT** required.
    
    * **APPLICATION_ID, BOT_TOKEN, PUBLIC_KEY** are environment variables
    provided by **Discord**. 
    Head to *discord developer portal* and create app in order to obtain it.

    * **GUILD_ID** is also obtainable from **Discord**, but from *server* where app is installed ofc.
```

#### Running

* Via ```make```

    * make build-dev ```dev environment```

    * make build-prod ```prod environment```

* Via ```docker compose```

    * docker compose run --rm --name DEBUG -e RUNTIME_ENVIRONMENT=DEV -d trader

    * docker compose up -d
