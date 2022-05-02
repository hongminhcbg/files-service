drop table if exist `templates`;
CREATE TABLE `file_service`.`templates` (
    id INT(20) primary key AUTO_INCREMENT,
    name VARCHAR(255) COMMENT 'name of file',
    content BLOB  NOT NULL COMMENT 'Raw content of file'
);

insert into `file_service`.`templates`(`name`, `content`) VALUES ('welcome template', LOAD_FILE('/var/lib/mysql-files/welcome.html'));
