CREATE TABLE countries (
    guid VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255),
    code VARCHAR(2),
    continent VARCHAR(2),
    created_at TIMESTAMP Default current_timestamp,
    updated_at TIMESTAMP
);


CREATE TABLE cities (
    guid VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255),
    country_id VARCHAR(36) REFERENCES countries(guid) ON DELETE CASCADE,
    city_code VARCHAR(255),
    latitude VARCHAR(255),
    longitude VARCHAR(255),
    "offset" VARCHAR(255),
    timezone_id VARCHAR(36),
    country_name VARCHAR(255),
    created_at TIMESTAMP Default current_timestamp,
    updated_at TIMESTAMP
);

CREATE TABLE buildings (
    guid uuid PRIMARY KEY,
    title VARCHAR(255) ,
    country_id VARCHAR(36)  ,
    city_id VARCHAR(36)  ,
    latitude DECIMAL(9, 6)  ,
    longitude DECIMAL(9, 6)  ,
    radius VARCHAR(233),
    image VARCHAR(255),
    address VARCHAR(255),
    timezone_id VARCHAR(36)  ,
    country VARCHAR(255)  ,
    city VARCHAR(255)  ,
    search_text VARCHAR(255)  ,
    code VARCHAR(255)  ,
    product_count INT  ,
    gmt VARCHAR(6),
    created_at TIMESTAMP default current_timestamp,
    updated_at TIMESTAMP
);
