// file upload script - adpated from https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequestUpload

const fileInput = document.getElementById("datafile");
const abortButton = document.getElementById("abort");

const overwrite = document.getElementById("overwrite");
const purge = document.getElementById("purge");

const output = document.querySelector("output");
output.style.display = "none" 

export function btnImport() {
  //fileInput.addEventListener("change", () => {
  output.style.display = "none"

  // Link abort button
  const controller = new AbortController();
  abortButton.addEventListener(
    "click",
    () => {
      //xhr.abort();
      controller.abort();
    },
    { once: true },
  );

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

  fetch("/import/"+route, {
    signal: controller.signal,
    method: 'POST',
    body: fileData
    //body: JSON.stringify(jsonData) //for JSON data
    
  })
  //.then((response) => response.json())
  .then((response) => response.text())

  .then((result) => {
    console.log('Message:', result);  
    abortButton.disabled = true;

    output.style.display = "block"
    output.textContent = result
  })

  .catch((error) => {
    console.error('Error:', error);
    abortButton.disabled = true;

    output.style.display = "block"
    output.textContent = error
  }); 
   
};

window.btnImport = btnImport;

