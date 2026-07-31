INSERT INTO projects (id, name, description)
VALUES ('00000000-0000-0000-0000-000000000000', 'Default Project', 'Main demo project for IDP')
ON CONFLICT (id) DO NOTHING;
