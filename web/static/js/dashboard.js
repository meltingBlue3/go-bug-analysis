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
    // Public API
    // ===========================
    window.BugAnalysis.renderDashboard = renderDashboard;

})();
