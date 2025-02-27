INSERT INTO gateways (name, data_format_supported, priority)
VALUES ('Stripe', 'JSON,SOAP', 0),
       ('PayPal', 'JSON', 1),
       ('Braintree', 'SOAP', 2),
       ('Adyen', 'JSON', 0),
       ('Authorize.Net', 'JSON,SOAP', 1),
       ('Worldpay', 'SOAP', 2),
       ('Checkout.com', 'JSON', 0),
       ('Square', 'JSON', 1),
       ('Payoneer', 'JSON,SOAP', 2),
       ('2Checkout', 'JSON', 0);


INSERT INTO gateway_countries (gateway_id, country_id)
VALUES (1, 1),  -- Stripe, USA
       (1, 2),  -- Stripe, Canada
       (2, 1),  -- PayPal, USA
       (2, 3),  -- PayPal, UK
       (3, 4),  -- Braintree, Germany
       (3, 5),  -- Braintree, France
       (4, 6),  -- Adyen, Japan
       (4, 7),  -- Adyen, Australia
       (5, 1),  -- Authorize.Net, USA
       (5, 8),  -- Authorize.Net, China
       (6, 3),  -- Worldpay, UK
       (6, 9),  -- Worldpay, India
       (7, 1),  -- Checkout.com, USA
       (7, 10), -- Checkout.com, Brazil
       (8, 2),  -- Square, Canada
       (8, 4),  -- Square, Germany
       (9, 1),  -- Payoneer, USA
       (9, 5),  -- Payoneer, France
       (10, 3), -- 2Checkout, UK
       (10, 6); -- 2Checkout, Japan