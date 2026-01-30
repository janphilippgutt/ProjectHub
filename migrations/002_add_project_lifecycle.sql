ALTER TABLE projects
ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE INDEX idx_projects_deleted_at
ON projects(deleted_at);
