## Trader Bot

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
* Need help? ask bot with /help command **CURRENTLY NOT AVAILABLE**

---

## Project overview

![overview](./arch_overview2.png)

--- 

## How to run project

### Prerequisities:

* Installed docker

* Installed ```make``` utility tool (not required, but preferred)

* Required environment variables:

    * You can create all needed env template files with ```make quick-setup```

    * Or ```chmod +x quick_env_setup && ./quick_env_setup```

    * Or manually, described below:

```md
**File:** .app.env

**Content**
GUILD_ID=
APPLICATION_ID=
BOT_TOKEN=
PUBLIC_KEY=
**EOF Content**

**File:** .db.name

**Content**
random
**EOF Content**

**File:** .db.password

**Content**
random
**EOF Content**

**File:** .db.user

**Content**
random
**EOF Content**

---

* Description of all needed environment variables / its content:
    ***.app.env***

    * In **DEV** environment, application commands are registered in **guild**, 
      therefore the **GUILD_ID** environment is **required**.

    * In **PROD** environment, application commands are registered **globally**,
      therefore the **GUILD_ID** environment is **NOT** required.
    
    * **APPLICATION_ID, BOT_TOKEN, PUBLIC_KEY** are environment variables
    provided by **Discord**. 
    Head to *discord developer portal* and create app in order to obtain it.

    * **GUILD_ID** is also obtainable from **Discord**, but from *server* where app is installed.

    ***.db.name***

    * Provides database name for db provider via docker secret

    ***.db.password***

    * Provides database password for db provider via docker secret

    ***.db.user***

    * Provides database user for db provider via docker secret
```

---

### Running

* Via ```make```

    * ```make build-debug``` Debug run-time environment

    * ```make build-prod``` Prod run-time environment

* Via ```docker compose```

    * ```docker compose -f compose.dev.yaml up -d`` Debug run-time environment

    * ```docker compose -f compose.prod.yaml up -d``` Prod run-time environment
