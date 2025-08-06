let batteryChart = null;
let currentPage = 0;
const pageSize = 5;

function chartConfig(labels, dataPoints) {
    return {
        type: 'line',
        data: {
            labels,
            datasets: [{
                label: 'Battery %',
                data: dataPoints,
                borderColor: 'rgba(75, 192, 192, 1)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                fill: true,
                tension: 0.4,
                pointRadius: 0,
                pointHoverRadius: 4,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { display: true },
                tooltip: { mode: 'index', intersect: false }
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
                        text: 'Battery %'
                    }
                },
                x: {
                    title: {
                        display: true,
                        text: 'Time'
                    }
                }
            }
        }
    };
}

function renderChart(logs) {
    const labels = logs.map(log => new Date(log.timestamp).toLocaleTimeString());
    const dataPoints = logs.map(log => log.percent);

    const ctx = document.getElementById('batteryChart').getContext('2d');
    if (batteryChart) batteryChart.destroy();
    batteryChart = new Chart(ctx, chartConfig(labels, dataPoints));
}

function fetchAndRenderData(date) {
    fetch(`/data?date=${date}`)
        .then(res => res.ok ? res.json() : [])
        .then(data => renderChart(data))
        .catch(err => console.error("Error fetching chart data:", err));
}

function updateDateSelector(dates) {
    const selector = document.getElementById('dateSelector');
    selector.innerHTML = '';

    dates.forEach(date => {
        const option = document.createElement('option');
        option.value = date;
        option.textContent = date;
        selector.appendChild(option);
    });

    if (dates.length > 0) {
        selector.value = dates[dates.length - 1];
        fetchAndRenderData(selector.value);
    }
}

function fetchDates() {
    fetch(`/dates?page=${currentPage}&size=${pageSize}`)
        .then(res => res.ok ? res.json() : [])
        .then(dates => updateDateSelector(dates))
        .catch(err => console.error("Error loading dates:", err));
}

function setupPaginationControls() {
    document.getElementById('prevPage').addEventListener('click', () => {
        if (currentPage > 0) {
            currentPage--;
            fetchDates();
        }
    });

    document.getElementById('nextPage').addEventListener('click', () => {
        currentPage++;
        fetchDates();
    });
}

function setupDateSelectorChange() {
    document.getElementById('dateSelector').addEventListener('change', (e) => {
        fetchAndRenderData(e.target.value);
    });
}

window.addEventListener('DOMContentLoaded', () => {
    setupPaginationControls();
    setupDateSelectorChange();
    fetchDates();
});
