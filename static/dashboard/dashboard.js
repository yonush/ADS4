// dashboard.js

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
    status: null,
};

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
    loadOfferingsAndUpdateTable("","");
}

function getFilterOptions() {
    //TODO change this to retrieve the years from the database
    //SELECT distinct year FROM examMetrics
    fetchAndPopulateSelect(
        "/yearlist",
        "yearFilter",
        "Year",
        "Year",
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
        // Loop through `buildingFilter` options to select the one with matching text
        for (const option of semesterFilter.options) {
            if (option.text === semester) {
                option.selected = true;
                break;
            }
        }
    } else {
        semester = semesterFilter.selectedOptions[0].text;
        year = document.getElementById("yearFilter").value;
        //year = activeFilters.year
         //activeFilters.semester = semester

    }

    if (semester === "All Semesters" || semesterFilter.value === "") {
        loadOfferingsAndUpdateTable(year,"");
    } else {
        loadOfferingsAndUpdateTable(year ,semester);
    }
}

//TODO need code to acquire the stored years
//based on SELECT distinct year FROM examMetrics
function setupYearFilter() {
    document.getElementById("yearFilter").addEventListener("change", () => {
        filterByYear();
    });
}

function filterByYear(year) {
    const yearFilter = document.getElementById("yearFilter").value;
 
    if (year) {
        // Loop through `buildingFilter` options to select the one with matching text
        for (const option of yearFilter.options) {
            if (option.text === semester) {
                option.selected = true;
                break;
            }
        }
    } else {
        year = yearFilter.selectedOptions[0].text;
        semester = document.getElementById("yearFilter").value;
        //semester = activeFilters.semester
        //activeFilters.year = year
    }


    if (year === "Current Year" || yearFilter.value === "") {
        loadOfferingsAndUpdateTable("",semester);
    } else {
        loadOfferingsAndUpdateTable(year,semester);
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

async function getAllExamOfferings(year="",semester = "") {
    try {
        let url = "/exammetrics";
        const params = new URLSearchParams();
        if (semester) params.append("semester", semester);
        if (year) params.append("year", year);
        if (params.toString()) url += `?${params.toString()}`;
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const offerings = await response.json();
        //console.error(offerings)
        return offerings; // Return the devices instead of storing in global variable
    } catch (err) {
        console.error("Failed to fetch offerings:", err);
        return []; // Return empty array in case of error
    }
}

export async function filteredTable() {
    console.log("> Y: "+activeFilters.year+" S: "+activeFilters.semester)
    getAllExamOfferings(activeFilters.year,activeFilters.semester)
}

async function loadOfferingsAndUpdateTable(year="",semester="") {
    
    const offerings = await getAllExamOfferings(year,semester);
    allOfferings = offerings; // Update global variable if needed
    filteredOfferings = offerings; // Initialize filtered offerings

    updateTable();

    if (offerings.length === 0) {
        const tbody = document.getElementById("exam-offering-body");
        if (tbody) {
            tbody.innerHTML = `<tr><td colspan="12" class="text-center">No exam offerings found.</td></tr>`;
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

    // Clear table if no exam offerings
    if (!Array.isArray(pageOfferings) || pageOfferings.length === 0) {
        tbody.innerHTML = `<tr><td colspan="12" class="text-center">No exam offerings found.</td></tr>`;
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
            (offering) => offering.year === selectedYear
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
            (offering) => offering.semester === selectedSemester
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
        activeFilters.year !== "Current Year" 
    ) {
        filteredOfferings = filteredOfferings.filter(
            (offering) => offering.year === activeFilters.year
        );
    }

    // Apply semester filter if active
    if (
        activeFilters.semester &&
        activeFilters.semester !== "All Semesters" 
    ) {
        filteredOfferings = filteredOfferings.filter(
            //TODO fix the filter application
            (offering) =>
                offering.semester === activeFilters.semester
        );
    }
    updateTable();
}



// Initial fetch without filtering
loadOfferingsAndUpdateTable(activeFilters.year,activeFilters.semester);

//TODO update to show the offerings
function formatOfferingRow(offering) {
    if (!offering) return "";

    //disable CRUD controls
    //const buttons = getActionButtons();
    const buttons = "";
    // Declare isAdmin within the function
    let isAdmin = false;

    // Ensure role is defined and check for "Admin"
    if (role === "Admin") {
        isAdmin = true;
    }

    //CourseCode,Description, Password, ExamID, Year, Semester,Ready, Active, Expired, Closed

    return `
        <tr>
            <td data-label="Course">${
                offering.coursecode
            }</td>
            <td data-label="Description">${
                offering.description
            }</td>
            <td data-label="Exam ID">${offering.examid}</td>
            <td data-label="Password">${offering.password}</td>

            <td data-label="Ready">${offering.ready}</td>
            <td data-label="Active">${offering.active}</td>
            <td data-label="Expired">${offering.expired}</td>
            <td data-label="Closed">${offering.closed}</td>

            <td>
                <div class="btn-group">
                    ${buttons}
                </div>
            </td>
        </tr>
    `;
}

export function getActionButtons() {
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
window.filteredTable = filteredTable;
window.getAllExamOfferings = getAllExamOfferings;