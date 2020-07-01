# General
## Overview
## Running the node
### Create the crypto keys
Either run `./cryptogen.sh` or place a PEM encoded private key in `config/key.pem`

If you want to have another keypair, just rename the file and run cryptogen again.

### Creathing HCS Topic with threshold keys
Example `A_PUB_KEY=302a300506032b6570032100086e579c72b037e72bddc3d5c8af5e7b5e5269ab6bb025792a480940ba501b16 B_PUB_KEY=302a300506032b65700321005685758381e67fdaf28b6a992e1e725a707e95280d89eae47fc5132271dc2b1a ./createhcstopic.sh`

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
