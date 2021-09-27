## Refsys
Referral system where user earn points for every 3 friends that:
- uses their referral code to register
- transfers over 200 points

## Usage
### Requirements
- Docker and docker-compose
- Make (Optional)
- [Migrate](https://github.com/golang-migrate/migrate) - used to manage the migrations and seeders. Optionally, you can `source` the sql files in the `.sql` directory.
### App setup
- Copy the example environment file (`.env.example`) to `.env` i.e.,
```bash
$ cp .env.example .env
```
and fill the `.env` file with your actual credentials.

- Start the application with:
```bash
$ docker-compose up --build
```
This will start the PostgreSQL server on `http://localhost:5432` and the API server on `http://localhost:8000` 
### Database setup/seeding
The migration files needed to setup the database are stored in the `sql` folder. Run the make command below to apply them:
```bash
$ make migrate-up
```
The command above will:
- create a users table to store users
- create a transactions table to store transactions
- create a wallets table to store user wallet information
- create a payouts table to store referral payouts that are due.
- insert a seed transaction for the user with id 1 (since no user would have enough balance to make a transfer initially).
## API Docs
### POST /api/register
Signs up a new user
#### Request parameters
- `username=[String]`: Username of the user. Will return an error if the username already exists
- `password=[String]`: The user's password
- `referrer=[String]`: Optional. Referral code of the user that referred this user.
#### Response

### POST /api/transaction
Initiates a P2P transfer
#### Request parameters
- `sender_id=[Number]`: user ID of the sender.
- `recipient_id=[Number]`: user ID of the recipient
- `amount=[Number]`: Amount to be transferred. 
- `description=[String]`: Optional. Description of the transaction