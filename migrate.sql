ALTER TABLE articles ADD COLUMN subcategory VARCHAR(50) DEFAULT '' AFTER category;
ALTER TABLE articles ADD COLUMN study_status VARCHAR(20) DEFAULT '' AFTER subcategory;
ALTER TABLE articles ADD COLUMN is_pinned TINYINT(1) DEFAULT 0 AFTER study_status;