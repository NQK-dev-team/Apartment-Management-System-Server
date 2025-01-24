# Appartment-Management-System-Server

# Khoa download postgreSQL version 16.6 from postgreSQL site

# Khoa guide to running the Server


## Setup

Make sure to install dependencies:
```bash
#Make sure to install go: https://medium.com/novai-go-programming-101/step-by-step-guide-to-installing-go-golang-on-windows-linux-and-mac-cab22d0320ef 
#Install Dependencies: Ensure you have all the required dependencies installed. You can do this by running:
go mod tidy
#Set up your database: https://drive.google.com/open?id=11YERYXUiPxKusKegHhlN2C42ltfAZA0l&usp=drive_fs
#Fill the .env file accordingly

#Set Up Environment Variables: Copy the .env.dist file to .env and fill in the necessary environment variables.
cp .env.dist .env
```

## Development Server

Start the development server on

```bash
#Run the Server: Use the following command to start the server:
go run cmd/server/main.go
```