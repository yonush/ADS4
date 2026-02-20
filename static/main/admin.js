// admin.js
// Fetch users from the server
fetch("/api/user")
    .then((response) => response.json())
    .then((users) => {
        // Convert current_user_id to a number
        const currentUserIdNumber = parseInt(current_user_id, 10);

        // Sort the users array to put the current user first
        users.sort((a, b) => {
            if (a.user_id === currentUserIdNumber) return -1;
            if (b.user_id === currentUserIdNumber) return 1;
            return a.username.localeCompare(b.username); // Sort others alphabetically
        });

        // Create a table row for each user
        const userRows = users.map((user) => {
            // Convert user.default_admin to a boolean
            var isAdmin = JSON.parse(user.default_admin);
            // Convert is_current_user_default_admin to a boolean
            var current_default_admin = JSON.parse(
                is_current_user_default_admin
            );

            const hideDelete = current_default_admin && isAdmin;

            // Determine whether to hide action buttons based on conditions
            const hideActions = !current_default_admin && isAdmin;

            // Generate the row HTML
            return `
<tr${user.user_id === currentUserIdNumber ? ' class="table-primary"' : ""}>
    <td data-label="Username">${user.username}</td>
    <td data-label="Email">${user.email}</td>
    <td data-label="Role">${user.role}</td>
    <td>
    <div class="btn-group">
    ${
        hideActions
            ? "<span class='text-muted'>No actions available</span>"
            : `
            <button class="btn btn-warning p-2 edit-user-button" onclick="editUser(${
                user.user_id
            })"
                    title="Edit User">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/>
                    <path d="m15 5 4 4"/>
                </svg>
            </button>
            ${
                hideDelete
                    ? ""
                    : `<button class="btn btn-danger p-2 delete-button" 
                            onclick="showDeleteModal(${user.user_id}, 'user', '${user.username}', '${currentUserIdNumber}')" 
                            data-id="${user.user_id}" 
                            title="Delete User">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M3 6h18"/>
                            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                            <line x1="10" y1="11" x2="10" y2="17"/>
                            <line x1="14" y1="11" x2="14" y2="17"/>
                        </svg>
                    </button>`
            }
            `
    }
</div>
    </td>
</tr>
`;
        });

        // Add the rows to the users table
        $("#users-table tbody").html(userRows.join(""));
    });

export async function editUser(userId) {
    const id = userId;
    console.log("Edit user ID:", id);
    // Handle edit
    // Fetch the user data from the nearest row
    const row = $(event.target).closest("tr");
    const username = row.find("td[data-label=Username]").text();
    const email = row.find("td[data-label=Email]").text();
    const role = row.find("td[data-label=Role]").text();

    const default_admin = await fetch(`/api/user/${username}`)
        .then((response) => response.json())
        .then((user) => {
            return user.default_admin.toString();
        });

    // Fill in the form with the user data
    $("#editUserForm")[0].reset();
    $("#editUserForm input[name=current_user_id]").val(current_user_id);
    $("#editUserForm input[name=user_id]").val(id);
    $("#editUserForm input[name=username]").val(username);
    $("#editUserForm input[name=email]").val(email);
    $("#editUserForm select[name=role]").val(role);
    $("#editUserForm input[name=default_admin]").val(default_admin);

    // Set the form action to the update endpoint for this user
    $("#editUserForm").attr("action", `/api/user/${id}`);

    // Get the user ID of the user being updated
    const updatedUserId = $("#editUserForm input[name=user_id]").val();

    // If the current user ID is equal to the user being updated, display the password field
    if (current_user_id === updatedUserId) {
        $("#passwordField").show();
    } else {
        $("#passwordField").hide();
    }

    // Show the modal
    $("#editUserModal").modal("show");
    // Clear validation classes
    $("#editUserForm").removeClass("was-validated");

    var editUserForm = document.getElementById("editUserForm");

    // Add event listener to the submit button
    $("#editUserBtn").click(function (event) {
        // Check if the form is valid
        if (!editUserForm.checkValidity()) {
            event.stopPropagation();
            editUserForm.classList.add("was-validated");
        } else {
            // If the form is valid, prepare to send the PUT request
            const formData = new FormData(editUserForm);
            const jsonData = {};
            for (const [key, value] of formData.entries()) {
                jsonData[key] = value;
            }
            fetch(`/api/user/${document.getElementById("editUserID").value}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(jsonData),
            })
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        window.location.href = data.redirectURL;
                    } else if (data.message) {
                        window.location.href = data.redirectURL;
                    } else {
                        console.error("Unexpected response:", data);
                        // Handle unexpected responses (e.g., show an error message)
                        throw new Error("Unexpected response");
                    }
                })
                .catch((error) => {
                    console.error("Fetch error:", error);
                    // Optionally display a user-friendly error message
                });
        }
    });
}

export function AddUser() {
    // Clear the form before showing it
    document.getElementById("addUserForm").reset();
    document
        .getElementById("addUserForm")
        .classList.remove("was-validated");


}

(function () {
    "use strict";

    // Fetch the form and the submit button
    var form = document.querySelector("#addUserForm");;
    var submitButton = document.querySelector("#addUserBtn");

    // Add event listener to the submit button
    // Add event listener to the submit button
    submitButton.addEventListener(
        "click",
        function (event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            } else {
                
                // If the form is valid, prepare to send the POST request
                const formData = new FormData(form);
                const jsonData = {};
                for (const [key, value] of formData.entries()) {
                    jsonData[key] = value;
                }
                console.log("Sending",JSON.stringify(jsonData));
                fetch(`/api/user`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(jsonData),
                })
                    .then((response) => response.json())
                    .then((data) => {
                        if (data.error) {
                            window.location.href = data.redirectURL;
                        } else if (data.message) {
                            window.location.href = data.redirectURL;
                        } else {
                            console.error("Unexpected response:", data);
                            // Handle unexpected responses (e.g., show an error message)
                            throw new Error("Unexpected response");
                        }
                    })
                    .catch((error) => {
                        console.error("Fetch error:", error);
                        // Optionally display a user-friendly error message
                    });     
            }

            form.classList.add("was-validated");
        },
        false
    );
})();


// Make functions available globally
window.editUser = editUser;
window.AddUser = AddUser;

