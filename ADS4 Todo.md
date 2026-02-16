# ADS4 Todo & summary of changes

- Consider upgrading from Echo v4 to Echo v5.
- Changed to SQlite for portablility but include support for PosgreSQL for scalability
- Consider using templ for embedding the templates 
- Include CORS and CSRF protections
  
- change main page to include basic metrics without logging into the system

- https://www.sqlite.org/pragma.html#pragma_synchronous
- https://github.com/mattn/go-sqlite3?tab=readme-ov-file

**LOC counter**
Get-ChildItem -recurse *.go |Get-Content | Measure-Object -line

## TODO

- add data importers under the admin area for the courses, learners, offerings and learner exams
  - start with the learners, easier with less fields and validation

- use for dateime handling https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go
- update the /dashboard page content  

- consider timestamping the actual questions in the assessment tool
  
### dashboard metrics
Exam progression ready->active->expired|closed->marked

Course code, description, examid, yr, semester, ready, active, expired, closed, marked

## ADS4/internal/app

Basic database Model


``` mermaid
  flowchart LR
    main-->app-->routes-->handlers
```

``` mermaid
block
	block:grp0:2
		columns 1
		main.go 
	end	
	block:grp1:3
		columns 1
		id1("app.go") 
		id2[("database.go")] 
		id3("config.go")
	end	
	block:grp2:2
		columns 1
		id4("routes.go")
		id5("*_handlers.go")
	end	
	
	

```

app/app.go - udpate to include certificate handling. Include domain in the config.

app/auth_handler.go - fix the password email sender and SMTP credentials
app/routes.go - update these to reflect the new routes

## ADS4/internal/database

- add DB tests

## ADS4/templates

- Update the templates to reflect the new views
- Keep the admin/main, admin/authentication and admin/dashboard
- Update admin/dashboard with the new metrics display. Remove refernces to device management.
- remove admin/inspection

**templates/admin**
- update admin/ with new views
- remove admin/building_management, admin/device_type_management, admin/room_management, admin/site_management
- keep admin/user_management
- update admin/admin.html with toast and templates
- keep admin/admin_navbar.html

## ADS4/static
**ADS4/static/admin**
- Update admin.css
- Update admin.js and remove old handler JS. Replace with new handler functionality

**ADS4/static/assets**
- add new user guide
- update app logo

**ADS4/static/authentication**
- keep this folder
  
**ADS4/static/dashboard**
- Update JS with new functionality

**ADS4/static/main**
- udpate notifications.js with new notification handlers
- remove inspections.js
- keep main.js

**ADS4/static/site_maps**
- remove this folder and contents


