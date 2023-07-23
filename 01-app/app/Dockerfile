# Use official Node.js as a parent image
FROM node:20-slim

# Set the working directory in the container to /app
WORKDIR /app

# Copy the package.json and package-lock.json to the working directory
COPY package*.json ./

# Install the application dependencies inside the container
RUN npm install

# Copy the rest of the application to the working directory
COPY . .

# Make port 3000 available to the world outside this container
EXPOSE 3000

# Define the entry point of the container
ENTRYPOINT [ "node", "app.js" ]
