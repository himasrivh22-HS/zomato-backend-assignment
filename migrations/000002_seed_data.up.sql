-- ================= USER =================
INSERT INTO users (id, name, email, password, created_at)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  'Test User',
  'test@example.com',
  '$2a$12$KbQiP0FvKpYfYh6lqV8rQeU7v3y3QX7sP6Zb0lY0Yz0l0QfKQk9lK',
  NOW()
);

-- ================= PROJECT =================
INSERT INTO projects (id, name, description, owner_id, created_at)
VALUES (
  '22222222-2222-2222-2222-222222222222',
  'Sample Project',
  'Demo project',
  '11111111-1111-1111-1111-111111111111',
  NOW()
);

-- ================= TASKS =================
INSERT INTO tasks (id, title, status, priority, project_id, assignee_id, creator_id, created_at, updated_at)
VALUES
(
  '33333333-3333-3333-3333-333333333333',
  'Task 1',
  'todo',
  'low',
  '22222222-2222-2222-2222-222222222222',
  '11111111-1111-1111-1111-111111111111',
  '11111111-1111-1111-1111-111111111111',
  NOW(),
  NOW()
),
(
  '44444444-4444-4444-4444-444444444444',
  'Task 2',
  'in_progress',
  'medium',
  '22222222-2222-2222-2222-222222222222',
  NULL,
  '11111111-1111-1111-1111-111111111111',
  NOW(),
  NOW()
),
(
  '55555555-5555-5555-5555-555555555555',
  'Task 3',
  'done',
  'high',
  '22222222-2222-2222-2222-222222222222',
  NULL,
  '11111111-1111-1111-1111-111111111111',
  NOW(),
  NOW()
);