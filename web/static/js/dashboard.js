(function () {
    'use strict';

    // ===========================
    // Module State
    // ===========================
    var severityChart = null;
    var severityData = null;
    var currentChartType = 'pie';
    var currentScope = 'all';
    var fixTimeDistChart = null;
    var workloadActiveChart = null;
    var workloadTotalChart = null;
    var moduleHeatmapChart = null;
    var moduleTrendChart = null;
    var moduleTrendData = null;
    var currentTrendRange = '30';
    var reportData = null;
    var currentReportFormat = 'plain';
    var reportControlsInitialized = false;

    // Color scheme for severity levels
    var severityColors = {
        '1': '#ff4d4f', // 致命 - red
        '2': '#faad14', // 严重 - orange
        '3': '#1890ff', // 一般 - blue
        '4': '#52c41a', // 轻微 - green
        '0': '#999999'  // 未设置 - gray
    };

    // ===========================
    // Dashboard Entry Point
    // ===========================
    function renderDashboard() {
        fetch('/api/analysis')
            .then(function (response) {
                if (!response.ok) {
                    return response.json().then(function (data) {
                        throw new Error(data.error || '分析请求失败');
                    });
                }
                return response.json();
            })
            .then(function (data) {
                renderKPI(data.kpi);
                severityData = data.severity;
                renderSeverityChart(data.severity);
                initToggles();
                renderAge(data.age);
                renderWorkload(data.workload);
                renderModule(data.module);
                reportData = data.report;
                renderReport(data.report);
            })
            .catch(function (err) {
                console.error('Dashboard error:', err);
            });
    }

    // ===========================
    // KPI Card Rendering
    // ===========================
    function renderKPI(kpi) {
        setKPICard('kpi-today-new', kpi.todayNew, '昨日: ' + kpi.yesterdayNew);
        setKPICard('kpi-today-fixed', kpi.todayFixed, '昨日: ' + kpi.yesterdayFixed);
        setKPICard('kpi-total', kpi.total, null);
        setKPICard('kpi-active', kpi.active, null);
        setKPICard('kpi-pending', kpi.pendingVerify, null);

        // Apply warning/danger classes
        var activeCard = document.getElementById('kpi-active');
        if (activeCard) {
            activeCard.classList.remove('placeholder-card');
            if (kpi.active > 0) {
                activeCard.classList.add('kpi-danger');
            }
        }

        var pendingCard = document.getElementById('kpi-pending');
        if (pendingCard) {
            pendingCard.classList.remove('placeholder-card');
            if (kpi.pendingVerify > 0) {
                pendingCard.classList.add('kpi-warning');
            }
        }
    }

    function setKPICard(id, value, subText) {
        var card = document.getElementById(id);
        if (!card) return;

        card.classList.remove('placeholder-card');

        var valueEl = card.querySelector('.kpi-value');
        if (valueEl) {
            valueEl.textContent = value;
        }

        var subEl = card.querySelector('.kpi-sub');
        if (subEl) {
            if (subText) {
                subEl.textContent = subText;
                subEl.style.display = '';
            } else {
                subEl.style.display = 'none';
            }
        }
    }

    // ===========================
    // Severity Chart
    // ===========================
    function renderSeverityChart(severity) {
        var container = document.getElementById('chart-severity');
        if (!container) return;

        // Clear placeholder content
        container.innerHTML = '';

        if (severityChart) {
            severityChart.dispose();
        }
        severityChart = echarts.init(container);

        updateChart();

        // Handle window resize
        window.addEventListener('resize', function () {
            if (severityChart) {
                severityChart.resize();
            }
        });
    }

    function updateChart() {
        if (!severityChart || !severityData) return;

        var data = currentScope === 'all' ? severityData.all : severityData.newOnly;
        var option;

        if (currentChartType === 'pie') {
            option = buildPieOption(data);
        } else {
            option = buildBarOption(data);
        }

        severityChart.setOption(option, true);
    }

    function buildPieOption(data) {
        var seriesData = [];
        for (var i = 0; i < data.length; i++) {
            seriesData.push({
                name: data[i].label,
                value: data[i].count,
                itemStyle: {
                    color: severityColors[data[i].level] || '#999999'
                }
            });
        }

        return {
            tooltip: {
                trigger: 'item',
                formatter: '{b}: {c} ({d}%)'
            },
            legend: {
                orient: 'vertical',
                right: '5%',
                top: 'center',
                textStyle: { fontSize: 13 }
            },
            series: [{
                type: 'pie',
                radius: '55%',
                center: ['40%', '50%'],
                roseType: 'radius',
                data: seriesData,
                label: {
                    formatter: '{b}: {c}',
                    fontSize: 12
                },
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.2)'
                    }
                },
                animationType: 'scale',
                animationEasing: 'elasticOut'
            }]
        };
    }

    function buildBarOption(data) {
        var labels = [];
        var values = [];
        var colors = [];
        for (var i = 0; i < data.length; i++) {
            labels.push(data[i].label);
            values.push(data[i].count);
            colors.push(severityColors[data[i].level] || '#999999');
        }

        return {
            tooltip: {
                trigger: 'axis',
                axisPointer: { type: 'shadow' }
            },
            grid: {
                left: '3%',
                right: '4%',
                bottom: '3%',
                containLabel: true
            },
            xAxis: {
                type: 'category',
                data: labels,
                axisTick: { alignWithLabel: true },
                axisLabel: { fontSize: 13 }
            },
            yAxis: {
                type: 'value',
                minInterval: 1
            },
            series: [{
                type: 'bar',
                data: values.map(function (v, idx) {
                    return {
                        value: v,
                        itemStyle: {
                            color: colors[idx],
                            borderRadius: [4, 4, 0, 0]
                        }
                    };
                }),
                barWidth: '50%',
                label: {
                    show: true,
                    position: 'top',
                    fontSize: 12
                }
            }]
        };
    }

    // ===========================
    // Toggle Controls
    // ===========================
    function initToggles() {
        // Chart type toggle (pie / bar)
        var typeToggle = document.getElementById('severity-type-toggle');
        if (typeToggle) {
            typeToggle.addEventListener('click', function (e) {
                var btn = e.target.closest('.toggle-btn');
                if (!btn || btn.classList.contains('active')) return;

                var siblings = typeToggle.querySelectorAll('.toggle-btn');
                for (var i = 0; i < siblings.length; i++) {
                    siblings[i].classList.remove('active');
                }
                btn.classList.add('active');

                currentChartType = btn.getAttribute('data-value');
                updateChart();
            });
        }

        // Data scope toggle (all / new-only)
        var scopeToggle = document.getElementById('severity-scope-toggle');
        if (scopeToggle) {
            scopeToggle.addEventListener('click', function (e) {
                var btn = e.target.closest('.toggle-btn');
                if (!btn || btn.classList.contains('active')) return;

                var siblings = scopeToggle.querySelectorAll('.toggle-btn');
                for (var i = 0; i < siblings.length; i++) {
                    siblings[i].classList.remove('active');
                }
                btn.classList.add('active');

                currentScope = btn.getAttribute('data-value');
                updateChart();
            });
        }
    }

    // ===========================
    // Age Analysis (Fix Time + Backlog)
    // ===========================
    function renderAge(ageData) {
        if (!ageData) return;
        renderFixTimeStats(ageData.fixTime);
        renderBacklogTable(ageData.backlog);
    }

    function renderFixTimeStats(fixTime) {
        var summaryEl = document.getElementById('fix-time-summary');
        var chartEl = document.getElementById('chart-fix-time-dist');
        if (!summaryEl) return;

        if (!fixTime) {
            summaryEl.innerHTML = '<div class="placeholder-text">暂无已解决 Bug 数据</div>';
            if (chartEl) chartEl.innerHTML = '<div class="placeholder-text">暂无数据</div>';
            return;
        }

        // Format value: if < 24h show hours, else show days
        function formatDuration(hours, days) {
            if (hours < 24) {
                return { value: hours.toFixed(1), unit: '小时' };
            }
            return { value: days.toFixed(1), unit: '天' };
        }

        var avg = formatDuration(fixTime.avgHours, fixTime.avgDays);
        var p50 = formatDuration(fixTime.p50Hours, fixTime.p50Days);

        summaryEl.innerHTML =
            '<div class="stat-box">' +
                '<div class="stat-label">平均修复时长</div>' +
                '<div class="stat-value">' + avg.value + '<span class="stat-unit">' + avg.unit + '</span></div>' +
            '</div>' +
            '<div class="stat-box">' +
                '<div class="stat-label">P50 修复时长</div>' +
                '<div class="stat-value">' + p50.value + '<span class="stat-unit">' + p50.unit + '</span></div>' +
            '</div>' +
            '<div class="stat-sub-text">共 ' + fixTime.totalResolved + ' 条已解决记录</div>';

        // Render distribution chart
        if (!chartEl) return;
        chartEl.innerHTML = '';

        if (fixTimeDistChart) {
            fixTimeDistChart.dispose();
        }
        fixTimeDistChart = echarts.init(chartEl);

        var labels = [];
        var values = [];
        for (var i = 0; i < fixTime.distribution.length; i++) {
            labels.push(fixTime.distribution[i].label);
            values.push(fixTime.distribution[i].count);
        }

        var barColors = ['#52c41a', '#1890ff', '#faad14', '#ff4d4f'];

        var option = {
            tooltip: {
                trigger: 'axis',
                axisPointer: { type: 'shadow' }
            },
            grid: {
                left: '3%',
                right: '8%',
                top: '8%',
                bottom: '3%',
                containLabel: true
            },
            xAxis: {
                type: 'value',
                minInterval: 1
            },
            yAxis: {
                type: 'category',
                data: labels.slice().reverse(),
                axisTick: { show: false },
                axisLabel: { fontSize: 13 }
            },
            series: [{
                type: 'bar',
                data: values.slice().reverse().map(function (v, idx) {
                    return {
                        value: v,
                        itemStyle: {
                            color: barColors.slice().reverse()[idx],
                            borderRadius: [0, 4, 4, 0]
                        }
                    };
                }),
                barWidth: '50%',
                label: {
                    show: true,
                    position: 'right',
                    fontSize: 12,
                    color: '#666'
                }
            }]
        };

        fixTimeDistChart.setOption(option);

        window.addEventListener('resize', function () {
            if (fixTimeDistChart) {
                fixTimeDistChart.resize();
            }
        });
    }

    function renderBacklogTable(backlog) {
        var countEl = document.getElementById('backlog-count');
        var tbody = document.getElementById('backlog-tbody');
        if (!tbody) return;

        if (countEl) {
            countEl.textContent = backlog ? '共 ' + backlog.length + ' 条' : '';
        }

        if (!backlog || backlog.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="placeholder-text">暂无未解决 Bug</td></tr>';
            return;
        }

        var html = '';
        for (var i = 0; i < backlog.length; i++) {
            var item = backlog[i];
            var title = item.title;
            if (title.length > 40) {
                title = title.substring(0, 40) + '…';
            }

            // Severity badge
            var sevClass = 's' + item.severity;
            var sevLabels = { '1': '致命', '2': '严重', '3': '一般', '4': '轻微' };
            var sevLabel = sevLabels[item.severity] || item.severity;

            // Age color coding
            var ageClass = '';
            if (item.ageDays > 30) {
                ageClass = 'age-danger';
            } else if (item.ageDays > 14) {
                ageClass = 'age-warning';
            }

            // Date: show first 10 chars
            var dateStr = item.createdDate ? item.createdDate.substring(0, 10) : '';

            html += '<tr>' +
                '<td>' + item.id + '</td>' +
                '<td title="' + item.title.replace(/"/g, '&quot;') + '">' + title + '</td>' +
                '<td><span class="severity-badge ' + sevClass + '">' + sevLabel + '</span></td>' +
                '<td>' + (item.assignee || '—') + '</td>' +
                '<td>' + dateStr + '</td>' +
                '<td class="' + ageClass + '">' + item.ageDays + ' 天</td>' +
                '</tr>';
        }

        tbody.innerHTML = html;
    }

    // ===========================
    // Workload Distribution Charts
    // ===========================
    function renderWorkload(workload) {
        if (!workload) return;

        workloadActiveChart = renderWorkloadChart(
            'chart-workload-active',
            workload.byActive,
            '#ff7a45', '#ff4d4f',
            workloadActiveChart
        );
        workloadTotalChart = renderWorkloadChart(
            'chart-workload-total',
            workload.byTotal,
            '#1890ff', '#096dd9',
            workloadTotalChart
        );
    }

    function renderWorkloadChart(containerId, data, colorStart, colorEnd, existingChart) {
        var container = document.getElementById(containerId);
        if (!container) return null;

        container.innerHTML = '';

        if (!data || data.length === 0) {
            container.innerHTML = '<div class="placeholder-text">暂无数据</div>';
            return null;
        }

        // Limit to top 15 assignees
        var maxItems = 15;
        var displayData = data.length > maxItems ? data.slice(0, maxItems) : data;

        // Set dynamic height before init
        var chartHeight = Math.max(300, displayData.length * 32);
        container.style.height = chartHeight + 'px';

        if (existingChart) {
            existingChart.dispose();
        }
        var chart = echarts.init(container);

        // Reverse for ECharts y-axis (highest count at top)
        var names = [];
        var values = [];
        for (var i = displayData.length - 1; i >= 0; i--) {
            names.push(displayData[i].name);
            values.push(displayData[i].count);
        }

        var option = {
            tooltip: {
                trigger: 'axis',
                axisPointer: { type: 'shadow' }
            },
            grid: {
                left: '100',
                right: '60',
                top: '8%',
                bottom: '3%',
                containLabel: false
            },
            xAxis: {
                type: 'value',
                minInterval: 1
            },
            yAxis: {
                type: 'category',
                data: names,
                axisTick: { show: false },
                axisLabel: {
                    fontSize: 12,
                    width: 80,
                    overflow: 'truncate',
                    ellipsis: '...'
                }
            },
            series: [{
                type: 'bar',
                data: values.map(function (v) {
                    return {
                        value: v,
                        itemStyle: {
                            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
                                { offset: 0, color: colorStart },
                                { offset: 1, color: colorEnd }
                            ]),
                            borderRadius: [0, 4, 4, 0]
                        }
                    };
                }),
                barWidth: '60%',
                label: {
                    show: true,
                    position: 'right',
                    fontSize: 12,
                    color: '#666'
                }
            }]
        };

        // Add subtitle if data was truncated
        if (data.length > maxItems) {
            option.title = {
                text: '显示前' + maxItems + '人（共' + data.length + '人）',
                textStyle: { fontSize: 12, color: '#999', fontWeight: 'normal' },
                right: '10',
                top: '0'
            };
        }

        chart.setOption(option);

        window.addEventListener('resize', function () {
            if (chart) {
                chart.resize();
            }
        });

        return chart;
    }

    // ===========================
    // Module Analysis
    // ===========================
    function renderModule(moduleData) {
        if (!moduleData) return;
        renderModuleTable(moduleData.stats);
        renderModuleHeatmap(moduleData.heatmap);
        moduleTrendData = moduleData.trend;
        renderModuleTrend(moduleData.trend);
        initTrendToggle();
    }

    function renderModuleTable(stats) {
        var countEl = document.getElementById('module-count');
        var tbody = document.getElementById('module-tbody');
        if (!tbody) return;

        if (countEl) {
            countEl.textContent = stats ? '共 ' + stats.length + ' 个模块' : '';
        }

        if (!stats || stats.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="placeholder-text">暂无模块数据</td></tr>';
            return;
        }

        // Clone stats for sorting
        var sortedStats = stats.slice();
        var currentSortKey = 'total';
        var currentSortDir = 'desc';

        function buildRows(data) {
            var html = '';
            for (var i = 0; i < data.length; i++) {
                var s = data[i];
                var activeClass = s.active > 10 ? ' class="age-danger"' : '';
                var rateClass = '';
                if (s.activeRate > 50) {
                    rateClass = ' class="rate-danger"';
                } else if (s.activeRate > 30) {
                    rateClass = ' class="rate-warning"';
                }
                html += '<tr>' +
                    '<td>' + s.name + '</td>' +
                    '<td>' + s.total + '</td>' +
                    '<td' + activeClass + '>' + s.active + '</td>' +
                    '<td' + rateClass + '>' + s.activeRate.toFixed(1) + '%</td>' +
                    '</tr>';
            }
            return html;
        }

        tbody.innerHTML = buildRows(sortedStats);

        // Click-to-sort on headers
        var table = document.getElementById('module-table');
        if (!table) return;
        var headers = table.querySelectorAll('th[data-sort]');
        for (var h = 0; h < headers.length; h++) {
            headers[h].addEventListener('click', (function (th) {
                return function () {
                    var key = th.getAttribute('data-sort');
                    if (currentSortKey === key) {
                        currentSortDir = currentSortDir === 'desc' ? 'asc' : 'desc';
                    } else {
                        currentSortKey = key;
                        currentSortDir = 'desc';
                    }

                    // Update header indicators
                    for (var j = 0; j < headers.length; j++) {
                        headers[j].classList.remove('sort-active', 'sort-desc');
                        var text = headers[j].textContent.replace(/ [▲▼]$/, '');
                        headers[j].textContent = text;
                    }
                    th.classList.add('sort-active');
                    if (currentSortDir === 'desc') th.classList.add('sort-desc');
                    th.textContent = th.textContent + (currentSortDir === 'desc' ? ' ▼' : ' ▲');

                    // Sort
                    sortedStats = stats.slice();
                    sortedStats.sort(function (a, b) {
                        var va, vb;
                        if (key === 'name') { va = a.name; vb = b.name; }
                        else if (key === 'total') { va = a.total; vb = b.total; }
                        else if (key === 'active') { va = a.active; vb = b.active; }
                        else if (key === 'rate') { va = a.activeRate; vb = b.activeRate; }
                        if (key === 'name') {
                            return currentSortDir === 'asc' ? va.localeCompare(vb) : vb.localeCompare(va);
                        }
                        return currentSortDir === 'asc' ? va - vb : vb - va;
                    });

                    tbody.innerHTML = buildRows(sortedStats);
                };
            })(headers[h]));
        }
    }

    function renderModuleHeatmap(heatmap) {
        var container = document.getElementById('chart-module-heatmap');
        if (!container) return;
        if (!heatmap || !heatmap.modules || !heatmap.modules.length) {
            container.innerHTML = '<div class="placeholder-text">暂无模块数据</div>';
            return;
        }

        container.innerHTML = '';
        container.style.height = Math.max(300, heatmap.modules.length * 36) + 'px';

        if (moduleHeatmapChart) {
            moduleHeatmapChart.dispose();
        }
        moduleHeatmapChart = echarts.init(container);

        // Reverse modules for top-down display (highest count at top)
        var reversedModules = heatmap.modules.slice().reverse();

        // Flatten 2D data to [xIdx, yIdx, value] tuples with reversed yIdx
        var seriesData = [];
        for (var mi = 0; mi < heatmap.modules.length; mi++) {
            for (var si = 0; si < heatmap.severities.length; si++) {
                var val = heatmap.data[mi][si];
                var reversedY = heatmap.modules.length - 1 - mi;
                seriesData.push([si, reversedY, val]);
            }
        }

        var option = {
            tooltip: {
                position: 'top',
                formatter: function (params) {
                    return '模块: ' + reversedModules[params.value[1]] +
                        ', ' + heatmap.severities[params.value[0]] +
                        ': ' + params.value[2] + '条';
                }
            },
            grid: {
                left: 120,
                right: 80,
                top: 30,
                bottom: 40
            },
            xAxis: {
                type: 'category',
                data: heatmap.severities,
                splitArea: { show: true }
            },
            yAxis: {
                type: 'category',
                data: reversedModules,
                splitArea: { show: true },
                axisLabel: {
                    fontSize: 12,
                    width: 100,
                    overflow: 'truncate',
                    ellipsis: '...'
                }
            },
            visualMap: {
                min: 0,
                max: heatmap.maxValue || 1,
                calculable: true,
                orient: 'vertical',
                right: 10,
                top: 'center',
                inRange: {
                    color: ['#f0f5ff', '#1890ff', '#003a8c']
                }
            },
            series: [{
                type: 'heatmap',
                data: seriesData,
                label: {
                    show: true,
                    formatter: function (params) {
                        return params.value[2] > 0 ? params.value[2] : '';
                    }
                },
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }]
        };

        moduleHeatmapChart.setOption(option);

        window.addEventListener('resize', function () {
            if (moduleHeatmapChart) {
                moduleHeatmapChart.resize();
            }
        });
    }

    var trendColors = ['#1890ff', '#ff4d4f', '#faad14', '#52c41a', '#722ed1',
                       '#eb2f96', '#13c2c2', '#fa8c16', '#2f54eb', '#a0d911'];

    function renderModuleTrend(trend) {
        var container = document.getElementById('chart-module-trend');
        if (!container) return;
        if (!trend || !trend.series || !trend.series.length) {
            container.innerHTML = '<div class="placeholder-text">暂无趋势数据</div>';
            return;
        }

        container.innerHTML = '';

        if (moduleTrendChart) {
            moduleTrendChart.dispose();
        }
        moduleTrendChart = echarts.init(container);

        // Determine date slice based on currentTrendRange
        var dates, slicedSeries;
        if (currentTrendRange === '7') {
            dates = trend.dates.slice(trend.days7);
            slicedSeries = [];
            for (var i = 0; i < trend.series.length; i++) {
                slicedSeries.push({
                    name: trend.series[i].name,
                    counts: trend.series[i].counts.slice(trend.days7)
                });
            }
        } else {
            dates = trend.dates;
            slicedSeries = trend.series;
        }

        var seriesOpt = [];
        for (var s = 0; s < slicedSeries.length; s++) {
            seriesOpt.push({
                type: 'line',
                name: slicedSeries[s].name,
                data: slicedSeries[s].counts,
                smooth: true,
                symbol: 'circle',
                symbolSize: 4,
                lineStyle: { width: 2 },
                color: trendColors[s % trendColors.length]
            });
        }

        var option = {
            tooltip: {
                trigger: 'axis'
            },
            legend: {
                bottom: 0,
                type: 'scroll'
            },
            grid: {
                left: '3%',
                right: '4%',
                bottom: '15%',
                top: '8%',
                containLabel: true
            },
            xAxis: {
                type: 'category',
                data: dates,
                axisLabel: {
                    rotate: dates.length > 15 ? 30 : 0
                }
            },
            yAxis: {
                type: 'value',
                minInterval: 1
            },
            series: seriesOpt
        };

        moduleTrendChart.setOption(option, true);

        window.addEventListener('resize', function () {
            if (moduleTrendChart) {
                moduleTrendChart.resize();
            }
        });
    }

    function initTrendToggle() {
        var toggleGroup = document.getElementById('trend-range-toggle');
        if (!toggleGroup) return;
        toggleGroup.addEventListener('click', function (e) {
            var btn = e.target.closest('.toggle-btn');
            if (!btn || btn.classList.contains('active')) return;

            var siblings = toggleGroup.querySelectorAll('.toggle-btn');
            for (var i = 0; i < siblings.length; i++) {
                siblings[i].classList.remove('active');
            }
            btn.classList.add('active');

            currentTrendRange = btn.getAttribute('data-value');
            renderModuleTrend(moduleTrendData);
        });
    }

    // ===========================
    // Report Section
    // ===========================
    function renderReport(report) {
        var contentEl = document.getElementById('report-content');
        if (!contentEl) return;

        if (!report) {
            contentEl.textContent = '暂无日报数据';
            return;
        }

        if (currentReportFormat === 'plain') {
            contentEl.textContent = report.plainText;
        } else {
            contentEl.textContent = report.markdown;
        }

        if (!reportControlsInitialized) {
            initReportControls();
            reportControlsInitialized = true;
        }
    }

    function initReportControls() {
        // Format toggle
        var toggleGroup = document.getElementById('report-format-toggle');
        if (toggleGroup) {
            toggleGroup.addEventListener('click', function (e) {
                var btn = e.target.closest('.toggle-btn');
                if (!btn || btn.classList.contains('active')) return;

                var siblings = toggleGroup.querySelectorAll('.toggle-btn');
                for (var i = 0; i < siblings.length; i++) {
                    siblings[i].classList.remove('active');
                }
                btn.classList.add('active');

                currentReportFormat = btn.getAttribute('data-value');
                renderReport(reportData);
            });
        }

        // Copy button
        var copyBtn = document.getElementById('btn-copy-report');
        if (copyBtn) {
            copyBtn.addEventListener('click', function () {
                if (!reportData) return;
                var text = currentReportFormat === 'plain' ? reportData.plainText : reportData.markdown;
                navigator.clipboard.writeText(text).then(function () {
                    var toast = document.getElementById('copy-toast');
                    if (toast) {
                        toast.classList.add('show');
                        setTimeout(function () { toast.classList.remove('show'); }, 2000);
                    }
                }).catch(function (err) {
                    console.error('复制失败:', err);
                });
            });
        }
    }

    // ===========================
    // Public API
    // ===========================
    window.BugAnalysis.renderDashboard = renderDashboard;

})();
