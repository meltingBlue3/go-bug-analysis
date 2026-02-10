(function () {
    'use strict';

    // ===========================
    // DOM References
    // ===========================
    var uploadSection = document.getElementById('upload-section');
    var dashboard = document.getElementById('dashboard');
    var csvFileInput = document.getElementById('csv-file');
    var uploadStatus = document.getElementById('upload-status');
    var uploadBtn = null; // resolved after DOM ready

    // ===========================
    // Constants
    // ===========================
    var MAX_FILE_SIZE = 100 * 1024 * 1024; // 100MB

    // ===========================
    // Status Display
    // ===========================
    function showStatus(message, type) {
        if (!uploadStatus) return;
        type = type || 'info';
        uploadStatus.innerHTML = '<div class="status-' + type + '">' + escapeHtml(message) + '</div>';
    }

    function showStatusHTML(html) {
        if (!uploadStatus) return;
        uploadStatus.innerHTML = html;
    }

    function clearStatus() {
        if (uploadStatus) {
            uploadStatus.innerHTML = '';
        }
    }

    // ===========================
    // HTML Escape Utility
    // ===========================
    function escapeHtml(text) {
        var div = document.createElement('div');
        div.appendChild(document.createTextNode(text));
        return div.innerHTML;
    }

    // ===========================
    // File Size Formatting
    // ===========================
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        var units = ['B', 'KB', 'MB', 'GB'];
        var i = Math.floor(Math.log(bytes) / Math.log(1024));
        return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i];
    }

    // ===========================
    // Dashboard Toggle
    // ===========================
    function showDashboard() {
        if (uploadSection) uploadSection.style.display = 'none';
        if (dashboard) dashboard.style.display = 'block';
    }

    function showUpload() {
        if (uploadSection) uploadSection.style.display = '';
        if (dashboard) dashboard.style.display = 'none';
    }

    // ===========================
    // Upload Button State
    // ===========================
    function setUploading(isUploading) {
        if (uploadBtn) {
            if (isUploading) {
                uploadBtn.classList.add('uploading');
                uploadBtn.style.pointerEvents = 'none';
                uploadBtn.style.opacity = '0.6';
            } else {
                uploadBtn.classList.remove('uploading');
                uploadBtn.style.pointerEvents = '';
                uploadBtn.style.opacity = '';
            }
        }
        if (csvFileInput) {
            csvFileInput.disabled = isUploading;
        }
    }

    // ===========================
    // CSV Upload Handler
    // ===========================
    function uploadCSV(file) {
        // Client-side file size check
        if (file.size > MAX_FILE_SIZE) {
            showStatus('文件过大（最大 100MB），当前文件 ' + formatFileSize(file.size), 'error');
            return;
        }

        // Show uploading state
        setUploading(true);
        showStatus('正在解析 CSV 文件...', 'info');

        var formData = new FormData();
        formData.append('file', file);

        fetch('/api/upload', {
            method: 'POST',
            body: formData
        })
        .then(function (response) {
            return response.json().then(function (data) {
                return { status: response.status, data: data };
            });
        })
        .then(function (result) {
            setUploading(false);
            var data = result.data;

            if (data.success) {
                // Build success message with summary
                var summary = data.summary;
                var html = '<div class="status-success">'
                    + escapeHtml('解析成功！共 ' + summary.totalBugs + ' 条 Bug 记录')
                    + '</div>';

                // Show sample bug info
                if (summary.sampleBug) {
                    html += '<div class="status-info" style="margin-top: 8px; font-size: 0.85em;">'
                        + '示例: #' + escapeHtml(summary.sampleBug.id)
                        + ' - ' + escapeHtml(summary.sampleBug.title)
                        + ' [' + escapeHtml(summary.sampleBug.status) + ']'
                        + '</div>';
                }

                // Show warnings if any
                if (summary.warnings && summary.warnings.length > 0) {
                    var warnCount = Math.min(summary.warnings.length, 5);
                    html += '<div class="status-warning" style="margin-top: 8px; font-size: 0.85em;">';
                    html += '⚠ ' + summary.warnings.length + ' 条警告：<br>';
                    for (var i = 0; i < warnCount; i++) {
                        html += escapeHtml(summary.warnings[i]) + '<br>';
                    }
                    if (summary.warnings.length > 5) {
                        html += '... 还有 ' + (summary.warnings.length - 5) + ' 条';
                    }
                    html += '</div>';
                }

                showStatusHTML(html);

                // Store data for other modules
                window.BugAnalysis.data = data.summary;

                // Switch to dashboard view after a brief delay
                setTimeout(function () {
                    showDashboard();
                    if (window.BugAnalysis.renderDashboard) {
                        window.BugAnalysis.renderDashboard();
                    }
                }, 1500);
            } else {
                showStatus(data.error || '上传失败，请重试', 'error');
            }
        })
        .catch(function (err) {
            setUploading(false);
            showStatus('网络错误，请检查服务是否运行', 'error');
        });
    }

    // ===========================
    // File Input Listener
    // ===========================
    function initFileInput() {
        if (!csvFileInput) return;

        // Resolve the upload button (label wrapping the input)
        uploadBtn = csvFileInput.closest('.upload-btn') || csvFileInput.parentElement;

        csvFileInput.addEventListener('change', function (e) {
            var file = e.target.files && e.target.files[0];
            if (!file) {
                clearStatus();
                return;
            }

            // Check file extension
            var name = file.name.toLowerCase();
            if (name.indexOf('.csv') !== name.length - 4) {
                showStatus('请选择 .csv 格式的文件', 'error');
                csvFileInput.value = '';
                return;
            }

            uploadCSV(file);
        });
    }

    // ===========================
    // Reimport Button
    // ===========================
    function initReimport() {
        var btn = document.getElementById('btn-reimport');
        if (btn) {
            btn.addEventListener('click', function () {
                showUpload();
                clearStatus();
                if (csvFileInput) {
                    csvFileInput.value = '';
                }
            });
        }
    }

    // ===========================
    // Initialization
    // ===========================
    function init() {
        initFileInput();
        initReimport();
    }

    // Wait for DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // ===========================
    // Public API
    // ===========================
    window.BugAnalysis = {
        showStatus: showStatus,
        clearStatus: clearStatus,
        showDashboard: showDashboard,
        showUpload: showUpload,
        data: null
    };

})();
