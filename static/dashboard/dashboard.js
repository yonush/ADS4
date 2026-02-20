// dashboard.js

//-----------------------------------------------------------------------------------------------------------
//filtering controls
export function clearFilters() {
    // Reset each filter dropdown to its first option (except site filter)
    document.getElementById("yearFilter").selectedIndex = 0;
    document.getElementById("semesterFilter").selectedIndex = 0;
    document.getElementById("searchInput").value = ""; // Clear search input

    // Reset active filters
    activeFilters = {
        year: null,
        semester: null,
        status: null,
    };

    filteredOfferings = [...allOfferings]; // Reset to original devices
    clearTableBody();
    loadOfferingsAndUpdateTable();
}

function getFilterOptions() {
    //TODO change this to retrieve the years from the database
    //SELECT distinct year FROM examMetrics
    fetchAndPopulateSelect(
        "/yearlist",
        "yearFilter",
        "year",
        "year",
        "Current Year"
    );
    setupSemesterFilter();
    setupYearFilter();
}

function fetchAndPopulateSelect(
    url,
    selectId,
    textKey,
    valueKey,
    defaultOptionText
) {
    fetch(url)
        .then((response) => response.json())
        .then((data) => {
            const select = document.getElementById(selectId);
            select.innerHTML = "";

            // Add the default option
            addDefaultOption(select, defaultOptionText);

            // Check if data is valid and is an array
            if (data && Array.isArray(data) && data.length > 0) {
                data.forEach((item) => {
                    const option = document.createElement("option");
                    option.text = item[textKey];
                    if (valueKey) option.value = item[valueKey];
                    select.add(option);
                });
            } else {
                console.log(`No data available for ${selectId}`);
            }
        })
        .catch((error) => {
            console.error(`Error fetching ${selectId} data:`, error);
        });
}

function addDefaultOption(select, text) {
    const defaultOption = document.createElement("option");
    defaultOption.text = text;
    defaultOption.selected = true;
    select.add(defaultOption);
}

function setupSemesterFilter() {
    document.getElementById("semesterFilter").addEventListener("change", () => {
        filterBySemester();
    });
}


//-----------------------------------------------------------------------------------------------------------
//filtering actions on the table
//TODO fix this filter
function filterBySemester(semester) {
    const semesterFilter = document.getElementById("semesterFilter");

    if (semester) {
        // Loop through `semesterFilter` options to select the one with matching text
        for (const option of semesterFilter.options) {
            if (option.text === semester) {
                option.selected = true;
                break;
            }
        }
    } else {
        semester = semesterFilter.selectedOptions[0].text;
    }

    // Fetch devices based on `buildingCode` and `siteId`
    if (semester === "All Semesters" || semesterFilter.value === "") {
        loadOfferingsAndUpdateTable("");
    } else {
        loadOfferingsAndUpdateTable(semester);
    }
}

//TODO need code to acquire the stored years
//based on SELECT distinct year FROM examMetrics
function setupYearFilter() {
    document.getElementById("yearFilter").addEventListener("change", () => {
        filterByYear();
    });
}

function filterByYear() {
    const selectedYear = document.getElementById("yearFilter").value;

    if (selectedYear != "All Years") {        
        //TODO add route for years
        if (selectedYear == 'Current Year') {
            selectedYear = new Date().getFullYear()  // returns the current year
        }
        fetch(`/api/offerings?year=${selectedYear}`)
            .then((response) => response.json())
            .then((data) => {
                const yearSelect = document.getElementById("yearFilter");
                yearSelect.innerHTML = "";

                // Add default "Curernt Year" option
                addDefaultOption(yearSelect, "Current Year");

                // Add room options
                data.forEach((room) => {
                    const option = document.createElement("option");
                    option.value = year.year // Store the ID as the value
                    option.text = room.year; // Show the code as the text
                    yearSelect.add(option);
                });
            });
        return;
    }
}

function clearTableBody() {
    const tableBody = document.getElementById("exam-offering-body");
    if (tableBody) {
        tableBody.innerHTML = "";
    } else {
        console.error("Table body element not found");
    }
}

// Initialize the map and populate filter options
getFilterOptions();

let currentPage = 1;
let rowsPerPage = 10;
let allOfferings = [];
let filteredOfferings = [];

// Keep track of active filters
let activeFilters = {
    year: null,
    semester: null,
};

// Add event listeners for the new filters
document.getElementById("yearFilter").addEventListener("change", () => {
    filterTableByYear();
    clearTableBody();
    updateTable();
});

document.getElementById("semesterFilter").addEventListener("change", () => {
    filterTableBySemester();
    clearTableBody();
    updateTable();
});




async function getAllExamOfferings(semester = "") {
    try {
        let url = "/offerings";
        const params = new URLSearchParams();
        if (semester) params.append("semester", semester);
        if (params.toString()) url += `?${params.toString()}`;

        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const offerings = await response.json();
        return offerings; // Return the devices instead of storing in global variable
    } catch (err) {
        console.error("Failed to fetch offerings:", err);
        return []; // Return empty array in case of error
    }
}

async function loadOfferingsAndUpdateTable(semester="") {
    const offerings = await getAllExamOfferings(semester);
    allOfferings = offerings; // Update global variable if needed
    filteredOfferings = offerings; // Initialize filtered devices

    updateTable();

    if (offerings.length === 0) {
        const tbody = document.getElementById("exam-offering-body");
        if (tbody) {
            tbody.innerHTML = `<tr><td colspan="12" class="text-center">No devices found.</td></tr>`;
        }
    }
}

function updateTable() {
    const tbody = document.getElementById("exam-offering-body");
    if (!tbody) {
        console.error("Table body element not found");
        return;
    }

    const startIndex = (currentPage - 1) * rowsPerPage;
    const endIndex = startIndex + rowsPerPage;
    const pageOfferings = filteredOfferings.slice(startIndex, endIndex);

    // Clear table if no devices
    if (!Array.isArray(pageOfferings) || pageOfferings.length === 0) {
        tbody.innerHTML = `<tr><td colspan="12" class="text-center">No devices found.</td></tr>`;
    } else {
        tbody.innerHTML = pageOfferings.map(formatOfferingRow).join("");
    }

    updatePaginationControls(); 
}

// Filter functions for each criteria
function filterTableByYear() {
    const yearSelect = document.getElementById("yearFilter");
    const selectedYear = yearSelect.value;
    const selectedYearText = yearSelect.selectedOptions[0].text;

    if (selectedYear === "Current Year") {
        filteredOfferings = [...allOfferings];
    } else {
        // Try matching against both the value and text of the selected room
        filteredOfferings = allOfferings.filter(
            //TODO: fix the year filter
            (device) =>
                device.room_code === selectedYearText ||
                device.room_id === selectedYear
        );
    }
    activeFilters.year = selectedYearText;
    applyFilters();
}

function filterTableBySemester() {
    const selectedSemester =
        document.getElementById("semesterFilter").selectedOptions[0].text;

    if (selectedSemester === "All Semesters") {
        filteredOfferings = [...allOfferings];
    } else {
        filteredOfferings = allOfferings.filter(
            //TODO: fix the semester filter
            (device) => device.emergency_device_type_name === selectedSemester
        );
    }
    activeFilters.semester = selectedSemester;
    applyFilters();
}


// Apply all active filters
function applyFilters() {
    // Start with all offerings
    filteredOfferings = [...allOfferings];

    // Apply year filter if active
    if (activeFilters.year && 
        activeFilters.year !== "Current year`" && 
        activeFilters.year !== "All Years`"
    ) {
        filteredOfferings = filteredOfferings.filter(
            //TODO fix the filter application
            (device) => device.room_code === activeFilters.year
        );
    }

    // Apply semester filter if active
    if (
        activeFilters.semester &&
        activeFilters.semester !== "Semester" &&
        activeFilters.semester !== "All Semesters"
    ) {
        filteredOfferings = filteredOfferings.filter(
            //TODO fix the filter application
            (device) =>
                device.emergency_device_type_name === activeFilters.semester
        );
    }
    updateTable();
}



// Initial fetch without filtering
loadOfferingsAndUpdateTable();
/*
<tr>
    <th>Course Code ▲</th>
    <th>Description ▲</th>
    <th>Exam ID ▲</th>
    <th>Password ▲</th>
    <th>Ready ▲</th>
    <th>Active ▲</th>
    <th>Expired ▲</th>
    <th>Closed ▲</th>
    <th>Actions</th>
</tr>
*/
//TODO update to show the offerings
function formatOfferingRow(offering) {
    if (!device) return "";
    const formatDateMonthYear = (dateString) =>
        formatDate(dateString, { month: "short", year: "numeric" });
    const formatDateFull = (dateString) =>
        formatDate(dateString, {
            year: "numeric",
            month: "short",
            day: "numeric",
            timeZone: "Pacific/Auckland",
        });

    const buttons = getActionButtons(offering);

    // Declare isAdmin within the function
    let isAdmin = false;

    // Ensure role is defined and check for "Admin"
    if (role === "Admin") {
        isAdmin = true;
    }

    return `
        <tr>
            <td data-label="Device Type">${
                device.emergency_device_type_name
            }</td>
            <td data-label="Extinguisher Type">${
                device.extinguisher_type_name.String
            }</td>
            <td data-label="Building">${device.building_code}</td>
            <td data-label="Room">${device.room_code}</td>
            <td data-label="Serial Number">${device.serial_number.String}</td>
            <td data-label="Manufacture Date">${formatDateMonthYear(
                device.manufacture_date.Time
            )}</td>
            <td data-label="Expire Date">${formatDateMonthYear(
                device.expire_date.Time
            )}</td>
            ${
                isAdmin
                    ? `<td data-label="Last Inspection Date">${formatDateFull(
                          device.last_inspection_datetime.Time
                      )}</td>`
                    : ""
            }
            ${
                isAdmin
                    ? `<td data-label="Next Inspection Date">${formatDateFull(
                          device.next_inspection_date.Time
                      )}</td>`
                    : ""
            }
            <td data-label="Size">${device.size.String}</td>

            <td>
                <div class="btn-group">
                    ${buttons}
                </div>
            </td>
        </tr>
    `;
}

function formatDate(dateString, options) {
    if (!dateString || dateString === "0001-01-01T00:00:00Z") {
        return "N/A";
    }
    return new Date(dateString).toLocaleString("en-NZ", {
        timeZone: "Pacific/Auckland", // Ensure the correct timezone
        ...options,
    });
}

export function getActionButtons(offering) {
    let buttons = `
        <button class="btn btn-primary p-2" 
                onclick="" 
                title="View Notes">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M8 2v4"/>
                <path d="M12 2v4"/>
                <path d="M16 2v4"/>
                <rect width="16" height="18" x="4" y="4" rx="2"/>
                <path d="M8 10h6"/>
                <path d="M8 14h8"/>
                <path d="M8 18h5"/>
            </svg>
        </button>`;

    if (role === "Admin") {
        buttons += `
            <button class="btn btn-warning p-2 ml-2" 
                    onclick=""
                    title="Edit Offering">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/>
                    <path d="m15 5 4 4"/>
                </svg>
            </button>
            <button class="btn btn-danger p-2 ml-2" 
                    onclick=""
                    title="Delete Offering">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M3 6h18"/>
                    <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                    <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                    <line x1="10" y1="11" x2="10" y2="17"/>
                    <line x1="14" y1="11" x2="14" y2="17"/>
                </svg>
            </button>
        `;
    }
    return buttons;
}



document.getElementById("searchInput").addEventListener("input", () => {
    searchOfferings();
});

// Updated search function to use combined filtering approach
//TODO update the search to look for offerings
export async function searchOfferings() {
    const searchInput = document.getElementById("searchInput");
    const searchValue = searchInput.value.toLowerCase();

    // First, reapply base filters to get fresh filtered state
    // Start with all devices
    let basefilteredOfferings = [...allOfferings];

    // Apply active filters to get our base filtered state
    if (activeFilters.room && activeFilters.room !== "All Rooms") {
        basefilteredOfferings = basefilteredOfferings.filter(
            (device) => device.room_code === activeFilters.room
        );
    }

    if (
        activeFilters.deviceType &&
        activeFilters.deviceType !== "Device Type" &&
        activeFilters.deviceType !== "All Device Types"
    ) {
        basefilteredOfferings = basefilteredOfferings.filter(
            (device) =>
                device.emergency_device_type_name === activeFilters.deviceType
        );
    }

    if (
        activeFilters.status &&
        activeFilters.status !== "Status" &&
        activeFilters.status !== "All Statuses"
    ) {
        basefilteredOfferings = basefilteredOfferings.filter(
            (device) => device.status.String === activeFilters.status
        );
    }

    // If search is empty, use just the filtered results
    if (!searchValue) {
        filteredOfferings = basefilteredOfferings;
    } else {
        // Apply search filter to the fresh filtered state
        filteredOfferings = basefilteredOfferings.filter((device) => {
            const baseSearch =
                device.emergency_device_type_name
                    .toLowerCase()
                    .includes(searchValue) ||
                device.extinguisher_type_name.String.toLowerCase().includes(
                    searchValue
                ) ||
                device.room_code.toLowerCase().includes(searchValue) ||
                device.serial_number.String.toLowerCase().includes(
                    searchValue
                ) ||
                device.manufacture_date.Time.toLowerCase().includes(
                    searchValue
                ) ||
                device.expire_date.Time.toLowerCase().includes(searchValue) ||
                device.size.String.toLowerCase().includes(searchValue) ||
                device.status.String.toLowerCase().includes(searchValue) ||
                device.description.String.toLowerCase().includes(searchValue);

            // Add admin-only fields if user is admin
            if (role === "Admin") {
                const lastInspectionFormatted = new Date(
                    device.last_inspection_datetime.Time
                )
                    .toLocaleDateString("en-NZ", {
                        day: "numeric",
                        month: "long",
                        year: "numeric",
                    })
                    .toLowerCase();

                const nextInspectionFormatted = new Date(
                    device.next_inspection_date.Time
                )
                    .toLocaleDateString("en-NZ", {
                        day: "numeric",
                        month: "long",
                        year: "numeric",
                    })
                    .toLowerCase();

                return (
                    baseSearch ||
                    lastInspectionFormatted.includes(searchValue) ||
                    nextInspectionFormatted.includes(searchValue)
                );
            }

            return baseSearch;
        });
    }

    updateTable();
}

//-----------------------------------------------------------------------------------------------------------
//pagination functions

export function updatePaginationControls() {
    const totalPages = Math.ceil(allOfferings.length / rowsPerPage);
    const paginationEl = document.querySelector(".pagination");
    const isMobile = window.innerWidth < 768; // Detect mobile devices

    let paginationHTML = `
        <li class="page-item ${currentPage === 1 ? "disabled" : ""}">
            <a class="page-link" href="#" data-page="${
                currentPage - 1
            }" aria-label="Previous">
                <span aria-hidden="true">&laquo;</span>
            </a>
        </li>
    `;

    function addPageNumber(pageNum) {
        paginationHTML += `
            <li class="page-item ${
                currentPage === pageNum ? "active" : ""
            }" aria-current="page">
                <a class="page-link" href="#" data-page="${pageNum}">${pageNum}</a>
            </li>
        `;
    }

    function addEllipsis() {
        paginationHTML += `
            <li class="page-item disabled">
                <span class="page-link">...</span>
            </li>
        `;
    }

    if (isMobile) {
        // Simplified pagination for mobile
        if (totalPages <= 3) {
            for (let i = 1; i <= totalPages; i++) {
                addPageNumber(i);
            }
        } else {
            addPageNumber(1);
            if (currentPage !== 1 && currentPage !== totalPages) {
                addPageNumber(currentPage);
            }
            addPageNumber(totalPages);
        }
    } else {
        if (totalPages <= 7) {
            for (let i = 1; i <= totalPages; i++) {
                addPageNumber(i);
            }
        } else {
            addPageNumber(1);
            if (currentPage > 3) addEllipsis();

            let start = Math.max(2, currentPage - 1);
            let end = Math.min(totalPages - 1, currentPage + 1);

            if (currentPage <= 3) {
                end = 4;
            } else if (currentPage >= totalPages - 2) {
                start = totalPages - 3;
            }

            for (let i = start; i <= end; i++) {
                addPageNumber(i);
            }

            if (currentPage < totalPages - 2) addEllipsis();
            addPageNumber(totalPages);
        }
    }

    paginationHTML += `
        <li class="page-item ${currentPage === totalPages ? "disabled" : ""}">
            <a class="page-link" href="#" data-page="${
                currentPage + 1
            }" aria-label="Next">
                <span aria-hidden="true">&raquo;</span>
            </a>
        </li>
    `;

    paginationEl.innerHTML = paginationHTML;

    function handlePaginationClick(e) {
        e.preventDefault();
        e.stopPropagation();

        let target = e.target.closest(".page-link");

        if (target && target.hasAttribute("data-page")) {
            const newPage = parseInt(target.getAttribute("data-page"), 10);

            if (
                newPage !== currentPage &&
                newPage > 0 &&
                newPage <= totalPages
            ) {
                currentPage = newPage;
                updateTable();
            }
        }
    }

    // Remove existing event listeners
    paginationEl.removeEventListener("click", handlePaginationClick);
    paginationEl.removeEventListener("touchstart", handlePaginationClick);

    // Add event listeners to the pagination container
    paginationEl.addEventListener("click", handlePaginationClick);
    paginationEl.addEventListener("touchstart", handlePaginationClick);
}

// Event listener for rows per page dropdown
document.getElementById("rowsPerPage").addEventListener("change", (e) => {
    rowsPerPage = parseInt(e.target.value);
    currentPage = 1; // Reset to first page when changing rows per page
    updateTable();
});

// Make functions available globally
window.clearFilters = clearFilters;

