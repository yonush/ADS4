-- +goose Up
-- +goose StatementBegin

-- Table to store courses created by admins, with a unique course code

CREATE TABLE "Courses" (
	"CourseCode"  VARCHAR(9),
	"Description" VARCHAR(255),
	"Level"  	  int, -- 1-9
    "Status"      VARCHAR(6) NOT NULL DEFAULT 'active',
	PRIMARY KEY("CourseCode"),
	UNIQUE("CourseCode"),
	CHECK (Status IN ('active','closed')),
	CHECK (Level IN (1,2,3,4,5,6,7,8,9))
);	
CREATE UNIQUE INDEX courses_byCourseCode ON courses(CourseCode);


-- Table to store exam offerings created by admins, with a unique ExamID that follows the format [year:4][semester:2][coursecode:9]

CREATE TABLE "Offerings" (
    "ExamID"      VARCHAR(15) NOT NULL, 
    "Year" 	      INTEGER,
    "Semester"    VARCHAR(2) NOT NULL DEFAULT 'S1',
    "CourseCode"  VARCHAR(9) NOT NULL,
    "Password"    VARCHAR(8) NOT NULL,
    "Status"      VARCHAR(6) NOT NULL DEFAULT 'active',
    "Coordinator" INTEGER,
    "OwnerID"     INTEGER,
    "Duration"    INTEGER,
    PRIMARY KEY("ExamID"),
    UNIQUE("ExamID"),
    FOREIGN KEY("Coordinator") REFERENCES "UserT"("UserID"),
    FOREIGN KEY("OwnerID") REFERENCES "UserT"("UserID"),
    FOREIGN KEY("CourseCode") REFERENCES "Courses"("CourseCode"),
	CHECK (Status IN ('active','closed')),
	CHECK ("Semester" IN ('S1','S2','S3')),
    CHECK ("Duration" > 29 and "Duration" < 241)
);

CREATE UNIQUE INDEX offering_byExamID ON offerings(ExamID);
CREATE INDEX offering_byCourseCode ON offerings(CourseCode);

-- Table to store exam learner information, with a unique StudentID for each learner (NON-SYSTEM USER)
CREATE TABLE "Learners" (
    "StudentID"   VARCHAR(8) NOT NULL,
    "Name"        VARCHAR NOT NULL,
    "Status"      VARCHAR NOT NULL DEFAULT 'active',
    UNIQUE("StudentID"),
    PRIMARY KEY("StudentID"),
    CHECK (Status IN ('active','inactive'))
);
CREATE UNIQUE INDEX learners_byStudentID ON learners(StudentID);

-- Table to store each learner's exam attempt, with a unique LearnerExamID for each attempt, and a foreign key reference to the Offerings table using ExamID

CREATE TABLE "Learnerexams" (
    "LearnerExamID" INTEGER,
    "StudentID"     VARCHAR(8) NOT NULL,
    "ExamID"        VARCHAR(15) NOT NULL,
    "StartTime"     TIME,
    "EndTime"       TIME,
    "Status"        VARCHAR(6) NOT NULL DEFAULT 'ready',
    "Grade"         INTEGER,
    PRIMARY KEY("LearnerExamID" AUTOINCREMENT),
    FOREIGN KEY("ExamID") REFERENCES "Offerings"("ExamID"),
    CHECK (Status IN ('ready', 'active', 'expire', 'closed', 'marked'))
    FOREIGN KEY("StudentID") REFERENCES "Learners"("StudentID")    
);
CREATE INDEX learnerexams_byCourseCode ON learnerexams(StudentID);
CREATE INDEX learnerexams_byExamID ON learnerexams(ExamID);


-- Table to store system user information, with a unique UserID for each user, and a role column to differentiate between admins and coordinators (NON EXAM USER)
CREATE TABLE "UserT" (
    "UserID"    INTEGER,
    "Username"  VARCHAR(50) NOT NULL,
    "Password"  VARCHAR(255) NOT NULL,
    "Email"     VARCHAR(255) UNIQUE NOT NULL,
    "Role"      VARCHAR(20) NOT NULL DEFAULT 'Learner',
    "DefaultAdmin"    BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY("UserID" AUTOINCREMENT)
    CHECK (Role IN ('Admin','Faculty','Learner'))
);

-- View to determine the curernt state of the exam sessions - past and present
-- used for the dashboard
CREATE VIEW examMetrics AS
SELECT c.CourseCode, c.Description, o.Password, 
       o.ExamID, o.Year, o.Semester,
       --SUBSTRING(o.ExamID,1,4) AS "Year",
       --SUBSTRING(o.ExamID,5,2) AS "Semester",       
	   COUNT(CASE l.status WHEN 'ready' THEN 1 END) AS Ready,
	   COUNT(CASE l.status WHEN 'active' THEN 1 END) AS Active,	   
	   COUNT(CASE l.status WHEN 'expired' THEN 1 END) AS Expired,
	   COUNT(CASE l.status WHEN 'closed' THEN 1 END) AS Closed		   
FROM courses c, offerings o, Learnerexams l
WHERE c.CourseCode = o.CourseCode 
	  AND o.ExamID = l.ExamID  
GROUP BY c.CourseCode, o.year
ORDER by o.year DESC;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop tables in reverse order of creation to avoid foreign key constraint violations
DROP VIEW IF EXISTS "examMetrics";
DROP TABLE IF EXISTS "Learnerexams";
DROP TABLE IF EXISTS "Learners";
DROP TABLE IF EXISTS "Offerings";
DROP TABLE IF EXISTS "Courses";
DROP TABLE IF EXISTS "UserT";
-- +goose StatementEnd