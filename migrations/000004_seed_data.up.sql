INSERT INTO companies (name, legal_name, industry, data_source) VALUES
('Google', 'Google LLC', 'Technology', 'manual'),
('Apple', 'Apple Inc.', 'Technology', 'manual'),
('Microsoft', 'Microsoft Corporation', 'Technology', 'manual');

INSERT INTO addresses (company_id, address_line1, city, state, postal_code, country, is_primary) VALUES
(1, '1600 Amphitheatre Parkway', 'Mountain View', 'CA', '94043', 'USA', 1),
(2, '1 Apple Park Way', 'Cupertino', 'CA', '95014', 'USA', 1),
(3, '1 Microsoft Way', 'Redmond', 'WA', '98052', 'USA', 1);