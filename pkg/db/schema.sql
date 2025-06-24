CREATE DATABASE demo_orders;
\c demo_orders

CREATE USER demo_user WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE demo_orders TO demo_user;
GRANT ALL ON ALL TABLES IN SCHEMA public TO demo_user;

CREATE TABLE delivery
(
    DeliveryId SERIAL PRIMARY KEY,
    FirstName VARCHAR(40),
    Phone VARCHAR(15),
    Zip VARCHAR(15),
    City VARCHAR(30),
    Address VARCHAR(30),
    Region VARCHAR(30),
    Email VARCHAR(30)
);

CREATE TABLE payment
(
    TransactionName CHARACTER VARYING(40) PRIMARY KEY,
    RequestId CHARACTER VARYING(30),
    Currency CHARACTER(3),
    ProviderName VARCHAR(15),
    Amount INTEGER,
    PaymentDt BIGINT,
    Bank VARCHAR(15),
    DeliveryCost INTEGER,
    GoodsTotal INTEGER,
    CustomFee INTEGER
);

CREATE TABLE item
(
    ChrtId SERIAL PRIMARY KEY,
    TrackNumber VARCHAR(30),
    Price INTEGER,
    RId VARCHAR(30),
    ItemName VARCHAR(15),
    Sale INTEGER,
    ItemSize VARCHAR(30),
    TotalPrice INTEGER,
    NmId INTEGER,
    Brand VARCHAR(30),
    ItemStatus INTEGER
);

CREATE TABLE orders
(
    OrderUid CHARACTER VARYING(30) PRIMARY KEY,
    TrackNumber CHARACTER VARYING(30),
    Entry CHARACTER VARYING(30),
    Locale CHARACTER(2),
    CustomerId CHARACTER VARYING(30),
    DeliveryService CHARACTER VARYING(30),
    ShardKey CHARACTER VARYING(30),
    SmId INTEGER,
    DataCreated TIMESTAMP,
    OofShard CHARACTER VARYING(30),
    DeliveryId INTEGER REFERENCES delivery (DeliveryId),
    PaymentId CHARACTER VARYING(40) REFERENCES payment (TransactionName)
);

CREATE TABLE items
(
    ItemId INTEGER REFERENCES item (ChrtId),
    OrderId CHARACTER VARYING(30) REFERENCES orders (OrderUid),
    PRIMARY KEY (ItemId, OrderId)
);