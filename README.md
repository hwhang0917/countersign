# Countersign CLI

> A simple CLI application to generate OTP-like countersigns to counter social engineering.

## Install

1. Clone repository

```bash
git clone https://github.com/hwhang0917/countersign.git
```

2. Build

```bash
make
```

3. Setup environment variables (`.env`)

```dotenv
# dist/.env
API_KEY=<your_api_key>
```

4. Run binary

```bash
./dist/countersign
```
