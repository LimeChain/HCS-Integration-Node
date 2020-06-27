# General
## Overview
## Running the node
### Create the crypto keys
Either run `./cryptogen.sh` or place a PEM encoded private key in `config/key.pem`

If you want to have another keypair, just rename the file and run cryptogen again.

### Creathing HCS Topic
Example `./createhcstopic.sh`

### Run your mongo database
`docker run --name hedera-mongo -d -p 27017:27017 -v ~/data:/data/db mongo`

### Create .env file
Create file named `.env` based on the `.env.example`

### Starting the node
Run `./start.sh`

## Running the second peer
1. Backup your `key.pem` and run `./cryptogen`. Rename the new `key.pem` to another name and restore the previous one.
2. Run mongo with different name port and data storage path
3. Create `peer2.env` based on `.env.example`
4. Run `./start-peer2.sh`
