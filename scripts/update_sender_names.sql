-- Update sender_name in cards table based on known user IDs
-- This script fixes existing data that has empty sender_name

UPDATE cards SET sender_name = 'Youshan Li (SZX)' WHERE sender_id = 'd3dc438c-1de2-49dc-ae8b-8dffb7fcb2cc' AND (sender_name IS NULL OR sender_name = '');
UPDATE cards SET sender_name = 'Zehui Lin (SZX)' WHERE sender_id = '02129f90-2224-486b-859f-b60c2da6b826' AND (sender_name IS NULL OR sender_name = '');
UPDATE cards SET sender_name = 'Christopher Li (SZX)' WHERE sender_id = 'cd7e9fb0-d844-441e-aac4-b8c21d36ef7c' AND (sender_name IS NULL OR sender_name = '');
UPDATE cards SET sender_name = 'Shawine Wu (SZX)' WHERE sender_id = '3ef2d96d-35b0-42ee-8ace-e4e7880f12c0' AND (sender_name IS NULL OR sender_name = '');
UPDATE cards SET sender_name = 'Alex Zu (SZX)' WHERE sender_id = 'f1f7264d-5a6c-4fef-82e9-a122f9a332d0' AND (sender_name IS NULL OR sender_name = '');
UPDATE cards SET sender_name = 'Bingye Li (SZX)' WHERE sender_id = 'dfaa0ed8-2473-403b-8fdb-7fbbb38dfd06' AND (sender_name IS NULL OR sender_name = '');

-- Verify the update
SELECT sender_id, sender_name, COUNT(*) as count FROM cards GROUP BY sender_id, sender_name ORDER BY count DESC;
