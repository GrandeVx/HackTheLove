CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Table SurveysResponse
CREATE TABLE IF NOT EXISTS surveys (
    id_survey UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    response VARCHAR(11) NOT NULL, -- Risposta, massimo 11 caratteri
    date_created TIMESTAMP DEFAULT NOW()
);

-- Table Users
CREATE TABLE IF NOT EXISTS users (
    email VARCHAR(100) PRIMARY KEY, -- Primary key
    id_survey UUID REFERENCES surveys(id_survey) ON DELETE SET NULL, -- Id of the survey
    name VARCHAR(255) NOT NULL, -- Name
    surname VARCHAR(255) NOT NULL, -- Surname
    phone VARCHAR(15), -- Phone
    sex BOOLEAN, -- Sex
    bio TEXT, -- Biography
    age INTEGER, -- Age
    section VARCHAR(1), -- Section
    date_created TIMESTAMP DEFAULT NOW() -- Date of creation
);

-- Table Images
CREATE TABLE IF NOT EXISTS images (
    id_image SERIAL PRIMARY KEY, -- Primary key
    email_user VARCHAR(100) REFERENCES users(email) ON DELETE CASCADE, -- User reference
    lo_oid OID NOT NULL, -- File reference
    uploaded_at TIMESTAMP DEFAULT NOW(), -- Date of upload
    metadata JSONB -- Extra information
);


CREATE TABLE IF NOT EXISTS matches (
    email_user1 VARCHAR(100) REFERENCES users(email) ON DELETE CASCADE,
    email_user2 VARCHAR(100) REFERENCES users(email) ON DELETE CASCADE,
    like_user1 BOOLEAN NOT NUll DEFAULT FALSE, -- True like - False dislike
    like_user2 BOOLEAN  NOT NUll DEFAULT FALSE, -- True like - False dislike
    compatibility FLOAT NOT NULL,
    date_created TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (email_user1, email_user2)
);
