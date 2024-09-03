# Use a Node.js base image
FROM node:22

# Set the working directory
WORKDIR /app

# Copy the package.json and yarn.lock/npm package files
COPY package*.json ./

# Install dependencies
RUN npm install -g expo-cli
RUN npm install -g @expo/ngrok

# Copy the project files
COPY . .

# Start the Expo server on the specified port
CMD ["npx", "expo", "start", "--tunnel", "--port", "19006"]
