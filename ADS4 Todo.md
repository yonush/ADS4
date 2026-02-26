# ADS4 Todo & summary of changes

- Upgrade from Echo v4 to Echo v5. v4 safe until 31 Dec 2026 
  - https://github.com/labstack/echo/blob/master/API_CHANGES_V5.md
- Use templ and go fs for embedding the templates 
- Include CORS and CSRF protections
  
- replace the XMLHTTP API with the more modern fetch - refer to static/main/xhrupload.js
  - https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch

- https://www.sqlite.org/pragma.html#pragma_synchronous
- https://github.com/mattn/go-sqlite3?tab=readme-ov-file
- https://go-sponge.com/
- https://docs.gofiber.io/template/html_v2.x.x/html/TEMPLATES_CHEATSHEET/
- https://gowebly.org/

- Update filter based on the semster period. In other words the dashboard only shows the current exams based on S1-S3 and year
- admin.html -> fix/remove the following as this is a security risk
          <script>
            var role = "{{.role}}";
            var current_user_id = "{{.user_id}}";
            var user_id = "{{.user_id}}";
            var is_current_user_default_admin = "{{.default_admin}}";
        </script>

- Update importers to include:
  - check if it is the correct file - contains header with correct fields and field count
  - check field types and structure
  - generate calculated data items - e.g examid
  - generate exception report
- create general search query
  - groupby/filterby by field (examid, studentid, course) then value. e.g route /groupby/:field/:value
  
- use for datetime handling https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go

- Timestamp the actual questions in the assessment tool
  
- Template with templ - https://github.com/a-h/templ
- HTMX for SSE with go-HTMX & HTMX - https://github.com/donseba/go-htmx, https://htmx.org/   
- AI text reviewer with LangChainGo - https://github.com/tmc/langchaingo
  - https://tmc.github.io/langchaingo/docs/tutorials/code-reviewer
  - https://eli.thegreenplace.net/2023/using-ollama-with-langchaingo/
  - https://eli.thegreenplace.net/2024/gemma-ollama-and-langchaingo/
  - https://github.com/eliben/code-for-blog/tree/main/2023/ollama-go-langchain
    - C:\devwork\code-for-blog\2023

### other

**LOC counter**
Get-ChildItem -recurse *.go |Get-Content | Measure-Object -line

### dashboard metrics

Exam progression ready->active->expired|closed->marked

Course code, description, examid, yr, semester, ready, active, expired, closed, marked

### timed events
https://developer.mozilla.org/en-US/docs/Web/API/Window/setInterval 
<script>
    function autoRefresh() {
        window.location = window.location.href;
    }
    setInterval('autoRefresh()', 5000);
</script>

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

- Add the CRUD templates for the new handlers

## ADS4/static
**ADS4/static/admin**
- Update admin.css
- Update admin.js and remove old handler JS. Replace with new handler functionality

**ADS4/static/dashboard**
- Update JS with new functionality

**ADS4/static/main**
- udpate notifications.js with new notification handlers
- remove inspections.js
- keep main.js
