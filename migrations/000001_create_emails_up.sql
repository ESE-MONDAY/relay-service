CREATE TABLE emails (

    id UUID PRIMARY KEY,

    sender TEXT NOT NULL,

    recipient TEXT NOT NULL,

    subject TEXT NOT NULL,

    text_body TEXT,

    html_body TEXT,

    status TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()

);