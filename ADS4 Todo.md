# ADS4 Todo & summary of changes

- Upgrade from Echo v4 to Echo v5. v4 safe until 31 Dec 2026 
  - https://github.com/labstack/echo/blob/master/API_CHANGES_V5.md
- Use templ and go fs for embedding the templates 
- Include CORS and CSRF protections - needs domain name
  
- https://www.sqlite.org/pragma.html#pragma_synchronous
- https://github.com/mattn/go-sqlite3?tab=readme-ov-file
- https://go-sponge.com/
- https://docs.gofiber.io/template/html_v2.x.x/html/TEMPLATES_CHEATSHEET/
- https://gowebly.org/

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
  - [option]generate exception report 
- use for datetime handling https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go
- [option] Timestamp the actual questions in the assessment tool
  
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
    main-->app-->routes-->handlers-->database
```

app/app.go - udpate to include certificate handling. Include domain in the config.

app/auth_handler.go - fix the password email sender and SMTP credentials
app/routes.go - update these to reflect the new routes

## ADS4/templates

- [paused] Add the javaScript for the CRUD

