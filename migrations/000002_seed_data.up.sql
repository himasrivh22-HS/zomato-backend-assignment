INSERT INTO users (id, name, email, password)
VALUES (
  uuid_generate_v4(),
  'Test User',
  'test@example.com',
  '$2a$12$KbQiP0FvKpYfYh6lqV8rQeU7v3y3QX7sP6Zb0lY0Yz0l0QfKQk9lK'
);

INSERT INTO projects (id, name, description, owner_id)
SELECT uuid_generate_v4(), 'Sample Project', 'Demo project', id FROM users LIMIT 1;

INSERT INTO tasks (id, title, status, priority, project_id)
SELECT uuid_generate_v4(), 'Task 1', 'todo', 'low', id FROM projects LIMIT 1;

INSERT INTO tasks (id, title, status, priority, project_id)
SELECT uuid_generate_v4(), 'Task 2', 'in_progress', 'medium', id FROM projects LIMIT 1;

INSERT INTO tasks (id, title, status, priority, project_id)
SELECT uuid_generate_v4(), 'Task 3', 'done', 'high', id FROM projects LIMIT 1;