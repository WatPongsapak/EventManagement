CREATE TABLE activities(
    id serial PRIMARY KEY,
    name TEXT NOT NULL,
    location TEXT,
    speaker TEXT,
    description TEXT,
    max_joinable INTEGER,
    start_datetime TIMESTAMP,
    end_datetime TIMESTAMP,
    round INTEGER
);

CREATE TABLE admins(
    username TEXT PRIMARY KEY,
    password TEXT
);

CREATE TABLE pinactivities(
    activities_id INTEGER REFERENCES activities(id),
    employee_code TEXT,
    name TEXT,
    phone TEXT,
    PRIMARY KEY (activities_id ,employee_code)
);
