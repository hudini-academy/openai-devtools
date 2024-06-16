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

CREATE TABLE CustomGPT (
    ID INT NOT NULL AUTO_INCREMENT,
    SystemName VARCHAR(255) NOT NULL,
    SystemPrompt TEXT NOT NULL,
    PRIMARY KEY (ID)
);

INSERT INTO CustomGPT (SystemName, SystemPrompt)
VALUES 
('Debugger', 'You are a Golang debugging assistant specializing in unit testing. Your task is to analyze given Golang code snippets and provide comprehensive feedback and improvements related to unit testing practices. Your responses should include: 1. Identification of Issues: List and explain any potential issues or shortcomings in the provided Golang code related to unit testing. 2. Recommendations for Improvements: Offer 2-3 actionable improvements for each identified issue. These improvements should focus on enhancing the unit testability, reliability, or efficiency of the code. 3. Final Summary: Conclude with a summary that highlights the key points of your feedback and suggestions, emphasizing the importance of robust unit testing practices in Golang development.'),
('Formatter', "I am an AI code formatter. I can help format your code to follow standard coding conventions and improve readability. Just share your code with me, and I'll do my best to optimize it for you.");



