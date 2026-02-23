// file upload script - adpated from https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequestUpload
//replce with fetch api https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API
/* TODO: replace the XMLHttpRequest with the fetch api
//syntax with authentication
fetch("url", {
    headers: {
      "Content-Type": "application/json",
      'Authorization': 'Basic ' + btoa(globalUsername + ":" + globalPassword),
    },
    method: "POST",
    body: moveBody
  })
  .then(response => console.log(response.status) || response) // output the status and return response
  .then(response => response.text()) // send response body to next then chain
  .then(body => console.log(body)) // you can use response body here

// Source - https://stackoverflow.com/a/36082038
// Posted by Damien, modified by community. See post 'Timeline' for change history
// Retrieved 2026-02-23, License - CC BY-SA 4.0

// Select your input type file and store it in a variable
const input = document.getElementById('fileinput');

// This will upload the file after having read it
const upload = (file) => {
  fetch('http://www.example.net', { // Your POST endpoint
    method: 'POST',
    headers: {
      // Content-Type may need to be completely **omitted**
      // or you may need something
      "Content-Type": "You will perhaps need to define a content-type here"
    },
    body: file // This is your file object
  }).then(
    response => response.json() // if the response is a JSON object
  ).then(
    success => console.log(success) // Handle the success response object
  ).catch(
    error => console.log(error) // Handle the error response object
  );
};

// Event handler executed when a file is selected
const onSelectFile = () => upload(input.files[0]);

// Add a listener on your input
// It will be triggered when a file will be selected
input.addEventListener('change', onSelectFile, false);

//----------------with form data
const formData = new FormData([html_form]);
//for JSON data
//const jsonData = {};
//for (const [key, value] of formData.entries()) {jsonData[key] = value;}
const fileField = document.querySelector('input[type="file"]');
// for multiple files
//for (const file of fileField.files) {formData.append('files',file,file.name)}

formData.append('username', 'abc123');
formData.append('avatar', fileField.files[0]);

fetch('https://example.com/profile/avatar', {
  method: 'PUT',
  body: formData
  //body: JSON.stringify(jsonData) //for JSON data
})
.then((response) => response.json())
.then((result) => {
  console.log('Success:', result);  
})
.catch((error) => {
  console.error('Error:', error);
});

//---------------------
const upload = (file) => {
    console.log(file);

    

    fetch('http://localhost:8080/files/uploadFile', { 
    method: 'POST',
    // headers: {
    //   //"Content-Disposition": "attachment; name='file'; filename='xml2.txt'",
    //   "Content-Type": "multipart/form-data; boundary=BbC04y " //"multipart/mixed;boundary=gc0p4Jq0M2Yt08jU534c0p" //  ή // multipart/form-data 
    // },
    body: file // This is your file object
  }).then(
    response => response.json() // if the response is a JSON object
  ).then(
    success => console.log(success) // Handle the success response object
  ).catch(
    error => console.log(error) // Handle the error response object
  );

  //cvForm.submit();
};
const onSelectFile = () => upload(uploadCvInput.files[0]);
uploadCvInput.addEventListener('change', onSelectFile, false);

<form id="cv_form" style="display: none;"
  enctype="multipart/form-data">
  <input id="uploadCV" type="file" name="file"/>
  <button type="submit" id="upload_btn">upload</button>
</form>
<ul class="dropdown-menu">
<li class="nav-item"><a class="nav-link" href="#" id="upload">UPLOAD CV</a></li>
<li class="nav-item"><a class="nav-link" href="#" id="download">DOWNLOAD CV</a></li>
</ul>
*/
const fileInput = document.getElementById("datafile");
const progressBar = document.querySelector("progress");
const log = document.querySelector("output");
const abortButton = document.getElementById("abort");
const overwrite = document.getElementById("overwrite");
const purge = document.getElementById("purge");

export function btnImport() {
//fileInput.addEventListener("change", () => {
  const xhr = new XMLHttpRequest();
  xhr.timeout = 2000; // 2 seconds

  // Link abort button
  abortButton.addEventListener(
    "click",
    () => {
      xhr.abort();
    },
    { once: true },
  );

  // When the upload starts, we display the progress bar
  xhr.upload.addEventListener("loadstart", (event) => {
    progressBar.classList.add("visible");
    progressBar.value = 0;
    progressBar.max = event.total;
    log.textContent = "Uploading (0%)…";
    abortButton.disabled = false;
  });

  // Each time a progress event is received, we update the bar
  xhr.upload.addEventListener("progress", (event) => {
    progressBar.value = event.loaded;
    log.textContent = `Uploading (${(
      (event.loaded / event.total) *
      100
    ).toFixed(2)}%)…`;
  });

  // When the upload is finished, we hide the progress bar.
  xhr.upload.addEventListener("loadend", (event) => {
    progressBar.classList.remove("visible");
    if (event.loaded !== 0) {
      log.textContent = "Upload finished.";
    }
    abortButton.disabled = true;
  });

  // In case of an error, an abort, or a timeout, we hide the progress bar
  function errorAction(event) {
    progressBar.classList.remove("visible");
    log.textContent = `Upload failed: ${event.type}`;
  }
  xhr.upload.addEventListener("error", errorAction);
  xhr.upload.addEventListener("abort", errorAction);
  xhr.upload.addEventListener("timeout", errorAction);
  xhr.addEventListener("error", errorAction);

  // Build the payload
  const fileData = new FormData();
  fileData.append("datafile", fileInput.files[0]);
  fileData.append("purge", purge.checked);
  fileData.append("overwrite", overwrite.checked);

  var route = ""
  var routes = document.getElementsByName("route"); 
  for (let i = 0; i < routes.length; i++) {
      if (routes[i].type == "radio" && routes[i].checked ){
          route = routes[i].value;
          break
      }
  }

  // Theoretically, event listeners could be set after the open() call
  switch (route) {
    case 'course':
        xhr.open("POST", "/importcourses", true);
        break;
    case 'learner':
        xhr.open("POST", "/importlearners", true);
        break;
    case 'learnerexam':
        xhr.open("POST", "/importlearnerexams", true);
        break;
    case 'offering':
        xhr.open("POST", "/importofferings", true);
        break; 

  }

  xhr.send(fileData);
};

window.btnImport = btnImport;
