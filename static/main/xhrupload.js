// file upload script - adpated from https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequestUpload

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
