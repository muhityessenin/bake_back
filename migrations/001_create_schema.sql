-- +goose Up
-- Drop existing tables to recreate with new structure
DROP TABLE IF EXISTS sales_data;
DROP TABLE IF EXISTS marketing_data;
DROP TABLE IF EXISTS report_dates;

CREATE TABLE marketing_sources (
                                   id SERIAL PRIMARY KEY,
                                   name VARCHAR(50) NOT NULL UNIQUE,
                                   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sales_teams (
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(50) NOT NULL UNIQUE,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE marketing_data (
                                id SERIAL PRIMARY KEY,
                                date DATE NOT NULL,
                                source_id INTEGER NOT NULL REFERENCES marketing_sources(id),
                                expense DECIMAL(15,2) NOT NULL DEFAULT 0,
                                leads INTEGER NOT NULL DEFAULT 0,
                                trials_scheduled INTEGER NOT NULL DEFAULT 0,
                                trials_conducted INTEGER NOT NULL DEFAULT 0,
                                payments INTEGER NOT NULL DEFAULT 0,
                                total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
                                is_saved BOOLEAN NOT NULL DEFAULT FALSE,
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                UNIQUE (date, source_id)
);

CREATE TABLE sales_data (
                            id SERIAL PRIMARY KEY,
                            date DATE NOT NULL,
                            team_id INTEGER NOT NULL REFERENCES sales_teams(id),
                            leads INTEGER NOT NULL DEFAULT 0,
                            trials_scheduled INTEGER NOT NULL DEFAULT 0,
                            trials_conducted INTEGER NOT NULL DEFAULT 0,
                            payments INTEGER NOT NULL DEFAULT 0,
                            total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
                            kaspi_refund DECIMAL(15,2) NOT NULL DEFAULT 0,
                            is_saved BOOLEAN NOT NULL DEFAULT FALSE,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            UNIQUE (date, team_id)
);

-- Initial data for sources and teams
INSERT INTO marketing_sources (name) VALUES
                                         ('Facebook-1'), ('Facebook-2'), ('Facebook-3'), ('Facebook-4'), ('Facebook-5'), ('Facebook-6'),
                                         ('Instagram'), ('TikTok'), ('Другие источники');

INSERT INTO sales_teams (name) VALUES
                                   ('Команда Нұрғиса'), ('Команда Даниял'), ('Команда Дамир'), ('Команда Абылай'), ('Команда Тоғжан'), ('Команда Айша Тараз');

-- +goose Down
DROP TABLE IF EXISTS sales_data;
DROP TABLE IF EXISTS marketing_data;
DROP TABLE IF EXISTS sales_teams;
DROP TABLE IF EXISTS marketing_sources;