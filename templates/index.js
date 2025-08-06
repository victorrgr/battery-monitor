let batteryChart = null;
let currentPage = 0;
let totalPages = 0;
let isLoading = false;

const pageSize = 10;

// ==============================
// UI Feedback
// ==============================

function showLoading() {
    document.getElementById('loadingSpinner').style.display = 'block';
    isLoading = true;
}

function hideLoading() {
    document.getElementById('loadingSpinner').style.display = 'none';
    isLoading = false;
}

function showError(message) {
    const errorElement = document.getElementById('errorMessage');
    errorElement.textContent = message;
    errorElement.style.display = 'block';

    setTimeout(() => {
        errorElement.style.display = 'none';
    }, 5000);
}

// ==============================
// Pagination Display
// ==============================

function updatePageInfo() {
    const pageInfo = document.getElementById('pageInfo');
    pageInfo.textContent = totalPages === 0
        ? 'No data'
        : `Page ${currentPage + 1} of ${totalPages}`;
}

function updatePaginationButtons() {
    const firstBtn = document.getElementById('firstPage');
    const prevBtn = document.getElementById('prevPage');
    const nextBtn = document.getElementById('nextPage');
    const lastBtn = document.getElementById('lastPage');

    const atFirstPage = currentPage === 0;
    const atLastPage = currentPage >= totalPages - 1;
    const noPages = totalPages === 0;

    firstBtn.disabled = atFirstPage || noPages;
    prevBtn.disabled = atFirstPage || noPages;
    nextBtn.disabled = atLastPage || noPages;
    lastBtn.disabled = atLastPage || noPages;
}

// ==============================
// Chart Configuration & Rendering
// ==============================

function chartConfig(labels, dataPoints) {
    return {
        type: 'line',
        data: {
            labels,
            datasets: [{
                label: 'Battery Level (%)',
                data: dataPoints,
                borderColor: '#667eea',
                backgroundColor: 'rgba(102, 126, 234, 0.1)',
                fill: true,
                tension: 0.4,
                pointRadius: 3,
                pointHoverRadius: 6,
                pointBackgroundColor: '#667eea',
                pointBorderColor: '#ffffff',
                pointBorderWidth: 2,
                borderWidth: 3,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: {
                duration: 1000,
                easing: 'easeInOutCubic'
            },
            plugins: {
                legend: {
                    display: true,
                    position: 'top',
                    labels: {
                        font: { size: 14, weight: '500' },
                        color: '#4a5568'
                    }
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    backgroundColor: 'rgba(255, 255, 255, 0.95)',
                    titleColor: '#2d3748',
                    bodyColor: '#4a5568',
                    borderColor: '#e2e8f0',
                    borderWidth: 1,
                    cornerRadius: 8,
                    displayColors: true,
                    titleFont: { weight: '600' },
                    bodyFont: { weight: '500' }
                }
            },
            interaction: {
                mode: 'nearest',
                intersect: false
            },
            scales: {
                y: {
                    min: 0,
                    max: 100,
                    title: {
                        display: true,
                        text: 'Battery Level (%)',
                        font: { size: 14, weight: '500' },
                        color: '#4a5568'
                    },
                    grid: { color: 'rgba(74, 85, 104, 0.1)' },
                    ticks: {
                        font: { size: 12 },
                        color: '#718096'
                    }
                },
                x: {
                    title: {
                        display: true,
                        text: 'Time',
                        font: { size: 14, weight: '500' },
                        color: '#4a5568'
                    },
                    grid: { color: 'rgba(74, 85, 104, 0.1)' },
                    ticks: {
                        font: { size: 12 },
                        color: '#718096',
                        maxTicksLimit: 10
                    }
                }
            }
        }
    };
}

function renderChart(logs) {
    if (!logs || logs.length === 0) {
        showError('No battery data available for the selected date.');
        return;
    }

    const labels = logs.map(log => new Date(log.timestamp).toLocaleTimeString([], {
        hour: '2-digit', minute: '2-digit', hour12: false
    }));

    const dataPoints = logs.map(log => log.percent);

    const ctx = document.getElementById('batteryChart').getContext('2d');

    if (batteryChart) batteryChart.destroy();

    batteryChart = new Chart(ctx, chartConfig(labels, dataPoints));
}

// ==============================
// Fetch Data
// ==============================

async function fetchAndRenderData(date) {
    if (isLoading) return;

    try {
        showLoading();
        const response = await fetch(`/data?date=${date}`);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const data = await response.json();
        renderChart(data);
    } catch (err) {
        console.error("Error fetching chart data:", err);
        showError(`Failed to load battery data: ${err.message}`);
    } finally {
        hideLoading();
    }
}

async function fetchDates() {
    if (isLoading) return;

    try {
        showLoading();
        const response = await fetch(`/dates?page=${currentPage}&size=${pageSize}`);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const result = await response.json();

        // Handle new format: { dates: [...], totalPages, totalItems }
        updateDateSelector(result.dates || []);
        totalPages = result.totalPages || 1;

        updatePageInfo();
        updatePaginationButtons();
    } catch (err) {
        console.error("Error loading dates:", err);
        showError(`Failed to load dates: ${err.message}`);
        totalPages = 0;

        updatePageInfo();
        updatePaginationButtons();

        const selector = document.getElementById('dateSelector');
        selector.innerHTML = '<option value="">No dates available</option>';
    } finally {
        // ✅ End loading first
        hideLoading();

        // ✅ Then trigger chart render if a date is selected
        const selector = document.getElementById('dateSelector');
        if (selector.value) {
            fetchAndRenderData(selector.value);
        }
    }
}

// ==============================
// DOM Handling
// ==============================

function updateDateSelector(dates) {
    const selector = document.getElementById('dateSelector');
    selector.innerHTML = '';

    if (!dates || dates.length === 0) {
        const option = document.createElement('option');
        option.value = '';
        option.textContent = 'No dates available';
        selector.appendChild(option);
        return;
    }

    dates.forEach(date => {
        const option = document.createElement('option');
        option.value = date;

        const formattedDate = new Date(date).toLocaleDateString([], {
            weekday: 'short', year: 'numeric', month: 'short', day: 'numeric'
        });

        option.textContent = formattedDate;
        selector.appendChild(option);
    });

    // Select most recent date
    selector.value = dates[dates.length - 1];
}

// ==============================
// Event Setup
// ==============================

function setupPaginationControls() {
    document.getElementById('firstPage').addEventListener('click', () => {
        if (currentPage !== 0) {
            currentPage = 0;
            fetchDates();
        }
    });

    document.getElementById('prevPage').addEventListener('click', () => {
        if (currentPage > 0) {
            currentPage--;
            fetchDates();
        }
    });

    document.getElementById('nextPage').addEventListener('click', () => {
        if (currentPage < totalPages - 1) {
            currentPage++;
            fetchDates();
        }
    });

    document.getElementById('lastPage').addEventListener('click', () => {
        if (currentPage !== totalPages - 1) {
            currentPage = totalPages - 1;
            fetchDates();
        }
    });
}

function setupDateSelectorChange() {
    document.getElementById('dateSelector').addEventListener('change', (e) => {
        if (e.target.value) {
            fetchAndRenderData(e.target.value);
        }
    });
}

// ==============================
// Initialization
// ==============================

window.addEventListener('DOMContentLoaded', () => {
    setupPaginationControls();
    setupDateSelectorChange();
    fetchDates();
});
