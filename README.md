## Trader Bot

* Bot for trading items
```
Want to create market for your discord? 

Would you like to declare what items ppl will trading on your server?

In this type of implementation, discord server / guild **administrator**,
needs to register items for which offers will be created with **/item-register**

You can declare up to 25 items (limit by discord), for each item drop down list
with items is displayed on discord for user when writing command 
to add / update / remove offer.

There is also rate limit for registering commands per guild, so be cautious.
When writing this project limit is **200** and resets after 24h.
```
* Simple interface, send commands with /command
    Actual supported commands:
        /item-register
            add
            remove
            list 
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
      therefore the **GUILD_ID** environment is **required** 
      and **WITH_ITEM_REGISTRAR_COMMAND_REGISTER=true** ONLY IF 
      you want to register item registrar command on start. 
      You should specify it once, then set it to false,
      because there is rate limit for registering commands on discord.

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
