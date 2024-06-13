-- Create the database
CREATE DATABASE IF NOT EXISTS openaiusers;

-- Use the newly created database
USE openaiusers;

-- Create the 'users' table to store user information
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, -- Unique ID for each user
    name VARCHAR(255) NOT NULL, -- Name of the user
    email VARCHAR(255) NOT NULL, -- Email of the user
    password CHAR(60) NOT NULL -- Password of the user
);

-- Create the 'messages' table to store messages
CREATE TABLE messages (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, -- Unique ID for each message
    user_id INTEGER NOT NULL, -- Foreign key to reference the user who created the message
    title VARCHAR(255) NOT NULL, -- Title of the message
    message TEXT NOT NULL, -- Content of the message
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Foreign key constraint linking to the users table
);