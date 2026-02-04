# SubTrack

SubTrack is a subscription tracker service that runs twice a day. It tracks subscription payments and notifies you via Telegram bot when the payment date is under 5 days.

## Features

- Track subscriptions with name, price, currency, cycle, and payment date
- Automatic notifications via Telegram for upcoming payments (< 5 days)
- Automatic payment date updates based on subscription cycle (monthly/yearly)
- CLI interface for managing subscriptions
- Background service for automated checking

## Setup

1. Clone the repository
2. Install dependencies: `make deps`
3. Copy `.env.example` to `.env` and fill in your Telegram bot token and chat ID
4. Build the project: `make build`

## Configuration

Create a `.env` file with the following variables:

```
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
DB_PATH=subtrack.db
```

## Usage

### CLI Commands

Add a subscription:
```bash
./bin/subtrack-cli add "Netflix" 15.99 USD monthly 15-02-2025
```

List all subscriptions:
```bash
./bin/subtrack-cli list
```

Update a subscription:
```bash
./bin/subtrack-cli update 1 "Netflix" 19.99 USD monthly 15-03-2025
```

Delete a subscription:
```bash
./bin/subtrack-cli delete 1
```

Manually check upcoming payments:
```bash
./bin/subtrack-cli check
```

Check Telegram bot health:
```bash
./bin/subtrack-cli health
```

### Running the Service

Start the background service:
```bash
./bin/subtrack-service
```

The service will run automatically and check payments twice daily (at 9:00 AM and 9:00 PM).

## Makefile Commands

- `make build` - Build CLI and service binaries
- `make install` - Install CLI and service to GOPATH/bin
- `make run-service` - Run the service (requires building first)
- `make deps` - Install and tidy Go dependencies
- `make clean` - Remove build artifacts

## Date Format

All dates use the format: `DD-MM-YYYY` (e.g., 15-02-2025)

## Subscription Cycles

- `monthly` - Payment recurs every month
- `yearly` - Payment recurs every year
